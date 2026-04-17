package config

type LogConfig struct {
	EnableConsole bool   `mapstructure:"enable_console"`
	EnableFile    bool   `mapstructure:"enable_file"`
	FilePath      string `mapstructure:"file_path"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}
