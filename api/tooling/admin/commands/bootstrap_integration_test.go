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

func bootstrapProofDB(t *testing.T) (*sqlx.DB, *logger.Logger) {
	t.Helper()
	dsn := os.Getenv("BOOTSTRAP_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("BOOTSTRAP_TEST_DATABASE_URL is required for isolated PostgreSQL proof")
	}
	base, err := sqlx.Open("pgx", dsn)
	if err != nil {
		t.Fatal("opening proof database failed")
	}
	schema := "bootstrap_proof_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := base.Exec("CREATE SCHEMA " + schema); err != nil {
		base.Close()
		t.Fatal("creating isolated proof schema failed")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal("invalid proof database URL")
	}
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		t.Fatal("opening isolated proof schema failed")
	}
	t.Cleanup(func() {
		db.Close()
		// Only this test's fresh UUID-named schema is removed.
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
	CREATE TABLE audit (
		id UUID PRIMARY KEY, obj_id UUID NOT NULL, obj_domain TEXT NOT NULL,
		obj_name TEXT NOT NULL, actor_id UUID NOT NULL, action TEXT NOT NULL,
		data JSONB, message TEXT, timestamp TIMESTAMP NOT NULL
	)`)
	if err != nil {
		t.Fatal("creating isolated proof tables failed")
	}
	return db, logger.New(io.Discard, logger.LevelInfo, "PROOF", func(context.Context) string { return "proof" })
}

func TestBootstrapPostgresConcurrencyAndReplay(t *testing.T) {
	db, log := bootstrapProofDB(t)
	nu, err := bootstrapUser("Owen Adirah", "proof@example.invalid", fixturePassword)
	if err != nil {
		t.Fatal("fixture validation failed")
	}
	var wg sync.WaitGroup
	results := make(chan uuid.UUID, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := bootstrapAdmin(context.Background(), log, db, nu)
			if err != nil {
				t.Error("concurrent bootstrap failed")
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
		} else if id != first {
			t.Error("concurrent bootstrap returned different accounts")
		}
	}
	var users, audits int
	if err := db.Get(&users, "SELECT count(*) FROM users"); err != nil {
		t.Fatal("reading proof users failed")
	}
	if err := db.Get(&audits, "SELECT count(*) FROM audit"); err != nil {
		t.Fatal("reading proof audits failed")
	}
	if users != 1 || audits != 1 {
		t.Fatal("concurrent retry duplicated a user or audit")
	}
	other, _ := bootstrapUser("Other Operator", "different@example.invalid", fixturePassword)
	if _, err := bootstrapAdmin(context.Background(), log, db, other); err == nil {
		t.Fatal("bootstrap changed a populated database")
	}
}

func TestBootstrapPostgresAuditRollback(t *testing.T) {
	db, log := bootstrapProofDB(t)
	if _, err := db.Exec("DROP TABLE audit"); err != nil {
		t.Fatal("preparing isolated audit failure failed")
	}
	nu, _ := bootstrapUser("Owen Adirah", "proof@example.invalid", fixturePassword)
	if _, err := bootstrapAdmin(context.Background(), log, db, nu); err == nil {
		t.Fatal("bootstrap accepted missing audit storage")
	}
	var users int
	if err := db.Get(&users, "SELECT count(*) FROM users"); err != nil || users != 0 {
		t.Fatal("audit failure committed a privileged account")
	}
}
