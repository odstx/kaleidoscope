package config

type MicroserviceConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	ServiceDomain string `mapstructure:"service_domain"`
}
