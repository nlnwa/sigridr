package worker

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
)

type Config struct {
	AccessToken string
}

// result storage format
type SearchResult struct {
	Id         string            `json:"id,omitempty"`
	CreateTime time.Time         `json:"create_time,omitempty"`
	Seq        int32             `json:"seq,omitempty"`
	Ref        string            `json:"ref,omitempty"`
	Metadata   *twitter.Metadata `json:"search_metadata,omitempty"`
	Statuses   *[]twitter.Tweet  `json:"statuses,omitempty"`
	Response   *twitter.Response `json:"response,omitempty"`
}

type worker struct {
	store  *database.Rethink
	client *twitter.Client
}

func NewApi(c Config) api.WorkerServer {
	httpClient := auth.HttpClient(c.AccessToken)
	httpClient.Timeout = 10 * time.Second

	return &worker{database.New(), twitter.New(httpClient)}
}

func (w *worker) Do(context context.Context, work *api.WorkRequest) (*api.WorkReply, error) {
	queuedSeed := work.QueuedSeed

	params := &twitter.Params{
		ResultType: "recent",
		TweetMode:  "extended",
		Count:      100,
		Query:      queuedSeed.Parameter.Query,
	}
	maxId, err := strconv.ParseInt(queuedSeed.Parameter.MaxId, 10, 64)
	if err == nil {
		params.MaxID = maxId
	}
	sinceId, err := strconv.ParseInt(queuedSeed.Parameter.SinceId, 10, 64)
	if err == nil {
		params.SinceID = sinceId
	}

	now := time.Now().UTC()

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

	search := &SearchResult{
		Seq:        queuedSeed.GetSeq(),
		Ref:        queuedSeed.GetRef(),
		CreateTime: now,
		Metadata:   result.Metadata,
		Statuses:   &result.Statuses,
		Response:   response,
	}
	id, err := w.saveSearchResult(search)
	if err != nil {
		log.WithError(err).Errorln("saving search result")
		return nil, err
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

func (w *worker) saveSearchResult(value interface{}) (string, error) {
	err := w.store.Connect(database.DefaultOptions())
	defer w.store.Disconnect()
	if err != nil {
		return "", err
	}
	return w.store.Insert("result", value)
}
