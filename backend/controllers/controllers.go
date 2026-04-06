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

func RegisterRoutes(router *gin.Engine, logger *zap.Logger, userService *services.UserService, rateLimiter *middleware.RateLimiter) {
	userController := NewUserController(logger, userService)
	systemController := NewSystemController(logger)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			if rateLimiter != nil {
				userGroup.POST("/register", rateLimiter.RateLimit(), userController.Register)
				userGroup.POST("/login", rateLimiter.RateLimit(), userController.Login)
				userGroup.GET("/info", middleware.JWTAuth(), rateLimiter.RateLimit(), userController.GetUserInfo)
				userGroup.POST("/totp/setup", middleware.JWTAuth(), rateLimiter.RateLimit(), userController.SetupTOTP)
				userGroup.POST("/totp/verify", middleware.JWTAuth(), rateLimiter.RateLimit(), userController.VerifyTOTP)
				userGroup.POST("/totp/enable", middleware.JWTAuth(), rateLimiter.RateLimit(), userController.EnableTOTP)
				userGroup.POST("/totp/disable", middleware.JWTAuth(), rateLimiter.RateLimit(), userController.DisableTOTP)
			} else {
				userGroup.POST("/register", userController.Register)
				userGroup.POST("/login", userController.Login)
				userGroup.GET("/info", middleware.JWTAuth(), userController.GetUserInfo)
				userGroup.POST("/totp/setup", middleware.JWTAuth(), userController.SetupTOTP)
				userGroup.POST("/totp/verify", middleware.JWTAuth(), userController.VerifyTOTP)
				userGroup.POST("/totp/enable", middleware.JWTAuth(), userController.EnableTOTP)
				userGroup.POST("/totp/disable", middleware.JWTAuth(), userController.DisableTOTP)
			}
		}

		systemGroup := v1.Group("/system")
		{
			if rateLimiter != nil {
				systemGroup.GET("/info", rateLimiter.RateLimit(), systemController.GetSystemInfo)
			} else {
				systemGroup.GET("/info", systemController.GetSystemInfo)
			}
		}
	}
}
