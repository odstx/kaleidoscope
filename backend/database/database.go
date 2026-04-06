package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kaleidoscope/config"
)

type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// Init initializes both PostgreSQL and Redis connections using the provided config
func Init(cfg *config.Config) (*Database, error) {
	// Initialize PostgreSQL with retry
	db, err := InitPostgreSQLWithRetry(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
		cfg.Database.MaxRetryAttempts,
		time.Duration(cfg.Database.RetryIntervalSeconds)*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL after %d attempts: %w", cfg.Database.MaxRetryAttempts, err)
	}

	// Initialize Redis with retry
	redisClient, err := InitRedisWithRetry(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.MaxRetryAttempts,
		time.Duration(cfg.Redis.RetryIntervalSeconds)*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis after %d attempts: %w", cfg.Redis.MaxRetryAttempts, err)
	}

	return &Database{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func InitPostgreSQL(host, port, user, password, dbname, sslmode string) (*gorm.DB, error) {
	// Use default retry values for backward compatibility
	return InitPostgreSQLWithRetry(host, port, user, password, dbname, sslmode, 5, 5*time.Second)
}

func InitPostgreSQLWithRetry(host, port, user, password, dbname, sslmode string, maxAttempts int, retryInterval time.Duration) (*gorm.DB, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Test connection with timeout
			sqlDB, err := db.DB()
			if err != nil {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()

			if err == nil {
				return db, nil
			}
		}

		if attempt < maxAttempts {
			time.Sleep(retryInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect to postgresql after %d attempts", maxAttempts)
}

func InitRedis(host, port, password string, db int) (*redis.Client, error) {
	// Use default retry values for backward compatibility
	return InitRedisWithRetry(host, port, password, db, 5, 5*time.Second)
}

func InitRedisWithRetry(host, port, password string, db int, maxAttempts int, retryInterval time.Duration) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := rdb.Ping(ctx).Result()
		cancel()

		if err == nil {
			return rdb, nil
		}

		if attempt < maxAttempts {
			time.Sleep(retryInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect to redis after %d attempts", maxAttempts)
}
