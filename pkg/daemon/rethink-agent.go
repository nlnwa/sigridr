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

package daemon

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/nlnwa/sigridr/pkg/db"
)

func Run() {
	log.Println("Daemon is doing it's job")
}

func rethinkAgent(name string, group *sync.WaitGroup) {
	var value interface{}

	db.Connect()
	defer db.Disconnect()

	res, err := db.Changes(name)
	if err != nil {
		log.WithError(err).Fatal("Tried to get changefeed")
	} else {
		log.WithField("table", name).Println("Database agent is listening on changes")
	}

	for res.Next(&value) {
		log.WithField("value", value).Println("Got new result")
	}

	log.Println("Database agent is done listening")

	group.Done()
}


type Work struct {
	Number int
}

/*
func start(cmd *cobra.Command, args []string) {
	nrOfWorkers := 5
	waitGroup := new(sync.WaitGroup)
	work := make(chan Work)

	go databaseAgent("results")

	waitGroup.Add(nrOfWorkers)
	for i := 0; i < nrOfWorkers; i++ {
		go func(work chan Work, group *sync.WaitGroup, index int) {
			defer group.Done()
			log.WithField("number", index).Println("Worker running")
			x := <-work
			log.WithField("work", x.Number).Println("Got something to do")
		}(work, waitGroup, i)
	}

	work <- Work{Number: 42}

	waitGroup.Wait()
	log.Println("Start is finisihed")
}
*/

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

type SearchRequest struct {
	query string
}

type SearchResponse struct {
	result string
}

func searchWorker(request <-chan SearchRequest, response chan<- SearchResponse, done chan struct{}, group *sync.WaitGroup) {
	for {
		select {
		// no value is actually sent over the done channel, closing it results in the channel returning it's types null value
		case <-done:
			break
		case req := <-request:
			// receive request
			log.WithField("request", req).Println("Search worker got request")

			// do work

			// send response
			response <- SearchResponse{result: "result"}
		}
	}
	group.Done()
}
