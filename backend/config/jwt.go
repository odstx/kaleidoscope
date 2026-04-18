package config

type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	ExpirationHours  int    `mapstructure:"expiration_hours"`
}
