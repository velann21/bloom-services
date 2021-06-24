package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/database"
	"github.com/velann21/bloom-services/users-service/pkg/routes"
	"github.com/velann21/bloom-services/users-service/pkg/workers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	helpers.DetectAppMode(os.Args)

	ctx, cancel := context.WithCancel(context.Background())
	rc := &database.RedisConnection{}
	err := rc.NewRedisConnection(ctx, helpers.GetEnv(helpers.REDIS), "")
	if err != nil{
		logrus.WithError(err).Error("Something went wrong in redis connection")
		os.Exit(1)
	}

	closeEvent := make(chan bool)
	worker := workers.NewWorker(rc.Client)
	go worker.ExpiredUserListener(closeEvent)

	muxRoutes := server.NewMux()
	indexRoutes := muxRoutes.PathPrefix("/users").Subrouter()
	usersRoutes := muxRoutes.PathPrefix("/users/api/v1").Subrouter()

	routes.Routes(server.NewRouter(usersRoutes), rc.Client)
	routes.IndexRoutes(server.NewRouter(indexRoutes))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	go func() {
		osSignal := <-c
		logrus.Info(osSignal)
		cancel()
	}()

	server.Server(ctx, muxRoutes, ":5000", closeEvent)
}

