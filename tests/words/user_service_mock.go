package words

import (
	"github.com/stretchr/testify/mock"

	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(user request.LoginRequest) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(user request.CreateUserRequest) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthService) GetUserId(token string) (int, error) {
	args := m.Called(token)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthService) FindUser(userId int) (response.UserResponse, error) {
	args := m.Called(userId)
	return args.Get(0).(response.UserResponse), args.Error(1)
}
