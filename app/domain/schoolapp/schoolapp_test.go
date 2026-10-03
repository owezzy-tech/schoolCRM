//go:build integration

package schoolapp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	authbuild "github.com/owezzy/schoolCRM/api/services/auth/build"
	schoolbuild "github.com/owezzy/schoolCRM/api/services/schoolcrm/build"
	"github.com/owezzy/schoolCRM/app/sdk/apitest"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	authhttp "github.com/owezzy/schoolCRM/app/sdk/authclient/http"
	"github.com/owezzy/schoolCRM/app/sdk/mux"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus/stores/schooldb"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestSchoolAPI(t *testing.T) {
	db := dbtest.New(t, "school-api")
	ctx := context.Background()
	newUser := func(role string) uuid.UUID {
		t.Helper()
		id := uuid.New()
		if _, err := db.DB.ExecContext(ctx, `INSERT INTO users (user_id,name,email,password_hash,roles,enabled,date_created,date_updated)
		VALUES ($1,'API User',$2,'unused',$3::text[],true,now(),now())`, id, id.String()+"@test.invalid", "{"+role+"}"); err != nil {
			t.Fatal(err)
		}
		return id
	}
	root, admin, teacher := newUser("SUPER_ADMIN"), newUser("SCHOOL_ADMIN"), newUser("TEACHER")
	tracer := noop.NewTracerProvider().Tracer("test")
	identity := auth.New(auth.Config{Log: db.Log, UserBus: db.BusDomain.User, KeyLookup: &apitest.KeyStore{}, Issuer: "school-test"})
	authServer := httptest.NewServer(mux.WebAPI(mux.Config{Log: db.Log, DB: db.DB, Tracer: tracer,
		BusConfig: mux.BusConfig{UserBus: db.BusDomain.User}, AuthConfig: mux.AuthConfig{Auth: identity, TokenKey: "test"}}, authbuild.Routes()))
	t.Cleanup(authServer.Close)
	client, err := authhttp.New(db.Log, authServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Error(err)
		}
	})
	handler := mux.WebAPI(mux.Config{Log: db.Log, DB: db.DB, Tracer: tracer,
		BusConfig:       mux.BusConfig{SchoolBus: schoolbus.NewBusiness(schooldb.NewStore(db.Log, db.DB)), UserBus: db.BusDomain.User, AuditBus: db.BusDomain.Audit},
		SchoolCRMConfig: mux.SchoolCRMConfig{AuthClient: client}}, schoolbuild.Routes())
	token := func(id uuid.UUID, roles ...string) string {
		t.Helper()
		value, err := identity.GenerateToken("test", auth.Claims{RegisteredClaims: jwt.RegisteredClaims{
			Subject: id.String(), Issuer: "school-test", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}, Roles: roles})
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	rootToken, adminToken, teacherToken := token(root, "SUPER_ADMIN"), token(admin, "SCHOOL_ADMIN"), token(teacher, "TEACHER")
	request := func(method, path, bearer, body string, status int) map[string]json.RawMessage {
		t.Helper()
		r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		if bearer != "" {
			r.Header.Set("Authorization", "Bearer "+bearer)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != status {
			t.Fatalf("%s %s: status %d, want %d: %s", method, path, w.Code, status, w.Body.String())
		}
		if status == http.StatusNoContent {
			return nil
		}
		if w.Header().Get("Content-Type") != "application/vnd.api+json" {
			t.Fatalf("content type %s", w.Header().Get("Content-Type"))
		}
		var doc map[string]json.RawMessage
		if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
			t.Fatal(err)
		}
		var version struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(doc["jsonapi"], &version); err != nil || version.Version != "1.1" {
			t.Fatalf("invalid version: %s", w.Body.String())
		}
		if status >= 400 {
			var failures []struct {
				Detail string `json:"detail"`
			}
			if err := json.Unmarshal(doc["errors"], &failures); err != nil || len(failures) == 0 || failures[0].Detail == "" {
				t.Fatalf("missing safe error: %s", w.Body.String())
			}
		} else if doc["data"] == nil {
			t.Fatal("missing data")
		}
		return doc
	}
	resourceID := func(doc map[string]json.RawMessage) string {
		t.Helper()
		var resource struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(doc["data"], &resource); err != nil || resource.ID == "" {
			t.Fatalf("missing resource ID: %s", doc["data"])
		}
		return resource.ID
	}

	t.Run("authentication and global roles", func(t *testing.T) {
		request("GET", "/v1/schools", "", "", 401)
		request("POST", "/v1/schools", adminToken, `{"name":"Denied"}`, 403)
		request("POST", "/v1/schools", token(teacher, "SUPER_ADMIN"), `{"name":"Stale claims"}`, 403)
	})
	schoolID := resourceID(request("POST", "/v1/schools", rootToken, `{"name":"School One"}`, 200))
	otherID := resourceID(request("POST", "/v1/schools", rootToken, `{"name":"School Two"}`, 200))
	scopePath := "/v1/schools/" + schoolID
	departmentID := resourceID(request("POST", scopePath+"/departments", rootToken, `{"name":"Sciences"}`, 200))
	managerID := resourceID(request("POST", scopePath+"/memberships", rootToken, `{"userID":"`+admin.String()+`","capability":"manage_members"}`, 200))

	t.Run("scoped administration and collection shape", func(t *testing.T) {
		request("POST", scopePath+"/departments", adminToken, `{"name":"Arts"}`, 200)
		request("GET", "/v1/schools/"+otherID+"/memberships", adminToken, "", 403)
		request("POST", scopePath+"/memberships", adminToken, `{"userID":"`+admin.String()+`","capability":"manage_members"}`, 403)
		request("DELETE", scopePath+"/memberships/"+managerID, adminToken, "", 403)
		doc := request("GET", "/v1/schools", adminToken, "", 200)
		var resources []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(doc["data"], &resources); err != nil {
			t.Fatal(err)
		}
		if len(resources) != 1 || resources[0].ID != schoolID || doc["meta"] == nil {
			t.Fatalf("invalid scoped collection: %s", doc["data"])
		}
	})
	t.Run("grant and revoke", func(t *testing.T) {
		body := `{"userID":"` + teacher.String() + `","departmentID":"` + departmentID + `","capability":"teach"}`
		memberID := resourceID(request("POST", scopePath+"/memberships", adminToken, body, 200))
		if duplicate := resourceID(request("POST", scopePath+"/memberships", adminToken, body, 200)); duplicate != memberID {
			t.Fatal("duplicate membership")
		}
		request("GET", scopePath+"/departments", teacherToken, "", 200)
		own := func() []struct {
			ID         string `json:"id"`
			Attributes struct {
				Capability string `json:"capability"`
			} `json:"attributes"`
		} {
			var items []struct {
				ID         string `json:"id"`
				Attributes struct {
					Capability string `json:"capability"`
				} `json:"attributes"`
			}
			if err := json.Unmarshal(request("GET", "/v1/me/school-memberships", teacherToken, "", 200)["data"], &items); err != nil {
				t.Fatal(err)
			}
			return items
		}
		if items := own(); len(items) != 1 || items[0].ID != memberID || items[0].Attributes.Capability != "teach" {
			t.Fatalf("own memberships: %+v", items)
		}
		request("DELETE", "/v1/schools/"+otherID+"/memberships/"+memberID, rootToken, "", 404)
		request("DELETE", scopePath+"/memberships/"+memberID, adminToken, "", 204)
		request("DELETE", scopePath+"/memberships/"+memberID, adminToken, "", 204)
		request("GET", scopePath+"/departments", teacherToken, "", 403)
		if items := own(); len(items) != 0 {
			t.Fatalf("revoked membership still listed: %+v", items)
		}
	})
	t.Run("boundary validation", func(t *testing.T) {
		cases := []struct{ name, path, body string }{
			{"nil school", "/v1/schools/" + uuid.Nil.String() + "/departments", `{"name":"Invalid"}`},
			{"unknown fields", "/v1/schools", `{"name":"Invalid","actorID":"` + root.String() + `"}`},
			{"multiple objects", "/v1/schools", `{"name":"Invalid"}{"name":"Other"}`},
			{"oversize payload", "/v1/schools", `{"name":"` + strings.Repeat("a", 17000) + `"}`},
			{"empty name", "/v1/schools", `{"name":" "}`},
			{"unknown capability", scopePath + "/memberships", `{"userID":"` + teacher.String() + `","capability":"root"}`},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) { request("POST", tc.path, rootToken, tc.body, 400) })
		}
	})
	t.Run("delegation revocation overrides existing token", func(t *testing.T) {
		request("DELETE", scopePath+"/memberships/"+managerID, rootToken, "", 204)
		request("POST", scopePath+"/departments", adminToken, `{"name":"Denied"}`, 403)
	})
}
