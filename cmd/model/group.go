package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	TeacherId int                `bson:"teacher_id"`
	Students  []int              `bson:"students"`
}
