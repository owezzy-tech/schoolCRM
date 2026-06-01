// Package admissionsdb contains admissions related persistence functionality.
package admissionsdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Store manages the set of APIs for admissions database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

const constituentColumns = `
		constituent_id,
		first_name,
		last_name,
		preferred_name,
		middle_name,
		suffix,
		date_of_birth,
		primary_email,
		primary_phone,
		external_sis_id,
		national_id,
		national_id_verified_at,
		national_id_verified_by_adapter,
		upi,
		upi_verified_at,
		upi_verified_by_adapter,
		kcse_index_number,
		kcse_index_verified_at,
		kcse_index_verified_by_adapter,
		lifecycle_stage,
		duplicate_status,
		duplicate_of_id,
		sis_synced_at,
		date_created,
		date_updated`

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (admissionsbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Health returns the current admissions store scaffold metadata.
func (s *Store) Health(ctx context.Context) (admissionsbus.Health, error) {
	s.log.Info(ctx, "admissions scaffold health", "status", "ready")

	return admissionsbus.Health{
		Context:    admissionsbus.DomainName,
		Status:     "ready",
		Aggregates: admissionsbus.AggregateNames(),
	}, nil
}

// CreateStaffProfile inserts a new admissions staff profile into the database.
func (s *Store) CreateStaffProfile(ctx context.Context, profile admissionsbus.StaffProfile) error {
	const q = `
	INSERT INTO admissions_staff_profiles
		(staff_profile_id, user_id, admissions_roles, is_active, date_created, date_updated)
	VALUES
		(:staff_profile_id, :user_id, :admissions_roles, :is_active, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBStaffProfile(profile)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateStaffProfile replaces mutable admissions staff profile data in the database.
func (s *Store) UpdateStaffProfile(ctx context.Context, profile admissionsbus.StaffProfile) error {
	const q = `
	UPDATE
		admissions_staff_profiles
	SET
		user_id = :user_id,
		admissions_roles = :admissions_roles,
		is_active = :is_active,
		date_updated = :date_updated
	WHERE
		staff_profile_id = :staff_profile_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBStaffProfile(profile)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryStaffProfiles retrieves a list of admissions staff profiles from the database.
func (s *Store) QueryStaffProfiles(ctx context.Context, filter admissionsbus.StaffProfileQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.StaffProfile, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		staff_profile_id, user_id, admissions_roles, is_active, date_created, date_updated
	FROM
		admissions_staff_profiles`

	buf := bytes.NewBufferString(q)
	s.applyStaffProfileFilter(filter, data, buf)

	orderByClause, err := staffProfileOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbProfiles []staffProfileDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbProfiles); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusStaffProfiles(dbProfiles)
}

// CountStaffProfiles returns the total number of admissions staff profiles in the database.
func (s *Store) CountStaffProfiles(ctx context.Context, filter admissionsbus.StaffProfileQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_staff_profiles`

	buf := bytes.NewBufferString(q)
	s.applyStaffProfileFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryStaffProfileByID finds an admissions staff profile by ID.
func (s *Store) QueryStaffProfileByID(ctx context.Context, profileID uuid.UUID) (admissionsbus.StaffProfile, error) {
	filter := admissionsbus.StaffProfileQueryFilter{ID: &profileID}
	profile, err := s.queryStaffProfile(ctx, filter)
	if err != nil {
		return admissionsbus.StaffProfile{}, err
	}

	return profile, nil
}

// QueryStaffProfileByUserID finds an admissions staff profile by identity user ID.
func (s *Store) QueryStaffProfileByUserID(ctx context.Context, userID uuid.UUID) (admissionsbus.StaffProfile, error) {
	filter := admissionsbus.StaffProfileQueryFilter{UserID: &userID}
	profile, err := s.queryStaffProfile(ctx, filter)
	if err != nil {
		return admissionsbus.StaffProfile{}, err
	}

	return profile, nil
}

func (s *Store) queryStaffProfile(ctx context.Context, filter admissionsbus.StaffProfileQueryFilter) (admissionsbus.StaffProfile, error) {
	data := map[string]any{}

	const q = `
	SELECT
		staff_profile_id, user_id, admissions_roles, is_active, date_created, date_updated
	FROM
		admissions_staff_profiles`

	buf := bytes.NewBufferString(q)
	s.applyStaffProfileFilter(filter, data, buf)

	var dbProfile staffProfileDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbProfile); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.StaffProfile{}, fmt.Errorf("db: %w", admissionsbus.ErrStaffProfileNotFound)
		}
		return admissionsbus.StaffProfile{}, fmt.Errorf("db: %w", err)
	}

	return toBusStaffProfile(dbProfile)
}

// CreateApplicantProfile inserts a new admissions applicant profile into the database.
func (s *Store) CreateApplicantProfile(ctx context.Context, profile admissionsbus.ApplicantProfile) error {
	const q = `
	INSERT INTO admissions_applicant_profiles
		(applicant_profile_id, user_id, constituent_id, is_active, date_created, date_updated)
	VALUES
		(:applicant_profile_id, :user_id, :constituent_id, :is_active, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBApplicantProfile(profile)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateApplicantProfile replaces mutable admissions applicant profile data in the database.
func (s *Store) UpdateApplicantProfile(ctx context.Context, profile admissionsbus.ApplicantProfile) error {
	const q = `
	UPDATE
		admissions_applicant_profiles
	SET
		user_id = :user_id,
		constituent_id = :constituent_id,
		is_active = :is_active,
		date_updated = :date_updated
	WHERE
		applicant_profile_id = :applicant_profile_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBApplicantProfile(profile)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryApplicantProfiles retrieves a list of admissions applicant profiles from the database.
func (s *Store) QueryApplicantProfiles(ctx context.Context, filter admissionsbus.ApplicantProfileQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicantProfile, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		applicant_profile_id, user_id, constituent_id, is_active, date_created, date_updated
	FROM
		admissions_applicant_profiles`

	buf := bytes.NewBufferString(q)
	s.applyApplicantProfileFilter(filter, data, buf)

	orderByClause, err := applicantProfileOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbProfiles []applicantProfileDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbProfiles); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusApplicantProfiles(dbProfiles), nil
}

// CountApplicantProfiles returns the total number of admissions applicant profiles in the database.
func (s *Store) CountApplicantProfiles(ctx context.Context, filter admissionsbus.ApplicantProfileQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_applicant_profiles`

	buf := bytes.NewBufferString(q)
	s.applyApplicantProfileFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryApplicantProfileByID finds an admissions applicant profile by ID.
func (s *Store) QueryApplicantProfileByID(ctx context.Context, profileID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	filter := admissionsbus.ApplicantProfileQueryFilter{ID: &profileID}
	profile, err := s.queryApplicantProfile(ctx, filter)
	if err != nil {
		return admissionsbus.ApplicantProfile{}, err
	}

	return profile, nil
}

// QueryApplicantProfileByUserID finds an admissions applicant profile by identity user ID.
func (s *Store) QueryApplicantProfileByUserID(ctx context.Context, userID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	filter := admissionsbus.ApplicantProfileQueryFilter{UserID: &userID}
	profile, err := s.queryApplicantProfile(ctx, filter)
	if err != nil {
		return admissionsbus.ApplicantProfile{}, err
	}

	return profile, nil
}

// QueryApplicantProfileByConstituentID finds an admissions applicant profile by constituent ID.
func (s *Store) QueryApplicantProfileByConstituentID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	filter := admissionsbus.ApplicantProfileQueryFilter{ConstituentID: &constituentID}
	profile, err := s.queryApplicantProfile(ctx, filter)
	if err != nil {
		return admissionsbus.ApplicantProfile{}, err
	}

	return profile, nil
}

func (s *Store) queryApplicantProfile(ctx context.Context, filter admissionsbus.ApplicantProfileQueryFilter) (admissionsbus.ApplicantProfile, error) {
	data := map[string]any{}

	const q = `
	SELECT
		applicant_profile_id, user_id, constituent_id, is_active, date_created, date_updated
	FROM
		admissions_applicant_profiles`

	buf := bytes.NewBufferString(q)
	s.applyApplicantProfileFilter(filter, data, buf)

	var dbProfile applicantProfileDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbProfile); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.ApplicantProfile{}, fmt.Errorf("db: %w", admissionsbus.ErrApplicantProfileNotFound)
		}
		return admissionsbus.ApplicantProfile{}, fmt.Errorf("db: %w", err)
	}

	return toBusApplicantProfile(dbProfile), nil
}

