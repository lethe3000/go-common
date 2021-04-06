package config

type serverOptions struct {
	SecretKey  string `mapstructure:"secret_key"`
	ServerName string `mapstructure:"server_name"`
	ServerUrl  string `mapstructure:"server_url"`
}

func (server serverOptions) setDefaults() {
}
