package main

import (
	"context"
	"fmt"
	"net"
	"syscall"

	log "github.com/sirupsen/logrus"
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
		dbHost, dbPort, dbName := globalFlags()

		port := agentViper.GetInt("port")
		workerPort := agentViper.GetInt("worker-port")
		workerHost := agentViper.GetString("worker-host")

		if err := act(port, workerHost, workerPort, dbHost, dbPort, dbName); err != nil {
			log.WithError(err).Fatal()
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

func act(port int, workerHost string, workerPort int, dbHost string, dbPort int, dbName string) error {
	config := agent.Config{
		WorkerAddress:   fmt.Sprintf("%s:%d", workerHost, workerPort),
		DatabaseAddress: fmt.Sprintf("%s:%d", dbHost, dbPort),
		DatabaseName:    dbName,
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
				log.WithField("port", port).Infoln("Agent API server listening")
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
			log.WithError(err).Error()
		}
		// stop gRPC server and wait for it to finish
		server.GracefulStop()
		if err := <-errc; err != nil {
			log.WithError(err).Error()
		}
	}
	return nil
}
