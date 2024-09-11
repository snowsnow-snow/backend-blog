package config

import (
	"flag"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v3"
	"os"
)

func init() {
	env := flag.String("env", os.Getenv("GO_ENV"), "Application environment (dev|staging|prod)")
	flag.Parse()
	var filepath string
	if *env == "" {
		*env = "dev"
		filepath = "config/config-" + *env + ".yaml"
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal("Error:", err)
			return
		}
		filepath = dir + "/config-" + *env + ".yaml"
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error read config file", err)
	}
	// yaml 文件内容影射到结构体中
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatal("Error init config", err)
	}

}

// Config 全剧配置
type Config struct {
	DBConfig DBConfig `yaml:"db"`
	Log      Log      `yaml:"log"`
	File     File     `yaml:"file"`
	ExifTool ExifTool `yaml:"exifTool"`
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
}

type ExifTool struct {
	Path string `yaml:"path"`
}

var (
	GlobalConfig Config
)
