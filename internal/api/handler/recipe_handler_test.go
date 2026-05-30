package handler

import "testing"

func TestBuildRecommendationAlgorithmDTO(t *testing.T) {
	meta := buildRecommendationAlgorithmDTO()
	if meta.Name != "item_cf_hybrid" {
		t.Fatalf("expected item_cf_hybrid algorithm name, got %s", meta.Name)
	}
	if !meta.CollaborativeFiltering {
		t.Fatalf("expected collaborative filtering flag to be enabled")
	}
	if meta.ScoreBoostCap != 15 {
		t.Fatalf("expected score boost cap 15, got %.1f", meta.ScoreBoostCap)
	}
	if len(meta.Constraints) == 0 {
		t.Fatalf("expected non-empty recommendation constraints")
	}
}
