package main

import (
	"context"
	"net"
	"strings"
	"fmt"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/agent"
	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/signal"
	"syscall"
)

var (
	workerAddress string
	port          int
)

func init() {
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&workerAddress, "worker-address", "localhost:10001", "worker service address")
	flag.IntVar(&port, "port", 10000, "gRPC server listening port")

	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("worker-address", flag.Lookup("worker-address"))

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
	var server *grpc.Server
	errc := make(chan error, 2)

	go func() {
		errc <- func() error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return fmt.Errorf("listening on %s failed", port)
			} else {
				log.WithField("port", port).Infoln("Agent API server listening")
			}
			server = grpc.NewServer(grpcOpts...)
			api.RegisterAgentServer(server, agent.NewApi())
			return server.Serve(listener)
		}()
	}()

	go func() {
		errc <- agent.QueueWorker(ctx, workerConfig)
	}()

	select {
	case err := <-errc:
		log.WithError(err).Errorln()
	case <-signal.Receive(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT):
		cancel()
		// wait for QueueWorker to finish what it is doing
		if err := <-errc; err != nil  {
			log.WithError(err).Error()
		}
		server.GracefulStop()
		if err := <-errc; err != nil {
			log.WithError(err).Errorln()
		}
	}
}
