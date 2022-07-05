package config

import "github.com/spf13/viper"

type dbOptions struct {
	Echo              bool
	DatabaseDriver    string `mapstructure:"database_driver"`
	DatabaseServer    string `mapstructure:"database_server"`
	DatabaseUser      string `mapstructure:"database_user"`
	DatabasePassword  string `mapstructure:"database_password"`
	DatabaseName      string `mapstructure:"database_name"`
	DatabaseDsn       string `mapstructure:"database_dsn"`
	DatabaseConnsIdle int    `mapstructure:"database_conns_idle"`
	DatabaseConns     int    `mapstructure:"database_conns"`
}

func (db dbOptions) setDefaults() {
	viper.SetDefault("ECHO", false)
	viper.SetDefault("DATABASE_DRIVER", "mysql")
	viper.SetDefault("DATABASE_SERVER", "localhost")
	viper.SetDefault("DATABASE_USER", "root")
	viper.SetDefault("DATABASE_PASSWORD", "root")
	viper.SetDefault("DATABASE_NAME", "workflow")
	viper.SetDefault("DATABASE_DSN", "")
	viper.SetDefault("DATABASE_CONNS_IDLE", 3)
	viper.SetDefault("DATABASE_CONNS", 10)
}

func (c *Config) DatabaseDriver() string {
	return c.options.DB.DatabaseDriver
}

func (c *Config) DatabaseServer() string {
	return c.options.DB.DatabaseServer
}

func (c *Config) DatabaseUser() string {
	return c.options.DB.DatabaseServer
}

func (c *Config) DatabasePassword() string {
	return c.options.DB.DatabasePassword
}

func (c *Config) DatabaseName() string {
	return c.options.DB.DatabaseName
}

func (c *Config) DatabaseDsn() string {
	return c.options.DB.DatabaseDsn
}

func (c *Config) DatabaseConnsIdle() int {
	return c.options.DB.DatabaseConnsIdle
}

func (c *Config) DatabaseConns() int {
	return c.options.DB.DatabaseConns
}

func (c *Config) Echo() bool {
	return c.options.DB.Echo
}
