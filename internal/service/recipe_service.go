package service

import (
	"NutriPlan/internal/client"
	"NutriPlan/internal/config"
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RecipeService 食谱推荐服务接口
type RecipeService interface {
	// RecommendRecipes 为用户推荐每日食谱计划
	// count: 推荐方案数量（默认3-5套）
	RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error)

	// SaveDailyPlan 保存每日食谱计划
	SaveDailyPlan(plan *models.DailyRecipePlan) error

	// GetSelectedPlan 获取用户当前选中的计划
	GetSelectedPlan(userID uint) (*models.DailyRecipePlan, error)

	// SelectDailyPlan 选择每日食谱计划并保存
	SelectDailyPlan(userID uint, plan *models.DailyRecipePlan) error

	// GetRecipeDetail 获取食谱详情
	GetRecipeDetail(recipeID uint, userID uint) (*models.Recipe, bool, error)

	// AddFavorite 添加收藏
	AddFavorite(userID, recipeID uint) error

	// RemoveFavorite 取消收藏
	RemoveFavorite(userID, recipeID uint) error

	// GetUserFavorites 获取用户收藏列表
	GetUserFavorites(userID uint) ([]models.Recipe, error)

	// IsFavorite 检查收藏状态
	IsFavorite(userID, recipeID uint) (bool, error)

	// RecognizeIngredientsFromImage 从图片识别可用食材列表
	RecognizeIngredientsFromImage(ctx context.Context, imageReader io.Reader) ([]string, error)

	// RegenerateConstrainedMeal 基于食材约束重构单餐方案
	RegenerateConstrainedMeal(ctx context.Context, user *models.User, req ConstrainedMealRegenerateRequest) (*ConstrainedMealRegenerateResponse, error)

	// AdoptRegeneratedMeal 采纳重构结果并仅更新目标餐次
	AdoptRegeneratedMeal(userID uint, req ConstrainedMealAdoptRequest) (*models.DailyRecipePlan, error)
}

// RecipeServiceImpl 食谱推荐服务实现
type RecipeServiceImpl struct {
	recipeRepo   dao.RecipeRepository
	nutriService NutriService
	zhipuClient  client.ZhipuAIClient
	rng          *rand.Rand
}

type dietModeProfile struct {
	calorieFactor      float64
	proteinFactor      float64
	carbohydrateFactor float64
	fatFactor          float64
	strictFiltering    bool
	forbiddenKeywords  []string
}

type recipeTextConfig struct {
	StapleKeywords            []string            `json:"staple_keywords"`
	BreakfastPriorityKeywords map[string][]string `json:"breakfast_priority_keywords"`
	DietModeForbiddenKeywords map[string][]string `json:"diet_mode_forbidden_keywords"`
}

var (
	recipeTextConfigOnce sync.Once
	recipeTextConfigData recipeTextConfig
)

func getRecipeTextConfig() recipeTextConfig {
	recipeTextConfigOnce.Do(func() {
		raw := config.GetRecipeTextConfigRaw()
		if err := json.Unmarshal(raw, &recipeTextConfigData); err != nil {
			log.Printf("[recipe] parse recipe_text_config.json failed: %v", err)
			recipeTextConfigData = recipeTextConfig{}
		}

		if recipeTextConfigData.BreakfastPriorityKeywords == nil {
			recipeTextConfigData.BreakfastPriorityKeywords = map[string][]string{}
		}
		if recipeTextConfigData.DietModeForbiddenKeywords == nil {
			recipeTextConfigData.DietModeForbiddenKeywords = map[string][]string{}
		}
	})

	return recipeTextConfigData
}

// NewRecipeService 创建食谱推荐服务实例
func NewRecipeService(recipeRepo dao.RecipeRepository, nutriService NutriService, zhipuClient client.ZhipuAIClient) RecipeService {
	return &RecipeServiceImpl{
		recipeRepo:   recipeRepo,
		nutriService: nutriService,
		zhipuClient:  zhipuClient,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func normalizePortionWeight(weight float64) float64 {
	if weight <= 0 {
		return 100
	}
	return weight
}

func scaleRecipeNutritionByPortion(recipe models.Recipe) models.Recipe {
	scaled := recipe
	portionWeight := normalizePortionWeight(recipe.PortionWeightG)
	scaled.PortionWeightG = portionWeight
	factor := portionWeight / 100
	scaled.Energy = recipe.Energy * factor
	scaled.Protein = recipe.Protein * factor
	scaled.Carbohydrate = recipe.Carbohydrate * factor
	scaled.Fat = recipe.Fat * factor
	return scaled
}

func getDietModeProfile(mode models.DietMode) dietModeProfile {
	cfg := getRecipeTextConfig()
	getForbidden := func(key string) []string {
		return cfg.DietModeForbiddenKeywords[key]
	}

	switch mode {
	case models.DietModeLightAdjust:
		return dietModeProfile{
			calorieFactor:      0.94,
			proteinFactor:      1.0,
			carbohydrateFactor: 0.96,
			fatFactor:          0.88,
			strictFiltering:    false,
			forbiddenKeywords:  getForbidden(string(models.DietModeLightAdjust)),
		}
	case models.DietModeBland:
		return dietModeProfile{
			calorieFactor:      0.92,
			proteinFactor:      1.02,
			carbohydrateFactor: 0.95,
			fatFactor:          0.82,
			strictFiltering:    true,
			forbiddenKeywords:  getForbidden(string(models.DietModeBland)),
		}
	case models.DietModeHeavyAdjust:
		return dietModeProfile{
			calorieFactor:      0.88,
			proteinFactor:      1.05,
			carbohydrateFactor: 0.92,
			fatFactor:          0.75,
			strictFiltering:    true,
			forbiddenKeywords:  getForbidden(string(models.DietModeHeavyAdjust)),
		}
	default:
		return dietModeProfile{
			calorieFactor:      1,
			proteinFactor:      1,
			carbohydrateFactor: 1,
			fatFactor:          1,
			strictFiltering:    false,
			forbiddenKeywords:  nil,
		}
	}
}

func (s *RecipeServiceImpl) isManualModeActive(user *models.User, now time.Time) bool {
	if user == nil {
		return false
	}
	if user.DietModeSource != "manual" {
		return false
	}
	if user.DietMode == "" || user.DietMode == models.DietModeNormal {
		return false
	}
	if user.DietModeUntil == nil {
		return true
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	until := time.Date(user.DietModeUntil.Year(), user.DietModeUntil.Month(), user.DietModeUntil.Day(), 0, 0, 0, 0, user.DietModeUntil.Location())
	return !until.Before(today)
}

func (s *RecipeServiceImpl) resolveDietMode(user *models.User, targetCalorie, targetFat float64) models.DietMode {
	if user == nil {
		return models.DietModeNormal
	}

	now := time.Now()
	if s.isManualModeActive(user, now) {
		return user.DietMode
	}

	// 手动模式已过期，回落 normal。
	if user.DietModeSource == "manual" && user.DietModeUntil != nil {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		until := time.Date(user.DietModeUntil.Year(), user.DietModeUntil.Month(), user.DietModeUntil.Day(), 0, 0, 0, 0, user.DietModeUntil.Location())
		if until.Before(today) {
			_ = dao.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]any{
				"diet_mode":        models.DietModeNormal,
				"diet_mode_source": "auto",
				"diet_mode_reason": "",
				"diet_mode_until":  nil,
			}).Error
			user.DietMode = models.DietModeNormal
			user.DietModeSource = "auto"
			user.DietModeReason = ""
			user.DietModeUntil = nil
		}
	}

	if dao.DB == nil || targetCalorie <= 0 || targetFat <= 0 {
		return user.DietMode
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	start := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	end := start.Add(24 * time.Hour)

	type intakeSummary struct {
		TotalEnergy float64
		TotalFat    float64
		Count       int64
	}
	var summary intakeSummary
	err := dao.DB.Model(&models.DailyIntakeRecord{}).
		Select("COALESCE(SUM(calculated_energy),0) as total_energy, COALESCE(SUM(calculated_fat),0) as total_fat, COUNT(*) as count").
		Where("user_id = ? AND record_date >= ? AND record_date < ?", user.ID, start, end).
		Scan(&summary).Error
	if err != nil || summary.Count == 0 {
		return models.DietModeNormal
	}

	energyRate := summary.TotalEnergy / targetCalorie
	fatRate := summary.TotalFat / targetFat

	resolved := models.DietModeNormal
	modeReason := ""
	modeUntil := (*time.Time)(nil)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if energyRate >= 1.25 || fatRate >= 1.35 {
		resolved = models.DietModeHeavyAdjust
		modeReason = "昨日重度偏离，自动重调"
		until := today.AddDate(0, 0, 2)
		modeUntil = &until
	} else if energyRate >= 1.10 || fatRate >= 1.15 {
		resolved = models.DietModeLightAdjust
		modeReason = "昨日轻度偏离，自动轻调"
		until := today.AddDate(0, 0, 1)
		modeUntil = &until
	}

	if user.DietMode != resolved || user.DietModeSource != "auto" {
		_ = dao.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]any{
			"diet_mode":        resolved,
			"diet_mode_source": "auto",
			"diet_mode_reason": modeReason,
			"diet_mode_until":  modeUntil,
		}).Error
	}

	user.DietMode = resolved
	user.DietModeSource = "auto"
	user.DietModeReason = modeReason
	user.DietModeUntil = modeUntil

	return resolved
}

func recipeSupportsMealType(recipe models.Recipe, mealType models.MealType) bool {
	if len(recipe.AllowedMealTypes) > 0 {
		target := strings.TrimSpace(string(mealType))
		for _, allowedMealType := range recipe.AllowedMealTypes {
			if strings.TrimSpace(allowedMealType) == target {
				return true
			}
		}
		return false
	}

	return recipe.MealType == mealType
}

func (s *RecipeServiceImpl) filterRecipesByAllowedMealType(recipes []models.Recipe, mealType models.MealType) []models.Recipe {
	if len(recipes) == 0 {
		return recipes
	}

	filtered := make([]models.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		if recipeSupportsMealType(recipe, mealType) {
			filtered = append(filtered, recipe)
		}
	}

	return filtered
}

func (s *RecipeServiceImpl) filterRecipesByDietMode(recipes []models.Recipe, profile dietModeProfile) []models.Recipe {
	if len(recipes) == 0 || len(profile.forbiddenKeywords) == 0 {
		return recipes
	}

	containsForbidden := func(recipe models.Recipe) bool {
		content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))
		for _, keyword := range profile.forbiddenKeywords {
			if keyword == "" {
				continue
			}
			if strings.Contains(content, strings.ToLower(keyword)) {
				return true
			}
		}
		return false
	}

	filtered := make([]models.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		if containsForbidden(recipe) {
			continue
		}
		if profile.strictFiltering && recipe.Fat > 25 {
			continue
		}
		filtered = append(filtered, recipe)
	}

	if len(filtered) < 3 {
		return recipes
	}

	return filtered
}

// 目前的推荐算法逻辑如下：

// 1. 目标设定
// 根据用户的档案（年龄、体重、目标等）计算出每日目标热量 (TDEE) 和 三大营养素目标（蛋白质、碳水、脂肪）。

