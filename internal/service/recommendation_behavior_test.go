package service

import (
	"NutriPlan/internal/repository/models"
	"math"
	"testing"
	"time"
)

func TestAggregateRecipePreferenceEventsCombinesWeightsWithDecay(t *testing.T) {
	now := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	events := []recipePreferenceEvent{
		{UserID: 1, RecipeID: 10, Weight: recipeFavoriteWeight, EventTime: now},
		{UserID: 1, RecipeID: 10, Weight: recipeIntakeWeight, EventTime: now.AddDate(0, 0, -30)},
		{UserID: 1, RecipeID: 11, Weight: recipeSelectedPlanWeight, EventTime: now},
		{UserID: 0, RecipeID: 11, Weight: recipeIntakeWeight, EventTime: now},
	}

	preferences := aggregateRecipePreferenceEvents(events, now)
	byRecipe := make(map[uint]UserRecipePreference)
	for _, preference := range preferences {
		byRecipe[preference.RecipeID] = preference
	}

	got := byRecipe[10]
	want := round4(recipeFavoriteWeight + recipeIntakeWeight*math.Exp(-1))
	if got.Score != want {
		t.Fatalf("expected score %.4f, got %.4f", want, got.Score)
	}
	if got.EventCount != 2 {
		t.Fatalf("expected event count 2, got %d", got.EventCount)
	}
	if !got.LastEventAt.Equal(now) {
		t.Fatalf("expected latest event time %v, got %v", now, got.LastEventAt)
	}
	if _, ok := byRecipe[11]; !ok {
		t.Fatal("expected recipe 11 preference")
	}
	if len(preferences) != 2 {
		t.Fatalf("expected 2 preferences after filtering invalid events, got %d", len(preferences))
	}
}

func TestPlanRecipeIDsDeduplicatesPrimaryAndComboItems(t *testing.T) {
	plan := &models.DailyRecipePlan{
		BreakfastRecipeID: 1,
		LunchRecipeID:     2,
		DinnerRecipeID:    3,
		SnackRecipeID:     0,
		BreakfastItemIDs:  []uint{1, 4},
		LunchItemIDs:      []uint{2, 4, 5},
		DinnerItemIDs:     []uint{3},
		SnackItemIDs:      []uint{0, 6},
	}

	ids := planRecipeIDs(plan)
	want := []uint{1, 4, 2, 5, 3, 6}
	if len(ids) != len(want) {
		t.Fatalf("expected %d ids, got %d: %#v", len(want), len(ids), ids)
	}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("expected ids %#v, got %#v", want, ids)
		}
	}
}

func TestNormalizeBehaviorWindowDays(t *testing.T) {
	if got := normalizeBehaviorWindowDays(0); got != 90 {
		t.Fatalf("expected default 90 days, got %d", got)
	}
	if got := normalizeBehaviorWindowDays(500); got != 365 {
		t.Fatalf("expected cap at 365 days, got %d", got)
	}
	if got := normalizeBehaviorWindowDays(30); got != 30 {
		t.Fatalf("expected 30 days, got %d", got)
	}
}
