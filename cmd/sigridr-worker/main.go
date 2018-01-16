package main

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/worker"
	"github.com/nlnwa/sigridr/api"
)

var (
	dbAddress   string
	dbName      string
	port        int
	accessToken string
	debug       bool
)

func init() {
	flag.IntVar(&port, "port", 10001, "gRPC server listening port")
	flag.StringVar(&accessToken, "access-token", "", "Twitter access token")
	flag.StringVar(&dbName, "db-name", "sigridr", "Database name")
	flag.StringVar(&dbAddress, "db", "", "Database host and port")
	flag.BoolVar(&debug, "debug", false, "Enable debug")

	viper.BindPFlag("db", flag.Lookup("db"))
	viper.BindPFlag("db-name", flag.Lookup("db-name"))
	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("access-token", flag.Lookup("access-token"))

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	dbAddress = viper.GetString("db")
	port = viper.GetInt("port")
	accessToken = viper.GetString("access-token")
	dbName = viper.GetString("db-name")

	workerConfig := worker.Config{
		AccessToken:     accessToken,
		DatabaseAddress: dbAddress,
		DatabaseName:    dbName,
	}

	var grpcOpts []grpc.ServerOption

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Errorf("listening on %s failed", port)
		return
	} else {
		log.WithField("port", port).Infoln("Worker API server listening")
	}
	server := grpc.NewServer(grpcOpts...)
	api.RegisterWorkerServer(server, worker.NewApi(workerConfig))

	server.Serve(listener)
}
