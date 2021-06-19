package controller

import (
	"github.com/velann21/bloom-users-service/pkg/service"
	"net/http"
)

type User struct {
	service service.UserInterface
}

func NewUserController(srv service.UserInterface)*User{
	return &User{service: srv}
}

func (user User) CreateUser(resp http.ResponseWriter, req *http.Request){

}
