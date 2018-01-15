package agent

import (
	"fmt"

	pb "github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"

)

type agentApi struct {
	store *database.Rethink
}

func NewApi() api.AgentServer {
	return &agentApi{database.New()}
}

// Implements AgentClient gRPC interface
func (a *agentApi) Do(ctx context.Context, req *api.DoJobRequest) (*pb.Empty, error) {
	seed := new(types.Seed).FromProto(req.Seed)

	if seed.Meta.Name == "" {
		log.WithField("description", seed.Meta.Description).Debugln("Not enqueuing seed (no query)")
		return new(pb.Empty), nil
	}

	log.WithField("description", seed.Meta.Description).Debugln("Enqueueing seed")

	queuedSeed := &api.QueuedSeed{
		SeedId:     seed.Id,
		Parameter: &api.Parameter{Query: seed.Meta.Name},
	}
	err := a.enqueueSeed(queuedSeed)
	if err != nil {
		return nil, fmt.Errorf("failed enqueuing seed: %v", err)
	}

	return new(pb.Empty), nil
}


func (a *agentApi) enqueueSeed(queuedSeed *api.QueuedSeed) error {
	err := a.store.Connect(database.DefaultOptions())
	defer a.store.Disconnect()
	if err != nil {
		return err
	}
	_, err = a.store.Insert("queue", queuedSeed)
	return err
}


