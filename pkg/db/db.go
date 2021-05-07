package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

const (
	MySQL           = "mysql"
	MariaDB         = "mariadb"
	SQLite          = "sqlite3"
	SqliteMemoryDsn = "file::memory:?cache=shared"
)

var (
	db                *gorm.DB
	once              sync.Once
	UnsupportedDriver = errors.New("不支持的数据库driver")
)

type dbConfig interface {
	DatabaseDriver() string
	DatabaseDsn() string
	DatabaseUser() string
	DatabasePassword() string
	DatabaseName() string
	DatabaseServer() string
	DatabaseConns() int
	DatabaseConnsIdle() int
	Echo() bool
}

func databaseDriver(c dbConfig) string {
	var driver string
	switch strings.ToLower(c.DatabaseDriver()) {
	case MySQL, MariaDB:
		driver = MySQL
	case SQLite, "sqlite", "sqllite", "file", "":
		driver = SQLite
	default:
		log.Printf("config: unsupported database driver %s, using sqlite", c.DatabaseDriver())
		driver = SQLite
	}

	return driver
}

func databaseDsn(c dbConfig) string {
	if c.DatabaseDsn() == "" {
		switch databaseDriver(c) {
		case MySQL, MariaDB:
			return fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
				c.DatabaseUser(),
				c.DatabasePassword(),
				c.DatabaseServer(),
				c.DatabaseName(),
			)
		case SQLite:
			return "app.db"
		default:
			log.Fatalf("config: empty database dsn")
		}
	}
	return c.DatabaseDsn()
}

func conns(c dbConfig) int {
	limit := c.DatabaseConns()
	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}
	if limit > 1024 {
		limit = 1024
	}
	return limit
}

func idle(c dbConfig) int {
	limit := c.DatabaseConnsIdle()
	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}
	if limit > c.DatabaseConns() {
		limit = c.DatabaseConns()
	}
	return limit
}

func getDialector(c dbConfig) (gorm.Dialector, error) {
	var dialector gorm.Dialector
	dsn := databaseDsn(c)
	switch databaseDriver(c) {
	case MySQL:
		dialector = mysql.Open(dsn)
	case SQLite:
		dialector = sqlite.Open(dsn)
	default:
		return nil, UnsupportedDriver
	}
	return dialector, nil
}

func InitDb(c dbConfig) {
	once.Do(func() {
		var err error

		gormConfig := gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
		dialector, err := getDialector(c)
		if err != nil {
			panic(err)
		}

		db, err = gorm.Open(dialector, &gormConfig)
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxIdleConns(idle(c))
		sqlDB.SetMaxOpenConns(conns(c))
		sqlDB.SetConnMaxLifetime(time.Hour)
	})
}

func InitTestDb() {
	var err error
	_ = os.Remove("test.db")
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("init test db err: %v", err)
	}
}

func InitTestMemoryDb() {
	var err error
	db, err = gorm.Open(sqlite.Open(SqliteMemoryDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("init test db err: %v", err)
	}
}

func DB() *gorm.DB {
	if db == nil {
		panic("db instance not initialized")
	}
	return db
}
