package httpapi

import (
	"testing"
	"time"
)

func TestLoginLimiter(t *testing.T) {
	limiter := newLoginLimiter(2, time.Minute)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }

	if !limiter.allow("client") || !limiter.allow("client") {
		t.Fatal("expected initial attempts to be allowed")
	}
	if limiter.allow("client") {
		t.Fatal("expected limit to reject the third attempt")
	}

	limiter.reset("client")
	if !limiter.allow("client") {
		t.Fatal("expected reset to allow another attempt")
	}
}
