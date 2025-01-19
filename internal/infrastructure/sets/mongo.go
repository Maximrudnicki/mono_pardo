package sets

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "mono_pardo/internal/domain/sets"
	"mono_pardo/pkg/data/request"
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

func NewMongoRepositoryImpl(collection *mongo.Collection) domain.Repository {
	return &repositoryImpl{collection: collection}
}

type repositoryImpl struct {
	collection *mongo.Collection
}

func (r *repositoryImpl) AddToList(ctx context.Context, wordSetId string, words []int) error {
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return errors.New("cannot parse word set ID")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$addToSet": bson.M{"words": bson.M{"$each": words}}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("cannot add word to word set")
	}

	if res.MatchedCount == 0 {
		return errors.New("word set not found")
	}

	return nil
}

func (r *repositoryImpl) Delete(ctx context.Context, wordSetId string) error {
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return errors.New("cannot parse word set ID")
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.New("cannot delete word set")
	}

	if res.DeletedCount == 0 {
		return errors.New("word set was not found")
	}

	return nil
}

func (r *repositoryImpl) FindById(ctx context.Context, wordSetId string) (domain.WordSet, error) {
	var data domain.WordSet
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return data, errors.New("cannot parse word set ID")
	}
	filter := bson.M{"_id": oid}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Decode(&data); err != nil {
		return data, errors.New("cannot find word set with specified ID")
	}

	return data, nil
}

func (r *repositoryImpl) FindByUserId(ctx context.Context, userId int) ([]domain.WordSet, error) {
	filter := bson.M{"user_id": userId}

	res, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("cannot find word sets by user's id")
	}
	defer res.Close(ctx)

	var sets []domain.WordSet

	for res.Next(ctx) {
		var set domain.WordSet

		if err := res.Decode(&set); err != nil {
			return nil, err
		}

		sets = append(sets, set)
	}

	if err := res.Err(); err != nil {
		log.Printf("Cannot find word sets by user's id. Error: %v\n", err.Error())
		return nil, errors.New("cannot find word sets by user's id")
	}

	return sets, nil
}

func (r *repositoryImpl) RemoveFromList(ctx context.Context, wordSetId string, words []int) error {
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return errors.New("cannot parse word set ID")
	}

	update := bson.M{"$pull": bson.M{"words": bson.M{"$in": words}}}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return errors.New("cannot remove words")
	}

	if res.MatchedCount == 0 {
		return errors.New("word set not found")
	}

	return nil
}

func (r *repositoryImpl) Save(ctx context.Context, wordSet domain.WordSet) error {
	if _, err := r.collection.InsertOne(ctx, wordSet); err != nil {
		return errors.New("cannot create new word set")
	}

	return nil
}

func (r *repositoryImpl) Update(ctx context.Context, wordSetId string, updates []request.FieldUpdate) error {
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return errors.New("cannot parse word set ID")
	}

	updateFields := bson.M{}

	for _, update := range updates {
		switch update.Field {
		case "words", "_id", "created_at":
			continue
		default:
			updateFields[update.Field] = update.Value
		}
	}

	if len(updateFields) == 0 {
		return errors.New("no valid fields to update")
	}

	update := bson.M{"$set": updateFields}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)

	if err != nil {
		return errors.New("cannot update word set")
	}

	if result.MatchedCount == 0 {
		return errors.New("word set not found")
	}

	return nil
}

func (r *repositoryImpl) IsOwnerOfWordSet(ctx context.Context, userId int, wordSetId string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(wordSetId)
	if err != nil {
		return false, errors.New("cannot parse word set ID")
	}

	filter := bson.M{"_id": oid, "user_id": userId}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.New("cannot check word set ownership")
	}

	return count > 0, nil
}
