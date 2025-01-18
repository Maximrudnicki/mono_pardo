package sets

import (
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
}

type Repository interface {
	Save(wordSet WordSet) error
	FindById(groupId string) (WordSet, error)
	FindByUserId(userId int) ([]WordSet, error)
	Update(groupId string, updates []request.FieldUpdate) error
	Delete(groupId string) error

	AddToList(groupId string, wordId int) error
	RemoveFromList(groupId string, wordId int) error
}
