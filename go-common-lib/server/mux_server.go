package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Server(ctx context.Context, router *mux.Router, port string, close chan bool) {
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("Failed to start server")
		}
	}()

	logrus.Info("Server Started")
	<-ctx.Done()
	logrus.Info("Server Stopped")
	close <- true

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	err := srv.Shutdown(ctxShutDown)
	if err != nil {
		if err == http.ErrServerClosed {
			err = nil
		}
		logrus.Fatal("server Shutdown Failed:%+s", err)
	}
	logrus.Info("server exited properly")
	return
}

func NewMux() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	return r
}

type Router struct {
	Router *mux.Router
}

func NewRouter(routes *mux.Router) *Router {
	return &Router{Router: routes}
}
