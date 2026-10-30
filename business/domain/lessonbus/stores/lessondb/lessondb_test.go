//go:build integration

package lessondb_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus/stores/lessondb"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus/stores/schooldb"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
)

func TestLessonPublishing(t *testing.T) {
	db := dbtest.New(t, "lesson-publishing")
	schools := schoolbus.NewBusiness(schooldb.NewStore(db.Log, db.DB))
	bus := lessonbus.NewBusiness(lessondb.NewStore(db.Log, db.DB))
	ctx := context.Background()
	newUser := func() uuid.UUID {
		t.Helper()
		id := uuid.New()
		if _, err := db.DB.ExecContext(ctx, `INSERT INTO users (user_id,name,email,password_hash,roles,enabled,date_created,date_updated)
		VALUES ($1,'Lesson User',$2,'unused','{TEACHER}',true,now(),now())`, id, id.String()+"@test.invalid"); err != nil {
			t.Fatal(err)
		}
		return id
	}
	root := uuid.New()
	if _, err := db.DB.ExecContext(ctx, `INSERT INTO users (user_id,name,email,password_hash,roles,enabled,date_created,date_updated)
	VALUES ($1,'Root',$2,'unused','{SUPER_ADMIN}',true,now(),now())`, root, root.String()+"@test.invalid"); err != nil {
		t.Fatal(err)
	}
	school, err := schools.CreateSchool(ctx, root, "School One")
	if err != nil {
		t.Fatal(err)
	}
	other, err := schools.CreateSchool(ctx, root, "School Two")
	if err != nil {
		t.Fatal(err)
	}
	newDepartment := func(schoolID uuid.UUID, label string) uuid.UUID {
		t.Helper()
		department, err := schools.CreateDepartment(ctx, root, schoolID, label)
		if err != nil {
			t.Fatal(err)
		}
		return department.ID
	}
	sciences, arts, otherSciences := newDepartment(school.ID, "Sciences"), newDepartment(school.ID, "Arts"), newDepartment(other.ID, "Sciences")
	grant := func(user, schoolID, departmentID uuid.UUID, capabilities ...schoolbus.Capability) schoolbus.Membership {
		t.Helper()
		var membership schoolbus.Membership
		for _, capability := range capabilities {
			var err error
			membership, err = schools.Grant(ctx, root, user, schoolbus.Scope{SchoolID: schoolID, DepartmentID: &departmentID}, capability)
			if err != nil {
				t.Fatal(err)
			}
		}
		return membership
	}
	teacher, hod, dean, allRoles := newUser(), newUser(), newUser(), newUser()
	teaching := grant(teacher, school.ID, sciences, schoolbus.Teach)
	grant(hod, school.ID, sciences, schoolbus.ReviewLessons)
	grant(dean, school.ID, sciences, schoolbus.ApproveLessons)
	grant(allRoles, school.ID, sciences, schoolbus.Teach, schoolbus.ReviewLessons, schoolbus.ApproveLessons)
	artsHOD, otherHOD := newUser(), newUser()
	grant(artsHOD, school.ID, arts, schoolbus.ReviewLessons)
	grant(otherHOD, other.ID, otherSciences, schoolbus.ReviewLessons)

	revision := func(title string) lessonbus.Revision {
		return lessonbus.Revision{Title: title, Content: json.RawMessage(`{"objectives":["` + title + `"]}`), ChangeSummary: "test"}
	}
	expect := func(t *testing.T, err, want error) {
		t.Helper()
		if !errors.Is(err, want) {
			t.Fatalf("got %v, want %v", err, want)
		}
	}
	must := func(t *testing.T) func(lessonbus.Version, error) lessonbus.Version {
		return func(v lessonbus.Version, err error) lessonbus.Version {
			t.Helper()
			if err != nil {
				t.Fatal(err)
			}
			return v
		}
	}
	approve := func(t *testing.T, planID uuid.UUID, number int) {
		t.Helper()
		must(t)(bus.Submit(ctx, teacher, planID, number))
		must(t)(bus.Review(ctx, hod, planID, number, lessonbus.Accept, ""))
		must(t)(bus.Approve(ctx, dean, planID, number, lessonbus.Accept, "good"))
	}
	auditCount := func(t *testing.T, planID uuid.UUID, action string) int {
		t.Helper()
		var count int
		if err := db.DB.GetContext(ctx, &count, `SELECT count(*) FROM audit WHERE obj_id = $1 AND action = $2`, planID, action); err != nil {
			t.Fatal(err)
		}
		return count
	}

	t.Run("only department teachers create plans", func(t *testing.T) {
		_, err := bus.Create(ctx, hod, school.ID, sciences, revision("HOD plan"))
		expect(t, err, lessonbus.ErrForbidden)
		_, err = bus.Create(ctx, teacher, school.ID, arts, revision("Wrong department"))
		expect(t, err, lessonbus.ErrForbidden)
		_, err = bus.Create(ctx, teacher, other.ID, sciences, revision("Department from another school"))
		expect(t, err, lessonbus.ErrNotFound)
		_, err = bus.Create(ctx, teacher, school.ID, sciences, lessonbus.Revision{Title: "Array", Content: json.RawMessage(`[]`)})
		expect(t, err, lessonbus.ErrInvalid)
	})

	t.Run("publication requires HOD review and dean approval of the exact version", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Photosynthesis"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = bus.Publish(ctx, teacher, plan.ID, 1)
		expect(t, err, lessonbus.ErrConflict)
		_, err = bus.Review(ctx, hod, plan.ID, 1, lessonbus.Accept, "")
		expect(t, err, lessonbus.ErrConflict)
		must(t)(bus.Submit(ctx, teacher, plan.ID, 1))
		_, err = bus.Publish(ctx, teacher, plan.ID, 1)
		expect(t, err, lessonbus.ErrConflict)
		_, err = bus.Approve(ctx, dean, plan.ID, 1, lessonbus.Accept, "")
		expect(t, err, lessonbus.ErrConflict)
		must(t)(bus.Review(ctx, hod, plan.ID, 1, lessonbus.Accept, ""))
		_, err = bus.Publish(ctx, teacher, plan.ID, 1)
		expect(t, err, lessonbus.ErrConflict)
		must(t)(bus.Approve(ctx, dean, plan.ID, 1, lessonbus.Accept, ""))
		_, err = bus.Publish(ctx, hod, plan.ID, 1)
		expect(t, err, lessonbus.ErrForbidden)
		published := must(t)(bus.Publish(ctx, teacher, plan.ID, 1))
		if published.Status != lessonbus.Published || *published.ReviewerID != hod || *published.ApproverID != dean {
			t.Fatalf("got %+v", published)
		}
		got, err := bus.Plan(ctx, teacher, plan.ID)
		if err != nil || got.PublishedVersion == nil || *got.PublishedVersion != 1 {
			t.Fatalf("got %+v, %v", got, err)
		}
	})

	t.Run("author, reviewer and approver are different people", func(t *testing.T) {
		plan, err := bus.Create(ctx, allRoles, school.ID, sciences, revision("Self review"))
		if err != nil {
			t.Fatal(err)
		}
		must(t)(bus.Submit(ctx, allRoles, plan.ID, 1))
		_, err = bus.Review(ctx, allRoles, plan.ID, 1, lessonbus.Accept, "")
		expect(t, err, lessonbus.ErrForbidden)

		plan, err = bus.Create(ctx, teacher, school.ID, sciences, revision("Same reviewer and approver"))
		if err != nil {
			t.Fatal(err)
		}
		must(t)(bus.Submit(ctx, teacher, plan.ID, 1))
		must(t)(bus.Review(ctx, allRoles, plan.ID, 1, lessonbus.Accept, ""))
		_, err = bus.Approve(ctx, allRoles, plan.ID, 1, lessonbus.Accept, "")
		expect(t, err, lessonbus.ErrForbidden)
		must(t)(bus.Approve(ctx, dean, plan.ID, 1, lessonbus.Accept, ""))
	})

	t.Run("other departments and schools cannot read or decide", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Scoped"))
		if err != nil {
			t.Fatal(err)
		}
		must(t)(bus.Submit(ctx, teacher, plan.ID, 1))
		for _, user := range []uuid.UUID{artsHOD, otherHOD} {
			_, err = bus.Review(ctx, user, plan.ID, 1, lessonbus.Accept, "")
			expect(t, err, lessonbus.ErrForbidden)
			_, err = bus.Plan(ctx, user, plan.ID)
			expect(t, err, lessonbus.ErrForbidden)
			_, err = bus.Versions(ctx, user, plan.ID)
			expect(t, err, lessonbus.ErrForbidden)
			_, err = bus.Plans(ctx, user, school.ID, sciences)
			expect(t, err, lessonbus.ErrForbidden)
		}
	})

	t.Run("revisions reset approval and leave the publication live", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Cells v1"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		must(t)(bus.Publish(ctx, teacher, plan.ID, 1))
		v2 := must(t)(bus.Revise(ctx, teacher, plan.ID, 1, revision("Cells v2")))
		if v2.Number != 2 || v2.Status != lessonbus.Draft || v2.ReviewerID != nil || v2.ApproverID != nil {
			t.Fatalf("revision inherited workflow state: %+v", v2)
		}
		_, err = bus.Publish(ctx, teacher, plan.ID, 2)
		expect(t, err, lessonbus.ErrConflict)
		_, err = bus.Revise(ctx, teacher, plan.ID, 1, revision("Stale edit"))
		expect(t, err, lessonbus.ErrConflict)
		_, err = bus.Revise(ctx, allRoles, plan.ID, 2, revision("Not the author"))
		expect(t, err, lessonbus.ErrForbidden)
		got, err := bus.Plan(ctx, teacher, plan.ID)
		if err != nil || *got.PublishedVersion != 1 || got.CurrentVersion != 2 || got.Title != "Cells v2" {
			t.Fatalf("got %+v, %v", got, err)
		}

		must(t)(bus.Submit(ctx, teacher, plan.ID, 2))
		must(t)(bus.Revise(ctx, teacher, plan.ID, 2, revision("Cells v3")))
		_, err = bus.Review(ctx, hod, plan.ID, 2, lessonbus.Accept, "")
		expect(t, err, lessonbus.ErrConflict)

		approve(t, plan.ID, 3)
		must(t)(bus.Publish(ctx, teacher, plan.ID, 3))
		versions, err := bus.Versions(ctx, dean, plan.ID)
		if err != nil || len(versions) != 3 {
			t.Fatalf("got %d versions, %v", len(versions), err)
		}
		if v1 := versions[2]; v1.Title != "Cells v1" || v1.Status != lessonbus.Published || string(v1.Content) != `{"objectives": ["Cells v1"]}` {
			t.Fatalf("earlier publication changed: %+v %s", v1, v1.Content)
		}
		got, err = bus.Plan(ctx, teacher, plan.ID)
		if err != nil || *got.PublishedVersion != 3 {
			t.Fatalf("got %+v, %v", got, err)
		}
	})

	t.Run("unsubmitted drafts are visible only to their author", func(t *testing.T) {
		hidden, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Private draft"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = bus.Plan(ctx, hod, hidden.ID)
		expect(t, err, lessonbus.ErrNotFound)
		_, err = bus.Versions(ctx, dean, hidden.ID)
		expect(t, err, lessonbus.ErrNotFound)
		for _, user := range []uuid.UUID{hod, allRoles} {
			plans, err := bus.Plans(ctx, user, school.ID, sciences)
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range plans {
				if p.ID == hidden.ID {
					t.Fatal("another author's draft-only plan is listed")
				}
			}
		}

		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Live v1"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		must(t)(bus.Publish(ctx, teacher, plan.ID, 1))
		must(t)(bus.Revise(ctx, teacher, plan.ID, 1, revision("Secret v2")))
		seen, err := bus.Plan(ctx, hod, plan.ID)
		if err != nil || seen.CurrentVersion != 1 || seen.Title != "Live v1" || seen.Status != lessonbus.Published {
			t.Fatalf("reviewer saw the draft: %+v, %v", seen, err)
		}
		versions, err := bus.Versions(ctx, hod, plan.ID)
		if err != nil || len(versions) != 1 || versions[0].Number != 1 {
			t.Fatalf("reviewer saw draft versions: %+v, %v", versions, err)
		}
		own, err := bus.Versions(ctx, teacher, plan.ID)
		if err != nil || len(own) != 2 {
			t.Fatalf("author lost own draft: %d, %v", len(own), err)
		}
		must(t)(bus.Submit(ctx, teacher, plan.ID, 2))
		if seen, err = bus.Plan(ctx, hod, plan.ID); err != nil || seen.CurrentVersion != 2 || seen.Status != lessonbus.HODReview {
			t.Fatalf("submitted version hidden: %+v, %v", seen, err)
		}
	})

	t.Run("requested changes need a new version", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Needs work"))
		if err != nil {
			t.Fatal(err)
		}
		must(t)(bus.Submit(ctx, teacher, plan.ID, 1))
		_, err = bus.Review(ctx, hod, plan.ID, 1, lessonbus.RequestChanges, " ")
		expect(t, err, lessonbus.ErrInvalid)
		must(t)(bus.Review(ctx, hod, plan.ID, 1, lessonbus.Accept, "fine"))
		v := must(t)(bus.Approve(ctx, dean, plan.ID, 1, lessonbus.RequestChanges, "add assessment"))
		if v.Status != lessonbus.ChangesRequested || v.ReviewFeedback != "fine" || v.ApprovalFeedback != "add assessment" {
			t.Fatalf("got %+v", v)
		}
		_, err = bus.Submit(ctx, teacher, plan.ID, 1)
		expect(t, err, lessonbus.ErrConflict)
		must(t)(bus.Revise(ctx, teacher, plan.ID, 1, revision("Reworked")))
		approve(t, plan.ID, 2)
	})

	t.Run("revoked teaching authority blocks publication", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Revoked"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		if err := schools.Revoke(ctx, root, school.ID, teaching.ID); err != nil {
			t.Fatal(err)
		}
		_, err = bus.Publish(ctx, teacher, plan.ID, 1)
		expect(t, err, lessonbus.ErrForbidden)
	})
	teaching = grant(teacher, school.ID, sciences, schoolbus.Teach)

	t.Run("concurrent publication retries publish once", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Concurrent"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		var wg sync.WaitGroup
		failures := make(chan error, 8)
		for range 8 {
			wg.Go(func() {
				if _, err := bus.Publish(ctx, teacher, plan.ID, 1); err != nil {
					failures <- err
				}
			})
		}
		wg.Wait()
		close(failures)
		for err := range failures {
			t.Fatal(err)
		}
		if n := auditCount(t, plan.ID, "lesson_published"); n != 1 {
			t.Fatalf("got %d publication events", n)
		}
	})

	t.Run("audit failure rolls back publication", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Audit"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		if _, err := db.DB.ExecContext(ctx, `ALTER TABLE audit ADD CONSTRAINT reject_lesson_audit CHECK (action <> 'lesson_published') NOT VALID`); err != nil {
			t.Fatal(err)
		}
		_, err = bus.Publish(ctx, teacher, plan.ID, 1)
		if _, dropErr := db.DB.ExecContext(ctx, `ALTER TABLE audit DROP CONSTRAINT reject_lesson_audit`); dropErr != nil {
			t.Fatal(dropErr)
		}
		if err == nil {
			t.Fatal("expected audit failure")
		}
		got, err := bus.Plan(ctx, teacher, plan.ID)
		if err != nil || got.PublishedVersion != nil || got.Status != lessonbus.Approved {
			t.Fatalf("publication committed without audit: %+v, %v", got, err)
		}
	})

	t.Run("database rejects content edits and ungated publication", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Immutable"))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.DB.ExecContext(ctx, `UPDATE lesson_plan_versions SET content = '{"edited":true}' WHERE plan_id = $1`, plan.ID); err == nil {
			t.Fatal("content update succeeded")
		}
		if _, err := db.DB.ExecContext(ctx, `UPDATE lesson_plan_versions SET status = 'published', published_at = now() WHERE plan_id = $1`, plan.ID); err == nil {
			t.Fatal("ungated publication succeeded")
		}
		if _, err := db.DB.ExecContext(ctx, `UPDATE lesson_plan_versions SET reviewer_id = author_id WHERE plan_id = $1`, plan.ID); err == nil {
			t.Fatal("author reviewed own version")
		}
	})

	t.Run("waiting publication sees committed revocation", func(t *testing.T) {
		plan, err := bus.Create(ctx, teacher, school.ID, sciences, revision("Race"))
		if err != nil {
			t.Fatal(err)
		}
		approve(t, plan.ID, 1)
		tx, err := db.DB.BeginTxx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = tx.Rollback() }()
		var pid int
		if err := tx.GetContext(ctx, &pid, `SELECT pg_backend_pid()`); err != nil {
			t.Fatal(err)
		}
		// Mirror schoolbus.Revoke: lock the school, then deactivate the membership.
		var locked uuid.UUID
		if err := tx.GetContext(ctx, &locked, `SELECT school_id FROM schools WHERE school_id = $1 FOR UPDATE`, school.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := tx.ExecContext(ctx, `UPDATE school_memberships SET active = false WHERE membership_id = $1`, teaching.ID); err != nil {
			t.Fatal(err)
		}
		result := make(chan error, 1)
		commandCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		go func() {
			_, err := bus.Publish(commandCtx, teacher, plan.ID, 1)
			result <- err
		}()
		deadline := time.Now().Add(5 * time.Second)
		for blocked := false; !blocked; {
			if time.Now().After(deadline) {
				t.Fatal("expected publication to wait on the school lock")
			}
			if err := db.DB.GetContext(ctx, &blocked, `SELECT EXISTS (SELECT 1 FROM pg_stat_activity WHERE $1 = ANY(pg_blocking_pids(pid)))`, pid); err != nil {
				t.Fatal(err)
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		expect(t, <-result, lessonbus.ErrForbidden)
	})
}
