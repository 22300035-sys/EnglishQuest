package services

import (
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"sort"
	"strings"
)

var Images = map[string]bool{"apple": true, "banana": true, "orange": true, "book": true, "pen": true, "school": true, "breakfast": true, "night": true}
var GameKinds = map[string]bool{"word_hunter": true, "sentence_builder": true, "memory_match": true, "sound_detective": true, "grammar_rescue": true, "spell_master": true, "conversation_quest": true, "story_explorer": true, "word_rain": true, "english_rush": true}

// ValidateActivity validates admin input before any content changes are committed.
func ValidateActivity(a models.Activity) error {
	bad := func() error {
		return apperr.NewAppError(400, "Contenido inválido: revisa título, instrucciones, opciones y soluciones.")
	}
	if len(strings.TrimSpace(a.Title)) < 1 || len(a.Title) > 150 || len(a.Instructions) < 5 || len(a.Instructions) > 2000 || len(a.Questions) < 1 || len(a.Questions) > 12 {
		return bad()
	}
	if !GameKinds[a.Kind] {
		return bad()
	}
	if a.VocabularyCount < 0 || a.VocabularyCount > 100 {
		return bad()
	}
	skills := map[string]bool{"vocabulary": true, "grammar": true, "listening": true, "reading": true, "writing": true, "mixed": true}
	if !skills[a.LearningSkill()] {
		return bad()
	}
	if a.Kind == "word_rain" || a.Kind == "english_rush" {
		if a.TimeLimitSeconds < 30 || a.TimeLimitSeconds > 300 {
			return bad()
		}
	} else if a.TimeLimitSeconds != 0 {
		return bad()
	}
	for _, q := range a.Questions {
		answerType := q.AnswerType(a.Kind)
		if a.Kind != "english_rush" && q.Response != "" && q.Response != (models.Question{}).AnswerType(a.Kind) {
			return bad()
		}
		if a.Kind == "english_rush" && answerType != "choice" && answerType != "text" && answerType != "order" && answerType != "letters" {
			return bad()
		}
		if len(q.Prompt) < 1 || len(q.Prompt) > 500 || len(q.Explanation) < 1 || len(q.Explanation) > 1500 || len(q.Options) > 12 || len(q.Passage) > 5000 {
			return bad()
		}
		if answerType == "text" {
			if len(q.Options) != 0 {
				return bad()
			}
		} else if len(q.Options) < 2 {
			return bad()
		}
		if a.Kind == "sound_detective" && q.AudioKey == "" {
			return bad()
		}
		if q.AudioKey != "" && !ResourceKey.MatchString(q.AudioKey) {
			return bad()
		}
		if q.AudioKey != "" && answerType != "choice" {
			return bad()
		}
		if (a.Kind == "story_explorer" || a.Kind == "conversation_quest") && len(q.Passage) < 10 {
			return bad()
		}
		ids := map[string]bool{}
		for _, o := range q.Options {
			if len(o.ID) < 1 || len(o.ID) > 40 || strings.Contains(o.ID, ":") || ids[o.ID] || len(o.Label) < 1 || len(o.Label) > 120 {
				return bad()
			}
			if o.Image != "" && !ResourceKey.MatchString(o.Image) {
				return bad()
			}
			ids[o.ID] = true
			if a.Kind == "word_hunter" && o.Image == "" {
				return bad()
			}
		}
		switch answerType {
		case "choice":
			var key string
			if json.Unmarshal(q.Solution, &key) != nil || !ids[key] {
				return bad()
			}
		case "order":
			var key []string
			if json.Unmarshal(q.Solution, &key) != nil || len(key) != len(ids) {
				return bad()
			}
			seen := map[string]bool{}
			for _, v := range key {
				if !ids[v] || seen[v] {
					return bad()
				}
				seen[v] = true
			}
		case "pairs":
			var pairs [][]string
			if json.Unmarshal(q.Solution, &pairs) != nil || len(pairs)*2 != len(ids) {
				return bad()
			}
			seen := map[string]bool{}
			for _, p := range pairs {
				if len(p) != 2 {
					return bad()
				}
				for _, v := range p {
					if !ids[v] || seen[v] {
						return bad()
					}
					seen[v] = true
				}
			}
		case "text":
			var key string
			if json.Unmarshal(q.Solution, &key) != nil || len(strings.TrimSpace(key)) < 1 || len(key) > 500 {
				return bad()
			}
		case "letters":
			var key string
			if json.Unmarshal(q.Solution, &key) != nil || len(key) != len(q.Options) {
				return bad()
			}
			letters := []string{}
			for _, option := range q.Options {
				if len(option.Label) != 1 || !((option.Label[0] >= 'a' && option.Label[0] <= 'z') || (option.Label[0] >= 'A' && option.Label[0] <= 'Z')) {
					return bad()
				}
				letters = append(letters, strings.ToLower(option.Label))
			}
			sort.Strings(letters)
			want := strings.Split(strings.ToLower(key), "")
			sort.Strings(want)
			if strings.Join(letters, "") != strings.Join(want, "") {
				return bad()
			}
		case "multi":
			var key []string
			if json.Unmarshal(q.Solution, &key) != nil || len(key) < 1 || len(key) >= len(ids) {
				return bad()
			}
			seen := map[string]bool{}
			for _, value := range key {
				if !ids[value] || seen[value] {
					return bad()
				}
				seen[value] = true
			}
		default:
			return bad()
		}
	}
	return nil
}
