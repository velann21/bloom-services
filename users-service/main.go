package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-users-service/app"
	"github.com/velann21/bloom-users-service/pkg/routes"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := mux.NewRouter().StrictSlash(false)
	usersRoutes := r.PathPrefix("users/api/v1").Subrouter()
	routes.Routes(usersRoutes)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		osSignal := <-c
		logrus.Info("system call:%+v", osSignal)
		cancel()
	}()
	app.Server(ctx, r)
}