// 2. 食谱初筛
// 	禁忌过滤：剔除过敏食材。
// 	历史去重：剔除最近 7 天使用过的食谱。
// 	候选池选取：
// 	使用多路召回技术：不再仅取热量最接近的菜品。我们同时调用多个独立查询拉取候选食谱并汇总去重：
//      - 适中热量的传统召回（沿用之前的逻辑）
//      - 极低热量与极高热量召回（用于高低搭配）
//      - 按用户偏好特征标签拉取相关的候选菜（如喜欢某种食材）
//      - 随机冷启动探索（捞一些不经常出现的盲盒菜谱）
// 3. 生成组合
// 4. 评分逻辑
// 	计算每种组合的总热量和宏量营养素。
// 	计算其与目标值的偏差
// 	真实评分：根据偏差计算 0-100 的分数。
// 	无同情分：您刚才删除了 <60 分强行补到 60 分的逻辑，现在返回的是真实反映差距的分数。
// 5. 排序与优选
// 	将所有方案按分数从高到低排序。
// 	加权随机：取前 100 名，进行加权随机抽取（分数越高的越容易被选中，但不是绝对选第一名）。
// 	硬门槛拦截（您刚才修改的）：
// 	代码行 if plan.MatchScore < 80 { break }
// 	如果不满 80 分，直接丢弃，不再推荐。
// 	多样性检查：如果新选出来的方案和已选方案在“午餐”和“晚餐”上由于过度重合，会被跳过。

// multiChannelRecall 执行单餐次的“多路召回”，并发拉取各渠道食谱并去重合并。
func (s *RecipeServiceImpl) multiChannelRecall(mealType models.MealType, targetCal float64, goal models.HealthGoal, userTags []string, forbidden []string) []models.Recipe {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 用于存储召回来的所有去重结果
	recipeMap := make(map[uint]models.Recipe)

	// 通用的“添加合并食谱”方法，内含加锁保护和重复判定
	addRecipes := func(recipes []models.Recipe) {
		mu.Lock()
		defer mu.Unlock()
		for _, rawRecipe := range recipes {
			if _, exists := recipeMap[rawRecipe.ID]; !exists {
				recipeMap[rawRecipe.ID] = rawRecipe
			}
		}
	}

	// 1. 传统适中热量路：沿用老接口拿100条基础菜，再选适中的。
	wg.Add(1)
	go func() {
		defer wg.Done()
		bases, _ := s.recipeRepo.FindByMealType(mealType, goal, append(forbidden, "")) // (简单取)
		// 初选排序拿出最接近目标的 30 道
		filtered := s.getTopCandidates(bases, targetCal, targetCal*0.3/4.0, targetCal*0.5/4.0, targetCal*0.2/9.0, 30, "") // 这是本服务自带现成的函数
		addRecipes(filtered)
	}()

	// 2. 极低热量路：拉取垫底的 15 道
	wg.Add(1)
	go func() {
		defer wg.Done()
		lowCals, _ := s.recipeRepo.FindLowCalorieRecipes(mealType, goal, 15)
		addRecipes(lowCals)
	}()

	// 3. 极高热量路：拉取封顶的 15 道
	wg.Add(1)
	go func() {
		defer wg.Done()
		highCals, _ := s.recipeRepo.FindHighCalorieRecipes(mealType, goal, 15)
		addRecipes(highCals)
	}()

	// 4. 偏好标签路：根据 userTags（如牛肉, 海鲜等）模糊拉取
	if len(userTags) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			liked, _ := s.recipeRepo.FindByIngredientTags(mealType, goal, userTags, 20)
			addRecipes(liked)
		}()
	}

	// 5. 随机探索路：增加库里的长尾曝光 10 道
	wg.Add(1)
	go func() {
		defer wg.Done()
		randoms, _ := s.recipeRepo.FindRandomRecipes(mealType, goal, 10)
		addRecipes(randoms)
	}()

	wg.Wait()

	// 把 Map 压平成 Slice
	result := make([]models.Recipe, 0, len(recipeMap))
	for _, v := range recipeMap {
		result = append(result, v)
	}

	return s.filterRecipesByAllowedMealType(result, mealType)
}

// RecommendRecipes 为用户推荐每日食谱计划
func (s *RecipeServiceImpl) recommendRecipesLegacy(user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	// 计算用户目标营养需求
	targetCalorie := s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
	targetProtein, targetCarb, targetFat := s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)
	dietMode := s.resolveDietMode(user, targetCalorie, targetFat)
	modeProfile := getDietModeProfile(dietMode)
	targetCalorie *= modeProfile.calorieFactor
	targetProtein *= modeProfile.proteinFactor
	targetCarb *= modeProfile.carbohydrateFactor
	targetFat *= modeProfile.fatFactor

	// 解析用户偏好与禁忌（过敏 / 健康状况）
	userTags, forbiddenIngredients := s.parseUserPreferences(user)

	// ---- 第一阶段：多路召回 (Multi-channel Recall) ----

	// 对于每一顿饭分别并发执行多路召回（这里不再依赖老式的单条语句硬捞）
	var wgRecall sync.WaitGroup
	var breakfastRaw, lunchRaw, dinnerRaw, snackRaw []models.Recipe

	wgRecall.Add(4)
	go func() {
		defer wgRecall.Done()
		breakfastRaw = s.multiChannelRecall(models.MealTypeBreakfast, targetCalorie*0.25, user.HealthGoal, userTags, forbiddenIngredients)
	}()
	go func() {
		defer wgRecall.Done()
		lunchRaw = s.multiChannelRecall(models.MealTypeLunch, targetCalorie*0.35, user.HealthGoal, userTags, forbiddenIngredients)
	}()
	go func() {
		defer wgRecall.Done()
		dinnerRaw = s.multiChannelRecall(models.MealTypeDinner, targetCalorie*0.30, user.HealthGoal, userTags, forbiddenIngredients)
	}()
	go func() {
		defer wgRecall.Done()
		// 加餐这里假设占目标极小比例或者是固定值，随便传一个较小热量目标
		snackRaw = s.multiChannelRecall(models.MealTypeSnack, targetCalorie*0.10, user.HealthGoal, userTags, forbiddenIngredients)
	}()
	wgRecall.Wait()

	// -- 回落：针对各路的返回，再经过一次硬性的 Filter 门槛确认 --
	// 比如：过敏原在 SQL OR 里可能被意外拉起，安全起见我们要确保全通过。以及用户历史过滤。

	// 过滤黑名单
	breakfastRecipes := s.filterRecipesByUserPreferences(breakfastRaw, userTags)
	lunchRecipes := s.filterRecipesByUserPreferences(lunchRaw, userTags)
	dinnerRecipes := s.filterRecipesByUserPreferences(dinnerRaw, userTags)
	snackRecipes := s.filterRecipesByUserPreferences(snackRaw, userTags)

	log.Printf("Recall sizes -> Breakfast: %d, Lunch: %d, Dinner: %d, Snack: %d", len(breakfastRaw), len(lunchRaw), len(dinnerRaw), len(snackRaw))
	log.Printf("After pref filter -> Breakfast: %d, Lunch: %d, Dinner: %d, Snack: %d", len(breakfastRecipes), len(lunchRecipes), len(dinnerRecipes), len(snackRecipes))

	// 过滤最近5天吃过的食谱
	prefPools := mealRecipePools{
		breakfast: breakfastRecipes,
		lunch:     lunchRecipes,
		dinner:    dinnerRecipes,
		snack:     snackRecipes,
	}

	recentPools := prefPools
	recentRecipeIDs, err := s.recipeRepo.FindRecentPlanRecipeIDs(user.ID, 5)
	if err == nil && len(recentRecipeIDs) > 0 {
		recentPools = mealRecipePools{
			breakfast: s.filterRecentRecipes(prefPools.breakfast, recentRecipeIDs),
			lunch:     s.filterRecentRecipes(prefPools.lunch, recentRecipeIDs),
			dinner:    s.filterRecentRecipes(prefPools.dinner, recentRecipeIDs),
			snack:     s.filterRecentRecipes(prefPools.snack, recentRecipeIDs),
		}
	}

	dietPools := mealRecipePools{
		breakfast: s.filterRecipesByDietMode(recentPools.breakfast, modeProfile),
		lunch:     s.filterRecipesByDietMode(recentPools.lunch, modeProfile),
		dinner:    s.filterRecipesByDietMode(recentPools.dinner, modeProfile),
		snack:     s.filterRecipesByDietMode(recentPools.snack, modeProfile),
	}

	strictPools := mealRecipePools{
		breakfast: s.filterRecipesByHealthGoal(dietPools.breakfast, user.HealthGoal),
		lunch:     s.filterRecipesByHealthGoal(dietPools.lunch, user.HealthGoal),
		dinner:    s.filterRecipesByHealthGoal(dietPools.dinner, user.HealthGoal),
		snack:     s.filterRecipesByHealthGoal(dietPools.snack, user.HealthGoal),
	}
	_ = strictPools

	// ---- 第二阶段：预估排序与 全局组合阶段 ----

	// 生成每日食谱计划
	targetNutrition := NutritionTarget{
		Energy:       targetCalorie,
		Protein:      targetProtein,
		Carbohydrate: targetCarb,
		Fat:          targetFat,
	}

	rankedPlans := s.generateRankedPlans(
		user.ID,
		breakfastRecipes, lunchRecipes, dinnerRecipes, snackRecipes,
		targetNutrition,
		user.HealthGoal,
	)
	log.Printf("generateRankedPlans returned %d plans", len(rankedPlans))
	if len(rankedPlans) > 0 {
		log.Printf("Highest score: %f", rankedPlans[0].MatchScore)
	}

	// 从计划里随机抽取,确保每次刷新结果不一致
	selectedPlans := s.shuffleByWeights(rankedPlans)

	// 多样性过滤与结果截取
	finalPlans := make([]*models.DailyRecipePlan, 0)
	minFinalScore := recommendationScoreFloor(len(breakfastRecipes), len(lunchRecipes), len(dinnerRecipes), len(snackRecipes))

	for _, plan := range selectedPlans {
		if plan.MatchScore < minFinalScore {
			break // 因为已经排好序了，后面更低，直接断开
		}

		// 多样性检查 (核心)
		// 确保新加入的 plan 和已经选中的 plan 不吃重复的主菜
		if !s.checkPlanDiversity(finalPlans, plan) {
			continue
		}

		// 不再入库,只返回推荐结果
		finalPlans = append(finalPlans, plan)

		if len(finalPlans) >= count {
			break
		}
	}

	if len(finalPlans) == 0 {
		for _, plan := range selectedPlans {
			if !s.checkPlanDiversity(finalPlans, plan) {
				continue
			}
			finalPlans = append(finalPlans, plan)
			if len(finalPlans) >= count {
				break
			}
		}
	}

	return finalPlans, nil
}

