package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"kaleidoscope/config"
	"kaleidoscope/utils"
)

// Worker handles background task processing
type Worker struct {
	server       *asynq.Server
	logger       *zap.Logger
	emailService *utils.EmailService
}

// NewWorker creates a new worker instance
func NewWorker(redisAddr, redisPassword string, redisDB int, emailConfig *config.EmailConfig, logger *zap.Logger) *Worker {
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

	emailService := utils.NewEmailService(emailConfig)

	return &Worker{
		server:       server,
		logger:       logger,
		emailService: emailService,
	}
}

// Start starts the worker to process tasks
func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(string(TaskSendWelcomeEmail), w.handleSendWelcomeEmail)
	mux.HandleFunc(string(TaskSendPasswordResetEmail), w.handleSendPasswordResetEmail)

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

	w.logger.Info("Sending welcome email",
		zap.Uint("user_id", payload.UserID),
		zap.String("username", payload.Username),
		zap.String("email", payload.Email))

	// Create welcome email content
	subject := "Welcome to Kaleidoscope!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome, %s!</h2>
			<p>Thank you for registering with Kaleidoscope. We're excited to have you on board!</p>
			<p>Your user ID is: %d</p>
			<p>Best regards,<br>The Kaleidoscope Team</p>
		</body>
		</html>
	`, payload.Username, payload.UserID)

	// Send the email using EmailService
	if err := w.emailService.SendEmail(payload.Email, subject, body); err != nil {
		w.logger.Error("Failed to send welcome email",
			zap.Uint("user_id", payload.UserID),
			zap.String("email", payload.Email),
			zap.Error(err))
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	w.logger.Info("Welcome email sent successfully",
		zap.Uint("user_id", payload.UserID),
		zap.String("email", payload.Email))

	return nil
}

func (w *Worker) handleSendPasswordResetEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendPasswordResetEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	w.logger.Info("Sending password reset email",
		zap.Uint("user_id", payload.UserID),
		zap.String("username", payload.Username),
		zap.String("email", payload.Email))

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", w.emailService.GetFrontendURL(), payload.Token)
	subject := "Reset Your Password - Kaleidoscope"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Hello, %s!</h2>
			<p>We received a request to reset your password.</p>
			<p>Click the link below to reset your password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you did not request a password reset, please ignore this email.</p>
			<p>Best regards,<br>The Kaleidoscope Team</p>
		</body>
		</html>
	`, payload.Username, resetLink)

	if err := w.emailService.SendEmail(payload.Email, subject, body); err != nil {
		w.logger.Error("Failed to send password reset email",
			zap.Uint("user_id", payload.UserID),
			zap.String("email", payload.Email),
			zap.Error(err))
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	w.logger.Info("Password reset email sent successfully",
		zap.Uint("user_id", payload.UserID),
		zap.String("email", payload.Email))

	return nil
}
