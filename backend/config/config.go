package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds all configuration for our application
type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Log          LogConfig          `mapstructure:"log"`
	CORS         CORSConfig         `mapstructure:"cors"`
	RateLimit    RateLimitConfig    `mapstructure:"rate_limit"`
	OTEL         OTELConfig         `mapstructure:"otel"`
	Hawk         HawkConfig         `mapstructure:"hawk"`
	Email        EmailConfig        `mapstructure:"email"`
	OIDC         OIDCConfig         `mapstructure:"oidc"`
	Microservice MicroserviceConfig `mapstructure:"microservice"`
	LLM          LLMConfig          `mapstructure:"llm"`
}

type ServerConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Environment     string `mapstructure:"environment"`
	StaticFilesPath string `mapstructure:"static_files_path"`
}

type DatabaseConfig struct {
	Host                 string `mapstructure:"host"`
	Port                 string `mapstructure:"port"`
	User                 string `mapstructure:"user"`
	Password             string `mapstructure:"password"`
	Name                 string `mapstructure:"name"`
	SSLMode              string `mapstructure:"sslmode"`
	MaxRetryAttempts     int    `mapstructure:"max_retry_attempts"`
	RetryIntervalSeconds int    `mapstructure:"retry_interval_seconds"`
	ConnectionTimeout    int    `mapstructure:"connection_timeout"`
}

type RedisConfig struct {
	Host                 string `mapstructure:"host"`
	Port                 string `mapstructure:"port"`
	Password             string `mapstructure:"password"`
	DB                   int    `mapstructure:"db"`
	MaxRetryAttempts     int    `mapstructure:"max_retry_attempts"`
	RetryIntervalSeconds int    `mapstructure:"retry_interval_seconds"`
	ConnectionTimeout    int    `mapstructure:"connection_timeout"`
}

type LogConfig struct {
	EnableConsole bool   `mapstructure:"enable_console"`
	EnableFile    bool   `mapstructure:"enable_file"`
	FilePath      string `mapstructure:"file_path"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
}

type OTELConfig struct {
	Enabled           bool               `mapstructure:"enabled"`
	ServiceName       string             `mapstructure:"service_name"`
	CollectorURL      string             `mapstructure:"collector_url"`
	TracesExporter    string             `mapstructure:"traces_exporter"`
	MetricsExporter   string             `mapstructure:"metrics_exporter"`
	LogsExporter      string             `mapstructure:"logs_exporter"`
	SamplingRate      float64            `mapstructure:"sampling_rate"`
	PropagationFormat string             `mapstructure:"propagation_format"`
	Headers           []OTELHeaderConfig `mapstructure:"headers"`
}

type OTELHeaderConfig struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type HawkConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	TimestampSkewSecs int  `mapstructure:"timestamp_skew_secs"`
}

type EmailConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	From        string `mapstructure:"from"`
	UseTLS      bool   `mapstructure:"use_tls"`
	FrontendURL string `mapstructure:"frontend_url"`
}

type OIDCConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	IssuerURL    string   `mapstructure:"issuer_url"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURI  string   `mapstructure:"redirect_uri"`
	Scopes       []string `mapstructure:"scopes"`
}

type MicroserviceConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	ServiceDomain string `mapstructure:"service_domain"`
}

type LLMConfig struct {
	URL          string `mapstructure:"url"`
	APIKey       string `mapstructure:"api_key"`
	Model        string `mapstructure:"model"`
	SystemPrompt string `mapstructure:"system_prompt"`
}

func generateDefaultConfig(path string) error {
	config := `server:
  host: ""
  port: "9000"
  environment: "development"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  name: "kaleidoscope"
  sslmode: "disable"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

log:
  enable_console: true
  enable_file: true
  file_path: "logs/app.log"
  max_size: 100
  max_backups: 3
  max_age: 30
  compress: true

cors:
  allow_origins:
    - "*"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Accept"
    - "Authorization"
  allow_credentials: true

rate_limit:
  enabled: true
  requests_per_minute: 60

otel:
  enabled: false
  service_name: "kaleidoscope"
  collector_url: "http://localhost:4318"
  traces_exporter: "otlp"
  metrics_exporter: "otlp"
  logs_exporter: "otlp"
  sampling_rate: 1.0
  propagation_format: "w3c"
  headers: []

email:
  host: "smtp.gmail.com"
  port: 587
  username: ""
  password: ""
  from: ""
  use_tls: true
  frontend_url: "http://localhost:5173"
`
	if err := os.WriteFile(path, []byte(config), 0644); err != nil {
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
	viper.SetDefault("log.enable_console", true)
	viper.SetDefault("log.enable_file", true)
	viper.SetDefault("log.file_path", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_minute", 60)
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

	return &config, nil
}

// InitLogger initializes the zap logger
func InitLogger(cfg *Config) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	if cfg.Log.EnableConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel))
	}

	if cfg.Log.EnableFile {
		if err := os.MkdirAll(cfg.Log.FilePath[:len(cfg.Log.FilePath)-len("app.log")], 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel))
	}

	if len(cores) == 0 {
		return zap.NewNop(), nil
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}