// NutritionTarget 营养目标
func (s *RecipeServiceImpl) RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	targetCalorie := s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
	targetProtein, targetCarb, targetFat := s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)
	dietMode := s.resolveDietMode(user, targetCalorie, targetFat)
	modeProfile := getDietModeProfile(dietMode)
	targetCalorie *= modeProfile.calorieFactor
	targetProtein *= modeProfile.proteinFactor
	targetCarb *= modeProfile.carbohydrateFactor
	targetFat *= modeProfile.fatFactor

	userTags, forbiddenIngredients := s.parseUserPreferences(user)

	var wg sync.WaitGroup
	var breakfastRaw, lunchRaw, dinnerRaw, snackRaw []models.Recipe

	loadMeal := func(mealType models.MealType, dest *[]models.Recipe) {
		defer wg.Done()
		recipes, err := s.recipeRepo.FindByMealType(mealType, user.HealthGoal, forbiddenIngredients)
		if err != nil {
			log.Printf("load recipes by meal type failed: meal=%s err=%v", mealType, err)
			return
		}
		*dest = s.filterRecipesByAllowedMealType(recipes, mealType)
	}

	wg.Add(4)
	go loadMeal(models.MealTypeBreakfast, &breakfastRaw)
	go loadMeal(models.MealTypeLunch, &lunchRaw)
	go loadMeal(models.MealTypeDinner, &dinnerRaw)
	go loadMeal(models.MealTypeSnack, &snackRaw)
	wg.Wait()

	breakfastRecipes := s.filterRecipesByUserPreferences(breakfastRaw, userTags)
	lunchRecipes := s.filterRecipesByUserPreferences(lunchRaw, userTags)
	dinnerRecipes := s.filterRecipesByUserPreferences(dinnerRaw, userTags)
	snackRecipes := s.filterRecipesByUserPreferences(snackRaw, userTags)

	log.Printf("Base pool sizes -> Breakfast: %d, Lunch: %d, Dinner: %d, Snack: %d", len(breakfastRaw), len(lunchRaw), len(dinnerRaw), len(snackRaw))
	log.Printf("After pref filter -> Breakfast: %d, Lunch: %d, Dinner: %d, Snack: %d", len(breakfastRecipes), len(lunchRecipes), len(dinnerRecipes), len(snackRecipes))

	prefPools := mealRecipePools{
		breakfast: breakfastRecipes,
		lunch:     lunchRecipes,
		dinner:    dinnerRecipes,
		snack:     snackRecipes,
	}

	recentPools := prefPools
	recentRecipeIDs, err := s.recipeRepo.FindRecentPlanRecipeIDs(user.ID, 5)
	if err == nil && len(recentRecipeIDs) > 0 {
		recentPools = mealRecipePools{
			breakfast: s.filterRecentRecipes(prefPools.breakfast, recentRecipeIDs),
			lunch:     s.filterRecentRecipes(prefPools.lunch, recentRecipeIDs),
			dinner:    s.filterRecentRecipes(prefPools.dinner, recentRecipeIDs),
			snack:     s.filterRecentRecipes(prefPools.snack, recentRecipeIDs),
		}
	}

	dietPools := mealRecipePools{
		breakfast: s.filterRecipesByDietMode(recentPools.breakfast, modeProfile),
		lunch:     s.filterRecipesByDietMode(recentPools.lunch, modeProfile),
		dinner:    s.filterRecipesByDietMode(recentPools.dinner, modeProfile),
		snack:     s.filterRecipesByDietMode(recentPools.snack, modeProfile),
	}

	strictPools := mealRecipePools{
		breakfast: s.filterRecipesByHealthGoal(dietPools.breakfast, user.HealthGoal),
		lunch:     s.filterRecipesByHealthGoal(dietPools.lunch, user.HealthGoal),
		dinner:    s.filterRecipesByHealthGoal(dietPools.dinner, user.HealthGoal),
		snack:     s.filterRecipesByHealthGoal(dietPools.snack, user.HealthGoal),
	}

	targetNutrition := NutritionTarget{
		Energy:       targetCalorie,
		Protein:      targetProtein,
		Carbohydrate: targetCarb,
		Fat:          targetFat,
	}

	stages := []recommendationStage{
		{label: "strict", pools: strictPools},
		{label: "relax_health_goal", pools: dietPools},
	}
	if user.HealthGoal == models.GoalWeightLoss {
		stages = append(stages,
			recommendationStage{label: "relax_diet_mode", pools: recentPools},
			recommendationStage{label: "relax_recent_history", pools: prefPools},
		)
	} else {
		stages = append(stages, recommendationStage{label: "relax_diet_mode", pools: prefPools})
	}

	finalPlans, stage := s.buildRecommendedPlansWithFallback(
		user.ID,
		stages,
		targetNutrition,
		user.HealthGoal,
		count,
		s.checkDiversity,
	)
	if stage != "" && stage != "strict" {
		log.Printf("recommendation fallback activated: goal=%s stage=%s", user.HealthGoal, stage)
	}

	return finalPlans, nil
}

type mealRecipePools struct {
	breakfast []models.Recipe
	lunch     []models.Recipe
	dinner    []models.Recipe
	snack     []models.Recipe
}

type recommendationStage struct {
	label string
	pools mealRecipePools
}

type planDiversityChecker func(existing []*models.DailyRecipePlan, candidate *models.DailyRecipePlan) bool

func (s *RecipeServiceImpl) selectRecommendedPlans(selectedPlans []*models.DailyRecipePlan, count int, minFinalScore float64, allow planDiversityChecker) []*models.DailyRecipePlan {
	if len(selectedPlans) == 0 || count <= 0 {
		return nil
	}

	finalPlans := make([]*models.DailyRecipePlan, 0, count)
	for _, plan := range selectedPlans {
		if plan.MatchScore < minFinalScore {
			break
		}
		if !allow(finalPlans, plan) {
			continue
		}

		finalPlans = append(finalPlans, plan)
		if len(finalPlans) >= count {
			return finalPlans
		}
	}

	if len(finalPlans) > 0 {
		return finalPlans
	}

	for _, plan := range selectedPlans {
		if !allow(finalPlans, plan) {
			continue
		}
		finalPlans = append(finalPlans, plan)
		if len(finalPlans) >= count {
			break
		}
	}

	return finalPlans
}

func (s *RecipeServiceImpl) buildRecommendedPlans(userID uint, pools mealRecipePools, target NutritionTarget, goal models.HealthGoal, count int, allow planDiversityChecker) []*models.DailyRecipePlan {
	rankedPlans := s.generateRankedPlans(
		userID,
		pools.breakfast, pools.lunch, pools.dinner, pools.snack,
		target,
		goal,
	)
	log.Printf("generateRankedPlans returned %d plans", len(rankedPlans))
	if len(rankedPlans) > 0 {
		log.Printf("Highest score: %f", rankedPlans[0].MatchScore)
	}

	selectedPlans := s.shuffleByWeights(rankedPlans)
	minFinalScore := recommendationScoreFloor(len(pools.breakfast), len(pools.lunch), len(pools.dinner), len(pools.snack))
	return s.selectRecommendedPlans(selectedPlans, count, minFinalScore, allow)
}

func (s *RecipeServiceImpl) buildRecommendedPlansWithFallback(userID uint, stages []recommendationStage, target NutritionTarget, goal models.HealthGoal, count int, allow planDiversityChecker) ([]*models.DailyRecipePlan, string) {
	for _, stage := range stages {
		plans := s.buildRecommendedPlans(userID, stage.pools, target, goal, count, allow)
		if len(plans) > 0 {
			return plans, stage.label
		}
	}
	return nil, ""
}

type NutritionTarget struct {
	Energy       float64 // 热量
	Protein      float64 // 蛋白质
	Carbohydrate float64 // 碳水化合物
	Fat          float64 // 脂肪
}

// mealCombo 表示“单餐组合”（可包含 1~N 道菜）。
// Primary 用于兼容旧字段（BreakfastRecipeID/LunchRecipeID...），
// Items 用于新 DTO 返回完整组合。
type mealCombo struct {
	Items          []models.Recipe
	Primary        models.Recipe
	Energy         float64
	Protein        float64
	Carbohydrate   float64
	Fat            float64
	BreakfastBoost float64
	ItemIDSet      map[uint]struct{}
	IngredientSet  map[string]struct{}
}

func buildIngredientSet(items []models.Recipe) map[string]struct{} {
	set := make(map[string]struct{})
	for _, recipe := range items {
		for _, ingredient := range recipe.Ingredients {
			normalized := strings.ToLower(strings.TrimSpace(ingredient))
			if normalized != "" {
				set[normalized] = struct{}{}
			}
		}
	}
	return set
}

func buildItemIDSet(items []models.Recipe) map[uint]struct{} {
	set := make(map[uint]struct{})
	for _, item := range items {
		if item.ID > 0 {
			set[item.ID] = struct{}{}
		}
	}
	return set
}

func buildComboKey(items []models.Recipe) (string, bool) {
	if len(items) == 0 {
		return "", false
	}

	ids := make([]uint, 0, len(items))
	for _, item := range items {
		if item.ID == 0 {
			return "", false
		}
		ids = append(ids, item.ID)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	var builder strings.Builder
	for index, id := range ids {
		if index > 0 {
			builder.WriteByte('-')
		}
		builder.WriteString(strconv.FormatUint(uint64(id), 10))
	}

	return builder.String(), true
}

func isStapleRecipe(recipe models.Recipe) bool {
	cfg := getRecipeTextConfig()
	content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))
	for _, keyword := range cfg.StapleKeywords {
		if strings.Contains(content, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func isMainStapleRecipe(recipe models.Recipe) bool {
	content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))
	keywords := []string{
		"炒饭", "蛋包饭", "盖饭", "烩饭", "焗饭", "拌饭", "咖喱饭", "丼", "risotto",
		"炒面", "拌面", "意面", "炒粉", "炒米粉", "河粉", "米线", "拉面", "刀削面", "焖面", "面条",
		"炒年糕", "馄饨", "饺子",
	}
	return containsKeyword(content, keywords)
}

func countStapleRecipes(items []models.Recipe) int {
	count := 0
	for _, item := range items {
		if isStapleRecipe(item) {
			count++
		}
	}
	return count
}

func countMainStapleRecipes(items []models.Recipe) int {
	count := 0
	for _, item := range items {
		if isMainStapleRecipe(item) {
			count++
		}
	}
	return count
}

func mealHasStaple(items []models.Recipe) bool {
	for _, item := range items {
		if isStapleRecipe(item) {
			return true
		}
	}
	return false
}

func buildMealCombo(items []models.Recipe, boost float64) (mealCombo, string, bool) {
	if len(items) == 0 {
		return mealCombo{}, "", false
	}

	key, ok := buildComboKey(items)
	if !ok {
		return mealCombo{}, "", false
	}

	totalE, totalP, totalC, totalF := 0.0, 0.0, 0.0, 0.0
	scaledItems := make([]models.Recipe, 0, len(items))
	for _, item := range items {
		scaled := scaleRecipeNutritionByPortion(item)
		totalE += scaled.Energy
		totalP += scaled.Protein
		totalC += scaled.Carbohydrate
		totalF += scaled.Fat
		scaledItems = append(scaledItems, scaled)
	}

	combo := mealCombo{
		Items:          scaledItems,
		Primary:        scaledItems[0],
		Energy:         totalE,
		Protein:        totalP,
		Carbohydrate:   totalC,
		Fat:            totalF,
		BreakfastBoost: boost,
		ItemIDSet:      buildItemIDSet(items),
		IngredientSet:  buildIngredientSet(items),
	}

	return combo, key, true
}

func comboContent(combo mealCombo) string {
	parts := make([]string, 0, len(combo.Items)*2)
	for _, item := range combo.Items {
		parts = append(parts, item.Name)
		parts = append(parts, strings.Join(item.Ingredients, " "))
	}
	return strings.ToLower(strings.Join(parts, " "))
}

