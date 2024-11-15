package words

import "mono_pardo/internal/model"

type WordRepository interface {
	Add(word model.Word) (int, error)
	Save(word model.Word) error
	Update(word model.Word) error
	UpdateStatus(word model.Word) error
	Delete(wordId int)
	FindByUserId(userId int) ([]model.Word, error)
	FindById(wordId int) (model.Word, error)
	ManageTrainings(res bool, training string, wordId int) error

	// utils
	IsOwnerOfWord(userId int, wordId int) bool
}
