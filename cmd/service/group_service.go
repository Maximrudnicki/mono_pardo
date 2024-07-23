package service

import (
	"context"
	"errors"
	"fmt"
	"mono_pardo/cmd/data/request"
	"mono_pardo/cmd/data/response"
	"mono_pardo/cmd/model"
	"mono_pardo/cmd/repository"

	"github.com/go-playground/validator"
)

type GroupService interface {
	AddStudent(addStudentRequest request.AddStudentRequest) error
	AddWordToUser(addWordToUserRequest request.AddWordToUserRequest) (response.AddWordToUserResponse, error)
	CreateGroup(createGroupRequest request.CreateGroupRequest) error
	DeleteGroup(deleteGroupRequest request.DeleteGroupRequest) error
	FindGroup(findGroupRequest request.FindGroupRequest) (response.GroupResponse, error)
	FindStudent(findStudentRequest request.FindStudentRequest) (response.StudentResponse, error)
	FindTeacher(findTeacherRequest request.FindTeacherRequest) (response.TeacherResponse, error)
	FindGroupsTeacher(findGroupsTeacherRequest request.FindGroupsTeacherRequest) ([]response.GroupResponse, error)
	FindGroupsStudent(findGroupsStudentRequest request.FindGroupsStudentRequest) ([]response.GroupResponse, error)
	GetStatistics(getStatisticsRequest request.GetStatisticsRequest) (response.StatisticsResponse, error)
	RemoveStudent(removeStudentRequest request.RemoveStudentRequest) error
}

type GroupServiceImpl struct {
	AuthenticationService AuthenticationService
	VocabService          VocabService
	Validate              *validator.Validate
	GroupRepository       repository.GroupRepository
	StatisticsRepository  repository.StatisticsRepository
}

// AddStudent implements GroupService.
func (g *GroupServiceImpl) AddStudent(addStudentRequest request.AddStudentRequest) error {
	userId, err := g.AuthenticationService.GetUserId(addStudentRequest.Token)
	if err != nil {
		return err
	}

	err = g.GroupRepository.AddStudent(context.Background(), userId, addStudentRequest.GroupId)
	if err != nil {
		return err
	}

	group, err := g.GroupRepository.FindById(context.Background(), addStudentRequest.GroupId)
	if err != nil {
		return err
	}

	stat := model.Statistics{
		Group_id:  group.Id,
		TeacherId: group.TeacherId,
		StudentId: userId,
	}

	err = g.StatisticsRepository.CreateStatistics(context.Background(), stat)
	if err != nil {
		return err
	}

	return nil
}

// AddWordToUser implements GroupService.
func (g *GroupServiceImpl) AddWordToUser(addWordToUserRequest request.AddWordToUserRequest) (response.AddWordToUserResponse, error) {
	teacherId, err := g.AuthenticationService.GetUserId(addWordToUserRequest.Token)
	if err != nil {
		return response.AddWordToUserResponse{}, err
	}

	IsStudent := g.GroupRepository.IsStudent(context.Background(), addWordToUserRequest.UserId, addWordToUserRequest.GroupId)
	IsTeacher := g.GroupRepository.IsTeacher(context.Background(), teacherId, addWordToUserRequest.GroupId)

	if IsStudent && IsTeacher {
		wordId, err := g.VocabService.AddWordToStudent(request.AddWordToStudentRequest{
			Word:       addWordToUserRequest.Word,
			Definition: addWordToUserRequest.Definition,
			UserId:     addWordToUserRequest.UserId,
		})
		if err != nil {
			return response.AddWordToUserResponse{}, err
		}
		statId, err := g.StatisticsRepository.GetId(
			context.Background(), addWordToUserRequest.GroupId, addWordToUserRequest.UserId)
		if err != nil {
			return response.AddWordToUserResponse{}, err
		}
		err = g.StatisticsRepository.AddWordToStatistics(context.Background(), statId, wordId)
		if err != nil {
			return response.AddWordToUserResponse{}, err
		}

		return response.AddWordToUserResponse{
			WordId: wordId,
		}, nil
	} else {
		return response.AddWordToUserResponse{}, errors.New("you are not allowed to delete the group")
	}

}

