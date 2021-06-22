package routes

import (
	"encoding/json"
	"github.com/velann21/bloom-services/common-lib/server"
	"net/http"
)

func IndexRoutes(router *server.Router) {
	router.Router.Path("/health").HandlerFunc(HealthStatus).Methods("GET")
}

func HealthStatus(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	_ = json.NewEncoder(resp).Encode(resp)
	return
}
