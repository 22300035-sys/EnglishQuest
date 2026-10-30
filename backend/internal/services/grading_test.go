package services

import (
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	"testing"
)

func TestGradeAllGames(t *testing.T) {
	for _, tc := range []struct{ kind, key, good, bad string }{
		{"word_hunter", `"random-option"`, `"random-option"`, `"another-option"`},
		{"sentence_builder", `["opaque-b","opaque-a"]`, `["opaque-b","opaque-a"]`, `["opaque-a","opaque-b"]`},
		{"memory_match", `[["a","b"],["c","d"]]`, `[["d","c"],["b","a"]]`, `[["a","b"],["a","d"]]`},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			a := models.Activity{Kind: tc.kind, Questions: []models.Question{{ID: "q", Solution: json.RawMessage(tc.key)}}}
			for _, answer := range []struct {
				value string
				score int
			}{{tc.good, 100}, {tc.bad, 0}} {
				result, err := Grade(a, models.Submission{Answers: map[string]json.RawMessage{"q": json.RawMessage(answer.value)}})
				if err != nil || result.Score != answer.score || result.XP != 0 {
					t.Fatalf("unexpected grade %+v %v", result, err)
				}
			}
			if _, err := Grade(a, models.Submission{}); err == nil {
				t.Fatal("missing answer accepted")
			}
		})
	}
}
func TestAssessmentThreshold(t *testing.T) {
	a := models.Activity{Kind: "word_hunter", Assessment: true, Questions: []models.Question{}}
	answers := models.Submission{Answers: map[string]json.RawMessage{}}
	for _, id := range []string{"a", "b", "c"} {
		a.Questions = append(a.Questions, models.Question{ID: id, Solution: json.RawMessage(`"yes"`)})
		answers.Answers[id] = json.RawMessage(`"yes"`)
	}
	answers.Answers["c"] = json.RawMessage(`"no"`)
	r, err := Grade(a, answers)
	if err != nil || r.Score != 66 || r.Passed {
		t.Fatal("assessment should require 80%")
	}
	a.Assessment = false
	r, err = Grade(a, answers)
	if err != nil || !r.Passed {
		t.Fatal("practice should pass at 60%")
	}
}
func TestCivilDateStreak(t *testing.T) {
	for _, tc := range []struct {
		last, today                          string
		current, best, wantCurrent, wantBest int
	}{
		{"", "2026-10-25", 0, 0, 1, 1}, {"2026-10-25", "2026-10-25", 4, 7, 4, 7},
		{"2026-10-24", "2026-10-25", 6, 6, 7, 7}, {"2026-10-22", "2026-10-25", 7, 10, 1, 10},
		{"2026-03-07", "2026-03-08", 2, 2, 3, 3}, {"2026-12-31", "2027-01-01", 2, 5, 3, 5},
	} {
		c, b := NextStreak(tc.last, tc.today, tc.current, tc.best)
		if c != tc.wantCurrent || b != tc.wantBest {
			t.Fatalf("streak %s -> %s: %d %d", tc.last, tc.today, c, b)
		}
	}
}
func TestAdminContentValidation(t *testing.T) {
	a := models.Activity{Title: "Words", Kind: "word_hunter", Instructions: "Read and select.", Questions: []models.Question{{Prompt: "Find: apple", Explanation: "Apple means manzana.", Solution: json.RawMessage(`"a"`), Options: []models.Option{{ID: "a", Label: "Manzana", Image: "apple"}, {ID: "b", Label: "Plátano", Image: "banana"}}}}}
	if err := ValidateActivity(a); err != nil {
		t.Fatal(err)
	}
	a.Questions[0].Options[0].Image = "https://evil.example/track"
	if ValidateActivity(a) == nil {
		t.Fatal("unapproved image accepted")
	}
}

func TestTextLettersAndTimedPartialAnswers(t *testing.T) {
	a := models.Activity{Kind: "grammar_rescue", Questions: []models.Question{{ID: "q", Solution: json.RawMessage(`"She goes to school."`)}}}
	r, err := Grade(a, models.Submission{Answers: map[string]json.RawMessage{"q": json.RawMessage(`"  she GOES   to school!? "`)}})
	if err != nil || r.Score != 100 {
		t.Fatalf("normalization failed: %+v %v", r, err)
	}
	a = models.Activity{Kind: "spell_master", Questions: []models.Question{{ID: "q", Solution: json.RawMessage(`"apple"`), Options: []models.Option{{ID: "a", Label: "a"}, {ID: "p1", Label: "p"}, {ID: "p2", Label: "p"}, {ID: "l", Label: "l"}, {ID: "e", Label: "e"}}}}}
	for _, tc := range []struct {
		answer string
		score  int
	}{{`["a","p2","p1","l","e"]`, 100}, {`["a","p1","p1","l","e"]`, 0}} {
		result, err := Grade(a, models.Submission{Answers: map[string]json.RawMessage{"q": json.RawMessage(tc.answer)}})
		if err != nil || result.Score != tc.score {
			t.Fatal("letter tile reuse or repeated letters mishandled")
		}
	}
	a = models.Activity{Kind: "english_rush", TimeLimitSeconds: 120, Questions: []models.Question{{ID: "q", Response: "choice", Solution: json.RawMessage(`"a"`)}}}
	r, err = Grade(a, models.Submission{})
	if err != nil || r.Score != 0 || r.Passed {
		t.Fatal("timed partial answer was not graded as incomplete")
	}
}
