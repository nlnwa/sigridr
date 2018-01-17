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