// CreateLeadScoreRule inserts a new admissions lead score rule into the database.
func (s *Store) CreateLeadScoreRule(ctx context.Context, rule admissionsbus.LeadScoreRule) error {
	const q = `
	INSERT INTO admissions_lead_score_rules
		(lead_score_rule_id, name, description, criteria, points, is_active, priority, date_created, date_updated)
	VALUES
		(:lead_score_rule_id, :name, :description, :criteria, :points, :is_active, :priority, :date_created, :date_updated)`

	dbRule, err := toDBLeadScoreRule(rule)
	if err != nil {
		return fmt.Errorf("todb: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbRule); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateLeadScoreRule replaces mutable admissions lead score rule data in the database.
func (s *Store) UpdateLeadScoreRule(ctx context.Context, rule admissionsbus.LeadScoreRule) error {
	const q = `
	UPDATE
		admissions_lead_score_rules
	SET
		name = :name,
		description = :description,
		criteria = :criteria,
		points = :points,
		is_active = :is_active,
		priority = :priority,
		date_updated = :date_updated
	WHERE
		lead_score_rule_id = :lead_score_rule_id`

	dbRule, err := toDBLeadScoreRule(rule)
	if err != nil {
		return fmt.Errorf("todb: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbRule); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryLeadScoreRules retrieves a list of admissions lead score rules from the database.
func (s *Store) QueryLeadScoreRules(ctx context.Context, filter admissionsbus.LeadScoreRuleQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.LeadScoreRule, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		lead_score_rule_id, name, description, criteria, points, is_active, priority, date_created, date_updated
	FROM
		admissions_lead_score_rules`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreRuleFilter(filter, data, buf)

	orderByClause, err := leadScoreRuleOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbRules []leadScoreRuleDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbRules); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusLeadScoreRules(dbRules)
}

// CountLeadScoreRules returns the total number of admissions lead score rules in the database.
func (s *Store) CountLeadScoreRules(ctx context.Context, filter admissionsbus.LeadScoreRuleQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_lead_score_rules`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreRuleFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryLeadScoreRuleByID finds an admissions lead score rule by ID.
func (s *Store) QueryLeadScoreRuleByID(ctx context.Context, ruleID uuid.UUID) (admissionsbus.LeadScoreRule, error) {
	filter := admissionsbus.LeadScoreRuleQueryFilter{ID: &ruleID}
	rule, err := s.queryLeadScoreRule(ctx, filter)
	if err != nil {
		return admissionsbus.LeadScoreRule{}, err
	}

	return rule, nil
}

func (s *Store) queryLeadScoreRule(ctx context.Context, filter admissionsbus.LeadScoreRuleQueryFilter) (admissionsbus.LeadScoreRule, error) {
	data := map[string]any{}

	const q = `
	SELECT
		lead_score_rule_id, name, description, criteria, points, is_active, priority, date_created, date_updated
	FROM
		admissions_lead_score_rules`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreRuleFilter(filter, data, buf)

	var dbRule leadScoreRuleDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbRule); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.LeadScoreRule{}, fmt.Errorf("db: %w", admissionsbus.ErrLeadScoreRuleNotFound)
		}
		return admissionsbus.LeadScoreRule{}, fmt.Errorf("db: %w", err)
	}

	return toBusLeadScoreRule(dbRule)
}

// UpsertLeadScore inserts or replaces a constituent's latest admissions lead score.
func (s *Store) UpsertLeadScore(ctx context.Context, score admissionsbus.LeadScore) error {
	const q = `
	INSERT INTO admissions_lead_scores
		(lead_score_id, constituent_id, total_score, score_band, breakdown, recalculated_at, date_created, date_updated)
	VALUES
		(:lead_score_id, :constituent_id, :total_score, :score_band, :breakdown, :recalculated_at, :date_created, :date_updated)
	ON CONFLICT (constituent_id) DO UPDATE SET
		total_score = EXCLUDED.total_score,
		score_band = EXCLUDED.score_band,
		breakdown = EXCLUDED.breakdown,
		recalculated_at = EXCLUDED.recalculated_at,
		date_updated = EXCLUDED.date_updated`

	dbScore, err := toDBLeadScore(score)
	if err != nil {
		return fmt.Errorf("todb: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbScore); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryLeadScores retrieves a list of admissions lead scores from the database.
func (s *Store) QueryLeadScores(ctx context.Context, filter admissionsbus.LeadScoreQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.LeadScore, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		lead_score_id, constituent_id, total_score, score_band, breakdown, recalculated_at, date_created, date_updated
	FROM
		admissions_lead_scores`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreFilter(filter, data, buf)

	orderByClause, err := leadScoreOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbScores []leadScoreDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbScores); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusLeadScores(dbScores)
}

// CountLeadScores returns the total number of admissions lead scores in the database.
func (s *Store) CountLeadScores(ctx context.Context, filter admissionsbus.LeadScoreQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_lead_scores`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryLeadScoreByID finds an admissions lead score by ID.
func (s *Store) QueryLeadScoreByID(ctx context.Context, scoreID uuid.UUID) (admissionsbus.LeadScore, error) {
	filter := admissionsbus.LeadScoreQueryFilter{ID: &scoreID}
	score, err := s.queryLeadScore(ctx, filter)
	if err != nil {
		return admissionsbus.LeadScore{}, err
	}

	return score, nil
}

// QueryLeadScoreByConstituentID finds an admissions lead score by constituent ID.
func (s *Store) QueryLeadScoreByConstituentID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.LeadScore, error) {
	filter := admissionsbus.LeadScoreQueryFilter{ConstituentID: &constituentID}
	score, err := s.queryLeadScore(ctx, filter)
	if err != nil {
		return admissionsbus.LeadScore{}, err
	}

	return score, nil
}

func (s *Store) queryLeadScore(ctx context.Context, filter admissionsbus.LeadScoreQueryFilter) (admissionsbus.LeadScore, error) {
	data := map[string]any{}

	const q = `
	SELECT
		lead_score_id, constituent_id, total_score, score_band, breakdown, recalculated_at, date_created, date_updated
	FROM
		admissions_lead_scores`

	buf := bytes.NewBufferString(q)
	s.applyLeadScoreFilter(filter, data, buf)

	var dbScore leadScoreDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbScore); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.LeadScore{}, fmt.Errorf("db: %w", admissionsbus.ErrLeadScoreNotFound)
		}
		return admissionsbus.LeadScore{}, fmt.Errorf("db: %w", err)
	}

	return toBusLeadScore(dbScore)
}

// CreateConstituent inserts a new Constituent into the database.
func (s *Store) CreateConstituent(ctx context.Context, cst admissionsbus.Constituent) error {
	const q = `
	INSERT INTO admissions_constituents
		(constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, national_id, national_id_verified_at, national_id_verified_by_adapter, upi, upi_verified_at, upi_verified_by_adapter, kcse_index_number, kcse_index_verified_at, kcse_index_verified_by_adapter, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated)
	VALUES
		(:constituent_id, :first_name, :last_name, :preferred_name, :middle_name, :suffix, :date_of_birth, :primary_email, :primary_phone, :external_sis_id, :national_id, :national_id_verified_at, :national_id_verified_by_adapter, :upi, :upi_verified_at, :upi_verified_by_adapter, :kcse_index_number, :kcse_index_verified_at, :kcse_index_verified_by_adapter, :lifecycle_stage, :duplicate_status, :duplicate_of_id, :sis_synced_at, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBConstituent(cst)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateConstituent replaces mutable constituent data in the database.
func (s *Store) UpdateConstituent(ctx context.Context, cst admissionsbus.Constituent) error {
	const q = `
	UPDATE
		admissions_constituents
	SET
		preferred_name = :preferred_name,
		middle_name = :middle_name,
		suffix = :suffix,
		primary_email = :primary_email,
		primary_phone = :primary_phone,
		national_id = :national_id,
		national_id_verified_at = :national_id_verified_at,
		national_id_verified_by_adapter = :national_id_verified_by_adapter,
		upi = :upi,
		upi_verified_at = :upi_verified_at,
		upi_verified_by_adapter = :upi_verified_by_adapter,
		kcse_index_number = :kcse_index_number,
		kcse_index_verified_at = :kcse_index_verified_at,
		kcse_index_verified_by_adapter = :kcse_index_verified_by_adapter,
		lifecycle_stage = :lifecycle_stage,
		duplicate_status = :duplicate_status,
		duplicate_of_id = :duplicate_of_id,
		sis_synced_at = :sis_synced_at,
		date_updated = :date_updated
	WHERE
		constituent_id = :constituent_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBConstituent(cst)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryConstituents retrieves a list of Constituents from the database.
func (s *Store) QueryConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Constituent, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		` + constituentColumns + `
	FROM
		admissions_constituents`

	buf := bytes.NewBufferString(q)
	s.applyConstituentFilter(filter, data, buf)

	orderByClause, err := constituentOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbConstituents []constituentDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbConstituents); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusConstituents(dbConstituents)
}

// CountConstituents returns the total number of Constituents in the database.
func (s *Store) CountConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_constituents`

	buf := bytes.NewBufferString(q)
	s.applyConstituentFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryConstituentByID finds a Constituent by ID.
func (s *Store) QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.Constituent, error) {
	filter := admissionsbus.ConstituentQueryFilter{ID: &constituentID}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

// QueryConstituentByPrimaryEmail finds a Constituent by primary email.
func (s *Store) QueryConstituentByPrimaryEmail(ctx context.Context, email string) (admissionsbus.Constituent, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return admissionsbus.Constituent{}, fmt.Errorf("parse email: %w", err)
	}

	filter := admissionsbus.ConstituentQueryFilter{PrimaryEmail: address}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

// QueryConstituentByExternalSISID finds a Constituent by SIS ID.
func (s *Store) QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Constituent, error) {
	filter := admissionsbus.ConstituentQueryFilter{ExternalSISID: &externalSISID}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

// QueryConstituentByNationalID finds a Constituent by Kenyan national ID.
func (s *Store) QueryConstituentByNationalID(ctx context.Context, nationalID string) (admissionsbus.Constituent, error) {
	filter := admissionsbus.ConstituentQueryFilter{NationalID: &nationalID}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

// QueryConstituentByUPI finds a Constituent by Kenyan UPI.
func (s *Store) QueryConstituentByUPI(ctx context.Context, upi string) (admissionsbus.Constituent, error) {
	filter := admissionsbus.ConstituentQueryFilter{UPI: &upi}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

// QueryConstituentByKCSEIndexNumber finds a Constituent by KCSE index number.
func (s *Store) QueryConstituentByKCSEIndexNumber(ctx context.Context, kcseIndexNumber string) (admissionsbus.Constituent, error) {
	filter := admissionsbus.ConstituentQueryFilter{KCSEIndexNumber: &kcseIndexNumber}
	cst, err := s.queryConstituent(ctx, filter)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return cst, nil
}

func (s *Store) queryConstituent(ctx context.Context, filter admissionsbus.ConstituentQueryFilter) (admissionsbus.Constituent, error) {
	data := map[string]any{}

	const q = `
	SELECT
		` + constituentColumns + `
	FROM
		admissions_constituents`

	buf := bytes.NewBufferString(q)
	s.applyConstituentFilter(filter, data, buf)

	var dbConstituent constituentDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbConstituent); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.Constituent{}, fmt.Errorf("db: %w", admissionsbus.ErrConstituentNotFound)
		}
		return admissionsbus.Constituent{}, fmt.Errorf("db: %w", err)
	}

	return toBusConstituent(dbConstituent)
}

// CreateInquiry inserts a new Inquiry into the database.
func (s *Store) CreateInquiry(ctx context.Context, inquiry admissionsbus.Inquiry) error {
	const q = `
	INSERT INTO admissions_inquiries
		(inquiry_id, constituent_id, first_name, last_name, date_of_birth, primary_email, primary_phone, program_of_interest, term_of_interest, source, utm_source, utm_medium, utm_campaign, message, status, date_created, date_updated)
	VALUES
		(:inquiry_id, :constituent_id, :first_name, :last_name, :date_of_birth, :primary_email, :primary_phone, :program_of_interest, :term_of_interest, :source, :utm_source, :utm_medium, :utm_campaign, :message, :status, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBInquiry(inquiry)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateInquiry replaces mutable inquiry data in the database.
func (s *Store) UpdateInquiry(ctx context.Context, inquiry admissionsbus.Inquiry) error {
	const q = `
	UPDATE
		admissions_inquiries
	SET
		constituent_id = :constituent_id,
		first_name = :first_name,
		last_name = :last_name,
		date_of_birth = :date_of_birth,
		primary_email = :primary_email,
		primary_phone = :primary_phone,
		program_of_interest = :program_of_interest,
		term_of_interest = :term_of_interest,
		source = :source,
		utm_source = :utm_source,
		utm_medium = :utm_medium,
		utm_campaign = :utm_campaign,
		message = :message,
		status = :status,
		date_updated = :date_updated
	WHERE
		inquiry_id = :inquiry_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBInquiry(inquiry)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryInquiries retrieves a list of Inquiries from the database.
func (s *Store) QueryInquiries(ctx context.Context, filter admissionsbus.InquiryQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Inquiry, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		inquiry_id, constituent_id, first_name, last_name, date_of_birth, primary_email, primary_phone, program_of_interest, term_of_interest, source, utm_source, utm_medium, utm_campaign, message, status, date_created, date_updated
	FROM
		admissions_inquiries`

	buf := bytes.NewBufferString(q)
	s.applyInquiryFilter(filter, data, buf)

	orderByClause, err := inquiryOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbInquiries []inquiryDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbInquiries); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusInquiries(dbInquiries)
}

// CountInquiries returns the total number of Inquiries in the database.
func (s *Store) CountInquiries(ctx context.Context, filter admissionsbus.InquiryQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_inquiries`

	buf := bytes.NewBufferString(q)
	s.applyInquiryFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryInquiryByID finds an Inquiry by ID.
func (s *Store) QueryInquiryByID(ctx context.Context, inquiryID uuid.UUID) (admissionsbus.Inquiry, error) {
	filter := admissionsbus.InquiryQueryFilter{ID: &inquiryID}
	inquiry, err := s.queryInquiry(ctx, filter)
	if err != nil {
		return admissionsbus.Inquiry{}, err
	}

	return inquiry, nil
}

func (s *Store) queryInquiry(ctx context.Context, filter admissionsbus.InquiryQueryFilter) (admissionsbus.Inquiry, error) {
	data := map[string]any{}

	const q = `
	SELECT
		inquiry_id, constituent_id, first_name, last_name, date_of_birth, primary_email, primary_phone, program_of_interest, term_of_interest, source, utm_source, utm_medium, utm_campaign, message, status, date_created, date_updated
	FROM
		admissions_inquiries`

	buf := bytes.NewBufferString(q)
	s.applyInquiryFilter(filter, data, buf)

	var dbInquiry inquiryDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbInquiry); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.Inquiry{}, fmt.Errorf("db: %w", admissionsbus.ErrInquiryNotFound)
		}
		return admissionsbus.Inquiry{}, fmt.Errorf("db: %w", err)
	}

	return toBusInquiry(dbInquiry)
}

// UpsertProgram creates or updates a Program by immutable external SIS ID.
func (s *Store) UpsertProgram(ctx context.Context, prg admissionsbus.Program) error {
	const q = `
	INSERT INTO admissions_programs
		(program_id, external_sis_id, name, code, description, degree_level, is_active, synced_at, date_created, date_updated)
	VALUES
		(:program_id, :external_sis_id, :name, :code, :description, :degree_level, :is_active, :synced_at, :date_created, :date_updated)
	ON CONFLICT (external_sis_id) DO UPDATE SET
		name = EXCLUDED.name,
		code = EXCLUDED.code,
		description = EXCLUDED.description,
		degree_level = EXCLUDED.degree_level,
		is_active = EXCLUDED.is_active,
		synced_at = EXCLUDED.synced_at,
		date_updated = EXCLUDED.date_updated`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProgram(prg)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryPrograms retrieves a list of Program reference records from the database.
func (s *Store) QueryPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Program, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		program_id, external_sis_id, name, code, description, degree_level, is_active, synced_at, date_created, date_updated
	FROM
		admissions_programs`

	buf := bytes.NewBufferString(q)
	s.applyProgramFilter(filter, data, buf)

	orderByClause, err := programOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbPrograms []programDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPrograms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusPrograms(dbPrograms), nil
}

// CountPrograms returns the total number of Program reference records in the database.
func (s *Store) CountPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_programs`

	buf := bytes.NewBufferString(q)
	s.applyProgramFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryProgramByID finds a Program by ID.
func (s *Store) QueryProgramByID(ctx context.Context, programID uuid.UUID) (admissionsbus.Program, error) {
	filter := admissionsbus.ProgramQueryFilter{ID: &programID}
	programs, err := s.queryProgram(ctx, filter)
	if err != nil {
		return admissionsbus.Program{}, err
	}

	return programs, nil
}

// QueryProgramByExternalSISID finds a Program by immutable SIS ID.
func (s *Store) QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Program, error) {
	filter := admissionsbus.ProgramQueryFilter{ExternalSISID: &externalSISID}
	program, err := s.queryProgram(ctx, filter)
	if err != nil {
		return admissionsbus.Program{}, err
	}

	return program, nil
}

func (s *Store) queryProgram(ctx context.Context, filter admissionsbus.ProgramQueryFilter) (admissionsbus.Program, error) {
	data := map[string]any{}

	const q = `
	SELECT
		program_id, external_sis_id, name, code, description, degree_level, is_active, synced_at, date_created, date_updated
	FROM
		admissions_programs`

	buf := bytes.NewBufferString(q)
	s.applyProgramFilter(filter, data, buf)

	var dbProgram programDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbProgram); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.Program{}, fmt.Errorf("db: %w", admissionsbus.ErrProgramNotFound)
		}
		return admissionsbus.Program{}, fmt.Errorf("db: %w", err)
	}

	return toBusProgram(dbProgram), nil
}

