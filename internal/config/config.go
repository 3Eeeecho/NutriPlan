package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 结构体包含所有配置
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Nutrients NutrientsConfig `mapstructure:"nutrients"`
}

// ServerConfig 包含了服务器相关的配置
type ServerConfig struct {
	Port      int    `mapstructure:"port"`
	Env       string `mapstructure:"env"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

// DatabaseConfig 包含了数据库连接相关的配置
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
}

// NutrientsConfig 包含了营养计算相关的配置
type NutrientsConfig struct {
	// 宏量营养素热量因子 (kcal/g)
	ProteinFactor float64 `mapstructure:"protein_factor"` // 蛋白质每克热量 (4.0 kcal)
	CarbFactor    float64 `mapstructure:"carb_factor"`    // 碳水每克热量 (4.0 kcal)
	FatFactor     float64 `mapstructure:"fat_factor"`     // 脂肪每克热量 (9.0 kcal)

	// 目标热量调整值 (用于 TDEE 调整)
	LossCalorieDeficit    float64 `mapstructure:"loss_calorie_deficit"`    // 减脂赤字 (如 400)
	GainCalorieSurplus    float64 `mapstructure:"gain_calorie_surplus"`    // 增肌盈余 (如 300)
	ControlCalorieDeficit float64 `mapstructure:"control_calorie_deficit"` // 控糖赤字 (如 200)
	MinCalorie            float64 `mapstructure:"min_calorie"`             // 最低安全热量 (如 1200)
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
	v.AddConfigPath(".")

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