// CreateGroup implements GroupService.
func (g *GroupServiceImpl) CreateGroup(createGroupRequest request.CreateGroupRequest) error {
	userId, err := g.AuthenticationService.GetUserId(createGroupRequest.Token)
	if err != nil {
		return err
	}

	newGroup := model.Group{
		Title:     createGroupRequest.Title,
		TeacherId: userId,
	}

	err = g.GroupRepository.CreateGroup(context.Background(), newGroup)
	if err != nil {
		return fmt.Errorf("internal error: %v", err)
	}

	return nil
}

// DeleteGroup implements GroupService.
func (g *GroupServiceImpl) DeleteGroup(deleteGroupRequest request.DeleteGroupRequest) error {
	userId, err := g.AuthenticationService.GetUserId(deleteGroupRequest.Token)
	if err != nil {
		return err
	}

	if IsTeacher := g.GroupRepository.IsTeacher(context.Background(), userId, deleteGroupRequest.GroupId); IsTeacher {
		err = g.GroupRepository.DeleteGroup(context.Background(), deleteGroupRequest.GroupId)
		if err != nil {
			return fmt.Errorf("internal error: %v", err)
		}
	} else {
		return errors.New("you are not allowed to delete the group")
	}

	return nil
}

// FindGroup implements GroupService.
func (g *GroupServiceImpl) FindGroup(findGroupRequest request.FindGroupRequest) (response.GroupResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(findGroupRequest.Token)
	if err != nil {
		return response.GroupResponse{}, err
	}

	IsStudent := g.GroupRepository.IsStudent(context.Background(), userId, findGroupRequest.GroupId)
	IsTeacher := g.GroupRepository.IsTeacher(context.Background(), userId, findGroupRequest.GroupId)

	if IsStudent || IsTeacher {
		group, err := g.GroupRepository.FindById(context.Background(), findGroupRequest.GroupId)
		if err != nil {
			return response.GroupResponse{}, fmt.Errorf("internal error: %v", err)
		}

		gr := response.GroupResponse{
			UserId:   userId,
			GroupId:  findGroupRequest.GroupId,
			Title:    group.Title,
			Students: group.Students,
		}

		return gr, nil
	} else {
		return response.GroupResponse{}, errors.New("you are not in the group")
	}
}

// FindGroupsStudent implements GroupService.
func (g *GroupServiceImpl) FindGroupsStudent(findGroupsStudentRequest request.FindGroupsStudentRequest) ([]response.GroupResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(findGroupsStudentRequest.Token)
	if err != nil {
		return nil, err
	}

	var groupsResponse []response.GroupResponse

	groups, groups_err := g.GroupRepository.FindByStudentId(context.Background(), userId)
	if groups_err != nil {
		return groupsResponse, groups_err
	}

	for _, teacher_group := range groups {
		groupsResponse = append(groupsResponse, response.GroupResponse{
			UserId:   userId,
			GroupId:  teacher_group.Id.Hex(),
			Title:    teacher_group.Title,
			Students: teacher_group.Students,
		})
	}

	return groupsResponse, nil
}

// FindGroupsTeacher implements GroupService.
func (g *GroupServiceImpl) FindGroupsTeacher(findGroupsTeacherRequest request.FindGroupsTeacherRequest) ([]response.GroupResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(findGroupsTeacherRequest.Token)
	if err != nil {
		return nil, err
	}

	var groupsResponse []response.GroupResponse

	teacher_groups, teacher_groups_err := g.GroupRepository.FindByTeacherId(context.Background(), userId)
	if teacher_groups_err != nil {
		return groupsResponse, teacher_groups_err
	}

	for _, teacher_group := range teacher_groups {
		groupsResponse = append(groupsResponse, response.GroupResponse{
			UserId:   userId,
			GroupId:  teacher_group.Id.Hex(),
			Title:    teacher_group.Title,
			Students: teacher_group.Students,
		})
	}

	return groupsResponse, nil
}

// FindStudent implements GroupService.
func (g *GroupServiceImpl) FindStudent(findStudentRequest request.FindStudentRequest) (response.StudentResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(findStudentRequest.Token)
	if err != nil {
		return response.StudentResponse{}, err
	}

	IsStudent := g.GroupRepository.IsStudent(context.Background(), userId, findStudentRequest.GroupId)
	IsTeacher := g.GroupRepository.IsTeacher(context.Background(), userId, findStudentRequest.GroupId)

	if !IsTeacher && !IsStudent {
		return response.StudentResponse{}, errors.New("you are not allowed")
	}

	student, err := g.AuthenticationService.FindUser(findStudentRequest.StudentId)
	if err != nil {
		return response.StudentResponse{}, fmt.Errorf("internal error: %v", err)
	}

	return response.StudentResponse(student), nil
}

