package config

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
