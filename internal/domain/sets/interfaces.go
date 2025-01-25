package sets

import (
	"context"

	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"
)

type Service interface {
	CreateSet(createSetRequest request.CreateSetRequest) error
	GetSets(getSetsRequest request.GetSetsRequest) ([]response.SetResponse, error)
	GetSet(getSetRequest request.GetSetRequest) (response.SetResponse, error)
	UpdateSet(updateSetRequest request.UpdateSetRequest) error
	DeleteSet(deleteSetRequest request.DeleteSetRequest) error

	AddWord(addWordRequest request.AddWordRequest) error
	RemoveWord(removeWordRequest request.RemoveWordRequest) error

	validateSetUpdates(updates []request.FieldUpdate) error
}

type Repository interface {
	Save(ctx context.Context, wordSet WordSet) error
	FindById(ctx context.Context, wordSetId string) (WordSet, error)
	FindByUserId(ctx context.Context, userId int) ([]WordSet, error)
	Update(ctx context.Context, wordSetId string, updates []request.FieldUpdate) error
	Delete(ctx context.Context, wordSetId string) error

	AddToList(ctx context.Context, wordSetId string, words []int) error
	RemoveFromList(ctx context.Context, wordSetId string, words []int) error

	// utils
	IsOwnerOfWordSet(ctx context.Context, userId int, wordSetId string) (bool, error)
}
