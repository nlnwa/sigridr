package agent

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type Client struct {
	address     string
	cc          *grpc.ClientConn
	agentClient api.AgentClient
}

func NewApiClient(address string) *Client {
	return &Client{address: address}
}

func (ac *Client) Dial() (err error) {
	if ac.cc, err = grpc.Dial(ac.address, grpc.WithInsecure()); err != nil {
		return errors.Wrapf(err, "failed to dial: %s", ac.address)
	} else {
		ac.agentClient = api.NewAgentClient(ac.cc)
		return
	}
}

func (ac *Client) Hangup() error {
	return ac.cc.Close()
}

func (ac *Client) Do(job *types.Job, seed *types.Seed) error {
	j, err := job.ToProto()
	if err != nil {
		return err
	}
	s, err := seed.ToProto()
	if err != nil {
		return err
	}
	request := api.DoJobRequest{
		Job:  j,
		Seed: s,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = ac.agentClient.Do(ctx, &request)
	return err
}
