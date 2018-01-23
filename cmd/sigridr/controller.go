package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	cron "github.com/nlnwa/gocron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/controller"
	"github.com/nlnwa/sigridr/signal"
)

var controllerViper = viper.New()

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Sigridr controller service",
	Long:  "Sigridr controller",
	Run: func(cmd *cobra.Command, args []string) {
		dbHost, dbPort, dbName := globalFlags()
		agentHost := controllerViper.GetString("agent-host")
		agentPort := controllerViper.GetInt("agent-port")

		if err := control(dbHost, dbPort, dbName, agentHost, agentPort); err != nil {
			logger.Error(err.Error())
			os.Exit(2)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		initViper(controllerViper)
	})

	cmd := controllerCmd

	rootCmd.AddCommand(cmd)

	cmd.Flags().String("agent-host", "localhost", "agent hostname")
	cmd.Flags().Int("agent-port", 10000, "agent port")

	controllerViper.BindPFlag("agent-host", cmd.Flags().Lookup("agent-host"))
	controllerViper.BindPFlag("agent-port", cmd.Flags().Lookup("agent-port"))
}

func control(dbHost string, dbPort int, dbName string, agentHost string, agentPort int) error {
	config := controller.Config{
		AgentAddress: fmt.Sprintf("%s:%d", agentHost, agentPort),
		DatabaseName: dbName,
		DatabaseHost: dbHost,
		DatabasePort: dbPort,
		Logger:       logger,
	}

	scheduler := cron.NewScheduler()

	scheduler.
		Location(time.UTC).
		Interval(time.Minute).
		AddCollector(controller.NewDbCollector(config))

	scheduler.Start()
	defer scheduler.Stop()

	logger.Info("Scheduler running", "interval", time.Minute)

	<-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return nil
}
