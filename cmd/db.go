// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/nlnwa/sigridr/db"
	log "github.com/sirupsen/logrus"
)

var databaseAddress string

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database test command",
	Long:  `Database test command`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := db.Options{Database: "twitter"}
		u, err := url.Parse(viper.GetString("database-address"))
		if err != nil {
			log.WithError(err).Fatal()
		} else {
			opts.Address = u.Path
		}
		db.Connect(opts)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Disconnect()
	},
}

// createCmd represents the db init subcommand
var createCmd = &cobra.Command{
	Use:   "create db|table name",
	Short: "Initialize database",
	Long:  `Initialize database`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		arg := args[1]

		switch command {
		case "database":
			fallthrough
		case "db":
			db.CreateDb(arg)
		case "table":
			db.CreateTable(arg)
		default:
			log.Println("No op ", command)
		}
	},
}

func init() {
	RootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(createCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dbCmd.PersistentFlags().StringVarP(&databaseAddress, "database-address", "d", "", "Address to the Database service")
	viper.BindPFlag("database-address", dbCmd.PersistentFlags().Lookup("database-address"))
}
