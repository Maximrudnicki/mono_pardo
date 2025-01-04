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

	// At the moment, we don't handle status change in case of all trainings are true
	UpdateWord(updateWordRequest request.UpdateWordRequest) error
	UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error
}

type Repository interface {
	// Add(word Word) (int, error)
	Save(word Word) error
	Update(word request.WordUpdate) error
	UpdateStatus(word Word) error
	Delete(wordId int)
	FindByUserId(userId int) ([]Word, error)
	FindById(wordId int) (Word, error)

	// utils
	IsOwnerOfWord(userId int, wordId int) (bool, error)
}
