package controller

import (
	"net/http"
	"strconv"

	"mono_pardo/internal/api/errors"
	wordsDomain "mono_pardo/internal/domain/words"
	"mono_pardo/internal/utils"
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
	token, err := utils.GetToken(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, errors.UnauthorizedError, "Login required")
		return
	}

	req := request.CreateWordRequest{Token: token}
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
	token, err := utils.GetToken(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, errors.UnauthorizedError, "Login required")
		return
	}

	wordId := ctx.Param("wordId")
	id, err := strconv.Atoi(wordId)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Cannot parse id from url")
		return
	}

	req := request.DeleteWordRequest{Token: token, WordId: id}

	if err = controller.vocabService.DeleteWord(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *VocabController) GetWords(ctx *gin.Context) {
	token, err := utils.GetToken(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, errors.UnauthorizedError, "Login required")
		return
	}

	vocabRequest := request.VocabRequest{TokenType: "Bearer", Token: token}

	res, err := controller.vocabService.GetWords(vocabRequest)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *VocabController) UpdateWord(ctx *gin.Context) {
	token, err := utils.GetToken(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, errors.UnauthorizedError, "Login required")
		return
	}

	var words []request.WordUpdate
	if !BindJSON(ctx, &words) {
		return
	}

	request := request.UpdateWordRequest{Token: token, Words: words}

	if len(request.Words) == 0 {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "No updates provided")
		return
	}

	if err := controller.vocabService.UpdateWord(request); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