// UpsertAcademicTerm creates or updates an AcademicTerm by immutable external SIS ID.
func (s *Store) UpsertAcademicTerm(ctx context.Context, term admissionsbus.AcademicTerm) error {
	const q = `
	INSERT INTO admissions_academic_terms
		(academic_term_id, external_sis_id, name, code, term_type, start_date, end_date, application_start_date, application_deadline, is_active, synced_at, date_created, date_updated)
	VALUES
		(:academic_term_id, :external_sis_id, :name, :code, :term_type, :start_date, :end_date, :application_start_date, :application_deadline, :is_active, :synced_at, :date_created, :date_updated)
	ON CONFLICT (external_sis_id) DO UPDATE SET
		name = EXCLUDED.name,
		code = EXCLUDED.code,
		term_type = EXCLUDED.term_type,
		start_date = EXCLUDED.start_date,
		end_date = EXCLUDED.end_date,
		application_start_date = EXCLUDED.application_start_date,
		application_deadline = EXCLUDED.application_deadline,
		is_active = EXCLUDED.is_active,
		synced_at = EXCLUDED.synced_at,
		date_updated = EXCLUDED.date_updated`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBAcademicTerm(term)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryAcademicTerms retrieves a list of AcademicTerm reference records from the database.
func (s *Store) QueryAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.AcademicTerm, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		academic_term_id, external_sis_id, name, code, term_type, start_date, end_date, application_start_date, application_deadline, is_active, synced_at, date_created, date_updated
	FROM
		admissions_academic_terms`

	buf := bytes.NewBufferString(q)
	s.applyAcademicTermFilter(filter, data, buf)

	orderByClause, err := academicTermOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbTerms []academicTermDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbTerms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusAcademicTerms(dbTerms), nil
}

// CountAcademicTerms returns the total number of AcademicTerm reference records in the database.
func (s *Store) CountAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_academic_terms`

	buf := bytes.NewBufferString(q)
	s.applyAcademicTermFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryAcademicTermByID finds an AcademicTerm by ID.
func (s *Store) QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (admissionsbus.AcademicTerm, error) {
	filter := admissionsbus.AcademicTermQueryFilter{ID: &termID}
	term, err := s.queryAcademicTerm(ctx, filter)
	if err != nil {
		return admissionsbus.AcademicTerm{}, err
	}

	return term, nil
}

// QueryAcademicTermByExternalSISID finds an AcademicTerm by immutable SIS ID.
func (s *Store) QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.AcademicTerm, error) {
	filter := admissionsbus.AcademicTermQueryFilter{ExternalSISID: &externalSISID}
	term, err := s.queryAcademicTerm(ctx, filter)
	if err != nil {
		return admissionsbus.AcademicTerm{}, err
	}

	return term, nil
}

