package databases

import (
	"context"
	redisV8 "github.com/go-redis/redis/v8"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}) (string, error)
	SetWithTTL(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Begin() redisV8.Pipeliner
	Commit(ctx context.Context, pipe redisV8.Pipeliner) ([]redisV8.Cmder, error)
	Watch(ctx context.Context, txFunc func(tx *redisV8.Tx) error, key string) error
	GetTransactionFunc(ctx context.Context, key string, value []byte, expiration time.Duration) func(tx *redisV8.Tx) error
	GetExpireTime(ctx context.Context, key string)(time.Duration, error)
	SetNX(ctx context.Context, key string, value []byte, expiration time.Duration)(*redisV8.BoolCmd,error)
	GetSet(ctx context.Context, key string, value []byte)([]byte, error)
	PipelineGet(pipe redisV8.Pipeliner, ctx context.Context, key string) ([]byte, error)
	PipelineSet(pipe redisV8.Pipeliner, ctx context.Context, key string, value interface{}) (string, error)
	PipelineGetSet(pipe redisV8.Pipeliner, ctx context.Context, key string, value []byte)([]byte, error)
	PipelineSetNX(pipe redisV8.Pipeliner, ctx context.Context, key string, value []byte, expiration time.Duration)(*redisV8.BoolCmd, error)
	GetTTL(ctx context.Context, key string)(time.Duration,error)
	PipelineGetTTL(pipe redisV8.Pipeliner, ctx context.Context, key string)(time.Duration,error)
	PipelineSetTTL(pipe redisV8.Pipeliner, ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	PipelineDelete(pipe redisV8.Pipeliner, ctx context.Context, key string)(*redisV8.IntCmd, error)
}

type Redis struct {
	Client *redisV8.Client
}

func NewRedis(ctx context.Context, address string, password string) (*Redis, error) {
	client := redisV8.NewClient(&redisV8.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{
		Client: client,
	}, nil
}

func (redis *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	result := redis.Client.Get(ctx, key)
	if result.Err() != nil {
		return nil, result.Err()
	}
	res, err := result.Bytes()
	if err != nil {
		return nil, result.Err()
	}
	return res, nil
}

func (redis *Redis) GetExpireTime(ctx context.Context, key string)(time.Duration, error){
	result := redis.Client.TTL(ctx, key)
	if result.Err() != nil {
		return 0, result.Err()
	}
	return result.Val(), nil
}

func (redis *Redis) Set(ctx context.Context, key string, value interface{}) (string, error) {
	result := redis.Client.Set(ctx, key, value, -1)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func (redis *Redis) SetWithTTL(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	result := redis.Client.Set(ctx, key, value, expiration)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (redis *Redis) Begin() redisV8.Pipeliner {
	pipe := redis.Client.TxPipeline()
	return pipe
}

func (redis *Redis) Commit(ctx context.Context, pipe redisV8.Pipeliner) ([]redisV8.Cmder, error) {
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (redis *Redis) Watch(ctx context.Context, txFunc func(tx *redisV8.Tx) error, key string) error {
	err := redis.Client.Watch(ctx, txFunc, key)
	if err != nil {
		return err
	}
	return nil
}

func (redis *Redis) GetTransactionFunc(ctx context.Context, key string, value []byte, expiration time.Duration) func(tx *redisV8.Tx) error {
	txf := func(tx *redisV8.Tx) error {
		_, err := tx.Get(ctx, key).Bytes()
		if err != nil{
			if err.Error() == helpers.RedisNil{
				return err
			}
		}
		_, err = tx.TxPipelined(ctx, func(pipe redisV8.Pipeliner) error {pipe.Set(ctx, key, value, expiration)
		return nil})
		if err != nil{
			return err
		}
		return nil
	}
	return txf
}

func (redis *Redis) SetNX(ctx context.Context, key string, value []byte, expiration time.Duration)(*redisV8.BoolCmd,error){
	boolCmd := redis.Client.SetNX(ctx, key, value, expiration)
	if boolCmd.Err() != nil{
		return boolCmd, boolCmd.Err()
	}
	return boolCmd, nil
}

func (redis *Redis) GetSet(ctx context.Context, key string, value []byte)([]byte, error){
	cmd := redis.Client.GetSet(ctx, key, value)
	if cmd.Err() != nil{
		return nil, cmd.Err()
	}
	byteData, err := cmd.Bytes()
	if err != nil{
		return nil, cmd.Err()
	}
	return byteData, nil
}

func (redis *Redis) PipelineGet(pipe redisV8.Pipeliner, ctx context.Context, key string) ([]byte, error) {
	result := pipe.Get(ctx, key)
	if result.Err() != nil {
		return nil, result.Err()
	}
	res, err := result.Bytes()
	if err != nil {
		return nil, result.Err()
	}
	return res, nil
}

func (redis *Redis)  PipelineSet(pipe redisV8.Pipeliner, ctx context.Context, key string, value interface{}) (string, error) {
	result := pipe.Set(ctx, key, value, -1)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func (redis *Redis)  PipelineSetTTL(pipe redisV8.Pipeliner, ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	result := pipe.Set(ctx, key, value, expiration)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func (redis *Redis) PipelineGetSet(pipe redisV8.Pipeliner, ctx context.Context, key string, value []byte)([]byte, error){
	cmd := pipe.GetSet(ctx, key, value)
	if cmd.Err() != nil{
		return nil, cmd.Err()
	}
	byteData, err := cmd.Bytes()
	if err != nil{
		return nil, cmd.Err()
	}
	return byteData, nil
}

func (redis *Redis) PipelineSetNX(pipe redisV8.Pipeliner, ctx context.Context, key string, value []byte, expiration time.Duration)(*redisV8.BoolCmd, error){
	boolCmd := pipe.SetNX(ctx, key, value, expiration)
	if boolCmd.Err() != nil{
		return nil, boolCmd.Err()
	}
	return boolCmd, nil
}

func (redis *Redis) GetTTL(ctx context.Context, key string)(time.Duration,error){
	output := redis.Client.TTL(ctx, key)
	if output.Err() != nil{
		return 0, output.Err()
	}
	return output.Val(), nil
}

func (redis *Redis) PipelineGetTTL(pipe redisV8.Pipeliner, ctx context.Context, key string)(time.Duration,error){
	output := pipe.TTL(ctx, key)
	if output.Err() != nil{
		return 0, output.Err()
	}
	return output.Val(), nil
}

func (redis *Redis) PipelineDelete(pipe redisV8.Pipeliner, ctx context.Context, key string)(*redisV8.IntCmd, error){
	deleteResult := pipe.Del(ctx, key)
	if deleteResult.Err() != nil{
		return nil, deleteResult.Err()
	}
	return deleteResult, nil
}