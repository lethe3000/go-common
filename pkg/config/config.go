package config

import (
	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

const YamlFileType = "yaml"

type Config struct {
	options *options
}

type options struct {
	Server serverOptions `mapstructure:"server"`
	DB     dbOptions     `mapstructure:"db"`
	Gin    ginOptions    `mapstructure:"gin"`
}

func NewConfig(configName, configType string) (*Config, error) {
	config := Config{}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New("cannot read configuration")
	}

	err = viper.Unmarshal(&config.options)
	if err != nil {
		return nil, errors.WithMessage(err, "environment can't be loaded: ")
	}

	return &config, nil
}
