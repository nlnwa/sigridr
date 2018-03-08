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
	"context"
	"fmt"
	"net"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/agent"
	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/signal"
)

var agentViper = viper.New()

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Sigridr agent service",
	Long:  "Sigridr agent",
	Run: func(cmd *cobra.Command, args []string) {
		dbHost, dbPort, dbName, dbUser, dbPassword := globalFlags()

		port := agentViper.GetInt("port")
		workerPort := agentViper.GetInt("worker-port")
		workerHost := agentViper.GetString("worker-host")

		if err := act(port, workerHost, workerPort, dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		initViper(agentViper)
	})

	cmd := agentCmd

	rootCmd.AddCommand(cmd)

	cmd.Flags().Int("port", 10000, "server listening port")
	cmd.Flags().String("worker-host", "localhost", "worker hostname")
	cmd.Flags().Int("worker-port", 10001, "worker port")

	agentViper.BindPFlag("port", cmd.Flags().Lookup("port"))
	agentViper.BindPFlag("worker-host", cmd.Flags().Lookup("worker-host"))
	agentViper.BindPFlag("worker-port", cmd.Flags().Lookup("worker-port"))
}

func act(port int, workerHost string, workerPort int, dbHost string, dbPort int, dbName string, dbUser string, dbPassword string) error {
	config := agent.Config{
		WorkerAddress:    fmt.Sprintf("%s:%d", workerHost, workerPort),
		DatabaseHost:     dbHost,
		DatabasePort:     dbPort,
		DatabaseName:     dbName,
		DatabaseUser:     dbUser,
		DatabasePassword: dbPassword,
		Logger:           logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var grpcOpts []grpc.ServerOption
	var server *grpc.Server
	errc := make(chan error, 2)

	go func() {
		errc <- func() error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return fmt.Errorf("listening on %d failed", port)
			} else {
				logger.Info("API server listening", "port", port)
			}
			server = grpc.NewServer(grpcOpts...)
			api.RegisterAgentServer(server, agent.NewApi(config))
			return server.Serve(listener)
		}()
	}()

	go func() {
		errc <- agent.NewQueueWorker(config).Run(ctx)
	}()

	select {
	case err := <-errc:
		return err
	case <-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT):
		cancel()
		// wait for QueueWorker to finish what it is doing
		if err := <-errc; err != nil {
			logger.Error(err.Error())
		}
		// stop gRPC server and wait for it to finish
		server.GracefulStop()
		if err := <-errc; err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}
