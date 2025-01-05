package words

import (
	"errors"
	"fmt"

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

func (r *repositoryImpl) Delete(wordId int) error {
	var word domain.Word

	if err := r.Db.Where("id = ?", wordId).Delete(&word).Error; err != nil {
		return fmt.Errorf("cannot delete word: %d", wordId)
	}

	return nil
}

func (r *repositoryImpl) FindByUserId(userId int) ([]domain.Word, error) {
	var words []domain.Word

	if err := r.Db.Where("user_id = ?", userId).Find(&words).Error; err != nil {
		return nil, errors.New("words is not found")
	}

	// Should return empty slice in case if user exists but have not added any words yet.
	return words, nil
}

func (r *repositoryImpl) FindById(wordId int) (domain.Word, error) {
	var word domain.Word

	if err := r.Db.Where("id = ?", wordId).Find(&word).Error; err != nil {
		return word, fmt.Errorf("cannot find word with id: %d", wordId)
	}

	return word, nil
}

func (r *repositoryImpl) Save(word domain.Word) error {
	if err := r.Db.Create(&word).Error; err != nil {
		return errors.New("cannot save word")
	}

	return nil
}

func (r *repositoryImpl) Update(wordUpdate request.WordUpdate) error {
	var word domain.Word

	updateMap := utils.ConvertFieldUpdatesToMap(wordUpdate.Updates)

	err := r.Db.Model(&word).Where("id = ?", wordUpdate.WordId).Updates(updateMap).Error
	if err != nil {
		return fmt.Errorf("cannot update word: %d", wordUpdate.WordId)
	}

	return nil
}

func (r *repositoryImpl) IsOwnerOfWord(userId int, wordId int) (bool, error) {
	var word domain.Word

	if err := r.Db.Where("id = ?", wordId).Find(&word).Error; err != nil {
		return false, fmt.Errorf("cannot check who is owner of the word: %d", wordId)
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
