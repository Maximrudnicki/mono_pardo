package controller

import (
	"log"
	"net/http"
	"strconv"

	wordsDomain "mono_pardo/internal/domain/words"
	"mono_pardo/internal/utils"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/gin-gonic/gin"
)

type VocabController struct {
	vocabService wordsDomain.Service
}

func NewVocabController(service wordsDomain.Service) *VocabController {
	return &VocabController{vocabService: service}
}

func (controller *VocabController) CreateWord(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	req := request.CreateWordRequest{Token: token}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add word",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err_cw := controller.vocabService.CreateWord(req)
	if err_cw != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot add word",
		}
		log.Printf("Cannot add: %v", err_cw)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller *VocabController) DeleteWord(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	wordId := ctx.Param("wordId")
	id, err_id := strconv.Atoi(wordId)

	req := request.DeleteWordRequest{
		Token:  token,
		WordId: id,
	}

	err_dw := controller.vocabService.DeleteWord(req)
	if err_dw != nil || err_id != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot delete word",
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *VocabController) GetWords(ctx *gin.Context) {
	token, err := utils.GetToken(ctx)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "Login required",
		}
		ctx.JSON(http.StatusUnauthorized, webResponse)
		return
	}

	vocabRequest := request.VocabRequest{
		TokenType: "Bearer",
		Token:     token,
	}

	res, err := controller.vocabService.GetWords(vocabRequest)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot get words",
		}
		log.Printf("err: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (controller *VocabController) UpdateWord(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	var words []request.WordUpdate

	err := ctx.ShouldBindJSON(&words)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot update word",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	request := request.UpdateWordRequest{Token: token, Words: words}

	if len(request.Words) == 0 {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "No updates provided",
		}
		log.Printf("No updates provided")
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err = controller.vocabService.UpdateWord(request)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot update words",
		}
		log.Printf("Cannot update: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *VocabController) UpdateWordStatus(ctx *gin.Context) {
	token, _ := utils.GetToken(ctx)

	wordId := ctx.Param("wordId")
	id, err_id := strconv.Atoi(wordId)

	req := request.UpdateWordStatusRequest{
		Token:  token,
		WordId: id,
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil || err_id != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot update status",
		}
		log.Printf("Cannot bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err_uws := controller.vocabService.UpdateWordStatus(req)

	if err_uws != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Cannot update status",
		}
		log.Printf("Cannot update status: %v", err_uws)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ctx.Status(http.StatusOK)
}
