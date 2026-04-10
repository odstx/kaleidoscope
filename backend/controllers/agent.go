package controllers

import (
	"net/http"

	"kaleidoscope/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgentController struct {
	logger       *zap.Logger
	agentService *services.AgentService
}

func NewAgentController(logger *zap.Logger, agentService *services.AgentService) *AgentController {
	return &AgentController{logger: logger, agentService: agentService}
}

func (ac *AgentController) Chat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUID := c.GetString("user_uid")
	resp, err := ac.agentService.Chat(c.Request.Context(), userUID, req.Message)
	if err != nil {
		ac.logger.Error("agent chat error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp})
}
