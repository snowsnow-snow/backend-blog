package config

// Config 全剧配置
type Config struct {
	DBConfig DBConfig `yaml:"db"`
	Log      Log      `yaml:"log"`
	File     File     `yaml:"file"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Type    string  `yaml:"type"`
	Sqlite3 Sqlite3 `yaml:"sqlite3"`
	Mysql   Mysql   `yaml:"mysql"`
}

// Mysql 配置
type Mysql struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Sqlite3 配置
type Sqlite3 struct {
	Path string `yaml:"path"`
}
type Log struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	MaxSize int64  `yaml:"maxSize"`
}

type File struct {
	Path Path `yaml:"path"`
}

type Path struct {
	Public string `yaml:"public"`
	Api    string `yaml:"api"`
}

var (
	GlobalConfig Config
)
