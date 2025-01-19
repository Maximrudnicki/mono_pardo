package api

import (
	"mono_pardo/internal/api/controller"
	"mono_pardo/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authenticationController *controller.AuthenticationController,
	vocabController *controller.VocabController,
	setsController *controller.SetsController) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(authenticationController.AuthenticationService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r := router.Group("/api/v1")
	authenticationRouter := r.Group("/authentication")
	authenticationRouter.POST("/login", authenticationController.Login)
	authenticationRouter.POST("/register", authenticationController.Register)

	vocabRouter := r.Group("/vocab", authMiddleware.Handle())
	vocabRouter.GET("", vocabController.GetWords)
	vocabRouter.POST("", vocabController.CreateWord)
	vocabRouter.PATCH("", vocabController.UpdateWords)
	vocabRouter.DELETE("/:wordId", vocabController.DeleteWord)

	setsRouter := r.Group("/sets", authMiddleware.Handle())
	setsRouter.POST("", setsController.CreateSet)
	setsRouter.GET("", setsController.GetSets)
	setsRouter.GET("/:setId", setsController.GetSet)
	setsRouter.PATCH("/:setId", setsController.UpdateSet)
	setsRouter.DELETE("/:setId", setsController.DeleteSet)

	setsRouter.POST("/:setId/words", setsController.AddWord)
	setsRouter.DELETE("/:setId/words", setsController.RemoveWord)

	return router
}
