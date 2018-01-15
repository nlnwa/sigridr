package agent

import (
	"context"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"

	"google.golang.org/grpc"
)

var (
	worker *workerClient
	client api.WorkerClient
	db     *queueStore
)

type Config struct {
	Worker string
}

type workerClient struct {
	address string
	cc      *grpc.ClientConn
}

func (wc *workerClient) dial() (api.WorkerClient, error) {
	conn, err := grpc.Dial(wc.address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	wc.cc = conn
	return api.NewWorkerClient(conn), nil
}

func (wc *workerClient) hangup() error {
	return wc.cc.Close()
}

func QueueWorker(ctx context.Context, c Config) error {
	timer := time.NewTimer(0)
	db = newStore()
	worker = &workerClient{address: c.Worker}

	defer timer.Stop()

	for {
		// wait for timer or return if done
		select {
		case <-timer.C:
			err := db.connect()
			if err != nil {
				timer.Reset(time.Minute)
				break
			}
			client, err = worker.dial()
			if err != nil {
				log.WithError(err).Errorln("failed to connect, will sleep and try again later")
				timer.Reset(time.Minute)
				break
			}
			for queuedSeed := range db.getNextToFetch(ctx) {
				if queuedSeed == nil {
					timer.Reset(time.Minute)
					break
				}
				rateLimit, err := dispatch(ctx, queuedSeed)
				if err != nil {
					log.WithError(err).Errorln("dispatching queued seed")
					timer.Reset(time.Minute)
					break
				}
				log.WithFields(log.Fields{
					"remaining": rateLimit.Remaining,
					"reset":     rateLimit.Reset,
					"timeout":   rateLimit.Timeout(),
				}).Debugln("Ratelimit")

				// stop fetching if ratelimit reached
				if rateLimit.Remaining < 1 {
					timer.Reset(rateLimit.Timeout() + 5*time.Second)
					break
				}
			}
			db.Disconnect()
			worker.hangup()
		case <-ctx.Done():
			return nil
		}
	}
}

// dispatch sends work to the client
func dispatch(ctx context.Context, queuedSeed *api.QueuedSeed) (*ratelimit.RateLimit, error) {
	param, err := db.parameter(queuedSeed.GetSeedId())
	if err != nil {
		return nil, err
	}
	if queuedSeed.GetSeq() == 0 && param.GetId() != "" {
		queuedSeed.Parameter.SinceId = param.GetSinceId()
	}
	work := &api.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	reply, err := client.Do(ctx, work)
	if err != nil {
		return nil, err
	}

	// remove seed from queue
	err = db.deleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return nil, err
	}

	// update queued seed parameters
	if reply.Count < twitter.MAX_STATUSES_PER_REQUEST {
		reply.QueuedSeed.Parameter.MaxId = ""
		reply.QueuedSeed.Parameter.SinceId = reply.GetSinceId()
	}

	// only save/update parameter if first in sequence
	if reply.QueuedSeed.GetSeq() == 0 {
		reply.QueuedSeed.Parameter.Id = queuedSeed.GetSeedId()
		if param.GetId() == "" {
			err = db.saveParameter(reply.QueuedSeed.Parameter)
		} else {
			err = db.updateParameter(reply.QueuedSeed.Parameter)
		}
		if err != nil {
			return nil, err
		}
	}

	// enqueue next fetch of queued seed
	if reply.Count >= twitter.MAX_STATUSES_PER_REQUEST {
		reply.QueuedSeed.Parameter.MaxId = reply.GetMaxId()
		reply.QueuedSeed.Seq++
		db.enqueueSeed(reply.QueuedSeed)
	}

	return new(ratelimit.RateLimit).FromProto(reply.RateLimit), nil
}
