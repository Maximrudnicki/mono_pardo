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
	UpdateWord(updateWordRequest request.UpdateWordRequest) error
	updateWordStatus(wordId int) error
}

type Repository interface {
	// Add(word Word) (int, error)
	Save(word Word) error
	Update(word request.WordUpdate) error
	Delete(wordId int)
	FindByUserId(userId int) ([]Word, error)
	FindById(wordId int) (Word, error)

	// utils
	IsOwnerOfWord(userId int, wordId int) (bool, error)
}
