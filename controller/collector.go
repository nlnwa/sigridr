package controller

import (
	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/controller/store"
)

type dbCollector struct {
	store *store.JobStore
}

func NewDbCollector() cron.Collector {
	return &dbCollector{store.New()}
}

// Collect jobs and schedule them with the scheduler
func (c *dbCollector) GetJobs() []*cron.Job {
	jobs := make([]*cron.Job, 0)

	err := c.store.Connect()
	defer c.store.Disconnect()
	if err != nil {
		log.WithError(err).Errorln("Connecting to database")
	}

	for _, job := range c.store.GetJobs() {
		if job.Disabled {
			log.WithField("job", job).Debugln("Job disabled")
			continue
		}
		if !job.IsValid() {
			log.WithField("job", job).Debugln("Job not valid")
			continue
		}

		cronJob, err := cron.NewCronJob(job.CronExpression)
		if err != nil {
			log.WithError(err).Error()
			continue
		}

		// Add tasks to job
		for i := range job.Seeds {
			cronJob.AddTask(dispatch, &job, &job.Seeds[i])
		}
		jobs = append(jobs, cronJob)
	}
	return jobs
}
