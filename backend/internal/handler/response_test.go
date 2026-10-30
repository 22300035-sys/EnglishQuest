package handler

import (
	"context"
	"errors"
	"fmt"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondError(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
		text   string
	}{
		{fmt.Errorf("wrapped: %w", apperr.NewAppError(409, "conflict")), 409, "conflict"},
		{errors.New("secret database details"), 500, "internal server error"},
	} {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		RespondError(c, tc.err)
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.text) || strings.Contains(w.Body.String(), "secret") {
			t.Fatalf("unexpected response %s", w.Body.String())
		}
	}
}

func TestCancelledRequestIsNotServerFailure(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	RespondError(c, context.Canceled)
	if w.Code != 499 {
		t.Fatalf("client cancellation returned %d", w.Code)
	}
}
