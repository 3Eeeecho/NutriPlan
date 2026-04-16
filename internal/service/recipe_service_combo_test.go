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
