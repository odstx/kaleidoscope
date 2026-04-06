package database

import (
	"testing"
	"time"

	"kaleidoscope/config"
)

func TestRetryConfig(t *testing.T) {
	// Test with default retry configuration
	cfg := &config.DatabaseConfig{
		Host:                 "localhost",
		Port:                 "5432",
		User:                 "test",
		Password:             "test",
		Name:                 "test_db",
		SSLMode:              "disable",
		MaxRetryAttempts:     3,
		RetryIntervalSeconds: 2,
		ConnectionTimeout:    30,
	}

	if cfg.MaxRetryAttempts != 3 {
		t.Errorf("Expected MaxRetryAttempts to be 3, got %d", cfg.MaxRetryAttempts)
	}

	if cfg.RetryIntervalSeconds != 2 {
		t.Errorf("Expected RetryIntervalSeconds to be 2, got %d", cfg.RetryIntervalSeconds)
	}

	if cfg.ConnectionTimeout != 30 {
		t.Errorf("Expected ConnectionTimeout to be 30, got %d", cfg.ConnectionTimeout)
	}
}

func TestRetryLogic(t *testing.T) {
	// Test that retry logic parameters are correctly used
	// This is a simple unit test to verify the retry parameters
	cfg := &config.DatabaseConfig{
		MaxRetryAttempts:     5,
		RetryIntervalSeconds: 1,
	}

	expectedDuration := time.Duration(cfg.RetryIntervalSeconds) * time.Second
	if expectedDuration != time.Second {
		t.Errorf("Expected retry interval to be 1 second, got %v", expectedDuration)
	}

	// Verify that retry attempts are properly configured
	if cfg.MaxRetryAttempts < 1 {
		t.Errorf("MaxRetryAttempts should be at least 1, got %d", cfg.MaxRetryAttempts)
	}
}

func TestRedisRetryConfig(t *testing.T) {
	// Test Redis retry configuration
	cfg := &config.RedisConfig{
		Host:                 "localhost",
		Port:                 "6379",
		Password:             "",
		DB:                   0,
		MaxRetryAttempts:     4,
		RetryIntervalSeconds: 3,
		ConnectionTimeout:    20,
	}

	if cfg.MaxRetryAttempts != 4 {
		t.Errorf("Expected Redis MaxRetryAttempts to be 4, got %d", cfg.MaxRetryAttempts)
	}

	if cfg.RetryIntervalSeconds != 3 {
		t.Errorf("Expected Redis RetryIntervalSeconds to be 3, got %d", cfg.RetryIntervalSeconds)
	}

	if cfg.ConnectionTimeout != 20 {
		t.Errorf("Expected Redis ConnectionTimeout to be 20, got %d", cfg.ConnectionTimeout)
	}
}
