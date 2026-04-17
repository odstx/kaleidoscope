package config

type SecurityConfig struct {
	MaxLoginAttempts    int `mapstructure:"max_login_attempts"`
	LockoutDurationMins int `mapstructure:"lockout_duration_mins"`
}
