package routes

import (
	"github.com/velann21/bloom-services/common-lib/databases"
	server "github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/controller"
	"github.com/velann21/bloom-services/users-service/pkg/repository"
	"github.com/velann21/bloom-services/users-service/pkg/service"
)

func Routes(router *server.Router, redisConnection *databases.Redis){
	userRepo := repository.NewUserRepo(redisConnection)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	router.Router.Path("/user").HandlerFunc(userController.GetUser).Methods("GET")
	router.Router.Path("/user").HandlerFunc(userController.CreateUser).Methods("POST")
	router.Router.Path("/user/optimistic").HandlerFunc(userController.UpdateUserWithOptimisticLock).Methods("PUT")
	router.Router.Path("/user/pessimistic").HandlerFunc(userController.UpdateUserWithPessimisticLock).Methods("PUT")
}
