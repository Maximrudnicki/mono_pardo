package controller

import (
	"net/http"
	"strconv"

	"mono_pardo/internal/api/errors"
	wordsDomain "mono_pardo/internal/domain/words"
	"mono_pardo/pkg/data/request"

	"github.com/gin-gonic/gin"
)

type VocabController struct {
	vocabService wordsDomain.Service
}

func NewVocabController(service wordsDomain.Service) *VocabController {
	return &VocabController{vocabService: service}
}

func (controller *VocabController) CreateWord(ctx *gin.Context) {
	req := request.CreateWordRequest{UserId: ctx.GetInt("userId")}
	if !BindJSON(ctx, &req) {
		return
	}

	if err := controller.vocabService.CreateWord(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller *VocabController) DeleteWord(ctx *gin.Context) {
	wordId := ctx.Param("wordId")
	id, err := strconv.Atoi(wordId)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Cannot parse id from url")
		return
	}

	req := request.DeleteWordRequest{UserId: ctx.GetInt("userId"), WordId: id}

	if err = controller.vocabService.DeleteWord(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *VocabController) GetWords(ctx *gin.Context) {
	vocabRequest := request.VocabRequest{UserId: ctx.GetInt("userId")}

	res, err := controller.vocabService.GetWords(vocabRequest)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *VocabController) UpdateWord(ctx *gin.Context) {
	var words []request.WordUpdate
	if !BindJSON(ctx, &words) {
		return
	}

	req := request.UpdateWordRequest{UserId: ctx.GetInt("userId"), Words: words}

	if len(req.Words) == 0 {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "No updates provided")
		return
	}

	if err := controller.vocabService.UpdateWord(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
