package services

import "kaleidoscope/config"

var globalConfig *config.Config

func SetConfig(cfg *config.Config) {
	globalConfig = cfg
}

func GetConfig() *config.Config {
	return globalConfig
}
