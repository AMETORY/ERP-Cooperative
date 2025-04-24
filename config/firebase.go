package config

type FirebaseConfig struct {
	Bucket         string `mapstructure:"bucket"`
	CredentialPath string `mapstructure:"credential_path"`
	ServerKey      string `mapstructure:"server_key"`
}
