package sets

import (
	"context"

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
	newSet, err := NewWordSet(createSetRequest.Name, createSetRequest.UserId)
	if err != nil {
		return err
	}

	if err = s.Repository.Save(context.Background(), *newSet); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) DeleteSet(deleteSetRequest request.DeleteSetRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) GetSet(getSetRequest request.GetSetRequest) (response.SetResponse, error) {
	panic("unimplemented")
}

func (s *serviceImpl) GetSets(getSetsRequest request.GetSetsRequest) ([]response.SetResponse, error) {
	var setsResponse []response.SetResponse

	sets, err := s.Repository.FindByUserId(context.Background(), getSetsRequest.UserId)
	if err != nil {
		return nil, err
	}

	for _, set := range sets {
		setsResponse = append(setsResponse, response.SetResponse{
			Id:        set.Id.Hex(),
			Name:      set.Name,
			CreatedAt: set.CreatedAt,
			Words:     set.Words,
		})
	}

	return setsResponse, nil
}

func (s *serviceImpl) RemoveWord(removeWordRequest request.RemoveWordRequest) error {
	panic("unimplemented")
}

func (s *serviceImpl) UpdateSet(updateSetRequest request.UpdateSetRequest) error {
	panic("unimplemented")
}
