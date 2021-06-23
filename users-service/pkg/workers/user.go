package workers

import (
	"context"
	"fmt"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/users-service/pkg/repository"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	"sync"
)

type Workers struct {
	service service.UserInterface
}

func NewWorker(redisConnection *databases.Redis)*Workers{
	userRepo := repository.NewUserRepo(redisConnection)
	userService := service.NewUserService(userRepo)
	return &Workers{service: userService}
}

func (Workers *Workers) ExpiredUserListener(eventCloser chan bool){
	eventStream := make(chan string)
	errorChan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(){
		defer wg.Done()
		Workers.service.UserExpiredEvent(context.Background(), eventStream, errorChan)
	}()

	go func(){
		defer wg.Done()
		for {
			select {
			case data, ok := <-eventStream:
				if !ok{
					eventStream = nil
				}
				fmt.Println(data)
			case err, ok := <-errorChan:
				if !ok{
					errorChan = nil
				}
				fmt.Println(err)
			case _,ok := <-eventCloser:
				if !ok{
					errorChan = nil
				}
				close(eventStream)
				close(errorChan)
			}
		}
	}()
	wg.Wait()
}


