package sets

import (
	"context"
	"errors"

	"github.com/go-playground/validator"

	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

const (
	ErrAccessForbidden = "access forbidden"
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
	if isOwner, err := s.Repository.IsOwnerOfWordSet(
		context.Background(), addWordRequest.UserId, addWordRequest.WordSetId); err != nil {
		return err
	} else if !isOwner {
		return errors.New(ErrAccessForbidden)
	}

	if err := s.Repository.AddToList(
		context.Background(), addWordRequest.WordSetId, addWordRequest.Words); err != nil {
		return err
	}

	return nil
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
	if isOwner, err := s.Repository.IsOwnerOfWordSet(
		context.Background(), deleteSetRequest.UserId, deleteSetRequest.WordSetId); err != nil {
		return err
	} else if !isOwner {
		return errors.New(ErrAccessForbidden)
	}

	if err := s.Repository.Delete(context.Background(), deleteSetRequest.WordSetId); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) GetSet(getSetRequest request.GetSetRequest) (response.SetResponse, error) {
	var setResponse response.SetResponse

	if isOwner, err := s.Repository.IsOwnerOfWordSet(
		context.Background(), getSetRequest.UserId, getSetRequest.WordSetId); err != nil {
		return setResponse, err
	} else if !isOwner {
		return setResponse, errors.New("cannot find word set with specified ID")
	}

	set, err := s.Repository.FindById(context.Background(), getSetRequest.WordSetId)
	if err != nil {
		return setResponse, err
	}

	return response.SetResponse{
		Id:        set.Id.Hex(),
		Name:      set.Name,
		CreatedAt: set.CreatedAt,
		Words:     set.Words,
	}, nil
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
	if isOwner, err := s.Repository.IsOwnerOfWordSet(
		context.Background(), removeWordRequest.UserId, removeWordRequest.WordSetId); err != nil {
		return err
	} else if !isOwner {
		return errors.New(ErrAccessForbidden)
	}

	if err := s.Repository.RemoveFromList(
		context.Background(), removeWordRequest.WordSetId, removeWordRequest.Words); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) UpdateSet(updateSetRequest request.UpdateSetRequest) error {
	if isOwner, err := s.Repository.IsOwnerOfWordSet(
		context.Background(), updateSetRequest.UserId, updateSetRequest.WordSetId); err != nil {
		return err
	} else if !isOwner {
		return errors.New(ErrAccessForbidden)
	}

	if err := s.Repository.Update(
		context.Background(), updateSetRequest.WordSetId, updateSetRequest.Updates); err != nil {
		return err
	}

	return nil
}
