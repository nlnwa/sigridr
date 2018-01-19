package controller

import (
	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/agent"
	"github.com/nlnwa/sigridr/types"
)

type dbCollector struct {
	store  *jobStore
	runner *jobRunner
}

func NewDbCollector(c Config) cron.Collector {
	return &dbCollector{newJobStore(c), newJobRunner(c)}
}

// Collect jobs and schedule them with the scheduler
func (c *dbCollector) GetJobs() []*cron.Job {
	jobs := make([]*cron.Job, 0)

	if err := c.store.connect(); err != nil {
		log.WithError(err).Errorln("Connecting to database")
		return jobs
	}
	defer c.store.disconnect()

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

		if log.GetLevel() == log.DebugLevel {
			log.WithFields(log.Fields{
				"cron": job.CronExpression,
				"name": job.Meta.Name,
			}).Debugln("Collect job")
		}

		cronJob.AddTask(c.runner.execute, &job)
		jobs = append(jobs, cronJob)
	}

	return jobs
}

type jobRunner struct {
	store       *jobStore
	agentClient *agent.Client
}

func newJobRunner(c Config) *jobRunner {
	return &jobRunner{newJobStore(c), agent.NewApiClient(c.AgentAddress)}
}

func (jr *jobRunner) execute(job *types.Job) {
	log.WithField("name", job.Meta.Name).Infoln("Running job")

	// connect db
	if err := jr.store.connect(); err != nil {
		log.WithError(err).Errorln("connecting to database")
		return
	}
	// get seeds of job
	seeds := jr.store.getSeeds(job)
	jr.store.disconnect()

	// connect agent
	if len(seeds) > 0 {
		jr.agentClient.Dial()
		defer jr.agentClient.Hangup()
	}

	// dispatch seeds to agent
	for i := range seeds {
		if seeds[i].Meta.Name == "" {
			continue
		}
		log.WithFields(log.Fields{
			"handle": seeds[i].Meta.Description,
			"query":  seeds[i].Meta.Name,
		}).Infoln("Dispatch")
		jr.agentClient.Do(job, &seeds[i])
	}
}
