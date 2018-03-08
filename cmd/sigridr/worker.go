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
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/worker"
)

var workerViper = viper.New()

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Sigridr worker service",
	Long:  `Sigridr worker`,
	Run: func(cmd *cobra.Command, args []string) {
		dbHost, dbPort, dbName, dbUser, dbPassword := globalFlags()
		port := workerViper.GetInt("port")
		accessToken := workerViper.GetString("access-token")

		if err := work(port, accessToken, dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		initViper(workerViper)
	})

	cmd := workerCmd

	rootCmd.AddCommand(cmd)

	cmd.Flags().Int("port", 10001, "listening port")
	cmd.Flags().String("access-token", "", "twitter access token")

	workerViper.BindPFlag("port", cmd.Flags().Lookup("port"))
	workerViper.BindPFlag("access-token", cmd.Flags().Lookup("access-token"))
}

func work(port int, accessToken string, dbHost string, dbPort int, dbName string, dbUser string, dbPassword string) error {
	apiConfig := worker.Config{
		AccessToken:      accessToken,
		DatabaseHost:     dbHost,
		DatabasePort:     dbPort,
		DatabaseName:     dbName,
		DatabaseUser:     dbUser,
		DatabasePassword: dbPassword,
		Logger:           logger,
	}

	var grpcOpts []grpc.ServerOption

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listening on %d failed", port)
	} else {
		logger.Info("API server listening", "port", port)
	}
	server := grpc.NewServer(grpcOpts...)
	api.RegisterWorkerServer(server, worker.NewApi(apiConfig))

	return server.Serve(listener)
}
