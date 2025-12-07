package main

import (
	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库 (保持不变)
func InitDB() (*gorm.DB, error) {
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

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用日志以提高插入速度
		Logger: nil,
	})
}

// ---------------- 基础数据池 ----------------

var (
	// 风味/形容词
	flavors = []string{"香辣", "蒜蓉", "红烧", "清淡", "黑椒", "柠檬", "咖喱", "照烧", "椒盐", "蜜汁", "宫保", "鱼香", "麻辣", "孜然"}

	// 中式烹饪
	cnMethods = []string{"爆炒", "清蒸", "红烧", "干煸", "回锅", "炖", "凉拌", "白灼", "滑蛋"}
	// 西式烹饪
	enMethods = []string{"香煎", "慢烤", "水煮", "炙烤", "低温慢煮"}

	// 主料池
	meats     = []string{"鸡胸肉", "牛里脊", "猪瘦肉", "去皮鸭腿", "羊排", "三文鱼", "龙利鱼", "虾仁", "牛肉丸"}
	veggies   = []string{"西兰花", "娃娃菜", "菠菜", "胡萝卜", "西葫芦", "茄子", "青椒", "口蘑", "芦笋", "秋葵", "生菜"}
	staples   = []string{"糙米饭", "荞麦面", "红薯", "玉米", "全麦意面", "藜麦", "紫薯"}
	breakfast = []string{"全麦面包", "燕麦片", "鸡蛋", "低脂牛奶", "无糖豆浆", "全麦卷饼"}
	fruits    = []string{"蓝莓", "香蕉", "苹果", "草莓", "猕猴桃", "火龙果", "橙子"}
	nuts      = []string{"巴旦木", "核桃仁", "腰果", "奇亚籽", "南瓜子"}
)

// ---------------- 工具函数 ----------------

func randElem(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func truncate(val float64) float64 {
	return math.Round(val*100) / 100
}

// 生成 JSON 字符串
func makeJSON(data []string) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// ---------------- 核心生成逻辑 ----------------

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}
	repo := dao.NewGormRecipeRepository(db)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	targetCount := 2000
	currentCount := 0
	uniqueNames := make(map[string]bool) // 用于去重

	log.Println("开始生成真实感食谱数据...")

	for currentCount < targetCount {
		var r *models.Recipe

		// 按照一定比例生成不同类型的菜谱
		dice := rng.Intn(100)
		if dice < 15 {
			r = genBreakfast() // 15% 早餐
		} else if dice < 25 {
			r = genSnack() // 10% 零食
		} else if dice < 60 {
			r = genChineseDish() // 35% 中式正餐
		} else if dice < 85 {
			r = genWesternDish() // 25% 西式/轻食
		} else {
			r = genStaple() // 15% 主食/面点
		}

		// 简单的去重检查
		if uniqueNames[r.Name] {
			continue
		}
		uniqueNames[r.Name] = true

		if err := repo.Create(r); err != nil {
			log.Printf("Insert Error: %v", err)
		} else {
			currentCount++
			if currentCount%100 == 0 {
				fmt.Printf("已生成: %d / %d\n", currentCount, targetCount)
			}
		}
	}
	log.Println("完成！")
}

// 1. 生成中式正餐 (注重炒、炖)
func genChineseDish() *models.Recipe {
	method := randElem(cnMethods)
	flavor := randElem(flavors)
	meat := randElem(meats)
	veg := randElem(veggies)

	name := fmt.Sprintf("%s%s%s", flavor, method, meat)
	if rand.Intn(2) == 0 {
		name = fmt.Sprintf("%s%s炒%s", flavor, meat, veg)
	}

	// 中餐营养模型：脂肪适中，碳水低(因为是纯菜)，蛋白高
	energy := 200 + rand.Float64()*400 // 200-600 kcal

	// 营养分配
	protein := energy * 0.4 / 4 // 高蛋白
	fat := energy * 0.45 / 9    // 脂肪略高(炒菜油)
	carb := energy * 0.15 / 4   // 低碳水(不含饭)

	return &models.Recipe{
		Name:         name,
		Description:  fmt.Sprintf("一道经典的%s，使用%s烹饪工艺，鲜香下饭。", flavor, method),
		ImageURL:     "https://placehold.co/600x400?text=Chinese+Food",
		MealType:     randMealType(true), // 午晚餐
		Energy:       truncate(energy),
		Protein:      truncate(protein),
		Carbohydrate: truncate(carb),
		Fat:          truncate(fat),
		Ingredients:  makeJSON([]string{meat, veg, "姜", "蒜", "食用油", "生抽"}),
		CookingTime:  10 + rand.Intn(30),
		Difficulty:   randDifficulty(),
	}
}

