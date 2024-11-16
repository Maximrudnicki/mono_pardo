package words

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
