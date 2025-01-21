package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mono_pardo/internal/api/errors"
	domain "mono_pardo/internal/domain/sets"
	"mono_pardo/pkg/data/request"
)

type SetsController struct {
	setsService domain.Service
}

func NewSetsController(service domain.Service) *SetsController {
	return &SetsController{setsService: service}
}

func (controller *SetsController) CreateSet(ctx *gin.Context) {
	var req request.CreateSetRequest
	if !BindJSON(ctx, &req) {
		return
	}

	req.UserId = ctx.GetInt("userId")

	if err := controller.setsService.CreateSet(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller *SetsController) GetSets(ctx *gin.Context) {
	req := request.GetSetsRequest{UserId: ctx.GetInt("userId")}

	res, err := controller.setsService.GetSets(req)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, errors.InternalError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *SetsController) GetSet(ctx *gin.Context) {
	req := request.GetSetRequest{
		UserId:    ctx.GetInt("userId"),
		WordSetId: ctx.Param("setId"),
	}

	res, err := controller.setsService.GetSet(req)
	if err != nil {
		SendError(ctx, http.StatusNotFound, errors.NotFoundError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *SetsController) UpdateSet(ctx *gin.Context) {
	var req request.UpdateSetRequest
	if !BindJSON(ctx, &req) {
		return
	}

	req.UserId = ctx.GetInt("userId")
	req.WordSetId = ctx.Param("setId")

	if len(req.Updates) == 0 {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "No updates provided")
		return
	}

	if err := controller.setsService.UpdateSet(req); err != nil {
		if err.Error() == domain.ErrAccessForbidden {
			SendError(ctx, http.StatusForbidden, errors.ForbiddenError, err.Error())
			return
		}
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *SetsController) DeleteSet(ctx *gin.Context) {
	req := request.DeleteSetRequest{
		UserId:    ctx.GetInt("userId"),
		WordSetId: ctx.Param("setId"),
	}

	if err := controller.setsService.DeleteSet(req); err != nil {
		if err.Error() == domain.ErrAccessForbidden {
			SendError(ctx, http.StatusForbidden, errors.ForbiddenError, err.Error())
			return
		}
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *SetsController) AddWord(ctx *gin.Context) {
	var req request.AddWordRequest
	if !BindJSON(ctx, &req) {
		return
	}

	req.UserId = ctx.GetInt("userId")
	req.WordSetId = ctx.Param("setId")

	if len(req.Words) == 0 {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "No words provided")
		return
	}

	if err := controller.setsService.AddWord(req); err != nil {
		if err.Error() == domain.ErrAccessForbidden {
			SendError(ctx, http.StatusForbidden, errors.ForbiddenError, err.Error())
			return
		}
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *SetsController) RemoveWord(ctx *gin.Context) {
	var req request.RemoveWordRequest
	if !BindJSON(ctx, &req) {
		return
	}

	req.UserId = ctx.GetInt("userId")
	req.WordSetId = ctx.Param("setId")

	if len(req.Words) == 0 {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "No words provided")
		return
	}

	if err := controller.setsService.RemoveWord(req); err != nil {
		if err.Error() == domain.ErrAccessForbidden {
			SendError(ctx, http.StatusForbidden, errors.ForbiddenError, err.Error())
			return
		}
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
