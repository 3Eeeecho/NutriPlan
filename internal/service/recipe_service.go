package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// RecipeService 食谱推荐服务接口
type RecipeService interface {
	// RecommendRecipes 为用户推荐每日食谱计划
	// count: 推荐方案数量（默认3-5套）
	RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error)

	// GetRecipesByMealType 根据餐次类型获取食谱
	GetRecipesByMealType(mealType models.MealType) ([]models.Recipe, error)

	// SaveDailyPlan 保存每日食谱计划
	SaveDailyPlan(plan *models.DailyRecipePlan) error

	// GetSelectedPlan 获取用户当前选中的计划
	GetSelectedPlan(userID uint) (*models.DailyRecipePlan, error)

	// SelectDailyPlan 选择每日食谱计划
	SelectDailyPlan(userID, planID uint) error
}

// RecipeServiceImpl 食谱推荐服务实现
type RecipeServiceImpl struct {
	recipeRepo   dao.RecipeRepository
	nutriService NutriService
	rng          *rand.Rand
}

// NewRecipeService 创建食谱推荐服务实例
func NewRecipeService(recipeRepo dao.RecipeRepository, nutriService NutriService) RecipeService {
	return &RecipeServiceImpl{
		recipeRepo:   recipeRepo,
		nutriService: nutriService,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RecommendRecipes 为用户推荐每日食谱计划
func (s *RecipeServiceImpl) RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	// 1. 计算用户目标营养需求
	targetCalorie := s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
	targetProtein, targetCarb, targetFat := s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)

	// 2. 获取所有餐次的食谱
	breakfastRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeBreakfast)
	if err != nil {
		return nil, err
	}
	lunchRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeLunch)
	if err != nil {
		return nil, err
	}
	dinnerRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeDinner)
	if err != nil {
		return nil, err
	}
	snackRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeSnack)
	if err != nil {
		return nil, err
	}

	// 3. 过滤禁忌食材
	forbiddenIngredients := s.parseForbiddenIngredients(user)
	breakfastRecipes = s.filterForbiddenRecipes(breakfastRecipes, forbiddenIngredients)
	lunchRecipes = s.filterForbiddenRecipes(lunchRecipes, forbiddenIngredients)
	dinnerRecipes = s.filterForbiddenRecipes(dinnerRecipes, forbiddenIngredients)
	snackRecipes = s.filterForbiddenRecipes(snackRecipes, forbiddenIngredients)

	// 4. 过滤最近3天吃过的食谱 (History-based Anti-Repetition)
	recentRecipeIDs, err := s.recipeRepo.FindRecentPlanRecipeIDs(user.ID, 3)
	if err == nil && len(recentRecipeIDs) > 0 {
		breakfastRecipes = s.filterRecentRecipes(breakfastRecipes, recentRecipeIDs)
		lunchRecipes = s.filterRecentRecipes(lunchRecipes, recentRecipeIDs)
		dinnerRecipes = s.filterRecentRecipes(dinnerRecipes, recentRecipeIDs)
		snackRecipes = s.filterRecentRecipes(snackRecipes, recentRecipeIDs)
	}

	// 5. 生成每日食谱计划
	plans := make([]*models.DailyRecipePlan, 0)
	targetNutrition := NutritionTarget{
		Energy:       targetCalorie,
		Protein:      targetProtein,
		Carbohydrate: targetCarb,
		Fat:          targetFat,
	}

	// 尝试生成多套方案
	maxAttempts := count * 20
	for i := 0; i < maxAttempts && len(plans) < count; i++ {
		plan := s.generateOneDailyPlan(
			user.ID,
			breakfastRecipes,
			lunchRecipes,
			dinnerRecipes,
			snackRecipes,
			targetNutrition,
		)

		// 只保留匹配度>=80%的方案
		if plan != nil && plan.MatchScore >= 80 {
			// 检查是否与已有方案重复 (完全重复)
			if s.isDuplicatePlan(plans, plan) {
				continue
			}

			// 多样性检查 (Diversity Filtering): 确保同一批推荐中主菜不重复
			if !s.checkDiversity(plans, plan) {
				continue
			}

			// 保存方案到数据库，以获取ID
			if err := s.recipeRepo.CreateDailyPlan(plan); err != nil {
				continue
			}
			plans = append(plans, plan)
		}
	}

	return plans, nil
}

// NutritionTarget 营养目标
type NutritionTarget struct {
	Energy       float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
}

