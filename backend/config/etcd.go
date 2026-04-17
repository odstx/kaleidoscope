package config

type EtcdConfig struct {
	Endpoints   []string `mapstructure:"endpoints"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	DialTimeout int      `mapstructure:"dial_timeout"`
}
