package main

import (
	"syscall"
	"time"

	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/controller"
	"github.com/nlnwa/sigridr/pkg/signal"
)

var (
	address = "localhost:10000"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	scheduler := cron.NewScheduler().Location(time.UTC).Interval(time.Minute)

	scheduler.AddCollector(controller.NewDbCollector())

	scheduler.Start()
	defer scheduler.Stop()

	<-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)
}
