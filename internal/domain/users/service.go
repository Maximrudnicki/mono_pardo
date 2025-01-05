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

type serviceImpl struct {
	Config     config.Config
	Validate   *validator.Validate
	Repository Repository
}

func NewServiceImpl(
	config config.Config,
	validate *validator.Validate,
	repository Repository) Service {
	return &serviceImpl{
		Config:     config,
		Validate:   validate,
		Repository: repository,
	}
}

func (s *serviceImpl) Login(user request.LoginRequest) (string, error) {
	new_user, user_err := s.Repository.FindByEmail(user.Email)
	if user_err != nil {
		return "", user_err
	}

	verify_error := utils.VerifyPassword(new_user.Password, user.Password)
	if verify_error != nil {
		return "", verify_error
	}

	token, err_token := utils.GenerateToken(s.Config.TokenExpiresIn, new_user.Id, s.Config.TokenSecret)
	if err_token != nil {
		return "", err_token
	}
	return token, nil
}

func (s *serviceImpl) Register(user request.CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	newUser := User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	save_err := s.Repository.Save(newUser)
	if save_err != nil {
		return save_err
	}

	return nil
}

func (s *serviceImpl) GetUserId(token string) (int, error) {
	user, err := utils.ValidateToken(token, s.Config.TokenSecret)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(fmt.Sprint(user))

	if err != nil {
		fmt.Errorf("failed to get id: %w\n", err)
	}

	return userId, nil
}

func (s *serviceImpl) FindUser(userId int) (response.UserResponse, error) {
	user, err := s.Repository.FindById(userId)
	if err != nil {
		log.Printf("find user error: %v\n", err)
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
