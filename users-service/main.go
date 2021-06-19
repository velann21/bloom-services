package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/routes"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	muxRoutes := server.NewMux()
	usersRoutes := muxRoutes.PathPrefix("users/api/v1").Subrouter()
	routes.Routes(server.NewRouter(usersRoutes))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		osSignal := <-c
		logrus.Info("system call:%+v", osSignal)
		cancel()
	}()
	server.Server(ctx, muxRoutes)
}
