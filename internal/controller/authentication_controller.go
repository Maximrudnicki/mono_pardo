package controller

import (
	"log"
	"net/http"

	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authenticationService usersDomain.AuthenticationService
}

func NewAuthenticationController(service usersDomain.AuthenticationService) *AuthenticationController {
	return &AuthenticationController{authenticationService: service}
}

func (controller *AuthenticationController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		panic(err)
	}

	token, err_token := controller.authenticationService.Login(loginRequest)
	if err_token != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid username or password",
		}
		log.Printf("Token err: %v", err_token)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	resp := response.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully log in!",
		Data:    resp,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *AuthenticationController) Register(ctx *gin.Context) {
	createUserRequest := request.CreateUserRequest{}
	err := ctx.ShouldBindJSON(&createUserRequest)
	if err != nil {
		panic(err)
	}

	reg_err := controller.authenticationService.Register(createUserRequest)
	if reg_err != nil {
		webResponse := response.Response{
			Code:    http.StatusForbidden,
			Status:  "Forbidden",
			Message: "Please use another email address",
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    201,
		Status:  "Created",
		Message: "Successfully created user!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
