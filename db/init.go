package db

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

func CreateDb(name string) {
	result, err := r.DBCreate(name).RunWrite(session)
	if err != nil {
		log.Error(err)
	}
	if n := result.DBsCreated; n > 0 {
		log.Debug("Database created: ", name)
	}
}

func CreateTable(name string) {
	result, err := r.TableCreate(name).RunWrite(session)
	if err != nil {
		log.Error(err)
	}
	if n := result.TablesCreated; n > 0 {
		log.Debug("Table created: ", name)
	}
}

func Insert(table string, data interface{}) {
	result, err := r.Table(table).Insert(data).RunWrite(session)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln(result)
}
