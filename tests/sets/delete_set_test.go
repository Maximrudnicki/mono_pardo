package sets

import (
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

func TestDeleteSet(t *testing.T) {
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

	firstSetURL := fmt.Sprintf("/api/v1/sets/%s", firstSet.Id.Hex())
	secondSetURL := fmt.Sprintf("/api/v1/sets/%s", secondSet.Id.Hex())
	thirdSetURL := fmt.Sprintf("/api/v1/sets/%s", thirdSet.Id.Hex())

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
	setsbGroup.DELETE("/:setId", setsController.DeleteSet)
	setsbGroup.GET("/:setId", setsController.GetSet)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", secondSetURL, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Set ID Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/sets/123", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Set Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", thirdSetURL, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Attempt to Delete Other User's WordSet", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", secondSetURL, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Success Delete WordSet", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", firstSetURL, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		checkW := httptest.NewRecorder()
		checkReq, _ := http.NewRequest("GET", firstSetURL, nil)
		checkReq.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(checkW, checkReq)

		assert.Equal(t, http.StatusNotFound, checkW.Code)
	})
}
