package workspace

import (
	"strings"
	"testing"
)

func TestGenerateInviteCode_Length(t *testing.T) {
	code, err := generateInviteCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 16 random bytes → 32 hex characters
	if len(code) != 32 {
		t.Errorf("got len %d, want 32", len(code))
	}
}

func TestGenerateInviteCode_IsHex(t *testing.T) {
	code, _ := generateInviteCode()
	const hexChars = "0123456789abcdef"
	for _, c := range code {
		if !strings.ContainsRune(hexChars, c) {
			t.Errorf("non-hex character %q in invite code %q", c, code)
		}
	}
}

func TestGenerateInviteCode_Unique(t *testing.T) {
	seen := make(map[string]struct{})
	for range 100 {
		code, err := generateInviteCode()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, dup := seen[code]; dup {
			t.Errorf("duplicate invite code generated: %q", code)
		}
		seen[code] = struct{}{}
	}
}
