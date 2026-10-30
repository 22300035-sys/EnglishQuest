package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/englishquest/backend/curriculum"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"math/big"
)

const AccessExamSeconds = 1200
const AccessExamQuestions = 12

// BuildAccessExam selects a fresh private snapshot; no bank keys reach the browser.
func BuildAccessExam(level string) (models.Activity, error) {
	var bank map[string][][]json.RawMessage
	a := models.Activity{Title: "Examen de acceso a " + level, Kind: "conversation_quest", Assessment: true, Published: true, Skill: "mixed", Instructions: "Responde 12 preguntas en 20 minutos. Se necesita al menos 80% (10 aciertos). Este examen solo abre el nivel: no completa actividades ni acredita un nivel CEFR."}
	if err := json.Unmarshal(curriculum.AccessExams, &bank); err != nil {
		return a, err
	}
	items, ok := bank[level]
	if !ok {
		return a, apperr.NewAppError(404, "Este nivel no tiene un examen de acceso.")
	}
	for i := len(items) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return a, err
		}
		j := int(n.Int64())
		items[i], items[j] = items[j], items[i]
	}
	if len(items) < AccessExamQuestions {
		return a, apperr.NewAppError(409, "Banco de examen incompleto.")
	}
	for _, item := range items[:AccessExamQuestions] {
		var prompt, passage, explanation string
		var labels []string
		var answer int
		if len(item) != 5 {
			return a, apperr.NewAppError(500, "Banco de examen inválido.")
		}
		targets := []interface{}{&prompt, &passage, &labels, &answer, &explanation}
		for i := range targets {
			if err := json.Unmarshal(item[i], targets[i]); err != nil {
				return a, err
			}
		}
		if answer < 0 || answer >= len(labels) {
			return a, apperr.NewAppError(500, "Solución de examen inválida.")
		}
		bytes := make([]byte, 12)
		if _, err := rand.Read(bytes); err != nil {
			return a, err
		}
		if passage == "" {
			passage = "Read the item carefully and choose the most accurate answer in this context."
		}
		q := models.Question{ID: hex.EncodeToString(bytes), Prompt: prompt, Passage: passage, Explanation: explanation, Options: []models.Option{}}
		for i, label := range labels {
			q.Options = append(q.Options, models.Option{ID: hex.EncodeToString([]byte{byte(i + 1)}), Label: label})
		}
		key, _ := json.Marshal(q.Options[answer].ID)
		q.Solution = key
		a.Questions = append(a.Questions, q)
	}
	if err := ValidateActivity(a); err != nil {
		return a, err
	}
	a.TimeLimitSeconds = AccessExamSeconds
	return a, nil
}
