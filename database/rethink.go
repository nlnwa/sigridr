package database

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

type ConnectOpts = r.ConnectOpts

type Rethink struct {
	*r.Session
}

func New() *Rethink {
	return new(Rethink)
}

func init() {
	r.SetTags("gorethink", "json", "url")
}

func (db *Rethink) Connect(opts *ConnectOpts) error {
	var err error
	db.Session, err = r.Connect(*opts)
	return err
}

func (db *Rethink) Disconnect() error {
	return db.Session.Close()
}

func (db *Rethink) DropDatabase(name string) error {
	_, err := r.DBDrop(name).RunWrite(db.Session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) CreateDatabase(name string) error {
	_, err := r.DBCreate(name).RunWrite(db.Session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) DropTable(name string) error {
	_, err := r.TableDrop(name).RunWrite(db.Session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) CreateTable(name string) error {
	_, err := r.TableCreate(name).RunWrite(db.Session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) Insert(table string, document interface{}) (string, error) {
	res, err := r.Table(table).Insert(document).RunWrite(db.Session)
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
	_, err := r.Table(table).Get(id).Update(value).RunWrite(db.Session)
	if err != nil {
		return err
	}
	return nil
}

func (db *Rethink) Delete(table string, id string) error {
	_, err := r.Table(table).Get(id).Delete().Run(db.Session)
	return err
}

func (db *Rethink) Changes(name string) (*r.Cursor, error) {
	return r.Table(name).Changes().Run(db.Session)
}

func (db *Rethink) Filter(table string, filterFunc interface{}) (*r.Cursor, error) {
	return r.Table(table).Filter(filterFunc).Run(db.Session)
}

func (db *Rethink) Get(table string, id string, value interface{}) error {
	cursor, err := r.Table(table).Get(id).Run(db.Session)
	if err != nil {
		return err
	}
	cursor.One(value)
	return nil
}

func (db *Rethink) FetchOne(table string, value interface{}) error {
	cursor, err := r.Table(table).Run(db.Session)
	if err != nil {
		return err
	}
	cursor.One(value)
	return nil
}

func (db *Rethink) GetCursor(name string) (*r.Cursor, error) {
	cursor, err := r.Table(name).Run(db.Session)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (db *Rethink) ListTable(name string, value interface{}) error {
	cursor, err := r.Table(name).Run(db.Session)
	if err != nil {
		return err
	}
	cursor.All(value)
	cursor.Close()
	return nil
}