func containsKeyword(content string, keywords []string) bool {
	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}
		if strings.Contains(content, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func (s *RecipeServiceImpl) filterCombosByMealProfile(mealType models.MealType, combos []mealCombo, targetE, targetP, targetC, targetF float64, goal models.HealthGoal) []mealCombo {
	if len(combos) == 0 {
		return combos
	}

	breakfastKeywords := []string{
		"早餐", "早点", "早饭", "粥", "燕麦", "鸡蛋", "酸奶", "豆浆", "牛奶",
		"吐司", "面包", "三明治", "包子", "馒头", "花卷", "饼",
	}
	breakfastHeavyKeywords := []string{
		"火锅", "鸡翅", "烤鱼", "羊肉", "羊排", "羊汤", "牛腩", "肘子", "排骨",
		"肥肠", "啤酒鸭", "烧烤", "炸鸡", "空气炸锅", "肉饼", "腊肠",
	}
	breakfastMainMealKeywords := []string{
		"炒饭", "蛋包饭", "盖饭", "烩饭", "焗饭", "拌饭", "咖喱饭",
		"炒面", "拌面", "意面", "炒粉", "炒米粉", "焖面",
	}
	snackBadKeywords := []string{
		"奶茶", "莫吉托", "mojito", "长岛冰茶", "鸡尾酒", "特调", "果茶", "饮料",
		"糖水", "杨枝甘露", "冰淇淋", "提拉米苏", "雪花酥", "奶冻", "布丁",
		"火锅", "鸡翅", "烤鱼", "排骨", "牛腩", "肘子", "盖饭", "炒饭", "面", "粉",
	}
	dinnerHeavyKeywords := []string{
		"鸡尾酒", "莫吉托", "奶茶", "提拉米苏", "雪花酥", "糖水",
	}

	filtered := make([]mealCombo, 0, len(combos))
	for _, combo := range combos {
		content := comboContent(combo)
		switch mealType {
		case models.MealTypeBreakfast:
			maxBreakfastEnergy := math.Min(targetE*1.25, 520)
			if combo.Energy > maxBreakfastEnergy {
				continue
			}
			if combo.Fat > math.Max(targetF*1.6, 18) {
				continue
			}
			if containsKeyword(content, breakfastHeavyKeywords) {
				continue
			}
			if containsKeyword(content, breakfastMainMealKeywords) {
				continue
			}
			if !containsKeyword(content, breakfastKeywords) {
				continue
			}
		case models.MealTypeSnack:
			maxSnackEnergy := math.Min(targetE*1.35, 220)
			if combo.Energy > maxSnackEnergy {
				continue
			}
			if combo.Fat > math.Max(targetF*1.8, 12) {
				continue
			}
			if containsKeyword(content, snackBadKeywords) {
				continue
			}
			if goal == models.GoalWeightLoss && combo.Carbohydrate > 20 {
				continue
			}
		case models.MealTypeLunch, models.MealTypeDinner:
			maxMealEnergy := targetE * 1.25
			if combo.Energy > maxMealEnergy {
				continue
			}
			if mealType == models.MealTypeDinner && containsKeyword(content, dinnerHeavyKeywords) {
				continue
			}
		}
		filtered = append(filtered, combo)
	}

	if len(filtered) == 0 {
		return combos
	}
	return filtered
}

// rankMealCombos 按单餐营养匹配度对组合进行评分并截断 Top-N。
// BreakfastBoost 仅在早餐场景生效，用于稳定提升“酸奶+鸡蛋”组合优先级。
func (s *RecipeServiceImpl) rankMealCombos(combos []mealCombo, targetE, targetP, targetC, targetF float64, goal models.HealthGoal, limit int) []mealCombo {
	if len(combos) == 0 {
		return nil
	}

	type rankedCombo struct {
		combo mealCombo
		score float64
	}
	ranked := make([]rankedCombo, 0, len(combos))
	for _, combo := range combos {
		score := s.calculateMatchScore(
			combo.Energy, combo.Protein, combo.Carbohydrate, combo.Fat,
			targetE, targetP, targetC, targetF,
			goal,
		) + combo.BreakfastBoost
		ranked = append(ranked, rankedCombo{combo: combo, score: score})
	}

	sort.Slice(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	actualLimit := min(len(ranked), limit)
	result := make([]mealCombo, 0, actualLimit)
	for i := range actualLimit {
		result = append(result, ranked[i].combo)
	}

	return result
}

// buildMealCombos 通用餐次组合生成器。
// candidateK: 先从候选池选出的菜品数量。
// comboLimit: 最终返回的组合数量上限。
// minItems: 单个组合最少包含的菜品数（例如午/晚餐=2）。
// maxItems: 单个组合最多包含的菜品数（例如午/晚餐=3，加餐=2）。
func (s *RecipeServiceImpl) buildMealCombos(recipes []models.Recipe, targetE, targetP, targetC, targetF float64, goal models.HealthGoal, candidateK, comboLimit, minItems, maxItems int, requireStaple bool) []mealCombo {
	if len(recipes) == 0 {
		return nil
	}

	candidates := s.selectDiverseCandidates(recipes, targetE, targetP, targetC, targetF, candidateK, goal)
	if len(candidates) == 0 {
		return nil
	}

	if requireStaple {
		hasStapleCandidate := false
		for _, item := range candidates {
			if isStapleRecipe(item) {
				hasStapleCandidate = true
				break
			}
		}

		if !hasStapleCandidate {
			staplePool := make([]models.Recipe, 0)
			for _, recipe := range recipes {
				if isStapleRecipe(recipe) {
					staplePool = append(staplePool, recipe)
				}
			}

			sort.Slice(staplePool, func(i, j int) bool {
				left := s.scoreSingleRecipe(staplePool[i], targetE, targetP, targetC, targetF, goal)
				right := s.scoreSingleRecipe(staplePool[j], targetE, targetP, targetC, targetF, goal)
				return left > right
			})

			existing := make(map[uint]struct{}, len(candidates))
			for _, item := range candidates {
				existing[item.ID] = struct{}{}
			}

			for i := 0; i < len(staplePool) && i < 4; i++ {
				item := staplePool[i]
				if _, ok := existing[item.ID]; ok {
					continue
				}
				candidates = append(candidates, item)
			}
		}
	}

	stapleAvailable := false
	if requireStaple {
		for _, item := range candidates {
			if isStapleRecipe(item) {
				stapleAvailable = true
				break
			}
		}
	}

	used := make(map[string]bool)
	combos := make([]mealCombo, 0, comboLimit*2)

	// pushCombo 负责：
	// 1) 汇总组合营养；2) 以菜品 ID 排序后的 key 去重；3) 构建 mealCombo。
	pushCombo := func(items []models.Recipe, boost float64) {
		if requireStaple && stapleAvailable && !mealHasStaple(items) {
			return
		}
		if requireStaple && len(items) > 1 && countMainStapleRecipes(items) > 1 {
			return
		}

		combo, key, ok := buildMealCombo(items, boost)
		if !ok {
			return
		}
		if used[key] {
			return
		}
		used[key] = true

		combos = append(combos, combo)
	}

	// 生成 1/2/3 道菜组合（受 maxItems 控制），并限制搜索空间避免组合爆炸。
	for i := 0; i < len(candidates); i++ {
		if minItems <= 1 {
			pushCombo([]models.Recipe{candidates[i]}, 0)
		}
		if maxItems < 2 {
			continue
		}

		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].ID == candidates[j].ID {
				continue
			}
			if minItems <= 2 {
				pushCombo([]models.Recipe{candidates[i], candidates[j]}, 0)
			}

			if maxItems < 3 {
				continue
			}

			for k := j + 1; k < len(candidates); k++ {
				if candidates[k].ID == candidates[i].ID || candidates[k].ID == candidates[j].ID {
					continue
				}
				pushCombo([]models.Recipe{candidates[i], candidates[j], candidates[k]}, 0)
			}
		}
	}

	return s.rankMealCombos(combos, targetE, targetP, targetC, targetF, goal, comboLimit)
}

// buildBreakfastCombos 早餐专用组合器。
// 在通用评分基础上，额外优先“酸奶+鸡蛋”作为稳定命中组合。
func (s *RecipeServiceImpl) buildBreakfastCombos(recipes []models.Recipe, targetE, targetP, targetC, targetF float64, goal models.HealthGoal) []mealCombo {
	if len(recipes) == 0 {
		return nil
	}

	breakfastSource := make([]models.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		if isMainStapleRecipe(recipe) {
			continue
		}
		breakfastSource = append(breakfastSource, recipe)
	}
	if len(breakfastSource) >= 6 {
		recipes = breakfastSource
	}

	candidates := s.selectDiverseCandidates(recipes, targetE, targetP, targetC, targetF, 24, goal)
	if len(candidates) == 0 {
		return nil
	}

	// 在候选集合中做关键词匹配，用于识别“酸奶/鸡蛋”类食物。
	findByKeyword := func(keys ...string) (models.Recipe, bool) {
		for _, recipe := range candidates {
			content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))
			matched := false
			for _, key := range keys {
				if strings.Contains(content, strings.ToLower(key)) {
					matched = true
					break
				}
			}
			if matched {
				return recipe, true
			}
		}
		return models.Recipe{}, false
	}

	used := make(map[string]bool)
	combos := make([]mealCombo, 0, 12)
	// 早餐组合去重与营养汇总。
	pushCombo := func(items []models.Recipe, boost float64) {
		if countMainStapleRecipes(items) > 0 {
			return
		}
		if len(items) > 1 && countStapleRecipes(items) > 1 {
			return
		}

		combo, key, ok := buildMealCombo(items, boost)
		if !ok {
			return
		}
		if used[key] {
			return
		}
		used[key] = true

		combos = append(combos, combo)
	}

	// 优先注入“酸奶+鸡蛋”组合，提高在测试数据下的稳定产出概率。
	cfg := getRecipeTextConfig()
	yogurtKeys := cfg.BreakfastPriorityKeywords["yogurt"]
	eggKeys := cfg.BreakfastPriorityKeywords["egg"]
	if yogurt, okY := findByKeyword(yogurtKeys...); okY {
		if egg, okE := findByKeyword(eggKeys...); okE && yogurt.ID != egg.ID {
			pushCombo([]models.Recipe{yogurt, egg}, 12)
		}
	}

	for i := 0; i < len(candidates); i++ {
		pushCombo([]models.Recipe{candidates[i]}, 0)
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].ID == candidates[j].ID {
				continue
			}
			pushCombo([]models.Recipe{candidates[i], candidates[j]}, 0)
		}
	}

	return s.rankMealCombos(combos, targetE, targetP, targetC, targetF, goal, 6)
}

// hasRecipeOverlap 判断两个餐次组合是否存在相同菜品 ID。
// 用于跨餐去重，避免同一道菜在不同餐次重复出现。
func (s *RecipeServiceImpl) hasRecipeOverlap(left mealCombo, right mealCombo) bool {
	leftSet := left.ItemIDSet
	if len(leftSet) == 0 {
		leftSet = buildItemIDSet(left.Items)
	}

	rightSet := right.ItemIDSet
	if len(rightSet) == 0 {
		rightSet = buildItemIDSet(right.Items)
	}

	for id := range leftSet {
		if _, ok := rightSet[id]; ok {
			return true
		}
	}

	return false
}

