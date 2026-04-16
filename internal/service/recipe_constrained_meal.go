package service

import (
	"NutriPlan/internal/repository/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"NutriPlan/internal/repository/dao"

	"gorm.io/gorm"
)

type ConstrainedMealRegenerateRequest struct {
	MealType        string   `json:"meal_type"`
	Ingredients     []string `json:"ingredients"`
	CoreIngredients []string `json:"core_ingredients,omitempty"`
}

type ConstrainedMealRegenerateResponse struct {
	Meal             GeneratedMeal    `json:"meal"`
	Gap              MealNutritionGap `json:"gap"`
	UsedIngredients  []string         `json:"used_ingredients"`
	SupplementaryTip string           `json:"supplementary_tip,omitempty"`
	RetryCount       int              `json:"retry_count"`
}

type MealNutritionGap struct {
	Energy       float64 `json:"energy"`
	Protein      float64 `json:"protein"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
}

type GeneratedMeal struct {
	MealName                 string           `json:"meal_name"`
	MealType                 string           `json:"meal_type"`
	CoreIngredientsUsed      []string         `json:"core_ingredients_used"`
	SupplementaryIngredients []string         `json:"supplementary_ingredients"`
	NutritionEstimate        MealNutritionGap `json:"nutrition_estimate"`
	Steps                    []string         `json:"steps"`
	DietitianTip             string           `json:"dietitian_tip"`
	FillReason               string           `json:"fill_reason"`
}

type ConstrainedMealAdoptRequest struct {
	MealType string        `json:"meal_type"`
	Meal     GeneratedMeal `json:"meal"`
}

func (s *RecipeServiceImpl) RecognizeIngredientsFromImage(ctx context.Context, imageReader io.Reader) ([]string, error) {
	if s.zhipuClient == nil {
		return nil, errors.New("食材识别服务暂不可用")
	}
	ingredients, err := s.zhipuClient.RecognizeIngredients(ctx, imageReader)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(ingredients))
	seen := make(map[string]struct{})
	for _, item := range ingredients {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}
	return result, nil
}

func (s *RecipeServiceImpl) RegenerateConstrainedMeal(ctx context.Context, user *models.User, req ConstrainedMealRegenerateRequest) (*ConstrainedMealRegenerateResponse, error) {
	if user == nil || user.ID == 0 {
		return nil, errors.New("用户信息无效")
	}
	if s.zhipuClient == nil {
		return nil, errors.New("重构服务暂不可用")
	}

	mealTypeCN, _, err := normalizeMealType(req.MealType)
	if err != nil {
		return nil, err
	}

	ingredients := uniqueIngredients(req.Ingredients)
	if len(ingredients) == 0 {
		return nil, errors.New("请至少提供一种食材")
	}
	coreIngredients := uniqueIngredients(req.CoreIngredients)
	if len(coreIngredients) == 0 {
		coreIngredients = ingredients
	}

	gap, err := s.computeMealGap(user, mealTypeCN)
	if err != nil {
		return nil, err
	}

	prompt := s.buildConstrainedMealPrompt(user, mealTypeCN, coreIngredients, ingredients, gap)

	var parsed GeneratedMeal
	retryCount := 0
	var lastErr error
	for i := 0; i < 3; i++ {
		timeoutCtx, cancel := context.WithTimeout(ctx, 18*time.Second)
		content, genErr := s.zhipuClient.GenerateConstrainedMeal(timeoutCtx, prompt)
		cancel()
		if genErr != nil {
			lastErr = genErr
			retryCount = i + 1
			continue
		}

		if err := unmarshalAndValidateGeneratedMeal(content, coreIngredients, mealTypeCN, &parsed); err != nil {
			lastErr = err
			retryCount = i + 1
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return nil, errors.New("食材组合过于复杂，请尝试减少种类")
	}

	supplementaryTip := ""
	if len(parsed.SupplementaryIngredients) == 0 && (len(ingredients) <= 1 || gap.Energy >= 300) {
		parsed.SupplementaryIngredients = []string{"米饭", "橄榄油"}
		if strings.TrimSpace(parsed.FillReason) == "" {
			parsed.FillReason = "核心食材能量不足，自动补齐基础主食与健康脂肪来源"
		}
		supplementaryTip = "已自动补齐基础辅料以满足热量与宏量营养目标"
	}

	if strings.TrimSpace(parsed.DietitianTip) == "" {
		parsed.DietitianTip = "建议搭配足量饮水，并根据饱腹感微调主食份量。"
	}

	return &ConstrainedMealRegenerateResponse{
		Meal:             parsed,
		Gap:              gap,
		UsedIngredients:  coreIngredients,
		SupplementaryTip: supplementaryTip,
		RetryCount:       retryCount,
	}, nil
}

func (s *RecipeServiceImpl) AdoptRegeneratedMeal(userID uint, req ConstrainedMealAdoptRequest) (*models.DailyRecipePlan, error) {
	if userID == 0 {
		return nil, errors.New("用户ID无效")
	}
	mealTypeCN, mealTypeEN, err := normalizeMealType(req.MealType)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Meal.MealName) == "" {
		return nil, errors.New("重构食谱名称不能为空")
	}
	if req.Meal.NutritionEstimate.Energy <= 0 {
		return nil, errors.New("重构食谱营养值无效")
	}

	if dao.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tx := dao.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var plan models.DailyRecipePlan
	if err := tx.Where("user_id = ? AND is_selected = ? AND plan_date = ?", userID, true, today).
		First(&plan).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("当前暂无已锁定计划，请先锁定当日计划")
	}

	allIngredients := mergeIngredients(req.Meal.CoreIngredientsUsed, req.Meal.SupplementaryIngredients)
	if len(allIngredients) == 0 {
		allIngredients = []string{"适量调味"}
	}
	steps := req.Meal.Steps
	if len(steps) == 0 {
		steps = []string{"按推荐食材与份量烹饪即可"}
	}

	newRecipe := models.Recipe{
		Name:           req.Meal.MealName,
		MealType:       mealTypeCN,
		Difficulty:     "中等",
		CookingTime:    25,
		PortionWeightG: 100,
		Energy:         req.Meal.NutritionEstimate.Energy,
		Protein:        req.Meal.NutritionEstimate.Protein,
		Carbohydrate:   req.Meal.NutritionEstimate.Carbohydrate,
		Fat:            req.Meal.NutritionEstimate.Fat,
		Ingredients:    allIngredients,
		CookingSteps:   steps,
	}
	if err := tx.Create(&newRecipe).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存重构食谱失败: %w", err)
	}

	switch mealTypeCN {
	case models.MealTypeBreakfast:
		plan.BreakfastRecipeID = newRecipe.ID
		plan.BreakfastItemIDs = []uint{newRecipe.ID}
	case models.MealTypeLunch:
		plan.LunchRecipeID = newRecipe.ID
		plan.LunchItemIDs = []uint{newRecipe.ID}
	case models.MealTypeDinner:
		plan.DinnerRecipeID = newRecipe.ID
		plan.DinnerItemIDs = []uint{newRecipe.ID}
	case models.MealTypeSnack:
		plan.SnackRecipeID = newRecipe.ID
		plan.SnackItemIDs = []uint{newRecipe.ID}
	}

	if err := s.recalculatePlanNutrition(tx, &plan); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(&plan).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新计划失败: %w", err)
	}

	if err := tx.Where("user_id = ? AND record_date = ? AND meal_type IN ?", userID, today, []string{mealTypeEN, string(mealTypeCN)}).
		Delete(&models.DailyIntakeRecord{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("清理旧餐次记录失败: %w", err)
	}

	newRecord := models.DailyIntakeRecord{
		UserID:            userID,
		RecordDate:        today,
		MealType:          mealTypeEN,
		FoodSource:        1,
		SourceID:          newRecipe.ID,
		FoodName:          newRecipe.Name,
		IntakeAmount:      100,
		CalculatedEnergy:  newRecipe.Energy,
		CalculatedProtein: newRecipe.Protein,
		CalculatedCarb:    newRecipe.Carbohydrate,
		CalculatedFat:     newRecipe.Fat,
	}
	if err := tx.Create(&newRecord).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("写入新餐次记录失败: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.GetSelectedPlan(userID)
}

func (s *RecipeServiceImpl) computeMealGap(user *models.User, requestedMeal models.MealType) (MealNutritionGap, error) {
	if user.TDEE <= 0 {
		return MealNutritionGap{}, errors.New("用户档案不完整，请先完善健康档案")
	}

	targetCalorie := s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
	targetProtein, targetCarb, targetFat := s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)

	otherTotals := MealNutritionGap{}
	ingestedByMeal, err := queryTodayMealIntake(user.ID)
	if err != nil {
		return MealNutritionGap{}, err
	}

	selectedPlan, _ := s.GetSelectedPlan(user.ID)
	for _, meal := range []models.MealType{models.MealTypeBreakfast, models.MealTypeLunch, models.MealTypeDinner, models.MealTypeSnack} {
		if meal == requestedMeal {
			continue
		}
		_, intakeKey, _ := normalizeMealType(string(meal))
		consumed := ingestedByMeal[intakeKey]
		if consumed.Energy > 0 {
			otherTotals = sumGap(otherTotals, consumed)
			continue
		}
		if selectedPlan != nil {
			locked := mealNutritionFromPlan(selectedPlan, meal)
			otherTotals = sumGap(otherTotals, locked)
		}
	}

	gap := MealNutritionGap{
		Energy:       targetCalorie - otherTotals.Energy,
		Protein:      targetProtein - otherTotals.Protein,
		Carbohydrate: targetCarb - otherTotals.Carbohydrate,
		Fat:          targetFat - otherTotals.Fat,
	}

	if gap.Energy < 120 {
		gap.Energy = 120
	}
	if gap.Protein < 8 {
		gap.Protein = 8
	}
	if gap.Carbohydrate < 10 {
		gap.Carbohydrate = 10
	}
	if gap.Fat < 3 {
		gap.Fat = 3
	}

	return gap, nil
}

func (s *RecipeServiceImpl) buildConstrainedMealPrompt(user *models.User, mealType models.MealType, coreIngredients, allIngredients []string, gap MealNutritionGap) string {
	return fmt.Sprintf(`你是一名临床营养师兼家庭厨师，请严格输出 JSON。
目标餐次：%s
核心必须使用食材：%s
可使用全部食材：%s
当前餐次营养缺口（请尽量贴近）：
- 热量：%.1f kcal
- 蛋白质：%.1f g
- 碳水：%.1f g
- 脂肪：%.1f g
用户目标：%s
要求：
1) 必须包含核心食材；
2) 若现有食材无法满足营养目标，可补齐基础辅料；
3) 返回结构必须是 JSON，不要 markdown。
JSON 结构：
{
  "meal_name":"",
  "meal_type":"",
  "core_ingredients_used":[],
  "supplementary_ingredients":[],
  "nutrition_estimate":{"energy":0,"protein":0,"carbohydrate":0,"fat":0},
  "steps":[],
  "dietitian_tip":"",
  "fill_reason":""
}
`, mealType, strings.Join(coreIngredients, "、"), strings.Join(allIngredients, "、"), gap.Energy, gap.Protein, gap.Carbohydrate, gap.Fat, user.HealthGoal)
}

func unmarshalAndValidateGeneratedMeal(content string, coreIngredients []string, mealType models.MealType, out *GeneratedMeal) error {
	jsonBody := extractJSONBody(content)
	if jsonBody == "" {
		return errors.New("模型返回非JSON")
	}
	if err := json.Unmarshal([]byte(jsonBody), out); err != nil {
		return err
	}
	if strings.TrimSpace(out.MealName) == "" || len(out.Steps) == 0 {
		return errors.New("关键字段缺失")
	}
	if out.NutritionEstimate.Energy <= 0 || out.NutritionEstimate.Protein <= 0 || out.NutritionEstimate.Carbohydrate <= 0 {
		return errors.New("营养估算字段无效")
	}
	if strings.TrimSpace(out.MealType) == "" {
		out.MealType = string(mealType)
	}

	usedSet := make(map[string]struct{})
	for _, item := range out.CoreIngredientsUsed {
		usedSet[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}
	for _, core := range coreIngredients {
		key := strings.ToLower(strings.TrimSpace(core))
		if key == "" {
			continue
		}
		if _, ok := usedSet[key]; !ok {
			return errors.New("核心食材未被使用")
		}
	}
	return nil
}

func extractJSONBody(raw string) string {
	trimmed := strings.TrimSpace(raw)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		return trimmed[start : end+1]
	}
	return ""
}

func normalizeMealType(value string) (models.MealType, string, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "早餐", "breakfast":
		return models.MealTypeBreakfast, "breakfast", nil
	case "午餐", "lunch":
		return models.MealTypeLunch, "lunch", nil
	case "晚餐", "dinner":
		return models.MealTypeDinner, "dinner", nil
	case "加餐", "snack":
		return models.MealTypeSnack, "snack", nil
	default:
		return "", "", errors.New("meal_type 仅支持 breakfast/lunch/dinner/snack")
	}
}

func uniqueIngredients(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{})
	for _, item := range values {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, v)
	}
	return result
}

func mergeIngredients(a, b []string) []string {
	merged := make([]string, 0, len(a)+len(b))
	merged = append(merged, a...)
	merged = append(merged, b...)
	return uniqueIngredients(merged)
}

func sumGap(a, b MealNutritionGap) MealNutritionGap {
	return MealNutritionGap{
		Energy:       a.Energy + b.Energy,
		Protein:      a.Protein + b.Protein,
		Carbohydrate: a.Carbohydrate + b.Carbohydrate,
		Fat:          a.Fat + b.Fat,
	}
}

func queryTodayMealIntake(userID uint) (map[string]MealNutritionGap, error) {
	if dao.DB == nil {
		return nil, errors.New("数据库未初始化")
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var records []models.DailyIntakeRecord
	if err := dao.DB.Where("user_id = ? AND record_date >= ? AND record_date < ?", userID, start, end).Find(&records).Error; err != nil {
		return nil, err
	}

	result := map[string]MealNutritionGap{
		"breakfast": {},
		"lunch":     {},
		"dinner":    {},
		"snack":     {},
	}
	for _, record := range records {
		_, key, err := normalizeMealType(record.MealType)
		if err != nil {
			continue
		}
		current := result[key]
		current.Energy += record.CalculatedEnergy
		current.Protein += record.CalculatedProtein
		current.Carbohydrate += record.CalculatedCarb
		current.Fat += record.CalculatedFat
		result[key] = current
	}
	return result, nil
}

func mealNutritionFromPlan(plan *models.DailyRecipePlan, mealType models.MealType) MealNutritionGap {
	if plan == nil {
		return MealNutritionGap{}
	}
	items := make([]models.Recipe, 0)
	switch mealType {
	case models.MealTypeBreakfast:
		items = plan.BreakfastItems
		if len(items) == 0 && plan.BreakfastRecipe.ID > 0 {
			items = []models.Recipe{plan.BreakfastRecipe}
		}
	case models.MealTypeLunch:
		items = plan.LunchItems
		if len(items) == 0 && plan.LunchRecipe.ID > 0 {
			items = []models.Recipe{plan.LunchRecipe}
		}
	case models.MealTypeDinner:
		items = plan.DinnerItems
		if len(items) == 0 && plan.DinnerRecipe.ID > 0 {
			items = []models.Recipe{plan.DinnerRecipe}
		}
	case models.MealTypeSnack:
		items = plan.SnackItems
		if len(items) == 0 && plan.SnackRecipe.ID > 0 {
			items = []models.Recipe{plan.SnackRecipe}
		}
	}
	result := MealNutritionGap{}
	for _, item := range items {
		result.Energy += item.Energy
		result.Protein += item.Protein
		result.Carbohydrate += item.Carbohydrate
		result.Fat += item.Fat
	}
	return result
}

func (s *RecipeServiceImpl) recalculatePlanNutrition(db *gorm.DB, plan *models.DailyRecipePlan) error {
	type mealIDs struct {
		ids      []uint
		fallback uint
	}
	groups := []mealIDs{
		{ids: plan.BreakfastItemIDs, fallback: plan.BreakfastRecipeID},
		{ids: plan.LunchItemIDs, fallback: plan.LunchRecipeID},
		{ids: plan.DinnerItemIDs, fallback: plan.DinnerRecipeID},
		{ids: plan.SnackItemIDs, fallback: plan.SnackRecipeID},
	}

	total := MealNutritionGap{}
	for _, group := range groups {
		ids := group.ids
		if len(ids) == 0 && group.fallback > 0 {
			ids = []uint{group.fallback}
		}
		if len(ids) == 0 {
			continue
		}

		var recipes []models.Recipe
		if err := db.Model(&models.Recipe{}).Where("id IN ?", ids).Find(&recipes).Error; err != nil {
			return err
		}
		recipeMap := make(map[uint]models.Recipe)
		for _, recipe := range recipes {
			recipeMap[recipe.ID] = recipe
		}
		for _, id := range ids {
			recipe, ok := recipeMap[id]
			if !ok {
				continue
			}
			scaled := scaleRecipeNutritionByPortion(recipe)
			total.Energy += scaled.Energy
			total.Protein += scaled.Protein
			total.Carbohydrate += scaled.Carbohydrate
			total.Fat += scaled.Fat
		}
	}

	plan.TotalEnergy = total.Energy
	plan.TotalProtein = total.Protein
	plan.TotalCarbohydrate = total.Carbohydrate
	plan.TotalFat = total.Fat
	return nil
}
