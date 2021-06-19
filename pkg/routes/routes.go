package routes

import (
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router){
	router.Path("/user").HandlerFunc(controller.CreateTodo).Methods("POST")
}
