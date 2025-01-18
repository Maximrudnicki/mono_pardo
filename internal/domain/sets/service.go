package sets

import (
	"github.com/go-playground/validator"

	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

func NewServiceImpl(validate *validator.Validate, repository Repository) Service {
	return &serviceImpl{
		Validate:   validate,
		Repository: repository,
	}
}

type serviceImpl struct {
	Validate   *validator.Validate
	Repository Repository
}

func (s *serviceImpl) AddWord(addWordRequest request.AddWordRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) CreateSet(createSetRequest request.CreateSetRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) DeleteSet(deleteSetRequest request.DeleteSetRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) GetSet(getSetRequest request.GetSetRequest) (response.SetResponse, error) {
	panic("unimplemented")
}

func (s *serviceImpl) GetSets(getSetsRequest request.GetSetsRequest) ([]response.SetResponse, error) {
	panic("unimplemented")
}

func (s *serviceImpl) RemoveWord(removeWordRequest request.RemoveWordRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) UpdateSet(updateSetRequest request.UpdateSetRequest) error {
	panic("unimplemented")
}
