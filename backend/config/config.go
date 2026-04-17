package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for our application
type Config struct {
	Server             ServerConfig       `mapstructure:"server"`
	Database           DatabaseConfig     `mapstructure:"database"`
	Redis              RedisConfig        `mapstructure:"redis"`
	Etcd               EtcdConfig         `mapstructure:"etcd"`
	Log                LogConfig          `mapstructure:"log"`
	CORS               CORSConfig         `mapstructure:"cors"`
	RateLimit          RateLimitConfig    `mapstructure:"rate_limit"`
	Security           SecurityConfig     `mapstructure:"security"`
	OTEL               OTELConfig         `mapstructure:"otel"`
	Hawk               HawkConfig         `mapstructure:"hawk"`
	Email              EmailConfig        `mapstructure:"email"`
	OIDC               OIDCConfig         `mapstructure:"oidc"`
	Microservice       MicroserviceConfig `mapstructure:"microservice"`
	LLM                LLMConfig          `mapstructure:"llm"`
	JWT                JWTConfig          `mapstructure:"jwt"`
	AgentEnabled       bool               `mapstructure:"agent_enabled"`
	EnableRegistration bool               `mapstructure:"enable_registration"`
}

func generateDefaultConfig(path string) error {
	defaultConfig := Config{
		Server: ServerConfig{
			Host:        "",
			Port:        "9000",
			Environment: "development",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Name:     "kaleidoscope",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		Log: LogConfig{
			EnableConsole: true,
			EnableFile:    true,
			FilePath:      "logs/app.log",
			MaxSize:       100,
			MaxBackups:    3,
			MaxAge:        30,
			Compress:      true,
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-UID", "X-User-Name", "X-Version", "X-Source"},
			AllowCredentials: true,
		},
		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 60,
		},
		Security: SecurityConfig{
			MaxLoginAttempts:    5,
			LockoutDurationMins: 15,
		},
		OTEL: OTELConfig{
			Enabled:           false,
			ServiceName:       "kaleidoscope",
			CollectorURL:      "http://localhost:4318",
			TracesExporter:    "otlp",
			MetricsExporter:   "otlp",
			LogsExporter:      "otlp",
			SamplingRate:      1.0,
			PropagationFormat: "w3c",
		},
		Email: EmailConfig{
			Host:        "smtp.gmail.com",
			Port:        587,
			Username:    "",
			Password:    "",
			From:        "",
			UseTLS:      true,
			FrontendURL: "http://localhost:5173",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}
	fmt.Printf("Generated default config at: %s\n", path)
	return nil
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// Set defaults
	viper.SetDefault("server.host", "")
	viper.SetDefault("server.port", "9000")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.static_files_path", "./frontend")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "kaleidoscope")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_retry_attempts", 5)
	viper.SetDefault("database.retry_interval_seconds", 5)
	viper.SetDefault("database.connection_timeout", 5)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retry_attempts", 5)
	viper.SetDefault("redis.retry_interval_seconds", 5)
	viper.SetDefault("redis.connection_timeout", 5)
	viper.SetDefault("etcd.endpoints", []string{"localhost:2379"})
	viper.SetDefault("etcd.username", "")
	viper.SetDefault("etcd.password", "")
	viper.SetDefault("etcd.dial_timeout", 5)
	viper.SetDefault("log.enable_console", true)
	viper.SetDefault("log.enable_file", true)
	viper.SetDefault("log.file_path", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"*"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_minute", 60)
	viper.SetDefault("security.max_login_attempts", 5)
	viper.SetDefault("security.lockout_duration_mins", 15)
	viper.SetDefault("otel.enabled", false)
	viper.SetDefault("otel.service_name", "kaleidoscope")
	viper.SetDefault("otel.collector_url", "http://localhost:4318")
	viper.SetDefault("otel.traces_exporter", "otlp")
	viper.SetDefault("otel.metrics_exporter", "otlp")
	viper.SetDefault("otel.logs_exporter", "otlp")
	viper.SetDefault("otel.sampling_rate", 1.0)
	viper.SetDefault("otel.propagation_format", "w3c")
	viper.SetDefault("otel.headers", []interface{}{})
	viper.SetDefault("email.host", "smtp.gmail.com")
	viper.SetDefault("email.port", 587)
	viper.SetDefault("email.username", "")
	viper.SetDefault("email.password", "")
	viper.SetDefault("email.from", "")
	viper.SetDefault("email.use_tls", true)
	viper.SetDefault("email.frontend_url", "http://localhost:5173")
	viper.SetDefault("oidc.enabled", false)
	viper.SetDefault("oidc.issuer_url", "")
	viper.SetDefault("oidc.client_id", "")
	viper.SetDefault("oidc.client_secret", "")
	viper.SetDefault("oidc.redirect_uri", "http://localhost:9000/api/v1/users/oidc/callback")
	viper.SetDefault("oidc.scopes", []string{"openid", "profile", "email"})
	viper.SetDefault("microservice.enabled", false)
	viper.SetDefault("microservice.service_domain", "service")
	viper.SetDefault("agent_enabled", true)
	viper.SetDefault("enable_registration", true)
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			generatePath := configPath
			if generatePath == "" {
				generatePath = "config.yaml"
			}
			if err := generateDefaultConfig(generatePath); err != nil {
				return nil, err
			}
			viper.SetConfigFile(generatePath)
			if err := viper.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read generated config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("config file found but another error occurred: %w", err)
		}
	}

	// Override with environment variables
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if config.JWT.Secret == "your-secret-key-change-in-production" {
		return nil, fmt.Errorf("jwt.secret must be changed from default value in production")
	}

	return &config, nil
}
