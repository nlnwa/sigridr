package controller

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type agentClient struct {
	address string
	cc      *grpc.ClientConn
}

func newAgentClient(c Config) *agentClient {
	return &agentClient{address: c.AgentAddress}
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

func (ac *agentClient) dispatch(job *types.Job, seed *types.Seed) {
	request := api.DoJobRequest{
		Job:  job.ToProto(),
		Seed: seed.ToProto(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	agent, err := ac.dial()
	if err != nil {
		log.WithError(err).Errorln()
	}
	defer ac.hangup()

	_, err = agent.Do(ctx, &request)
	if err != nil {
		log.WithError(err).Error()
	} else {
		log.WithField("seed", seed.Meta.Description).Infoln("Start fetching")
	}
}
