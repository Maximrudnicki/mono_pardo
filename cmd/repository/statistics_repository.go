package repository

import (
	"context"
	"errors"
	"mono_pardo/cmd/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsRepository interface {
	AddWordToStatistics(ctx context.Context, statId string, word int) error
	CreateStatistics(ctx context.Context, stat model.Statistics) error
	GetStatistics(ctx context.Context, statId string) (model.Statistics, error)
	GetId(ctx context.Context, groupId string, studentId int) (string, error)
	DeleteStatistics(ctx context.Context, statId string) error
}

type StatisticsRepositoryImpl struct {
	collection *mongo.Collection
}

// AddWordsToStatistics implements StatisticsRepository.
func (s *StatisticsRepositoryImpl) AddWordToStatistics(ctx context.Context, statId string, word int) error {
	oid, err := primitive.ObjectIDFromHex(statId)
	if err != nil {
		return errors.New("cannot parse stat ID")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$addToSet": bson.M{"words": word}}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("cannot add words to statistics")
	}

	return nil
}

// CreateStatistics implements StatisticsRepository.
func (s *StatisticsRepositoryImpl) CreateStatistics(ctx context.Context, stat model.Statistics) error {
	stat.Words = []int{}
	_, err := s.collection.InsertOne(ctx, stat)
	if err != nil {
		return errors.New("cannot create stats")
	}

	return nil
}

// DeleteStatistics implements StatisticsRepository.
func (s *StatisticsRepositoryImpl) DeleteStatistics(ctx context.Context, statId string) error {
	oid, err := primitive.ObjectIDFromHex(statId)
	if err != nil {
		return errors.New("cannot get id to delete stats")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.New("cannot delete statistics")
	}

	if res.DeletedCount == 0 {
		return errors.New("statistics was not found")
	}

	return nil
}

// GetStatistics implements StatisticsRepository.
func (s *StatisticsRepositoryImpl) GetStatistics(ctx context.Context, statId string) (model.Statistics, error) {
	var data model.Statistics
	oid, err := primitive.ObjectIDFromHex(statId)
	if err != nil {
		return data, errors.New("cannot parse stat ID")
	}
	filter := bson.M{"_id": oid}

	res := s.collection.FindOne(ctx, filter)
	if err := res.Decode(&data); err != nil {
		return data, errors.New("cannot find statistics with specified ID")
	}

	return data, nil
}

// GetStatistics implements StatisticsRepository.
func (s *StatisticsRepositoryImpl) GetId(ctx context.Context, groupId string, studentId int) (string, error) {
	oid, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return "", errors.New("cannot parse group ID")
	}

	filter := bson.M{
		"group_id":   oid,
		"student_id": studentId,
	}

	var data model.Statistics
	res := s.collection.FindOne(ctx, filter)
	if err := res.Decode(&data); err != nil {
		return "", errors.New("cannot find statistics")
	}

	return data.Id.Hex(), nil
}

func NewStatisticsRepositoryImpl(collection *mongo.Collection) StatisticsRepository {
	return &StatisticsRepositoryImpl{collection: collection}
}
