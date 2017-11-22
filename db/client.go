package db

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

type Options = r.ConnectOpts

var session *r.Session

func Connect(opts r.ConnectOpts) {
	var err error

	// https://godoc.org/github.com/GoRethink/gorethink#ConnectOpts
	session, err = r.Connect(opts)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Debugf("Database session: connected to %+v", opts.Address)
	}
}

func Disconnect() {
	if session.IsConnected() {
		session.Close()
		log.Debugln("Database session: closed")
	}
}