// CreateDuplicateReview inserts a new duplicate review into the queue.
func (s *Store) CreateDuplicateReview(ctx context.Context, review admissionsbus.DuplicateReview) error {
	const q = `
	INSERT INTO admissions_duplicate_reviews
		(duplicate_review_id, source_constituent_id, candidate_constituent_id, match_type, match_score, match_reason, status, resolved_by, resolved_at, resolution_note, date_created, date_updated)
	VALUES
		(:duplicate_review_id, :source_constituent_id, :candidate_constituent_id, :match_type, :match_score, :match_reason, :status, :resolved_by, :resolved_at, :resolution_note, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBDuplicateReview(review)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateDuplicateReview replaces mutable duplicate review data in the database.
func (s *Store) UpdateDuplicateReview(ctx context.Context, review admissionsbus.DuplicateReview) error {
	const q = `
	UPDATE
		admissions_duplicate_reviews
	SET
		status = :status,
		resolved_by = :resolved_by,
		resolved_at = :resolved_at,
		resolution_note = :resolution_note,
		date_updated = :date_updated
	WHERE
		duplicate_review_id = :duplicate_review_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBDuplicateReview(review)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryDuplicateReviews retrieves a list of duplicate reviews from the database.
func (s *Store) QueryDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.DuplicateReview, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		duplicate_review_id, source_constituent_id, candidate_constituent_id, match_type, match_score, match_reason, status, resolved_by, resolved_at, resolution_note, date_created, date_updated
	FROM
		admissions_duplicate_reviews`

	buf := bytes.NewBufferString(q)
	s.applyDuplicateReviewFilter(filter, data, buf)

	orderByClause, err := duplicateReviewOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbReviews []duplicateReviewDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbReviews); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusDuplicateReviews(dbReviews), nil
}

// CountDuplicateReviews returns the total number of duplicate reviews in the database.
func (s *Store) CountDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_duplicate_reviews`

	buf := bytes.NewBufferString(q)
	s.applyDuplicateReviewFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryDuplicateReviewByID finds a duplicate review by ID.
func (s *Store) QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (admissionsbus.DuplicateReview, error) {
	filter := admissionsbus.DuplicateReviewQueryFilter{ID: &reviewID}
	review, err := s.queryDuplicateReview(ctx, filter)
	if err != nil {
		return admissionsbus.DuplicateReview{}, err
	}

	return review, nil
}

func (s *Store) queryDuplicateReview(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter) (admissionsbus.DuplicateReview, error) {
	data := map[string]any{}

	const q = `
	SELECT
		duplicate_review_id, source_constituent_id, candidate_constituent_id, match_type, match_score, match_reason, status, resolved_by, resolved_at, resolution_note, date_created, date_updated
	FROM
		admissions_duplicate_reviews`

	buf := bytes.NewBufferString(q)
	s.applyDuplicateReviewFilter(filter, data, buf)

	var dbReview duplicateReviewDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbReview); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.DuplicateReview{}, fmt.Errorf("db: %w", admissionsbus.ErrDuplicateReviewNotFound)
		}
		return admissionsbus.DuplicateReview{}, fmt.Errorf("db: %w", err)
	}

	return toBusDuplicateReview(dbReview), nil
}

