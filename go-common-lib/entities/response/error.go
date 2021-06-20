package response

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"net/http"
)

type ErrorResponse struct {
	Success bool `json:"success"`
	Errors  []Error `json:"errors"`
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

type Error struct {
	Message   string `json:"message"`
	ErrorCode int `json:"error_code"`
}

func (err *ErrorResponse) HandleError(er error, w http.ResponseWriter) {
	if er == nil {
		logrus.Error("invalid error")
		return
	}
	errList := make([]Error, 0)
	w.Header().Set("Content-Type", "application/json")
	switch er {
	case helpers.InvalidRequest:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 1,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(400)

		_ = json.NewEncoder(w).Encode(resp)

	case helpers.SomethingWrong:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 0,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(resp)
	case helpers.NoresultFound:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 2,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(resp)
	default:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 0,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(resp)
	}

}
