package controller

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
)

type jobStore struct {
	*database.Rethink
	*database.ConnectOpts
}

func newJobStore(c Config) *jobStore {
	return &jobStore{
		Rethink: database.New(),
		ConnectOpts: &database.ConnectOpts{
			Address:  c.DatabaseAddress,
			Database: c.DatabaseName,
		},
	}
}

func (js *jobStore) connect() error {
	return js.Rethink.Connect(js.ConnectOpts)
}

func (js *jobStore) disconnect() {
	if err := js.Rethink.Disconnect(); err != nil {
		log.WithError(err).Errorln("disconnecting from database")
	}
}

func (js *jobStore) getJobs() []types.Job {
	var jobs []types.Job

	if err := js.ListTable("job", &jobs); err != nil {
		log.WithError(err).Error("Getting jobs from database")
		return make([]types.Job, 0)
	} else {
		return jobs
	}
}

func (js *jobStore) getSeeds(job *types.Job) []types.Seed {
	var seeds []types.Seed

	cursor, err := js.Filter("seed", func(seed r.Term) r.Term {
		return seed.Field("jobId").Contains(job.Id)
	})
	if err != nil {
		log.WithError(err).Errorln("Getting seeds with jobId from database")
	}
	err = cursor.All(&seeds)
	if err != nil {
		log.WithError(err).Errorln("Getting seeds with jobId from database")
	}
	return seeds
}