// CreateApplication inserts a new Application into the database.
func (s *Store) CreateApplication(ctx context.Context, app admissionsbus.Application) error {
	const q = `
	INSERT INTO admissions_applications
		(application_id, constituent_id, program_id, academic_term_id, application_type, status, kuccps_placement, kcse_result, assigned_reviewer_id, submitted_at, date_created, date_updated)
	VALUES
		(:application_id, :constituent_id, :program_id, :academic_term_id, :application_type, :status, :kuccps_placement, :kcse_result, :assigned_reviewer_id, :submitted_at, :date_created, :date_updated)`

	dbApp, err := toDBApplication(app)
	if err != nil {
		return fmt.Errorf("to db application: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbApp); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", admissionsbus.ErrDuplicateApplication)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateApplication replaces mutable application data in the database.
func (s *Store) UpdateApplication(ctx context.Context, app admissionsbus.Application) error {
	const q = `
	UPDATE
		admissions_applications
	SET
		status = :status,
		kuccps_placement = :kuccps_placement,
		kcse_result = :kcse_result,
		assigned_reviewer_id = :assigned_reviewer_id,
		submitted_at = :submitted_at,
		date_updated = :date_updated
	WHERE
		application_id = :application_id`

	dbApp, err := toDBApplication(app)
	if err != nil {
		return fmt.Errorf("to db application: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbApp); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", admissionsbus.ErrDuplicateApplication)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryApplications retrieves a list of Applications from the database.
func (s *Store) QueryApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Application, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		application_id, constituent_id, program_id, academic_term_id, application_type, status, kuccps_placement, kcse_result, assigned_reviewer_id, submitted_at, date_created, date_updated
	FROM
		admissions_applications`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFilter(filter, data, buf)

	orderByClause, err := applicationOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbApplications []applicationDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbApplications); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusApplications(dbApplications)
}

// CountApplications returns the total number of Applications in the database.
func (s *Store) CountApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_applications`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryApplicationByID finds an Application by ID.
func (s *Store) QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (admissionsbus.Application, error) {
	filter := admissionsbus.ApplicationQueryFilter{ID: &applicationID}
	app, err := s.queryApplication(ctx, filter)
	if err != nil {
		return admissionsbus.Application{}, err
	}

	return app, nil
}

// QueryActiveApplicationByTuple finds an active Application for a constituent, term, and program.
func (s *Store) QueryActiveApplicationByTuple(ctx context.Context, constituentID uuid.UUID, academicTermID uuid.UUID, programID uuid.UUID) (admissionsbus.Application, error) {
	activeOnly := true
	filter := admissionsbus.ApplicationQueryFilter{
		ConstituentID:  &constituentID,
		AcademicTermID: &academicTermID,
		ProgramID:      &programID,
		ActiveOnly:     &activeOnly,
	}
	app, err := s.queryApplication(ctx, filter)
	if err != nil {
		return admissionsbus.Application{}, err
	}

	return app, nil
}

func (s *Store) queryApplication(ctx context.Context, filter admissionsbus.ApplicationQueryFilter) (admissionsbus.Application, error) {
	data := map[string]any{}

	const q = `
	SELECT
		application_id, constituent_id, program_id, academic_term_id, application_type, status, kuccps_placement, kcse_result, assigned_reviewer_id, submitted_at, date_created, date_updated
	FROM
		admissions_applications`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFilter(filter, data, buf)

	var dbApplication applicationDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbApplication); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.Application{}, fmt.Errorf("db: %w", admissionsbus.ErrApplicationNotFound)
		}
		return admissionsbus.Application{}, fmt.Errorf("db: %w", err)
	}

	return toBusApplication(dbApplication)
}

// CreateApplicationFormTemplate inserts a new application form template.
func (s *Store) CreateApplicationFormTemplate(ctx context.Context, template admissionsbus.ApplicationFormTemplate) error {
	const q = `
	INSERT INTO admissions_application_form_templates
		(form_template_id, program_id, academic_term_id, application_type, name, description, version, required_fields, checklist_items, is_active, priority, date_created, date_updated)
	VALUES
		(:form_template_id, :program_id, :academic_term_id, :application_type, :name, :description, :version, :required_fields, :checklist_items, :is_active, :priority, :date_created, :date_updated)`

	dbTemplate, err := toDBApplicationFormTemplate(template)
	if err != nil {
		return fmt.Errorf("to db application form template: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbTemplate); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateApplicationFormTemplate replaces mutable application form template data.
func (s *Store) UpdateApplicationFormTemplate(ctx context.Context, template admissionsbus.ApplicationFormTemplate) error {
	const q = `
	UPDATE
		admissions_application_form_templates
	SET
		program_id = :program_id,
		academic_term_id = :academic_term_id,
		application_type = :application_type,
		name = :name,
		description = :description,
		version = :version,
		required_fields = :required_fields,
		checklist_items = :checklist_items,
		is_active = :is_active,
		priority = :priority,
		date_updated = :date_updated
	WHERE
		form_template_id = :form_template_id`

	dbTemplate, err := toDBApplicationFormTemplate(template)
	if err != nil {
		return fmt.Errorf("to db application form template: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbTemplate); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryApplicationFormTemplates retrieves application form templates.
func (s *Store) QueryApplicationFormTemplates(ctx context.Context, filter admissionsbus.ApplicationFormTemplateQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicationFormTemplate, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		form_template_id, program_id, academic_term_id, application_type, name, description, version, required_fields, checklist_items, is_active, priority, date_created, date_updated
	FROM
		admissions_application_form_templates`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFormTemplateFilter(filter, data, buf)

	orderByClause, err := applicationFormTemplateOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbTemplates []applicationFormTemplateDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbTemplates); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusApplicationFormTemplates(dbTemplates)
}

// CountApplicationFormTemplates returns the total number of application form templates.
func (s *Store) CountApplicationFormTemplates(ctx context.Context, filter admissionsbus.ApplicationFormTemplateQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_application_form_templates`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFormTemplateFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryApplicationFormTemplateByID finds an application form template by ID.
func (s *Store) QueryApplicationFormTemplateByID(ctx context.Context, templateID uuid.UUID) (admissionsbus.ApplicationFormTemplate, error) {
	filter := admissionsbus.ApplicationFormTemplateQueryFilter{ID: &templateID}
	template, err := s.queryApplicationFormTemplate(ctx, filter)
	if err != nil {
		return admissionsbus.ApplicationFormTemplate{}, err
	}

	return template, nil
}

func (s *Store) queryApplicationFormTemplate(ctx context.Context, filter admissionsbus.ApplicationFormTemplateQueryFilter) (admissionsbus.ApplicationFormTemplate, error) {
	data := map[string]any{}

	const q = `
	SELECT
		form_template_id, program_id, academic_term_id, application_type, name, description, version, required_fields, checklist_items, is_active, priority, date_created, date_updated
	FROM
		admissions_application_form_templates`

	buf := bytes.NewBufferString(q)
	s.applyApplicationFormTemplateFilter(filter, data, buf)

	var dbTemplate applicationFormTemplateDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbTemplate); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.ApplicationFormTemplate{}, fmt.Errorf("db: %w", admissionsbus.ErrFormTemplateNotFound)
		}
		return admissionsbus.ApplicationFormTemplate{}, fmt.Errorf("db: %w", err)
	}

	return toBusApplicationFormTemplate(dbTemplate)
}

// CreateCustomFieldDefinition inserts a custom field definition.
func (s *Store) CreateCustomFieldDefinition(ctx context.Context, definition admissionsbus.CustomFieldDefinition) error {
	const q = `
	INSERT INTO admissions_custom_field_definitions
		(custom_field_definition_id, owner, field_key, label, description, data_type, is_required, options, validation, is_searchable, is_reportable, is_importable, is_exportable, display_order, is_active, date_created, date_updated)
	VALUES
		(:custom_field_definition_id, :owner, :field_key, :label, :description, :data_type, :is_required, :options, :validation, :is_searchable, :is_reportable, :is_importable, :is_exportable, :display_order, :is_active, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCustomFieldDefinition(definition)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateCustomFieldDefinition replaces mutable custom field definition data.
func (s *Store) UpdateCustomFieldDefinition(ctx context.Context, definition admissionsbus.CustomFieldDefinition) error {
	const q = `
	UPDATE
		admissions_custom_field_definitions
	SET
		owner = :owner,
		field_key = :field_key,
		label = :label,
		description = :description,
		data_type = :data_type,
		is_required = :is_required,
		options = :options,
		validation = :validation,
		is_searchable = :is_searchable,
		is_reportable = :is_reportable,
		is_importable = :is_importable,
		is_exportable = :is_exportable,
		display_order = :display_order,
		is_active = :is_active,
		date_updated = :date_updated
	WHERE
		custom_field_definition_id = :custom_field_definition_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCustomFieldDefinition(definition)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryCustomFieldDefinitions retrieves custom field definitions.
func (s *Store) QueryCustomFieldDefinitions(ctx context.Context, filter admissionsbus.CustomFieldDefinitionQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.CustomFieldDefinition, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		custom_field_definition_id, owner, field_key, label, description, data_type, is_required, options, validation, is_searchable, is_reportable, is_importable, is_exportable, display_order, is_active, date_created, date_updated
	FROM
		admissions_custom_field_definitions`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldDefinitionFilter(filter, data, buf)

	orderByClause, err := customFieldDefinitionOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbDefinitions []customFieldDefinitionDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbDefinitions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCustomFieldDefinitions(dbDefinitions), nil
}

// CountCustomFieldDefinitions returns the total number of custom field definitions.
func (s *Store) CountCustomFieldDefinitions(ctx context.Context, filter admissionsbus.CustomFieldDefinitionQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_custom_field_definitions`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldDefinitionFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryCustomFieldDefinitionByID finds a custom field definition by ID.
func (s *Store) QueryCustomFieldDefinitionByID(ctx context.Context, definitionID uuid.UUID) (admissionsbus.CustomFieldDefinition, error) {
	filter := admissionsbus.CustomFieldDefinitionQueryFilter{ID: &definitionID}
	definition, err := s.queryCustomFieldDefinition(ctx, filter)
	if err != nil {
		return admissionsbus.CustomFieldDefinition{}, err
	}

	return definition, nil
}

func (s *Store) queryCustomFieldDefinition(ctx context.Context, filter admissionsbus.CustomFieldDefinitionQueryFilter) (admissionsbus.CustomFieldDefinition, error) {
	data := map[string]any{}

	const q = `
	SELECT
		custom_field_definition_id, owner, field_key, label, description, data_type, is_required, options, validation, is_searchable, is_reportable, is_importable, is_exportable, display_order, is_active, date_created, date_updated
	FROM
		admissions_custom_field_definitions`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldDefinitionFilter(filter, data, buf)

	var dbDefinition customFieldDefinitionDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbDefinition); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.CustomFieldDefinition{}, fmt.Errorf("db: %w", admissionsbus.ErrCustomFieldDefinitionNotFound)
		}
		return admissionsbus.CustomFieldDefinition{}, fmt.Errorf("db: %w", err)
	}

	return toBusCustomFieldDefinition(dbDefinition), nil
}

// SetCustomFieldValue creates or replaces a custom field value for an owner record.
func (s *Store) SetCustomFieldValue(ctx context.Context, value admissionsbus.CustomFieldValue) error {
	const q = `
	INSERT INTO admissions_custom_field_values
		(custom_field_value_id, custom_field_definition_id, owner, owner_id, value, date_created, date_updated)
	VALUES
		(:custom_field_value_id, :custom_field_definition_id, :owner, :owner_id, :value, :date_created, :date_updated)
	ON CONFLICT (custom_field_definition_id, owner, owner_id) DO UPDATE SET
		value = EXCLUDED.value,
		date_updated = EXCLUDED.date_updated`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCustomFieldValue(value)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryCustomFieldValues retrieves custom field values.
func (s *Store) QueryCustomFieldValues(ctx context.Context, filter admissionsbus.CustomFieldValueQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.CustomFieldValue, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		custom_field_value_id, custom_field_definition_id, owner, owner_id, value, date_created, date_updated
	FROM
		admissions_custom_field_values`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldValueFilter(filter, data, buf)

	orderByClause, err := customFieldValueOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbValues []customFieldValueDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbValues); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCustomFieldValues(dbValues), nil
}

// CountCustomFieldValues returns the total number of custom field values.
func (s *Store) CountCustomFieldValues(ctx context.Context, filter admissionsbus.CustomFieldValueQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_custom_field_values`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldValueFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryCustomFieldValueByID finds a custom field value by ID.
func (s *Store) QueryCustomFieldValueByID(ctx context.Context, valueID uuid.UUID) (admissionsbus.CustomFieldValue, error) {
	filter := admissionsbus.CustomFieldValueQueryFilter{ID: &valueID}
	value, err := s.queryCustomFieldValue(ctx, filter)
	if err != nil {
		return admissionsbus.CustomFieldValue{}, err
	}

	return value, nil
}

func (s *Store) queryCustomFieldValue(ctx context.Context, filter admissionsbus.CustomFieldValueQueryFilter) (admissionsbus.CustomFieldValue, error) {
	data := map[string]any{}

	const q = `
	SELECT
		custom_field_value_id, custom_field_definition_id, owner, owner_id, value, date_created, date_updated
	FROM
		admissions_custom_field_values`

	buf := bytes.NewBufferString(q)
	s.applyCustomFieldValueFilter(filter, data, buf)

	var dbValue customFieldValueDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbValue); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.CustomFieldValue{}, fmt.Errorf("db: %w", admissionsbus.ErrCustomFieldValueNotFound)
		}
		return admissionsbus.CustomFieldValue{}, fmt.Errorf("db: %w", err)
	}

	return toBusCustomFieldValue(dbValue), nil
}

// CreateApplicationTransition inserts immutable Application transition history.
func (s *Store) CreateApplicationTransition(ctx context.Context, transition admissionsbus.ApplicationTransition) error {
	const q = `
	INSERT INTO admissions_application_transitions
		(application_transition_id, application_id, from_status, to_status, actor_id, reason, note, metadata, date_created)
	VALUES
		(:application_transition_id, :application_id, :from_status, :to_status, :actor_id, :reason, :note, :metadata, :date_created)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBApplicationTransition(transition)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryApplicationTransitions retrieves a list of Application transitions from the database.
func (s *Store) QueryApplicationTransitions(ctx context.Context, filter admissionsbus.ApplicationTransitionQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicationTransition, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		application_transition_id, application_id, from_status, to_status, actor_id, reason, note, metadata, date_created
	FROM
		admissions_application_transitions`

	buf := bytes.NewBufferString(q)
	s.applyApplicationTransitionFilter(filter, data, buf)

	orderByClause, err := applicationTransitionOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbTransitions []applicationTransitionDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbTransitions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusApplicationTransitions(dbTransitions), nil
}

// CountApplicationTransitions returns the total number of Application transitions in the database.
func (s *Store) CountApplicationTransitions(ctx context.Context, filter admissionsbus.ApplicationTransitionQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_application_transitions`

	buf := bytes.NewBufferString(q)
	s.applyApplicationTransitionFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// CreateChecklistItem inserts a new application checklist item.
func (s *Store) CreateChecklistItem(ctx context.Context, item admissionsbus.ChecklistItem) error {
	const q = `
	INSERT INTO admissions_checklist_items
		(checklist_item_id, application_id, item_key, document_name, description, is_required, status, display_order, date_created, date_updated)
	VALUES
		(:checklist_item_id, :application_id, :item_key, :document_name, :description, :is_required, :status, :display_order, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBChecklistItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateChecklistItem replaces mutable application checklist item data.
func (s *Store) UpdateChecklistItem(ctx context.Context, item admissionsbus.ChecklistItem) error {
	const q = `
	UPDATE
		admissions_checklist_items
	SET
		application_id = :application_id,
		item_key = :item_key,
		document_name = :document_name,
		description = :description,
		is_required = :is_required,
		status = :status,
		display_order = :display_order,
		date_updated = :date_updated
	WHERE
		checklist_item_id = :checklist_item_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBChecklistItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryChecklistItems retrieves application checklist items.
func (s *Store) QueryChecklistItems(ctx context.Context, filter admissionsbus.ChecklistItemQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ChecklistItem, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		checklist_item_id, application_id, item_key, document_name, description, is_required, status, display_order, date_created, date_updated
	FROM
		admissions_checklist_items`

	buf := bytes.NewBufferString(q)
	s.applyChecklistItemFilter(filter, data, buf)

	orderByClause, err := checklistItemOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbItems []checklistItemDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbItems); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusChecklistItems(dbItems), nil
}

// CountChecklistItems returns the total number of application checklist items.
func (s *Store) CountChecklistItems(ctx context.Context, filter admissionsbus.ChecklistItemQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_checklist_items`

	buf := bytes.NewBufferString(q)
	s.applyChecklistItemFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryChecklistItemByID finds an application checklist item by ID.
func (s *Store) QueryChecklistItemByID(ctx context.Context, itemID uuid.UUID) (admissionsbus.ChecklistItem, error) {
	filter := admissionsbus.ChecklistItemQueryFilter{ID: &itemID}
	item, err := s.queryChecklistItem(ctx, filter)
	if err != nil {
		return admissionsbus.ChecklistItem{}, err
	}

	return item, nil
}

func (s *Store) queryChecklistItem(ctx context.Context, filter admissionsbus.ChecklistItemQueryFilter) (admissionsbus.ChecklistItem, error) {
	data := map[string]any{}

	const q = `
	SELECT
		checklist_item_id, application_id, item_key, document_name, description, is_required, status, display_order, date_created, date_updated
	FROM
		admissions_checklist_items`

	buf := bytes.NewBufferString(q)
	s.applyChecklistItemFilter(filter, data, buf)

	var dbItem checklistItemDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbItem); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.ChecklistItem{}, fmt.Errorf("db: %w", admissionsbus.ErrChecklistItemNotFound)
		}
		return admissionsbus.ChecklistItem{}, fmt.Errorf("db: %w", err)
	}

	return toBusChecklistItem(dbItem), nil
}

// CreateDocument inserts uploaded document metadata.
func (s *Store) CreateDocument(ctx context.Context, document admissionsbus.Document) error {
	const q = `
	INSERT INTO admissions_documents
		(document_id, application_id, checklist_item_id, file_name, content_type, size_bytes, storage_key, status, reviewer_id, reviewer_notes, uploaded_by_id, uploaded_at, reviewed_at, date_created, date_updated)
	VALUES
		(:document_id, :application_id, :checklist_item_id, :file_name, :content_type, :size_bytes, :storage_key, :status, :reviewer_id, :reviewer_notes, :uploaded_by_id, :uploaded_at, :reviewed_at, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBDocument(document)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateDocument replaces mutable uploaded document metadata.
func (s *Store) UpdateDocument(ctx context.Context, document admissionsbus.Document) error {
	const q = `
	UPDATE
		admissions_documents
	SET
		application_id = :application_id,
		checklist_item_id = :checklist_item_id,
		file_name = :file_name,
		content_type = :content_type,
		size_bytes = :size_bytes,
		storage_key = :storage_key,
		status = :status,
		reviewer_id = :reviewer_id,
		reviewer_notes = :reviewer_notes,
		uploaded_by_id = :uploaded_by_id,
		uploaded_at = :uploaded_at,
		reviewed_at = :reviewed_at,
		date_updated = :date_updated
	WHERE
		document_id = :document_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBDocument(document)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryDocuments retrieves uploaded document metadata.
func (s *Store) QueryDocuments(ctx context.Context, filter admissionsbus.DocumentQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Document, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		document_id, application_id, checklist_item_id, file_name, content_type, size_bytes, storage_key, status, reviewer_id, reviewer_notes, uploaded_by_id, uploaded_at, reviewed_at, date_created, date_updated
	FROM
		admissions_documents`

	buf := bytes.NewBufferString(q)
	s.applyDocumentFilter(filter, data, buf)

	orderByClause, err := documentOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbDocuments []documentDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbDocuments); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusDocuments(dbDocuments), nil
}

// CountDocuments returns the total number of uploaded document metadata records.
func (s *Store) CountDocuments(ctx context.Context, filter admissionsbus.DocumentQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_documents`

	buf := bytes.NewBufferString(q)
	s.applyDocumentFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryDocumentByID finds uploaded document metadata by ID.
func (s *Store) QueryDocumentByID(ctx context.Context, documentID uuid.UUID) (admissionsbus.Document, error) {
	filter := admissionsbus.DocumentQueryFilter{ID: &documentID}
	document, err := s.queryDocument(ctx, filter)
	if err != nil {
		return admissionsbus.Document{}, err
	}

	return document, nil
}

// CreateSyncJob is a persistence seam for SIS sync jobs. Durable storage is added with the SIS API integration slice.
func (s *Store) CreateSyncJob(context.Context, admissionsbus.SyncJob) error {
	return nil
}

// UpdateSyncJob is a persistence seam for SIS sync jobs. Durable storage is added with the SIS API integration slice.
func (s *Store) UpdateSyncJob(context.Context, admissionsbus.SyncJob) error {
	return nil
}

// QuerySyncJobs is a persistence seam for SIS sync jobs. Durable storage is added with the SIS API integration slice.
func (s *Store) QuerySyncJobs(context.Context, admissionsbus.SyncJobQueryFilter, order.By, page.Page) ([]admissionsbus.SyncJob, error) {
	return []admissionsbus.SyncJob{}, nil
}

// CountSyncJobs is a persistence seam for SIS sync jobs. Durable storage is added with the SIS API integration slice.
func (s *Store) CountSyncJobs(context.Context, admissionsbus.SyncJobQueryFilter) (int, error) {
	return 0, nil
}

// QuerySyncJobByID is a persistence seam for SIS sync jobs. Durable storage is added with the SIS API integration slice.
func (s *Store) QuerySyncJobByID(context.Context, uuid.UUID) (admissionsbus.SyncJob, error) {
	return admissionsbus.SyncJob{}, admissionsbus.ErrSyncJobNotFound
}

// CreateSyncEvent is a persistence seam for SIS sync events. Durable storage is added with the SIS API integration slice.
func (s *Store) CreateSyncEvent(context.Context, admissionsbus.SyncEvent) error {
	return nil
}

// UpdateSyncEvent is a persistence seam for SIS sync events. Durable storage is added with the SIS API integration slice.
func (s *Store) UpdateSyncEvent(context.Context, admissionsbus.SyncEvent) error {
	return nil
}

// QuerySyncEvents is a persistence seam for SIS sync events. Durable storage is added with the SIS API integration slice.
func (s *Store) QuerySyncEvents(context.Context, admissionsbus.SyncEventQueryFilter, order.By, page.Page) ([]admissionsbus.SyncEvent, error) {
	return []admissionsbus.SyncEvent{}, nil
}

// CountSyncEvents is a persistence seam for SIS sync events. Durable storage is added with the SIS API integration slice.
func (s *Store) CountSyncEvents(context.Context, admissionsbus.SyncEventQueryFilter) (int, error) {
	return 0, nil
}

// QuerySyncEventByID is a persistence seam for SIS sync events. Durable storage is added with the SIS API integration slice.
func (s *Store) QuerySyncEventByID(context.Context, uuid.UUID) (admissionsbus.SyncEvent, error) {
	return admissionsbus.SyncEvent{}, admissionsbus.ErrSyncEventNotFound
}

func (s *Store) queryDocument(ctx context.Context, filter admissionsbus.DocumentQueryFilter) (admissionsbus.Document, error) {
	data := map[string]any{}

	const q = `
	SELECT
		document_id, application_id, checklist_item_id, file_name, content_type, size_bytes, storage_key, status, reviewer_id, reviewer_notes, uploaded_by_id, uploaded_at, reviewed_at, date_created, date_updated
	FROM
		admissions_documents`

	buf := bytes.NewBufferString(q)
	s.applyDocumentFilter(filter, data, buf)

	var dbDocument documentDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbDocument); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.Document{}, fmt.Errorf("db: %w", admissionsbus.ErrDocumentNotFound)
		}
		return admissionsbus.Document{}, fmt.Errorf("db: %w", err)
	}

	return toBusDocument(dbDocument), nil
}

// CreateImportBatch inserts an admissions import batch record.
func (s *Store) CreateImportBatch(ctx context.Context, batch admissionsbus.ImportBatch) error {
	const q = `
	INSERT INTO admissions_import_batches
		(import_batch_id, source, file_type, target, status, file_name, storage_key, uploaded_by_id, total_rows, valid_rows, invalid_rows, duplicate_rows, field_mapping, invalid_report_key, validation_summary, committed_at, date_created, date_updated)
	VALUES
		(:import_batch_id, :source, :file_type, :target, :status, :file_name, :storage_key, :uploaded_by_id, :total_rows, :valid_rows, :invalid_rows, :duplicate_rows, :field_mapping, :invalid_report_key, :validation_summary, :committed_at, :date_created, :date_updated)`

	dbBatch, err := toDBImportBatch(batch)
	if err != nil {
		return fmt.Errorf("to db import batch: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbBatch); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// UpdateImportBatch replaces mutable admissions import batch fields.
func (s *Store) UpdateImportBatch(ctx context.Context, batch admissionsbus.ImportBatch) error {
	const q = `
	UPDATE
		admissions_import_batches
	SET
		source = :source,
		file_type = :file_type,
		target = :target,
		status = :status,
		file_name = :file_name,
		storage_key = :storage_key,
		uploaded_by_id = :uploaded_by_id,
		total_rows = :total_rows,
		valid_rows = :valid_rows,
		invalid_rows = :invalid_rows,
		duplicate_rows = :duplicate_rows,
		field_mapping = :field_mapping,
		invalid_report_key = :invalid_report_key,
		validation_summary = :validation_summary,
		committed_at = :committed_at,
		date_updated = :date_updated
	WHERE
		import_batch_id = :import_batch_id`

	dbBatch, err := toDBImportBatch(batch)
	if err != nil {
		return fmt.Errorf("to db import batch: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbBatch); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryImportBatches retrieves admissions import batch records.
func (s *Store) QueryImportBatches(ctx context.Context, filter admissionsbus.ImportBatchQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ImportBatch, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		import_batch_id, source, file_type, target, status, file_name, storage_key, uploaded_by_id, total_rows, valid_rows, invalid_rows, duplicate_rows, field_mapping, invalid_report_key, validation_summary, committed_at, date_created, date_updated
	FROM
		admissions_import_batches`

	buf := bytes.NewBufferString(q)
	s.applyImportBatchFilter(filter, data, buf)

	orderByClause, err := importBatchOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbBatches []importBatchDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbBatches); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusImportBatches(dbBatches)
}

// CountImportBatches returns the total number of admissions import batch records.
func (s *Store) CountImportBatches(ctx context.Context, filter admissionsbus.ImportBatchQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_import_batches`

	buf := bytes.NewBufferString(q)
	s.applyImportBatchFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryImportBatchByID finds an admissions import batch by ID.
func (s *Store) QueryImportBatchByID(ctx context.Context, batchID uuid.UUID) (admissionsbus.ImportBatch, error) {
	filter := admissionsbus.ImportBatchQueryFilter{ID: &batchID}
	batch, err := s.queryImportBatch(ctx, filter)
	if err != nil {
		return admissionsbus.ImportBatch{}, err
	}

	return batch, nil
}

func (s *Store) queryImportBatch(ctx context.Context, filter admissionsbus.ImportBatchQueryFilter) (admissionsbus.ImportBatch, error) {
	data := map[string]any{}

	const q = `
	SELECT
		import_batch_id, source, file_type, target, status, file_name, storage_key, uploaded_by_id, total_rows, valid_rows, invalid_rows, duplicate_rows, field_mapping, invalid_report_key, validation_summary, committed_at, date_created, date_updated
	FROM
		admissions_import_batches`

	buf := bytes.NewBufferString(q)
	s.applyImportBatchFilter(filter, data, buf)

	var dbBatch importBatchDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbBatch); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.ImportBatch{}, fmt.Errorf("db: %w", admissionsbus.ErrImportBatchNotFound)
		}
		return admissionsbus.ImportBatch{}, fmt.Errorf("db: %w", err)
	}

	return toBusImportBatch(dbBatch)
}

// CreateImportInvalidRows inserts invalid import rows for correction reports.
func (s *Store) CreateImportInvalidRows(ctx context.Context, rows []admissionsbus.ImportInvalidRow) error {
	if len(rows) == 0 {
		return nil
	}

	const q = `
	INSERT INTO admissions_import_invalid_rows
		(import_invalid_row_id, import_batch_id, row_number, field_name, raw_data, error_code, error_detail, date_created)
	VALUES
		(:import_invalid_row_id, :import_batch_id, :row_number, :field_name, :raw_data, :error_code, :error_detail, :date_created)`

	dbRows, err := toDBImportInvalidRows(rows)
	if err != nil {
		return fmt.Errorf("to db import invalid rows: %w", err)
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbRows); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryImportInvalidRows retrieves invalid import rows for correction reports.
func (s *Store) QueryImportInvalidRows(ctx context.Context, filter admissionsbus.ImportInvalidRowQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ImportInvalidRow, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		import_invalid_row_id, import_batch_id, row_number, field_name, raw_data, error_code, error_detail, date_created
	FROM
		admissions_import_invalid_rows`

	buf := bytes.NewBufferString(q)
	s.applyImportInvalidRowFilter(filter, data, buf)

	orderByClause, err := importInvalidRowOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbRows []importInvalidRowDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbRows); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusImportInvalidRows(dbRows)
}

// CountImportInvalidRows returns the total number of invalid import rows.
func (s *Store) CountImportInvalidRows(ctx context.Context, filter admissionsbus.ImportInvalidRowQueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		admissions_import_invalid_rows`

	buf := bytes.NewBufferString(q)
	s.applyImportInvalidRowFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryImportInvalidRowByID finds an invalid import row by ID.
func (s *Store) QueryImportInvalidRowByID(ctx context.Context, rowID uuid.UUID) (admissionsbus.ImportInvalidRow, error) {
	filter := admissionsbus.ImportInvalidRowQueryFilter{ID: &rowID}
	row, err := s.queryImportInvalidRow(ctx, filter)
	if err != nil {
		return admissionsbus.ImportInvalidRow{}, err
	}

	return row, nil
}

func (s *Store) queryImportInvalidRow(ctx context.Context, filter admissionsbus.ImportInvalidRowQueryFilter) (admissionsbus.ImportInvalidRow, error) {
	data := map[string]any{}

	const q = `
	SELECT
		import_invalid_row_id, import_batch_id, row_number, field_name, raw_data, error_code, error_detail, date_created
	FROM
		admissions_import_invalid_rows`

	buf := bytes.NewBufferString(q)
	s.applyImportInvalidRowFilter(filter, data, buf)

	var dbRow importInvalidRowDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbRow); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.ImportInvalidRow{}, fmt.Errorf("db: %w", admissionsbus.ErrImportInvalidRowNotFound)
		}
		return admissionsbus.ImportInvalidRow{}, fmt.Errorf("db: %w", err)
	}

	return toBusImportInvalidRow(dbRow)
}

func (s *Store) queryAcademicTerm(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter) (admissionsbus.AcademicTerm, error) {
	data := map[string]any{}

	const q = `
	SELECT
		academic_term_id, external_sis_id, name, code, term_type, start_date, end_date, application_start_date, application_deadline, is_active, synced_at, date_created, date_updated
	FROM
		admissions_academic_terms`

	buf := bytes.NewBufferString(q)
	s.applyAcademicTermFilter(filter, data, buf)

	var dbTerm academicTermDB
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &dbTerm); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return admissionsbus.AcademicTerm{}, fmt.Errorf("db: %w", admissionsbus.ErrAcademicTermNotFound)
		}
		return admissionsbus.AcademicTerm{}, fmt.Errorf("db: %w", err)
	}

	return toBusAcademicTerm(dbTerm), nil
}
