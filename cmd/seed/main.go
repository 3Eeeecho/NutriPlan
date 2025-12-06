package main

import (
	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接 (简化版，直接连接)
func InitDB() (*gorm.DB, error) {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	repo := dao.NewGormRecipeRepository(db)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 基础数据
	mealTypes := []models.MealType{
		models.MealTypeBreakfast,
		models.MealTypeLunch,
		models.MealTypeDinner,
		models.MealTypeSnack,
	}

	cookingMethods := []string{"炒", "蒸", "煮", "烤", "煎", "凉拌", "炖"}
	difficulties := []string{"简单", "中等", "困难"}

	// 食材库
	proteins := []string{"鸡胸肉", "牛肉", "猪瘦肉", "虾仁", "豆腐", "鸡蛋", "三文鱼", "鳕鱼"}
	vegetables := []string{"西兰花", "胡萝卜", "菠菜", "青椒", "洋葱", "番茄", "黄瓜", "生菜", "蘑菇"}
	carbs := []string{"糙米", "燕麦", "全麦面包", "红薯", "玉米", "荞麦面", "土豆"}
	fruits := []string{"苹果", "香蕉", "蓝莓", "橙子", "猕猴桃"}
	nuts := []string{"杏仁", "核桃", "腰果"}

	// 生成 100 个食谱
	for i := 0; i < 100; i++ {
		mealType := mealTypes[rng.Intn(len(mealTypes))]
		method := cookingMethods[rng.Intn(len(cookingMethods))]

		var mainIngredient, sideIngredient, carbIngredient string
		var ingredients []string

		// 根据餐次选择食材
		switch mealType {
		case models.MealTypeBreakfast:
			mainIngredient = proteins[rng.Intn(len(proteins))] // 主要是蛋奶
			if rng.Intn(2) == 0 {
				mainIngredient = "鸡蛋"
			}
			sideIngredient = fruits[rng.Intn(len(fruits))]
			carbIngredient = carbs[rng.Intn(len(carbs))]
			ingredients = []string{mainIngredient, sideIngredient, carbIngredient, "牛奶"}
		case models.MealTypeSnack:
			mainIngredient = fruits[rng.Intn(len(fruits))]
			sideIngredient = nuts[rng.Intn(len(nuts))]
			ingredients = []string{mainIngredient, sideIngredient, "酸奶"}
			method = "即食"
		default:
			mainIngredient = proteins[rng.Intn(len(proteins))]
			sideIngredient = vegetables[rng.Intn(len(vegetables))]
			carbIngredient = carbs[rng.Intn(len(carbs))]
			ingredients = []string{mainIngredient, sideIngredient, carbIngredient, "植物油", "盐"}
		}

		name := fmt.Sprintf("%s%s配%s", method, mainIngredient, sideIngredient)
		if mealType == models.MealTypeSnack {
			name = fmt.Sprintf("%s坚果杯", mainIngredient)
		}

		// 随机营养成分 (基于食材估算)
		energy := 300 + rng.Float64()*500
		if mealType == models.MealTypeSnack {
			energy = 100 + rng.Float64()*200
		}

		protein := energy * (0.15 + rng.Float64()*0.2) / 4
		fat := energy * (0.2 + rng.Float64()*0.15) / 9
		carb := (energy - protein*4 - fat*9) / 4

		ingredientsJSON, _ := json.Marshal(ingredients)

		recipe := &models.Recipe{
			Name:         name,
			Description:  fmt.Sprintf("这是一道美味的%s，富含营养。", name),
			ImageURL:     "https://placehold.co/600x400?text=Recipe", // 占位图
			MealType:     mealType,
			Energy:       truncate(energy),
			Protein:      truncate(protein),
			Carbohydrate: truncate(carb),
			Fat:          truncate(fat),
			Ingredients:  string(ingredientsJSON),
			CookingTime:  10 + rng.Intn(50),
			Difficulty:   difficulties[rng.Intn(len(difficulties))],
		}

		if err := repo.Create(recipe); err != nil {
			log.Printf("创建食谱失败: %v", err)
		} else {
			fmt.Printf("已创建食谱: %s (%s)\n", recipe.Name, recipe.MealType)
		}
	}
}

func truncate(val float64) float64 {
	return float64(int(val*100)) / 100
}
