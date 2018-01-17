package controller

import (
	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"
)

type dbCollector struct {
	store       *jobStore
	agentClient *agentClient
}

func NewDbCollector(c Config) cron.Collector {
	return &dbCollector{newJobStore(c), newAgentClient(c)}
}

// Collect jobs and schedule them with the scheduler
func (c *dbCollector) GetJobs() []*cron.Job {
	nrOfJobs := 0
	nrOfSeeds := 0

	jobs := make([]*cron.Job, 0)

	err := c.store.connect()
	defer c.store.disconnect()
	if err != nil {
		log.WithError(err).Errorln("Connecting to database")
	}

	for _, job := range c.store.getJobs() {
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
			cronJob.AddTask(c.agentClient.dispatch, &job, &job.Seeds[i])
			nrOfSeeds++
		}
		jobs = append(jobs, cronJob)
		nrOfJobs++
	}
	log.WithFields(log.Fields{
		"jobs":  nrOfJobs,
		"seeds": nrOfSeeds,
	}).Infoln("Job collector")
	return jobs
}
