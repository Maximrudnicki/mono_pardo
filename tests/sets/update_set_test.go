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
	"mono_pardo/pkg/data/request"
	"mono_pardo/tests"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSet(t *testing.T) {
	env, _ := tests.NewTestMongoEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	mockAuthService := &tests.MockAuthService{}
	mockAuthService.On("GetUserId", "test-token").Return(1, nil)
	mockAuthService.On("GetUserId", "").Return(0, fmt.Errorf("empty token"))

	createdAt := time.Now().UTC().Truncate(time.Millisecond)

	firstSet, _ := setsDomain.NewWordSet("first Set", 1)
	secondSet, _ := setsDomain.NewWordSet("second Set", 2)
	thirdSet, _ := setsDomain.NewWordSet("third Set", 1) // don't save this set to the database

	firstSet.Words = []int{1, 3}
	secondSet.Words = []int{2, 4}
	firstSet.CreatedAt, secondSet.CreatedAt = createdAt, createdAt

	fixture := &tests.WordSetFixture{
		Sets: []setsDomain.WordSet{*firstSet, *secondSet},
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
	setsbGroup.GET("/:setId", setsController.GetSet)
	setsbGroup.PATCH("/:setId", setsController.UpdateSet)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", firstSet.Id.Hex()), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Input", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", firstSet.Id.Hex()), nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Values", func(t *testing.T) {
		payload := request.UpdateSetRequest{
			Updates: []request.FieldUpdate{
				{
					Field: "name",
					Value: "",
				},
			},
		}

		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", firstSet.Id.Hex()), bytes.NewBuffer(jsonData))

		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Set Not Found", func(t *testing.T) {
		payload := request.UpdateSetRequest{
			Updates: []request.FieldUpdate{
				{
					Field: "name",
					Value: "new name",
				},
			},
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", thirdSet.Id.Hex()), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Other User's Set", func(t *testing.T) {
		payload := request.UpdateSetRequest{
			Updates: []request.FieldUpdate{
				{
					Field: "name",
					Value: "new name",
				},
			},
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", secondSet.Id.Hex()), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), `"access forbidden"`)
	})

	t.Run("Invalid Value Type", func(t *testing.T) {
		payload := request.UpdateSetRequest{
			Updates: []request.FieldUpdate{
				{
					Field: "name",
					Value: true,
				},
			},
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", firstSet.Id.Hex()), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success Update Set", func(t *testing.T) {
		payload := request.UpdateSetRequest{
			Updates: []request.FieldUpdate{
				{
					Field: "name",
					Value: "new name",
				},
			},
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH", fmt.Sprintf("/api/v1/sets/%v", firstSet.Id.Hex()), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
