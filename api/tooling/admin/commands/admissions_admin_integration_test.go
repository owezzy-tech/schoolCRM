//go:build integration

package commands

import (
	"context"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

func admissionsProofDB(t *testing.T) (*sqlx.DB, *logger.Logger, uuid.UUID) {
	t.Helper()
	dsn := os.Getenv("ADMISSIONS_ADMIN_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("ADMISSIONS_ADMIN_TEST_DATABASE_URL required for isolated PostgreSQL proof")
	}
	base, err := sqlx.Open("pgx", dsn)
	if err != nil {
		t.Fatal("opening proof database failed")
	}
	schema := "admissions_admin_proof_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := base.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal("creating proof schema failed")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal("invalid proof database URL")
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		t.Fatal("opening isolated proof failed")
	}
	t.Cleanup(func() {
		db.Close()
		if _, err := base.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
			t.Error("owned proof schema cleanup failed")
		}
		base.Close()
	})
	_, err = db.Exec(`CREATE TABLE users (
		user_id UUID PRIMARY KEY, name TEXT NOT NULL, email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL, roles TEXT[] NOT NULL, department TEXT,
		enabled BOOL NOT NULL, date_created TIMESTAMP NOT NULL, date_updated TIMESTAMP NOT NULL
	);
	CREATE TABLE admissions_staff_profiles (
		staff_profile_id UUID PRIMARY KEY, user_id UUID UNIQUE NOT NULL REFERENCES users(user_id),
		admissions_roles TEXT[] NOT NULL, is_active BOOL NOT NULL,
		date_created TIMESTAMP NOT NULL, date_updated TIMESTAMP NOT NULL
	);
	CREATE TABLE audit (
		id UUID PRIMARY KEY, obj_id UUID NOT NULL, obj_domain TEXT NOT NULL,
		obj_name TEXT NOT NULL, actor_id UUID NOT NULL, action TEXT NOT NULL,
		data JSONB, message TEXT, timestamp TIMESTAMP NOT NULL
	)`)
	if err != nil {
		t.Fatal("creating proof tables failed")
	}
	id := uuid.New()
	_, err = db.Exec(`INSERT INTO users VALUES ($1, 'Proof Admin', 'proof@example.invalid', 'unchanged-fixture-hash', '{SUPER_ADMIN}', NULL, true, now(), now())`, id)
	if err != nil {
		t.Fatal("creating proof identity failed")
	}
	return db, logger.New(io.Discard, logger.LevelInfo, "PROOF", func(context.Context) string { return "proof" }), id
}

func TestAdmissionsAdminConcurrentReplay(t *testing.T) {
	db, log, _ := admissionsProofDB(t)
	var wg sync.WaitGroup
	results := make(chan uuid.UUID, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := provisionAdmissionsAdmin(context.Background(), log, db, "proof@example.invalid")
			if err != nil {
				t.Error(err)
				return
			}
			results <- id
		}()
	}
	wg.Wait()
	close(results)
	var first uuid.UUID
	for id := range results {
		if first == uuid.Nil {
			first = id
		}
		if id != first {
			t.Error("replay returned a different profile")
		}
	}
	var valid bool
	if err := db.Get(&valid, `SELECT
		(SELECT count(*) FROM admissions_staff_profiles WHERE is_active AND admissions_roles=ARRAY['ADMISSIONS_ADMIN'])=1
		AND (SELECT count(*) FROM audit WHERE action='bootstrap_admissions_admin' AND obj_domain='ADMISSIONS' AND actor_id='00000000-0000-0000-0000-000000000000')=1
		AND (SELECT count(*) FROM users WHERE roles=ARRAY['SUPER_ADMIN'] AND password_hash='unchanged-fixture-hash')=1`); err != nil || !valid {
		t.Fatal("grant, audit, or unchanged identity invariant failed")
	}
}

func TestAdmissionsAdminRejectsUnsafeTargets(t *testing.T) {
	for _, scenario := range []string{"missing", "ordinary", "disabled", "different_profile", "audit_failure"} {
		t.Run(scenario, func(t *testing.T) {
			db, log, id := admissionsProofDB(t)
			email := "proof@example.invalid"
			var err error
			switch scenario {
			case "missing":
				email = "missing@example.invalid"
			case "ordinary":
				_, err = db.Exec("UPDATE users SET roles=ARRAY['SCHOOL_ADMIN']")
			case "disabled":
				_, err = db.Exec("UPDATE users SET enabled=false")
			case "different_profile":
				_, err = db.Exec(`INSERT INTO admissions_staff_profiles VALUES ($1,$2,'{REPORT_VIEWER}',false,now(),now())`, uuid.New(), id)
			case "audit_failure":
				_, err = db.Exec("DROP TABLE audit")
			}
			if err != nil {
				t.Fatal("preparing proof scenario failed")
			}
			if _, err := provisionAdmissionsAdmin(context.Background(), log, db, email); err == nil {
				t.Fatal("unsafe provisioning accepted")
			}
			var count int
			if err := db.Get(&count, "SELECT count(*) FROM admissions_staff_profiles WHERE is_active"); err != nil || count != 0 {
				t.Fatal("failed grant committed an active profile")
			}
		})
	}
}
