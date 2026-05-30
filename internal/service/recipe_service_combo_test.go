package service

import (
	"NutriPlan/internal/repository/models"
	"context"
	"math/rand"
	"testing"

	"gorm.io/gorm"
)

type fakeCollaborativeFilter map[models.MealType][]RecipeScore

func (f fakeCollaborativeFilter) RecommendRecipes(ctx context.Context, userID uint, mealType models.MealType, limit int) ([]RecipeScore, error) {
	scores := f[mealType]
	if len(scores) > limit {
		return scores[:limit], nil
	}
	return scores, nil
}

func assertItemCountInRange(t *testing.T, meal string, items []models.Recipe, minItems, maxItems int) {
	t.Helper()
	if len(items) < minItems || len(items) > maxItems {
		t.Fatalf("expected %s item count in [%d,%d], got %d: %+v", meal, minItems, maxItems, len(items), items)
	}
}

func TestBuildBreakfastCombos_PreferYogurtAndEgg(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(1))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 1}, Name: "测试-无糖酸奶杯", Energy: 120, Protein: 9, Carbohydrate: 12, Fat: 3, Ingredients: []string{"无糖酸奶", "蓝莓"}},
		{Model: gorm.Model{ID: 2}, Name: "测试-水煮蛋", Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
		{Model: gorm.Model{ID: 3}, Name: "测试-全麦吐司片", Energy: 95, Protein: 4, Carbohydrate: 18, Fat: 1.2, Ingredients: []string{"全麦面包"}},
	}

	combos := svc.buildBreakfastCombos(recipes, 260, 20, 30, 10, models.GoalWeightLoss, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty breakfast combos")
	}

	first := combos[0]
	if len(first.Items) < 2 {
		t.Fatalf("expected first combo to contain at least 2 items, got %d", len(first.Items))
	}

	hasYogurt := false
	hasEgg := false
	for _, item := range first.Items {
		if item.ID == 1 {
			hasYogurt = true
		}
		if item.ID == 2 {
			hasEgg = true
		}
	}

	if !hasYogurt || !hasEgg {
		t.Fatalf("expected first combo to include yogurt and egg, got %+v", first.Items)
	}
}

func TestBuildMealCombos_GeneratesMultiItemCombo(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(7))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 10}, Name: "测试-鸡胸肉", Energy: 220, Protein: 34, Carbohydrate: 2, Fat: 8, Ingredients: []string{"鸡胸肉", "黑胡椒"}},
		{Model: gorm.Model{ID: 11}, Name: "测试-糙米饭", Energy: 180, Protein: 4, Carbohydrate: 38, Fat: 1.5, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 12}, Name: "测试-西兰花", Energy: 80, Protein: 4, Carbohydrate: 10, Fat: 1, Ingredients: []string{"西兰花"}},
		{Model: gorm.Model{ID: 13}, Name: "测试-南瓜", Energy: 100, Protein: 2, Carbohydrate: 22, Fat: 0.5, Ingredients: []string{"南瓜"}},
	}

	combos := svc.buildMealCombos(recipes, 520, 35, 50, 12, models.GoalWeightLoss, 10, 8, 2, 3, false, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty meal combos")
	}

	for _, combo := range combos {
		if len(combo.Items) < 2 {
			t.Fatalf("expected all combos to contain at least 2 items, got %+v", combo.Items)
		}
	}
}

func TestBuildMealCombos_RequireStapleWhenAvailable(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(11))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 20}, Name: "测试-虾仁蒸蛋", Energy: 200, Protein: 24, Carbohydrate: 4, Fat: 9, Ingredients: []string{"鸡蛋", "虾仁"}},
		{Model: gorm.Model{ID: 21}, Name: "测试-清炒时蔬", Energy: 120, Protein: 5, Carbohydrate: 14, Fat: 4, Ingredients: []string{"西兰花", "胡萝卜"}},
		{Model: gorm.Model{ID: 22}, Name: "测试-糙米饭", Energy: 180, Protein: 4, Carbohydrate: 38, Fat: 1.5, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 23}, Name: "测试-蒸南瓜", Energy: 100, Protein: 2, Carbohydrate: 22, Fat: 0.6, Ingredients: []string{"南瓜"}},
	}

	combos := svc.buildMealCombos(recipes, 560, 35, 60, 16, models.GoalWeightLoss, 10, 8, 2, 3, true, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty meal combos")
	}

	for _, combo := range combos {
		if !mealHasStaple(combo.Items) {
			t.Fatalf("expected combo to include staple when staple is available, got %+v", combo.Items)
		}
	}
}

