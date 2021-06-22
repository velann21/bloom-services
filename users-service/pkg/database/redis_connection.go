package database

import (
	"context"
	"github.com/velann21/bloom-services/common-lib/databases"
	"sync"
)

type RedisConnection struct {
	Client *databases.Redis
}

func (connection *RedisConnection) NewRedisConnection(ctx context.Context, address string, password string) error {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	if connection.Client == nil {
		redisConn, err := databases.NewRedis(ctx, address, password)
		if err != nil {
			return err
		}
		connection.Client = redisConn
		return nil
	}
	return nil
}