// ingredientDiversityPenalty 计算跨餐食材重叠惩罚。
// 惩罚越大代表组合越不多样，用于从总分中扣减。
func (s *RecipeServiceImpl) ingredientDiversityPenalty(combos ...mealCombo) float64 {
	sets := make([]map[string]struct{}, 0, len(combos))
	for _, combo := range combos {
		set := combo.IngredientSet
		if len(set) == 0 {
			set = buildIngredientSet(combo.Items)
		}
		sets = append(sets, set)
	}

	overlapCount := 0
	// 统计任意两餐之间的食材交集数量。
	for i := 0; i < len(sets); i++ {
		for j := i + 1; j < len(sets); j++ {
			for ingredient := range sets[i] {
				if _, ok := sets[j][ingredient]; ok {
					overlapCount++
				}
			}
		}
	}

	return float64(overlapCount) * 1.5
}

func recommendationScoreFloor(comboSizes ...int) float64 {
	floor := 60.0
	minComboSize := math.MaxInt
	for _, size := range comboSizes {
		if size <= 0 {
			continue
		}
		if size < minComboSize {
			minComboSize = size
		}
	}

	if minComboSize == math.MaxInt {
		return floor
	}
	if minComboSize <= 4 {
		return 35
	}
	if minComboSize <= 8 {
		return 45
	}
	if minComboSize <= 12 {
		return 50
	}
	return floor
}

// generateOneDailyPlan 生成一套每日食谱计划
func (s *RecipeServiceImpl) generateRankedPlans(
	userID uint,
	breakfastRecipes, lunchRecipes, dinnerRecipes, snackRecipes []models.Recipe,
	target NutritionTarget, userGoal models.HealthGoal,
) []*models.DailyRecipePlan {

	// 1. 分配每餐的目标营养 (按比例)
	// 早餐 25%, 午餐 35%, 晚餐 30%, 加餐 10%
	ratios := []float64{0.25, 0.35, 0.30, 0.10}

	// 辅助函数：计算单餐目标
	getTarget := func(ratio float64) (float64, float64, float64, float64) {
		return target.Energy * ratio, target.Protein * ratio, target.Carbohydrate * ratio, target.Fat * ratio
	}

	// 2. 优化后的初筛：传入所有营养目标
	te, tp, tc, tf := getTarget(ratios[0])
	breakfastCombos := s.buildBreakfastCombos(breakfastRecipes, te, tp, tc, tf, userGoal)
	breakfastCombos = s.filterCombosByMealProfile(models.MealTypeBreakfast, breakfastCombos, te, tp, tc, tf, userGoal)

	te, tp, tc, tf = getTarget(ratios[1])
	// 午餐使用通用组合生成：最多 3 道菜。
	lunchCombos := s.buildMealCombos(lunchRecipes, te, tp, tc, tf, userGoal, 28, 8, 2, 3, true)
	lunchCombos = s.filterCombosByMealProfile(models.MealTypeLunch, lunchCombos, te, tp, tc, tf, userGoal)

	te, tp, tc, tf = getTarget(ratios[2])
	// 晚餐使用通用组合生成：最多 3 道菜。
	dinnerCombos := s.buildMealCombos(dinnerRecipes, te, tp, tc, tf, userGoal, 28, 8, 2, 3, true)
	dinnerCombos = s.filterCombosByMealProfile(models.MealTypeDinner, dinnerCombos, te, tp, tc, tf, userGoal)

	te, tp, tc, tf = getTarget(ratios[3])
	// 加餐允许小组合：最多 2 道菜。
	snackCombos := s.buildMealCombos(snackRecipes, te, tp, tc, tf, userGoal, 16, 6, 1, 2, false)
	snackCombos = s.filterCombosByMealProfile(models.MealTypeSnack, snackCombos, te, tp, tc, tf, userGoal)
	if len(snackCombos) == 0 {
		// 兜底空加餐，保持旧流程兼容。
		snackCombos = append(snackCombos, mealCombo{Items: []models.Recipe{{Name: ""}}, Primary: models.Recipe{Name: ""}})
	}

	log.Printf("Meal combo sizes -> Breakfast: %d, Lunch: %d, Dinner: %d, Snack: %d", len(breakfastCombos), len(lunchCombos), len(dinnerCombos), len(snackCombos))
	minPlanScore := recommendationScoreFloor(len(breakfastCombos), len(lunchCombos), len(dinnerCombos), len(snackCombos))

	if len(breakfastCombos) == 0 || len(lunchCombos) == 0 || len(dinnerCombos) == 0 {
		return nil
	}

	allPlans := make([]*models.DailyRecipePlan, 0, 2000)

	// 3. 循环与剪枝
	// 设定热量容忍阈值 (例如 ±2000大卡明显是不可能的，早点break)
	// 这里设定一个宽松的剪枝范围：如果当前热量已经超过目标值的 1.3倍，就没必要继续加菜了
	maxEnergyLimit := target.Energy * 1.3

	for _, bCombo := range breakfastCombos {
		for _, lCombo := range lunchCombos {
			// 跨餐菜品去重：早餐与午餐不能重复同菜。
			if s.hasRecipeOverlap(bCombo, lCombo) {
				continue
			}

			// [剪枝 1]：如果早+午已经热量爆表，跳过晚餐
			currentE := bCombo.Energy + lCombo.Energy
			if currentE > maxEnergyLimit {
				continue
			}

			for _, dCombo := range dinnerCombos {
				// 跨餐菜品去重：午/晚与早/晚不能重复同菜。
				if s.hasRecipeOverlap(lCombo, dCombo) || s.hasRecipeOverlap(bCombo, dCombo) {
					continue
				}

				// [剪枝 2]：早+午+晚 热量爆表
				if currentE+dCombo.Energy > maxEnergyLimit {
					continue
				}

				for _, snCombo := range snackCombos {
					// 加餐也参与跨餐去重（空加餐除外）。
					if snCombo.Primary.ID != 0 && (s.hasRecipeOverlap(snCombo, bCombo) || s.hasRecipeOverlap(snCombo, lCombo) || s.hasRecipeOverlap(snCombo, dCombo)) {
						continue
					}

					// 计算当前组合的总营养
					totalE := bCombo.Energy + lCombo.Energy + dCombo.Energy + snCombo.Energy
					totalP := bCombo.Protein + lCombo.Protein + dCombo.Protein + snCombo.Protein
					totalC := bCombo.Carbohydrate + lCombo.Carbohydrate + dCombo.Carbohydrate + snCombo.Carbohydrate
					totalF := bCombo.Fat + lCombo.Fat + dCombo.Fat + snCombo.Fat

					matchScore := s.calculateMatchScore(
						totalE, totalP, totalC, totalF,
						target.Energy, target.Protein, target.Carbohydrate, target.Fat,
						userGoal,
					)
					// 对食材重复进行扣分，鼓励全日多样性。
					matchScore -= s.ingredientDiversityPenalty(bCombo, lCombo, dCombo, snCombo)
					if matchScore < 0 {
						matchScore = 0
					}

					// 小样本手工库下允许适度放宽门槛，避免组合被全部筛空。
					if matchScore < minPlanScore {
						continue
					}

					lPrimary := lCombo.Primary
					dPrimary := dCombo.Primary
					snPrimary := snCombo.Primary

					plan := &models.DailyRecipePlan{
						UserID:            userID,
						BreakfastRecipe:   bCombo.Primary,
						BreakfastRecipeID: bCombo.Primary.ID,
						LunchRecipe:       lPrimary,
						LunchRecipeID:     lPrimary.ID,
						DinnerRecipe:      dPrimary,
						DinnerRecipeID:    dPrimary.ID,
						SnackRecipe:       snPrimary, // 如果是空占位，name为空
						SnackRecipeID:     snPrimary.ID,
						BreakfastItems:    bCombo.Items,
						LunchItems:        lCombo.Items,
						DinnerItems:       dCombo.Items,
						SnackItems:        snCombo.Items,

						TotalEnergy:       totalE,
						TotalProtein:      totalP,
						TotalCarbohydrate: totalC,
						TotalFat:          totalF,

						TargetEnergy:       target.Energy,
						TargetProtein:      target.Protein,
						TargetCarbohydrate: target.Carbohydrate,
						TargetFat:          target.Fat,
						MatchScore:         matchScore,
						IsSelected:         false,
						PlanDate:           time.Now(),
					}
					allPlans = append(allPlans, plan)
				}

			}
		}
	}

	// 排序
	sort.Slice(allPlans, func(i, j int) bool {
		return allPlans[i].MatchScore > allPlans[j].MatchScore
	})

	if len(allPlans) == 0 {
		fallbackPlans := s.buildFallbackPlansFromCombos(userID, breakfastCombos, lunchCombos, dinnerCombos, snackCombos, target, userGoal)
		if len(fallbackPlans) > 0 {
			log.Printf("generateRankedPlans fallback activated: built %d plan(s) from existing meal combos", len(fallbackPlans))
			return fallbackPlans
		}
	}

	return allPlans
}

func countComboRecipeOverlap(left, right mealCombo) int {
	leftSet := left.ItemIDSet
	if len(leftSet) == 0 {
		leftSet = buildItemIDSet(left.Items)
	}
	rightSet := right.ItemIDSet
	if len(rightSet) == 0 {
		rightSet = buildItemIDSet(right.Items)
	}

	count := 0
	for id := range leftSet {
		if _, ok := rightSet[id]; ok {
			count++
		}
	}
	return count
}

func firstRecipeOfCombo(combo mealCombo) models.Recipe {
	if combo.Primary.ID > 0 || combo.Primary.Name != "" {
		return combo.Primary
	}
	if len(combo.Items) > 0 {
		return combo.Items[0]
	}
	return models.Recipe{}
}

