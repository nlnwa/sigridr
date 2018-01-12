package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api/sigridr"
)

func EnqueueSeed(queuedSeed *sigridr.QueuedSeed) error {
	Connect()
	defer Disconnect()

	_, err := Insert("seed_queue", queuedSeed)
	return err
}

func DeleteQueuedSeed(id string) error {
	Connect()
	defer Disconnect()
	return Delete("seed_queue", id)
}

func GetNextToFetch(ctx context.Context) <-chan *sigridr.QueuedSeed {
	out := make(chan *sigridr.QueuedSeed)
	go func() {
		defer close(out)
		Connect()
		defer Disconnect()

		cursor, err := GetCursor("seed_queue")
		defer cursor.Close()
		if err != nil {
			log.WithError(err).Errorln("getting cursor to seed_queue")
			out <- nil;
			return
		}

		for {
			queuedSeed := new(sigridr.QueuedSeed)
			if ok := cursor.Next(queuedSeed); !ok {
				if err = cursor.Err(); err != nil {
					log.WithError(err).Errorln("failed getting next row in seed_queue")
					return
				} else {
					queuedSeed = nil
				}
			}

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

func SaveSearchParameters(params *sigridr.SearchParameters) (string, error) {
	Connect()
	defer Disconnect()
	return Insert("search_parameters", params)
}

func SearchParameters(id string) (*sigridr.SearchParameters, error) {
	Connect()
	defer Disconnect()
	var params *sigridr.SearchParameters
	err := Get("search_parameters", id, params)
	return params, err
}
