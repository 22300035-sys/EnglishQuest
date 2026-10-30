package services

import (
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

func DiplomaName(input string) (string, error) {
	bad := func() (string, error) {
		return "", apperr.NewAppError(400, "Escribe un nombre de 2 a 120 caracteres usando letras, espacios, guiones, puntos o apóstrofos.")
	}
	if utf8.RuneCountInString(input) > 120 {
		return bad()
	}
	letter := false
	for _, r := range input {
		if unicode.IsControl(r) {
			return bad()
		}
		if unicode.IsLetter(r) {
			letter = true
			continue
		}
		if unicode.IsMark(r) || r == ' ' || r == '-' || r == '.' || r == '\'' || r == '’' {
			continue
		}
		return bad()
	}
	name := strings.Join(strings.Fields(input), " ")
	if !letter || utf8.RuneCountInString(name) < 2 {
		return bad()
	}
	return name, nil
}
