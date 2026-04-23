package main

import (
	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type manualRecipeInput struct {
	Name                   string   `json:"name"`
	ImageURL               string   `json:"image_url"`
	MealType               string   `json:"meal_type"`
	AllowedMealTypes       []string `json:"allowed_meal_types"`
	Difficulty             string   `json:"difficulty"`
	CookingTime            int      `json:"cooking_time"`
	PortionWeightG         float64  `json:"portion_weight_g"`
	Energy                 float64  `json:"energy"`
	Protein                float64  `json:"protein"`
	Carbohydrate           float64  `json:"carbohydrate"`
	Fat                    float64  `json:"fat"`
	Ingredients            []string `json:"ingredients"`
	CookingSteps           []string `json:"cooking_steps"`
	TargetUsers            []string `json:"target_users"`
	ForbiddenUsers         []string `json:"forbidden_users"`
	IsWeightLossFriendly   *bool    `json:"is_weight_loss_friendly"`
	IsMuscleGainFriendly   *bool    `json:"is_muscle_gain_friendly"`
	IsSugarControlFriendly *bool    `json:"is_sugar_control_friendly"`
	IsGeneralFriendly      *bool    `json:"is_general_friendly"`
}

func main() {
	filePath := flag.String("file", "./pkg/scripts/import_manual_recipes/manual_recipes_template.json", "菜谱 JSON 文件路径")
	dryRun := flag.Bool("dry-run", false, "仅校验和预览，不写入数据库")
	upsert := flag.Bool("upsert", false, "名称相同则更新（默认跳过）")
	imagePrefix := flag.String("image-prefix", "", "为非 http/非绝对路径的 image_url 增加前缀，例如 http://127.0.0.1:3000")
	flag.Parse()

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	if err := dao.InitDatabase(config.AppConfig.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	inputs, err := readInputs(*filePath)
	if err != nil {
		log.Fatalf("读取 JSON 失败: %v", err)
	}
	if len(inputs) == 0 {
		log.Fatalf("JSON 中没有可导入菜谱")
	}

	created := 0
	updated := 0
	skipped := 0

	for index, item := range inputs {
		recipe, err := toRecipe(item, strings.TrimSpace(*imagePrefix))
		if err != nil {
			log.Printf("跳过第 %d 条: %v", index+1, err)
			skipped++
			continue
		}

		var existing models.Recipe
		err = dao.DB.Where("name = ?", recipe.Name).First(&existing).Error
		exists := err == nil
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("跳过 %s: 查询失败: %v", recipe.Name, err)
			skipped++
			continue
		}

		if *dryRun {
			action := "创建"
			if exists {
				if *upsert {
					action = "更新"
				} else {
					action = "跳过"
				}
			}
			log.Printf("[dry-run][%s] %s | %s | image=%s", action, recipe.Name, recipe.MealType, recipe.ImageURL)
			if exists && !*upsert {
				skipped++
				continue
			}
			if exists {
				updated++
			} else {
				created++
			}
			continue
		}

		if exists {
			if !*upsert {
				log.Printf("跳过重复菜名: %s", recipe.Name)
				skipped++
				continue
			}
			recipe.Model = existing.Model
			if err := dao.DB.Save(&recipe).Error; err != nil {
				log.Printf("更新失败 %s: %v", recipe.Name, err)
				skipped++
				continue
			}
			updated++
			continue
		}

		if err := dao.DB.Table("recipes").Create(&recipe).Error; err != nil {
			log.Printf("创建失败 %s: %v", recipe.Name, err)
			skipped++
			continue
		}
		created++
	}

	log.Printf("导入完成: 创建=%d, 更新=%d, 跳过=%d", created, updated, skipped)
}

func readInputs(path string) ([]manualRecipeInput, error) {
	cleanPath := filepath.Clean(path)
	bytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	var list []manualRecipeInput
	if err := json.Unmarshal(bytes, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func toRecipe(item manualRecipeInput, imagePrefix string) (models.Recipe, error) {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		return models.Recipe{}, fmt.Errorf("name 不能为空")
	}

	mealType, err := normalizeMealType(item.MealType)
	if err != nil {
		return models.Recipe{}, err
	}

	ingredients := cleanStringList(item.Ingredients)
	if len(ingredients) == 0 {
		return models.Recipe{}, fmt.Errorf("%s: ingredients 不能为空", name)
	}

	steps := cleanStringList(item.CookingSteps)
	if len(steps) == 0 {
		return models.Recipe{}, fmt.Errorf("%s: cooking_steps 不能为空", name)
	}

	if item.Energy <= 0 || item.Protein < 0 || item.Carbohydrate < 0 || item.Fat < 0 {
		return models.Recipe{}, fmt.Errorf("%s: 营养字段不合法", name)
	}

	portionWeight := item.PortionWeightG
	if portionWeight <= 0 {
		portionWeight = 100
	}

	imageURL := strings.TrimSpace(item.ImageURL)
	if imageURL != "" {
		if !isAbsoluteURL(imageURL) && !strings.HasPrefix(imageURL, "/") && imagePrefix != "" {
			imageURL = strings.TrimRight(imagePrefix, "/") + "/" + strings.TrimLeft(imageURL, "/")
		}
	}

	allowedMealTypes := cleanStringList(item.AllowedMealTypes)
	if len(allowedMealTypes) == 0 {
		allowedMealTypes = []string{string(mealType)}
	}

	recipe := models.Recipe{
		Name:                   name,
		ImageURL:               imageURL,
		MealType:               mealType,
		AllowedMealTypes:       allowedMealTypes,
		Difficulty:             strings.TrimSpace(item.Difficulty),
		CookingTime:            item.CookingTime,
		PortionWeightG:         portionWeight,
		Energy:                 item.Energy,
		Protein:                item.Protein,
		Carbohydrate:           item.Carbohydrate,
		Fat:                    item.Fat,
		Ingredients:            ingredients,
		CookingSteps:           steps,
		TargetUsers:            cleanStringList(item.TargetUsers),
		ForbiddenUsers:         cleanStringList(item.ForbiddenUsers),
		IsWeightLossFriendly:   defaultBool(item.IsWeightLossFriendly, false),
		IsMuscleGainFriendly:   defaultBool(item.IsMuscleGainFriendly, false),
		IsSugarControlFriendly: defaultBool(item.IsSugarControlFriendly, false),
		IsGeneralFriendly:      defaultBool(item.IsGeneralFriendly, true),
	}

	if recipe.Difficulty == "" {
		recipe.Difficulty = "简单"
	}
	if recipe.CookingTime < 0 {
		recipe.CookingTime = 0
	}

	return recipe, nil
}

func normalizeMealType(input string) (models.MealType, error) {
	value := strings.TrimSpace(input)
	switch value {
	case string(models.MealTypeBreakfast):
		return models.MealTypeBreakfast, nil
	case string(models.MealTypeLunch):
		return models.MealTypeLunch, nil
	case string(models.MealTypeDinner):
		return models.MealTypeDinner, nil
	case string(models.MealTypeSnack):
		return models.MealTypeSnack, nil
	default:
		return "", fmt.Errorf("meal_type 不合法: %s（可选：早餐/午餐/晚餐/加餐）", value)
	}
}

func cleanStringList(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func defaultBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func isAbsoluteURL(value string) bool {
	lowerValue := strings.ToLower(value)
	return strings.HasPrefix(lowerValue, "http://") || strings.HasPrefix(lowerValue, "https://")
}
