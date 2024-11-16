package users

import (
	"fmt"
	"log"
	"strconv"

	"mono_pardo/internal/utils"
	"mono_pardo/pkg/config"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/go-playground/validator"
)

type ServiceImpl struct {
	Config     config.Config
	Validate   *validator.Validate
	Repository Repository
}

func NewServiceImpl(
	config config.Config,
	validate *validator.Validate,
	repository Repository) Service {
	return &ServiceImpl{
		Config:     config,
		Validate:   validate,
		Repository: repository,
	}
}

// Login implements Service.
func (a *ServiceImpl) Login(user request.LoginRequest) (string, error) {
	// Find username in database
	new_user, user_err := a.Repository.FindByEmail(user.Email)
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
func (a *ServiceImpl) Register(user request.CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	newUser := User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	save_err := a.Repository.Save(newUser)
	if save_err != nil {
		return save_err
	}

	return nil
}

func (a *ServiceImpl) GetUserId(token string) (int, error) {
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

func (a *ServiceImpl) FindUser(userId int) (response.UserResponse, error) {
	user, err := a.Repository.FindById(userId)
	if err != nil {
		log.Printf("find user error: %v\n", err)
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
