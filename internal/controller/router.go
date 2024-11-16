package controller

import (
	"mono_pardo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authenticationController *AuthenticationController,
	vocabController *VocabController) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r := router.Group("/api/v1")
	authenticationRouter := r.Group("/authentication")
	authenticationRouter.POST("/login/", authenticationController.Login)
	authenticationRouter.POST("/register", authenticationController.Register)

	vocabRouter := r.Group("/vocab")
	vocabRouter.GET("/", vocabController.GetWords)
	vocabRouter.POST("/", vocabController.CreateWord)
	vocabRouter.DELETE("/:wordId", vocabController.DeleteWord)
	vocabRouter.PATCH("/:wordId", vocabController.UpdateWord)
	vocabRouter.PATCH("/:wordId/status", vocabController.UpdateWordStatus)
	vocabRouter.PATCH("/:wordId/trainings", vocabController.ManageTrainings)

	return router
}
