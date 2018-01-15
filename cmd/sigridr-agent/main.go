package main

import (
	"context"
	"net"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/agent"
)

var (
	workerAddress *string = flag.String("worker-address", "localhost:10001", "worker service address")
	port          int
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
	workerAddress = viper.GetString("worker-address")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerConfig := agent.Config{
		Worker: workerAddress,
	}
	var grpcOpts []grpc.ServerOption

	errc := make(chan error, 2)

	go func() {
		errc <- func() error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return log.WithError(err).Errorf("listening on %s failed", port)
			}
			server := grpc.NewServer(grpcOpts...)
			api.RegisterAgentServer(server, agent.NewApi())
			err = server.Serve(listener)
			log.WithError(err).Error()
			return err
		}()
	}()

	go func() {
		errc <- agent.Worker(ctx, workerConfig)
	}()

	return <-errc
}
