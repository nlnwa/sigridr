package worker

import (
	"github.com/nlnwa/sigridr/database"
)

type workerStore struct {
	*database.Rethink
}

func newStore(c Config) *workerStore {
	db := database.New(database.WithAddress(c.DatabaseHost, c.DatabasePort), database.WithName(c.DatabaseName))
	db.SetTags("json")

	return &workerStore{
		Rethink: db,
	}
}

func (ws *workerStore) saveSearchResult(value interface{}) (string, error) {
	if err := ws.Connect(); err != nil {
		return "", err
	} else {
		defer ws.Disconnect()
	}
	if id, err := ws.Insert("result", value); err != nil {
		return "", err
	} else {
		return id, nil
	}
}
