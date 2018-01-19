package agent

import (
	"context"
	"fmt"
	"time"

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
		return fmt.Errorf("dial: %v", err)
	} else {
		ac.agentClient = api.NewAgentClient(ac.cc)
		return
	}
}

func (ac *Client) Hangup() error {
	return ac.cc.Close()
}

func (ac *Client) Do(job *types.Job, seed *types.Seed) error {
	request := api.DoJobRequest{
		Job:  job.ToProto(),
		Seed: seed.ToProto(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ac.agentClient.Do(ctx, &request)
	return err
}
