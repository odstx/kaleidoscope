package config

type ServerConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Environment     string `mapstructure:"environment"`
	StaticFilesPath string `mapstructure:"static_files_path"`
}
