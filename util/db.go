package util

import (
	"backend-blog/config"
	"backend-blog/logger"
	"backend-blog/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// DataBaser 对数据库设置的接口
type DataBaser interface {
	init(dbConfig config.DBConfig) error
}

type SqliteInit struct {
}

// InitSQLCon 初始化持久层的设置，此方法会根据配置文件选择 sqlite3 数据库还是 Mysql 数据库
// 初始化完成后会对数据库的操作都由 DB 完成
func InitSQLCon() {
	var dataBaser DataBaser
	if config.GlobalConfig.DBConfig.Type == "sqlite" {
		dataBaser = new(SqliteInit)
	}
	err := dataBaser.init(config.GlobalConfig.DBConfig)
	if err != nil {
		logger.Error.Panic("Error init database %v", err)
	}
	err = initTables()
	if err != nil {
		logger.Error.Panic("Error init tables %v", err)
	}
}

// init sqlite
func (sqliteInit SqliteInit) init(dbConfig config.DBConfig) error {
	db, err := gorm.Open(sqlite.Open(dbConfig.Sqlite3.Path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: gormLogger.New(
			logger.Info, // 日志输出使用配置好的日志插件
			gormLogger.Config{
				LogLevel: gormLogger.Info,
			},
		),
	})
	DB = db
	return err
}

// initTables 在数据库中新建表
func initTables() error {
	var err error
	var tables []interface{}
	if !DB.Migrator().HasTable(&models.User{}) {
		tables = append(tables, &models.User{})
		err = DB.AutoMigrate(&models.User{})
	}
	if !DB.Migrator().HasTable(&models.ImgInfo{}) {
		tables = append(tables, &models.ImgInfo{})
		err = DB.AutoMigrate(&models.ImgInfo{})
	}
	if !DB.Migrator().HasTable(&models.VideoInfo{}) {
		tables = append(tables, &models.VideoInfo{})
		err = DB.AutoMigrate(&models.VideoInfo{})
	}
	if !DB.Migrator().HasTable(&models.ContentInfo{}) {
		tables = append(tables, &models.ContentInfo{})
		err = DB.AutoMigrate(&models.ContentInfo{})
	}
	if !DB.Migrator().HasTable(&models.ResourceDesc{}) {
		tables = append(tables, &models.ResourceDesc{})
		err = DB.AutoMigrate(&models.ResourceDesc{})
	}
	if !DB.Migrator().HasTable(&models.MarkdownInfo{}) {
		tables = append(tables, &models.MarkdownInfo{})
		err = DB.AutoMigrate(&models.MarkdownInfo{})
	}
	return err
}