func fallbackPlanKey(bCombo, lCombo, dCombo, snCombo mealCombo) string {
	keyOf := func(combo mealCombo) string {
		ids := make([]uint, 0, len(combo.Items))
		for _, item := range combo.Items {
			if item.ID > 0 {
				ids = append(ids, item.ID)
			}
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		if len(ids) == 0 {
			return "0"
		}

		parts := make([]string, 0, len(ids))
		for _, id := range ids {
			parts = append(parts, strconv.FormatUint(uint64(id), 10))
		}
		return strings.Join(parts, "-")
	}

	return strings.Join([]string{keyOf(bCombo), keyOf(lCombo), keyOf(dCombo), keyOf(snCombo)}, "|")
}

func (s *RecipeServiceImpl) chooseFallbackCombo(candidates []mealCombo, existing []mealCombo) mealCombo {
	if len(candidates) == 0 {
		return mealCombo{}
	}

	best := candidates[0]
	bestRecipeOverlap := math.MaxInt
	bestPenalty := math.MaxFloat64
	bestIndex := math.MaxInt

	for index, candidate := range candidates {
		recipeOverlap := 0
		for _, chosen := range existing {
			recipeOverlap += countComboRecipeOverlap(candidate, chosen)
		}

		all := append(append([]mealCombo{}, existing...), candidate)
		penalty := s.ingredientDiversityPenalty(all...)
		if recipeOverlap < bestRecipeOverlap ||
			(recipeOverlap == bestRecipeOverlap && penalty < bestPenalty) ||
			(recipeOverlap == bestRecipeOverlap && penalty == bestPenalty && index < bestIndex) {
			best = candidate
			bestRecipeOverlap = recipeOverlap
			bestPenalty = penalty
			bestIndex = index
		}
	}

	return best
}

func (s *RecipeServiceImpl) buildPlanFromCombos(userID uint, bCombo, lCombo, dCombo, snCombo mealCombo, target NutritionTarget, userGoal models.HealthGoal) *models.DailyRecipePlan {
	bPrimary := firstRecipeOfCombo(bCombo)
	lPrimary := firstRecipeOfCombo(lCombo)
	dPrimary := firstRecipeOfCombo(dCombo)
	snPrimary := firstRecipeOfCombo(snCombo)

	totalE := bCombo.Energy + lCombo.Energy + dCombo.Energy + snCombo.Energy
	totalP := bCombo.Protein + lCombo.Protein + dCombo.Protein + snCombo.Protein
	totalC := bCombo.Carbohydrate + lCombo.Carbohydrate + dCombo.Carbohydrate + snCombo.Carbohydrate
	totalF := bCombo.Fat + lCombo.Fat + dCombo.Fat + snCombo.Fat

	matchScore := s.calculateMatchScore(
		totalE, totalP, totalC, totalF,
		target.Energy, target.Protein, target.Carbohydrate, target.Fat,
		userGoal,
	)
	matchScore -= s.ingredientDiversityPenalty(bCombo, lCombo, dCombo, snCombo)
	if matchScore < 0 {
		matchScore = 0
	}

	return &models.DailyRecipePlan{
		UserID:             userID,
		BreakfastRecipe:    bPrimary,
		BreakfastRecipeID:  bPrimary.ID,
		LunchRecipe:        lPrimary,
		LunchRecipeID:      lPrimary.ID,
		DinnerRecipe:       dPrimary,
		DinnerRecipeID:     dPrimary.ID,
		SnackRecipe:        snPrimary,
		SnackRecipeID:      snPrimary.ID,
		BreakfastItems:     bCombo.Items,
		LunchItems:         lCombo.Items,
		DinnerItems:        dCombo.Items,
		SnackItems:         snCombo.Items,
		TotalEnergy:        totalE,
		TotalProtein:       totalP,
		TotalCarbohydrate:  totalC,
		TotalFat:           totalF,
		TargetEnergy:       target.Energy,
		TargetProtein:      target.Protein,
		TargetCarbohydrate: target.Carbohydrate,
		TargetFat:          target.Fat,
		MatchScore:         matchScore,
		IsSelected:         false,
		PlanDate:           time.Now(),
	}
}

func (s *RecipeServiceImpl) buildFallbackPlansFromCombos(userID uint, breakfastCombos, lunchCombos, dinnerCombos, snackCombos []mealCombo, target NutritionTarget, userGoal models.HealthGoal) []*models.DailyRecipePlan {
	if len(breakfastCombos) == 0 || len(lunchCombos) == 0 || len(dinnerCombos) == 0 {
		return nil
	}
	if len(snackCombos) == 0 {
		snackCombos = []mealCombo{{Items: []models.Recipe{{Name: ""}}, Primary: models.Recipe{Name: ""}}}
	}

	maxBreakfast := min(len(breakfastCombos), 3)
	maxLunch := min(len(lunchCombos), 3)
	plans := make([]*models.DailyRecipePlan, 0, maxBreakfast*maxLunch)
	seen := make(map[string]struct{})

	for bi := 0; bi < maxBreakfast; bi++ {
		bCombo := breakfastCombos[bi]
		for li := 0; li < maxLunch; li++ {
			lCombo := lunchCombos[li]
			existing := []mealCombo{bCombo, lCombo}
			dCombo := s.chooseFallbackCombo(dinnerCombos, existing)
			existing = append(existing, dCombo)
			snCombo := s.chooseFallbackCombo(snackCombos, existing)

			key := fallbackPlanKey(bCombo, lCombo, dCombo, snCombo)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			plans = append(plans, s.buildPlanFromCombos(userID, bCombo, lCombo, dCombo, snCombo, target, userGoal))
		}
	}

	sort.Slice(plans, func(i, j int) bool {
		return plans[i].MatchScore > plans[j].MatchScore
	})
	return plans
}

// shuffleByWeights 使用加权 Fisher-Yates 变种算法对方案进行洗牌
// 算法逻辑：
// 1. 遍历位置 i 从 0 到 N-1
// 2. 计算从 i 到 N-1 所有剩余元素的总权重
// 3. 在剩余总权重中随机生成一个值 r
// 4. 找到 r 对应的元素（累加权重法）
// 5. 将该元素与当前位置 i 的元素交换
func (s *RecipeServiceImpl) shuffleByWeights(candidates []*models.DailyRecipePlan) []*models.DailyRecipePlan {
	n := len(candidates)
	if n == 0 {
		return nil
	}

	// 只对前 Top N 进行洗牌（比如前100个），太后面的还是截断掉比较好
	// 如果列表过长，计算权重的开销会增大，限制范围是明智的
	limit := min(n, 100)

	// 直接在 candidates 原切片上操作的前 limit 个元素
	pool := candidates[:limit]

	// 预计算所有元素的权重，避免在双重循环中重复调用 math.Pow
	// 空间换时间：O(N) 空间换取 O(N) 次 Pow 计算
	weights := make([]float64, limit)
	for i, p := range pool {
		// 使用立方让高分项被选中的概率显著增加
		weights[i] = math.Pow(p.MatchScore, 3)
	}

	// 开始 Fisher-Yates 加权洗牌
	// i 代表当前要填充的位置
	for i := 0; i < limit-1; i++ {
		// 1. 计算剩余部分的总权重 (从 i 到 end)
		// 这一步虽然是 O(N)，但在 N=100 时非常快，且比维护树状数组简单
		remainingWeight := 0.0
		for j := i; j < limit; j++ {
			remainingWeight += weights[j]
		}

		// 2. 生成随机阈值
		r := s.rng.Float64() * remainingWeight

		// 3. 在剩余元素中寻找命中者
		currentSum := 0.0
		winnerIdx := -1

		for j := i; j < limit; j++ {
			currentSum += weights[j]
			if r <= currentSum {
				winnerIdx = j
				break
			}
		}

		// 浮点数兜底：如果没找到（极少数情况），默认选最后一个
		if winnerIdx == -1 {
			winnerIdx = limit - 1
		}

		// 4. 交换：将中奖者放到当前位置 i
		if winnerIdx != i {
			// 交换元素
			pool[i], pool[winnerIdx] = pool[winnerIdx], pool[i]
			// 同步交换对应的权重，保证下次循环逻辑正确
			weights[i], weights[winnerIdx] = weights[winnerIdx], weights[i]
		}
	}

	return pool
}

// 获取不同健康目标的权重向量
// 增肌：非常看重蛋白质
// 减脂：非常看重热量和脂肪控制
func getNutrientWeights(goal models.HealthGoal) (wE, wP, wC, wF float64) {
	switch goal {
	case models.GoalMuscleGain:
		return 0.4, 0.4, 0.1, 0.1 // 能量40%，蛋白40%
	case models.GoalWeightLoss:
		return 0.5, 0.1, 0.1, 0.3 // 能量50%，脂肪30%
	case models.GoalSugarControl:
		return 0.4, 0.2, 0.3, 0.1 // 能量40%，碳水30%
	default:
		return 0.4, 0.2, 0.2, 0.2 // 均衡
	}
}

// scoreRecipe 单个食谱评分：计算单个食谱与单餐目标的匹配度
func (s *RecipeServiceImpl) scoreSingleRecipe(r models.Recipe, targetE, targetP, targetC, targetF float64, goal models.HealthGoal) float64 {
	r = scaleRecipeNutritionByPortion(r)

	// 获取权重（复用你现有的逻辑）
	wE, wP, wC, wF := getNutrientWeights(goal)

	// 计算偏差率 (使用相对偏差)
	// 避免分母为0
	safeDiv := func(val, target float64) float64 {
		if target == 0 {
			return 0
		}
		diff := math.Abs(val - target)
		return diff / target // 偏差百分比
	}

	scoreE := safeDiv(r.Energy, targetE)
	scoreP := safeDiv(r.Protein, targetP)
	scoreC := safeDiv(r.Carbohydrate, targetC)
	scoreF := safeDiv(r.Fat, targetF)

	// 综合偏差 (越小越好)
	totalDiff := scoreE*wE + scoreP*wP + scoreC*wC + scoreF*wF

	// 转化为得分 (100分制，偏差越大分越低)
	return math.Max(0, 100*(1-totalDiff))
}

// 获取 Top-K 候选食谱,基于综合评分筛选
func (s *RecipeServiceImpl) getTopCandidates(recipes []models.Recipe,
	targetE, targetP, targetC, targetF float64, // 传入该餐的具体营养目标
	k int, goal models.HealthGoal) []models.Recipe {
	if len(recipes) == 0 {
		return nil
	}

	type candidate struct {
		r     models.Recipe
		score float64
	}

	var cs []candidate
	for _, r := range recipes {
		score := s.scoreSingleRecipe(r, targetE, targetP, targetC, targetF, goal)
		cs = append(cs, candidate{r, score})
	}

	// 分数从高到低排序
	sort.Slice(cs, func(i, j int) bool { return cs[i].score > cs[j].score })

	// 不直接取前 k 个，而是取前 3*k 个作为“候选池”，然后从中随机选 k 个
	// 直接取排序后的 Top-K，避免把质量较差的随机候选混入早餐/加餐。
	limit := min(len(cs), k)

	result := make([]models.Recipe, limit)
	for i := range limit {
		result[i] = cs[i].r
	}
	return result
}

// calculateMatchScore 计算营养匹配度（0-100）
func (s *RecipeServiceImpl) calculateMatchScore(
	actualE, actualP, actualC, actualF,
	targetE, targetP, targetC, targetF float64,
	userGoal models.HealthGoal,
) float64 {
	// 避免除以零
	if targetE == 0 || targetP == 0 || targetC == 0 || targetF == 0 {
		return 0
	}

	// 计算百分比偏差 (Delta Percentage)
	diffE := math.Abs(actualE-targetE) / targetE
	diffP := math.Abs(actualP-targetP) / targetP
	diffC := math.Abs(actualC-targetC) / targetC
	diffF := math.Abs(actualF-targetF) / targetF
	// 获取权重
	wE, wP, wC, wF := getNutrientWeights(userGoal)

	dist := math.Sqrt(wE*diffE*diffE + wP*diffP*diffP + wC*diffC*diffC + wF*diffF*diffF)

	var score float64
	if dist <= 0.10 {
		// 偏差极小：95 ~ 100 分 (线性映射)
		// dist=0 -> 100; dist=0.1 -> 95
		score = 100 - (dist * 50)
	} else if dist <= 0.25 {
		// 偏差中等：85 ~ 95 分
		// dist=0.1 -> 95; dist=0.25 -> 85
		score = 95 - ((dist - 0.1) * 66.6)
	} else {
		// 偏差较大：分数快速下降
		// score = 85 * e^(-2 * (dist-0.25))
		score = 85 * math.Exp(-2*(dist-0.25))
	}

	return math.Round(score)
}

// normalizePlanRecipeIDs 将计划中的餐次ID归一化为 recipes 表ID，避免外键约束失败。
func (s *RecipeServiceImpl) normalizePlanRecipeIDs(plan *models.DailyRecipePlan) error {
	if plan == nil {
		return fmt.Errorf("计划不能为空")
	}

	normalizeIDs := func(ids []uint) ([]uint, error) {
		if len(ids) == 0 {
			return ids, nil
		}
		result := make([]uint, 0, len(ids))
		for _, id := range ids {
			mappedID, mapErr := s.recipeRepo.ResolveRecipeIDForPlan(id)
			if mapErr != nil {
				return nil, mapErr
			}
			if mappedID > 0 {
				result = append(result, mappedID)
			}
		}
		return result, nil
	}

	var err error
	plan.BreakfastItemIDs, err = normalizeIDs(plan.BreakfastItemIDs)
	if err != nil {
		return fmt.Errorf("早餐组合ID映射失败: %w", err)
	}
	plan.LunchItemIDs, err = normalizeIDs(plan.LunchItemIDs)
	if err != nil {
		return fmt.Errorf("午餐组合ID映射失败: %w", err)
	}
	plan.DinnerItemIDs, err = normalizeIDs(plan.DinnerItemIDs)
	if err != nil {
		return fmt.Errorf("晚餐组合ID映射失败: %w", err)
	}
	plan.SnackItemIDs, err = normalizeIDs(plan.SnackItemIDs)
	if err != nil {
		return fmt.Errorf("加餐组合ID映射失败: %w", err)
	}

	breakfastID, err := s.recipeRepo.ResolveRecipeIDForPlan(plan.BreakfastRecipeID)
	if err != nil {
		return fmt.Errorf("早餐食谱ID映射失败: %w", err)
	}
	lunchID, err := s.recipeRepo.ResolveRecipeIDForPlan(plan.LunchRecipeID)
	if err != nil {
		return fmt.Errorf("午餐食谱ID映射失败: %w", err)
	}
	dinnerID, err := s.recipeRepo.ResolveRecipeIDForPlan(plan.DinnerRecipeID)
	if err != nil {
		return fmt.Errorf("晚餐食谱ID映射失败: %w", err)
	}
	snackID, err := s.recipeRepo.ResolveRecipeIDForPlan(plan.SnackRecipeID)
	if err != nil {
		return fmt.Errorf("加餐食谱ID映射失败: %w", err)
	}

	if breakfastID == 0 || lunchID == 0 || dinnerID == 0 {
		return fmt.Errorf("计划主餐食谱ID无效(早:%d 午:%d 晚:%d)", breakfastID, lunchID, dinnerID)
	}

	plan.BreakfastRecipeID = breakfastID
	plan.LunchRecipeID = lunchID
	plan.DinnerRecipeID = dinnerID
	plan.SnackRecipeID = snackID

	if len(plan.BreakfastItemIDs) == 0 && breakfastID > 0 {
		plan.BreakfastItemIDs = []uint{breakfastID}
	}
	if len(plan.LunchItemIDs) == 0 && lunchID > 0 {
		plan.LunchItemIDs = []uint{lunchID}
	}
	if len(plan.DinnerItemIDs) == 0 && dinnerID > 0 {
		plan.DinnerItemIDs = []uint{dinnerID}
	}
	if len(plan.SnackItemIDs) == 0 && snackID > 0 {
		plan.SnackItemIDs = []uint{snackID}
	}

	return nil
}

// SaveDailyPlan 保存每日食谱计划
func (s *RecipeServiceImpl) SaveDailyPlan(plan *models.DailyRecipePlan) error {
	if err := s.normalizePlanRecipeIDs(plan); err != nil {
		return err
	}
	return s.recipeRepo.CreateDailyPlan(plan)
}

// GetSelectedPlan 获取用户当前选中的计划
func (s *RecipeServiceImpl) GetSelectedPlan(userID uint) (*models.DailyRecipePlan, error) {
	// FindSelectedPlan已经只查询今天的计划了,如果查不到就是没有或已过期
	plan, err := s.recipeRepo.FindSelectedPlan(userID)
	if err != nil {
		log.Printf("FindSelectedPlan error: %v", err)
		return nil, fmt.Errorf("当前计划已过期,请重新获取推荐")
	}

	buildItems := func(ids []uint, fallback models.Recipe) []models.Recipe {
		if len(ids) == 0 {
			if fallback.ID > 0 {
				return []models.Recipe{scaleRecipeNutritionByPortion(fallback)}
			}
			return []models.Recipe{}
		}

		recipes, loadErr := s.recipeRepo.FindByIDsFromRuntimeTable(ids)
		if loadErr != nil {
			if fallback.ID > 0 {
				return []models.Recipe{scaleRecipeNutritionByPortion(fallback)}
			}
			return []models.Recipe{}
		}

		recipeMap := make(map[uint]models.Recipe)
		for _, recipe := range recipes {
			recipeMap[recipe.ID] = recipe
		}

		ordered := make([]models.Recipe, 0, len(ids))
		for _, id := range ids {
			if recipe, ok := recipeMap[id]; ok {
				ordered = append(ordered, scaleRecipeNutritionByPortion(recipe))
			}
		}

		if len(ordered) == 0 && fallback.ID > 0 {
			ordered = append(ordered, scaleRecipeNutritionByPortion(fallback))
		}

		return ordered
	}

	plan.BreakfastItems = buildItems(plan.BreakfastItemIDs, plan.BreakfastRecipe)
	plan.LunchItems = buildItems(plan.LunchItemIDs, plan.LunchRecipe)
	plan.DinnerItems = buildItems(plan.DinnerItemIDs, plan.DinnerRecipe)
	plan.SnackItems = buildItems(plan.SnackItemIDs, plan.SnackRecipe)

	if len(plan.BreakfastItems) > 0 {
		plan.BreakfastRecipe = plan.BreakfastItems[0]
	} else {
		plan.BreakfastRecipe = scaleRecipeNutritionByPortion(plan.BreakfastRecipe)
	}
	if len(plan.LunchItems) > 0 {
		plan.LunchRecipe = plan.LunchItems[0]
	} else {
		plan.LunchRecipe = scaleRecipeNutritionByPortion(plan.LunchRecipe)
	}
	if len(plan.DinnerItems) > 0 {
		plan.DinnerRecipe = plan.DinnerItems[0]
	} else {
		plan.DinnerRecipe = scaleRecipeNutritionByPortion(plan.DinnerRecipe)
	}
	if len(plan.SnackItems) > 0 {
		plan.SnackRecipe = plan.SnackItems[0]
	} else {
		plan.SnackRecipe = scaleRecipeNutritionByPortion(plan.SnackRecipe)
	}

	return plan, nil
}

// SelectDailyPlan 选择每日食谱计划
func (s *RecipeServiceImpl) SelectDailyPlan(userID uint, plan *models.DailyRecipePlan) error {
	// 0. 先做ID归一化，避免后续清理成功但保存失败
	if err := s.normalizePlanRecipeIDs(plan); err != nil {
		return fmt.Errorf("计划食谱ID校验失败: %w", err)
	}

	// 1. 先取消所有历史的is_selected=true记录(解决多条选中记录问题)
	if err := s.recipeRepo.ClearAllSelectedPlans(userID); err != nil {
		return fmt.Errorf("取消历史选中记录失败: %w", err)
	}

	// 2. 删除用户当天所有计划(包括刚取消选中的,保证每天只有一份)
	if err := s.recipeRepo.DeleteTodayAllPlans(userID); err != nil {
		return fmt.Errorf("删除当天计划失败: %w", err)
	}

	// 3. 设置新计划的属性
	plan.UserID = userID
	plan.IsSelected = true
	// 只保存日期部分,不包含时间(与type:date字段匹配)
	now := time.Now()
	plan.PlanDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 4. 保存新计划
	return s.recipeRepo.CreateDailyPlan(plan)
}

// GetRecipeDetail 获取食谱详情（含收藏状态）
func (s *RecipeServiceImpl) GetRecipeDetail(recipeID uint, userID uint) (*models.Recipe, bool, error) {
	// 查询食谱
	recipe, err := s.recipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, false, err
	}
	scaledRecipe := scaleRecipeNutritionByPortion(*recipe)
	recipe = &scaledRecipe

	// 检查是否已收藏
	isFavorite := false
	if userID > 0 {
		isFavorite, _ = s.recipeRepo.IsFavorite(userID, recipeID)
	}

	return recipe, isFavorite, nil
}

