package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mono_pardo/internal/api/errors"
)

type BaseController struct{}

func (b *BaseController) SendError(c *gin.Context, status int, errType errors.ErrorType, message string) {
	c.JSON(status, errors.NewAPIError(errType, message))
}

func (b *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		b.SendError(c, http.StatusBadRequest, errors.ValidationError, "Invalid request format")
		return false
	}
	return true
}

func NewBaseController() *BaseController {
	return &BaseController{}
}
