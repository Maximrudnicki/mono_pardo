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

func TestGetWords(t *testing.T) {
	env, _ := tests.NewTestEnv(t)
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
				UserId:          1,
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

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/vocab", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Success Get Words", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/vocab", nil)

		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []resp.VocabResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		words := response
		assert.Len(t, words, 2)

		assert.Equal(t, "hello", words[0].Word)
		assert.Equal(t, "planet earth", words[1].Definition)
		assert.Equal(t, false, words[0].IsLearned)
		assert.WithinDuration(t, createdAt, words[0].CreatedAt, time.Second)
	})
}
