package migrate

import (
	"strings"
	"testing"
)

func TestKenyaIdentityMigrationPreservesExternalSISID(t *testing.T) {
	t.Parallel()

	migration := migrationBlock(t, "-- Version: 1.21")

	if strings.Contains(migration, "DROP COLUMN external_sis_id") {
		t.Fatal("identity migration must preserve external_sis_id during deprecation window")
	}

	if strings.Contains(migration, "UPDATE admissions_constituents") {
		t.Fatal("identity migration must not infer Kenyan identifiers from external_sis_id")
	}

	for _, column := range []string{
		"ADD COLUMN national_id TEXT NULL",
		"ADD COLUMN upi TEXT NULL",
		"ADD COLUMN kcse_index_number TEXT NULL",
	} {
		if !strings.Contains(migration, column) {
			t.Fatalf("migration missing additive column %q", column)
		}
	}
}

func TestKenyaIdentityMigrationFlagsManualBackfills(t *testing.T) {
	t.Parallel()

	migration := migrationBlock(t, "-- Version: 1.21")

	checks := []string{
		"CREATE TABLE admissions_identity_backfill_reviews",
		"INSERT INTO admissions_identity_backfill_reviews",
		"WHERE\n    external_sis_id IS NOT NULL",
		"external_sis_id requires manual classification",
		"ON CONFLICT DO NOTHING",
	}

	for _, check := range checks {
		if !strings.Contains(migration, check) {
			t.Fatalf("manual backfill migration missing %q", check)
		}
	}
}

func TestKenyaIdentityMigrationAddsIndependentUniqueIndexes(t *testing.T) {
	t.Parallel()

	migration := migrationBlock(t, "-- Version: 1.21")

	checks := []string{
		"CREATE UNIQUE INDEX idx_admissions_constituents_national_id ON admissions_constituents (national_id) WHERE national_id IS NOT NULL",
		"CREATE UNIQUE INDEX idx_admissions_constituents_upi ON admissions_constituents (upi) WHERE upi IS NOT NULL",
		"CREATE UNIQUE INDEX idx_admissions_constituents_kcse_index_number ON admissions_constituents (kcse_index_number) WHERE kcse_index_number IS NOT NULL",
	}

	for _, check := range checks {
		if !strings.Contains(migration, check) {
			t.Fatalf("migration missing uniqueness guarantee %q", check)
		}
	}
}

func migrationBlock(t *testing.T, marker string) string {
	t.Helper()

	start := strings.Index(migrateDoc, marker)
	if start == -1 {
		t.Fatalf("migration marker %q not found", marker)
	}

	remaining := migrateDoc[start+len(marker):]
	next := strings.Index(remaining, "-- Version:")
	if next == -1 {
		return migrateDoc[start:]
	}

	return migrateDoc[start : start+len(marker)+next]
}
