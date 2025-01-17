package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"

	"mono_pardo/internal/api/controller"
	usersDomain "mono_pardo/internal/domain/users"
	usersInfra "mono_pardo/internal/infrastructure/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/tests"
)

func TestRegister(t *testing.T) {
	env, testConfig := tests.NewTestEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	userRepository := usersInfra.NewPostgresRepositoryImpl(env.DB.DB)
	validate := validator.New()
	authenticationService := usersDomain.NewServiceImpl(testConfig, validate, userRepository)
	authenticationController := controller.NewAuthenticationController(authenticationService)

	router := env.Router
	authenticationGroup := router.Group("/api/v1/authentication")
	authenticationGroup.POST("/register", authenticationController.Register)

	t.Run("Successful Registration", func(t *testing.T) {
		payload := request.CreateUserRequest{
			Username: "new user",
			Email:    "new@email.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdUser usersDomain.User
		err := env.DB.DB.Where("email = ?", payload.Email).First(&createdUser).Error
		assert.NoError(t, err)
		assert.Equal(t, payload.Username, createdUser.Username)
		assert.Equal(t, payload.Email, createdUser.Email)
		assert.NotEqual(t, payload.Password, createdUser.Password)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		payload := request.CreateUserRequest{
			Username: "another user",
			Email:    "new@email.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := []byte(`{"username": "valid name", "email": "valid@email.com", "password": }`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/register", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		payload := request.CreateUserRequest{
			Username: "valid name",
			Email:    "invalid-email-address",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
