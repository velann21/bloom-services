package service

type UserInterface interface {
	CreateUser()
}

type UserService struct {

}

func (users UserService) CreateUser() {}


