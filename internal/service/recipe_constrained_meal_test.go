package service

import (
	"NutriPlan/internal/repository/models"
	"testing"
)

func TestNormalizeMealType(t *testing.T) {
	tests := []struct {
		input      string
		expectCN   models.MealType
		expectEN   string
		expectFail bool
	}{
		{input: "breakfast", expectCN: models.MealTypeBreakfast, expectEN: "breakfast"},
		{input: "午餐", expectCN: models.MealTypeLunch, expectEN: "lunch"},
		{input: "dinner", expectCN: models.MealTypeDinner, expectEN: "dinner"},
		{input: "invalid", expectFail: true},
	}

	for _, tc := range tests {
		mealCN, mealEN, err := normalizeMealType(tc.input)
		if tc.expectFail {
			if err == nil {
				t.Fatalf("expected error for %s", tc.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tc.input, err)
		}
		if mealCN != tc.expectCN || mealEN != tc.expectEN {
			t.Fatalf("unexpected normalize result for %s: got (%s,%s)", tc.input, mealCN, mealEN)
		}
	}
}

func TestUnmarshalAndValidateGeneratedMeal(t *testing.T) {
	raw := `{"meal_name":"番茄牛肉饭","meal_type":"午餐","core_ingredients_used":["番茄","牛肉"],"supplementary_ingredients":["米饭"],"nutrition_estimate":{"energy":520,"protein":32,"carbohydrate":48,"fat":18},"steps":["处理食材","翻炒并焖煮"],"dietitian_tip":"控制盐分","fill_reason":"补充主食"}`

	var out GeneratedMeal
	err := unmarshalAndValidateGeneratedMeal(raw, []string{"番茄", "牛肉"}, models.MealTypeLunch, &out)
	if err != nil {
		t.Fatalf("expected valid result, got error: %v", err)
	}
	if out.MealName == "" || len(out.Steps) == 0 {
		t.Fatalf("expected non-empty parsed meal")
	}
}

func TestUnmarshalAndValidateGeneratedMeal_MissingCoreIngredient(t *testing.T) {
	raw := `{"meal_name":"黄瓜沙拉","meal_type":"午餐","core_ingredients_used":["黄瓜"],"supplementary_ingredients":[],"nutrition_estimate":{"energy":180,"protein":6,"carbohydrate":20,"fat":8},"steps":["切配并拌匀"],"dietitian_tip":"可搭配蛋白质来源","fill_reason":""}`

	var out GeneratedMeal
	err := unmarshalAndValidateGeneratedMeal(raw, []string{"牛肉"}, models.MealTypeLunch, &out)
	if err == nil {
		t.Fatalf("expected validation error when core ingredient is missing")
	}
}
