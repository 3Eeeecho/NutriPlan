package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const ()

// RecipeService 食谱推荐服务接口
type RecipeService interface {
	// RecommendRecipes 为用户推荐每日食谱计划
	// count: 推荐方案数量（默认3-5套）
	RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error)

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

// 目前的推荐算法逻辑如下：

// 1. 目标设定
// 根据用户的档案（年龄、体重、目标等）计算出每日目标热量 (TDEE) 和 三大营养素目标（蛋白质、碳水、脂肪）。

// 2. 食谱初筛 (Candidate Selection)
// 	禁忌过滤：剔除过敏食材。
// 	历史去重：剔除最近 7 天使用过的食谱（之前是 2 天），避免近期重复。
// 	候选池选取：
// 	对早/午/晚/加餐，分别找出热量最接近目标的 45道菜 (3倍候选数量)。
// 	随机洗牌：从中随机抽取 15 道进入下一轮计算。
// 	优化点：这是为了解决“每次都选完全一样的那几道菜”的问题。
// 3. 生成组合 (Combination)
// 	采用笛卡尔积（Breakfast × Lunch × Dinner × Snack）生成上万种可能的“一日食谱组合”。
// 4. 评分逻辑 (Scoring)
// 	计算每种组合的总热量和宏量营养素。
// 	计算其与目标值的偏差（Deviation）。
// 	真实评分：根据偏差计算 0-100 的分数。
// 	无同情分：您刚才删除了 <60 分强行补到 60 分的逻辑，现在返回的是真实反映差距的分数。
// 5. 排序与优选 (Ranking & Selection)
// 	将所有方案按分数从高到低排序。
// 	加权随机：取前 100 名，进行加权随机抽取（分数越高的越容易被选中，但不是绝对选第一名）。
// 	硬门槛拦截（您刚才修改的）：
// 	代码行 if plan.MatchScore < 80 { break }
// 	如果不满 80 分，直接丢弃，不再推荐。
// 	多样性检查：如果新选出来的方案和已选方案在“午餐”和“晚餐”上由于过度重合，会被跳过。

// RecommendRecipes 为用户推荐每日食谱计划
func (s *RecipeServiceImpl) RecommendRecipes(user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	// 计算用户目标营养需求
	targetCalorie := s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
	targetProtein, targetCarb, targetFat := s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)

	// 获取所有餐次过滤掉禁忌食材后的食谱
	forbiddenIngredients := s.parseForbiddenIngredients(user)
	breakfastRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeBreakfast, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	lunchRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeLunch, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	dinnerRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeDinner, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	snackRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeSnack, forbiddenIngredients)
	if err != nil {
		return nil, err
	}

	// 过滤最近7天吃过的食谱 (扩大去重范围，避免近期重复)
	recentRecipeIDs, err := s.recipeRepo.FindRecentPlanRecipeIDs(user.ID, 7)
	if err == nil && len(recentRecipeIDs) > 0 {
		breakfastRecipes = s.filterRecentRecipes(breakfastRecipes, recentRecipeIDs)
		lunchRecipes = s.filterRecentRecipes(lunchRecipes, recentRecipeIDs)
		dinnerRecipes = s.filterRecentRecipes(dinnerRecipes, recentRecipeIDs)
		snackRecipes = s.filterRecentRecipes(snackRecipes, recentRecipeIDs)
	}

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

	// 从计划里随机抽取,确保每次刷新结果不一致
	selectedPlans := s.shuffleByWeights(rankedPlans)

	// 多样性过滤与结果截取
	finalPlans := make([]*models.DailyRecipePlan, 0)

	for _, plan := range selectedPlans {
		// 匹配度门槛
		if plan.MatchScore < 80 {
			break // 因为已经排好序了，后面更低，直接断开
		}

		// 多样性检查 (核心)
		// 确保新加入的 plan 和已经选中的 plan 不吃重复的主菜
		if !s.checkDiversity(finalPlans, plan) {
			continue
		}

		// 入库并保存
		if err := s.recipeRepo.CreateDailyPlan(plan); err != nil {
			continue
		}
		finalPlans = append(finalPlans, plan)

		if len(finalPlans) >= count {
			break
		}
	}

	return finalPlans, nil
}

// NutritionTarget 营养目标
type NutritionTarget struct {
	Energy       float64 // 热量
	Protein      float64 // 蛋白质
	Carbohydrate float64 // 碳水化合物
	Fat          float64 // 脂肪
}

