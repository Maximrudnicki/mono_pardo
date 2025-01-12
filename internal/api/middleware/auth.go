package middleware

import (
	"net/http"

	"mono_pardo/internal/api/errors"
	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService usersDomain.Service
}

func NewAuthMiddleware(authService usersDomain.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetToken(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, errors.NewAPIError(errors.UnauthorizedError, "Login required"))
			return
		}

		userId, err := m.authService.GetUserId(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, errors.NewAPIError(errors.UnauthorizedError, "Invalid token"))
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
