package main

import (
	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type recipeIndexItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}

type recipeDetail struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Category    string                  `json:"category"`
	Difficulty  string                  `json:"difficulty"`
	Tools       []string                `json:"tools"`
	Ingredients []projKitchenIngredient `json:"ingredients"`
	Steps       []string                `json:"steps"`
	Tips        string                  `json:"tips"`
}

type projKitchenIngredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type recipeCandidate struct {
	Recipe  models.Recipe
	Scores  map[models.MealType]int
	Primary models.MealType
}

func main() {
	apiBase := flag.String("api-base", "https://proj.kitchen/api", "Proj.Kitchen API base URL")
	dryRun := flag.Bool("dry-run", false, "Fetch and select without writing to database")
	timeoutSec := flag.Int("timeout-sec", 20, "HTTP timeout in seconds")
	flag.Parse()

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	if err := dao.InitDatabase(config.AppConfig.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	client := &http.Client{Timeout: time.Duration(*timeoutSec) * time.Second}
	apiRoot := strings.TrimRight(*apiBase, "/")

	indexItems, err := fetchRecipeIndex(client, apiRoot)
	if err != nil {
		log.Fatalf("获取菜谱索引失败: %v", err)
	}
	log.Printf("菜谱索引数量: %d", len(indexItems))

	existingNames, err := loadExistingRecipeNames()
	if err != nil {
		log.Fatalf("加载现有 recipes 失败: %v", err)
	}
	log.Printf("现有 recipes 去重基数: %d", len(existingNames))
	seenNames := make(map[string]struct{}, len(indexItems))
	imported := 0
	skipped := 0

	for _, item := range indexItems {
		detail, err := fetchRecipeDetail(client, apiRoot, item.ID)
		if err != nil {
			log.Printf("跳过 %s: 获取详情失败: %v", item.Name, err)
			skipped++
			continue
		}

		candidate, ok, reason := buildCandidate(detail)
		if !ok {
			log.Printf("跳过 %s: %s", detail.Name, reason)
			skipped++
			continue
		}

		nameKey := normalizeName(candidate.Recipe.Name)
		if _, exists := existingNames[nameKey]; exists {
			continue
		}
		if _, exists := seenNames[nameKey]; exists {
			continue
		}

		seenNames[nameKey] = struct{}{}

		if *dryRun {
			log.Printf("[dry-run] 导入 %s | %s | steps=%d | ingredients=%d",
				candidate.Recipe.Name,
				candidate.Recipe.MealType,
				len(candidate.Recipe.CookingSteps),
				len(candidate.Recipe.Ingredients),
			)
			imported++
			continue
		}

		if err := dao.DB.Table("recipes").Create(&candidate.Recipe).Error; err != nil {
			log.Printf("导入失败 %s: %v", candidate.Recipe.Name, err)
			skipped++
			continue
		}
		imported++
	}

	log.Printf("Proj.Kitchen 全量导入完成: 导入=%d, 跳过=%d", imported, skipped)
}

func fetchRecipeIndex(client *http.Client, apiBase string) ([]recipeIndexItem, error) {
	var items []recipeIndexItem
	if err := getJSON(client, apiBase+"/recipes", &items); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchRecipeDetail(client *http.Client, apiBase, id string) (*recipeDetail, error) {
	var detail recipeDetail
	if err := getJSON(client, apiBase+"/recipes/"+id, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

func getJSON(client *http.Client, url string, out interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("状态码 %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func loadExistingRecipeNames() (map[string]struct{}, error) {
	var rows []struct {
		Name string `gorm:"column:name"`
	}
	if err := dao.DB.Table("recipes").Select("name").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		key := normalizeName(row.Name)
		if key == "" {
			continue
		}
		result[key] = struct{}{}
	}
	return result, nil
}

func buildCandidate(detail *recipeDetail) (recipeCandidate, bool, string) {
	if detail == nil {
		return recipeCandidate{}, false, "详情为空"
	}

	name := strings.TrimSpace(detail.Name)
	if name == "" {
		return recipeCandidate{}, false, "菜名为空"
	}

	ingredients := buildIngredientList(detail.Ingredients)
	if len(ingredients) == 0 {
		ingredients = []string{"见原始菜谱"}
	}

	steps := cleanSteps(detail.Steps, detail.Tips)
	if len(steps) == 0 {
		steps = []string{"按原始菜谱制作"}
	}

	scores := scoreMealTypes(detail)
	primary := strongestMealType(scores)

	recipe := models.Recipe{
		Name:           name,
		ImageURL:       "",
		MealType:       primary,
		Difficulty:     normalizeDifficulty(detail.Difficulty),
		CookingTime:    estimateCookingTime(detail.Steps, detail.Difficulty),
		PortionWeightG: 0,
		Energy:         0,
		Protein:        0,
		Carbohydrate:   0,
		Fat:            0,
		Ingredients:    ingredients,
		CookingSteps:   steps,
		TargetUsers:    []string{},
		ForbiddenUsers: inferForbiddenUsers(ingredients),
		IsWeightLossFriendly:   false,
		IsMuscleGainFriendly:   false,
		IsSugarControlFriendly: false,
		IsGeneralFriendly:      true,
	}

	return recipeCandidate{
		Recipe:  recipe,
		Scores:  scores,
		Primary: primary,
	}, true, ""
}

func buildIngredientList(items []projKitchenIngredient) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		amount := strings.TrimSpace(item.Amount)
		if name == "" {
			continue
		}
		if amount == "" {
			result = append(result, name)
			continue
		}
		result = append(result, fmt.Sprintf("%s %s", name, amount))
	}
	return uniqueStrings(result)
}

func cleanSteps(steps []string, tips string) []string {
	result := make([]string, 0, len(steps)+1)
	for _, step := range steps {
		step = strings.TrimSpace(step)
		if step == "" {
			continue
		}
		result = append(result, step)
	}
	if tip := strings.TrimSpace(tips); tip != "" {
		result = append(result, "提示："+tip)
	}
	return result
}

func scoreMealTypes(detail *recipeDetail) map[models.MealType]int {
	content := normalizeText(detail.Name + " " + detail.Category + " " + strings.Join(detail.Tools, " "))
	for _, ingredient := range detail.Ingredients {
		content += " " + normalizeText(ingredient.Name) + " " + normalizeText(ingredient.Amount)
	}

	scores := map[models.MealType]int{
		models.MealTypeBreakfast: 8,
		models.MealTypeLunch:     30,
		models.MealTypeDinner:    24,
		models.MealTypeSnack:     6,
	}

	addScore := func(mealType models.MealType, score int, keywords ...string) {
		if containsAny(content, keywords...) {
			scores[mealType] += score
		}
	}

	addScore(models.MealTypeBreakfast, 90,
		"早餐", "早点", "早饭", "粥", "包子", "馒头", "花卷", "豆浆", "牛奶", "吐司",
		"面包", "三明治", "鸡蛋饼", "蛋饼", "燕麦", "煎蛋", "荷包蛋", "小米粥", "馄饨", "早餐饼",
	)
	addScore(models.MealTypeBreakfast, 35,
		"鸡蛋", "酸奶", "面条", "米线", "卷饼", "玉米", "紫薯", "南瓜", "红薯", "豆腐脑",
	)

	addScore(models.MealTypeSnack, 100,
		"加餐", "甜品", "点心", "小吃", "零食", "饮品", "饮料", "沙拉", "凉菜", "凉拌",
		"奶昔", "酸奶", "蛋糕", "冰粉", "布丁", "水果", "果汁", "奶茶", "甜汤", "糖水",
	)
	addScore(models.MealTypeSnack, 40,
		"饼干", "果酱", "冰淇淋", "果盘", "果杯", "沙冰", "小食",
	)

	addScore(models.MealTypeDinner, 80,
		"晚餐", "红烧", "炖", "煲", "焖", "卤", "锅", "火锅", "排骨", "牛腩", "肘子",
		"烧肉", "啤酒鸭", "咖喱", "烤鸡", "烧鸭", "鱼煲", "羊腩", "大盘鸡",
	)
	addScore(models.MealTypeDinner, 35,
		"牛肉", "羊肉", "鸭肉", "五花肉", "猪蹄", "鸡块", "砂锅", "汤", "炖饭",
	)

	addScore(models.MealTypeLunch, 65,
		"午餐", "家常菜", "下饭", "盖饭", "便当", "炒饭", "炒面", "宫保", "回锅", "麻婆",
		"豆腐", "鸡丁", "肉丝", "时蔬", "小炒", "番茄炒蛋", "青椒", "地三鲜", "凉皮", "凉粉",
	)
	addScore(models.MealTypeLunch, 20,
		"鸡胸", "虾仁", "炒蛋", "米饭", "面", "土豆", "番茄", "黄瓜", "木耳", "花菜",
	)

	if containsAny(content, "凉拌", "沙拉", "水果", "点心", "蛋糕") {
		scores[models.MealTypeDinner] -= 15
	}
	if containsAny(content, "红烧", "焖", "炖", "火锅", "排骨", "牛腩") {
		scores[models.MealTypeSnack] -= 20
		scores[models.MealTypeBreakfast] -= 12
	}
	if containsAny(content, "粥", "包子", "豆浆", "吐司", "燕麦") {
		scores[models.MealTypeDinner] -= 12
	}

	return scores
}

func strongestMealType(scores map[models.MealType]int) models.MealType {
	order := []models.MealType{
		models.MealTypeBreakfast,
		models.MealTypeLunch,
		models.MealTypeDinner,
		models.MealTypeSnack,
	}

	bestMeal := models.MealTypeLunch
	bestScore := -1 << 30
	for _, mealType := range order {
		score := scores[mealType]
		if score > bestScore {
			bestMeal = mealType
			bestScore = score
		}
	}
	return bestMeal
}


func normalizeDifficulty(raw string) string {
	switch strings.TrimSpace(raw) {
	case "简单", "中等", "困难":
		return strings.TrimSpace(raw)
	default:
		return "中等"
	}
}

func estimateCookingTime(steps []string, difficulty string) int {
	switch normalizeDifficulty(difficulty) {
	case "简单":
		return max(10, len(steps)*3)
	case "困难":
		return max(30, len(steps)*6)
	default:
		return max(20, len(steps)*4)
	}
}

func inferForbiddenUsers(ingredients []string) []string {
	content := normalizeText(strings.Join(ingredients, " "))
	tags := make([]string, 0, 6)

	add := func(tag string) {
		for _, existing := range tags {
			if existing == tag {
				return
			}
		}
		tags = append(tags, tag)
	}

	if containsAny(content, "虾", "蟹", "贝", "鱼", "海参", "牡蛎", "鱿鱼", "海鲜") {
		add(models.TagSeafood)
	}
	if containsAny(content, "牛肉", "牛腩", "肥牛") {
		add(models.TagBeef)
	}
	if containsAny(content, "猪肉", "五花肉", "排骨", "里脊", "猪蹄") {
		add(models.TagPork)
	}
	if containsAny(content, "鸡蛋", "鸭蛋", "鹌鹑蛋", "蛋液") {
		add(models.TagEgg)
	}
	if containsAny(content, "牛奶", "酸奶", "奶酪", "芝士", "黄油", "奶油") {
		add(models.TagDairy)
	}
	if containsAny(content, "面粉", "面包", "吐司", "面条", "饺子皮", "馄饨皮", "面片", "全麦") {
		add(models.TagGluten)
	}
	if containsAny(content, "花生", "腰果", "核桃", "杏仁", "坚果", "松子") {
		add(models.TagNut)
	}

	return tags
}

func normalizeName(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(name)), " "))
}

func normalizeText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "（", "(")
	value = strings.ReplaceAll(value, "）", ")")
	value = strings.ReplaceAll(value, "：", ":")
	return value
}

func uniqueStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func containsAny(value string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
