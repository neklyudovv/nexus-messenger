package middleware

import (
	"context"
	"testing"
	"time"
)

func TestRateLimit_AllowsUnderLimit(t *testing.T) {
	rl := newRateLimiter(context.Background(), 3, time.Minute)
	for i := range 3 {
		if !rl.allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
}

func TestRateLimit_BlocksOverLimit(t *testing.T) {
	rl := newRateLimiter(context.Background(), 3, time.Minute)
	for range 3 {
		rl.allow("1.2.3.4")
	}
	if rl.allow("1.2.3.4") {
		t.Error("4th request should be blocked")
	}
}

func TestRateLimit_IsolatesByIP(t *testing.T) {
	rl := newRateLimiter(context.Background(), 1, time.Minute)
	rl.allow("1.1.1.1")
	rl.allow("1.1.1.1") // second for 1.1.1.1 — blocked

	if !rl.allow("2.2.2.2") {
		t.Error("different IP should not be blocked")
	}
}

func TestRateLimit_ResetsAfterWindow(t *testing.T) {
	rl := newRateLimiter(context.Background(), 1, 50*time.Millisecond)
	rl.allow("1.2.3.4")
	if rl.allow("1.2.3.4") {
		t.Error("should be blocked before window expires")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.allow("1.2.3.4") {
		t.Error("should be allowed after window reset")
	}
}
