package sets

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WordSet struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UserId    int                `bson:"user_id"`
	Words     []int              `bson:"words"` // list of IDs
}

func NewWordSet(name string, userId int) (*WordSet, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if userId <= 0 {
		return nil, errors.New("userId must be a positive integer")
	}

	return &WordSet{
		Id:        primitive.NewObjectID(),
		Name:      name,
		CreatedAt: time.Now(),
		UserId:    userId,
		Words:     []int{},
	}, nil
}
