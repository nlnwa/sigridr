// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"context"
	"strconv"
	"time"

	"github.com/nlnwa/pkg/log"
	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
)

type Config struct {
	AccessToken      string
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	Logger           log.Logger
}

type searchResult struct {
	Id          string            `json:"id,omitempty"`
	CreateTime  time.Time         `json:"create_time,omitempty"`
	Seq         int32             `json:"seq,omitempty"`
	ExecutionId string            `json:"execution_id,omitempty"`
	Metadata    *twitter.Metadata `json:"search_metadata,omitempty"`
	Statuses    *[]twitter.Tweet  `json:"statuses,omitempty"`
	Response    *twitter.Response `json:"response,omitempty"`
	Params      *api.Parameter    `json:"params,omitempty"`
}

type worker struct {
	store  *workerStore
	client twitter.Client
	log.Logger
}

func NewApi(c Config) api.WorkerServer {
	httpClient := auth.HttpClient(c.AccessToken)
	httpClient.Timeout = 10 * time.Second

	return &worker{
		store:  newStore(c),
		client: twitter.New(httpClient),
		Logger: c.Logger,
	}
}

// parameters returns the parameters to be used in the request to twitter by
// mixing default values with the parameters sent via the API.
func parameters(parameter *api.Parameter) *twitter.Params {
	p := &twitter.Params{
		ResultType: "recent",
		TweetMode:  "extended",
		Count:      twitter.MaxStatusesPerRequest,
		Query:      parameter.Query,
	}

	maxId, err := strconv.ParseInt(parameter.GetMaxId(), 10, 64)
	if err == nil {
		p.MaxID = maxId
	}

	sinceId, err := strconv.ParseInt(parameter.GetSinceId(), 10, 64)
	if err == nil {
		p.SinceID = sinceId
	}
	return p
}

func (w *worker) Do(context context.Context, work *api.WorkRequest) (*api.WorkReply, error) {
	queuedSeed := work.QueuedSeed
	now := time.Now().UTC()
	params := parameters(queuedSeed.Parameter)

	result, response, err := w.client.Search(params)
	if err != nil {
		return nil, err
	} else {
		w.Logger.Debug("Search", "query", params.Query)
	}

	// Rate limit
	rl, err := ratelimit.New().FromHttpHeaders(response.Header)
	if err != nil {
		return nil, err
	}
	rateLimit, err := rl.ToProto()
	if err != nil {
		return nil, err
	}
	// number of statuses
	count := int32(len(result.Statuses))

	// don't save result if no statuses
	if count == 0 {
		return &api.WorkReply{
			QueuedSeed: queuedSeed,
			Count:      count,
			RateLimit:  rateLimit,
		}, nil
	} else {
		// save search result
		search := &searchResult{
			Seq:         queuedSeed.GetSeq(),
			CreateTime:  now,
			Metadata:    result.Metadata,
			Statuses:    &result.Statuses,
			ExecutionId: queuedSeed.GetExecutionId(),
			Response:    response,
			Params:      params.ToProto(),
		}
		if _, err = w.store.saveSearchResult(search); err != nil {
			return nil, err
		}
	}
	maxIdStr, err := twitter.MaxId(result.Metadata)
	if err != nil {
		return nil, err
	}
	sinceIdStr, err := twitter.SinceId(result.Metadata)
	if err != nil {
		return nil, err
	}
	return &api.WorkReply{
		QueuedSeed: queuedSeed,
		Count:      count,
		MaxId:      maxIdStr,
		SinceId:    sinceIdStr,
		RateLimit:  rateLimit,
	}, nil

}
