package response

import (
	"github.com/velann21/bloom-services/common-lib/entities/response"
)

type Response struct {
	Success *response.SuccessResponse
}

func (resp *Response) CreateUserResponse(email string){
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	responseData = append(responseData, data)
	resp.Success.Data = responseData
	resp.Success.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Success"
	metaData["email"] = email
}

