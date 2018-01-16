package agent

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/database"
)

type queueStore struct {
	*database.Rethink
}

func newStore(c Config) *queueStore {
	db := database.New()
	db.ConnectOpts = database.DefaultOptions()
	db.ConnectOpts.Database = c.DatabaseName
	db.ConnectOpts.Address = c.DatabaseAddress

	return &queueStore{db}
}

func (qs *queueStore) connect() error {
	return qs.Rethink.Connect()
}

func (qs *queueStore) enqueueSeed(queuedSeed *api.QueuedSeed) error {
	_, err := qs.Insert("queue", queuedSeed)
	return err
}

func (qs *queueStore) updateParameter(param *api.Parameter) error {
	return qs.Update("parameter", param.Id, param)
}

func (qs *queueStore) saveParameter(param *api.Parameter) error {
	_, err := qs.Insert("parameter", param)
	return err
}

func (qs *queueStore) parameter(id string) (*api.Parameter, error) {
	param := new(api.Parameter)
	err := qs.Get("parameter", id, param)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (qs *queueStore) deleteQueuedSeed(id string) error {
	return qs.Delete("queue", id)
}

func (qs *queueStore) getNextToFetch(ctx context.Context) <-chan *api.QueuedSeed {
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
