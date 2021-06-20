package repository

import (
	"context"
	"github.com/velann21/bloom-services/common-lib/databases"
)

const USERKEY = "User"

type UserRepoInterface interface {
	CreateUser(ctx context.Context,key string,  value []byte)error
	GetUser(ctx context.Context, key string)([]byte, error)
}

type UserRepo struct {
	redisClient *databases.Redis
}

func NewUserRepo(redisClient *databases.Redis)UserRepoInterface{
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



