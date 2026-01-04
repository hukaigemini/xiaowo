package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置结构
type Config struct {
	DSN        string
	MaxOpen    int
	MaxIdle    int
	Lifetime   time.Duration
	IdleTime   time.Duration
	EnableWAL  bool
	CacheSize  int
	BusyTimeout time.Duration
}

// DefaultConfig 默认数据库配置
func DefaultConfig() *Config {
	return &Config{
		DSN:        getDatabaseDSN(),
		MaxOpen:    25,
		MaxIdle:    10,
		Lifetime:   5 * time.Minute,
		IdleTime:   2 * time.Minute,
		EnableWAL:  true,
		CacheSize:  10000,
		BusyTimeout: 5 * time.Second,
	}
}

// getDatabaseDSN 获取数据库DSN
func getDatabaseDSN() string {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "xiaowo.db"
	}
	
	return fmt.Sprintf("%s?"+
		"cache=shared&"+
		"mode=rwc&"+
		"_journal_mode=WAL&"+
		"_busy_timeout=5000&"+
		"_synchronous=NORMAL&"+
		"_cache_size=10000&"+
		"_temp_store=memory&"+
		"_mmap_size=268435456",
		dbSource)
}

// InitOptimizedDB 初始化优化的数据库连接
func InitOptimizedDB() (*gorm.DB, error) {
	config := DefaultConfig()
	return InitDBWithConfig(config)
}

// InitDBWithConfig 使用指定配置初始化数据库
func InitDBWithConfig(config *Config) (*gorm.DB, error) {
	// 创建数据库连接
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // 生产环境减少日志输出
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层 sql.DB 对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 连接池优化配置
	sqlDB.SetMaxOpenConns(config.MaxOpen)                    // 最大连接数
	sqlDB.SetMaxIdleConns(config.MaxIdle)                    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(config.Lifetime)                // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(config.IdleTime)                // 连接最大空闲时间

	// 额外的SQLite优化配置
	if config.EnableWAL {
		// 设置WAL模式检查点
		_, err = sqlDB.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
		if err != nil {
			log.Printf("WAL checkpoint failed: %v", err)
		}
	}

	// 设置busy timeout
	if config.BusyTimeout > 0 {
		_, err = sqlDB.Exec(fmt.Sprintf("PRAGMA busy_timeout=%d", int(config.BusyTimeout.Milliseconds())))
		if err != nil {
			log.Printf("Set busy timeout failed: %v", err)
		}
	}

	return db, nil
}

// Ping 验证数据库连接
func Ping(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 使用context设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// PingWithRetry 带重试的数据库连接验证
func PingWithRetry(db *gorm.DB, maxRetries int, retryDelay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := Ping(db); err == nil {
			return nil
		} else {
			lastErr = err
			if i < maxRetries-1 {
				log.Printf("数据库连接验证失败 (尝试 %d/%d): %v, %s 后重试", 
					i+1, maxRetries, err, retryDelay)
				time.Sleep(retryDelay)
			}
		}
	}

	return fmt.Errorf("数据库连接验证失败 (尝试 %d 次): %w", maxRetries, lastErr)
}

// GetConnectionStats 获取连接池统计信息
func GetConnectionStats(db *gorm.DB) (*sql.DBStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return &stats, nil
}

// CheckDatabaseHealth 检查数据库健康状态
type DatabaseHealth struct {
	IsHealthy   bool    `json:"is_healthy"`
	Message     string  `json:"message"`
	PingLatency string  `json:"ping_latency,omitempty"`
	ConnStats   *sql.DBStats `json:"connection_stats,omitempty"`
}

// HealthCheck 执行数据库健康检查
func HealthCheck(db *gorm.DB) *DatabaseHealth {
	health := &DatabaseHealth{
		IsHealthy: false,
		Message:   "数据库连接失败",
	}

	if db == nil {
		health.Message = "数据库连接对象为nil"
		return health
	}

	// 测量ping延迟
	start := time.Now()
	if err := Ping(db); err != nil {
		health.Message = fmt.Sprintf("数据库ping失败: %v", err)
		return health
	}
	pingLatency := time.Since(start)

	// 获取连接池统计
	stats, err := GetConnectionStats(db)
	if err != nil {
		log.Printf("获取连接池统计失败: %v", err)
		stats = nil
	}

	// 健康检查通过
	health.IsHealthy = true
	health.Message = "数据库连接正常"
	health.PingLatency = pingLatency.String()
	health.ConnStats = stats

	return health
}

// Close 优雅关闭数据库连接
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// MigrateDatabase 执行数据库迁移
func MigrateDatabase(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 自动迁移数据库模式
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Println("数据库迁移完成")
	return nil
}

// ValidateSchema 验证数据库模式
func ValidateSchema(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 检查关键表是否存在
	tables := []string{
		"user_sessions",
		"rooms", 
		"room_members",
		"room_messages",
		"room_state_history",
		"system_configs",
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			return fmt.Errorf("required table '%s' does not exist", table)
		}
	}

	log.Println("数据库模式验证通过")
	return nil
}