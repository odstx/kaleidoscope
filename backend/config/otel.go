package config

type OTELConfig struct {
	Enabled           bool               `mapstructure:"enabled"`
	ServiceName       string             `mapstructure:"service_name"`
	CollectorURL      string             `mapstructure:"collector_url"`
	TracesExporter    string             `mapstructure:"traces_exporter"`
	MetricsExporter   string             `mapstructure:"metrics_exporter"`
	LogsExporter      string             `mapstructure:"logs_exporter"`
	SamplingRate      float64            `mapstructure:"sampling_rate"`
	PropagationFormat string             `mapstructure:"propagation_format"`
	Headers           []OTELHeaderConfig `mapstructure:"headers"`
}

type OTELHeaderConfig struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}
