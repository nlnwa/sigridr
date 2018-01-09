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

type Options = r.ConnectOpts

var (
	session        *r.Session
	defaultOptions = Options{Database: "sigridr"}
)

func init() {
	// https://github.com/GoRethink/gorethink#encodingdecoding
	r.SetTags("gorethink", "json", "url")
}

func NoOptions() Options {
	return Options{}
}

func DefaultOptions() Options {
	return defaultOptions
}

func ConnectWithOptions(opts Options) {
	if IsConnected() {
		return
	}
	var err error
	// https://godoc.org/github.com/GoRethink/gorethink#ConnectOpts
	session, err = r.Connect(opts)
	if err != nil {
		log.WithError(err).Fatal("Connecting to db")
	} else {
		log.WithField("connected", session.IsConnected()).Debugln("Database session")
	}
}

func Connect() {
	ConnectWithOptions(defaultOptions)
}

func ReConnect() {
	if IsConnected() {
		return
	}
	err := session.Reconnect(r.CloseOpts{})
	if err != nil {
		log.WithError(err).Fatal("Reconnecting to db")
	}
}

func IsConnected() bool {
	return session != nil && session.IsConnected()
}

func Disconnect() {
	if IsConnected() {
		session.Close()
		log.WithField("connected", session.IsConnected()).Debugln("Database session")
	}
}

func DropDatabase(name string) {
	result, err := r.DBDrop(name).RunWrite(session)
	if err != nil {
		log.WithError(err).Errorln("Dropping database")
	}
	if n := result.DBsDropped; n > 0 {
		log.WithField("name", name).Debug("Database dropped")
	}
}

func CreateDatabase(name string) {
	result, err := r.DBCreate(name).RunWrite(session)
	if err != nil {
		log.WithError(err).Errorln("Creating database")
	}
	if n := result.DBsCreated; n > 0 {
		log.WithField("name", name).Debug("Database created")
	}
}

func DropTable(name string) {
	result, err := r.TableDrop(name).RunWrite(session)
	if err != nil {
		log.WithError(err).WithField("name", name).Errorln("Dropping table")
	}
	if n := result.TablesDropped; n > 0 {
		log.WithField("name", name).Debug("Database table dropped")
	}

}

func CreateTable(name string) error {
	result, err := r.TableCreate(name).RunWrite(session)
	if err != nil {
		log.WithError(err).Error("Creating table")
		return err
	}
	if n := result.TablesCreated; n > 0 {
		log.WithField("name", name).Debug("Database table created")
	}
	return nil
}

func Insert(table string, document interface{}) (string, error) {
	res, err := r.Table(table).Insert(document).RunWrite(session)
	if err != nil {
		return "", err
	}
	id  := ""
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	log.WithFields(log.Fields{
		"id": id,
		"table": table,
	}).Debugln("Database insert document")
	return id, nil
}

func Delete(table string, id string) error {
	_, err := r.Table(table).Get(id).Delete().Run(session)
	if err == nil {
		log.WithFields(log.Fields{
			"id": id,
			"table": table,
		}).Debugln("Database delete document")
	}
	return err
}

func Use(name string) {
	session.Use(name)
}

func Changes(name string) (*r.Cursor, error) {
	return r.Table(name).Changes().Run(session)
}

func Get(table string, id string, value interface{}) error {
	cursor, err := r.Table(table).Get(id).Run(session)
	if err != nil {
		return err
	}
	cursor.One(value)
	cursor.Close()
	return nil
}

func GetCursor(name string) (*r.Cursor, error) {
	cursor, err := r.Table(name).Run(session)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func ListTable(name string, value interface{}) error {
	cursor, err := r.Table(name).Run(session)
	if err != nil {
		return err
	}
	cursor.All(value)
	cursor.Close()
	return nil
}

func FetchOne(name string, value interface{}) error {
	cursor, err := r.Table(name).Run(session)
	if err != nil {
		return err
	}
	cursor.One(value)
	cursor.Close()
	return nil
}
