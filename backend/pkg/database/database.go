package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initializes the SQLite database with performance optimizations
func Init(dbPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Configure custom logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// --- SQLite Performance Optimization (Crucial for Go + GORM) ---
	
	// 1. Enable WAL Mode (Write-Ahead Logging)
	// This allows concurrent readers and writers, significantly improving performance.
	if err := DB.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		return err
	}

	// 2. Set Synchronous Mode to NORMAL
	// In WAL mode, NORMAL is safe and much faster than FULL.
	if err := DB.Exec("PRAGMA synchronous = NORMAL;").Error; err != nil {
		return err
	}

	// 3. Configure Connection Pool
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// SetMaxOpenConns(1) is recommended for SQLite to avoid "database is locked" errors
	// during concurrent write operations. WAL mode allows non-blocking reads even with this setting.
	sqlDB.SetMaxOpenConns(1)
	
	// SetMaxIdleConns should match MaxOpenConns or be slightly higher to keep the connection open
	sqlDB.SetMaxIdleConns(1)
	
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database initialized successfully with WAL mode and single-writer configuration")
	return nil
}
