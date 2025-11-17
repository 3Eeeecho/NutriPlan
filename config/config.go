package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 结构体包含所有配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig 包含了服务器相关的配置
type ServerConfig struct {
	Port      int    `mapstructure:"port"`
	Env       string `mapstructure:"env"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

// DatabaseConfig 包含了数据库连接相关的配置
type DatabaseConfig struct {
	Driver    string `mapstructure:"driver"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Host      string `mapstructure:"host"`
	Name      string `mapstructure:"name"`
	Charset   string `mapstructure:"charset"`
	ParseTime string `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}

// 全局配置实例
var AppConfig Config

// LoadConfig 用于初始化并读取配置文件
func LoadConfig() error {
	v := viper.New()

	// 设置配置文件的名称 (无扩展名)
	v.SetConfigName("config")
	// 设置配置文件的类型
	v.SetConfigType("yaml")
	// 设置查找配置文件的路径
	v.AddConfigPath("./config")

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("无法读取配置文件: %w", err)
	}

	// 将配置绑定到 Struct
	if err := v.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("配置解析失败: %w", err)
	}

	return nil
}