// generateOneDailyPlan 生成一套每日食谱计划
func (s *RecipeServiceImpl) generateOneDailyPlan(
	userID uint,
	breakfastRecipes, lunchRecipes, dinnerRecipes, snackRecipes []models.Recipe,
	target NutritionTarget,
) *models.DailyRecipePlan {
	if len(breakfastRecipes) == 0 || len(lunchRecipes) == 0 || len(dinnerRecipes) == 0 {
		return nil
	}

	// 随机选择各餐次的食谱
	// TODO: 实际应用中应使用更智能的选择算法
	breakfast := s.selectBestRecipe(breakfastRecipes, target.Energy*0.25)
	lunch := s.selectBestRecipe(lunchRecipes, target.Energy*0.35)
	dinner := s.selectBestRecipe(dinnerRecipes, target.Energy*0.30)

	var snack models.Recipe
	if len(snackRecipes) > 0 {
		selectedSnack := s.selectBestRecipe(snackRecipes, target.Energy*0.10)
		snack = selectedSnack
	}

	// 计算总营养
	totalEnergy := breakfast.Energy + lunch.Energy + dinner.Energy
	totalProtein := breakfast.Protein + lunch.Protein + dinner.Protein
	totalCarb := breakfast.Carbohydrate + lunch.Carbohydrate + dinner.Carbohydrate
	totalFat := breakfast.Fat + lunch.Fat + dinner.Fat

	snackRecipeID := snack.ID
	totalEnergy += snack.Energy
	totalProtein += snack.Protein
	totalCarb += snack.Carbohydrate
	totalFat += snack.Fat

	// 计算匹配度
	matchScore := s.calculateMatchScore(
		totalEnergy, totalProtein, totalCarb, totalFat,
		target.Energy, target.Protein, target.Carbohydrate, target.Fat,
	)

	plan := &models.DailyRecipePlan{
		UserID:             userID,
		BreakfastRecipeID:  breakfast.ID,
		BreakfastRecipe:    breakfast, // Populate full recipe
		LunchRecipeID:      lunch.ID,
		LunchRecipe:        lunch, // Populate full recipe
		DinnerRecipeID:     dinner.ID,
		DinnerRecipe:       dinner, // Populate full recipe
		SnackRecipeID:      snackRecipeID,
		SnackRecipe:        snack,
		TotalEnergy:        totalEnergy,
		TotalProtein:       totalProtein,
		TotalCarbohydrate:  totalCarb,
		TotalFat:           totalFat,
		TargetEnergy:       target.Energy,
		TargetProtein:      target.Protein,
		TargetCarbohydrate: target.Carbohydrate,
		TargetFat:          target.Fat,
		MatchScore:         matchScore,
		IsSelected:         false,
		PlanDate:           time.Now(),
	}

	return plan
}

// selectBestRecipe 从食谱列表中选择最佳食谱（Top-K 随机策略）
func (s *RecipeServiceImpl) selectBestRecipe(recipes []models.Recipe, targetEnergy float64) models.Recipe {
	if len(recipes) == 0 {
		return models.Recipe{}
	}

	type candidate struct {
		recipe models.Recipe
		score  float64
	}

	candidates := make([]candidate, 0, len(recipes))

	for _, recipe := range recipes {
		// 计算能量接近度
		energyMatch := 1 - math.Abs(recipe.Energy-targetEnergy)/targetEnergy
		if energyMatch < 0 {
			energyMatch = 0
		}

		candidates = append(candidates, candidate{recipe: recipe, score: energyMatch})
	}

	// 按分数降序排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// 从 Top-K 中随机选择
	k := min(5, len(candidates))

	// 随机选择一个
	idx := s.rng.Intn(k)
	return candidates[idx].recipe
}

// calculateMatchScore 计算营养匹配度（0-100）
func (s *RecipeServiceImpl) calculateMatchScore(
	actualE, actualP, actualC, actualF,
	targetE, targetP, targetC, targetF float64,
) float64 {
	// 避免除以零
	if targetE == 0 || targetP == 0 || targetC == 0 || targetF == 0 {
		return 0
	}

	// 计算各营养素的匹配度
	energyMatch := 1 - math.Abs(actualE-targetE)/targetE
	proteinMatch := 1 - math.Abs(actualP-targetP)/targetP
	carbMatch := 1 - math.Abs(actualC-targetC)/targetC
	fatMatch := 1 - math.Abs(actualF-targetF)/targetF

	// 确保匹配度在0-1之间
	energyMatch = math.Max(0, math.Min(1, energyMatch))
	proteinMatch = math.Max(0, math.Min(1, proteinMatch))
	carbMatch = math.Max(0, math.Min(1, carbMatch))
	fatMatch = math.Max(0, math.Min(1, fatMatch))

	// 加权计算总匹配度（能量占40%，其他各占20%）
	totalMatch := energyMatch*0.4 + proteinMatch*0.2 + carbMatch*0.2 + fatMatch*0.2

	// 转换为百分比
	return math.Round(totalMatch * 100)
}

