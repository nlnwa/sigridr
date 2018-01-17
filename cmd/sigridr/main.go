// Copyright © 2017 National Library of Norway
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
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "sigridrctl",
	Short: "Twitter API client",
	Long:  `Twitter API client`,
}

func init() {
	cobra.OnInitialize(func() {
		initViper(viper.GetViper())
	})

	rootCmd.PersistentFlags().String("db-host", "localhost", "Database hostname")
	rootCmd.PersistentFlags().Int("db-port", 28015, "Database port")
	rootCmd.PersistentFlags().String("db-name", "sigridr", "Database name")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug")

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
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal()
	}
}
