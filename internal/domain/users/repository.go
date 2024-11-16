package users

type UsersRepository interface {
	Save(user User) error
	Delete(usersId int)
	FindById(usersId int) (User, error)
	FindAll() []User
	FindByEmail(email string) (User, error)
}
