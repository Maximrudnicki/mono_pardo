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
	BaseController
	authenticationService usersDomain.Service
}

func NewAuthenticationController(service usersDomain.Service) *AuthenticationController {
	return &AuthenticationController{BaseController: BaseController{}, authenticationService: service}
}

func (controller *AuthenticationController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	if !controller.BindJSON(ctx, &loginRequest) {
		return
	}

	token, err := controller.authenticationService.Login(loginRequest)
	if err != nil {
		controller.SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Invalid username or password")
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
	if !controller.BindJSON(ctx, &createUserRequest) {
		return
	}

	if err := controller.authenticationService.Register(createUserRequest); err != nil {
		controller.SendError(ctx, http.StatusBadRequest, errors.ValidationError, "Please use another email address")
		return
	}

	ctx.Status(http.StatusCreated)
}
