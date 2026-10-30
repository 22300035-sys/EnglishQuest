package services

import (
	"github.com/englishquest/backend/internal/models"
	"testing"
)

func TestBoundedRewardMultipliers(t *testing.T) {
	for _, tc := range []struct {
		assessment, perfect bool
		factor, want        int
	}{
		{false, false, 1, 20}, {false, true, 1, 30}, {false, true, 2, 50},
		{true, true, 2, 90}, {true, false, 2, 80}, {false, true, 99, 30},
	} {
		a := models.Activity{Assessment: tc.assessment}
		if got := a.RewardXP(tc.perfect, tc.factor); got != tc.want {
			t.Fatalf("reward %d wanted %d", got, tc.want)
		}
	}
}
