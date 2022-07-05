package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ginOptions struct {
	Debug    bool
	HttpMode string `mapstructure:"http_mode"` // debug | release
	HttpHost string `mapstructure:"http_host"`
	HttpPort int    `mapstructure:"http_port"`
}

func (g ginOptions) setDefaults() {
	viper.SetDefault("HTTP_MODE", gin.ReleaseMode)
	viper.SetDefault("HTTP_HOST", "0.0.0.0")
	viper.SetDefault("HTTP_PORT", 8080)
}

func (c *Config) HttpMode() string {
	return c.options.Gin.HttpMode
}

func (c *Config) HttpHost() string {
	return c.options.Gin.HttpHost
}

func (c *Config) HttpPort() int {
	return c.options.Gin.HttpPort
}
