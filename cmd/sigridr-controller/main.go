package main

import (
	"syscall"
	"strings"
	"time"

	cron "github.com/nlnwa/gocron"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/controller"
	"github.com/nlnwa/sigridr/signal"
)

var (
	agentAddress string
	dbAddress    string
	dbName       string
	debug        bool
)

func init() {
	flag.StringVar(&dbAddress, "db", "", "Database host and port")
	flag.StringVar(&dbName, "db-name", "sigridr", "Database name")
	flag.StringVar(&agentAddress, "agent-address", "localhost:10000", "Address to sigridr-agent")
	flag.BoolVar(&debug, "debug", false, "Enable debug")

	viper.BindPFlag("db", flag.Lookup("db"))
	viper.BindPFlag("db-name", flag.Lookup("db-name"))
	viper.BindPFlag("agent-address", flag.Lookup("agent-address"))

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	dbAddress = viper.GetString("db")
	agentAddress = viper.GetString("agent-address")
	dbName = viper.GetString("db-name")

	config := controller.Config{
		AgentAddress:    agentAddress,
		DatabaseName:    dbName,
		DatabaseAddress: dbAddress,
	}

	scheduler := cron.NewScheduler().Location(time.UTC).Interval(time.Minute)

	scheduler.AddCollector(controller.NewDbCollector(config))

	scheduler.Start()
	defer scheduler.Stop()

	log.WithField("interval", time.Minute).Infoln("Scheduler running")

	<-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
}
