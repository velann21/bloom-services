package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool                     `json:"success"`
	Status  string                   `json:"status"`
	Data    []map[string]interface{} `json:"data"`
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{}
}

func (resp *SuccessResponse) SuccessResponse(rw http.ResponseWriter, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")

	switch statusCode {
	case http.StatusOK:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
		resp.Success = true
	case http.StatusCreated:
		rw.WriteHeader(http.StatusCreated)
		resp.Status = http.StatusText(http.StatusCreated)
		resp.Success = true
	case http.StatusAccepted:
		rw.WriteHeader(http.StatusAccepted)
		resp.Status = http.StatusText(http.StatusAccepted)
		resp.Success = true
	default:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
		resp.Success = true
	}
	// send response
	_ = json.NewEncoder(rw).Encode(resp)
	return
}
