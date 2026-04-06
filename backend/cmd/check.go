package cmd

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kaleidoscope/config"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check email, Redis, and PostgreSQL configurations",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Test email configuration
		if err := testEmailConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Email configuration test failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Email configuration is valid!")

		// Test PostgreSQL configuration
		if err := testPostgreSQLConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "PostgreSQL configuration test failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ PostgreSQL configuration is valid!")

		// Test Redis configuration
		if err := testRedisConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Redis configuration test failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Redis configuration is valid!")

		fmt.Println("\nAll configurations are valid! 🎉")
	},
}

func testEmailConfig(cfg *config.Config) error {
	// Create SMTP authentication
	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.Host)

	// Create test email message with proper RFC2822 format
	from := cfg.Email.From
	to := []string{cfg.Email.Username} // Send to the same email used for SMTP
	msg := []byte("From: " + from + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: Kaleidoscope Email Configuration Test\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		"This is a test email to verify your email configuration is working correctly.\r\n")

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Email.Host, cfg.Email.Port),
		auth,
		cfg.Email.Username,
		to,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	return nil
}

func testPostgreSQLConfig(cfg *config.Config) error {
	// Build DSN from config
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)

	// Attempt to connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test the connection with timeout
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.ConnectionTimeout)*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return nil
}

func testRedisConfig(cfg *config.Config) error {
	// Create Redis client
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Redis.ConnectionTimeout)*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
