package database

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/carp-cobain/tracker-pg/database/model"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectAndMigrate connects to a database and runs migrations using project models.
func ConnectAndMigrate() (*gorm.DB, error) {
	dialect, dsn := lookupConnectParams()
	db, err := Connect(dialect, dsn, max(4, runtime.NumCPU()))
	if err != nil {
		return nil, err
	}
	if err = RunMigrations(db); err != nil {
		return nil, err
	}
	return db, nil
}

// Connect to a sqlite3 database.
func Connect(dialect, dsn string, maxConns int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	config := &gorm.Config{
		Logger: logger.Discard, //Default.LogMode(logger.Info),
	}
	if dialect == "sqlite3" {
		db, err = gorm.Open(sqlite.Open(dsn), config)
	} else {
		db, err = gorm.Open(postgres.Open(dsn), config)
	}
	if err != nil {
		return nil, err
	}
	if dialect == "sqlite3" {
		if err = execPragmas(db); err != nil {
			log.Printf("unable to execute PRAGMA statements: %+v", err)
		}
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(maxConns)
	}
	return db, nil
}

// Run migrations on a database using project models.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&model.Campaign{}, &model.Referral{})
}

// Optimizations for running SQLite in production.
func execPragmas(db *gorm.DB) error {
	stmts := []string{
		"journal_mode = WAL",
		"busy_timeout = 5000",
		"synchronous = NORMAL",
		"cache_size = 1000000000",
		"foreign_keys = true",
		"temp_store = memory",
		"wal_autocheckpoint = 0",
	}
	for _, stmt := range stmts {
		if err := db.Exec(fmt.Sprintf("PRAGMA %s;", stmt)).Error; err != nil {
			return err
		}
	}
	return nil
}

// Lookup db connection params from env vars
func lookupConnectParams() (dialect string, dsn string) {
	var ok bool
	if dialect, ok = os.LookupEnv("DB_DIALECT"); !ok {
		log.Panicf("DB_DIALECT not set")
	}
	if dsn, ok = os.LookupEnv("DB_DSN"); !ok {
		log.Panicf("DB_DSN not set")
	}
	return
}
