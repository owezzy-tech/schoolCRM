package authapp

import (
	"encoding/json"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-cmp/cmp"
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

func TestAppendUniqueRolesAddsAdmissionsRolesWithoutDuplicates(t *testing.T) {
	t.Parallel()

	claims := auth.Claims{
		Roles: []string{"SCHOOL_ADMIN", "ADMISSIONS_ADMIN"},
	}

	got := appendUniqueRoles(
		claims,
		[]string{"ADMISSIONS_ADMIN", "APPLICATION_REVIEWER", "REPORT_VIEWER"},
	)

	want := []string{"SCHOOL_ADMIN", "ADMISSIONS_ADMIN", "APPLICATION_REVIEWER", "REPORT_VIEWER"}
	if diff := cmp.Diff(want, got.Roles); diff != "" {
		t.Fatalf("roles mismatch (-want +got):\n%s", diff)
	}
}

func TestAuthenticateRespEncodeSupportsApplicantPortalClaims(t *testing.T) {
	t.Parallel()

	resp := authenticateResp{
		UserID: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Claims: auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Issuer:  "service project",
			},
			Roles: []string{"APPLICANT"},
			Portal: &auth.PortalClaims{
				Scope:         auth.PortalScopeApplicant,
				ApplicationID: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44",
				ConstituentID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Email:         "applicant@example.com",
			},
		},
		User: loginUser{
			ID:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			Name:  "John Applicant",
			Email: "applicant@example.com",
			Roles: []string{"APPLICANT"},
		},
	}

	data, _, err := resp.Encode()
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	var got struct {
		UserID string `json:"userID"`
		Claims struct {
			Subject string `json:"sub"`
			Portal  struct {
				Scope         string `json:"scope"`
				ApplicationID string `json:"applicationID"`
				ConstituentID string `json:"constituentID"`
			} `json:"portal"`
		} `json:"claims"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got.UserID != resp.UserID {
		t.Fatalf("userID = %q, want %q", got.UserID, resp.UserID)
	}
	if got.Claims.Subject != resp.UserID {
		t.Fatalf("claims.subject = %q, want applicant user id", got.Claims.Subject)
	}
	if got.Claims.Portal.Scope != auth.PortalScopeApplicant {
		t.Fatalf("portal.scope = %q, want %q", got.Claims.Portal.Scope, auth.PortalScopeApplicant)
	}
	if got.Claims.Portal.ApplicationID != "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44" {
		t.Fatalf("portal.applicationID = %q, want application id", got.Claims.Portal.ApplicationID)
	}
}
