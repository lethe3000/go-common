package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

const YamlFileType = "yaml"

var (
	settings *config
	once     sync.Once
)

type config struct {
	options *options
}

type options struct {
	Server *serverOptions `mapstructure:"server"`
	DB     *dbOptions     `mapstructure:"db"`
	Gin    *ginOptions    `mapstructure:"gin"`
}

func InitSettings(configName string) *config {
	o := options{
		Server: &serverOptions{},
		DB:     &dbOptions{},
		Gin:    &ginOptions{},
	}
	viper.SetConfigName(configName)
	viper.SetConfigType(YamlFileType)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	o.Server.setDefaults()
	o.DB.setDefaults()
	o.Gin.setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err read config: %v\n", err)
	}
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&o); err != nil {
		log.Fatalf("unmarshal config fail: %v", err)
	}
	return &config{options: &o}
}

func SetConfig(c *config) {
	once.Do(func() {
		settings = c
	})
}

func Settings() *config {
	if settings == nil {
		panic("settings not initialized")
	}
	return settings
}
