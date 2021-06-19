package repository

import (
	redis "github.com/go-redis/redis/v8"
)

type UserRepoInterface interface {

}

type UserRepo struct {
	Client *redis.Client
}


