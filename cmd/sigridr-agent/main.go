package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	workerAddress *string = flag.String("worker-address", "localhost:10001", "worker service address")
	port          int
	wg            sync.WaitGroup
)

func init() {
	log.SetLevel(log.DebugLevel)

	flag.IntVar(&port, "port", 10000, "gRPC server listening port")
	viper.BindPFlag("port", flag.Lookup("port"))

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	port = viper.GetInt("port")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	new(agent).register(server)


	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal()
	} else {
		log.WithField("port", port).Debugln("Listening")
	}

	wg.Add(1)
	go queueWorker(ctx, &wg)

	server.Serve(listener)

	wg.Wait()
}
