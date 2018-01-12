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
	conn      *grpc.ClientConn
	client    pb.WorkerClient
	rateLimit *ratelimit.RateLimit
)

func queueWorker(ctx context.Context, wg *sync.WaitGroup) {
	needSleep := false
	enqueued := false
	timeout := time.Duration(0)
	timer := time.NewTimer(timeout)
	rateLimit = ratelimit.New()

	defer func() {
		if !timer.Stop() {
			<-timer.C // drain timer channel
		}
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if rateLimit.Remaining < 1 {
				timeout = rateLimit.Timeout() + 5 * time.Second
			} else if enqueued {
				// queuedSeed was enqueued last dispatch
				timeout = 0
			} else if needSleep {
				// seed queue empty or dispatch error
				timeout = time.Minute
			}
			log.WithField("timeout", timeout).Debugln("sleeping")
			timer.Reset(timeout)
		}

		select {
		case <-timer.C:
			err := connect()
			if err != nil {
				log.WithError(err).Errorln("failed to connect, will sleep and try again later")
				needSleep = true // try to sleep the connection error away
				break
			}
			for queuedSeed := range db.GetNextToFetch(ctx) {
				if queuedSeed == nil {
					needSleep = true
					break
				}
				enqueued, err = dispatch(ctx, queuedSeed)
				if err != nil {
					log.WithError(err).Errorln("dispatch")
					needSleep = true
					break
				}
				log.WithFields(log.Fields{
					"remaining": rateLimit.Remaining,
					"reset": rateLimit.Reset,
					"timeout": rateLimit.Timeout(),
					}).Debugln("Ratelimit")
				if rateLimit.Remaining < 1 {
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
func dispatch(ctx context.Context, queuedSeed *pb.QueuedSeed) (bool, error) {
	if queuedSeed.Parameters.GetMaxId() == "" {
		params, err := db.SearchParameters(queuedSeed.GetSeedId())
		if err != nil {
			return false, err
		}
		queuedSeed.Parameters.MaxId = params.GetMaxId()
		queuedSeed.Parameters.SinceId = params.GetSinceId()
	}

	log.WithField("parameters", queuedSeed.Parameters).Debugln("dispatch")
	work := &pb.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := netcontext.WithTimeout(ctx, time.Minute)
	defer cancel()
	reply, err := client.Do(ctx, work)
	if err != nil {
		return false, err
	}

	err = db.DeleteQueuedSeed(queuedSeed.Id)
	if err != nil {
		return false, err
	}

	rateLimit.FromProto(reply.RateLimit)

	reply.QueuedSeed.Parameters.Id = queuedSeed.GetSeedId()
	db.SaveSearchParameters(reply.QueuedSeed.Parameters)

	// if count equals maximum per request (from twitter) enqueue new request to
	// possibly fetch more
	if reply.GetCount() >= 100 {
		reply.QueuedSeed.Seq++
		db.EnqueueSeed(reply.QueuedSeed)
		return true, nil
	} else {
		return false, nil
	}
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