func TestBuildMealCombos_AvoidDoubleMainStaple(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(19))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 30}, Name: "测试-扬州炒饭", Energy: 190, Protein: 7, Carbohydrate: 34, Fat: 5, Ingredients: []string{"米饭", "鸡蛋", "胡萝卜"}},
		{Model: gorm.Model{ID: 31}, Name: "测试-肉蛋盖饭", Energy: 210, Protein: 12, Carbohydrate: 32, Fat: 6, Ingredients: []string{"米饭", "猪肉", "鸡蛋"}},
		{Model: gorm.Model{ID: 32}, Name: "测试-清炒西兰花", Energy: 90, Protein: 4, Carbohydrate: 8, Fat: 3, Ingredients: []string{"西兰花"}},
		{Model: gorm.Model{ID: 33}, Name: "测试-虾仁蒸蛋", Energy: 160, Protein: 18, Carbohydrate: 3, Fat: 8, Ingredients: []string{"鸡蛋", "虾仁"}},
	}

	combos := svc.buildMealCombos(recipes, 560, 35, 60, 16, models.GoalWeightLoss, 10, 8, 2, 3, true, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty meal combos")
	}

	for _, combo := range combos {
		if countMainStapleRecipes(combo.Items) > 1 {
			t.Fatalf("expected combo to avoid double main staples, got %+v", combo.Items)
		}
	}
}

func TestBuildMealCombos_AvoidDoubleGrainStaple(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(21))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 34}, Name: "测试-糙米饭", Energy: 150, Protein: 3, Carbohydrate: 32, Fat: 1, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 35}, Name: "测试-荞麦面", Energy: 160, Protein: 6, Carbohydrate: 31, Fat: 1.2, Ingredients: []string{"荞麦面"}},
		{Model: gorm.Model{ID: 36}, Name: "测试-青椒牛柳", Energy: 180, Protein: 22, Carbohydrate: 6, Fat: 7, Ingredients: []string{"牛肉", "青椒"}},
		{Model: gorm.Model{ID: 37}, Name: "测试-清炒菜心", Energy: 80, Protein: 3, Carbohydrate: 8, Fat: 3, Ingredients: []string{"菜心"}},
	}

	combos := svc.buildMealCombos(recipes, 560, 35, 60, 16, models.GoalWeightLoss, 10, 8, 2, 3, true, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty meal combos")
	}

	for _, combo := range combos {
		if countGrainStapleRecipes(combo.Items) > 1 {
			t.Fatalf("expected combo to avoid double grain staples, got %+v", combo.Items)
		}
	}
}

func TestBuildBreakfastCombos_RejectMainMealStaples(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(23))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 40}, Name: "测试-蛋包饭", Energy: 220, Protein: 10, Carbohydrate: 30, Fat: 8, Ingredients: []string{"米饭", "鸡蛋"}},
		{Model: gorm.Model{ID: 41}, Name: "测试-无糖酸奶杯", Energy: 110, Protein: 9, Carbohydrate: 10, Fat: 3, Ingredients: []string{"无糖酸奶", "蓝莓"}},
		{Model: gorm.Model{ID: 42}, Name: "测试-水煮蛋", Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
	}

	combos := svc.buildBreakfastCombos(recipes, 260, 20, 30, 10, models.GoalWeightLoss, nil)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty breakfast combos")
	}

	for _, combo := range combos {
		for _, item := range combo.Items {
			if item.ID == 40 {
				t.Fatalf("expected breakfast combos to reject main-meal staple recipe, got %+v", combo.Items)
			}
		}
	}
}

