package service

import (
	"bytes"
	"context"
	"encoding/gob"
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
		if err.Error() == "redis: nil"{
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
	err = gob.NewEncoder(reqBodyBytes).Encode(userModel)
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
		return nil, helpers.NoresultFound
	}
	userModel := &models.User{}
	err = userModel.PopulateUser(result)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

