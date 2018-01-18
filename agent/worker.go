package agent

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"

	"google.golang.org/grpc"
)

var client api.WorkerClient

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

type QueueWorker interface {
	Run(context.Context) error
}

type queueWorker struct {
	store  *agentStore
	worker *workerClient
}

func NewQueueWorker(c Config) QueueWorker {
	return &queueWorker{
		store:  newStore(c),
		worker: &workerClient{address: c.WorkerAddress},
	}
}

func (qw *queueWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		// wait for timer or return if done
		select {
		case <-timer.C:
			err := qw.store.connect()
			if err != nil {
				timer.Reset(time.Minute)
				break
			}
			client, err = qw.worker.dial()
			if err != nil {
				log.WithError(err).Errorln("failed to connect, will sleep and try again later")
				timer.Reset(time.Minute)
				break
			}
			for queuedSeed := range qw.store.getNextToFetch(ctx) {
				if queuedSeed == nil {
					timer.Reset(time.Minute)
					break
				}
				rateLimit, err := qw.dispatch(ctx, queuedSeed)
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
			qw.store.Disconnect()
			qw.worker.hangup()
		case <-ctx.Done():
			return nil
		}
	}
}

// dispatch sends work to the client
func (qw *queueWorker) dispatch(ctx context.Context, queuedSeed *api.QueuedSeed) (*ratelimit.RateLimit, error) {
	seq := queuedSeed.GetSeq()

	if queuedSeed.Parameter.GetId() == "" {
		param, err := qw.store.parameter(queuedSeed.GetSeedId())
		if err != nil {
			return nil, err
		}
		if seq == 0 {
			queuedSeed.Parameter.Id = param.GetId()
			queuedSeed.Parameter.SinceId = param.GetSinceId()
		}
	}
	work := &api.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	reply, err := client.Do(ctx, work)
	if err != nil {
		return nil, err
	}

	// remove seed from queue
	err = qw.store.deleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return nil, err
	}

	// if possibly more to fetch enqueue next fetch
	if reply.Count >= twitter.MaxStatusesPerRequest {
		reply.QueuedSeed.Parameter.MaxId = reply.GetMaxId()
		reply.QueuedSeed.Seq++
		qw.store.enqueueSeed(reply.QueuedSeed)
	}

	// only save/update parameters if first in sequence
	if seq == 0 {
		queuedSeed.Parameter.MaxId = "" // don't need to save this
		queuedSeed.Parameter.SinceId = reply.GetSinceId()

		if queuedSeed.Parameter.GetId() == "" {
			queuedSeed.Parameter.Id = queuedSeed.GetSeedId()
			err = qw.store.saveParameter(queuedSeed.Parameter)
		} else {
			err = qw.store.updateParameter(queuedSeed.Parameter)
		}
		if err != nil {
			return nil, err
		}
	}

	return new(ratelimit.RateLimit).FromProto(reply.RateLimit), nil
}
