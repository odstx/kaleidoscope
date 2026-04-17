package config

type LLMConfig struct {
	URL          string `mapstructure:"url"`
	APIKey       string `mapstructure:"api_key"`
	Model        string `mapstructure:"model"`
	SystemPrompt string `mapstructure:"system_prompt"`
}