// 2. 生成西式/轻食 (注重煎、烤、沙拉)
func genWesternDish() *models.Recipe {
	method := randElem(enMethods)
	meat := randElem(meats)

	name := fmt.Sprintf("%s%s能量碗", method, meat)
	isSalad := false
	if rand.Intn(3) == 0 {
		name = fmt.Sprintf("%s%s考伯沙拉", method, meat)
		isSalad = true
	}

	energy := 300 + rand.Float64()*500 // 300-800 kcal

	// 西餐模型：如果是沙拉，脂肪低；如果是煎肉，脂肪高
	var pRatio, fRatio, cRatio float64
	if isSalad {
		pRatio, fRatio, cRatio = 0.3, 0.2, 0.5 // 高碳水(有玉米/土豆等配菜)
	} else {
		pRatio, fRatio, cRatio = 0.4, 0.5, 0.1 // 纯肉排
	}

	return &models.Recipe{
		Name:         name,
		Description:  "低卡饱腹，优质蛋白来源，健身人士首选。",
		ImageURL:     "https://placehold.co/600x400?text=Healthy+Bowl",
		MealType:     randMealType(true),
		Energy:       truncate(energy),
		Protein:      truncate(energy * pRatio / 4),
		Carbohydrate: truncate(energy * cRatio / 4),
		Fat:          truncate(energy * fRatio / 9),
		Ingredients:  makeJSON([]string{meat, "生菜", "圣女果", "橄榄油", "黑胡椒"}),
		CookingTime:  15 + rand.Intn(40),
		Difficulty:   randDifficulty(),
	}
}

// 3. 生成早餐
func genBreakfast() *models.Recipe {
	main := randElem(breakfast)
	method := "搭配"

	name := fmt.Sprintf("%s%s%s", main, method, randElem(breakfast))
	if rand.Intn(2) == 0 {
		name = fmt.Sprintf("活力%s早餐盘", main)
	}

	energy := 250 + rand.Float64()*350 // 250-600 kcal

	// 早餐：碳水较高
	return &models.Recipe{
		Name:         name,
		Description:  "唤醒清晨活力的营养早餐，开启元气满满的一天。",
		ImageURL:     "https://placehold.co/600x400?text=Breakfast",
		MealType:     models.MealTypeBreakfast,
		Energy:       truncate(energy),
		Protein:      truncate(energy * 0.2 / 4),
		Carbohydrate: truncate(energy * 0.55 / 4),
		Fat:          truncate(energy * 0.25 / 9),
		Ingredients:  makeJSON([]string{main, "牛奶", "水果"}),
		CookingTime:  5 + rand.Intn(15),
		Difficulty:   "简单",
	}
}

// 4. 生成主食 (面/饭/薯)
func genStaple() *models.Recipe {
	staple := randElem(staples)
	meat := randElem(meats)

	name := fmt.Sprintf("%s配%s", meat, staple)
	if strings.Contains(staple, "面") {
		name = fmt.Sprintf("%s炒%s", meat, staple)
	}

	energy := 400 + rand.Float64()*400

	// 主食：高碳水
	return &models.Recipe{
		Name:         name,
		Description:  "优质碳水与蛋白的结合，提供持久饱腹感。",
		ImageURL:     "https://placehold.co/600x400?text=Staple+Food",
		MealType:     randMealType(true),
		Energy:       truncate(energy),
		Protein:      truncate(energy * 0.2 / 4),
		Carbohydrate: truncate(energy * 0.6 / 4),
		Fat:          truncate(energy * 0.2 / 9),
		Ingredients:  makeJSON([]string{staple, meat, "青菜"}),
		CookingTime:  15 + rand.Intn(30),
		Difficulty:   "中等",
	}
}

// 5. 生成零食/加餐
func genSnack() *models.Recipe {
	fruit := randElem(fruits)
	nut := randElem(nuts)

	name := fmt.Sprintf("%s%s杯", fruit, nut)
	if rand.Intn(2) == 0 {
		name = fmt.Sprintf("低脂%s酸奶碗", fruit)
	}

	energy := 100 + rand.Float64()*200 // 100-300 kcal

	return &models.Recipe{
		Name:         name,
		Description:  "解馋不胖的健康加餐，补充微量元素。",
		ImageURL:     "https://placehold.co/600x400?text=Snack",
		MealType:     models.MealTypeSnack,
		Energy:       truncate(energy),
		Protein:      truncate(energy * 0.15 / 4),
		Carbohydrate: truncate(energy * 0.5 / 4),  // 水果碳水
		Fat:          truncate(energy * 0.35 / 9), // 坚果脂肪
		Ingredients:  makeJSON([]string{fruit, nut, "希腊酸奶"}),
		CookingTime:  0, // 即食
		Difficulty:   "简单",
	}
}

// 辅助：随机返回午餐或晚餐
func randMealType(mainMeal bool) models.MealType {
	if !mainMeal {
		return models.MealTypeSnack
	}
	if rand.Intn(2) == 0 {
		return models.MealTypeLunch
	}
	return models.MealTypeDinner
}

func randDifficulty() string {
	r := rand.Intn(10)
	if r < 6 {
		return "简单"
	} else if r < 9 {
		return "中等"
	}
	return "困难"
}
