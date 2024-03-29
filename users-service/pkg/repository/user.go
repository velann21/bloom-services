package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"time"
)

const (
	LockPrefix        = "lock."
	LockExpiry        = time.Minute * 5
	ReleaseLockScript = "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, key string, value []byte, expiration time.Duration) error
	GetUser(ctx context.Context, key string) ([]byte, error)
	UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration) error
	UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte, uid string, expiration time.Duration) error
	GetUserLock(ctx context.Context, key string) error
	DeleteUserLock(ctx context.Context, key string) error
	SubscribeForKeyExpireChannel(ctx context.Context, eventStream chan string, errChan chan error)
}

type UserRepo struct {
	redisClient databases.Cache
}

func NewUserRepo(redisClient databases.Cache) UserRepoInterface {
	return &UserRepo{redisClient: redisClient}
}

func (userRepo UserRepo) CreateUser(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	logrus.Debug("Inside the CreateUser Repository func")

	_, err := userRepo.redisClient.SetWithTTL(ctx, key, value, expiration)
	if err != nil {
		return err
	}

	logrus.Debug("Complete the CreateUser Repository func")
	return nil
}

func (userRepo UserRepo) GetUser(ctx context.Context, key string) ([]byte, error) {
	logrus.Debug("Inside the GetUser Repository func")

	result, err := userRepo.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	logrus.Debug("Completed the GetUser Repository func")
	return result, nil
}

func (userRepo UserRepo) UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte, uid string, expiration time.Duration) error {
	logrus.Debug("Inside the UpdateUserWithPessimisticLocking Repository func")

	lockKey := fmt.Sprintf(LockPrefix + key)

	res, err := userRepo.redisClient.SetNX(ctx, lockKey, []byte(uid), LockExpiry)
	if err != nil {
		return err
	}

	if res.Val() == false {
		ttl, _ := userRepo.redisClient.GetTTL(ctx, lockKey)
		logrus.Debug("Lock Exist with TTL ", ttl)
		return helpers.TryLater
	}

	pipe := userRepo.redisClient.Begin()

	_, err = userRepo.redisClient.PipelineSetTTL(pipe, ctx, key, value, expiration)
	if err != nil {
		// TODO: Add the retry if required here
		return err
	}

	_, err = userRepo.redisClient.PipelineInjectScripts(pipe, ctx, ReleaseLockScript, []string{lockKey}, []byte(uid))
	if err != nil {
		// TODO: Add the retry if required here
		return err
	}

	execResult, err := userRepo.redisClient.Commit(ctx, pipe)
	if err != nil {
		// TODO: Add the retry if required here
		return err
	}

	for _, result := range execResult {
		logrus.Debug(result.String())
		if result.Err() != nil {
			// Todo Add the errors in list and send it to client
			logrus.WithError(result.Err()).Error("Partially failed ")
		}
	}

	// TODO: Check the transaction completed with delete commands if delete failed retry.
	logrus.Debug("Inside the UpdateUserWithPessimisticLocking Repository func")
	return nil
}

func (userRepo UserRepo) UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	logrus.Debug("Inside the UpdateUserWithOptimisticLocking Repository func")

	err := userRepo.redisClient.Watch(ctx,
		userRepo.redisClient.RunTransactionFunc(ctx, key, value, expiration),
		key)
	if err != nil {
		if err.Error() == helpers.RedisNil {
			return helpers.UserMayDeleted
		}
		// TODO: Add retry mechanism here
		logrus.WithError(err).Debug("Transaction failed with error")
		return helpers.ConflictUpdate
	}

	logrus.Debug("Completed the UpdateUserWithOptimisticLocking Repository func")
	return nil
}

func (userRepo UserRepo) SubscribeForKeyExpireChannel(ctx context.Context, eventStream chan string, errChan chan error) {
	stream := userRepo.redisClient.Subscribe(ctx, "__keyevent@0__:expired")
	for {
		logrus.Info("Start Listening for message")
		msg, err := stream.ReceiveMessage(ctx)
		if err != nil {
			logrus.WithError(err).Error("Error while ReceiveMessage for event key expired")
			errChan <- err
		}
		logrus.Info("Event Received")
		eventStream <- msg.String()
	}
}

// This is just write integration test
func (userRepo UserRepo) GetUserLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf(LockPrefix + key)
	timeout := time.Minute * 5

	res, err := userRepo.redisClient.SetNX(ctx, lockKey, nil, timeout)
	if err != nil {
		return err
	}

	if res.Val() == false {
		ttl, _ := userRepo.redisClient.GetTTL(ctx, lockKey)
		logrus.Debug("Lock Exist with TTL ", ttl)
		return helpers.TryLater
	}
	return nil
}

// This is just write integration test
func (userRepo UserRepo) DeleteUserLock(ctx context.Context, key string) error {
	return userRepo.redisClient.Delete(ctx, key)
}