func TestFilterRecipesByAllowedMealType_UsesExplicitAllowedList(t *testing.T) {
	svc := &RecipeServiceImpl{}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 60}, Name: "legacy breakfast", MealType: models.MealTypeBreakfast},
		{Model: gorm.Model{ID: 61}, Name: "lunch dinner only", MealType: models.MealTypeBreakfast, AllowedMealTypes: []string{string(models.MealTypeLunch), string(models.MealTypeDinner)}},
		{Model: gorm.Model{ID: 62}, Name: "breakfast allowed", MealType: models.MealTypeLunch, AllowedMealTypes: []string{string(models.MealTypeBreakfast), string(models.MealTypeLunch)}},
	}

	filtered := svc.filterRecipesByAllowedMealType(recipes, models.MealTypeBreakfast)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 breakfast recipes after filtering, got %d", len(filtered))
	}

	ids := map[uint]bool{}
	for _, recipe := range filtered {
		ids[recipe.ID] = true
	}

	if !ids[60] {
		t.Fatalf("expected legacy meal type fallback recipe to remain")
	}
	if !ids[62] {
		t.Fatalf("expected explicitly allowed breakfast recipe to remain")
	}
	if ids[61] {
		t.Fatalf("expected recipe excluded by explicit allowed list to be removed")
	}
}

func TestMergeRecipesByID_PreservesBaseAndAppendsUniqueExtras(t *testing.T) {
	base := []models.Recipe{
		{Model: gorm.Model{ID: 1}, Name: "base 1"},
		{Model: gorm.Model{ID: 2}, Name: "base 2"},
	}
	extras := []models.Recipe{
		{Model: gorm.Model{ID: 2}, Name: "duplicate 2"},
		{Model: gorm.Model{ID: 3}, Name: "extra 3"},
		{Name: "missing id"},
	}

	merged := mergeRecipesByID(base, extras)
	if len(merged) != 3 {
		t.Fatalf("expected 3 merged recipes, got %d", len(merged))
	}

	expectedIDs := []uint{1, 2, 3}
	for i, expectedID := range expectedIDs {
		if merged[i].ID != expectedID {
			t.Fatalf("expected merged[%d] id=%d, got %d", i, expectedID, merged[i].ID)
		}
	}
	if merged[1].Name != "base 2" {
		t.Fatalf("expected base recipe to win duplicate merge, got %s", merged[1].Name)
	}
}

func TestNormalizeCollaborativeScoreBoosts_CapsAndScalesScores(t *testing.T) {
	boosts := normalizeCollaborativeScoreBoosts([]RecipeScore{
		{RecipeID: 1, Score: 2},
		{RecipeID: 2, Score: 1},
		{RecipeID: 2, Score: 0.5},
		{RecipeID: 3, Score: -1},
	})

	if boosts[1] != maxCollaborativeCandidateBoost {
		t.Fatalf("expected top collaborative boost %.1f, got %.1f", maxCollaborativeCandidateBoost, boosts[1])
	}
	if boosts[2] != maxCollaborativeCandidateBoost/2 {
		t.Fatalf("expected scaled collaborative boost %.1f, got %.1f", maxCollaborativeCandidateBoost/2, boosts[2])
	}
	if _, ok := boosts[3]; ok {
		t.Fatalf("expected non-positive collaborative score to be ignored")
	}
	if got := collaborativeCandidateBoost(1, map[uint]float64{1: 99}); got != maxCollaborativeCandidateBoost {
		t.Fatalf("expected collaborative boost cap %.1f, got %.1f", maxCollaborativeCandidateBoost, got)
	}
}