// AddFavorite 添加收藏
func (s *RecipeServiceImpl) AddFavorite(userID, recipeID uint) error {
	// 检查食谱是否存在
	_, err := s.recipeRepo.FindByID(recipeID)
	if err != nil {
		return err
	}

	// 已收藏则直接返回，避免重复插入
	if ok, _ := s.recipeRepo.IsFavorite(userID, recipeID); ok {
		return nil
	}

	// 添加收藏
	return s.recipeRepo.AddFavorite(userID, recipeID)
}

// RemoveFavorite 取消收藏
func (s *RecipeServiceImpl) RemoveFavorite(userID, recipeID uint) error {
	return s.recipeRepo.RemoveFavorite(userID, recipeID)
}

// GetUserFavorites 获取用户收藏列表
func (s *RecipeServiceImpl) GetUserFavorites(userID uint) ([]models.Recipe, error) {
	return s.recipeRepo.GetUserFavorites(userID)
}

// IsFavorite 检查收藏状态
func (s *RecipeServiceImpl) IsFavorite(userID, recipeID uint) (bool, error) {
	return s.recipeRepo.IsFavorite(userID, recipeID)
}

// filterRecentRecipes 过滤最近吃过的食谱
func (s *RecipeServiceImpl) filterRecentRecipes(recipes []models.Recipe, recentIDs []uint) []models.Recipe {
	recentMap := make(map[uint]bool)
	for _, id := range recentIDs {
		recentMap[id] = true
	}

	filtered := make([]models.Recipe, 0)
	for _, recipe := range recipes {
		if !recentMap[recipe.ID] {
			filtered = append(filtered, recipe)
		}
	}

	// 如果过滤后数量太少（少于3个），则不过滤，避免无方案可选
	if len(filtered) < 3 {
		return recipes
	}
	return filtered
}

// parseUserPreferences 解析用户的标签与禁忌，用于后续筛选
// 返回 userTags: 用户拥有的标签（如健康状况、饮食偏好、目标等）
// 返回 forbidden: 需要在仓库层初筛的禁忌食材（过敏源）
func (s *RecipeServiceImpl) parseUserPreferences(user *models.User) (userTags []string, forbidden []string) {
	tags := make([]string, 0)
	forb := make([]string, 0)

	if user == nil {
		return tags, forb
	}

	// 健康状况（例如: 高血压, 糖尿病）
	if strings.TrimSpace(user.HealthConditions) != "" {
		parts := strings.SplitSeq(user.HealthConditions, ",")
		for p := range parts {
			v := strings.TrimSpace(p)
			if v != "" {
				tags = append(tags, v)
				// 健康状况既可能是用户的标签（用于 targetUsers 匹配），也可能是禁忌（例如某些菜禁忌高血压）
				forb = append(forb, v)
			}
		}
	}

	// 饮食偏好
	if strings.TrimSpace(user.DietaryPrefs) != "" {
		parts := strings.SplitSeq(user.DietaryPrefs, ",")
		for p := range parts {
			v := strings.TrimSpace(p)
			if v != "" {
				tags = append(tags, v)
			}
		}
	}

	// 过敏源作为严格的 forbidden
	if strings.TrimSpace(user.Allergies) != "" {
		parts := strings.SplitSeq(user.Allergies, ",")
		for p := range parts {
			v := strings.TrimSpace(p)
			if v != "" {
				forb = append(forb, v)
			}
		}
	}

	tags = append(tags, string(user.HealthGoal))
	return tags, forb
}

