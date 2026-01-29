package dao

import (
	"backend-blog/config"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/pkg/logger" // 确保这里引用了你刚才改写的包含 GormLogger 的包
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitSQLCon(cfg config.DBConfig) *gorm.DB {
	// 1. 使用我们自定义的 GormLogger 适配器 (对接 slog)
	newLogger := logger.NewGormLogger()
	// 你可以根据环境动态调整级别，例如：
	// if !cfg.Debug { newLogger.LogLevel = glogger.Warn }

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 关键点：直接传入实现了 glogger.Interface 的结构体
		Logger: newLogger,
		// 开启错误翻译，将 DB 错误转换为 GORM 错误
		TranslateError: true,
		// 开启预编译语句缓存
		PrepareStmt: false,
	}

	// 2. 开启连接
	// 增加 _busy_timeout=5000 (5秒) 避免 Database Locked 错误
	// 增加 _journal_mode=WAL 确保并发读写
	dsn := cfg.Sqlite3.Path + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=1"
	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		slog.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}

	// 3. 配置底层 SQL 连接池
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(5)

	sqlDB.SetConnMaxLifetime(0)

	DB = db
	initTables()
	return DB
}

func initTables() {
	// 禁用外键检查，以便 GORM 可以删除/重建表以进行迁移 (SQLite 限制)
	DB.Exec("PRAGMA foreign_keys = OFF")

	err := DB.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Post{},
		&entity.MediaAsset{},
	)

	// 重新启用外键检查
	DB.Exec("PRAGMA foreign_keys = ON")

	if err != nil {
		slog.Error("Database migration failed", "error", err)
		os.Exit(1)
	}
}
