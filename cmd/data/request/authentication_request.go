package request

type CreateUserRequest struct {
	Username string `validate:"required,min=2,max=100" json:"username"`
	Email    string `validate:"required,min=2,max=100" json:"email"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}

type LoginRequest struct {
	Email string `validate:"required,max=200,min=2" json:"email"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}
