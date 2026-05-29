package authapp

import (
	"encoding/json"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
)

func TestAuthenticateRespEncodeIncludesSanitizedUser(t *testing.T) {
	t.Parallel()

	resp := authenticateResp{
		UserID: "5cf37266-3473-4006-984f-9325122678b7",
		Claims: auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "5cf37266-3473-4006-984f-9325122678b7",
				Issuer:  "service project",
			},
			Roles: []string{"SCHOOL_ADMIN"},
		},
		User: loginUser{
			ID:    "5cf37266-3473-4006-984f-9325122678b7",
			Name:  "Admin Gopher",
			Email: "admin@example.com",
			Roles: []string{"SCHOOL_ADMIN"},
		},
	}

	data, contentType, err := resp.Encode()
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	if contentType != "application/json" {
		t.Fatalf("contentType = %q, want application/json", contentType)
	}

	var got struct {
		UserID string `json:"userID"`
		Claims struct {
			Subject string   `json:"sub"`
			Issuer  string   `json:"iss"`
			Roles   []string `json:"roles"`
		} `json:"claims"`
		User loginUser `json:"user"`
	}

	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got.UserID != resp.UserID {
		t.Fatalf("userID = %q, want %q", got.UserID, resp.UserID)
	}

	if got.Claims.Subject != resp.UserID {
		t.Fatalf("claims.subject = %q, want %q", got.Claims.Subject, resp.UserID)
	}

	if len(got.Claims.Roles) != 1 || got.Claims.Roles[0] != "SCHOOL_ADMIN" {
		t.Fatalf("claims.roles = %v, want [SCHOOL_ADMIN]", got.Claims.Roles)
	}

	if got.User.Email != "admin@example.com" {
		t.Fatalf("user.email = %q, want admin@example.com", got.User.Email)
	}
}
