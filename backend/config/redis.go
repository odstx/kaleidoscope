package config

type RedisConfig struct {
	Host                 string `mapstructure:"host"`
	Port                 string `mapstructure:"port"`
	Password             string `mapstructure:"password"`
	DB                   int    `mapstructure:"db"`
	MaxRetryAttempts     int    `mapstructure:"max_retry_attempts"`
	RetryIntervalSeconds int    `mapstructure:"retry_interval_seconds"`
	ConnectionTimeout    int    `mapstructure:"connection_timeout"`
}
