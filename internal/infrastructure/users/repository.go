package users

import "mono_pardo/internal/domain/model"

type UsersRepository interface {
	Save(user model.User) error
	Delete(usersId int)
	FindById(usersId int) (model.User, error)
	FindAll() []model.User
	FindByEmail(email string) (model.User, error)
}
