// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agent

import (
	"context"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
)

type agentStore struct {
	db *database.Rethink
}

func newStore(c Config) *agentStore {
	db := database.New(
		database.WithAddress(c.DatabaseHost, c.DatabasePort),
		database.WithName(c.DatabaseName),
		database.WithCredentials(c.DatabaseUser, c.DatabasePassword))
	db.SetTags("json")

	return &agentStore{db}
}

func (qs *agentStore) connect() error {
	return qs.db.Connect()
}

func (qs *agentStore) disconnect() error {
	return qs.db.Disconnect()
}

func (qs *agentStore) enqueueSeed(queuedSeed *api.QueuedSeed) error {
	_, err := qs.db.Insert("queue", queuedSeed)
	return err
}

func (qs *agentStore) updateParameter(param *api.Parameter) error {
	return qs.db.Update("parameter", param.Id, param)
}

func (qs *agentStore) saveParameter(param *api.Parameter) error {
	_, err := qs.db.Insert("parameter", param)
	return err
}

func (qs *agentStore) parameter(id string) (*api.Parameter, error) {
	param := new(api.Parameter)
	if err := qs.db.Get("parameter", id, param); err != nil {
		return nil, err
	} else {
		return param, nil
	}
}

func (qs *agentStore) status(id string) (*status, error) {
	exec := new(types.Execution)
	if err := qs.db.Get("execution", id, exec); err != nil {
		return nil, err
	} else {
		return new(status).fromExecution(exec), nil
	}
}

func (qs *agentStore) saveStatus(s *status) (string, error) {
	if s.Execution.Id == "" {
		return qs.db.Insert("execution", s.Execution)
	} else {
		return s.Execution.Id, qs.db.Update("execution", s.Execution.Id, s.Execution)
	}

}

func (qs *agentStore) deleteQueuedSeed(id string) error {
	return qs.db.Delete("queue", id)
}

func (qs *agentStore) getNextToFetch(ctx context.Context) (<-chan *api.QueuedSeed, <-chan error) {
	out := make(chan *api.QueuedSeed)
	errc := make(chan error)
	go func() {
		defer close(out)
		defer close(errc)

		cursor, err := qs.db.GetCursor("queue")
		if err != nil {
			errc <- err
			return
		}
		defer cursor.Close()

		for {
			queuedSeed := new(api.QueuedSeed)
			if ok := cursor.Next(queuedSeed); !ok {
				if err = cursor.Err(); err != nil {
					errc <- err
					return
				} else {
					out <- nil
					return
				}
			}

			// return if done else send next to fetch on channel
			select {
			case <-ctx.Done():
				return
			case out <- queuedSeed:
				break
			}
		}
	}()
	return out, errc
}
