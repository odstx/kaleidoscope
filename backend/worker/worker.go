package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// Worker handles background task processing
type Worker struct {
	server *asynq.Server
	logger *zap.Logger
}

// NewWorker creates a new worker instance
func NewWorker(redisAddr, redisPassword string, redisDB int, logger *zap.Logger) *Worker {
	redisConnOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	server := asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Worker{
		server: server,
		logger: logger,
	}
}

// Start starts the worker to process tasks
func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(string(TaskSendWelcomeEmail), w.handleSendWelcomeEmail)

	w.logger.Info("Starting Asynq worker...")
	if err := w.server.Run(mux); err != nil {
		return fmt.Errorf("failed to start Asynq worker: %w", err)
	}
	return nil
}

// Stop gracefully stops the worker
func (w *Worker) Stop() error {
	w.logger.Info("Stopping Asynq worker...")
	w.server.Shutdown()
	return nil
}

// handleSendWelcomeEmail processes the welcome email task
func (w *Worker) handleSendWelcomeEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendWelcomeEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// TODO: Implement actual email sending logic here
	// For now, just log the action
	w.logger.Info("Sending welcome email",
		zap.Uint("user_id", payload.UserID),
		zap.String("username", payload.Username),
		zap.String("email", payload.Email))

	// In a real implementation, you would:
	// 1. Use an email service/client to send the email
	// 2. Handle errors and retries appropriately
	// 3. Log success/failure

	return nil
}
