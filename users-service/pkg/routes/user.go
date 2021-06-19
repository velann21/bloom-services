package routes

import (
	server "github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/controller"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	)

func Routes(router *server.Router){
	userService := service.UserService{}
	userController := controller.NewUserController(userService)
	router.Router.Path("/user").HandlerFunc(userController.CreateUser).Methods("POST")
}
