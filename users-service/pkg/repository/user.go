package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"sync"
	"time"
)


type UserRepoInterface interface {
	CreateUser(ctx context.Context,key string,  value []byte)error
	GetUser(ctx context.Context, key string)([]byte, error)
	UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte)error
	UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration)error
}

type UserRepo struct {
	redisClient databases.Cache
	mu *sync.Mutex
}

func NewUserRepo(redisClient databases.Cache)UserRepoInterface{
	return &UserRepo{redisClient: redisClient}
}

func (userRepo *UserRepo) CreateUser(ctx context.Context,key string,  value []byte)error{
	_, err := userRepo.redisClient.Set(ctx, key, value)
	if err != nil{
		return err
	}
	return nil
}


func (userRepo *UserRepo) GetUser(ctx context.Context, key string)([]byte, error){
	result, err := userRepo.redisClient.Get(ctx, key)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (userRepo *UserRepo) UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration)error{
	lockKey := fmt.Sprintf("lock."+ key)
	timeout := time.Minute*5
	res, err := userRepo.redisClient.SetNX(ctx, lockKey, nil, timeout)
	if err != nil{
		fmt.Println(err)
	}
	if res.Val() == false{
		ttl, _ :=  userRepo.redisClient.GetTTL(ctx, lockKey)
		return errors.New("Update after"+ttl.String())
	}

	pipe := userRepo.redisClient.Begin()
	_, err = userRepo.redisClient.PipelineSetTTL(pipe, ctx, key, value, expiration)
	if err != nil{
		// TODO: Add the retry if required here
		return err
	}
	_, err = userRepo.redisClient.PipelineDelete(pipe, ctx, lockKey)
	if err != nil{
		// TODO: Add the retry if required here
		return err
	}
	outputs ,err := userRepo.redisClient.Commit(ctx, pipe)
	if err != nil{
	// TODO: Add the retry if required here
	return err
	}

	for _, output := range outputs{
		if output.Err() != nil{

		}
	}

	return nil
}

func (userRepo *UserRepo) UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte)error{
	err := userRepo.redisClient.Watch(ctx,
		userRepo.redisClient.GetTransactionFunc(ctx, key, value, 0),
		key)
	if err != nil{
		if err.Error() == helpers.RedisNil{
			return err
		}
		// TODO: Add retry mechanism here
		return err
	}
	return nil
}
