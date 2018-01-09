package main

import (
	"syscall"
	"time"

	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/signal"
)

var (
	address = "localhost:10000"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	scheduler := cron.NewScheduler().Location(time.UTC).Interval(time.Nanosecond)

	scheduler.AddCollector(new(dbCollector))

	scheduler.Start()
	defer scheduler.Stop()

	<-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)
}

type dbCollector struct{}

// Collect jobs and schedule them with the scheduler
func (c *dbCollector) GetJobs() []*cron.Job {
	jobs := make([]*cron.Job, 0)

	for _, job := range db.GetJobs() {
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
