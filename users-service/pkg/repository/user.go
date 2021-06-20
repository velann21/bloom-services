package repository

import (
	"context"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/common-lib/helpers"
)


type UserRepoInterface interface {
	CreateUser(ctx context.Context,key string,  value []byte)error
	GetUser(ctx context.Context, key string)([]byte, error)
	UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte)error
	UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte)error
}

type UserRepo struct {
	redisClient databases.Cache
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

func (userRepo *UserRepo) UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte)error{

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