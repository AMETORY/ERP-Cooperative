package config

type WhatsappConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	MockNumber string `mapstructure:"mock_number"`
	IsMock     bool   `mapstructure:"is_mock"`
}