// filterRecipesByUserPreferences 基于 Recipe.TargetUsers 与 Recipe.ForbiddenUsers 进行精筛
func (s *RecipeServiceImpl) filterRecipesByUserPreferences(recipes []models.Recipe, userTags []string) []models.Recipe {
	if len(recipes) == 0 {
		return recipes
	}

	normalize := func(value string) string {
		return strings.ToLower(strings.TrimSpace(value))
	}

	tagSet := make(map[string]bool)
	for _, t := range userTags {
		t = normalize(t)
		if t != "" {
			tagSet[t] = true
		}
	}

	for _, tag := range []string{"大众", "通用", "维持健康"} {
		tagSet[normalize(tag)] = true
	}

	generalTargetSet := map[string]bool{
		normalize("大众"):   true,
		normalize("通用"):   true,
		normalize("维持健康"): true,
	}

	allowedRecipes := make([]models.Recipe, 0, len(recipes))
	preferredRecipes := make([]models.Recipe, 0, len(recipes))
	preferredIDs := make(map[uint]bool, len(recipes))
	for _, r := range recipes {
		skip := false
		for _, fu := range r.ForbiddenUsers {
			fu = normalize(fu)
			if fu == "" {
				continue
			}
			if tagSet[fu] {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		allowedRecipes = append(allowedRecipes, r)

		if len(r.TargetUsers) == 0 {
			preferredRecipes = append(preferredRecipes, r)
			preferredIDs[r.ID] = true
			continue
		}

		matched := false
		for _, tu := range r.TargetUsers {
			tu = normalize(tu)
			if tu == "" {
				continue
			}
			if generalTargetSet[tu] || tagSet[tu] {
				matched = true
				break
			}
		}
		if matched {
			preferredRecipes = append(preferredRecipes, r)
			preferredIDs[r.ID] = true
		}
	}

	if len(preferredRecipes) >= 3 {
		return preferredRecipes
	}
	if len(preferredRecipes) > 0 {
		result := make([]models.Recipe, 0, len(allowedRecipes))
		result = append(result, preferredRecipes...)
		for _, recipe := range allowedRecipes {
			if preferredIDs[recipe.ID] {
				continue
			}
			result = append(result, recipe)
		}
		return result
	}
	return allowedRecipes
}

func (s *RecipeServiceImpl) filterRecipesByHealthGoal(recipes []models.Recipe, goal models.HealthGoal) []models.Recipe {
	if len(recipes) == 0 {
		return recipes
	}

	highSugarDrinkKeywords := []string{
		"奶茶", "莫吉托", "mojito", "长岛冰茶", "鸡尾酒", "特调", "咖啡特调",
		"果茶", "水果茶", "气泡饮", "冰茶", "甜饮", "饮料", "奶昔", "糖水",
		"杨枝甘露", "冰淇淋", "奶冻", "提拉米苏", "雪花酥", "布丁",
	}
	containsKeyword := func(content string, keywords []string) bool {
		for _, keyword := range keywords {
			if keyword == "" {
				continue
			}
			if strings.Contains(content, strings.ToLower(keyword)) {
				return true
			}
		}
		return false
	}

	shouldExclude := func(recipe models.Recipe) bool {
		content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))
		switch goal {
		case models.GoalWeightLoss, models.GoalSugarControl:
			if containsKeyword(content, highSugarDrinkKeywords) {
				return true
			}
			if recipe.MealType == models.MealTypeSnack && recipe.Carbohydrate >= 18 && recipe.Energy >= 180 {
				return true
			}
		}
		return false
	}

	filtered := make([]models.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		if shouldExclude(recipe) {
			continue
		}
		filtered = append(filtered, recipe)
	}

	if len(filtered) < 3 {
		return recipes
	}

	return filtered
}

// checkDiversity 多样性检查
func (s *RecipeServiceImpl) checkDiversity(existingPlans []*models.DailyRecipePlan, newPlan *models.DailyRecipePlan) bool {
	for _, plan := range existingPlans {
		// 策略：只有当 午餐 AND 晚餐 的主菜都重复时，才判定为重复方案
		// 这样允许早餐重复（大家早餐常吃一样的），也允许只有一顿主菜重复
		lunchDup := plan.LunchRecipeID == newPlan.LunchRecipeID
		dinnerDup := plan.DinnerRecipeID == newPlan.DinnerRecipeID

		if lunchDup && dinnerDup {
			return false // 太像了，不要
		}

		// 进阶：如果 4 顿里有 3 顿完全一样，也不要
		sameCount := 0
		if plan.BreakfastRecipeID == newPlan.BreakfastRecipeID {
			sameCount++
		}
		if lunchDup {
			sameCount++
		}
		if dinnerDup {
			sameCount++
		}
		if plan.SnackRecipeID == newPlan.SnackRecipeID {
			sameCount++
		}

		if sameCount >= 3 {
			return false
		}
	}
	return true
}

func recipeVarietyKey(recipe models.Recipe) string {
	content := strings.ToLower(recipe.Name + " " + strings.Join(recipe.Ingredients, " "))

	mainStapleKeywords := []string{
		"炒饭", "蛋包饭", "盖饭", "烩饭", "焗饭", "拌饭",
		"炒面", "拌面", "意面", "炒粉", "炒米粉", "河粉", "米线", "拉面", "焖面",
		"饺子", "馄饨",
	}
	for _, keyword := range mainStapleKeywords {
		if strings.Contains(content, keyword) {
			return "main_staple:" + keyword
		}
	}

	proteinKeywords := []string{
		"鸡胸", "鸡腿", "牛肉", "牛腱", "虾", "鱼", "三文鱼", "鸡蛋", "豆腐", "豆干",
	}
	for _, keyword := range proteinKeywords {
		if strings.Contains(content, keyword) {
			return "protein:" + keyword
		}
	}

	cfg := getRecipeTextConfig()
	for _, keyword := range cfg.StapleKeywords {
		normalized := strings.ToLower(keyword)
		if normalized != "" && strings.Contains(content, normalized) {
			return "staple:" + normalized
		}
	}

	if len(recipe.Ingredients) > 0 {
		first := strings.ToLower(strings.TrimSpace(recipe.Ingredients[0]))
		if first != "" {
			return "ingredient:" + first
		}
	}

	return "name:" + strings.ToLower(strings.TrimSpace(recipe.Name))
}

func (s *RecipeServiceImpl) selectDiverseCandidates(recipes []models.Recipe, targetE, targetP, targetC, targetF float64, k int, goal models.HealthGoal) []models.Recipe {
	if len(recipes) == 0 || k <= 0 {
		return nil
	}

	type candidate struct {
		r     models.Recipe
		score float64
		key   string
	}

	cs := make([]candidate, 0, len(recipes))
	for _, recipe := range recipes {
		cs = append(cs, candidate{
			r:     recipe,
			score: s.scoreSingleRecipe(recipe, targetE, targetP, targetC, targetF, goal),
			key:   recipeVarietyKey(recipe),
		})
	}

	sort.Slice(cs, func(i, j int) bool {
		if cs[i].score == cs[j].score {
			return cs[i].r.ID < cs[j].r.ID
		}
		return cs[i].score > cs[j].score
	})

	limit := min(len(cs), k)
	maxPerKey := 2
	if limit <= 6 {
		maxPerKey = 1
	}

	selected := make([]models.Recipe, 0, limit)
	selectedIDs := make(map[uint]struct{}, limit)
	keyCount := make(map[string]int)

	for _, item := range cs {
		if len(selected) >= limit {
			break
		}
		if keyCount[item.key] >= maxPerKey {
			continue
		}
		selected = append(selected, item.r)
		selectedIDs[item.r.ID] = struct{}{}
		keyCount[item.key]++
	}

	for _, item := range cs {
		if len(selected) >= limit {
			break
		}
		if _, exists := selectedIDs[item.r.ID]; exists {
			continue
		}
		selected = append(selected, item.r)
	}

	return selected
}

func normalizedMealIDs(ids []uint, fallback uint) []uint {
	if len(ids) == 0 {
		if fallback == 0 {
			return nil
		}
		return []uint{fallback}
	}

	out := append([]uint(nil), ids...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func sameMealSelection(leftIDs []uint, leftFallback uint, rightIDs []uint, rightFallback uint) bool {
	left := normalizedMealIDs(leftIDs, leftFallback)
	right := normalizedMealIDs(rightIDs, rightFallback)
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return len(left) > 0
}

func mealOverlapCount(leftIDs []uint, leftFallback uint, rightIDs []uint, rightFallback uint) int {
	left := normalizedMealIDs(leftIDs, leftFallback)
	right := normalizedMealIDs(rightIDs, rightFallback)
	if len(left) == 0 || len(right) == 0 {
		return 0
	}

	rightSet := make(map[uint]struct{}, len(right))
	for _, id := range right {
		rightSet[id] = struct{}{}
	}

	count := 0
	for _, id := range left {
		if _, ok := rightSet[id]; ok {
			count++
		}
	}
	return count
}

func (s *RecipeServiceImpl) checkPlanDiversity(existingPlans []*models.DailyRecipePlan, newPlan *models.DailyRecipePlan) bool {
	for _, plan := range existingPlans {
		sameBreakfast := sameMealSelection(plan.BreakfastItemIDs, plan.BreakfastRecipeID, newPlan.BreakfastItemIDs, newPlan.BreakfastRecipeID)
		sameLunch := sameMealSelection(plan.LunchItemIDs, plan.LunchRecipeID, newPlan.LunchItemIDs, newPlan.LunchRecipeID)
		sameDinner := sameMealSelection(plan.DinnerItemIDs, plan.DinnerRecipeID, newPlan.DinnerItemIDs, newPlan.DinnerRecipeID)
		sameSnack := sameMealSelection(plan.SnackItemIDs, plan.SnackRecipeID, newPlan.SnackItemIDs, newPlan.SnackRecipeID)

		if sameLunch && sameDinner {
			return false
		}

		sameCount := 0
		if sameBreakfast {
			sameCount++
		}
		if sameLunch {
			sameCount++
		}
		if sameDinner {
			sameCount++
		}
		if sameSnack {
			sameCount++
		}
		if sameCount >= 3 {
			return false
		}

		lunchOverlap := mealOverlapCount(plan.LunchItemIDs, plan.LunchRecipeID, newPlan.LunchItemIDs, newPlan.LunchRecipeID)
		dinnerOverlap := mealOverlapCount(plan.DinnerItemIDs, plan.DinnerRecipeID, newPlan.DinnerItemIDs, newPlan.DinnerRecipeID)
		if lunchOverlap >= 2 && dinnerOverlap >= 2 {
			return false
		}
	}
	return true
}
