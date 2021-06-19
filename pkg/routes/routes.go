package routes

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Routes(router *mux.Router){
	logrus.Info(router)
}
