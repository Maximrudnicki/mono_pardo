package sets

import (
	domain "mono_pardo/internal/domain/sets"
	"mono_pardo/pkg/data/request"

	"go.mongodb.org/mongo-driver/mongo"
)

/*
	Use Mongo DB because of convenient $push $pull system
	to manage lists of objects. In our case - lists of words.

	We don't to make any joins because front-end manages this on client's side.
	FE should fetch all user's words and check the IDs of them here - in the mongo collection.
	Then FE shows set with the words in it to the user.

	We don't care too much about consistency, so FE can just skip IDs that is missing in user's vocab.
	The same in case of validation of foreign keys. FE just skips missing words
*/

func NewMongoRepositoryImpl() domain.Repository {
	return &repositoryImpl{}
}

type repositoryImpl struct {
	collection *mongo.Collection
}

func (r *repositoryImpl) AddToList(groupId string, wordId int) error {
	panic("unimplemented")
}

func (r *repositoryImpl) Delete(groupId string) error {
	panic("unimplemented")
}

func (r *repositoryImpl) FindById(groupId string) (domain.WordSet, error) {
	panic("unimplemented")
}

func (r *repositoryImpl) FindByUserId(userId int) ([]domain.WordSet, error) {
	panic("unimplemented")
}

func (r *repositoryImpl) RemoveFromList(groupId string, wordId int) error {
	panic("unimplemented")
}

func (r *repositoryImpl) Save(wordSet domain.WordSet) error {
	panic("unimplemented")
}

func (r *repositoryImpl) Update(groupId string, updates []request.FieldUpdate) error {
	panic("unimplemented")
}
