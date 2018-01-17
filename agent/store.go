package agent

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/database"
)

type agentStore struct {
	*database.Rethink
	*database.ConnectOpts
}

func newStore(c Config) *agentStore {
	return &agentStore{
		Rethink: database.New(),
		ConnectOpts: &database.ConnectOpts{
			Database: c.DatabaseName,
			Address:  c.DatabaseAddress,
		},
	}
}

func (qs *agentStore) connect() error {
	return qs.Rethink.Connect(qs.ConnectOpts)
}

func (qs *agentStore) enqueueSeed(queuedSeed *api.QueuedSeed) error {
	_, err := qs.Insert("queue", queuedSeed)
	return err
}

func (qs *agentStore) updateParameter(param *api.Parameter) error {
	return qs.Update("parameter", param.Id, param)
}

func (qs *agentStore) saveParameter(param *api.Parameter) error {
	_, err := qs.Insert("parameter", param)
	return err
}

func (qs *agentStore) parameter(id string) (*api.Parameter, error) {
	param := new(api.Parameter)
	err := qs.Get("parameter", id, param)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (qs *agentStore) deleteQueuedSeed(id string) error {
	return qs.Delete("queue", id)
}

func (qs *agentStore) getNextToFetch(ctx context.Context) <-chan *api.QueuedSeed {
	out := make(chan *api.QueuedSeed)
	go func() {
		defer close(out)

		cursor, err := qs.GetCursor("queue")
		defer cursor.Close()
		if err != nil {
			log.WithError(err).Errorln("failed getting cursor to seed queue")
			out <- nil
			return
		}
		for {
			queuedSeed := new(api.QueuedSeed)
			if ok := cursor.Next(queuedSeed); !ok {
				if err = cursor.Err(); err != nil {
					log.WithError(err).Errorln("failed getting next row in seed queue")
				}
				queuedSeed = nil
			}

			// return if done else send next to fetch on channel
			select {
			case <-ctx.Done():
				return
			case out <- queuedSeed:
				if queuedSeed == nil {
					return
				}
			}
		}
	}()
	return out
}
