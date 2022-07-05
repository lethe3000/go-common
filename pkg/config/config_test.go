package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("test env", func(t *testing.T) {
		_ = os.Setenv("DATABASE_DRIVER", "sqlite3")
		_ = os.Setenv("ECHO", "true")
		_ = os.Setenv("HTTP_PORT", "8888")
		c, err := NewConfig("config", "yaml")
		assert.NoError(t, err)
		assert.Equal(t, c.options.Gin.HttpPort, 8888)
	})
}
