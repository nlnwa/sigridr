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
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
)

var (
	databaseAddress string
	db              *database.Rethink
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database test command",
	Long:  `Database test command`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		db = database.New()
		db.ConnectOpts.Database = "sigridr"
		db.ConnectOpts.Address = viper.GetString("database-address")
		db.Connect()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Disconnect()
	},

	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().UTC()

		db.DropDatabase("sigridr")
		db.CreateDatabase("sigridr")
		db.CreateTable("result")
		db.CreateTable("job")
		db.CreateTable("entity")
		db.CreateTable("seed")
		db.CreateTable("queue")
		db.CreateTable("parameter")

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
		db.Insert("job", job)
	},
}

func init() {
	RootCmd.AddCommand(dbCmd)

	dbCmd.PersistentFlags().StringVarP(&databaseAddress, "database-address", "d", "localhost:28015", "Database address")
	viper.BindPFlag("database-address", dbCmd.PersistentFlags().Lookup("database-address"))
}
