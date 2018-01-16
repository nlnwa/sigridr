package controller

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
	"fmt"
)

type jobStore struct {
	*database.Rethink
}

func newJobStore(c Config) *jobStore {
	db := database.New()
	db.ConnectOpts.Address = c.DatabaseAddress
	db.ConnectOpts.Database = c.DatabaseName

	fmt.Println("DatabaseName", c.DatabaseName)
	return &jobStore{db}
}

func (js *jobStore) Connect() error {
	return js.Rethink.Connect()
}

func (js *jobStore) GetJobs() []types.Job {
	var jobs []types.Job

	err := js.ListTable("job", &jobs)
	if err != nil {
		log.WithError(err).Error("Getting jobs from database")
		return make([]types.Job, 0)
	}

	// Add tasks to job
	for i := range jobs {
		jobs[i].Seeds = js.getSeeds(&jobs[i])
	}

	return jobs
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
