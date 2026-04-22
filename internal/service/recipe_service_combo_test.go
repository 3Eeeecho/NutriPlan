package service

import (
	"NutriPlan/internal/repository/models"
	"math/rand"
	"testing"

	"gorm.io/gorm"
)

func TestBuildBreakfastCombos_PreferYogurtAndEgg(t *testing.T) {
	svc := &RecipeServiceImpl{rng: rand.New(rand.NewSource(1))}

	recipes := []models.Recipe{
		{Model: gorm.Model{ID: 1}, Name: "测试-无糖酸奶杯", Energy: 120, Protein: 9, Carbohydrate: 12, Fat: 3, Ingredients: []string{"无糖酸奶", "蓝莓"}},
		{Model: gorm.Model{ID: 2}, Name: "测试-水煮蛋", Energy: 78, Protein: 6.5, Carbohydrate: 0.6, Fat: 5.3, Ingredients: []string{"鸡蛋"}},
		{Model: gorm.Model{ID: 3}, Name: "测试-全麦吐司片", Energy: 95, Protein: 4, Carbohydrate: 18, Fat: 1.2, Ingredients: []string{"全麦面包"}},
	}

	combos := svc.buildBreakfastCombos(recipes, 260, 20, 30, 10, models.GoalWeightLoss)
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

	combos := svc.buildMealCombos(recipes, 520, 35, 50, 12, models.GoalWeightLoss, 10, 8, 2, 3, false)
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

	combos := svc.buildMealCombos(recipes, 560, 35, 60, 16, models.GoalWeightLoss, 10, 8, 2, 3, true)
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

	combos := svc.buildMealCombos(recipes, 560, 35, 60, 16, models.GoalWeightLoss, 10, 8, 2, 3, true)
	if len(combos) == 0 {
		t.Fatalf("expected non-empty meal combos")
	}

	for _, combo := range combos {
		if countMainStapleRecipes(combo.Items) > 1 {
			t.Fatalf("expected combo to avoid double main staples, got %+v", combo.Items)
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

	combos := svc.buildBreakfastCombos(recipes, 260, 20, 30, 10, models.GoalWeightLoss)
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

	selected := svc.selectDiverseCandidates(recipes, 180, 24, 12, 8, 4, models.GoalWeightLoss)
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
