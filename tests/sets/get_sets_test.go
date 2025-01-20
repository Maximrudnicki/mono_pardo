package sets

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"mono_pardo/internal/api/controller"
	"mono_pardo/internal/api/middleware"
	setsDomain "mono_pardo/internal/domain/sets"
	setsInfra "mono_pardo/internal/infrastructure/sets"
	"mono_pardo/tests"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestGetSets(t *testing.T) {
	env, _ := tests.NewTestMongoEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	mockAuthService := &tests.MockAuthService{}
	mockAuthService.On("GetUserId", "test-token").Return(1, nil)
	mockAuthService.On("GetUserId", "").Return(0, fmt.Errorf("empty token"))

	firstSet, _ := setsDomain.NewWordSet("first Set", 1)
	secondSet, _ := setsDomain.NewWordSet("second Set", 2)

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
	setsbGroup.GET("", setsController.GetSets)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/sets", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
