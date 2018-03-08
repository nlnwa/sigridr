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

package worker

import (
	"github.com/nlnwa/sigridr/database"
)

type workerStore struct {
	*database.Rethink
}

func newStore(c Config) *workerStore {
	db := database.New(
		database.WithAddress(c.DatabaseHost, c.DatabasePort),
		database.WithName(c.DatabaseName),
		database.WithCredentials(c.DatabaseUser, c.DatabasePassword))
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
