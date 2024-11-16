package words

import (
	"errors"

	domain "mono_pardo/internal/domain/words"

	"gorm.io/gorm"
)

type RepositoryImpl struct {
	Db *gorm.DB
}

func NewPostgresRepositoryImpl(Db *gorm.DB) domain.Repository {
	return &RepositoryImpl{Db: Db}
}

func (r *RepositoryImpl) ManageTrainings(res bool, training string, wordId int) error {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Find(&word)
	if result.Error != nil {
		return errors.New("cannot find word")
	}

	switch training {
	case "cards":
		word.Cards = res
	case "word_translation":
		word.WordTranslation = res
	case "constructor":
		word.Constructor = res
	case "word_audio":
		word.WordAudio = res
	default:
		return errors.New("unknow training")
	}

	if word.Cards && word.WordTranslation && word.Constructor && word.WordAudio {
		word.IsLearned = true
	} else {
		word.IsLearned = false
	}

	result = r.Db.Save(&word)
	if result.Error != nil {
		return errors.New("cannot save word")
	}
	return nil
}

func (r *RepositoryImpl) Delete(wordId int) {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Delete(&word)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (r *RepositoryImpl) FindByUserId(userId int) ([]domain.Word, error) {
	var words []domain.Word
	result := r.Db.Where("user_id = ?", userId).Find(&words)
	// Should return empty slice in case if user exists.
	// It's made for case if user exists but have not added any words yet
	if result != nil {
		return words, nil
	} else {
		return nil, errors.New("words is not found")
	}
}

func (r *RepositoryImpl) FindById(wordId int) (domain.Word, error) {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Find(&word)
	if result != nil {
		return word, nil
	} else {
		return word, errors.New("word is not found")
	}
}

func (r *RepositoryImpl) Add(word domain.Word) (int, error) {
	result := r.Db.Create(&word)
	if result.Error != nil {
		return 0, errors.New("cannot add word")
	}
	return word.Id, nil
}

func (r *RepositoryImpl) Save(word domain.Word) error {
	result := r.Db.Create(&word)
	if result.Error != nil {
		return errors.New("cannot save word")
	}
	return nil
}

func (r *RepositoryImpl) Update(word domain.Word) error {
	var updatedWord = &domain.Word{
		Definition: word.Definition,
	}

	result := r.Db.Model(&word).Where("id = ?", word.Id).Updates(updatedWord)
	if result.Error != nil {
		return errors.New("cannot update word")
	}
	return nil
}

func (r *RepositoryImpl) UpdateStatus(word domain.Word) error {
	result := r.Db.Model(&word).Where("id = ?", word.Id).Update("is_learned", word.IsLearned)
	if result.Error != nil {
		return errors.New("cannot update status")
	}
	return nil
}

func (r *RepositoryImpl) IsOwnerOfWord(userId int, wordId int) bool {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Find(&word)
	if result.Error != nil {
		panic(result.Error)
	}

	if word.UserId == userId {
		return true
	} else {
		return false
	}
}
