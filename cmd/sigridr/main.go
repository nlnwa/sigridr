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

package main

import (
	"os"
	"strings"

	log "github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/pkg/logfmt"
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
