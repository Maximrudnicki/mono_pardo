package users

import (
	"errors"

	domain "mono_pardo/internal/domain/users"

	"gorm.io/gorm"
)

type RepositoryImpl struct {
	Db *gorm.DB
}

func NewPostgresRepositoryImpl(Db *gorm.DB) domain.Repository {
	return &RepositoryImpl{Db: Db}
}

// Save implements UsersRepository
func (r *RepositoryImpl) Save(user domain.User) error {
	result := r.Db.Create(&user)
	if result.Error != nil {
		return errors.New("please use different email")
	}
	return nil
}

// Delete implements UsersRepository
func (r *RepositoryImpl) Delete(usersId int) {
	var user domain.User
	result := r.Db.Where("id = ?", usersId).Delete(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

// FindAll implements UsersRepository
func (r *RepositoryImpl) FindAll() []domain.User {
	var user []domain.User
	results := r.Db.Find(&user)
	if results.Error != nil {
		panic(results.Error)
	}
	return user
}

// FindById implements UsersRepository
func (r *RepositoryImpl) FindById(userId int) (domain.User, error) {
	var user domain.User
	result := r.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

// FindByUsername implements UsersRepository
func (r *RepositoryImpl) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	result := r.Db.First(&user, "email = ?", email)

	if result.Error != nil {
		return user, errors.New("invalid email or Password")
	}
	return user, nil
}
