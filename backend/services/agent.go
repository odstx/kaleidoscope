package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"kaleidoscope/config"
	"kaleidoscope/models"
)

var agentCfg *config.Config

type AgentService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAgentService(db *gorm.DB, cfg *config.Config) *AgentService {
	return &AgentService{db: db, cfg: cfg}
}

func (s *AgentService) Chat(ctx context.Context, userUID, userMessage string) (string, error) {
	var agent models.Agent
	if err := s.db.WithContext(ctx).Where("user_uid = ?", userUID).First(&agent).Error; err != nil {
		agent = models.Agent{
			UserUID:  userUID,
			Messages: "[]",
		}
		if err := s.db.WithContext(ctx).Create(&agent).Error; err != nil {
			return "", err
		}
	}

	var messages []map[string]string
	if err := json.Unmarshal([]byte(agent.Messages), &messages); err != nil {
		messages = []map[string]string{}
	}

	if s.cfg.LLM.SystemPrompt != "" {
		messages = append([]map[string]string{{
			"role":    "system",
			"content": s.cfg.LLM.SystemPrompt,
		}}, messages...)
	}

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": userMessage,
	})

	llmResponse, err := s.callLLM(ctx, messages)
	if err != nil {
		return "", err
	}

	messages = append(messages, map[string]string{
		"role":    "assistant",
		"content": llmResponse,
	})

	messagesJSON, _ := json.Marshal(messages)
	agent.Messages = string(messagesJSON)
	s.db.WithContext(ctx).Save(&agent)

	return llmResponse, nil
}

func (s *AgentService) callLLM(ctx context.Context, messages []map[string]string) (string, error) {
	cfg := s.cfg
	llmURL := cfg.LLM.URL
	apiKey := cfg.LLM.APIKey

	if llmURL == "" || apiKey == "" {
		return "LLM not configured. Please configure LLM URL and API key in the settings.", nil
	}

	model := "qwen-plus"
	if s.cfg.LLM.Model != "" {
		model = s.cfg.LLM.Model
	}
	reqBody := map[string]interface{}{
		"model":    model,
		"messages": messages,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", llmURL+"/chat/completions", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm api error: %v", result)
	}

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", fmt.Errorf("no response from llm")
	}

	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	return message["content"].(string), nil
}
