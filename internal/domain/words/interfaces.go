package words

import (
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

type VocabService interface {
	CreateWord(createWordRequest request.CreateWordRequest) error
	DeleteWord(deleteWordRequest request.DeleteWordRequest) error
	GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error)
	FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error)
	UpdateWord(updateWordRequest request.UpdateWordRequest) error
	UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error
	ManageTrainings(manageTrainingsRequest request.ManageTrainingsRequest) error
}

type WordRepository interface {
	Add(word Word) (int, error)
	Save(word Word) error
	Update(word Word) error
	UpdateStatus(word Word) error
	Delete(wordId int)
	FindByUserId(userId int) ([]Word, error)
	FindById(wordId int) (Word, error)
	ManageTrainings(res bool, training string, wordId int) error

	// utils
	IsOwnerOfWord(userId int, wordId int) bool
}
