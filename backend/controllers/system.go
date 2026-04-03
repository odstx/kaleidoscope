package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kaleidoscope/version"
)

type SystemController struct {
	logger *zap.Logger
}

func NewSystemController(logger *zap.Logger) *SystemController {
	return &SystemController{
		logger: logger,
	}
}

// GetSystemInfo godoc
// @Summary Get system information
// @Description Get backend version and build information
// @Tags system
// @Produce json
// @Success 200 {object} version.Info
// @Router /system/info [get]
func (sc *SystemController) GetSystemInfo(c *gin.Context) {
	info := version.GetInfo()
	c.JSON(http.StatusOK, info)
}
