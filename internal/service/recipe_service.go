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

	// SelectDailyPlan 选择每日食谱计划并保存
	SelectDailyPlan(userID, planID uint) error

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

	// 解析用户偏好与禁忌（过敏 / 健康状况）
	userTags, forbiddenIngredients := s.parseUserPreferences(user)

	// 获取所有餐次，后续再基于 TargetUsers/ForbiddenUsers 做精筛
	breakfastRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeBreakfast, user.HealthGoal, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	lunchRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeLunch, user.HealthGoal, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	dinnerRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeDinner, user.HealthGoal, forbiddenIngredients)
	if err != nil {
		return nil, err
	}
	snackRecipes, err := s.recipeRepo.FindByMealType(models.MealTypeSnack, user.HealthGoal, forbiddenIngredients)
	if err != nil {
		return nil, err
	}

	// 基于 Recipe.TargetUsers / Recipe.ForbiddenUsers 做精确过滤
	breakfastRecipes = s.filterRecipesByUserPreferences(breakfastRecipes, userTags)
	lunchRecipes = s.filterRecipesByUserPreferences(lunchRecipes, userTags)
	dinnerRecipes = s.filterRecipesByUserPreferences(dinnerRecipes, userTags)
	snackRecipes = s.filterRecipesByUserPreferences(snackRecipes, userTags)

	// 过滤最近5天吃过的食谱 (扩大去重范围，避免近期重复)
	recentRecipeIDs, err := s.recipeRepo.FindRecentPlanRecipeIDs(user.ID, 5)
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
	// 获取用户当前已选的计划（如果有）
	currentSelected, err := s.recipeRepo.FindSelectedPlan(userID)
	if err == nil && currentSelected != nil {
		// 取消选中
		if err := s.recipeRepo.UpdatePlanSelection(currentSelected.ID, false); err != nil {
			return err
		}
	}

	// 选中新计划
	return s.recipeRepo.UpdatePlanSelection(planID, true)
}

// GetRecipeDetail 获取食谱详情（含收藏状态）
func (s *RecipeServiceImpl) GetRecipeDetail(recipeID uint, userID uint) (*models.Recipe, bool, error) {
	// 查询食谱
	recipe, err := s.recipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, false, err
	}

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

	tagSet := make(map[string]bool)
	for _, t := range userTags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t != "" {
			tagSet[t] = true
		}
	}

	filtered := make([]models.Recipe, 0, len(recipes))
	for _, r := range recipes {
		// 如果菜谱标注了 ForbiddenUsers，且与用户标签有交集，则排除
		skip := false
		for _, fu := range r.ForbiddenUsers {
			if fu == "" {
				continue
			}
			if tagSet[strings.TrimSpace(fu)] {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// 如果菜谱有 TargetUsers，则要求与用户标签有至少一个交集；否则认为通用
		if len(r.TargetUsers) > 0 {
			matched := false
			for _, tu := range r.TargetUsers {
				if tu == "" {
					continue
				}
				if tagSet[strings.TrimSpace(tu)] {
					matched = true
					break
				}
			}
			if !matched {
				// 如果没有匹配的标签，跳过这个菜谱
				continue
			}
		}

		filtered = append(filtered, r)
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