func TestLoadCollaborativeMealCandidatesGroupsByMealAndBuildsBoosts(t *testing.T) {
	svc := &RecipeServiceImpl{
		collaborativeFilter: fakeCollaborativeFilter{
			models.MealTypeBreakfast: {
				{RecipeID: 80, Recipe: models.Recipe{Model: gorm.Model{ID: 80}, Name: "cf breakfast"}, Score: 4},
			},
			models.MealTypeLunch: {
				{RecipeID: 81, Recipe: models.Recipe{Model: gorm.Model{ID: 81}, Name: "cf lunch"}, Score: 2},
			},
		},
	}

	pools, boosts := svc.loadCollaborativeMealCandidates(context.Background(), 1, 10)
	if len(pools.breakfast) != 1 || pools.breakfast[0].ID != 80 {
		t.Fatalf("expected breakfast CF candidate 80, got %+v", pools.breakfast)
	}
	if len(pools.lunch) != 1 || pools.lunch[0].ID != 81 {
		t.Fatalf("expected lunch CF candidate 81, got %+v", pools.lunch)
	}
	if len(pools.dinner) != 0 || len(pools.snack) != 0 {
		t.Fatalf("expected empty dinner/snack CF pools, got dinner=%+v snack=%+v", pools.dinner, pools.snack)
	}
	if boosts[80] != maxCollaborativeCandidateBoost {
		t.Fatalf("expected top CF boost %.1f, got %.1f", maxCollaborativeCandidateBoost, boosts[80])
	}
	if boosts[81] != maxCollaborativeCandidateBoost/2 {
		t.Fatalf("expected scaled lunch CF boost %.1f, got %.1f", maxCollaborativeCandidateBoost/2, boosts[81])
	}
}

func TestSelectDiverseCandidates_AvoidsSingleVarietyDominance(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(29))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 50}, Name: "测试-鸡胸肉A", Energy: 160, Protein: 25, Carbohydrate: 2, Fat: 5, Ingredients: []string{"鸡胸"}},
		{Model: gorm.Model{ID: 51}, Name: "测试-鸡胸肉B", Energy: 158, Protein: 24, Carbohydrate: 2, Fat: 5, Ingredients: []string{"鸡胸"}},
		{Model: gorm.Model{ID: 52}, Name: "测试-鸡胸肉C", Energy: 162, Protein: 25, Carbohydrate: 2, Fat: 5, Ingredients: []string{"鸡胸"}},
		{Model: gorm.Model{ID: 53}, Name: "测试-牛肉片", Energy: 170, Protein: 23, Carbohydrate: 3, Fat: 7, Ingredients: []string{"牛肉"}},
		{Model: gorm.Model{ID: 54}, Name: "测试-虾仁滑蛋", Energy: 155, Protein: 20, Carbohydrate: 4, Fat: 6, Ingredients: []string{"虾仁", "鸡蛋"}},
		{Model: gorm.Model{ID: 55}, Name: "测试-豆腐蒸蛋", Energy: 140, Protein: 18, Carbohydrate: 5, Fat: 5, Ingredients: []string{"豆腐", "鸡蛋"}},
	}

	selected := svc.selectDiverseCandidates(recipes, 180, 24, 12, 8, 4, models.GoalWeightLoss, nil)
	if len(selected) != 4 {
		t.Fatalf("expected 4 selected candidates, got %d", len(selected))
	}

	keyCount := make(map[string]int)
	for _, recipe := range selected {
		keyCount[recipeVarietyKey(recipe)]++
	}
	for key, count := range keyCount {
		if count > 2 {
			t.Fatalf("expected diversified candidates, key %s count=%d", key, count)
		}
	}
}

func TestSelectDiverseCandidates_AppliesCollaborativeBoost(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(43))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 70}, Name: "exact match", Energy: 100, Protein: 10, Carbohydrate: 10, Fat: 5, Ingredients: []string{"exact"}},
		{Model: gorm.Model{ID: 71}, Name: "boosted near match", Energy: 110, Protein: 10, Carbohydrate: 10, Fat: 5, Ingredients: []string{"boosted"}},
	}

	selectedWithoutBoost := svc.selectDiverseCandidates(recipes, 100, 10, 10, 5, 1, models.GoalWeightLoss, nil)
	if len(selectedWithoutBoost) != 1 || selectedWithoutBoost[0].ID != 70 {
		t.Fatalf("expected exact nutrition match without CF boost, got %+v", selectedWithoutBoost)
	}

	selectedWithBoost := svc.selectDiverseCandidates(recipes, 100, 10, 10, 5, 1, models.GoalWeightLoss, map[uint]float64{71: maxCollaborativeCandidateBoost})
	if len(selectedWithBoost) != 1 || selectedWithBoost[0].ID != 71 {
		t.Fatalf("expected collaborative boost to promote recipe 71, got %+v", selectedWithBoost)
	}
}

