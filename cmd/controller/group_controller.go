package controller

import (
	"log"
	"net/http"

	"mono_pardo/cmd/data/request"
	"mono_pardo/cmd/data/response"
	"mono_pardo/cmd/service"
	"mono_pardo/cmd/utils"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupService service.GroupService
	vocabService service.VocabService
}

func NewGroupController(service service.GroupService, vs service.VocabService) *GroupController {
	return &GroupController{groupService: service, vocabService: vs}
}

func (controller *GroupController) AddStudent(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.AddStudentRequest{Token: token}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add student",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err_as := controller.groupService.AddStudent(req)
	if err_as != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add student",
		}
		log.Printf("Cannot add student: %v", err_as)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully added!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) AddWordToUser(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.AddWordToUserRequest{Token: token}
	ctx.ShouldBindJSON(&req)

	res, err_fg := controller.groupService.AddWordToUser(req)
	if err_fg != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add word",
		}
		log.Printf("Cannot add word: %v", err_fg)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully added word!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) CreateGroup(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.CreateGroupRequest{Token: token}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot create group",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err_cg := controller.groupService.CreateGroup(req)
	if err_cg != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot create group",
		}
		log.Printf("Cannot create: %v", err_cg)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully created!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) DeleteGroup(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.DeleteGroupRequest{
		Token:   token,
		GroupId: ctx.Param("groupId"),
	}

	err_dg := controller.groupService.DeleteGroup(req)
	if err_dg != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add word",
		}
		log.Printf("Cannot delete group: %v", err_dg)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully deleted!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) FindGroup(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.FindGroupRequest{Token: token}
	ctx.ShouldBindJSON(&req)

	groupResponse, err_fg := controller.groupService.FindGroup(req)
	if err_fg != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot find group",
		}
		log.Printf("Cannot finds group: %v", err_fg)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	students := make([]response.StudentInformation, 0, len(groupResponse.Students))
	for _, studentId := range groupResponse.Students {
		findStudentRequest := request.FindStudentRequest{
			Token:     token,
			StudentId: studentId,
			GroupId:   req.GroupId,
		}
		getStatisticsRequest := request.GetStatisticsRequest{
			Token:     token,
			StudentId: studentId,
			GroupId:   req.GroupId,
		}

		statResp, err_gs := controller.groupService.GetStatistics(getStatisticsRequest)
		if err_gs != nil {
			webResponse := response.Response{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Cannot get stats",
			}
			log.Printf("Cannot get stats: %v", err_gs)
			ctx.JSON(http.StatusBadRequest, webResponse)
			return
		}

		words := make([]response.VocabResponse, 0, len(statResp.Words))
		for _, wordId := range statResp.Words {
			findWordRequest := request.FindWordRequest{
				WordId: wordId,
			}
			word, err := controller.vocabService.FindWord(findWordRequest)
			if err != nil {
				webResponse := response.Response{
					Code:    http.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Cannot find word",
				}
				log.Printf("Cannot find word with id %d: %v", wordId, err)
				ctx.JSON(http.StatusInternalServerError, webResponse)
				return
			}
			if word.ID != 0 {
				words = append(words, word)
			}
		}

		student, err := controller.groupService.FindStudent(findStudentRequest)
		studentInfo := response.StudentInformation{
			StudentId: studentId,
			Email:     student.Email,
			Username:  student.Username,
			Words:     words,
		}
		if err != nil {
			webResponse := response.Response{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Cannot find student",
			}
			log.Printf("cannot find studnent with id %v: %v", student, err)
			ctx.JSON(http.StatusInternalServerError, webResponse)
			return
		}
		students = append(students, studentInfo)
	}

	res := struct {
		UserId   int                           `json:"user_id"`
		GroupId  string                        `json:"group_id"`
		Title    string                        `json:"title"`
		Students []response.StudentInformation `json:"students"`
	}{
		UserId:   groupResponse.UserId,
		GroupId:  groupResponse.GroupId,
		Title:    groupResponse.Title,
		Students: students,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found groups!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) FindStudent(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	findStudentRequest := request.FindStudentRequest{Token: token}
	ctx.ShouldBindJSON(&findStudentRequest)

	res, err_fs := controller.groupService.FindStudent(findStudentRequest)
	if err_fs != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot find student",
		}
		log.Printf("Cannot find student: %v", err_fs)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found student!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) FindTeacher(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.FindTeacherRequest{Token: token}
	ctx.ShouldBindJSON(&req)

	res, err_ft := controller.groupService.FindTeacher(req)
	if err_ft != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot find teacher",
		}
		log.Printf("Cannot find teacher: %v", err_ft)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found teacher!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) FindGroupsTeacher(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.FindGroupsTeacherRequest{Token: token}

	res, err_fgt := controller.groupService.FindGroupsTeacher(req)
	if err_fgt != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot find groups",
		}
		log.Printf("Cannot finds groups: %v", err_fgt)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found groups!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) FindGroupsStudent(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.FindGroupsStudentRequest{Token: token}

	res, err_fgs := controller.groupService.FindGroupsStudent(req)
	if err_fgs != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot find groups",
		}
		log.Printf("Cannot finds groups: %v", err_fgs)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found groups!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) GetStatistics(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.GetStatisticsRequest{Token: token}
	ctx.ShouldBindJSON(&req)

	statResp, err_gs := controller.groupService.GetStatistics(req)
	if err_gs != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot get stats",
		}
		log.Printf("Cannot get stats: %v", err_gs)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	words := make([]response.VocabResponse, 0, len(statResp.Words))
	for _, wordId := range statResp.Words {
		findWordRequest := request.FindWordRequest{
			WordId: wordId,
		}
		word, err := controller.vocabService.FindWord(findWordRequest)
		if err != nil {
			webResponse := response.Response{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Cannot find word",
			}
			log.Printf("Cannot find word with id %d: %v", wordId, err)
			ctx.JSON(http.StatusInternalServerError, webResponse)
			return
		}
		if word.ID != 0 {
			words = append(words, word)
		}
	}

	student, err := controller.groupService.FindStudent(request.FindStudentRequest{
		Token: token, StudentId: statResp.StudentId, GroupId: statResp.GroupId,
	})
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Cannot find student",
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	res := struct {
		StatId    string                   `json:"statistics_id"`
		GroupId   string                   `json:"group_id"`
		TeacherId int                      `json:"teacher_id"`
		Student   response.StudentInfo     `json:"student"`
		Words     []response.VocabResponse `json:"words"`
	}{
		StatId:    statResp.StatId,
		GroupId:   statResp.GroupId,
		TeacherId: statResp.TeacherId,
		Student: response.StudentInfo{
			StudentId: statResp.StudentId,
			Email:     student.Email,
			Username:  student.Username,
		},
		Words: words,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully found student!",
		Data:    res,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *GroupController) RemoveStudent(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.RemoveStudentRequest{Token: token}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add student",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err_rs := controller.groupService.RemoveStudent(req)
	if err_rs != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot remove student",
		}
		log.Printf("Cannot remove student: %v", err_rs)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully removed!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
