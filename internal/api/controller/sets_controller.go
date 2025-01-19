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
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *SetsController) GetSet(ctx *gin.Context) {}

func (controller *SetsController) UpdateSet(ctx *gin.Context) {}

func (controller *SetsController) DeleteSet(ctx *gin.Context) {}

func (controller *SetsController) AddWord(ctx *gin.Context) {}

func (controller *SetsController) RemoveWord(ctx *gin.Context) {}
