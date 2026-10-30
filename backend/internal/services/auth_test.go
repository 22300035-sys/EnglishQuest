package services

import (
	"context"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"strings"
	"testing"
)

func TestRegistrationRejectsInvalidDataBeforePersistence(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*Registration)
	}{
		{"email", func(in *Registration) { in.Email = "not-an-email" }},
		{"display email", func(in *Registration) { in.Email = "Name <user@example.org>" }},
		{"short password", func(in *Registration) { in.Password = "short" }},
		{"long password", func(in *Registration) { in.Password = strings.Repeat("x", 73) }},
		{"empty name", func(in *Registration) { in.Name = "   " }},
		{"timezone", func(in *Registration) { in.Timezone = "Invalid/Timezone" }},
		{"server local timezone", func(in *Registration) { in.Timezone = "Local" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := Registration{Email: "user@example.org", Password: "long-test-password", Name: "Test", Timezone: "UTC"}
			tc.change(&in)
			_, err := (Auth{}).Register(context.Background(), in)
			if err == nil || apperr.HTTPStatus(err) != 400 {
				t.Fatalf("invalid data not rejected: %v", err)
			}
		})
	}
}
