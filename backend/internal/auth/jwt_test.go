package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseToken_Valid(t *testing.T) {
	secret := "test-secret-key-32-bytes-longggg"
	token, err := generateToken(42, time.Hour, secret)
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("got UserID %d, want 42", claims.UserID)
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, _ := generateToken(1, time.Hour, "secret-a")
	_, err := ParseToken(token, "secret-b")
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, _ := generateToken(1, -time.Minute, "secret") // already expired
	_, err := ParseToken(token, "secret")
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestParseToken_AlgorithmValidation(t *testing.T) {
	// RS256-signed token must be rejected even if the payload looks valid.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		UserID:           1,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	})
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign RS256 token: %v", err)
	}

	_, err = ParseToken(signed, "any-hmac-secret")
	if err == nil {
		t.Error("expected error for non-HMAC algorithm, got nil")
	}
}

func TestParseToken_Malformed(t *testing.T) {
	_, err := ParseToken("not.a.jwt", "secret")
	if err == nil {
		t.Error("expected error for malformed token, got nil")
	}
}
