// Copyright 2018 National Library of Norway
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

package controller

import (
	cron "github.com/nlnwa/gocron"

	"github.com/nlnwa/pkg/log"
	"github.com/nlnwa/sigridr/agent"
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
	defer func() {
		_ = c.store.disconnect()
	}()

	jobs, err := c.store.getJobs()
	if err != nil {
		c.Error("failed getting jobs from store", "error", err)
		return cronJobs
	}
	for _, job := range jobs {
		if job.Disabled {
			c.Debug("Job disabled", "job", job.Meta.Name)
			continue
		}
		if !job.IsValid() {
			c.Debug("Job not valid", "name", job.Meta.Name, "validFrom", job.ValidFrom.String(), "validTo", job.ValidTo.String())
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
	defer func() {
		_ = jr.store.disconnect()
	}()
	if err != nil {
		jr.Error("failed getting seeds from store", "error", err.Error())
		return
	}
	// connect agent
	if len(seeds) > 0 {
		err = jr.agentClient.Dial()
		if err != nil {
			jr.Error("failed to dial agent", "error", err.Error())
			return
		}
		defer jr.agentClient.Hangup()
	}

	// dispatch seeds to agent
	for i := range seeds {
		if seeds[i].Meta.Name == "" {
			continue
		}
		if err := jr.agentClient.Do(job, &seeds[i]); err != nil {
			jr.Error("failed to call agent client method: Do", "err", err)
		}
	}
}
