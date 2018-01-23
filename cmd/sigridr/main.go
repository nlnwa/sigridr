package main

import (
	"os"
	"strings"

	log "github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/logfmt"
	"github.com/nlnwa/sigridr/version"
)

var debug bool

var logger = log.New()

var rootCmd = &cobra.Command{
	Use:   "sigridrctl",
	Short: "Twitter API client",
	Long:  `Twitter API client`,
}

func init() {
	cobra.OnInitialize(func() {
		initViper(viper.GetViper())
	})

	rootCmd.PersistentFlags().String("db-host", "localhost", "database hostname")
	rootCmd.PersistentFlags().Int("db-port", 28015, "database port")
	rootCmd.PersistentFlags().String("db-name", "sigridr", "database name")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug")

	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
}

func initViper(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
}

func globalFlags() (dbHost string, dbPort int, dbName string) {
	dbHost = viper.GetString("db-host")
	dbPort = viper.GetInt("db-port")
	dbName = viper.GetString("db-name")
	return
}

func main() {
	logHandler := log.CallerFuncHandler(log.StreamHandler(os.Stdout, logfmt.LogbackFormat()))
	if debug {
		logger.SetHandler(log.CallerStackHandler("%+v", logHandler))
	} else {
		logger.SetHandler(log.LvlFilterHandler(log.LvlInfo, logHandler))
	}
	logger.Info(version.String())

	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
