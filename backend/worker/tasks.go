package worker

// TaskType represents the type of task
type TaskType string

const (
	// TaskSendWelcomeEmail represents a task to send a welcome email
	TaskSendWelcomeEmail TaskType = "send_welcome_email"
)

// Payload fields for TaskSendWelcomeEmail
type SendWelcomeEmailPayload struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
