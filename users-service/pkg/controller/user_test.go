package controller

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/velann21/bloom-services/common-lib/entities/requests"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"github.com/velann21/bloom-services/common-lib/server"
	"github.com/velann21/bloom-services/users-service/pkg/entities/models"
	"github.com/velann21/bloom-services/users-service/pkg/service"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"
)

func TestUser_CreateUserController(t *testing.T) {
	tests := []struct {
		body               string
		serviceTYpe        string
		expectedStatusCode string
	}{
		{body: ProcessUserRequestBody(Rand{rand.Intn(10000)}), serviceTYpe: "200", expectedStatusCode: "201 Created"},
		{body: invalidRequests, serviceTYpe: "400", expectedStatusCode: "400 Bad Request"},
		{body: "", serviceTYpe: "400", expectedStatusCode: "400 Bad Request"},
		{body: ProcessUserRequestBody(Rand{rand.Intn(10000)}), serviceTYpe: "400", expectedStatusCode: "500 Internal Server Error"},
		{body: ProcessUserRequestBody(Rand{rand.Intn(10000)}), serviceTYpe: "500", expectedStatusCode: "500 Internal Server Error"},
	}
	for _, test := range tests {
		request, _ := http.NewRequest(http.MethodPost, "/users/api/v1/user", strings.NewReader(test.body))
		response := httptest.NewRecorder()
		MockRouter(NewMockUserService(test.serviceTYpe)).ServeHTTP(response, request)
		assert.Equal(t, response.Result().Status, test.expectedStatusCode)
	}
}

func TestUser_GetUserController(t *testing.T) {
	tests := []struct {
		email              string
		expectedStatusCode string
		serviceTYpe        string
	}{
		{email: "jack@gmail.com", expectedStatusCode: "200 OK", serviceTYpe: "200"},
		{email: "", expectedStatusCode: "400 Bad Request", serviceTYpe: "400"},
		{email: "", expectedStatusCode: "400 Bad Request", serviceTYpe: "500"},
	}
	for _, test := range tests {
		request, _ := http.NewRequest(http.MethodGet, "/users/api/v1/user?email="+test.email, nil)
		response := httptest.NewRecorder()
		MockRouter(NewMockUserService(test.serviceTYpe)).ServeHTTP(response, request)
		assert.Equal(t, test.expectedStatusCode, response.Result().Status)
	}
}

type MockUserService struct {
	expectedOutput string
}

func (service MockUserService) UserExpiredEvent(ctx context.Context, eventStream chan string, errChan chan error) {

}

func NewMockUserService(expectedOutput string) *MockUserService {
	return &MockUserService{expectedOutput: expectedOutput}
}
func (service MockUserService) CreateUser(ctx context.Context, data *requests.User) error {
	if service.expectedOutput == "400" {
		return helpers.UserAlreadyExist
	}
	if service.expectedOutput == "500" {
		return helpers.SomethingWrong
	}
	return nil
}
func (service MockUserService) GetUser(ctx context.Context, email string) (*models.User, error) {
	if service.expectedOutput == "400" {
		return nil, helpers.NoResultFound
	}
	if service.expectedOutput == "500" {
		return nil, helpers.SomethingWrong
	}
	return nil, nil
}
func (service MockUserService) UpdateUserWithOptimisticLock(ctx context.Context, data *requests.User) error {
	return nil
}
func (service MockUserService) UpdateUserWithPessimisticLock(ctx context.Context, data *requests.User) error {
	return nil
}

func MockRouter(userService service.UserInterface) *mux.Router {
	router := server.NewMux()
	userRoutes := router.PathPrefix("/users/api/v1").Subrouter()
	controller := NewUserController(userService)
	userRoutes.Path("/user").HandlerFunc(controller.GetUser).Methods("GET")
	userRoutes.Path("/user").HandlerFunc(controller.CreateUser).Methods("POST")
	userRoutes.Path("/user/optimistic").HandlerFunc(controller.UpdateUserWithOptimisticLock).Methods("PUT")
	userRoutes.Path("/user/pessimistic").HandlerFunc(controller.UpdateUserWithPessimisticLock).Methods("PUT")
	return router
}

func ProcessUserRequestBody(Random Rand) string {
	var tmplBytes bytes.Buffer
	t := template.New("User template")
	t, err := t.Parse(bodySuccess)
	if err != nil {
		return ""
	}
	err = t.Execute(&tmplBytes, Random)
	if err != nil {
		return ""
	}
	return tmplBytes.String()
}

type Rand struct {
	Random int
}

const bodySuccess = `{
    "name": "velan",
    "email":"jack{{.Random}}@gmail.com",
    "address":{
        "zip_code":"1096BM",
        "street_name":"pietersbergweg",
        "house_number": 102
    },
    "dob": {
        "year": 1992,
        "month": 10,
        "day":21
    }
}`

const invalidRequests = `{
    "name": "",
    "email":"jack@gmail.com",
    "address":{
        "zip_code":"1096BM",
        "street_name":"pietersbergweg",
        "house_number": 102
    },
    "dob": {
        "year": 1992,
        "month": 10,
        "day":21
    }
}`
