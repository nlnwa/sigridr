package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	netcontext "golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/nlnwa/sigridr/api/sigridr"
	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/twitter/ratelimit"
)

var (
	conn   *grpc.ClientConn
	client pb.WorkerClient
)

func queueWorker(ctx context.Context, wg *sync.WaitGroup) {
	timer := time.NewTimer(0)

	defer func() {
		timer.Stop()
		wg.Done()
	}()

	for {
		// wait for timer or return if done
		select {
		case <-timer.C:
			err := connect()
			if err != nil {
				log.WithError(err).Errorln("failed to connect, will sleep and try again later")
				timer.Reset(time.Minute)
				break
			}
			for queuedSeed := range db.GetNextToFetch(ctx) {
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

				// stop fetching if ratelimit
				if rateLimit.Remaining < 1 {
					timer.Reset(rateLimit.Timeout() + 5*time.Second)
					break
				}
			}
			disconnect()
		case <-ctx.Done():
			return
		}
	}
}

// dispatch sends work to the client
func dispatch(ctx context.Context, queuedSeed *pb.QueuedSeed) (*ratelimit.RateLimit, error) {
	if queuedSeed.GetSeq() == 0 {
		params, err := db.SearchParameters(queuedSeed.GetSeedId())
		if err != nil {
			return nil, err
		}
		queuedSeed.Parameters.SinceId = params.GetSinceId()
	}
	work := &pb.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := netcontext.WithTimeout(ctx, time.Minute)
	defer cancel()

	reply, err := client.Do(ctx, work)
	if err != nil {
		return nil, err
	}

	// remove seed from queue
	err = db.DeleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return nil, err
	}

	if reply.QueuedSeed.GetSeq() == 0 {
		reply.QueuedSeed.Parameters.Id = queuedSeed.GetSeedId()
	}

	if reply.Count < ratelimit.MAX_STATUSES_PER_REQUEST {
		reply.QueuedSeed.Parameters.MaxId = ""
		reply.QueuedSeed.Parameters.SinceId = reply.GetSinceId()
		db.SaveSearchParameters(reply.QueuedSeed.Parameters)
	} else {
		reply.QueuedSeed.Parameters.MaxId = reply.GetMaxId()
		db.SaveSearchParameters(reply.QueuedSeed.Parameters)

		reply.QueuedSeed.Seq++
		db.EnqueueSeed(reply.QueuedSeed)
	}
	return new(ratelimit.RateLimit).FromProto(reply.RateLimit), nil
}

// connect establishes a connection to gRPC server and creates a new client which use the connection.
func connect() error {
	var err error
	opts := grpc.WithInsecure()
	conn, err = grpc.Dial(*workerAddress, opts)
	if err != nil {
		return err
	}
	client = pb.NewWorkerClient(conn)
	return nil
}

// disconnect closes the connection to the gRPC server.
func disconnect() {
	conn.Close()
}
