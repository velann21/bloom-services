package response

import (
	"github.com/velann21/bloom-services/common-lib/entities/response"
	"github.com/velann21/bloom-services/users-service/pkg/entities/models"
)

type Response struct {
	Success *response.SuccessResponse
}

func (resp *Response) CreateUserResponse(email string){
	responseData := make([]map[string]interface{}, 0)
	metaData := make(map[string]interface{})
	metaData["message"] = "Success"
	metaData["email"] = email
	responseData = append(responseData, metaData)
	resp.Success.Data = responseData
	resp.Success.Success = true

}

func (resp *Response) UpdateUserResponse(email string){
	responseData := make([]map[string]interface{}, 0)
	metaData := make(map[string]interface{})
	metaData["email"] = email
	responseData = append(responseData, metaData)
	resp.Success.Data = responseData
	resp.Success.Success = true
}

func (resp *Response) GetUserResponse(user *models.User){
	responseData := make([]map[string]interface{}, 0)
	metaData := make(map[string]interface{})
	metaData["user"] = user
	responseData = append(responseData, metaData)
	resp.Success.Data = responseData
	resp.Success.Success = true
}

