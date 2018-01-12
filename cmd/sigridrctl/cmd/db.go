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

package cmd

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/pkg/db"
	"github.com/nlnwa/sigridr/pkg/types"
)

var databaseAddress string

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database test command",
	Long:  `Database test command`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		opts := db.Options{Database: "sigridr"}
		u, err := url.Parse(viper.GetString("database-address"))
		if err != nil {
			log.WithError(err).Fatal()
		} else {
			opts.Address = u.Path
		}
		db.ConnectWithOptions(opts)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Disconnect()
	},
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().UTC()

		db.DropDatabase("sigridr")
		db.CreateDatabase("sigridr")
		db.Use("sigridr")
		db.CreateTable("results")
		db.CreateTable("jobs")
		db.CreateTable("entities")
		db.CreateTable("seeds")
		db.CreateTable("seed_queue")
		db.CreateTable("search_parameters")

		jobMeta := &types.Meta{
			Name:           "Default",
			Description:    "Default job",
			CreatedBy:      "anonymous",
			Created:        now,
			LastModified:   now,
			LastModifiedBy: "anonymous",
		}
		job := &types.Job{
			Id:             uuid.New().String(),
			CronExpression: "1 * * * *",
			ValidFrom:      time.Unix(0, 0).UTC(),
			ValidTo:        time.Date(2018, 12, 22, 14, 3, 0, 0, time.Now().Location()),
			Meta:           jobMeta,
		}
		db.Insert("jobs", job)
	},
}

// createCmd represents the db create subcommand
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
			db.CreateDatabase(arg)
		case "table":
			db.CreateTable(arg)
		default:
			log.Println("No op", command)
		}
	},
}

func init() {
	RootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(createCmd)

	dbCmd.PersistentFlags().StringVarP(&databaseAddress, "database-address", "d", "", "Address to the Database service")
	viper.BindPFlag("database-address", dbCmd.PersistentFlags().Lookup("database-address"))
}
