package services

import (
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"reflect"
	"sort"
	"strings"
	"time"
)

// Grade ignores client score/XP; only server snapshot keys determine the result.
func Grade(activity models.Activity, answers models.Submission) (models.Result, error) {
	result := models.Result{ActivityID: activity.ID, Feedback: []models.Feedback{}}
	if !GameKinds[activity.Kind] {
		return result, apperr.NewAppError(400, "Tipo de actividad no disponible.")
	}
	known := map[string]bool{}
	for _, q := range activity.Questions {
		known[q.ID] = true
	}
	for id := range answers.Answers {
		if !known[id] {
			return result, apperr.NewAppError(400, "Pregunta ajena al intento.")
		}
	}
	if len(activity.Questions) == 0 || len(answers.Answers) > len(activity.Questions) || (activity.TimeLimitSeconds == 0 && len(answers.Answers) != len(activity.Questions)) {
		return result, apperr.NewAppError(400, "Responde todas las preguntas.")
	}
	correct := 0
	for _, q := range activity.Questions {
		raw, exists := answers.Answers[q.ID]
		if !exists && activity.TimeLimitSeconds == 0 {
			return result, apperr.NewAppError(400, "Falta una respuesta.")
		}
		match := false
		if !exists {
			result.Feedback = append(result.Feedback, models.Feedback{QuestionID: q.ID, Explanation: q.Explanation, Solution: q.Solution})
			continue
		}
		switch q.AnswerType(activity.Kind) {
		case "choice":
			var got, want string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil {
				return result, apperr.NewAppError(400, "Respuesta inválida.")
			}
			match = got == want
		case "order":
			var got, want []string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil {
				return result, apperr.NewAppError(400, "Orden de palabras inválido.")
			}
			match = reflect.DeepEqual(got, want)
		case "pairs":
			var got, want [][]string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil {
				return result, apperr.NewAppError(400, "Parejas inválidas.")
			}
			canon := func(pairs [][]string) ([]string, bool) {
				out := []string{}
				seen := map[string]bool{}
				for _, p := range pairs {
					if len(p) != 2 || p[0] == p[1] || seen[p[0]] || seen[p[1]] {
						return nil, false
					}
					seen[p[0]] = true
					seen[p[1]] = true
					sort.Strings(p)
					out = append(out, p[0]+":"+p[1])
				}
				sort.Strings(out)
				return out, true
			}
			g, valid := canon(got)
			w, _ := canon(want)
			match = valid && reflect.DeepEqual(g, w)
		case "text":
			var got, want string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil || len(got) > 1000 {
				return result, apperr.NewAppError(400, "Respuesta escrita inválida.")
			}
			match = NormalizeAnswer(got) == NormalizeAnswer(want)
		case "letters":
			var got []string
			var want string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil {
				return result, apperr.NewAppError(400, "Letras inválidas.")
			}
			labels := map[string]string{}
			for _, o := range q.Options {
				labels[o.ID] = o.Label
			}
			seen := map[string]bool{}
			word := ""
			valid := len(got) == len(q.Options)
			for _, id := range got {
				if labels[id] == "" || seen[id] {
					valid = false
				}
				word += labels[id]
				seen[id] = true
			}
			match = valid && NormalizeAnswer(word) == NormalizeAnswer(want)
		case "multi":
			var got, want []string
			if json.Unmarshal(raw, &got) != nil || json.Unmarshal(q.Solution, &want) != nil {
				return result, apperr.NewAppError(400, "Selección inválida.")
			}
			seen := map[string]bool{}
			valid := true
			for _, id := range got {
				if seen[id] {
					valid = false
				}
				seen[id] = true
			}
			sort.Strings(got)
			sort.Strings(want)
			match = valid && reflect.DeepEqual(got, want)
		default:
			return result, apperr.NewAppError(400, "Tipo de actividad no disponible.")
		}
		if match {
			correct++
		}
		result.Feedback = append(result.Feedback, models.Feedback{QuestionID: q.ID, Correct: match, Explanation: q.Explanation, Solution: q.Solution})
	}
	result.Score = correct * 100 / len(activity.Questions)
	threshold := 60
	if activity.Assessment {
		threshold = 80
	}
	result.Passed = result.Score >= threshold
	if result.Passed {
		result.Stars = 1
		if result.Score >= 80 {
			result.Stars = 2
		}
		if result.Score == 100 {
			result.Stars = 3
		}
	}
	return result, nil
}

func NormalizeAnswer(value string) string {
	value = strings.ReplaceAll(value, "’", "'")
	value = strings.TrimRight(strings.TrimSpace(value), ".!?")
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}

// NextStreak uses civil dates, not elapsed 24-hour periods (DST-safe).
func NextStreak(last, today string, current, best int) (int, int) {
	if last == today {
		return current, best
	}
	date, _ := time.Parse("2006-01-02", today)
	if last == date.AddDate(0, 0, -1).Format("2006-01-02") {
		current++
	} else {
		current = 1
	}
	if current > best {
		best = current
	}
	return current, best
}