// parseForbiddenIngredients 解析用户禁忌食材
func (s *RecipeServiceImpl) parseForbiddenIngredients(user *models.User) []string {
	forbidden := make([]string, 0)

	// 解析过敏源
	if user.Allergies != "" {
		allergies := strings.Split(user.Allergies, ",")
		for _, a := range allergies {
			forbidden = append(forbidden, strings.TrimSpace(a))
		}
	}

	return forbidden
}

// filterForbiddenRecipes 过滤包含禁忌食材的食谱
func (s *RecipeServiceImpl) filterForbiddenRecipes(recipes []models.Recipe, forbidden []string) []models.Recipe {
	if len(forbidden) == 0 {
		return recipes
	}

	filtered := make([]models.Recipe, 0)
	for _, recipe := range recipes {
		// 解析食谱的食材列表
		var ingredients []string
		if err := json.Unmarshal([]byte(recipe.Ingredients), &ingredients); err != nil {
			// 如果解析失败，尝试按逗号分隔
			ingredients = strings.Split(recipe.Ingredients, ",")
		}

		// 检查是否包含禁忌食材
		hasForbidden := false
		for _, ingredient := range ingredients {
			ingredientLower := strings.ToLower(strings.TrimSpace(ingredient))
			for _, f := range forbidden {
				if strings.Contains(ingredientLower, strings.ToLower(f)) {
					hasForbidden = true
					break
				}
			}
			if hasForbidden {
				break
			}
		}

		if !hasForbidden {
			filtered = append(filtered, recipe)
		}
	}

	return filtered
}

// isDuplicatePlan 检查是否为重复方案
func (s *RecipeServiceImpl) isDuplicatePlan(plans []*models.DailyRecipePlan, newPlan *models.DailyRecipePlan) bool {
	for _, plan := range plans {
		if plan.BreakfastRecipeID == newPlan.BreakfastRecipeID &&
			plan.LunchRecipeID == newPlan.LunchRecipeID &&
			plan.DinnerRecipeID == newPlan.DinnerRecipeID {
			return true
		}
	}
	return false
}

// GetRecipesByMealType 根据餐次类型获取食谱
func (s *RecipeServiceImpl) GetRecipesByMealType(mealType models.MealType) ([]models.Recipe, error) {
	return s.recipeRepo.FindByMealType(mealType)
}

// SaveDailyPlan 保存每日食谱计划
func (s *RecipeServiceImpl) SaveDailyPlan(plan *models.DailyRecipePlan) error {
	return s.recipeRepo.CreateDailyPlan(plan)
}

// GetSelectedPlan 获取用户当前选中的计划
func (s *RecipeServiceImpl) GetSelectedPlan(userID uint) (*models.DailyRecipePlan, error) {
	return s.recipeRepo.FindSelectedPlan(userID)
}

// SelectDailyPlan 选择每日食谱计划
func (s *RecipeServiceImpl) SelectDailyPlan(userID, planID uint) error {
	// 1. 获取用户当前已选的计划（如果有）
	currentSelected, err := s.recipeRepo.FindSelectedPlan(userID)
	if err == nil && currentSelected != nil {
		// 取消选中
		if err := s.recipeRepo.UpdatePlanSelection(currentSelected.ID, false); err != nil {
			return err
		}
	}

	// 2. 选中新计划
	return s.recipeRepo.UpdatePlanSelection(planID, true)
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

// checkDiversity 检查方案多样性（确保主菜不重复）
func (s *RecipeServiceImpl) checkDiversity(existingPlans []*models.DailyRecipePlan, newPlan *models.DailyRecipePlan) bool {
	for _, plan := range existingPlans {
		// 检查早餐是否重复
		if plan.BreakfastRecipeID == newPlan.BreakfastRecipeID {
			return false
		}
		// 检查午餐是否重复
		if plan.LunchRecipeID == newPlan.LunchRecipeID {
			return false
		}
		// 检查晚餐是否重复
		if plan.DinnerRecipeID == newPlan.DinnerRecipeID {
			return false
		}
	}
	return true
}
