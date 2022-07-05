package db

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"gorm.io/gorm/logger"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	MySQL           = "mysql"
	MariaDB         = "mariadb"
	SQLite          = "sqlite3"
	SqliteMemoryDsn = "file::memory:?cache=shared"
	SqliteFileDsn   = "file:test.db?cache=shared"
)

type DatabaseConfig interface {
	DatabaseConns() int
	DatabaseConnsIdle() int
	DatabaseDsn() string
	DatabaseDriver() string

	DatabaseUser() string
	DatabasePassword() string
	DatabaseServer() string
	DatabaseName() string
}

type Database struct {
	*gorm.DB
	config DatabaseConfig
	logger logger.Interface
}

func NewTestSqlite(logger logger.Interface) (Database, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return Database{}, err
	}
	db = db.Debug()
	db.Exec("PRAGMA foreign_keys = ON;")
	return Database{DB: db, logger: logger}, nil
}

func NewTestMemoryDataBase(logger logger.Interface) (Database, error) {
	db, err := gorm.Open(sqlite.Open(SqliteMemoryDsn), &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return Database{}, errors.WithMessage(err, "init test db err: %v")
	}
	db = db.Debug()
	db.Exec("PRAGMA foreign_keys = ON")
	return Database{db, nil, logger}, nil
}

func NewDatabase(config DatabaseConfig, logger logger.Interface) (Database, error) {
	d := Database{
		config: config,
		logger: logger,
	}
	ctx := context.TODO()
	d.logger.Info(ctx, "Connecting to database")
	dialector, err := d.getDialector()
	if err != nil {
		d.logger.Error(ctx, "Failed to get database dialector: %+v", err)
		return d, nil
	}

	if db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		d.logger.Error(ctx, "failed to connect to database: %+v", err)
		return d, nil
	} else {
		session, _ := db.DB()
		session.SetMaxIdleConns(d.idle())
		session.SetMaxOpenConns(d.conns())
		d.DB = db
	}
	return d, nil
}

// databaseDriver returns the database driver name
func (d Database) databaseDriver() string {
	var driver string
	switch strings.ToLower(d.config.DatabaseDriver()) {
	case MySQL, MariaDB:
		driver = MySQL
	case SQLite, "sqlite", "sqllite", "file", "":
		driver = SQLite
	default:
		driver = SQLite
	}

	return driver
}

// conns calculates the number of connections to open
func (d Database) conns() int {
	limit := d.config.DatabaseConns()
	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}
	if limit > 1024 {
		limit = 1024
	}
	return limit
}

func (d Database) idle() int {
	limit := d.config.DatabaseConnsIdle()
	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}
	if limit > d.config.DatabaseConns() {
		limit = d.config.DatabaseConns()
	}
	return limit
}

func (d Database) databaseDsn() string {
	if d.config.DatabaseDsn() == "" {
		switch d.databaseDriver() {
		case MySQL, MariaDB:
			return fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
				d.config.DatabaseUser(),
				d.config.DatabasePassword(),
				d.config.DatabaseServer(),
				d.config.DatabaseName(),
			)
		case SQLite:
			return SqliteMemoryDsn
		}
	}
	return d.config.DatabaseDsn()
}

func (d Database) getDialector() (gorm.Dialector, error) {
	var dialector gorm.Dialector
	dsn := d.databaseDsn()
	switch driver := d.databaseDriver(); driver {
	case MySQL:
		dialector = mysql.Open(dsn)
	case SQLite:
		dialector = sqlite.Open(dsn)
	default:
		return nil, errors.Errorf("不支持的数据库driver: %s", driver)
	}
	return dialector, nil
}