func TestCheckPlanDiversity_RejectsDuplicateLunchDinnerCombos(t *testing.T) {
	svc := &RecipeServiceImpl{}

	existing := []*models.DailyRecipePlan{
		{
			LunchRecipeID:  60,
			DinnerRecipeID: 70,
			LunchItemIDs:   []uint{60, 61},
			DinnerItemIDs:  []uint{70, 71},
		},
	}

	newPlan := &models.DailyRecipePlan{
		LunchRecipeID:  60,
		DinnerRecipeID: 70,
		LunchItemIDs:   []uint{60, 61},
		DinnerItemIDs:  []uint{70, 71},
	}

	if svc.checkPlanDiversity(existing, newPlan) {
		t.Fatalf("expected duplicate lunch and dinner combos to be rejected")
	}
}

func TestGenerateRankedPlans_ReturnsPlansForSmallManualPools(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(31))}

	breakfast := []models.Recipe{
		{Model: gorm.Model{ID: 100}, Name: "测试-水煮蛋", MealType: models.MealTypeBreakfast, Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
		{Model: gorm.Model{ID: 101}, Name: "测试-无糖豆浆", MealType: models.MealTypeBreakfast, Energy: 31, Protein: 3.0, Carbohydrate: 1.8, Fat: 1.6, Ingredients: []string{"黄豆"}},
		{Model: gorm.Model{ID: 102}, Name: "测试-蒸玉米", MealType: models.MealTypeBreakfast, Energy: 112, Protein: 3.2, Carbohydrate: 22.8, Fat: 1.2, Ingredients: []string{"玉米"}},
	}

	lunch := []models.Recipe{
		{Model: gorm.Model{ID: 110}, Name: "测试-番茄鸡胸肉", MealType: models.MealTypeLunch, Energy: 118, Protein: 18.6, Carbohydrate: 4.2, Fat: 3.1, Ingredients: []string{"鸡胸肉", "番茄"}},
		{Model: gorm.Model{ID: 111}, Name: "测试-糙米饭", MealType: models.MealTypeLunch, Energy: 116, Protein: 2.6, Carbohydrate: 25.1, Fat: 0.9, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 112}, Name: "测试-清炒油麦菜", MealType: models.MealTypeLunch, Energy: 29, Protein: 1.4, Carbohydrate: 3.7, Fat: 1.0, Ingredients: []string{"油麦菜"}},
	}

	dinner := []models.Recipe{
		{Model: gorm.Model{ID: 120}, Name: "测试-清蒸鲈鱼", MealType: models.MealTypeDinner, Energy: 106, Protein: 17.9, Carbohydrate: 1.2, Fat: 3.4, Ingredients: []string{"鲈鱼"}},
		{Model: gorm.Model{ID: 121}, Name: "测试-荞麦面", MealType: models.MealTypeDinner, Energy: 99, Protein: 5.1, Carbohydrate: 19.8, Fat: 0.8, Ingredients: []string{"荞麦面"}},
		{Model: gorm.Model{ID: 122}, Name: "测试-蒜蓉西兰花", MealType: models.MealTypeDinner, Energy: 36, Protein: 2.7, Carbohydrate: 5.1, Fat: 1.1, Ingredients: []string{"西兰花"}},
	}

	snack := []models.Recipe{
		{Model: gorm.Model{ID: 130}, Name: "测试-苹果", MealType: models.MealTypeSnack, Energy: 53, Protein: 0.3, Carbohydrate: 13.7, Fat: 0.2, Ingredients: []string{"苹果"}},
		{Model: gorm.Model{ID: 131}, Name: "测试-无糖酸奶", MealType: models.MealTypeSnack, Energy: 72, Protein: 3.8, Carbohydrate: 6.2, Fat: 3.4, Ingredients: []string{"酸奶"}},
	}

	target := NutritionTarget{
		Energy:       1600,
		Protein:      110,
		Carbohydrate: 160,
		Fat:          45,
	}

	plans := svc.generateRankedPlans(1, breakfast, lunch, dinner, snack, target, models.GoalWeightLoss, nil)
	if len(plans) == 0 {
		t.Fatalf("expected non-empty plans for small manual recipe pools")
	}

	for _, plan := range plans {
		assertItemCountInRange(t, "breakfast", plan.BreakfastItems, breakfastComboMinItems, breakfastComboMaxItems)
		assertItemCountInRange(t, "lunch", plan.LunchItems, lunchComboMinItems, lunchComboMaxItems)
		assertItemCountInRange(t, "dinner", plan.DinnerItems, dinnerComboMinItems, dinnerComboMaxItems)
		assertItemCountInRange(t, "snack", plan.SnackItems, snackComboMinItems, snackComboMaxItems)
	}
}

