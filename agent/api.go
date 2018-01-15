package agent

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/types"
)

type agentApi struct {
	store *db.Database
}

func NewApi() api.AgentServer {
	return &agentApi{db.New()}
}

// Implements AgentClient gRPC interface
func (a *agentApi) Do(ctx context.Context, req *api.DoJobRequest) (*google_protobuf.Empty, error) {
	seed := new(types.Seed).FromProto(req.Seed)

	if seed.Meta.Name == "" {
		log.WithField("description", seed.Meta.Description).Debugln("Not enqueuing seed (no query)")
		return new(google_protobuf.Empty), nil
	}

	log.WithField("description", seed.Meta.Description).Debugln("Enqueueing seed")

	queuedSeed := &api.QueuedSeed{
		SeedId:     seed.Id,
		Parameters: &api.SearchParameters{Query: seed.Meta.Name},
	}
	err := a.enqueueSeed(queuedSeed)
	if err != nil {
		log.WithError(err).Errorln()
		return nil, err
	}

	return new(google_protobuf.Empty), nil
}

func (a *agentApi) enqueueSeed(queuedSeed *api.QueuedSeed) error {
	err := a.store.Connect()
	defer a.store.Disconnect()
	if err != nil {
		return err
	}

	_, err = a.store.Insert("seed_queue", queuedSeed)
	return err
}
