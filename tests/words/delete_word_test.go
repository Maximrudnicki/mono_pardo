package words

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"

	"mono_pardo/internal/api/controller"
	"mono_pardo/internal/api/middleware"
	wordsDomain "mono_pardo/internal/domain/words"
	wordsInfra "mono_pardo/internal/infrastructure/words"
	resp "mono_pardo/pkg/data/response"
	"mono_pardo/tests"
)

func TestDeleteWord(t *testing.T) {
	env := tests.NewTestEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	mockAuthService := &MockAuthService{}
	mockAuthService.On("GetUserId", "test-token").Return(1, nil)
	mockAuthService.On("GetUserId", "").Return(0, fmt.Errorf("empty token"))

	createdAt := time.Now()

	fixture := &tests.WordFixture{
		Words: []wordsDomain.Word{
			{
				UserId:          1,
				Word:            "hello",
				Definition:      "greeting",
				IsLearned:       false,
				WordTranslation: false,
				Cards:           true,
				Constructor:     false,
				WordAudio:       false,
				CreatedAt:       createdAt,
			},
			{
				UserId:          2,
				Word:            "world",
				Definition:      "planet earth",
				IsLearned:       true,
				WordTranslation: true,
				Cards:           true,
				Constructor:     true,
				WordAudio:       false,
				CreatedAt:       createdAt,
			},
		},
	}
	cleanup := env.WithFixture(t, fixture)
	defer cleanup()

	wordRepository := wordsInfra.NewPostgresRepositoryImpl(env.DB.DB)
	validate := validator.New()
	vocabService := wordsDomain.NewServiceImpl(validate, wordRepository)
	vocabController := controller.NewVocabController(vocabService)

	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)

	router := env.Router
	vocabGroup := router.Group("/api/v1/vocab")
	vocabGroup.Use(authMiddleware.Handle())
	vocabGroup.GET("", vocabController.GetWords)
	vocabGroup.DELETE("/:wordId", vocabController.DeleteWord)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/vocab/1", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Word ID Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/vocab/abc", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Word Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/vocab/999", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Attempt to Delete Other User's Word", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/vocab/2", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success Delete Word", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/vocab/1", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// check that user doesn't have any words now
		checkW := httptest.NewRecorder()
		checkReq, _ := http.NewRequest("GET", "/api/v1/vocab", nil)
		checkReq.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(checkW, checkReq)

		var response []resp.VocabResponse
		err := json.Unmarshal(checkW.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
	})
}
