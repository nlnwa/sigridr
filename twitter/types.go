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

package twitter

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type Metadata = twitter.SearchMetadata
type Tweet = twitter.Tweet
type Result = twitter.Search
type Response = types.Response
type Params twitter.SearchTweetParams

func (p *Params) FromProto(parameter *api.Parameter) *Params {
	sinceId, _ := strconv.ParseInt(parameter.SinceId, 10, 64)
	maxId, _ := strconv.ParseInt(parameter.MaxId, 10, 64)

	return &Params{
		Query:           parameter.Query,
		Geocode:         parameter.Geocode,
		Lang:            parameter.Geocode,
		Locale:          parameter.Locale,
		ResultType:      parameter.ResultType,
		Count:           int(parameter.Count),
		SinceID:         sinceId,
		MaxID:           maxId,
		Until:           parameter.Until,
		IncludeEntities: &parameter.IncludeEntities,
		TweetMode:       parameter.TweetMode,
	}
}

func (p *Params) ToProto() *api.Parameter {
	includeEntities := false
	maxId := ""
	sinceId := ""

	if p.IncludeEntities != nil {
		includeEntities = *p.IncludeEntities
	}
	if p.MaxID > 0 {
		maxId = strconv.FormatInt(p.MaxID, 10)
	}
	if p.SinceID > 0 {
		sinceId = strconv.FormatInt(p.SinceID, 10)
	}

	return &api.Parameter{
		Query:           p.Query,
		Geocode:         p.Geocode,
		Lang:            p.Lang,
		Locale:          p.Locale,
		ResultType:      p.ResultType,
		Count:           int32(p.Count),
		MaxId:           maxId,
		SinceId:         sinceId,
		Until:           p.Until,
		IncludeEntities: includeEntities,
		TweetMode:       p.TweetMode,
	}
}
