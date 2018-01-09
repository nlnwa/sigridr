package main

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	port        int
	accessToken string
)

func init() {
	flag.IntVar(&port, "port", 10001, "gRPC server listening port")
	flag.StringVar(&accessToken, "access-token", "", "Access token")

	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("access-token", flag.Lookup("access-token"))

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()

	port = viper.GetInt("port")
	accessToken = viper.GetString("access-token")

	opts := make([]grpc.ServerOption, 0)
	server := grpc.NewServer(opts...)

	new(worker).register(server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal()
	} else {
		log.WithField("port", port).Infoln("Listening")
	}

	server.Serve(listener)
}
