package controller

import (
	"github.com/gin-gonic/gin"

	domain "mono_pardo/internal/domain/sets"
)

type SetsController struct {
	setsService domain.Service
}

func NewSetsController(service domain.Service) *SetsController {
	return &SetsController{setsService: service}
}

func (controller *SetsController) CreateSet(ctx *gin.Context) {}

func (controller *SetsController) GetSets(ctx *gin.Context) {}

func (controller *SetsController) UpdateSet(ctx *gin.Context) {}

func (controller *SetsController) DeleteSet(ctx *gin.Context) {}

func (controller *SetsController) GetSet(ctx *gin.Context) {}

func (controller *SetsController) AddWord(ctx *gin.Context) {}

func (controller *SetsController) RemoveWord(ctx *gin.Context) {}
