package words

import (
	"bytes"
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
	"mono_pardo/pkg/data/request"
	resp "mono_pardo/pkg/data/response"
	"mono_pardo/tests"
)

func TestCreateWord(t *testing.T) {
	env := tests.NewTestEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	mockAuthService := &MockAuthService{}
	mockAuthService.On("GetUserId", "test-token").Return(1, nil)
	mockAuthService.On("GetUserId", "test-token-user2").Return(2, nil)
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
	vocabGroup.POST("", vocabController.CreateWord)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Input", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Definition", func(t *testing.T) {
		payload := map[string]interface{}{
			"word": "test",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Word", func(t *testing.T) {
		payload := map[string]interface{}{
			"definition": "test definition",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Duplicate Word for Same User", func(t *testing.T) {
		payload := request.CreateWordRequest{
			Word:       "hello",
			Definition: "new definition",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Same Word Different User", func(t *testing.T) {
		payload := request.CreateWordRequest{
			Word:       "hello",
			Definition: "greeting in English",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token-user2")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Success Create Word", func(t *testing.T) {
		payload := request.CreateWordRequest{
			Word:       "newword",
			Definition: "brand new word definition",
		}
		jsonData, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/vocab", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		checkW := httptest.NewRecorder()
		checkReq, _ := http.NewRequest("GET", "/api/v1/vocab", nil)
		checkReq.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(checkW, checkReq)

		var response []resp.VocabResponse
		err := json.Unmarshal(checkW.Body.Bytes(), &response)
		assert.NoError(t, err)

		found := false
		for _, word := range response {
			if word.Word == "newword" {
				found = true
				assert.Equal(t, "brand new word definition", word.Definition)
				break
			}
		}
		assert.True(t, found, "Created word not found in response")
	})
}