func TestBuildRecommendedPlansWithFallback_UsesRelaxedStageWhenStrictStageIsEmpty(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(37))}

	breakfast := []models.Recipe{
		{Model: gorm.Model{ID: 200}, Name: "测试-水煮蛋", MealType: models.MealTypeBreakfast, Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
		{Model: gorm.Model{ID: 201}, Name: "测试-无糖豆浆", MealType: models.MealTypeBreakfast, Energy: 31, Protein: 3.0, Carbohydrate: 1.8, Fat: 1.6, Ingredients: []string{"黄豆"}},
		{Model: gorm.Model{ID: 202}, Name: "测试-蒸玉米", MealType: models.MealTypeBreakfast, Energy: 112, Protein: 3.2, Carbohydrate: 22.8, Fat: 1.2, Ingredients: []string{"玉米"}},
	}
	lunch := []models.Recipe{
		{Model: gorm.Model{ID: 210}, Name: "测试-番茄鸡胸肉", MealType: models.MealTypeLunch, Energy: 118, Protein: 18.6, Carbohydrate: 4.2, Fat: 3.1, Ingredients: []string{"鸡胸肉", "番茄"}},
		{Model: gorm.Model{ID: 211}, Name: "测试-糙米饭", MealType: models.MealTypeLunch, Energy: 116, Protein: 2.6, Carbohydrate: 25.1, Fat: 0.9, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 212}, Name: "测试-清炒油麦菜", MealType: models.MealTypeLunch, Energy: 29, Protein: 1.4, Carbohydrate: 3.7, Fat: 1.0, Ingredients: []string{"油麦菜"}},
	}
	dinner := []models.Recipe{
		{Model: gorm.Model{ID: 220}, Name: "测试-清蒸鲈鱼", MealType: models.MealTypeDinner, Energy: 106, Protein: 17.9, Carbohydrate: 1.2, Fat: 3.4, Ingredients: []string{"鲈鱼"}},
		{Model: gorm.Model{ID: 221}, Name: "测试-荞麦面", MealType: models.MealTypeDinner, Energy: 99, Protein: 5.1, Carbohydrate: 19.8, Fat: 0.8, Ingredients: []string{"荞麦面"}},
		{Model: gorm.Model{ID: 222}, Name: "测试-蒜蓉西兰花", MealType: models.MealTypeDinner, Energy: 36, Protein: 2.7, Carbohydrate: 5.1, Fat: 1.1, Ingredients: []string{"西兰花"}},
	}
	snack := []models.Recipe{
		{Model: gorm.Model{ID: 230}, Name: "测试-苹果", MealType: models.MealTypeSnack, Energy: 53, Protein: 0.3, Carbohydrate: 13.7, Fat: 0.2, Ingredients: []string{"苹果"}},
	}

	target := NutritionTarget{
		Energy:       1600,
		Protein:      110,
		Carbohydrate: 160,
		Fat:          45,
	}

	stages := []recommendationStage{
		{
			label: "strict",
			pools: mealRecipePools{
				breakfast: nil,
				lunch:     lunch,
				dinner:    dinner,
				snack:     snack,
			},
		},
		{
			label: "relaxed",
			pools: mealRecipePools{
				breakfast: breakfast,
				lunch:     lunch,
				dinner:    dinner,
				snack:     snack,
			},
		},
	}

	plans, stage := svc.buildRecommendedPlansWithFallback(1, stages, target, models.GoalWeightLoss, 2, svc.checkDiversity)
	if len(plans) == 0 {
		t.Fatalf("expected fallback stage to produce plans")
	}
	if stage != "relaxed" {
		t.Fatalf("expected fallback to use relaxed stage, got %s", stage)
	}
}

