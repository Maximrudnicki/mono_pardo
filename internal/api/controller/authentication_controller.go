package controller

import (
	"net/http"

	"mono_pardo/internal/api/errors"
	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	AuthenticationService usersDomain.Service
}

func NewAuthenticationController(service usersDomain.Service) *AuthenticationController {
	return &AuthenticationController{AuthenticationService: service}
}

func (controller *AuthenticationController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	if !BindJSON(ctx, &loginRequest) {
		return
	}

	token, err := controller.AuthenticationService.Login(loginRequest)
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
	createUserRequest := request.CreateUserRequest{}
	if !BindJSON(ctx, &createUserRequest) {
		return
	}

	if err := controller.AuthenticationService.Register(createUserRequest); err != nil {
		SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Please use another email address")
		return
	}

	ctx.Status(http.StatusCreated)
}
