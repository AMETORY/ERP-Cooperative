package config

type GoogleConfig struct {
	APIKey       string `mapstructure:"api_key"`
	GeminiAPIKey string `mapstructure:"gemini_api_key"`
}
