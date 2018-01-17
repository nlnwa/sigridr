package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/worker"
	"github.com/nlnwa/sigridr/api"
)

var workerViper = viper.New()

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Sigridr worker service",
	Long:  `Sigridr worker`,
	Run: func(cmd *cobra.Command, args []string) {
		dbHost, dbPort, dbName := globalFlags()
		port := workerViper.GetInt("port")
		accessToken := workerViper.GetString("access-token")

		if err := work(dbHost, dbPort, dbName, port, accessToken); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	cmd := workerCmd

	rootCmd.AddCommand(cmd)

	cmd.Flags().Int("port", 10001, "Listening port")
	cmd.Flags().String("access-token", "", "Twitter access token")

	workerViper.BindPFlag("port", cmd.Flags().Lookup("port"))
	workerViper.BindPFlag("access-token", cmd.Flags().Lookup("access-token"))
}

func work(dbHost string, dbPort int, dbName string, port int, accessToken string) error {
	apiConfig := worker.Config{
		AccessToken:     accessToken,
		DatabaseAddress: fmt.Sprintf("%s:%d", dbHost, dbPort),
		DatabaseName:    dbName,
	}

	var grpcOpts []grpc.ServerOption

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listening on %s failed", port)
	} else {
		log.WithField("port", port).Infoln("API server listening")
	}
	server := grpc.NewServer(grpcOpts...)
	api.RegisterWorkerServer(server, worker.NewApi(apiConfig))

	return server.Serve(listener)
}
