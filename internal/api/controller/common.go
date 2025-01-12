package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mono_pardo/internal/api/errors"
)

func SendError(c *gin.Context, status int, errType errors.ErrorType, message string) {
	c.JSON(status, errors.NewAPIError(errType, message))
}

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		SendError(c, http.StatusBadRequest, errors.ValidationError, "Invalid request format")
		return false
	}
	return true
}
