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

// Alias gorethink connect options
type Options = r.ConnectOpts

var session *r.Session

func Connect(opts r.ConnectOpts) {
	var err error

	// https://godoc.org/github.com/GoRethink/gorethink#ConnectOpts
	session, err = r.Connect(opts)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.WithFields(log.Fields{
			"address": opts.Address,
		}).Debugln("Database session: connected")
	}
}

func Disconnect() {
	if session.IsConnected() {
		session.Close()
		log.Debugln("Database session: closed")
	}
}
