package auth_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

func TestFailedAuthenticationDoesNotLogToken(t *testing.T) {
	var logs bytes.Buffer
	log := logger.New(&logs, logger.LevelInfo, "TEST", func(context.Context) string { return "proof" })
	ath := auth.New(auth.Config{Log: log, KeyLookup: &keyStore{}, Issuer: "expected-issuer"})
	token, err := ath.GenerateToken(kid, auth.Claims{RegisteredClaims: jwt.RegisteredClaims{
		Issuer: "wrong-issuer", Subject: "15953eb5-2b0a-490b-b193-cbb2f3d90bc4",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}})
	if err != nil {
		t.Fatal("fixture token generation failed")
	}
	if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
		t.Fatal("wrong issuer was accepted")
	}
	if strings.Contains(logs.String(), token) {
		t.Fatal("authentication failure logged the bearer credential")
	}
	if !strings.Contains(logs.String(), "Authenticate-FAILED") {
		t.Fatal("credential-safe failure diagnostic missing")
	}
}
