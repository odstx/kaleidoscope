package controllers

import (
	"kaleidoscope/config"
	_ "kaleidoscope/docs"
	"kaleidoscope/middleware"
	"kaleidoscope/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, logger *zap.Logger, userService *services.UserService, oidcService *services.OIDCService, rateLimiter *middleware.RateLimiter, cfg *config.Config, db *gorm.DB) {
	userController := NewUserController(logger, userService, oidcService)
	systemController := NewSystemController(logger, cfg)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			if rateLimiter != nil {
				userGroup.POST("/register", rateLimiter.RateLimit(), userController.Register)
				userGroup.POST("/login", rateLimiter.RateLimit(), userController.Login)
				userGroup.POST("/forgot-password", rateLimiter.RateLimit(), userController.ForgotPassword)
				userGroup.POST("/reset-password", rateLimiter.RateLimit(), userController.ResetPassword)
				userGroup.GET("/info", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.GetUserInfo)
				userGroup.POST("/totp/setup", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.SetupTOTP)
				userGroup.POST("/totp/verify", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.VerifyTOTP)
				userGroup.POST("/totp/enable", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.EnableTOTP)
				userGroup.POST("/totp/disable", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.DisableTOTP)
				userGroup.POST("/hawk/setup", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.SetupHawk)
				userGroup.POST("/hawk/enable", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.EnableHawk)
				userGroup.POST("/hawk/disable", middleware.CombinedAuth(cfg, db), rateLimiter.RateLimit(), userController.DisableHawk)

				userGroup.GET("/oidc/login", rateLimiter.RateLimit(), userController.OIDCLogin)
				userGroup.POST("/oidc/callback", rateLimiter.RateLimit(), userController.OIDCCallback)
			} else {
				userGroup.POST("/register", userController.Register)
				userGroup.POST("/login", userController.Login)
				userGroup.POST("/forgot-password", userController.ForgotPassword)
				userGroup.POST("/reset-password", userController.ResetPassword)
				userGroup.GET("/info", middleware.CombinedAuth(cfg, db), userController.GetUserInfo)
				userGroup.POST("/totp/setup", middleware.CombinedAuth(cfg, db), userController.SetupTOTP)
				userGroup.POST("/totp/verify", middleware.CombinedAuth(cfg, db), userController.VerifyTOTP)
				userGroup.POST("/totp/enable", middleware.CombinedAuth(cfg, db), userController.EnableTOTP)
				userGroup.POST("/totp/disable", middleware.CombinedAuth(cfg, db), userController.DisableTOTP)
				userGroup.POST("/hawk/setup", middleware.CombinedAuth(cfg, db), userController.SetupHawk)
				userGroup.POST("/hawk/enable", middleware.CombinedAuth(cfg, db), userController.EnableHawk)
				userGroup.POST("/hawk/disable", middleware.CombinedAuth(cfg, db), userController.DisableHawk)

				userGroup.GET("/oidc/login", userController.OIDCLogin)
				userGroup.POST("/oidc/callback", userController.OIDCCallback)
			}
		}

		systemGroup := v1.Group("/system")
		{
			if rateLimiter != nil {
				systemGroup.GET("/info", rateLimiter.RateLimit(), systemController.GetSystemInfo)
				systemGroup.GET("/config", rateLimiter.RateLimit(), systemController.GetConfig)
			} else {
				systemGroup.GET("/info", systemController.GetSystemInfo)
				systemGroup.GET("/config", systemController.GetConfig)
			}
		}
	}
}
