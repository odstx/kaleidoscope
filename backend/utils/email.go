package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"kaleidoscope/config"
)

// EmailService handles email sending functionality
type EmailService struct {
	config *config.EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// GetFrontendURL returns the frontend URL from the email config
func (e *EmailService) GetFrontendURL() string {
	return e.config.FrontendURL
}

// SendEmail sends an email using the configured SMTP settings
func (e *EmailService) SendEmail(to, subject, body string) error {
	if e.config.Username == "" || e.config.Password == "" || e.config.From == "" {
		return fmt.Errorf("email configuration is incomplete - username, password, and from address are required")
	}

	// Create the authentication
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)

	// Create the message
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", e.config.From, to, subject, body))

	var err error
	// For QQ mail and other services that require SSL from the start (port 465)
	if e.config.Port == 465 {
		// Connect with SSL (implicit TLS)
		tlsConfig := &tls.Config{
			ServerName: e.config.Host,
		}
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", e.config.Host, e.config.Port), tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server with SSL: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, e.config.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
		}

		if err = client.Mail(e.config.From); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to send data command: %w", err)
		}

		if _, err = w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err = w.Close(); err != nil {
			return fmt.Errorf("failed to close message writer: %w", err)
		}

		if err = client.Quit(); err != nil {
			return fmt.Errorf("failed to quit SMTP client: %w", err)
		}
	} else if e.config.UseTLS {
		// Connect with STARTTLS (explicit TLS, typically port 587)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", e.config.Host, e.config.Port))
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, e.config.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		// Start TLS
		if err = client.StartTLS(&tls.Config{ServerName: e.config.Host}); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
		}

		if err = client.Mail(e.config.From); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to send data command: %w", err)
		}

		if _, err = w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err = w.Close(); err != nil {
			return fmt.Errorf("failed to close message writer: %w", err)
		}

		if err = client.Quit(); err != nil {
			return fmt.Errorf("failed to quit SMTP client: %w", err)
		}
	} else {
		// Connect without TLS
		err = smtp.SendMail(fmt.Sprintf("%s:%d", e.config.Host, e.config.Port), auth, e.config.From, []string{to}, msg)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}
