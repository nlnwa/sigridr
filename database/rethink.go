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

package database

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

type Rethink struct {
	session     *r.Session
	ConnectOpts *r.ConnectOpts
}

func New() *Rethink {
	rethink := new(Rethink)
	rethink.ConnectOpts = DefaultOptions()
	return rethink
}

func init() {
	r.SetTags("gorethink", "json", "url")
}

func (db *Rethink) Connect() error {
	var err error
	db.session, err = r.Connect(*db.ConnectOpts)
	return err
}

func (db *Rethink) Disconnect() {
	db.session.Close()
}

func (db *Rethink) DropDatabase(name string) error {
	_, err := r.DBDrop(name).RunWrite(db.session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) CreateDatabase(name string) error {
	_, err := r.DBCreate(name).RunWrite(db.session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) DropTable(name string) error {
	_, err := r.TableDrop(name).RunWrite(db.session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) CreateTable(name string) error {
	_, err := r.TableCreate(name).RunWrite(db.session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) Insert(table string, document interface{}) (string, error) {
	res, err := r.Table(table).Insert(document).RunWrite(db.session)
	if err != nil {
		return "", err
	}
	id := ""
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return id, nil
}

func (db *Rethink) Update(table string, id string, value interface{}) error {
	_, err := r.Table(table).Get(id).Update(value).RunWrite(db.session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) Delete(table string, id string) error {
	_, err := r.Table(table).Get(id).Delete().Run(db.session)
	return err
}

func (db *Rethink) Changes(name string) (*r.Cursor, error) {
	return r.Table(name).Changes().Run(db.session)
}

func (db *Rethink) Filter(table string, filterFunc interface{}) (*r.Cursor, error) {
	return r.Table(table).Filter(filterFunc).Run(db.session)
}

func (db *Rethink) Get(table string, id string, value interface{}) error {
	cursor, err := r.Table(table).Get(id).Run(db.session)
	if err != nil {
		return err
	}
	cursor.One(value)
	return nil
}

func (db *Rethink) FetchOne(table string, value interface{}) error {
	cursor, err := r.Table(table).Run(db.session)
	if err != nil {
		return err
	}
	cursor.One(value)
	return nil
}

func (db *Rethink) GetCursor(name string) (*r.Cursor, error) {
	cursor, err := r.Table(name).Run(db.session)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (db *Rethink) ListTable(name string, value interface{}) error {
	cursor, err := r.Table(name).Run(db.session)
	if err != nil {
		return err
	}
	cursor.All(value)
	cursor.Close()
	return nil
}
