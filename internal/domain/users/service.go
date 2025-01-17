package users

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	foundUser, err := s.Repository.FindByEmail(strings.TrimSpace(user.Email))
	if err != nil {
		return "", err
	}

	if err = utils.VerifyPassword(foundUser.Password, strings.TrimSpace(user.Password)); err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(s.Config.TokenExpiresIn, foundUser.Id, s.Config.TokenSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *serviceImpl) Register(user request.CreateUserRequest) error {
	newUser, err := NewUser(user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	if err = s.Repository.Save(*newUser); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) GetUserId(token string) (int, error) {
	user, err := utils.ValidateToken(token, s.Config.TokenSecret)
	if err != nil {
		return 0, errors.New("cannot validate token")
	}

	userId, err := strconv.Atoi(fmt.Sprint(user))
	if err != nil {
		return 0, fmt.Errorf("failed to get id: %w\n", err)
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
