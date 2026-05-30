package service

import (
	"context"
	"testing"
	"time"
)

type countingBehaviorRepository struct {
	userCalls int
	allCalls  int
	userPrefs []UserRecipePreference
	allPrefs  []UserRecipePreference
}

func (r *countingBehaviorRepository) GetUserRecipePreferences(ctx context.Context, userID uint, days int) ([]UserRecipePreference, error) {
	r.userCalls++
	return r.userPrefs, nil
}

func (r *countingBehaviorRepository) GetAllRecipePreferences(ctx context.Context, days int) ([]UserRecipePreference, error) {
	r.allCalls++
	return r.allPrefs, nil
}

func TestCosineSimilarity(t *testing.T) {
	left := map[uint]float64{1: 3, 2: 4}
	right := map[uint]float64{1: 3, 2: 4}

	if got := cosineSimilarity(left, right); got != 1 {
		t.Fatalf("expected identical vectors to have similarity 1, got %.4f", got)
	}

	if got := cosineSimilarity(left, map[uint]float64{3: 5}); got != 0 {
		t.Fatalf("expected disjoint vectors to have similarity 0, got %.4f", got)
	}
}

func TestRankCollaborativeRecipeScoresRanksSimilarUnseenRecipes(t *testing.T) {
	userPrefs := []UserRecipePreference{
		{UserID: 1, RecipeID: 10, Score: 5},
		{UserID: 1, RecipeID: 11, Score: 2},
	}
	allPrefs := []UserRecipePreference{
		{UserID: 1, RecipeID: 10, Score: 5},
		{UserID: 1, RecipeID: 11, Score: 2},
		{UserID: 2, RecipeID: 10, Score: 4},
		{UserID: 2, RecipeID: 20, Score: 4},
		{UserID: 3, RecipeID: 10, Score: 3},
		{UserID: 3, RecipeID: 20, Score: 3},
		{UserID: 4, RecipeID: 11, Score: 3},
		{UserID: 4, RecipeID: 30, Score: 3},
		{UserID: 5, RecipeID: 11, Score: 2},
		{UserID: 5, RecipeID: 30, Score: 2},
		{UserID: 6, RecipeID: 40, Score: 5},
	}

	ranked := rankCollaborativeRecipeScores(userPrefs, allPrefs, 10)
	if len(ranked) != 2 {
		t.Fatalf("expected two collaborative candidates, got %#v", ranked)
	}
	if ranked[0].RecipeID != 20 {
		t.Fatalf("expected recipe 20 to rank first, got %#v", ranked)
	}
	if ranked[0].Score <= ranked[1].Score {
		t.Fatalf("expected first score to be greater than second: %#v", ranked)
	}
	for _, candidate := range ranked {
		if candidate.RecipeID == 10 || candidate.RecipeID == 11 {
			t.Fatalf("expected seen recipe to be excluded: %#v", ranked)
		}
	}
}

func TestRankCollaborativeRecipeScoresReturnsEmptyForColdStart(t *testing.T) {
	ranked := rankCollaborativeRecipeScores(nil, []UserRecipePreference{{UserID: 2, RecipeID: 20, Score: 4}}, 10)
	if len(ranked) != 0 {
		t.Fatalf("expected no candidates for cold start, got %#v", ranked)
	}
}

func TestBuildItemUserPreferenceVectorsAggregatesDuplicatePreferences(t *testing.T) {
	vectors := buildItemUserPreferenceVectors([]UserRecipePreference{
		{UserID: 1, RecipeID: 10, Score: 2},
		{UserID: 1, RecipeID: 10, Score: 3},
		{UserID: 2, RecipeID: 10, Score: 4},
		{UserID: 0, RecipeID: 10, Score: 9},
	})

	if got := vectors[10][1]; got != 5 {
		t.Fatalf("expected duplicate user-item score to aggregate to 5, got %.2f", got)
	}
	if got := vectors[10][2]; got != 4 {
		t.Fatalf("expected user 2 score 4, got %.2f", got)
	}
}

func TestItemCollaborativeFilterCachesRankedCandidates(t *testing.T) {
	repo := &countingBehaviorRepository{
		userPrefs: []UserRecipePreference{
			{UserID: 1, RecipeID: 10, Score: 5},
		},
		allPrefs: []UserRecipePreference{
			{UserID: 1, RecipeID: 10, Score: 5},
			{UserID: 2, RecipeID: 10, Score: 4},
			{UserID: 2, RecipeID: 20, Score: 4},
		},
	}
	filter := &ItemCollaborativeFilter{
		behaviorRepo: repo,
		days:         defaultCollaborativeFilterDays,
		cacheTTL:     time.Minute,
		cache:        make(map[uint]collaborativeRankCacheEntry),
	}

	first, err := filter.rankRecipeCandidates(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected first ranking error: %v", err)
	}
	second, err := filter.rankRecipeCandidates(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected second ranking error: %v", err)
	}

	if len(first) != 1 || len(second) != 1 || first[0].RecipeID != second[0].RecipeID {
		t.Fatalf("expected cached ranked candidates to match, first=%#v second=%#v", first, second)
	}
	if repo.userCalls != 1 || repo.allCalls != 1 {
		t.Fatalf("expected behavior repository to be called once, userCalls=%d allCalls=%d", repo.userCalls, repo.allCalls)
	}
}
