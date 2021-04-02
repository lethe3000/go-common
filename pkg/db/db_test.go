package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dbconfig struct {
	driver   string
	dsn      string
	user     string
	password string
	name     string
	server   string
}

func (d dbconfig) DatabaseDriver() string {
	return d.driver
}

func (d dbconfig) DatabaseDsn() string {
	return d.dsn
}

func (d dbconfig) DatabaseUser() string {
	return d.user
}

func (d dbconfig) DatabasePassword() string {
	return d.password
}

func (d dbconfig) DatabaseName() string {
	return d.name
}

func (d dbconfig) DatabaseServer() string {
	return d.server
}

func (d dbconfig) DatabaseConns() int {
	return 1
}

func (d dbconfig) DatabaseConnsIdle() int {
	return 5
}

func (d dbconfig) Echo() bool {
	return true
}

func TestDbInit(t *testing.T) {
	config := dbconfig{
		driver:   MySQL,
		dsn:      "",
		user:     "root",
		password: "root",
		name:     "feishu",
		server:   "localhost:3306",
	}
	InitDb(config)
	assert.NotNil(t, db)
}

func TestInitTestDb(t *testing.T) {
	config := dbconfig{
		driver: SQLite,
		dsn:    SqliteMemoryDsn,
	}
	InitDb(config)
}
