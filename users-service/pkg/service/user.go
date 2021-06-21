package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/velann21/bloom-services/common-lib/entities/requests"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"github.com/velann21/bloom-services/users-service/pkg/entities/models"
	"github.com/velann21/bloom-services/users-service/pkg/repository"
	"time"
)

type UserInterface interface {
	CreateUser(ctx context.Context, data *requests.User)error
	GetUser(ctx context.Context, email string)(*models.User, error)
	UpdateUserWithOptimisticLock(ctx context.Context, data *requests.User)error
	UpdateUserPessimisticLock(ctx context.Context, email string, data *requests.User)
}

type UserService struct {
	userRepo repository.UserRepoInterface
}

func NewUserService(userRepo repository.UserRepoInterface)UserInterface{
	return &UserService{userRepo: userRepo}
}

func (users UserService) CreateUser(ctx context.Context, data *requests.User)error{
	result, err := users.userRepo.GetUser(ctx, data.Email)
	if err != nil{
		fmt.Println(err.Error())
		if err.Error() == helpers.RedisNil{
			goto createUser
		}
		return err
	}
	if result != nil{
		return helpers.UserAlreadyExist
	}

	createUser:
	userModel := &models.User{}
	userModel.Email = data.Email
	userModel.Name = data.Name
	userModel.Address.ZipCode = data.Address.ZipCode
	userModel.Address.StreetName = data.Address.StreetName
	userModel.Address.HouseNumber = data.Address.HouseNumber
	userModel.DOB.Month = data.DOB.Month
	userModel.DOB.Year = data.DOB.Year
	userModel.DOB.Day = data.DOB.Day
	userModel.CreatedAt = time.Now()
	userModel.UpdatedAt = time.Now()

	reqBodyBytes := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBytes).Encode(userModel)
	if err != nil{
		return err
	}

	err = users.userRepo.CreateUser(ctx, userModel.Email, reqBodyBytes.Bytes())
	if err != nil{
		return err
	}
	return nil
}

func (users UserService) GetUser(ctx context.Context, email string)(*models.User, error){
	result, err := users.userRepo.GetUser(ctx, email)
	if err != nil{
		return nil, err
	}
	if result == nil{
		return nil, helpers.NoResultFound
	}
	userModel := &models.User{}
	err = userModel.PopulateUser(result)
	if err != nil{
		return nil, err
	}
	fmt.Println(userModel)
	return userModel, nil
}

func (users UserService) UpdateUserPessimisticLock(ctx context.Context, email string, data *requests.User){

}

func (users UserService) UpdateUserWithOptimisticLock(ctx context.Context, data *requests.User)error{
	user, err := users.GetUser(ctx, data.Email)
	if err != nil{
		return err
	}
	userModel := &models.User{}
	userModel.Email = data.Email
	userModel.Name = data.Name
	userModel.Address.ZipCode = data.Address.ZipCode
	userModel.Address.StreetName = data.Address.StreetName
	userModel.Address.HouseNumber = data.Address.HouseNumber
	userModel.DOB.Month = data.DOB.Month
	userModel.DOB.Year = data.DOB.Year
	userModel.DOB.Day = data.DOB.Day
	userModel.CreatedAt = user.CreatedAt
	userModel.UpdatedAt = time.Now()
	reqBodyBytes := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBytes).Encode(userModel)
	if err != nil{
		return err
	}
	err = users.userRepo.UpdateUserWithOptimisticLocking(ctx, data.Email, reqBodyBytes.Bytes())
	if err != nil{
		return err
	}
	return nil
}