// FindTeacher implements GroupService.
func (g *GroupServiceImpl) FindTeacher(findTeacherRequest request.FindTeacherRequest) (response.TeacherResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(findTeacherRequest.Token)
	if err != nil {
		return response.TeacherResponse{}, err
	}

	group, err := g.GroupRepository.FindById(context.Background(), findTeacherRequest.GroupId)
	if err != nil {
		return response.TeacherResponse{}, err
	}

	IsStudent := g.GroupRepository.IsStudent(context.Background(), userId, findTeacherRequest.GroupId)

	if IsStudent {
		teacher, err := g.AuthenticationService.FindUser(group.TeacherId)
		if err != nil {
			return response.TeacherResponse{}, fmt.Errorf("internal error: %v", err)
		}

		return response.TeacherResponse{
			TeacherId: group.TeacherId,
			Email:     teacher.Email,
			Username:  teacher.Username,
		}, nil
	} else {
		return response.TeacherResponse{}, errors.New("you are not allowed")
	}
}

// GetStatistics implements GroupService.
func (g *GroupServiceImpl) GetStatistics(getStatisticsRequest request.GetStatisticsRequest) (response.StatisticsResponse, error) {
	userId, err := g.AuthenticationService.GetUserId(getStatisticsRequest.Token)
	if err != nil {
		return response.StatisticsResponse{}, err
	}

	IsStudent := g.GroupRepository.IsStudent(context.Background(), userId, getStatisticsRequest.GroupId)
	IsTeacher := g.GroupRepository.IsTeacher(context.Background(), userId, getStatisticsRequest.GroupId)

	if IsStudent || IsTeacher {
		statId, err := g.StatisticsRepository.GetId(context.Background(), getStatisticsRequest.GroupId, getStatisticsRequest.StudentId)
		if err != nil {
			return response.StatisticsResponse{}, err
		}
		res, err := g.StatisticsRepository.GetStatistics(context.Background(), statId)
		if err != nil {
			return response.StatisticsResponse{}, fmt.Errorf("internal error: %v", err)
		}
		return response.StatisticsResponse{
			StatId:    res.Id.Hex(),
			GroupId:   res.Group_id.Hex(),
			TeacherId: res.TeacherId,
			StudentId: res.StudentId,
			Words:     res.Words,
		}, nil
	} else {
		return response.StatisticsResponse{}, errors.New("you are not allowed")
	}
}

// RemoveStudent implements GroupService.
func (g *GroupServiceImpl) RemoveStudent(removeStudentRequest request.RemoveStudentRequest) error {
	userId, err := g.AuthenticationService.GetUserId(removeStudentRequest.Token)
	if err != nil {
		return err
	}

	if IsTeacher := g.GroupRepository.IsTeacher(context.Background(), userId, removeStudentRequest.GroupId); IsTeacher {
		err = g.GroupRepository.RemoveStudent(context.Background(), removeStudentRequest.UserId, removeStudentRequest.GroupId)
		if err != nil {
			statId, err := g.StatisticsRepository.GetId(context.Background(), removeStudentRequest.GroupId, removeStudentRequest.UserId)
			if err != nil {
				return err
			}
			err = g.StatisticsRepository.DeleteStatistics(context.Background(), statId)
			if err != nil {
				return err
			}

			return fmt.Errorf("internal error: %v", err)
		}
	} else {
		return errors.New("you are not allowed to remove the student")
	}

	return nil
}

func NewGroupServiceImpl(
	authenticationService AuthenticationService,
	vocabService VocabService,
	validate *validator.Validate,
	groupRepository repository.GroupRepository,
	statisticsRepository repository.StatisticsRepository,
) GroupService {
	return &GroupServiceImpl{
		AuthenticationService: authenticationService,
		VocabService:          vocabService,
		Validate:              validate,
		GroupRepository:       groupRepository,
		StatisticsRepository:  statisticsRepository,
	}
}
