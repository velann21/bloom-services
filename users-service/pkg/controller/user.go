package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/users-service/pkg/entities/requests"
	"github.com/velann21/bloom-services/users-service/pkg/entities/response"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	"net/http"
)

type User struct {
	service service.UserInterface
}

func NewUserController(srv service.UserInterface)*User {
	return &User{service: srv}
}

func (user User) CreateUser(resp http.ResponseWriter, req *http.Request){
	logrus.Debug("Inside the CreateUser Controller")
	errorResponse := response.NewErrorResponse()
	userEntity := requests.NewUserEntity()
	err := userEntity.PopulateUser(req.Body)
	if err != nil{
		logrus.WithError(err).Error("Error in PopulateUser()")
		errorResponse.HandleError(err, resp)
		return
	}
	err = userEntity.ValidateUser()
	if err != nil{
		logrus.WithError(err).Error("Error while validate the requests")
		errorResponse.HandleError(err, resp)
		return
	}
	user.service.CreateUser()
	return
}
