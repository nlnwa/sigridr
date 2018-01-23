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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
)

// dbCmd represents the db command
var initDbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Initialize database",
	Long:  `Initialize database`,
	Run: func(cmd *cobra.Command, args []string) {
		dbHost, dbPort, dbName := globalFlags()

		if err := initDb(dbHost, dbPort, dbName); err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(initDbCmd)
}

func initDb(dbHost string, dbPort int, dbName string) error {
	logger.Info("Initializing database", "dbHost", dbHost, "dbPort", dbPort, "dbName", dbName)

	db := database.New(database.WithAddress(dbHost, dbPort), database.WithName(dbName))

	if err := db.Connect(); err != nil {
		return err
	}
	defer db.Disconnect()

	now := time.Now().UTC()

	tables := []string{"result", "job", "entity", "seed", "queue", "parameter"}

	if err := db.CreateDatabase(dbName); err != nil {
		return err
	} else {
		logger.Info("Created database", "name", dbName)
	}

	for _, table := range tables {
		if err := db.CreateTable(table); err != nil {
			return err
		} else {
			logger.Info("Created table", "name", table)
		}
	}

	job := &types.Job{
		Id:             uuid.New().String(),
		CronExpression: "* * * * *",
		ValidFrom:      time.Unix(0, 0).UTC(),
		ValidTo:        time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC),
		Meta: &types.Meta{
			Name:           "Default",
			Description:    "Default job",
			CreatedBy:      "anonymous",
			Created:        now,
			LastModified:   now,
			LastModifiedBy: "anonymous",
		},
		Disabled: true,
	}

	if _, err := db.Insert("job", job); err != nil {
		return fmt.Errorf("inserting job %s: %v", job.Meta.Name, err)
	} else {
		logger.Info("Inserted job", "name", job.Meta.Name, "cron", job.CronExpression)
	}

	return nil
}
