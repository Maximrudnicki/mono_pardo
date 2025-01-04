package words

import (
	"errors"

	domain "mono_pardo/internal/domain/words"
	"mono_pardo/internal/utils"
	"mono_pardo/pkg/data/request"

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

func (r *repositoryImpl) Update(wordUpdate request.WordUpdate) error {
	var word domain.Word

	updateMap := utils.ConvertFieldUpdatesToMap(wordUpdate.Updates)

	result := r.Db.Model(&word).Where("id = ?", wordUpdate.WordId).Updates(updateMap)
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
