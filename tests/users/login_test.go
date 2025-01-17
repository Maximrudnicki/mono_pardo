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
	"mono_pardo/internal/utils"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
	"mono_pardo/tests"
)

func TestLogin(t *testing.T) {
	env, testConfig := tests.NewTestEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	hashedPassword, _ := utils.HashPassword("test_password")
	fixture := &tests.UserFixture{
		Users: []usersDomain.User{
			{
				Username: "test username",
				Email:    "test@email.com",
				Password: hashedPassword,
			},
		},
	}

	cleanup := env.WithFixture(t, fixture)
	defer cleanup()

	userRepository := usersInfra.NewPostgresRepositoryImpl(env.DB.DB)
	validate := validator.New()
	authenticationService := usersDomain.NewServiceImpl(testConfig, validate, userRepository)
	authenticationController := controller.NewAuthenticationController(authenticationService)

	router := env.Router
	authenticationGroup := router.Group("/api/v1/authentication")
	authenticationGroup.POST("/login", authenticationController.Login)

	t.Run("Invalid Password", func(t *testing.T) {
		payload := request.LoginRequest{
			Email:    "test@email.com",
			Password: "wrong_password",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Non-existent Email", func(t *testing.T) {
		payload := request.LoginRequest{
			Email:    "nonexistent@email.com",
			Password: "test_password",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		payload := request.LoginRequest{
			Email:    "invalid-email",
			Password: "test_password",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Password", func(t *testing.T) {
		payload := request.LoginRequest{
			Email:    "test@email.com",
			Password: "",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := []byte(`{"email": "test@email.com", "password": }`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Successful Login", func(t *testing.T) {
		payload := request.LoginRequest{
			Email:    "test@email.com",
			Password: "test_password",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authentication/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response response.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.NotEmpty(t, response.Token)
	})
}
