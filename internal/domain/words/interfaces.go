package words

import (
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

type Service interface {
	CreateWord(createWordRequest request.CreateWordRequest) error
	DeleteWord(deleteWordRequest request.DeleteWordRequest) error
	GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error)
	FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error)
	UpdateWord(token string, wordId int, updates map[string]interface{}) error  // replace with PATCH-like method
	UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error
	ManageTrainings(manageTrainingsRequest request.ManageTrainingsRequest) error
}

type Repository interface {
	// Add(word Word) (int, error)
	Save(word Word) error
	Update(wordId int, updates map[string]interface{}) error
	Delete(wordId int)
	FindByUserId(userId int) ([]Word, error)
	FindById(wordId int) (Word, error)
	ManageTrainings(res bool, training string, wordId int) error

	// utils
	IsOwnerOfWord(userId int, wordId int) (bool, error)
}
