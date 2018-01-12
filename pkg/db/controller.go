package db

import (
	"github.com/nlnwa/sigridr/pkg/types"
	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

func GetJobs() []types.Job {
	var jobs []types.Job

	Connect()
	defer Disconnect()

	err := ListTable("jobs", &jobs)
	if err != nil {
		log.WithError(err).Error("Getting jobs from database")
		return make([]types.Job, 0)
	}

	// Add tasks to job
	for i := range jobs {
		jobs[i].Seeds = GetSeeds(&jobs[i])
	}

	return jobs
}

func GetSeeds(job *types.Job) []types.Seed {
	var seeds []types.Seed

	Connect()
	defer Disconnect()

	cursor, err := r.Table("seeds").Filter(func(seed r.Term) r.Term {
		return seed.Field("jobId").Contains(job.Id)
	}).Run(session)
	if err != nil {
		log.WithError(err).Errorln("Getting seeds with jobId from database")
	}

	err = cursor.All(&seeds)
	if err != nil {
		log.WithError(err).Errorln("Getting seeds with jobId from database")
	}
	return seeds
}
