package config

import (
	"os"
	"testing"

	"github.com/lethe3000/go-common/pkg/db"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("test env", func(t *testing.T) {
		_ = os.Setenv("DATABASE_DRIVER", "sqlite3")
		_ = os.Setenv("ECHO", "true")
		_ = os.Setenv("HTTP_PORT", "8888")
		c := InitSettings("config.yaml")
		assert.Equal(t, c.options.DB.DatabaseDriver, db.SQLite)
		assert.Equal(t, c.options.Gin.HttpPort, 8888)
		assert.Equal(t, c.DatabaseDriver(), db.SQLite)
		assert.Equal(t, c.HttpPort(), 8888)

		SetConfig(c)
		assert.NotNil(t, Settings())
	})
}
