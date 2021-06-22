package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/entities/requests"
	"github.com/velann21/bloom-services/common-lib/entities/response"
	"github.com/velann21/bloom-services/common-lib/helpers"
	userResponse "github.com/velann21/bloom-services/users-service/pkg/entities/response"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	"net/http"
)

type User struct {
	service service.UserInterface
}

func NewUserController(srv service.UserInterface) *User {
	return &User{service: srv}
}

func (user User) CreateUser(resp http.ResponseWriter, req *http.Request) {
	logrus.Debug("Inside the CreateUser Controller")
	commonSuccessResponse := response.NewSuccessResponse()
	successResponse := userResponse.Response{Success: commonSuccessResponse}

	errorResponse := response.NewErrorResponse()
	userEntity := requests.NewUserEntity()
	err := userEntity.PopulateUser(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Error in PopulateUser() CreateUser")
		errorResponse.HandleError(err, resp)
		return
	}
	err = userEntity.ValidateUser()
	if err != nil {
		logrus.WithError(err).Error("Error while validate the requests CreateUser")
		errorResponse.HandleError(err, resp)
		return
	}
	err = user.service.CreateUser(req.Context(), userEntity)
	if err != nil {
		logrus.WithError(err).Error("Error while CreateUser service")
		errorResponse.HandleError(err, resp)
		return
	}
	successResponse.CreateUserResponse(userEntity.Email)
	successResponse.Success.SuccessResponse(resp, http.StatusCreated)
	logrus.Debug("Completed the CreateUser Controller")
	return
}

func (user User) UpdateUserWithOptimisticLock(resp http.ResponseWriter, req *http.Request) {
	logrus.Debug("Inside the UpdateUserWithOptimisticLock Controller")
	commonSuccessResponse := response.NewSuccessResponse()
	successResponse := userResponse.Response{Success: commonSuccessResponse}
	errorResponse := response.NewErrorResponse()
	userEntity := requests.NewUserEntity()
	err := userEntity.PopulateUser(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Error in PopulateUser() UpdateUserWithOptimisticLock")
		errorResponse.HandleError(err, resp)
		return
	}
	err = userEntity.ValidateUser()
	if err != nil {
		logrus.WithError(err).Error("Error while validate the requests UpdateUserWithOptimisticLock")
		errorResponse.HandleError(err, resp)
		return
	}
	err = user.service.UpdateUserWithOptimisticLock(req.Context(), userEntity)
	if err != nil {
		logrus.WithError(err).Error("Error while UpdateUserWithOptimisticLock service")
		errorResponse.HandleError(err, resp)
		return
	}
	successResponse.UpdateUserResponse(userEntity.Email)
	successResponse.Success.SuccessResponse(resp, http.StatusOK)
	logrus.Debug("Completed the UpdateUserWithOptimisticLock Controller")
}

func (user User) UpdateUserWithPessimisticLock(resp http.ResponseWriter, req *http.Request) {
	logrus.Debug("Inside the UpdateUserWithPessimisticLock Controller")
	commonSuccessResponse := response.NewSuccessResponse()
	successResponse := userResponse.Response{Success: commonSuccessResponse}
	errorResponse := response.NewErrorResponse()
	userEntity := requests.NewUserEntity()
	err := userEntity.PopulateUser(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Error in PopulateUser() UpdateUserWithPessimisticLock")
		errorResponse.HandleError(err, resp)
		return
	}
	err = userEntity.ValidateUser()
	if err != nil {
		logrus.WithError(err).Error("Error while validate the requests UpdateUserWithPessimisticLock")
		errorResponse.HandleError(err, resp)
		return
	}
	err = user.service.UpdateUserWithPessimisticLock(req.Context(), userEntity)
	if err != nil {
		logrus.WithError(err).Error("Error while UpdateUserWithOptimisticLock service")
		errorResponse.HandleError(err, resp)
		return
	}
	successResponse.UpdateUserResponse(userEntity.Email)
	successResponse.Success.SuccessResponse(resp, http.StatusOK)
	logrus.Debug("Completed the UpdateUserWithPessimisticLock Controller")
}

func (user User) GetUser(resp http.ResponseWriter, req *http.Request) {
	errorResponse := response.NewErrorResponse()
	emailID := req.URL.Query().Get("email")
	if emailID == "" {
		logrus.Error("No Email provided GetUser()")
		errorResponse.HandleError(helpers.InvalidRequest, resp)
		return
	}
	commonSuccessResponse := response.NewSuccessResponse()
	successResponse := userResponse.Response{Success: commonSuccessResponse}

	userData, err := user.service.GetUser(req.Context(), emailID)
	if err != nil {
		logrus.WithError(err).Error("Error while GetUser service")
		errorResponse.HandleError(err, resp)
		return
	}

	successResponse.GetUserResponse(userData)
	successResponse.Success.SuccessResponse(resp, http.StatusOK)
	logrus.Debug("Completed the GetUser Controller")
}
