package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistics struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Group_id  primitive.ObjectID `bson:"group_id"`
	TeacherId int                `bson:"teacher_id"`
	StudentId int                `bson:"student_id"`
	Words     []int              `bson:"words"`
}
