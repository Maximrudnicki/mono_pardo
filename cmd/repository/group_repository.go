package repository

import (
	"context"
	"errors"
	"mono_pardo/cmd/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepository interface {
	AddStudent(ctx context.Context, userId int, groupId string) error
	CreateGroup(ctx context.Context, group model.Group) error
	DeleteGroup(ctx context.Context, groupId string) error
	FindById(ctx context.Context, groupId string) (model.Group, error)
	FindByStudentId(ctx context.Context, userId int) ([]model.Group, error)
	FindByTeacherId(ctx context.Context, userId int) ([]model.Group, error)
	IsStudent(ctx context.Context, userId int, groupId string) bool
	IsTeacher(ctx context.Context, userId int, groupId string) bool
	RemoveStudent(ctx context.Context, userId int, groupId string) error
}

type GroupRepositoryImpl struct {
	collection *mongo.Collection
}

// AddStudent implements GroupRepository.
func (g *GroupRepositoryImpl) AddStudent(ctx context.Context, userId int, groupId string) error {
	oid, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return errors.New("cannot parse group ID")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$addToSet": bson.M{"students": userId}}

	res, err := g.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("cannot add student to group")
	}

	if res.MatchedCount == 0 {
		return errors.New("group not found")
	}

	return nil
}

// CreateGroup implements GroupRepository.
func (g *GroupRepositoryImpl) CreateGroup(ctx context.Context, group model.Group) error {
	group.Students = []int{} // if the list of students is not initialized, we cannot add students to it
	_, err := g.collection.InsertOne(ctx, group)
	if err != nil {
		return errors.New("cannot create group")
	}

	return nil
}

// DeleteGroup implements GroupRepository.
func (g *GroupRepositoryImpl) DeleteGroup(ctx context.Context, groupId string) error {
	oid, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return errors.New("cannot get id to delete group")
	}

	res, err := g.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.New("cannot delete group")
	}

	if res.DeletedCount == 0 {
		return errors.New("group was not found")
	}

	return nil
}

// FindById implements GroupRepository.
func (g *GroupRepositoryImpl) FindById(ctx context.Context, groupId string) (model.Group, error) {
	var data model.Group
	oid, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return data, errors.New("cannot parse id to delete group")
	}
	filter := bson.M{"_id": oid}

	res := g.collection.FindOne(ctx, filter)
	if err := res.Decode(&data); err != nil {
		return data, errors.New("cannot find group with specified ID")
	}

	return data, nil
}

// FindByStudentId implements GroupRepository.
func (g *GroupRepositoryImpl) FindByStudentId(ctx context.Context, userId int) ([]model.Group, error) {
	filter := bson.M{"students": userId}

	res, err := g.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("cannot find groups by student id")
	}
	defer res.Close(ctx)

	var groups []model.Group
	for res.Next(ctx) {
		var group model.Group
		err := res.Decode(&group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// FindByTeacherId implements GroupRepository.
func (g *GroupRepositoryImpl) FindByTeacherId(ctx context.Context, userId int) ([]model.Group, error) {
	filter := bson.M{"teacher_id": userId}

	res, err := g.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("cannot find groups by teacher id")
	}
	defer res.Close(ctx)

	var groups []model.Group
	for res.Next(ctx) {
		var group model.Group
		err := res.Decode(&group)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

// IsStudent implements GroupRepository.
func (g *GroupRepositoryImpl) IsStudent(ctx context.Context, userId int, groupId string) bool {
	group, err := g.FindById(ctx, groupId)
	if err != nil {
		return false
	}

	for _, studentId := range group.Students {
		if studentId == userId {
			return true
		}
	}

	return false
}

// IsTeacher implements GroupRepository.
func (g *GroupRepositoryImpl) IsTeacher(ctx context.Context, userId int, groupId string) bool {
	group, err := g.FindById(ctx, groupId)
	if err != nil {
		return false
	}

	if group.TeacherId == userId {
		return true
	} else {
		return false
	}
}

// RemoveStudent implements GroupRepository.
func (g *GroupRepositoryImpl) RemoveStudent(ctx context.Context, userId int, groupId string) error {
	group, err := g.FindById(ctx, groupId)
	if err != nil {
		return errors.New("cannot find group by id")
	}

	update := bson.M{"$pull": bson.M{"students": userId}}
	res, err := g.collection.UpdateOne(ctx, bson.M{"_id": group.Id}, update)
	if err != nil {
		return errors.New("cannot remove student")
	}

	if res.ModifiedCount == 0 {
		return errors.New("group not found or student not in group")
	}

	return nil
}

func NewGroupRepositoryImpl(collection *mongo.Collection) GroupRepository {
	return &GroupRepositoryImpl{collection: collection}
}
