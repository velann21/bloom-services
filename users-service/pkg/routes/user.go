package routes

import (
	"github.com/gorilla/mux"
	"github.com/velann21/bloom-users-service/pkg/controller"
	"github.com/velann21/bloom-users-service/pkg/service"
)

func Routes(router *mux.Router){
	userService := service.UserService{}
	userController := controller.NewUserController(userService)
	router.Path("/user").HandlerFunc(userController.CreateUser).Methods("POST")
}
