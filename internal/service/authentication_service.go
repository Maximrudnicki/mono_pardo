package service

import (
	"fmt"
	"log"

	usersInfra "mono_pardo/internal/infrastructure/users"
	"mono_pardo/internal/model"
	"mono_pardo/internal/utils"
	"mono_pardo/pkg/config"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
	"strconv"

	"github.com/go-playground/validator"
)

type AuthenticationService interface {
	Login(user request.LoginRequest) (string, error)
	Register(user request.CreateUserRequest) error
	GetUserId(token string) (int, error)
	FindUser(userId int) (response.UserResponse, error)
}

type AuthenticationServiceImpl struct {
	Config         config.Config
	Validate       *validator.Validate
	UserRepository usersInfra.UsersRepository
}

func NewAuthenticationServiceImpl(
	config config.Config,
	validate *validator.Validate,
	userRepository usersInfra.UsersRepository) AuthenticationService {
	return &AuthenticationServiceImpl{
		Config:         config,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

// Login implements AuthenticationService.
func (a *AuthenticationServiceImpl) Login(user request.LoginRequest) (string, error) {
	// Find username in database
	new_user, user_err := a.UserRepository.FindByEmail(user.Email)
	if user_err != nil {
		return "", user_err
	}

	verify_error := utils.VerifyPassword(new_user.Password, user.Password)
	if verify_error != nil {
		return "", verify_error
	}

	// Generate Token
	token, err_token := utils.GenerateToken(a.Config.TokenExpiresIn, new_user.Id, a.Config.TokenSecret)
	if err_token != nil {
		return "", err_token
	}
	return token, nil
}

// Register implements AuthenticationService.
func (a *AuthenticationServiceImpl) Register(user request.CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	newUser := model.User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	save_err := a.UserRepository.Save(newUser)
	if save_err != nil {
		return save_err
	}

	return nil
}

func (a *AuthenticationServiceImpl) GetUserId(token string) (int, error) {
	user, err := utils.ValidateToken(token, a.Config.TokenSecret)
	if err != nil {
		return 0, err
	}

	userId, err_id := strconv.Atoi(fmt.Sprint(user))

	if err_id != nil {
		log.Printf("Failed to get id: %v\n", err_id)
	}

	return userId, nil
}

func (a *AuthenticationServiceImpl) FindUser(userId int) (response.UserResponse, error) {
	user, err := a.UserRepository.FindById(userId)
	if err != nil {
		log.Printf("find user error: %v\n", err)
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