func TestGenerateRankedPlans_FallsBackWhenStrictMealAssemblyProducesNoPlans(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(41))}

	breakfast := []models.Recipe{
		{Model: gorm.Model{ID: 300}, Name: "测试-水煮蛋", MealType: models.MealTypeBreakfast, Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
		{Model: gorm.Model{ID: 301}, Name: "测试-无糖豆浆", MealType: models.MealTypeBreakfast, Energy: 31, Protein: 3.0, Carbohydrate: 1.8, Fat: 1.6, Ingredients: []string{"黄豆"}},
		{Model: gorm.Model{ID: 302}, Name: "测试-蒸玉米", MealType: models.MealTypeBreakfast, Energy: 112, Protein: 3.2, Carbohydrate: 22.8, Fat: 1.2, Ingredients: []string{"玉米"}},
	}
	lunch := []models.Recipe{
		{Model: gorm.Model{ID: 310}, Name: "测试-番茄鸡胸肉", MealType: models.MealTypeLunch, Energy: 118, Protein: 18.6, Carbohydrate: 4.2, Fat: 3.1, Ingredients: []string{"鸡胸肉", "番茄"}},
		{Model: gorm.Model{ID: 311}, Name: "测试-糙米饭", MealType: models.MealTypeLunch, Energy: 116, Protein: 2.6, Carbohydrate: 25.1, Fat: 0.9, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 312}, Name: "测试-清炒油麦菜", MealType: models.MealTypeLunch, Energy: 29, Protein: 1.4, Carbohydrate: 3.7, Fat: 1.0, Ingredients: []string{"油麦菜"}},
	}
	dinner := []models.Recipe{
		{Model: gorm.Model{ID: 320}, Name: "测试-清蒸鲈鱼", MealType: models.MealTypeDinner, Energy: 106, Protein: 17.9, Carbohydrate: 1.2, Fat: 3.4, Ingredients: []string{"鲈鱼"}},
		{Model: gorm.Model{ID: 311}, Name: "测试-糙米饭", MealType: models.MealTypeDinner, Energy: 116, Protein: 2.6, Carbohydrate: 25.1, Fat: 0.9, Ingredients: []string{"糙米"}},
		{Model: gorm.Model{ID: 322}, Name: "测试-蒜蓉西兰花", MealType: models.MealTypeDinner, Energy: 36, Protein: 2.7, Carbohydrate: 5.1, Fat: 1.1, Ingredients: []string{"西兰花"}},
	}
	snack := []models.Recipe{
		{Model: gorm.Model{ID: 330}, Name: "测试-苹果", MealType: models.MealTypeSnack, Energy: 53, Protein: 0.3, Carbohydrate: 13.7, Fat: 0.2, Ingredients: []string{"苹果"}},
	}

	target := NutritionTarget{
		Energy:       1600,
		Protein:      110,
		Carbohydrate: 160,
		Fat:          45,
	}

	plans := svc.generateRankedPlans(1, breakfast, lunch, dinner, snack, target, models.GoalWeightLoss, nil)
	if len(plans) == 0 {
		t.Fatalf("expected fallback plans when strict meal assembly is empty")
	}
}

func TestIngredientDiversityPenalty_WithOverlap(t *testing.T) {
	svc := &RecipeServiceImpl{}

	comboA := mealCombo{Items: []models.Recipe{{Ingredients: []string{"鸡蛋", "番茄"}}}}
	comboB := mealCombo{Items: []models.Recipe{{Ingredients: []string{"鸡蛋", "生菜"}}}}
	comboC := mealCombo{Items: []models.Recipe{{Ingredients: []string{"牛肉"}}}}

	penalty := svc.ingredientDiversityPenalty(comboA, comboB, comboC)
	if penalty <= 0 {
		t.Fatalf("expected positive diversity penalty, got %f", penalty)
	}
}
