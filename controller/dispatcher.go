package controller

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

func dispatch(job *types.Job, seed *types.Seed) {
	request := api.DoJobRequest{
		Job:  job.ToProto(),
		Seed: seed.ToProto(),
	}

	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(address, opts)
	if err != nil {
		log.WithError(err).Error()
	}
	defer conn.Close()

	agent := api.NewAgentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = agent.Do(ctx, &request)
	if err != nil {
		log.WithError(err).Error()
	} else {
		log.WithField("seed", seed.Meta.Description).Infoln("Dispatch")
	}
}
