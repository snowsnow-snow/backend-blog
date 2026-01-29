package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v3"
)

// GlobalConfig 全局变量
var GlobalConfig Config

type Config struct {
	Server   Server   `yaml:"server"`
	DBConfig DBConfig `yaml:"db"`
	Log      Log      `yaml:"log"`
	File     File     `yaml:"file"`
	ExifTool ExifTool `yaml:"exifTool"`
	JWT      JWT      `yaml:"jwt"`
}

type Server struct {
	Port uint16 `yaml:"port"`
}
type DBConfig struct {
	Type    string  `yaml:"type"`
	Sqlite3 Sqlite3 `yaml:"sqlite3"`
	Mysql   Mysql   `yaml:"mysql"`
}
type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
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
	System   string `yaml:"system"`
	Resource string `yaml:"resource"`
}
type ExifTool struct {
	Path string `yaml:"path"`
}
type JWT struct {
	ExpireHour int    `yaml:"expire_hour"`
	PrivateKey string `yaml:"private_key"`
}

// InitConfig 初始化配置
func InitConfig(env string) {
	// 1. 默认值处理
	if env == "" {
		env = "dev"
	}
	// 建议：统一放在 configs 目录下
	configFileName := fmt.Sprintf("config-%s.yaml", env)
	configPath := filepath.Join("config", configFileName)
	log.Infof("Loading config from: %s", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	log.Info("Config loaded successfully")
}
