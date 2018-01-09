package main

import (
	"time"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/nlnwa/sigridr/api/sigridr"
	"github.com/nlnwa/sigridr/pkg/types"
	log "github.com/sirupsen/logrus"
)

func dispatch(job *types.Job, seed *types.Seed) {
	request := sigridr.DoJobRequest{
		Job:  job.ToProto(),
		Seed: seed.ToProto(),
	}

	log.WithField("request", request).Debugln("Dispatch request to agent")

	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(address, opts)
	if err != nil {
		log.WithError(err).Fatal()
	}
	defer conn.Close()

	agent := sigridr.NewAgentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = agent.Do(ctx, &request)
	if err != nil {
		log.WithError(err).Error()
	}
}
