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

package database

import (
	"fmt"

	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v3"
)

type Rethink struct {
	connectOpts r.ConnectOpts
	*r.Session
}

type ConnectOption func(*r.ConnectOpts)

func WithName(name string) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Database = name
	}
}

func WithAddress(host string, port int) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Address = fmt.Sprintf("%s:%d", host, port)
	}
}

func WithCredentials(username string, password string) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Username = username
		opts.Password = password
	}
}

func New(options ...ConnectOption) *Rethink {
	var db Rethink
	for _, option := range options {
		option(&db.connectOpts)
	}
	return &db
}

func (db *Rethink) Connect() error {
	var err error
	if db.Session, err = r.Connect(db.connectOpts); err != nil {
		return errors.Wrap(err, "failed to connect to database")
	} else {
		return nil
	}
}

func (db *Rethink) Disconnect() error {
	if err := db.Session.Close(); err != nil {
		return errors.Wrap(err, "failed to disconnect from database")
	} else {
		return nil
	}
}

func (db *Rethink) SetTags(tags ...string) {
	r.SetTags(tags...)
}

func (db *Rethink) DropDatabase(name string) error {
	if _, err := r.DBDrop(name).RunWrite(db.Session); err != nil {
		return errors.Wrapf(err, "failed to drop database: %s", name)
	} else {
		return nil
	}
}

func (db *Rethink) CreateDatabase(name string) error {
	if _, err := r.DBCreate(name).RunWrite(db.Session); err != nil {
		return errors.Wrapf(err, "failed to create database: %s", name)
	} else {
		return nil
	}
}

func (db *Rethink) DropTable(name string) error {
	if _, err := r.TableDrop(name).RunWrite(db.Session); err != nil {
		return errors.Wrapf(err, "failed to drop table: %s", name)
	}
	return nil
}

func (db *Rethink) CreateTable(name string) error {
	if _, err := r.TableCreate(name).RunWrite(db.Session); err != nil {
		return errors.Wrapf(err, "failed to create table: %s", name)
	} else {
		return nil
	}
}

func (db *Rethink) Insert(table string, document interface{}) (string, error) {
	res, err := r.Table(table).Insert(document).RunWrite(db.Session)
	if err != nil {
		return "", errors.Wrapf(err, "failed to insert document into table: %s", table)
	}
	id := ""
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return id, nil
}

func (db *Rethink) Update(table string, id string, value interface{}) error {
	if _, err := r.Table(table).Get(id).Update(value).RunWrite(db.Session); err != nil {
		return errors.Wrapf(err, "failed to update document (id: %s) in table: %s", id, table)
	} else {
		return nil
	}
}

func (db *Rethink) Delete(table string, id string) error {
	if _, err := r.Table(table).Get(id).Delete().Run(db.Session); err != nil {
		return errors.Wrapf(err, "failed to delete document (id: %s) in table: %s", id, table)
	} else {
		return nil
	}
}

func (db *Rethink) Changes(table string) (*r.Cursor, error) {
	if cursor, err := r.Table(table).Changes().Run(db.Session); err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to changes in table: %s", table)
	} else {
		return cursor, nil
	}
}

func (db *Rethink) Filter(table string, filterFunc interface{}) (*r.Cursor, error) {
	if cursor, err := r.Table(table).Filter(filterFunc).Run(db.Session); err != nil {
		return nil, errors.Wrapf(err, "failed to filter table: %s", table)
	} else {
		return cursor, nil
	}
}

func (db *Rethink) Get(table string, id string, value interface{}) error {
	if cursor, err := r.Table(table).Get(id).Run(db.Session); err != nil {
		return errors.Wrapf(err, "failed to get document (id: %s) in table: %s", id, table)
	} else {
		cursor.One(value)
		return nil
	}
}

func (db *Rethink) FetchOne(table string, value interface{}) error {
	if cursor, err := r.Table(table).Run(db.Session); err != nil {
		return errors.Wrapf(err, "failed to fetch single row from table: %s", table)
	} else {
		cursor.One(value)
		return nil
	}
}

func (db *Rethink) GetCursor(table string) (*r.Cursor, error) {
	if cursor, err := r.Table(table).Run(db.Session); err != nil {
		return nil, errors.Wrapf(err, "failed to get curser to table: %s", table)
	} else {
		return cursor, nil
	}
}

func (db *Rethink) ListTable(name string, value interface{}) error {
	if cursor, err := r.Table(name).Run(db.Session); err != nil {
		return errors.Wrapf(err, "failed to list table: %s", name)
	} else {
		cursor.All(value)
		cursor.Close()
		return nil
	}
}
