package users

import (
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

type Service interface {
	Login(user request.LoginRequest) (string, error)
	Register(user request.CreateUserRequest) error
	GetUserId(token string) (int, error)
	FindUser(userId int) (response.UserResponse, error)
}

type Repository interface {
	Save(user User) error
	Delete(usersId int)
	FindById(usersId int) (User, error)
	FindAll() []User
	FindByEmail(email string) (User, error)
}
