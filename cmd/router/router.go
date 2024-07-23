package router

import (
	"mono_pardo/cmd/controller"
	"mono_pardo/cmd/middleware"
	
	"github.com/gin-gonic/gin"
)

func NewRouter(
	authenticationController *controller.AuthenticationController,
	vocabController *controller.VocabController,
	groupController *controller.GroupController) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r := router.Group("/api/v1")
	authenticationRouter := r.Group("/authentication")
	authenticationRouter.POST("/login/", authenticationController.Login)
	authenticationRouter.POST("/register", authenticationController.Register)
	
	groupRouter := r.Group("/group")
	groupRouter.POST("/add", groupController.AddStudent)
	groupRouter.POST("/add_word", groupController.AddWordToUser)
	groupRouter.POST("/", groupController.CreateGroup)
	groupRouter.DELETE("/:groupId", groupController.DeleteGroup)
	groupRouter.POST("/find", groupController.FindGroup)
	groupRouter.GET("/find_teacher", groupController.FindGroupsTeacher)
	groupRouter.GET("/find_student", groupController.FindGroupsStudent)
	groupRouter.POST("/find_teacher_info", groupController.FindTeacher)
	groupRouter.POST("/find_student_info", groupController.FindStudent)
	groupRouter.POST("/get_statistics", groupController.GetStatistics)
	groupRouter.PATCH("/remove", groupController.RemoveStudent)

	vocabRouter := r.Group("/vocab")
	vocabRouter.GET("/", vocabController.GetWords)
	vocabRouter.POST("/", vocabController.CreateWord)
	vocabRouter.DELETE("/:wordId", vocabController.DeleteWord)
	vocabRouter.PATCH("/:wordId", vocabController.UpdateWord)
	vocabRouter.PATCH("/:wordId/status", vocabController.UpdateWordStatus)
    vocabRouter.PATCH("/:wordId/trainings", vocabController.ManageTrainings)

	return router
}
