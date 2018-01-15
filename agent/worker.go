package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/pkg/twitter/ratelimit"
)

var (
	conn       *grpc.ClientConn
	client     pb.WorkerClient
	queueStore *store.QueueStore
)

func Worker(ctx context.Context, c Config) error {
	timer := time.NewTimer(0)
	queueStore = store.New()

	defer timer.Stop()

	for {
		// wait for timer or return if done
		select {
		case <-timer.C:
			err := queueStore.Connect()
			if err != nil {
				timer.Reset(time.Minute)
				break
			}
			err := connect(c)
			if err != nil {
				log.WithError(err).Errorln("failed to connect, will sleep and try again later")
				timer.Reset(time.Minute)
				break
			}
			for queuedSeed := range queueStore.GetNextToFetch(ctx) {
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
			disconnect()
		case <-ctx.Done():
			return nil
		}
	}
}

// dispatch sends work to the client
func dispatch(ctx context.Context, queuedSeed *pb.QueuedSeed) (*ratelimit.RateLimit, error) {
	if queuedSeed.GetSeq() == 0 {
		params, err := queueStore.SearchParameters(queuedSeed.GetSeedId())
		if err != nil {
			return nil, err
		}
		queuedSeed.Parameters.SinceId = params.GetSinceId()
	}
	work := &pb.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	reply, err := client.Do(ctx, work)
	if err != nil {
		return nil, err
	}

	// remove seed from queue
	err = queueStore.DeleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return nil, err
	}

	if reply.QueuedSeed.GetSeq() == 0 {
		reply.QueuedSeed.Parameters.Id = queuedSeed.GetSeedId()
	}

	if reply.Count < ratelimit.MAX_STATUSES_PER_REQUEST {
		reply.QueuedSeed.Parameters.MaxId = ""
		reply.QueuedSeed.Parameters.SinceId = reply.GetSinceId()
	} else {
		reply.QueuedSeed.Parameters.MaxId = reply.GetMaxId()

		reply.QueuedSeed.Seq++
		queueStore.EnqueueSeed(reply.QueuedSeed)
	}
	err = queueStore.SaveSearchParameters(reply.QueuedSeed.Parameters)
	if err != nil {
		return nil, err
	}

	return new(ratelimit.RateLimit).FromProto(reply.RateLimit), nil
}

// connect establishes:
// - a connection to gRPC server and creates a new client which use the connection.
// - a database session
func connect(c Config) error {
	var err error
	opts := grpc.WithInsecure()
	conn, err = grpc.Dial(c.Worker, opts)
	if err != nil {
		return err
	}
	client = pb.NewWorkerClient(conn)

	err = queueStore.Connect()
	if err != nil {
		return err
	}
	return nil
}

// disconnect closes the connection to the gRPC server and the database session
func disconnect() {
	queueStore.Disconnect()
	conn.Close()
}
