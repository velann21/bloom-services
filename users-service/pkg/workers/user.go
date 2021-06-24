package workers

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/users-service/pkg/repository"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	"sync"
)

type Workers struct {
	service service.UserInterface
}

func NewWorker(redisConnection *databases.Redis) *Workers {
	userRepo := repository.NewUserRepo(redisConnection)
	userService := service.NewUserService(userRepo)
	return &Workers{service: userService}
}

func (Workers Workers) ExpiredUserListener(eventCloser chan bool) {
	eventStream := make(chan string)
	errorChan := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		Workers.service.UserExpiredEvent(context.Background(), eventStream, errorChan)
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-eventStream:
				if !ok {
					logrus.Info("Channel eventStream closed")
					eventStream = nil
				}
				logrus.Info(data)
			case err, ok := <-errorChan:
				if !ok {
					logrus.Info("Channel errorChan closed")
					errorChan = nil
				}
				logrus.Info(err)
			case _, ok := <-eventCloser:
				if !ok {
					logrus.Info("Channel eventCloser closed")
					eventCloser = nil
				}
				logrus.Info("eventCloser")
				close(eventStream)
				close(errorChan)
				close(eventCloser)
			}
		}
	}()

	wg.Wait()
}
