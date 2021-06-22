package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/velann21/bloom-services/common-lib/databases"
	"github.com/velann21/bloom-services/common-lib/entities/requests"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"github.com/velann21/bloom-services/users-service/pkg/database"
	"github.com/velann21/bloom-services/users-service/pkg/repository"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		requests           *requests.User
		ctx                context.Context
		expectedResultType string
		err                error
	}{
		{
			requests:           &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}},
			ctx:                context.Background(),
			expectedResultType: CreateUserSuccess,
			err:                nil,
		},
		{
			requests:           &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}},
			ctx:                context.Background(),
			expectedResultType: CreateUserAlreadyExist,
			err:                helpers.UserAlreadyExist,
		},
	}

	for _, test := range tests {
		service := NewUserService(NewMockUserRepository(test.expectedResultType))
		err := service.CreateUser(test.ctx, test.requests)
		assert.IsType(t, test.err, err)
	}
}

func TestUserService_UpdateUserWithOptimisticLock(t *testing.T) {
	tests := []struct {
		requests           *requests.User
		ctx                context.Context
		expectedResultType string
		err                error
	}{
		{
			requests:           &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}},
			ctx:                context.Background(),
			expectedResultType: UpdateUserSuccess,
			err:                nil,
		},
	}
	for _, test := range tests {
		service := NewUserService(NewMockUserRepository(test.expectedResultType))
		err := service.UpdateUserWithOptimisticLock(test.ctx, test.requests)
		assert.IsType(t, test.err, err)
	}
}

func TestUserService_UpdateUserWithPessimisticLock(t *testing.T) {
	tests := []struct {
		requests           *requests.User
		ctx                context.Context
		expectedResultType string
		err                error
	}{
		{
			requests:           &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}},
			ctx:                context.Background(),
			expectedResultType: UpdateUserSuccess,
			err:                nil,
		},
	}
	for _, test := range tests {
		service := NewUserService(NewMockUserRepository(test.expectedResultType))
		err := service.UpdateUserWithPessimisticLock(test.ctx, test.requests)
		assert.IsType(t, test.err, err)
	}
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		requests           string
		ctx                context.Context
		expectedResultType string
		email              string
		err                error
	}{
		{
			requests:           "velann21@gmail.com",
			ctx:                context.Background(),
			expectedResultType: GetUserSuccess,
			email:              "velann21@gmail.com",
			err:                nil,
		},
	}

	for _, test := range tests {
		service := NewUserService(NewMockUserRepository(test.expectedResultType))
		resp, err := service.GetUser(test.ctx, test.requests)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.email, resp.Email)
	}
}

func Test_UserService_CreateUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []struct {
		ctx      context.Context
		requests *requests.User
		err      error
	}{
		{ctx: context.Background(), requests: &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}}, err: nil},
	}
	for _, test := range tests {
		userRepo := repository.NewUserRepo(RedisConnection())
		service := NewUserService(userRepo)
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		email := fmt.Sprintf("velan" + strconv.Itoa(r1.Intn(100)) + "@gmail.com")
		test.requests.Email = email
		err := service.CreateUser(test.ctx, test.requests)
		fmt.Println(err)
		fmt.Println(test.requests.Email)
		assert.IsType(t, test.err, err)
	}
}

func Test_UserService_UpdateUserWithPessimisticLockIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []struct {
		ctx      context.Context
		requests *requests.User
		err      error
		lock     bool
	}{
		{lock: true, ctx: context.Background(), requests: &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}}, err: errors.New("Update after ")},
		{lock: false, ctx: context.Background(), requests: &requests.User{Name: "velan", Email: "velann21@gmail.com", Address: requests.Address{ZipCode: "1096GJ", StreetName: "Piterberge", HouseNumber: 74}, DOB: requests.DOB{Day: 21, Year: 1992, Month: 10}}, err: nil},
	}
	for _, test := range tests {
		userRepo := repository.NewUserRepo(RedisConnection())
		_ = userRepo.DeleteUserLock(test.ctx, test.requests.Email)
		service := NewUserService(userRepo)
		_ = service.CreateUser(test.ctx, test.requests)
		if test.lock == true {
			_ = userRepo.GetUserLock(test.ctx, test.requests.Email)
		}
		err := service.UpdateUserWithPessimisticLock(test.ctx, test.requests)
		if err != nil {
			assert.Contains(t, err.Error(), test.err.Error())
		}
		if test.lock == true {
			_ = userRepo.DeleteUserLock(test.ctx, repository.LockPrefix+test.requests.Email)
		}
	}
}


func RedisConnection() *databases.Redis {
	helpers.SetEnvUsersDevelopmentMode()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rc := &database.RedisConnection{}
	err := rc.NewRedisConnection(ctx, helpers.GetEnv(helpers.REDIS), "")
	if err != nil {
		logrus.WithError(err).Error("Something went wrong in redis connection")
		os.Exit(1)
	}
	return rc.Client
}

type MockUserRepository struct {
	expectedResultType string
}

func (repo MockUserRepository) DeleteUserLock(ctx context.Context, key string) error {
	return nil
}

const (
	CreateUserSuccess      = "CreateSuccess"
	CreateUserAlreadyExist = "UserAlreadyExist"
	GetUserSuccess         = "GetSuccess"
	UpdateUserSuccess      = "UpdateSuccess"
)

func NewMockUserRepository(expectedResultType string) *MockUserRepository {
	return &MockUserRepository{expectedResultType: expectedResultType}
}

func (repo MockUserRepository) CreateUser(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if repo.expectedResultType == CreateUserSuccess {
		return nil
	}
	return nil
}
func (repo MockUserRepository) GetUser(ctx context.Context, key string) ([]byte, error) {
	if repo.expectedResultType == CreateUserSuccess {
		return nil, nil
	} else if repo.expectedResultType == CreateUserAlreadyExist || repo.expectedResultType == UpdateUserSuccess {
		return []byte("{\"name\":\"velan\",\"email\":\"velann21@gmail.com\",\"address\":{\"zip_code\":\"1096BM\",\"street_name\":\"pietersbergweg\",\"house_number\":74},\"dob\":{\"year\":1992,\"month\":10,\"day\":21},\"created_at\":\"2021-06-22T08:07:44.092025+02:00\",\"updated_at\":\"2021-06-22T08:08:10.617286+02:00\"}\n"), nil
	}
	return []byte("{\"name\":\"velan\",\"email\":\"velann21@gmail.com\",\"address\":{\"zip_code\":\"1096BM\",\"street_name\":\"pietersbergweg\",\"house_number\":74},\"dob\":{\"year\":1992,\"month\":10,\"day\":21},\"created_at\":\"2021-06-22T08:07:44.092025+02:00\",\"updated_at\":\"2021-06-22T08:08:10.617286+02:00\"}\n"), nil
}
func (repo MockUserRepository) UpdateUserWithOptimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return nil
}
func (repo MockUserRepository) UpdateUserWithPessimisticLocking(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return nil
}
func (repo MockUserRepository) GetUserLock(ctx context.Context, key string) error {
	return nil
}
