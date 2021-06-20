package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/database"
	"github.com/velann21/bloom-services/users-service/pkg/routes"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	rc := &database.RedisConnection{}
	err := rc.NewRedisConnection(ctx, "127.0.0.1:6379", "")
	if err != nil{
		logrus.WithError(err).Error("Something went wrong in redis connection")
		logrus.Fatal(err.Error())
	}
	muxRoutes := server.NewMux()
	usersRoutes := muxRoutes.PathPrefix("/users/api/v1").Subrouter()
	routes.Routes(server.NewRouter(usersRoutes), rc.Client)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	go func() {
		osSignal := <-c
		logrus.Info("system call:%+v", osSignal)
		cancel()
	}()
	server.Server(ctx, muxRoutes)
}
