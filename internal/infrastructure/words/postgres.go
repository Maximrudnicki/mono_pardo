package words

import (
	"errors"

	domain "mono_pardo/internal/domain/words"

	"gorm.io/gorm"
)

type repositoryImpl struct {
	Db *gorm.DB
}

func NewPostgresRepositoryImpl(Db *gorm.DB) domain.Repository {
	return &repositoryImpl{Db: Db}
}

func (r *repositoryImpl) Delete(wordId int) {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Delete(&word)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (r *repositoryImpl) FindByUserId(userId int) ([]domain.Word, error) {
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

func (r *repositoryImpl) FindById(wordId int) (domain.Word, error) {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Find(&word)
	if result != nil {
		return word, nil
	} else {
		return word, errors.New("word is not found")
	}
}

func (r *repositoryImpl) Save(word domain.Word) error {
	result := r.Db.Create(&word)
	if result.Error != nil {
		return errors.New("cannot save word")
	}
	return nil
}

func (r *repositoryImpl) Update(wordId int, updates map[string]interface{}) error {
	var word domain.Word

	if len(updates) == 0 {
		return nil
	}

	result := r.Db.Model(&word).Where("id = ?", wordId).Updates(updates)
	if result.Error != nil {
		return errors.New("cannot update word")
	}

	return nil
}

func (r *repositoryImpl) IsOwnerOfWord(userId int, wordId int) (bool, error) {
	var word domain.Word
	result := r.Db.Where("id = ?", wordId).Find(&word)
	if result.Error != nil {
		return false, result.Error
	}

	return word.UserId == userId, nil
}

// func (r *repositoryImpl) Add(word domain.Word) (int, error) {
// 	result := r.Db.Create(&word)
// 	if result.Error != nil {
// 		return 0, errors.New("cannot add word")
// 	}
// 	return word.Id, nil
// }

func (r *repositoryImpl) ManageTrainings(res bool, training string, wordId int) error {
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

// func (r *repositoryImpl) UpdateStatus(word domain.Word) error {
// 	result := r.Db.Model(&word).Where("id = ?", word.Id).Update("is_learned", word.IsLearned)
// 	if result.Error != nil {
// 		return errors.New("cannot update status")
// 	}
// 	return nil
// }
