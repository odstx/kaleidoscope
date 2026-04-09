package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"kaleidoscope/services"
)

type AppController struct {
	logger     *zap.Logger
	appService *services.AppService
}

func NewAppController(logger *zap.Logger, appService *services.AppService) *AppController {
	return &AppController{
		logger:     logger,
		appService: appService,
	}
}

// GetApps godoc
// @Summary      Get all apps
// @Description  Get all enabled apps ordered by order field
// @Tags         apps
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {array}   models.App  "Apps retrieved successfully"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /apps [get]
func (ac *AppController) GetApps(c *gin.Context) {
	apps, err := ac.appService.GetAllApps()
	if err != nil {
		ac.logger.Error("Failed to get apps", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get apps"})
		return
	}

	c.JSON(http.StatusOK, apps)
}
