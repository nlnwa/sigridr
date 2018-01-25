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

package agent

import (
	"context"
	"time"

	"github.com/nlnwa/pkg/log"
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
	log.Logger
}

func NewQueueWorker(c Config) QueueWorker {
	return &queueWorker{
		store:        newStore(c),
		workerClient: worker.NewApiClient(c.WorkerAddress),
		Logger:       c.Logger,
	}
}

func (qw *queueWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {

		// wait for timer or return if done
	timer:
		select {
		case <-timer.C:
			sleep := time.Minute

			err := qw.store.connect()
			if err != nil {
				qw.Error("failed to connect to store", "error", err, "sleep", sleep.String())
				timer.Reset(sleep)
				break timer
			}
			err = qw.workerClient.Dial()
			if err != nil {
				qw.Error("failed to connect to worker", "error", err, "sleep", sleep.String())
				timer.Reset(time.Minute)
				break timer
			}

			out, errc := qw.store.getNextToFetch(ctx)
			for {
				select {
				case queuedSeed := <-out:
					if queuedSeed == nil {
						timer.Reset(sleep)
						break timer
					}
					rateLimit, err := qw.dispatch(ctx, queuedSeed)
					if err != nil {
						qw.Error("failed to dispatch queued seed", "error", err, "sleep", sleep.String())
						timer.Reset(sleep)
						break timer
					}
					qw.Debug("Ratelimit",
						"remaining", rateLimit.Remaining,
						"reset", rateLimit.Reset,
						"timeout", rateLimit.Timeout())

					if rateLimit.Remaining < 1 {
						sleep = rateLimit.Timeout() + 5*time.Second
						qw.Info("Ratelimit reached", "sleep", sleep.String())
						timer.Reset(sleep)
						break timer
					}
				case err := <-errc:
					qw.Error("failed fetching next queued seed", "error", err)
					timer.Reset(time.Minute * 10)
					break timer
				}

			}
		case <-ctx.Done():
			return nil
		}
		qw.store.Disconnect()
		qw.workerClient.Hangup()
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
	if err = qw.store.deleteQueuedSeed(queuedSeed.Id); err != nil {
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

	return new(ratelimit.RateLimit).FromProto(reply.RateLimit)
}
