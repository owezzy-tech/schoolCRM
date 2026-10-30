//go:build integration

package lessondb_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
	"sync"
	"testing"
)

func TestCreationRetriesAndReuse(t *testing.T) {
	db := dbtest.New(t, "creation-retries")
	ctx := context.Background()
	root, author, reader, outsider := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	for _, id := range []uuid.UUID{root, author, reader, outsider} {
		roles := "{TEACHER}"
		if id == root {
			roles = "{SUPER_ADMIN}"
		}
		if _, err := db.DB.ExecContext(ctx, `INSERT INTO users(user_id,name,email,password_hash,roles,enabled,date_created,date_updated)
            VALUES($1,'Test',$2,'unused',$3::text[],true,now(),now())`, id, id.String()+"@test.invalid", roles); err != nil {
			t.Fatal(err)
		}
	}
	school, err := db.BusDomain.School.CreateSchool(ctx, root, "Retry School")
	if err != nil {
		t.Fatal(err)
	}
	department, err := db.BusDomain.School.CreateDepartment(ctx, root, school.ID, "Languages")
	if err != nil {
		t.Fatal(err)
	}
	memberships := map[uuid.UUID]schoolbus.Membership{}
	for _, actor := range []uuid.UUID{author, reader} {
		m, err := db.BusDomain.School.Grant(ctx, root, actor,
			schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, schoolbus.Teach)
		if err != nil {
			t.Fatal(err)
		}
		memberships[actor] = m
	}
	draft := lessonbus.Revision{Title: "School vocabulary", Content: json.RawMessage(`{"objectives":["Name objects"],"citations":[{"page":17}]}`), RequestID: uuid.New()}
	var wg sync.WaitGroup
	ids := make(chan uuid.UUID, 8)
	failures := make(chan error, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, e := db.BusDomain.Lesson.Create(ctx, author, school.ID, department.ID, draft)
			ids <- p.ID
			failures <- e
		}()
	}
	wg.Wait()
	close(ids)
	close(failures)
	for err := range failures {
		if err != nil {
			t.Fatal(err)
		}
	}
	var planID uuid.UUID
	for id := range ids {
		if planID == uuid.Nil {
			planID = id
		}
		if planID != id {
			t.Fatal("retry created different plan")
		}
	}
	var count int
	if err := db.DB.GetContext(ctx, &count, `SELECT count(*) FROM lesson_plan_versions WHERE plan_id=$1`, planID); err != nil || count != 1 {
		t.Fatalf("versions=%d err=%v", count, err)
	}
	if err := db.DB.GetContext(ctx, &count, `SELECT count(*) FROM audit WHERE obj_id=$1`, planID); err != nil || count != 1 {
		t.Fatalf("audits=%d err=%v", count, err)
	}
	changed := draft
	canonical := draft
	canonical.Content = json.RawMessage(`{ "citations": [ { "page": 17 } ], "objectives": ["Name objects"] }`)
	replayed, err := db.BusDomain.Lesson.Create(ctx, author, school.ID, department.ID, canonical)
	if err != nil || replayed.ID != planID {
		t.Fatalf("canonical replay: %v", err)
	}
	changed.Title = "Changed input"
	if _, err := db.BusDomain.Lesson.Create(ctx, author, school.ID, department.ID, changed); !errors.Is(err, lessonbus.ErrConflict) {
		t.Fatalf("changed request accepted: %v", err)
	}
	if _, err := db.BusDomain.Lesson.Reuse(ctx, reader, planID, 1, "", uuid.New()); !errors.Is(err, lessonbus.ErrNotFound) {
		t.Fatalf("private draft reused: %v", err)
	}
	if _, err := db.BusDomain.Lesson.Submit(ctx, author, planID, 1); err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New()
	reused, err := db.BusDomain.Lesson.Reuse(ctx, reader, planID, 1, "Reused vocabulary", requestID)
	if err != nil {
		t.Fatal(err)
	}
	again, err := db.BusDomain.Lesson.Reuse(ctx, reader, planID, 1, "Reused vocabulary", requestID)
	if err != nil || again.ID != reused.ID {
		t.Fatalf("reuse replay: %v", err)
	}
	if reused.AuthorID != reader || reused.CurrentVersion != 1 || reused.Status != lessonbus.Draft || reused.PublishedVersion != nil {
		t.Fatal("reuse inherited workflow or wrong author")
	}
	versions, err := db.BusDomain.Lesson.Versions(ctx, reader, reused.ID)
	if err != nil {
		t.Fatal(err)
	}
	var content map[string]json.RawMessage
	if err := json.Unmarshal(versions[0].Content, &content); err != nil {
		t.Fatal(err)
	}
	var citations []struct {
		Page int `json:"page"`
	}
	if err := json.Unmarshal(content["citations"], &citations); err != nil || len(citations) != 1 || citations[0].Page != 17 {
		t.Fatalf("citations changed: %s, err=%v", content["citations"], err)
	}
	if len(content["reusedFrom"]) == 0 || versions[0].ReviewerID != nil || versions[0].ApproverID != nil {
		t.Fatal("invalid reuse provenance or approvals")
	}
	if _, err := db.BusDomain.Lesson.Reuse(ctx, outsider, planID, 1, "", uuid.New()); !errors.Is(err, lessonbus.ErrForbidden) {
		t.Fatalf("outsider reused: %v", err)
	}
	if err := db.BusDomain.School.Revoke(ctx, root, school.ID, memberships[author].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.BusDomain.Lesson.Create(ctx, author, school.ID, department.ID, draft); !errors.Is(err, lessonbus.ErrForbidden) {
		t.Fatalf("revoked actor replayed: %v", err)
	}
}
