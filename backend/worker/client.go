package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Client handles task enqueuing
type Client struct {
	client *asynq.Client
}

// NewClient creates a new task client
func NewClient(redisAddr, redisPassword string, redisDB int) *Client {
	redisConnOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	client := asynq.NewClient(redisConnOpt)

	return &Client{
		client: client,
	}
}

// EnqueueSendWelcomeEmail enqueues a welcome email task
func (c *Client) EnqueueSendWelcomeEmail(ctx context.Context, userID uint, username, email string) error {
	payload := SendWelcomeEmailPayload{
		UserID:   userID,
		Username: username,
		Email:    email,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(string(TaskSendWelcomeEmail), payloadBytes, asynq.Queue("default"))

	_, err = c.client.EnqueueContext(ctx, task)
	return err
}

// EnqueueSendPasswordResetEmail enqueues a password reset email task
func (c *Client) EnqueueSendPasswordResetEmail(ctx context.Context, userID uint, username, email, token string) error {
	payload := SendPasswordResetEmailPayload{
		UserID:   userID,
		Username: username,
		Email:    email,
		Token:    token,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(string(TaskSendPasswordResetEmail), payloadBytes, asynq.Queue("default"))

	_, err = c.client.EnqueueContext(ctx, task)
	return err
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.client.Close()
}
