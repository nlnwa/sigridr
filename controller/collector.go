package controller

import (
	cron "github.com/nlnwa/gocron"

	"github.com/nlnwa/sigridr/agent"
	"github.com/nlnwa/sigridr/log"
	"github.com/nlnwa/sigridr/types"
)

type dbCollector struct {
	store  *jobStore
	runner *jobRunner
	log.Logger
}

func NewDbCollector(c Config) cron.Collector {
	return &dbCollector{newJobStore(c), newJobRunner(c), c.Logger}
}

// Collect jobs and schedule them with the scheduler
func (c *dbCollector) GetJobs() []*cron.Job {
	cronJobs := make([]*cron.Job, 0)

	if err := c.store.connect(); err != nil {
		c.Error("failed to connect to database", "error", err)
		return cronJobs
	}
	defer c.store.disconnect()

	jobs, err := c.store.getJobs()
	if err != nil {
		c.Info("failed getting jobs from store", "error", err)
		return cronJobs
	}
	for _, job := range jobs {
		if job.Disabled {
			c.Info("Job disabled", "job", job.Meta.Name)
			continue
		}
		if !job.IsValid() {
			c.Info("Job not valid", "name", job.Meta.Name, "validFrom", job.ValidFrom.String(), "validTo", job.ValidTo.String())
			continue
		}

		cronJob, err := cron.NewCronJob(job.CronExpression)
		if err != nil {
			c.Error("failed to create new cron job", "error", err)
			continue
		}
		c.Debug("Collect job",
			"cron", job.CronExpression,
			"name", job.Meta.Name)

		cronJob.AddTask(c.runner.execute, &job)
		cronJobs = append(cronJobs, cronJob)
	}

	return cronJobs
}

type jobRunner struct {
	store       *jobStore
	agentClient *agent.Client
	log.Logger
}

func newJobRunner(c Config) *jobRunner {
	return &jobRunner{newJobStore(c), agent.NewApiClient(c.AgentAddress), c.Logger}
}

func (jr *jobRunner) execute(job *types.Job) {
	jr.Info("Running job", "name", job.Meta.Name)

	// connect db
	if err := jr.store.connect(); err != nil {
		jr.Info("failed connecting to database")
		return
	}
	// get seeds of job
	seeds, err := jr.store.getSeeds(job)
	defer jr.store.disconnect()
	if err != nil {
		jr.Error("failed getting seeds from store", "error", err.Error())
		return
	}
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
		if err := jr.agentClient.Do(job, &seeds[i]); err != nil {
			jr.Error("failed to call agent client method: Do")
		}
	}
}
