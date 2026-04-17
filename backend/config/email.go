package config

type EmailConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	From        string `mapstructure:"from"`
	UseTLS      bool   `mapstructure:"use_tls"`
	FrontendURL string `mapstructure:"frontend_url"`
}
