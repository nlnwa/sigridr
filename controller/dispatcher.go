package controller

import (
	"fmt"
	"context"
	"time"

	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type agentClient struct {
	address string
	cc      *grpc.ClientConn
}

func (ac *agentClient) dial() (api.AgentClient, error) {
	var err error
	ac.cc, err = grpc.Dial(ac.address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return api.NewAgentClient(ac.cc), nil
}

func (ac *agentClient) hangup() error {
	return ac.cc.Close()
}

type dispatcher struct {
	client *agentClient
}

func newDispatcher(c Config) *dispatcher {
	return &dispatcher{&agentClient{address: c.AgentAddress}}
}

func (d *dispatcher) dispatch(job *types.Job, seed *types.Seed) {
	request := api.DoJobRequest{
		Job:  job.ToProto(),
		Seed: seed.ToProto(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, err := d.client.dial()
	if err != nil {
		log.WithError(err).Errorln()
	}
	defer d.client.hangup()

	_, err = client.Do(ctx, &request)
	if err != nil {
		log.WithError(err).Error()
	} else {
		log.WithField("seed", seed.Meta.Description).Debugln("Dispatch")
	}
}
