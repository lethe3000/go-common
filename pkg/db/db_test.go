package db

import (
	"testing"

	"gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	t.Parallel()
	t.Run("test db", func(t *testing.T) {
		db, err := NewTestMemoryDataBase(logger.Default)
		assert.NoError(t, err)
		var one int
		db.Raw("select 1").Scan(&one)
		assert.Equal(t, one, 1)
	})
}
