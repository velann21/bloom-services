package databases

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string)([]byte, error)
	Set(ctx context.Context, key string, value interface{})error
	SetWithTTL(ctx context.Context, key string, value interface{}, expiration time.Duration)error
}

type Redis struct {
	Client *redis.Client
}

func NewRedis(ctx context.Context, address string, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: 0,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{
		Client: client,
	}, nil
}

func (redis *Redis) Get(ctx context.Context, key string)([]byte, error){
	result := redis.Client.Get(ctx, key)
	if result.Err() != nil{
		return nil, result.Err()
	}
	res, err := result.Bytes()
	if err != nil{
		return nil, result.Err()
	}
	return res, nil
}

func (redis *Redis) Set(ctx context.Context, key string, value interface{})(string,error){
	result := redis.Client.Set(ctx, key, value, -1)
	if result.Err() != nil{
		return "",result.Err()
	}
	return result.Val(), nil
}

func (redis *Redis) SetWithTTL(ctx context.Context, key string, value interface{}, expiration time.Duration)error{
	result := redis.Client.Set(ctx, key, value, -1)
	if result.Err() != nil{
		return result.Err()
	}
	return nil
}

func (redis *Redis) Begin()redis.Pipeliner{
	pipe := redis.Client.TxPipeline()
	return pipe
}

func (redis *Redis) Commit(ctx context.Context, pipe redis.Pipeliner)([]redis.Cmder,error){
	results, err := pipe.Exec(ctx)
	if err != nil{
		return nil, err
	}
	return results, nil
}


