package users

import (
	"errors"
	"mono_pardo/internal/domain/model"
	"mono_pardo/internal/domain/users"

	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewPostgresRepositoryImpl(Db *gorm.DB) users.UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

// Save implements UsersRepository
func (u *UsersRepositoryImpl) Save(user model.User) error {
	result := u.Db.Create(&user)
	if result.Error != nil {
		return errors.New("please use different email")
	}
	return nil
}

// Delete implements UsersRepository
func (u *UsersRepositoryImpl) Delete(usersId int) {
	var user model.User
	result := u.Db.Where("id = ?", usersId).Delete(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

// FindAll implements UsersRepository
func (u *UsersRepositoryImpl) FindAll() []model.User {
	var user []model.User
	results := u.Db.Find(&user)
	if results.Error != nil {
		panic(results.Error)
	}
	return user
}

// FindById implements UsersRepository
func (u *UsersRepositoryImpl) FindById(userId int) (model.User, error) {
	var user model.User
	result := u.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

// FindByUsername implements UsersRepository
func (u *UsersRepositoryImpl) FindByEmail(email string) (model.User, error) {
	var user model.User
	result := u.Db.First(&user, "email = ?", email)

	if result.Error != nil {
		return user, errors.New("invalid email or Password")
	}
	return user, nil
}
