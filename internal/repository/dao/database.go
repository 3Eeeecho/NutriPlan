package dao

import (
	"fmt"
	"log"
	"time"

	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 变量将持有 GORM 数据库连接实例
var DB *gorm.DB

// InitDatabase 初始化 GORM 数据库连接
func InitDatabase(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Name,
	)

	// 使用 GORM 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 获取底层的 *sql.DB，并设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层DB失败: %w", err)
	}

	// 设定连接池的最大闲置连接数
	sqlDB.SetMaxIdleConns(10)
	// 设定连接池的最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设定连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 尝试 PING 数据库以确认连接有效
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接 PING 失败: %w", err)
	}

	log.Println("数据库连接成功!")
	DB = db // 将连接实例赋值给全局变量

	// 数据库模型迁移
	err = db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.DailyRecipePlan{},
		&models.UserWeightLog{},
		&models.DailyIntakeRecord{},
		&models.UserFavoriteRecipe{},
		&models.ShoppingList{},
	)
	if err != nil {
		log.Fatalf("数据库模型迁移失败: %v", err)
		return err
	}
	log.Println("数据库模型迁移完成 (User, Recipe, DailyRecipePlan, UserFavoriteRecipe)。")

	return nil
}
