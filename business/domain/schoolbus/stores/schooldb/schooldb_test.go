//go:build integration

package schooldb_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus/stores/schooldb"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
)

func TestSchoolAccess(t *testing.T) {
	db := dbtest.New(t, "school-access")
	bus := schoolbus.NewBusiness(schooldb.NewStore(db.Log, db.DB))
	ctx := context.Background()
	newUser := func(role string) uuid.UUID {
		t.Helper()
		id := uuid.New()
		_, err := db.DB.ExecContext(ctx, `INSERT INTO users (user_id,name,email,password_hash,roles,enabled,date_created,date_updated)
		VALUES ($1,'Test User',$2,'unused',$3::text[],true,now(),now())`, id, id.String()+"@test.invalid", "{"+role+"}")
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	root, admin, teacher := newUser("SUPER_ADMIN"), newUser("SCHOOL_ADMIN"), newUser("TEACHER")
	school, err := bus.CreateSchool(ctx, root, "School One")
	if err != nil {
		t.Fatal(err)
	}
	other, err := bus.CreateSchool(ctx, root, "School Two")
	if err != nil {
		t.Fatal(err)
	}
	department, err := bus.CreateDepartment(ctx, root, school.ID, "Sciences")
	if err != nil {
		t.Fatal(err)
	}
	otherDepartment, err := bus.CreateDepartment(ctx, root, other.ID, "Sciences")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("global admin role does not grant school authority", func(t *testing.T) {
		_, err := bus.CreateDepartment(ctx, admin, school.ID, "Arts")
		if !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	manager, err := bus.Grant(ctx, root, admin, schoolbus.Scope{SchoolID: school.ID}, schoolbus.ManageMembers)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("delegated admin manages only their school", func(t *testing.T) {
		if _, err := bus.CreateDepartment(ctx, admin, school.ID, "Arts"); err != nil {
			t.Fatal(err)
		}
		if _, err := bus.CreateDepartment(ctx, admin, other.ID, "Arts"); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if _, err := bus.Memberships(ctx, admin, other.ID); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if _, err := bus.Grant(ctx, admin, teacher, schoolbus.Scope{SchoolID: other.ID, DepartmentID: &otherDepartment.ID}, schoolbus.Teach); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("admin cannot delegate or revoke management", func(t *testing.T) {
		if _, err := bus.Grant(ctx, admin, admin, schoolbus.Scope{SchoolID: school.ID}, schoolbus.ManageMembers); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if err := bus.Revoke(ctx, admin, school.ID, manager.ID); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("scope and capability validation", func(t *testing.T) {
		cases := []struct {
			name  string
			scope schoolbus.Scope
			cap   schoolbus.Capability
			want  error
		}{
			{"unknown capability", schoolbus.Scope{SchoolID: school.ID}, "root", schoolbus.ErrInvalid},
			{"teacher missing department", schoolbus.Scope{SchoolID: school.ID}, schoolbus.Teach, schoolbus.ErrInvalid},
			{"manager with department", schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, schoolbus.ManageMembers, schoolbus.ErrInvalid},
			{"department from another school", schoolbus.Scope{SchoolID: school.ID, DepartmentID: &otherDepartment.ID}, schoolbus.Teach, schoolbus.ErrNotFound},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if _, err := bus.Grant(ctx, root, teacher, tc.scope, tc.cap); !errors.Is(err, tc.want) {
					t.Fatalf("got %v, want %v", err, tc.want)
				}
			})
		}
		if _, err := bus.Grant(ctx, root, teacher, schoolbus.Scope{SchoolID: school.ID}, schoolbus.ManageMembers); !errors.Is(err, schoolbus.ErrInvalid) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("all lesson responsibilities are department scoped", func(t *testing.T) {
		for _, cap := range []schoolbus.Capability{schoolbus.Teach, schoolbus.ReviewLessons, schoolbus.ApproveLessons} {
			member, err := bus.Grant(ctx, admin, teacher, schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, cap)
			if err != nil {
				t.Fatal(err)
			}
			if member.SchoolID != school.ID || *member.DepartmentID != department.ID || member.Capability != cap {
				t.Fatalf("wrong scope: %+v", member)
			}
		}
	})
	t.Run("reads expose only authorised schools", func(t *testing.T) {
		schools, err := bus.Schools(ctx, teacher)
		if err != nil {
			t.Fatal(err)
		}
		if len(schools) != 1 || schools[0].ID != school.ID {
			t.Fatalf("got %+v", schools)
		}
		if _, err := bus.Departments(ctx, teacher, school.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := bus.Departments(ctx, teacher, other.ID); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if _, err := bus.Memberships(ctx, teacher, school.ID); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("concurrent duplicate grants reuse one membership", func(t *testing.T) {
		target := newUser("TEACHER")
		var wg sync.WaitGroup
		ids := make(chan uuid.UUID, 8)
		failures := make(chan error, 8)
		for range 8 {
			wg.Go(func() {
				member, err := bus.Grant(ctx, admin, target, schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, schoolbus.Teach)
				if err != nil {
					failures <- err
					return
				}
				ids <- member.ID
			})
		}
		wg.Wait()
		close(ids)
		close(failures)
		for err := range failures {
			t.Fatal(err)
		}
		var first uuid.UUID
		for id := range ids {
			if first == uuid.Nil {
				first = id
			}
			if id != first {
				t.Fatal("duplicate membership created")
			}
		}
	})
	t.Run("audit failure rolls back membership", func(t *testing.T) {
		target := newUser("TEACHER")
		if _, err := db.DB.ExecContext(ctx, `ALTER TABLE audit ADD CONSTRAINT reject_school_audit CHECK (action <> 'membership_granted') NOT VALID`); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if _, err := db.DB.ExecContext(ctx, `ALTER TABLE audit DROP CONSTRAINT reject_school_audit`); err != nil {
				t.Error(err)
			}
		}()
		if _, err := bus.Grant(ctx, admin, target, schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, schoolbus.Teach); err == nil {
			t.Fatal("expected audit failure")
		}
		var count int
		if err := db.DB.GetContext(ctx, &count, `SELECT count(*) FROM school_memberships WHERE user_id = $1`, target); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatal("membership committed without audit")
		}
	})
	t.Run("removed global role is effective immediately", func(t *testing.T) {
		if _, err := db.DB.ExecContext(ctx, `UPDATE users SET roles = '{TEACHER}' WHERE user_id = $1`, admin); err != nil {
			t.Fatal(err)
		}
		if _, err := bus.CreateDepartment(ctx, admin, school.ID, "Denied"); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if _, err := db.DB.ExecContext(ctx, `UPDATE users SET roles = '{SCHOOL_ADMIN}' WHERE user_id = $1`, admin); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("disabled users cannot act or receive grants", func(t *testing.T) {
		target := newUser("TEACHER")
		if _, err := db.DB.ExecContext(ctx, `UPDATE users SET enabled = false WHERE user_id = $1`, target); err != nil {
			t.Fatal(err)
		}
		if _, err := bus.Grant(ctx, admin, target, schoolbus.Scope{SchoolID: school.ID, DepartmentID: &department.ID}, schoolbus.Teach); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		if _, err := bus.Schools(ctx, target); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("revocation is idempotent and effective immediately", func(t *testing.T) {
		if err := bus.Revoke(ctx, root, school.ID, manager.ID); err != nil {
			t.Fatal(err)
		}
		if err := bus.Revoke(ctx, root, school.ID, manager.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := bus.CreateDepartment(ctx, admin, school.ID, "Denied"); !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
		var events int
		if err := db.DB.GetContext(ctx, &events, `SELECT count(*) FROM audit WHERE obj_id = $1 AND action = 'membership_revoked'`, manager.ID); err != nil {
			t.Fatal(err)
		}
		if events != 1 {
			t.Fatalf("got %d revocation events", events)
		}
		var data []byte
		if err := db.DB.GetContext(ctx, &data, `SELECT data FROM audit WHERE obj_id = $1 AND action = 'membership_revoked'`, manager.ID); err != nil {
			t.Fatal(err)
		}
		var snapshot schoolbus.Membership
		if err := json.Unmarshal(data, &snapshot); err != nil {
			t.Fatal(err)
		}
		if snapshot.ID != manager.ID || snapshot.SchoolID != school.ID || snapshot.UserID != admin || snapshot.Capability != schoolbus.ManageMembers || snapshot.Active {
			t.Fatalf("incomplete revocation audit: %+v", snapshot)
		}
	})

	waitForBlock := func(t *testing.T, pid int) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			var blocked bool
			if err := db.DB.GetContext(ctx, &blocked, `SELECT EXISTS (SELECT 1 FROM pg_stat_activity WHERE $1 = ANY(pg_blocking_pids(pid)))`, pid); err != nil {
				t.Fatal(err)
			}
			if blocked {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatal("expected operation to wait on database lock")
	}
	t.Run("waiting management command sees committed revocation", func(t *testing.T) {
		manager, err := bus.Grant(ctx, root, admin, schoolbus.Scope{SchoolID: school.ID}, schoolbus.ManageMembers)
		if err != nil {
			t.Fatal(err)
		}
		tx, err := db.DB.BeginTxx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = tx.Rollback() }()
		var pid int
		if err := tx.GetContext(ctx, &pid, `SELECT pg_backend_pid()`); err != nil {
			t.Fatal(err)
		}
		var locked uuid.UUID
		if err := tx.GetContext(ctx, &locked, `SELECT school_id FROM schools WHERE school_id = $1 FOR UPDATE`, school.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := tx.ExecContext(ctx, `UPDATE school_memberships SET active = false WHERE membership_id = $1`, manager.ID); err != nil {
			t.Fatal(err)
		}
		result := make(chan error, 1)
		commandCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		go func() {
			_, err := bus.CreateDepartment(commandCtx, admin, school.ID, "Concurrent denied")
			result <- err
		}()
		waitForBlock(t, pid)
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		if err := <-result; !errors.Is(err, schoolbus.ErrForbidden) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("fresh user read holds role updates until transaction ends", func(t *testing.T) {
		store := schooldb.NewStore(db.Log, db.DB)
		updateResult := make(chan error, 1)
		err := store.WithinTx(ctx, func(s schoolbus.Storer) error {
			actor, err := s.Actor(ctx, admin)
			if err != nil {
				return err
			}
			if !actor.SchoolAdmin {
				t.Fatal("expected school admin")
			}
			// The transaction connection is the only idle-in-transaction backend for this database.
			var pid int
			if err := db.DB.GetContext(ctx, &pid, `SELECT pid FROM pg_stat_activity WHERE datname = current_database() AND state = 'idle in transaction'`); err != nil {
				return err
			}
			updateCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			t.Cleanup(cancel)
			go func() {
				_, err := db.DB.ExecContext(updateCtx, `UPDATE users SET roles = '{TEACHER}' WHERE user_id = $1`, admin)
				updateResult <- err
			}()
			waitForBlock(t, pid)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := <-updateResult; err != nil {
			t.Fatal(err)
		}
	})
}
