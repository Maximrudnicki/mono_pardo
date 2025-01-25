package sets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mono_pardo/internal/api/controller"
	"mono_pardo/internal/api/middleware"
	setsDomain "mono_pardo/internal/domain/sets"
	setsInfra "mono_pardo/internal/infrastructure/sets"
	"mono_pardo/tests"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestCreateSet(t *testing.T) {
	env, _ := tests.NewTestMongoEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	mockAuthService := &tests.MockAuthService{}
	mockAuthService.On("GetUserId", "test-token").Return(1, nil)
	mockAuthService.On("GetUserId", "").Return(0, fmt.Errorf("empty token"))

	createdAt := time.Now().UTC().Truncate(time.Millisecond)

	firstSet, _ := setsDomain.NewWordSet("first Set", 1)

	firstSet.Words = []int{1, 3}
	firstSet.CreatedAt = createdAt

	fixture := &tests.WordSetFixture{
		Sets: []setsDomain.WordSet{*firstSet},
	}
	cleanup := env.WithMongoFixture(t, fixture)
	defer cleanup()

	setsRepository := setsInfra.NewMongoRepositoryImpl(env.MongoDB.Collection)
	validate := validator.New()
	setsService := setsDomain.NewServiceImpl(validate, setsRepository)
	setsController := controller.NewSetsController(setsService)

	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)

	router := env.Router
	setsbGroup := router.Group("/api/v1/sets")
	setsbGroup.Use(authMiddleware.Handle())
	setsbGroup.GET("", setsController.GetSets)
	setsbGroup.POST("", setsController.CreateSet)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST", "/api/v1/sets", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Input", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/sets", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Name", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST", "/api/v1/sets", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Duplicate WordSet", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "first Set",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/sets", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Success Create WordSet", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "second Set",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/sets", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		checkW := httptest.NewRecorder()
		checkReq, _ := http.NewRequest("GET", "/api/v1/sets", nil)
		checkReq.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(checkW, checkReq)

		assert.Equal(t, http.StatusOK, checkW.Code)
		assert.Contains(t, checkW.Body.String(), `"name":"second Set"`)
	})
}
