package tests

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	setsDomain "mono_pardo/internal/domain/sets"
	usersDomain "mono_pardo/internal/domain/users"
	wordsDomain "mono_pardo/internal/domain/words"
)

type UserFixture struct {
	Users []usersDomain.User
}

func (f *UserFixture) Setup(db *gorm.DB) error {
	return db.Create(&f.Users).Error
}

func (f *UserFixture) Teardown(db *gorm.DB) error {
	return db.Unscoped().Delete(&f.Users).Error
}

type WordFixture struct {
	Words []wordsDomain.Word
}

func (f *WordFixture) Setup(db *gorm.DB) error {
	return db.Create(&f.Words).Error
}

func (f *WordFixture) Teardown(db *gorm.DB) error {
	return db.Unscoped().Delete(&f.Words).Error
}

type WordSetFixture struct {
	Sets []setsDomain.WordSet
}

func (f *WordSetFixture) Setup(db *mongo.Database) error {
	if len(f.Sets) == 0 {
		return nil
	}

	collection := db.Collection("sets")

	docs := make([]interface{}, len(f.Sets))
	for i, set := range f.Sets {
		docs[i] = set
	}

	_, err := collection.InsertMany(context.Background(), docs)
	return err
}

func (f *WordSetFixture) Teardown(db *mongo.Database) error {
	collection := db.Collection("sets")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	return err
}
