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
	// resp "mono_pardo/pkg/data/response"
	"mono_pardo/tests"
)

func TestUpdateWord(t *testing.T) {
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
				WordAudio:       true,
				CreatedAt:       createdAt,
			},
			{
				UserId:          1,
				Word:            "test2",
				Definition:      "test def",
				Cards:           true,
				WordTranslation: true,
				Constructor:     true,
				WordAudio:       true,
				IsLearned:       true,
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
	vocabGroup.PATCH("", vocabController.UpdateWord)

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/vocab", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Input", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/vocab", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Values", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 1,
				"updates": []map[string]interface{}{
					{
						"field": "word",
						"value": "",
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 2,
				"updates": []map[string]interface{}{
					{
						"field": "definition",
						"value": "new definition",
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Field Name", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 1,
				"updates": []map[string]interface{}{
					{
						"field": "nonexistent_field",
						"value": true,
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Value Type", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 1,
				"updates": []map[string]interface{}{
					{
						"field": "cards",
						"value": "not a boolean",
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Training Failure", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 3,
				"updates": []map[string]interface{}{
					{
						"field": "word_audio",
						"value": false,
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify IsLearned was set to false
		var updatedWord wordsDomain.Word
		err := env.DB.DB.First(&updatedWord, 3).Error
		assert.NoError(t, err)
		assert.False(t, updatedWord.IsLearned)
	})

	t.Run("Training Completion Success", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 3,
				"updates": []map[string]interface{}{
					{
						"field": "word_audio",
						"value": true,
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify IsLearned was set to true
		var updatedWord wordsDomain.Word
		err := env.DB.DB.First(&updatedWord, 3).Error
		assert.NoError(t, err)
		assert.True(t, updatedWord.IsLearned)
	})

	t.Run("Successful Update", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"id": 1,
				"updates": []map[string]interface{}{
					{
						"field": "definition",
						"value": "new greeting",
					},
					{
						"field": "word_audio",
						"value": true,
					},
				},
			},
		}
		w := httptest.NewRecorder()
		req := createJSONRequest(t, "PATCH", "/api/v1/vocab", payload)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the changes
		var updatedWord wordsDomain.Word
		err := env.DB.DB.First(&updatedWord, 1).Error
		assert.NoError(t, err)
		assert.False(t, updatedWord.IsLearned)
		assert.NotNil(t, updatedWord)
		assert.Equal(t, "new greeting", updatedWord.Definition)
		assert.True(t, updatedWord.WordAudio)
	})
}

// Helper function to create JSON requests
func createJSONRequest(t *testing.T, method, url string, payload interface{}) *http.Request {
	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	return req
}
