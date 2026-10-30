package models

import "encoding/json"

type Option struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Image    string `json:"image,omitempty"`
	ImageURL string `json:"image_url,omitempty" gorm:"column:image_url;->"`
}

// Solution and Explanation are private until an attempt is completed.
type Question struct {
	ID          string          `json:"id"`
	Prompt      string          `json:"prompt"`
	Options     []Option        `json:"options" gorm:"-"`
	Solution    json.RawMessage `json:"solution,omitempty" gorm:"-"`
	Explanation string          `json:"explanation,omitempty"`
	Response    string          `json:"response,omitempty"`
	Passage     string          `json:"passage,omitempty"`
	AudioKey    string          `json:"audio_key,omitempty"`
	AudioURL    string          `json:"audio_url,omitempty"`
}

type Activity struct {
	ID               string     `json:"id"`
	LessonID         string     `json:"lesson_id"`
	Title            string     `json:"title"`
	Kind             string     `json:"kind"`
	Instructions     string     `json:"instructions"`
	Assessment       bool       `json:"assessment"`
	Published        bool       `json:"published"`
	Questions        []Question `json:"questions" gorm:"-"`
	TimeLimitSeconds int        `json:"time_limit_seconds"`
	VocabularyCount  int        `json:"vocabulary_count"`
	Skill            string     `json:"skill"`
}

type Submission struct {
	Answers map[string]json.RawMessage `json:"answers"`
}

type Feedback struct {
	QuestionID  string          `json:"question_id"`
	Correct     bool            `json:"correct"`
	Explanation string          `json:"explanation"`
	Solution    json.RawMessage `json:"solution"`
}

type Result struct {
	AttemptID      string     `json:"attempt_id"`
	ActivityID     string     `json:"activity_id"`
	Score          int        `json:"score"`
	Passed         bool       `json:"passed"`
	Stars          int        `json:"stars"`
	XP             int        `json:"xp"`
	NewCompletion  bool       `json:"new_completion"`
	Feedback       []Feedback `json:"feedback"`
	Expired        bool       `json:"expired"`
	ElapsedSeconds int        `json:"elapsed_seconds"`
	Multiplier     int        `json:"multiplier"`
}

func (q Question) AnswerType(kind string) string {
	if q.Response != "" {
		return q.Response
	}
	switch kind {
	case "sentence_builder":
		return "order"
	case "memory_match":
		return "pairs"
	case "grammar_rescue":
		return "text"
	case "spell_master":
		return "letters"
	case "word_rain":
		return "multi"
	default:
		return "choice"
	}
}

func (a Activity) LearningSkill() string {
	if a.Skill != "" {
		return a.Skill
	}
	switch a.Kind {
	case "sentence_builder", "grammar_rescue":
		return "grammar"
	case "sound_detective":
		return "listening"
	case "conversation_quest", "story_explorer":
		return "reading"
	case "spell_master":
		return "writing"
	case "english_rush":
		return "mixed"
	default:
		return "vocabulary"
	}
}

// Only the base activity reward is doubled. Perfect and mission bonuses never stack.
func (a Activity) RewardXP(perfect bool, multiplier int) int {
	base := 20
	if a.Assessment {
		base = 40
	}
	if multiplier != 2 {
		multiplier = 1
	}
	xp := base * multiplier
	if perfect {
		xp += 10
	}
	return xp
}

type LessonSummary struct {
	ID              string `json:"id"`
	UnitID          string `json:"unit_id"`
	UnitTitle       string `json:"unit_title"`
	UnitDescription string `json:"unit_description"`
	HasProject      bool   `json:"has_project"`
	LevelCode       string `json:"level_code"`
	Title           string `json:"title"`
	Objective       string `json:"objective"`
	Content         string `json:"content"`
	ActivityID      string `json:"activity_id"`
	Kind            string `json:"kind"`
	Assessment      bool   `json:"assessment"`
	Completed       bool   `json:"completed"`
	BestScore       int    `json:"best_score"`
	Locked          bool   `json:"locked"`
	LevelLocked     bool   `json:"level_locked"`
}
