package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	_ "kaleidoscope/docs"
	"kaleidoscope/middleware"
	"kaleidoscope/services"
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(router *gin.Engine, logger *zap.Logger, userService *services.UserService) {
	userController := NewUserController(logger, userService)
	systemController := NewSystemController(logger)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", userController.Register)
			userGroup.POST("/login", userController.Login)
			userGroup.GET("/info", middleware.JWTAuth(), userController.GetUserInfo)
		}

		systemGroup := v1.Group("/system")
		{
			systemGroup.GET("/info", systemController.GetSystemInfo)
		}
	}
}
