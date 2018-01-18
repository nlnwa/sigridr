package worker

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
	"fmt"
)

type Config struct {
	AccessToken     string
	DatabaseAddress string
	DatabaseName    string
}

type searchResult struct {
	Id         string            `json:"id,omitempty"`
	CreateTime time.Time         `json:"create_time,omitempty"`
	Seq        int32             `json:"seq,omitempty"`
	Ref        string            `json:"ref,omitempty"`
	Metadata   *twitter.Metadata `json:"search_metadata,omitempty"`
	Statuses   *[]twitter.Tweet  `json:"statuses,omitempty"`
	Response   *twitter.Response `json:"response,omitempty"`
	Params     *api.Parameter    `json:"params,omitempty"`
}

type worker struct {
	store  *workerStore
	client twitter.Client
}

func NewApi(c Config) api.WorkerServer {
	httpClient := auth.HttpClient(c.AccessToken)
	httpClient.Timeout = 10 * time.Second

	return &worker{
		store:  newStore(c),
		client: twitter.New(httpClient),
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
		log.WithError(err).Errorln("searching twitter")
		return nil, err
	} else {
		log.WithFields(log.Fields{
			"query":   params.Query,
			"ref":     queuedSeed.GetRef(),
			"seq":     queuedSeed.GetSeq(),
			"maxId":   queuedSeed.GetParameter().GetMaxId(),
			"sinceId": queuedSeed.GetParameter().GetSinceId(),
			"tweets":  len(result.Statuses),
		}).Infoln("search")
	}

	search := &searchResult{
		Seq:        queuedSeed.GetSeq(),
		Ref:        queuedSeed.GetRef(),
		CreateTime: now,
		Metadata:   result.Metadata,
		Statuses:   &result.Statuses,
		Response:   response,
		Params:     params.ToProto(),
	}
	id, err := w.store.saveSearchResult(search)
	if err != nil {
		return nil, fmt.Errorf("saving search result: %v", err)
	}

	if queuedSeed.GetSeq() == 0 {
		queuedSeed.Ref = id
	}

	maxIdStr, err := twitter.MaxId(search.Metadata)
	if err != nil {
		return nil, err
	}
	sinceIdStr, err := twitter.SinceId(search.Metadata)
	if err != nil {
		return nil, err
	}

	return &api.WorkReply{
		QueuedSeed: queuedSeed,
		Count:      int32(len(result.Statuses)),
		MaxId:      maxIdStr,
		SinceId:    sinceIdStr,
		RateLimit:  ratelimit.New().FromHttpHeaders(response.Header).ToProto(),
	}, nil

}
