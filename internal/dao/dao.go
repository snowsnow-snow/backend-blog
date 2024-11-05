package dao

import (
	"backend-blog/config"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
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
	initDb(dbConfig config.DBConfig) error
}

type SqliteInit struct {
}

func init() {
	InitSQLCon()
}

// InitSQLCon 初始化持久层的设置，此方法会根据配置文件选择 sqlite3 数据库还是 Mysql 数据库
// 初始化完成后会对数据库的操作都由 DB 完成
func InitSQLCon() {
	var dataBaser DataBaser
	if config.GlobalConfig.DBConfig.Type == "sqlite" {
		dataBaser = new(SqliteInit)
	}
	err := dataBaser.initDb(config.GlobalConfig.DBConfig)
	if err != nil {
		logger.Error.Panic("Error init database ", err)
	}
	err = initTables()
	if err != nil {
		logger.Error.Panic("Error init tables ", err)
	}
}

// init sqlite
func (sqliteInit SqliteInit) initDb(dbConfig config.DBConfig) error {
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
	if !DB.Migrator().HasTable(&entity.User{}) {
		tables = append(tables, &entity.User{})
		err = DB.AutoMigrate(&entity.User{})
	}
	if !DB.Migrator().HasTable(&entity.BlogImage{}) {
		tables = append(tables, &entity.BlogImage{})
		err = DB.AutoMigrate(&entity.BlogImage{})
	}
	if !DB.Migrator().HasTable(&entity.BlogVideo{}) {
		tables = append(tables, &entity.BlogVideo{})
		err = DB.AutoMigrate(&entity.BlogVideo{})
	}
	if !DB.Migrator().HasTable(&entity.BlogContent{}) {
		tables = append(tables, &entity.BlogContent{})
		err = DB.AutoMigrate(&entity.BlogContent{})
	}
	if !DB.Migrator().HasTable(&entity.File{}) {
		tables = append(tables, &entity.File{})
		err = DB.AutoMigrate(&entity.File{})
	}
	if !DB.Migrator().HasTable(&entity.BlogMarkdown{}) {
		tables = append(tables, &entity.BlogMarkdown{})
		err = DB.AutoMigrate(&entity.BlogMarkdown{})
	}
	return err
}
