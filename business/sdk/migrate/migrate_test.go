package migrate

import (
	"fmt"
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

func TestKenyaReferenceCatalogCompletionMigrationCreatesReadOnlyTables(t *testing.T) {
	t.Parallel()

	migration := migrationBlock(t, "-- Version: 1.27")

	checks := []string{
		"CREATE TABLE wards",
		"CREATE TABLE universities",
		"CREATE TABLE programme_clusters",
		"CREATE TABLE knqf_levels",
		"CREATE TABLE programmes",
		"FOREIGN KEY (sub_county_code) REFERENCES sub_counties(code)",
		"FOREIGN KEY (university_code) REFERENCES universities(code)",
		"FOREIGN KEY (cluster_code) REFERENCES programme_clusters(code)",
		"FOREIGN KEY (knqf_level_code) REFERENCES knqf_levels(code)",
		"CONSTRAINT knqf_levels_level_range CHECK (level BETWEEN 1 AND 10)",
	}

	for _, check := range checks {
		if !strings.Contains(migration, check) {
			t.Fatalf("reference catalog migration missing %q", check)
		}
	}
}

func TestKenyaReferenceCatalogSeedCoversCanonicalLookups(t *testing.T) {
	t.Parallel()

	checks := []string{
		"INSERT INTO wards",
		"INSERT INTO universities",
		"INSERT INTO programme_clusters",
		"INSERT INTO knqf_levels",
		"INSERT INTO programmes",
		"('KNQF-1', 1, 'KNQF Level 1'",
		"('KNQF-10', 10, 'KNQF Level 10'",
		"('CL02', 'Engineering and Technology'",
		"('JKUAT-BENG-CIVIL', 'JKUAT', 'CL02', 'KNQF-7'",
	}

	for _, check := range checks {
		if !strings.Contains(seedDoc, check) {
			t.Fatalf("reference catalog seed missing %q", check)
		}
	}
}

func TestKenyaReferenceCatalogSeedCoversCanonicalCounts(t *testing.T) {
	t.Parallel()

	checks := map[string]int{
		"counties":           47,
		"sub_counties":       321,
		"wards":              10,
		"universities":       5,
		"programme_clusters": 4,
		"knqf_levels":        10,
		"programmes":         5,
	}

	for table, want := range checks {
		if got := seedInsertRowCount(seedDoc, table); got != want {
			t.Fatalf("expected %d %s seed rows, got %d", want, table, got)
		}
	}
}

func TestAdmissionsAdminSourceDataMigrationCreatesCampaignAndCommunicationTables(t *testing.T) {
	t.Parallel()

	migration := migrationBlock(t, "-- Version: 1.29")

	checks := []string{
		"CREATE TABLE admissions_campaigns",
		"CREATE TABLE admissions_campaign_audit_events",
		"CREATE TABLE admissions_communications",
		"FOREIGN KEY (campaign_id) REFERENCES admissions_campaigns(campaign_id)",
		"FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id)",
		"CREATE INDEX idx_admissions_campaigns_status",
		"CREATE INDEX idx_admissions_communications_status",
		"DROP CONSTRAINT sync_events_event_type",
		"'SMS_OUTBOUND'",
		"'WHATSAPP_WEBHOOK_INBOUND'",
	}

	for _, check := range checks {
		if !strings.Contains(migration, check) {
			t.Fatalf("admissions admin source data migration missing %q", check)
		}
	}
}

func TestAdmissionsAdminSeedCoversCampaignCommunicationAndProviderStatuses(t *testing.T) {
	t.Parallel()

	checks := []string{
		"INSERT INTO admissions_campaigns",
		"INSERT INTO admissions_campaign_audit_events",
		"INSERT INTO admissions_communications",
		"INSERT INTO admissions_sync_jobs",
		"INSERT INTO admissions_sync_events",
		"INSERT INTO audit",
		"Fall 2026 Open House Invite",
		"Missing Documents Reminder",
		"Financial Aid Deadline",
		"Spring 2026 Yield Campaign",
		"MSG-8492",
		"MSG-8486-WA",
		"RATE_LIMITED",
		"INVALID_SIGNATURE",
	}

	for _, check := range checks {
		if !strings.Contains(seedDoc, check) {
			t.Fatalf("admissions admin seed missing %q", check)
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

func seedInsertRowCount(seed string, table string) int {
	marker := fmt.Sprintf("INSERT INTO %s", table)
	start := strings.Index(seed, marker)
	if start == -1 {
		return 0
	}

	remaining := seed[start+len(marker):]
	end := strings.Index(remaining, "ON CONFLICT DO NOTHING;")
	if end == -1 {
		return 0
	}

	block := remaining[:end]
	return strings.Count(block, "\t('")
}
