package services

import (
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	"testing"
)

func TestAccessExamBankAndThreshold(t *testing.T) {
	for _, level := range []string{"B2", "C1", "C2"} {
		for n := 0; n < 30; n++ {
			activity, err := BuildAccessExam(level)
			if err != nil {
				t.Fatalf("%s: %v", level, err)
			}
			if len(activity.Questions) != 12 || activity.TimeLimitSeconds != 1200 {
				t.Fatal("wrong exam size or deadline")
			}
			for _, right := range []int{9, 10} {
				answers := models.Submission{Answers: map[string]json.RawMessage{}}
				for i, q := range activity.Questions {
					if i < right {
						answers.Answers[q.ID] = q.Solution
					}
				}
				result, err := Grade(activity, answers)
				if err != nil {
					t.Fatal(err)
				}
				if result.Passed != (right == 10) {
					t.Fatalf("wrong 80%% threshold: %+v", result)
				}
			}
		}
	}
	if _, err := BuildAccessExam("A1"); err == nil {
		t.Fatal("initial level unexpectedly has an access exam")
	}
}
func TestDiplomaNames(t *testing.T) {
	for _, name := range []string{"Ana María López", "O'Connor", "Jean-Luc Martin", "D’Arcy", "李明"} {
		if _, err := DiplomaName(name); err != nil {
			t.Fatalf("valid name %q: %v", name, err)
		}
	}
	for _, name := range []string{"", "A", "12345", "<script>alert(1)</script>", "Ana\nMaria", "Ana\x00Maria"} {
		if _, err := DiplomaName(name); err == nil {
			t.Fatalf("unsafe or invalid name accepted: %q", name)
		}
	}
	name, err := DiplomaName("  Ana   María  ")
	if err != nil || name != "Ana María" {
		t.Fatal("name normalization failed")
	}
}
