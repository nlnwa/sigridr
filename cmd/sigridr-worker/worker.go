package main

import (
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/nlnwa/sigridr/api/sigridr"
	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/twitter"
	"github.com/nlnwa/sigridr/pkg/twitter/auth"
	"github.com/nlnwa/sigridr/pkg/twitter/ratelimit"
)

var (
	once   sync.Once
	client *twitter.Client
	wg     sync.WaitGroup
)

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

type worker struct{}

func (s *worker) Do(context context.Context, work *pb.WorkRequest) (*pb.WorkReply, error) {
	queuedSeed := work.QueuedSeed

	params := &twitter.Params{
		ResultType: "recent",
		TweetMode:  "extended",
		Count:      100,
		Query:      queuedSeed.Parameters.Query,
	}
	maxId, err := strconv.ParseInt(queuedSeed.Parameters.MaxId, 10, 64)
	if err == nil {
		params.MaxID = maxId
	}
	sinceId, err := strconv.ParseInt(queuedSeed.Parameters.SinceId, 10, 64)
	if err == nil {
		params.SinceID = sinceId
	}

	now := time.Now().UTC()

	result, response, err := twitterClient().Search(params)
	if err != nil {
		log.WithError(err).Errorln("searching twitter")
		return nil, err
	} else {
		log.WithFields(log.Fields{
			"query":   params.Query,
			"ref":     queuedSeed.GetRef(),
			"seq":     queuedSeed.GetSeq(),
			"maxId":   queuedSeed.GetParameters().GetMaxId(),
			"sinceId": queuedSeed.GetParameters().GetSinceId(),
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
	id, err := db.SaveSearchResult(search)
	if err != nil {
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

	return &pb.WorkReply{
		QueuedSeed: queuedSeed,
		Count:      int32(len(result.Statuses)),
		MaxId:      maxIdStr,
		SinceId:    sinceIdStr,
		RateLimit:  ratelimit.New().FromHttpHeaders(response.Header).ToProto(),
	}, nil

}

// register worker with server.
func (s *worker) register(server *grpc.Server) {
	pb.RegisterWorkerServer(server, s)
}

// twitterClient returns a twitter client.
func twitterClient() *twitter.Client {
	// init client only once
	once.Do(func() {
		httpClient := auth.HttpClient(accessToken)
		client = twitter.NewClient(httpClient)
	})
	return client
}
