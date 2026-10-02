//go:build integration

package lessonapp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestLessonAPI(t *testing.T) {
	db := dbtest.New(t, "lesson-api")
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
	root, teacher, hod, dean, outsider := newUser("SUPER_ADMIN"), newUser("TEACHER"), newUser("TEACHER"), newUser("TEACHER"), newUser("TEACHER")
	school, err := db.BusDomain.School.CreateSchool(ctx, root, "School One")
	if err != nil {
		t.Fatal(err)
	}
	department, err := db.BusDomain.School.CreateDepartment(ctx, root, school.ID, "Sciences")
	if err != nil {
		t.Fatal(err)
	}
	for user, capability := range map[uuid.UUID]schoolbus.Capability{teacher: schoolbus.Teach, hod: schoolbus.ReviewLessons, dean: schoolbus.ApproveLessons} {
		if _, err := db.BusDomain.School.Grant(ctx, root, user, schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, capability); err != nil {
			t.Fatal(err)
		}
	}

	tracer := noop.NewTracerProvider().Tracer("test")
	identity := auth.New(auth.Config{Log: db.Log, UserBus: db.BusDomain.User, KeyLookup: &apitest.KeyStore{}, Issuer: "lesson-test"})
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
		BusConfig:       mux.BusConfig{LessonBus: db.BusDomain.Lesson, SchoolBus: db.BusDomain.School, UserBus: db.BusDomain.User, AuditBus: db.BusDomain.Audit},
		SchoolCRMConfig: mux.SchoolCRMConfig{AuthClient: client}}, schoolbuild.Routes())
	token := func(id uuid.UUID) string {
		t.Helper()
		value, err := identity.GenerateToken("test", auth.Claims{RegisteredClaims: jwt.RegisteredClaims{
			Subject: id.String(), Issuer: "lesson-test", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}, Roles: []string{"TEACHER"}})
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	request := func(method, path string, user uuid.UUID, body string, status int) map[string]json.RawMessage {
		t.Helper()
		r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		if user != uuid.Nil {
			r.Header.Set("Authorization", "Bearer "+token(user))
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != status {
			t.Fatalf("%s %s: status %d, want %d: %s", method, path, w.Code, status, w.Body.String())
		}
		if w.Header().Get("Content-Type") != "application/vnd.api+json" {
			t.Fatalf("content type %s", w.Header().Get("Content-Type"))
		}
		var doc map[string]json.RawMessage
		if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
			t.Fatal(err)
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
	type resource struct {
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		Attributes map[string]any `json:"attributes"`
	}
	single := func(doc map[string]json.RawMessage) resource {
		t.Helper()
		var r resource
		if err := json.Unmarshal(doc["data"], &r); err != nil || r.ID == "" {
			t.Fatalf("invalid resource: %s", doc["data"])
		}
		return r
	}

	lessons := fmt.Sprintf("/v1/schools/%s/departments/%s/lessons", school.ID, department.ID)
	t.Run("authentication and boundary validation", func(t *testing.T) {
		request("GET", lessons, uuid.Nil, "", 401)
		request("POST", lessons, teacher, `{"title":"x","content":[]}`, 400)
		request("POST", lessons, teacher, `{"title":"x","content":{},"baseVersion":1}`, 400)
		request("POST", lessons, teacher, `{"title":"x","content":{"a":"`+strings.Repeat("a", 300*1024)+`"}}`, 400)
		request("POST", "/v1/lessons/"+uuid.NewString()+"/versions/0/submit", teacher, "", 400)
		request("POST", lessons, hod, `{"title":"x","content":{}}`, 403)
	})

	plan := single(request("POST", lessons, teacher, `{"title":"Forces","content":{"objectives":["Newton"]},"changeSummary":"first"}`, 200))
	if plan.Type != "lessonplan" || plan.Attributes["status"] != "draft" {
		t.Fatalf("got %+v", plan)
	}
	version := "/v1/lessons/" + plan.ID + "/versions/1"

	t.Run("workflow", func(t *testing.T) {
		request("POST", version+"/publish", teacher, "", 409)
		request("POST", version+"/submit", teacher, "", 200)
		request("POST", version+"/review", hod, `{"decision":"approve","feedback":"clear"}`, 200)
		request("POST", version+"/approval", hod, `{"decision":"approve"}`, 403)
		request("POST", version+"/approval", dean, `{"decision":"maybe"}`, 400)
		request("POST", version+"/approval", dean, `{"decision":"approve"}`, 200)
		published := single(request("POST", version+"/publish", teacher, "", 200))
		if published.Type != "lessonversion" || published.ID != plan.ID+":1" || published.Attributes["status"] != "published" {
			t.Fatalf("got %+v", published)
		}
		request("POST", version+"/publish", teacher, "", 200)
	})

	t.Run("revision keeps publication live", func(t *testing.T) {
		request("POST", "/v1/lessons/"+plan.ID+"/versions", teacher, `{"baseVersion":1,"title":"Forces v2","content":{"objectives":["Newton","Hooke"]}}`, 200)
		request("POST", "/v1/lessons/"+plan.ID+"/versions", teacher, `{"baseVersion":1,"title":"Stale","content":{}}`, 409)
		got := single(request("GET", "/v1/lessons/"+plan.ID, dean, "", 200))
		if got.Attributes["publishedVersion"] != float64(1) || got.Attributes["currentVersion"] != float64(2) {
			t.Fatalf("got %+v", got.Attributes)
		}
		doc := request("GET", "/v1/lessons/"+plan.ID+"/versions", hod, "", 200)
		var versions []resource
		if err := json.Unmarshal(doc["data"], &versions); err != nil || len(versions) != 2 || doc["meta"] == nil {
			t.Fatalf("invalid versions collection: %s", doc["data"])
		}
		doc = request("GET", lessons, teacher, "", 200)
		var plans []resource
		if err := json.Unmarshal(doc["data"], &plans); err != nil || len(plans) != 1 || plans[0].ID != plan.ID {
			t.Fatalf("invalid plans collection: %s", doc["data"])
		}
	})

	t.Run("non-members are denied", func(t *testing.T) {
		request("GET", "/v1/lessons/"+plan.ID, outsider, "", 403)
		request("GET", lessons, outsider, "", 403)
		request("POST", "/v1/lessons/"+plan.ID+"/versions/2/submit", outsider, "", 403)
		request("GET", "/v1/lessons/"+uuid.NewString(), teacher, "", 404)
	})
}
