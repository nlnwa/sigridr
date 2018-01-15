package store

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/types"
)

type JobStore struct {
	*db.Database
}

func New() *JobStore {
	return &JobStore{db.New()}
}

func (js *JobStore) GetJobs() []types.Job {
	var jobs []types.Job

	err := js.ListTable("jobs", &jobs)
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

func (js *JobStore) getSeeds(job *types.Job) []types.Seed {
	var seeds []types.Seed

	cursor, err := js.Filter("seeds", func(seed r.Term) r.Term {
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
