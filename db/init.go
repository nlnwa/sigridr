// Copyright © 2017 National Library of Norway
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

package db

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

func CreateDb(name string) {
	result, err := r.DBCreate(name).RunWrite(session)
	if err != nil {
		log.WithError(err).Errorln("Creating database")
	}
	if n := result.DBsCreated; n > 0 {
		log.WithFields(log.Fields{"name": name}).Debug("Database created")
	}
}

func CreateTable(name string) {
	result, err := r.TableCreate(name).RunWrite(session)
	if err != nil {
		log.WithError(err).Error("Creating table")
	}
	if n := result.TablesCreated; n > 0 {
		log.WithFields(log.Fields{"name": name}).Debug("Table created")
	}
}

func Insert(table string, data interface{}) {
	result, err := r.Table(table).Insert(data).RunWrite(session)
	if err != nil {
		log.WithError(err).Errorln("Inserting data")
	}
	log.Debugln(result)
}
