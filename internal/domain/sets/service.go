package sets

import "github.com/go-playground/validator"

type serviceImpl struct {
	Validate   *validator.Validate
	Repository Repository
}

func NewServiceImpl(validate *validator.Validate, repository Repository) Service {
	return &serviceImpl{
		Validate:   validate,
		Repository: repository,
	}
}
