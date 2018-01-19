package agent

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
	"github.com/nlnwa/sigridr/worker"
)

type QueueWorker interface {
	Run(context.Context) error
}

type queueWorker struct {
	store        *agentStore
	workerClient *worker.ApiClient
}

func NewQueueWorker(c Config) QueueWorker {
	return &queueWorker{
		store:        newStore(c),
		workerClient: worker.NewApiClient(c.WorkerAddress),
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
			err = qw.workerClient.Dial()
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
			qw.workerClient.Hangup()
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

	// call worker
	reply, err := qw.workerClient.Do(ctx, queuedSeed)
	if err != nil {
		return nil, err
	}

	// remove seed from queue
	err = qw.store.deleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return nil, err
	}

	// enqueue next fetch if possibly more statuses to fetch
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
