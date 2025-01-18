package sets

import domain "mono_pardo/internal/domain/sets"

/*
	Use Mongo DB because of convinient $push $pull system
	to manage lists of objects. In our case - lists of words.

	We don't to make any joins because front-end manages this on client's side.
	FE should fetch all user's words and check the IDs of them here - in the mongo collection.
	Then FE shows set with the words in it to the user.

	We don't care too much about consistency, so FE can just skip IDs that is missing in user's vocab.
	The same in case of validation of foreign keys. FE just skips missing words
*/

type repositoryImpl struct{}

func NewMongoRepositoryImpl() domain.Repository {
	return &repositoryImpl{}
}
