package controller

import (
	"net/http"

	"mono_pardo/internal/api/errors"
	domain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	AuthenticationService domain.Service
}

func NewAuthenticationController(service domain.Service) *AuthenticationController {
	return &AuthenticationController{AuthenticationService: service}
}

func (controller *AuthenticationController) Login(ctx *gin.Context) {
	req := request.LoginRequest{}
	if !BindJSON(ctx, &req) {
		return
	}

	token, err := controller.AuthenticationService.Login(req)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Invalid username or password")
		return
	}

	resp := response.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (controller *AuthenticationController) Register(ctx *gin.Context) {
	req := request.CreateUserRequest{}
	if !BindJSON(ctx, &req) {
		return
	}

	if err := controller.AuthenticationService.Register(req); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Please use another email address")
		return
	}

	ctx.Status(http.StatusCreated)
}
