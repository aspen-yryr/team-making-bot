package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"

	matchpb "github.com/aspen-yryr/team-making-bot/proto/match"
	"github.com/aspen-yryr/team-making-bot/service/match"
	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const defaultEnvFile = "./env/.env"

func main() {
	e := flag.String("env_file", defaultEnvFile, "env variables definition file")
	flag.Parse()

	err := godotenv.Load(*e)
	if err != nil {
		glog.Errorf("Cannot load env file: %v", err)
		return
	}

	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	mt := match.NewMatchService()
	s := grpc.NewServer()
	matchpb.RegisterMatchSvcServer(s, mt)
	glog.V(1).Infof("server listening at %v", listener.Addr())

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-stopServer
		glog.Info("Stopping Server")
		s.GracefulStop()
	}()

	if err := s.Serve(listener); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
