package worker

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
)

type ApiClient struct {
	address      string
	cc           *grpc.ClientConn
	workerClient api.WorkerClient
}

func NewApiClient(address string) *ApiClient {
	return &ApiClient{address: address}
}

func (c *ApiClient) Dial() (err error) {
	if c.cc, err = grpc.Dial(c.address, grpc.WithInsecure()); err != nil {
		return errors.Wrapf(err, "failed to dial: %s", c.address)
	} else {
		c.workerClient = api.NewWorkerClient(c.cc)
		return
	}
}

func (c *ApiClient) Hangup() error {
	return c.cc.Close()
}

func (c *ApiClient) Do(ctx context.Context, queuedSeed *api.QueuedSeed) (*api.WorkReply, error) {
	work := &api.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	if reply, err := c.workerClient.Do(ctx, work); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}
