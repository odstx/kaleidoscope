package config

type HawkConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	TimestampSkewSecs int  `mapstructure:"timestamp_skew_secs"`
}
