package agent

import (
	"context"

	pb "github.com/golang/protobuf/ptypes/empty"

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
	seed, err := new(types.Seed).FromProto(req.Seed)
	if err != nil {
		return nil, err
	}

	if seed.Meta.Name == "" {
		return new(pb.Empty), nil
	}

	queuedSeed := &api.QueuedSeed{
		SeedId:    seed.Id,
		Parameter: &api.Parameter{Query: seed.Meta.Name},
	}
	if err := a.store.connect(); err != nil {
		return nil, err
	}
	if err := a.store.enqueueSeed(queuedSeed); err != nil {
		return nil, err
	}

	return new(pb.Empty), a.store.Disconnect()
}
