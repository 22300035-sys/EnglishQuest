package services

import (
	"github.com/englishquest/backend/internal/models"
	"testing"
)

func TestResourceSafety(t *testing.T) {
	for _, tc := range []struct {
		url   string
		valid bool
	}{
		{"/images/apple.svg", true}, {"/images/nonexistent.svg", false}, {"/images/../../secret.svg", false},
		{"https://example.org/authorized.png", true}, {"http://example.org/image.png", false},
		{"https://127.0.0.1/image.png", false}, {"https://10.0.0.1/image.png", false},
		{"https://user:password@example.org/image.png", false}, {"javascript:alert(1)", false},
	} {
		r := models.Resource{Key: "test", Title: "Test resource", Kind: "image", URL: tc.url, Attribution: "Original artwork", Active: true}
		if (ValidateResource(r) == nil) != tc.valid {
			t.Fatalf("unexpected URL validation: %s", tc.url)
		}
	}
}
