package agent

import (
	"fmt"

	pb "github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type agentApi struct {
	store *agentStore
}

func NewApi(c Config) api.AgentServer {
	return &agentApi{
		store: newStore(c),
	}
}

func (a *agentApi) Do(ctx context.Context, req *api.DoJobRequest) (*pb.Empty, error) {
	seed := new(types.Seed).FromProto(req.Seed)

	if seed.Meta.Name == "" {
		log.WithField("description", seed.Meta.Description).Debugln("Not enqueuing seed (no query)")
		return new(pb.Empty), nil
	}

	log.WithField("description", seed.Meta.Description).Debugln("Enqueueing seed")

	queuedSeed := &api.QueuedSeed{
		SeedId:    seed.Id,
		Parameter: &api.Parameter{Query: seed.Meta.Name},
	}
	if err := a.store.connect(); err != nil {
		return nil, fmt.Errorf("failed connecting to database: %v", err)
	}
	if err := a.store.enqueueSeed(queuedSeed); err != nil {
		return nil, fmt.Errorf("failed enqueuing seed: %v", err)
	}

	return new(pb.Empty), a.store.Disconnect()
}
