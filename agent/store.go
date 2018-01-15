package store

import (
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/nlnwa/sigridr/api/sigridr"
	"github.com/nlnwa/sigridr/pkg/db"
	// "github.com/nlnwa/sigridr/pkg/types"
)

type QueueStore struct {
	*db.Database
}

func New() *QueueStore {
	return &QueueStore{db.New()}
}

func (qs *QueueStore) EnqueueSeed(queuedSeed *pb.QueuedSeed) error {
	_, err := qs.Insert("seed_queue", queuedSeed)
	return err
}

func (qs *QueueStore) SaveSearchParameters(params *pb.SearchParameters) error {
	_, err := qs.Insert("search_parameters", params)
	return err
}

func (qs *QueueStore) SearchParameters(id string) (*pb.SearchParameters, error) {
	var params *pb.SearchParameters
	err := qs.Get("search_parameters", id, params)
	return params, err
}

func (qs *QueueStore) DeleteQueuedSeed(id string) error {
	return qs.Delete("seed_queue", id)
}

func (qs *QueueStore) GetNextToFetch(ctx context.Context) <-chan *pb.QueuedSeed {
	out := make(chan *pb.QueuedSeed)
	go func() {
		defer close(out)

		cursor, err := qs.GetCursor("seed_queue")
		defer cursor.Close()
		if err != nil {
			log.WithError(err).Errorln("getting cursor to seed_queue")
			out <- nil
			return
		}
		queuedSeed := new(pb.QueuedSeed)
		for {
			if ok := cursor.Next(queuedSeed); !ok {
				if err = cursor.Err(); err != nil {
					log.WithError(err).Errorln("failed getting next row in seed_queue")
					return
				} else {
					queuedSeed = nil
				}
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
