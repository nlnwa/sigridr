package worker

import (
	"fmt"

	"github.com/nlnwa/sigridr/database"
)

type workerStore struct {
	*database.Rethink
	*database.ConnectOpts
}

func newStore(c Config) *workerStore {
	return &workerStore{
		Rethink: database.New(),
		ConnectOpts: &database.ConnectOpts{
			Database: c.DatabaseName,
			Address:  c.DatabaseAddress,
		},
	}
}

func (ws *workerStore) saveSearchResult(value interface{}) (string, error) {
	if err := ws.Connect(ws.ConnectOpts); err != nil {
		return "", err
	} else {
		defer ws.Disconnect()
	}
	if id, err := ws.Insert("result", value); err != nil {
		return "", fmt.Errorf("inserting search result into database: %v", err)
	} else {
		return id, nil
	}
}
