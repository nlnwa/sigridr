package main

import (
	"syscall"
	"time"

	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/controller"
	"github.com/nlnwa/sigridr/signal"
)

var (
	address = "localhost:10000"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	config := controller.Config{
		Agent: address,
	}

	scheduler := cron.NewScheduler().Location(time.UTC).Interval(time.Minute)

	scheduler.AddCollector(controller.NewDbCollector(config))

	scheduler.Start()
	defer scheduler.Stop()

	log.WithField("interval", time.Minute).Infoln("Scheduler running")

	<-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
}