// generateOneDailyPlan 生成一套每日食谱计划
func (s *RecipeServiceImpl) generateRankedPlans(
	userID uint,
	breakfastRecipes, lunchRecipes, dinnerRecipes, snackRecipes []models.Recipe,
	target NutritionTarget, userGoal models.HealthGoal,
) []*models.DailyRecipePlan {
	// 初筛：每餐选出最接近单餐热量目标的候选
	// 这样可以避免后面组合爆炸，同时保证候选质量
	cBreakfast := s.getTopCandidates(breakfastRecipes, target.Energy*0.25, 15)
	cLunch := s.getTopCandidates(lunchRecipes, target.Energy*0.35, 15)
	cDinner := s.getTopCandidates(dinnerRecipes, target.Energy*0.30, 15)

	// 加餐可选，没有就塞一个空食谱进去方便循环
	cSnack := s.getTopCandidates(snackRecipes, target.Energy*0.10, 10)
	if len(cSnack) == 0 {
		cSnack = append(cSnack, models.Recipe{Name: ""}) // 空占位
	}

	if len(cBreakfast) == 0 || len(cLunch) == 0 || len(cDinner) == 0 {
		return nil // 无法生成
	}

	allPlans := make([]*models.DailyRecipePlan, 0, 2560)

	// 笛卡尔积循环：遍历所有组合
	for _, b := range cBreakfast {
		for _, l := range cLunch {
			for _, d := range cDinner {
				for _, sn := range cSnack {

					// 计算当前组合的总营养
					totalE := b.Energy + l.Energy + d.Energy + sn.Energy
					totalP := b.Protein + l.Protein + d.Protein + sn.Protein
					totalC := b.Carbohydrate + l.Carbohydrate + d.Carbohydrate + sn.Carbohydrate
					totalF := b.Fat + l.Fat + d.Fat + sn.Fat

					matchScore := s.calculateMatchScore(
						totalE, totalP, totalC, totalF,
						target.Energy, target.Protein, target.Carbohydrate, target.Fat,
						userGoal,
					)

					plan := &models.DailyRecipePlan{
						UserID:            userID,
						BreakfastRecipe:   b,
						BreakfastRecipeID: b.ID,
						LunchRecipe:       l,
						LunchRecipeID:     l.ID,
						DinnerRecipe:      d,
						DinnerRecipeID:    d.ID,
						SnackRecipe:       sn, // 如果是空占位，name为空
						SnackRecipeID:     sn.ID,

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

	return allPlans
}

// shuffleByWeights 对前 N 个方案进行加权洗牌
func (s *RecipeServiceImpl) shuffleByWeights(candidates []*models.DailyRecipePlan) []*models.DailyRecipePlan {
	if len(candidates) == 0 {
		return nil
	}

	// 1. 只取前 100 名 (或者全部)，保证质量
	// 剩下的分数太低，就不参与随机了，避免推荐烂方案
	topN := min(len(candidates), 100)
	pool := candidates[:topN]

	// 如果池子很小，直接返回，没必要随机
	if len(pool) < 3 {
		return pool
	}

	// --- 更稳健的实现：加权构建法 ---
	// 创建一个临时列表用于抽取
	source := make([]*models.DailyRecipePlan, len(pool))
	copy(source, pool)
	result := make([]*models.DailyRecipePlan, 0, len(pool))

	for len(source) > 0 {
		// 计算总权重
		totalWeight := 0.0
		for _, p := range source {
			// 权重放大：分数^3，让高分优势更明显
			totalWeight += math.Pow(p.MatchScore, 3)
		}

		r := rand.Float64() * totalWeight
		curr := 0.0
		foundIdx := -1

		for i, p := range source {
			curr += math.Pow(p.MatchScore, 3)
			if r <= curr {
				foundIdx = i
				break
			}
		}

		if foundIdx == -1 {
			foundIdx = len(source) - 1
		} // 防止精度误差

		// 放入结果集
		result = append(result, source[foundIdx])

		// 从源列表中删除已选的 (避免重复)
		source = append(source[:foundIdx], source[foundIdx+1:]...)
	}

	return result
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

// 获取 Top-K 候选食谱（只基于热量粗筛，为了减少计算量）
func (s *RecipeServiceImpl) getTopCandidates(recipes []models.Recipe, targetEnergy float64, k int) []models.Recipe {
	if len(recipes) == 0 {
		return nil
	}

	// 简单算一下热量差，排序
	type candidate struct {
		r    models.Recipe
		diff float64
	}

	var cs []candidate
	for _, r := range recipes {
		diff := math.Abs(r.Energy - targetEnergy)
		cs = append(cs, candidate{r, diff})
	}

	// 排序
	sort.Slice(cs, func(i, j int) bool { return cs[i].diff < cs[j].diff })

	// 不直接取前 k 个，而是取前 3*k 个作为“候选池”，然后从中随机选 k 个
	// 这样可以避免每次都选出完全一样的“最优解”
	poolSize := min(k*3, len(cs))

	pool := cs[:poolSize]

	// 洗牌 pool
	s.rng.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	// 取前 k 个
	limit := min(len(pool), k)

	result := make([]models.Recipe, limit)
	for i := range limit {
		result[i] = pool[i].r
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
