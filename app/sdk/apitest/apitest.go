// Package apitest provides support for excuting api test logic.
package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
	"github.com/owezzy/schoolCRM/business/types/role"
)

type testOption struct {
	skip    bool
	skipMsg string
}

type OptionFunc func(*testOption)

// WithSkip can be used to skip running a test.
func WithSkip(skip bool, msg string) OptionFunc {
	return func(to *testOption) {
		to.skip = skip
		to.skipMsg = msg
	}
}

// Test contains functions for executing an api test.
type Test struct {
	DB   *dbtest.Database
	Auth *auth.Auth
	mux  http.Handler
}

// Run performs the actual test logic based on the table data.
func (at *Test) Run(t *testing.T, table []Table, testName string, options ...OptionFunc) {
	to := new(testOption)
	for _, f := range options {
		f(to)
	}

	if to.skip {
		t.Skipf("%v: %v", testName, to.skipMsg)
	}

	for _, tt := range table {
		f := func(t *testing.T) {
			r := httptest.NewRequest(tt.Method, tt.URL, nil)
			w := httptest.NewRecorder()

			if tt.Input != nil {
				d, err := json.Marshal(tt.Input)
				if err != nil {
					t.Fatalf("Should be able to marshal the model : %s", err)
				}

				r = httptest.NewRequest(tt.Method, tt.URL, bytes.NewBuffer(d))
			}

			r.Header.Set("Authorization", "Bearer "+tt.Token)
			at.mux.ServeHTTP(w, r)

			if w.Code != tt.StatusCode {
				t.Fatalf("%s: Should receive a status code of %d for the response : %d", tt.Name, tt.StatusCode, w.Code)
			}

			if tt.StatusCode == http.StatusNoContent {
				return
			}

			if got, exp := w.Header().Get("Content-Type"), "application/vnd.api+json"; got != exp {
				t.Fatalf("%s: Should receive content type %q for the response : %q", tt.Name, exp, got)
			}

			body, err := unwrapJSONAPI(w.Body.Bytes(), tt.StatusCode)
			if err != nil {
				t.Fatalf("Should be able to unwrap the JSON:API response : %s", err)
			}

			if err := json.Unmarshal(body, tt.GotResp); err != nil {
				t.Fatalf("Should be able to unmarshal the response : %s", err)
			}

			diff := tt.CmpFunc(tt.GotResp, tt.ExpResp)
			if diff != "" {
				t.Log("DIFF")
				t.Logf("%s", diff)
				t.Log("GOT")
				t.Logf("%#v", tt.GotResp)
				t.Log("EXP")
				t.Logf("%#v", tt.ExpResp)
				t.Fatalf("Should get the expected response")
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

func unwrapJSONAPI(body []byte, statusCode int) ([]byte, error) {
	var doc struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Status string `json:"status"`
			Code   string `json:"code"`
			Detail string `json:"detail"`
		} `json:"errors"`
		Meta map[string]json.RawMessage `json:"meta"`
	}

	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("decode document: %w", err)
	}

	if statusCode >= http.StatusBadRequest {
		if len(doc.Errors) == 0 {
			return nil, fmt.Errorf("missing JSON:API errors")
		}

		legacy := struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    doc.Errors[0].Code,
			Message: doc.Errors[0].Detail,
		}

		data, err := json.Marshal(legacy)
		if err != nil {
			return nil, fmt.Errorf("marshal legacy error: %w", err)
		}

		return data, nil
	}

	if doc.Data == nil {
		return nil, fmt.Errorf("missing JSON:API data")
	}

	if doc.Meta != nil {
		items, err := unwrapJSONAPICollection(doc.Data)
		if err != nil {
			return nil, err
		}

		query := map[string]json.RawMessage{
			"items": items,
		}
		for key, value := range doc.Meta {
			query[key] = value
		}

		data, err := json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("marshal query result: %w", err)
		}

		return data, nil
	}

	var resource struct {
		ID         string          `json:"id"`
		Attributes json.RawMessage `json:"attributes"`
	}
	if err := json.Unmarshal(doc.Data, &resource); err != nil {
		return nil, fmt.Errorf("decode resource: %w", err)
	}

	if resource.Attributes == nil {
		return nil, fmt.Errorf("missing JSON:API resource attributes")
	}

	return mergeResourceID(resource.ID, resource.Attributes)
}

func unwrapJSONAPICollection(data json.RawMessage) (json.RawMessage, error) {
	var resources []struct {
		ID         string          `json:"id"`
		Attributes json.RawMessage `json:"attributes"`
	}
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("decode collection resources: %w", err)
	}

	items := make([]json.RawMessage, 0, len(resources))
	for _, resource := range resources {
		if resource.Attributes == nil {
			return nil, fmt.Errorf("missing collection resource attributes")
		}

		item, err := mergeResourceID(resource.ID, resource.Attributes)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	encoded, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("marshal collection items: %w", err)
	}

	return encoded, nil
}

func mergeResourceID(id string, attributes json.RawMessage) (json.RawMessage, error) {
	if id == "" {
		return attributes, nil
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(attributes, &body); err != nil {
		return nil, fmt.Errorf("decode resource attributes: %w", err)
	}

	encodedID, err := json.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("marshal resource id: %w", err)
	}
	body["id"] = encodedID

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal resource attributes: %w", err)
	}

	return data, nil
}

// =============================================================================

// Token generates an authenticated token for a user.
func Token(userBus userbus.ExtBusiness, ath *auth.Auth, email string) string {
	addr, _ := mail.ParseAddress(email)

	dbUsr, err := userBus.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    ath.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: role.ParseToString(dbUsr.Roles),
	}

	token, err := ath.GenerateToken(kid, claims)
	if err != nil {
		return ""
	}

	return token
}
