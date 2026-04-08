package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kaleidoscope/config"
	"kaleidoscope/version"
)

type SystemController struct {
	logger *zap.Logger
	cfg    *config.Config
}

func NewSystemController(logger *zap.Logger, cfg *config.Config) *SystemController {
	return &SystemController{
		logger: logger,
		cfg:    cfg,
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

type FrontendConfig struct {
	OIDCID          string `json:"oidcClientId"`
	OIDCEnabled     bool   `json:"oidcEnabled"`
	OIDCIssuerURL   string `json:"oidcIssuerUrl"`
	OIDCRedirectURI string `json:"oidcRedirectUri"`
}

// GetConfig godoc
// @Summary Get frontend configuration
// @Description Get frontend configuration including OIDC settings
// @Tags system
// @Produce json
// @Success 200 {object} FrontendConfig
// @Router /system/config [get]
func (sc *SystemController) GetConfig(c *gin.Context) {
	cfg := FrontendConfig{
		OIDCID:          sc.cfg.OIDC.ClientID,
		OIDCEnabled:     sc.cfg.OIDC.Enabled,
		OIDCIssuerURL:   sc.cfg.OIDC.IssuerURL,
		OIDCRedirectURI: sc.cfg.OIDC.RedirectURI,
	}
	c.JSON(http.StatusOK, cfg)
}
