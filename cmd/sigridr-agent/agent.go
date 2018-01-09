package main

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nlnwa/sigridr/api/sigridr"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/types"
)

type agent struct{}

// Implements AgentClient gRPC interface
func (a *agent) Do(ctx context.Context, req *pb.DoJobRequest) (*google_protobuf.Empty, error) {
	seed := new(types.Seed).FromProto(req.Seed)

	queuedSeed := &pb.QueuedSeed{
		SeedId:     seed.Id,
		Parameters: &pb.SearchParameters{Query: seed.Meta.Name},
	}
	err := db.EnqueueSeed(queuedSeed)
	if err != nil {
		log.WithError(err).Errorln()
		return nil, err
	}

	return new(google_protobuf.Empty), nil
}

func (a *agent) register(server *grpc.Server) {
	pb.RegisterAgentServer(server, a)
}
