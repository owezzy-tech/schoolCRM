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
		(constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated)
	VALUES
		(:constituent_id, :first_name, :last_name, :preferred_name, :middle_name, :suffix, :date_of_birth, :primary_email, :primary_phone, :external_sis_id, :lifecycle_stage, :duplicate_status, :duplicate_of_id, :sis_synced_at, :date_created, :date_updated)`

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
		constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated
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

func (s *Store) queryConstituent(ctx context.Context, filter admissionsbus.ConstituentQueryFilter) (admissionsbus.Constituent, error) {
	data := map[string]any{}

	const q = `
	SELECT
		constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated
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
		(application_id, constituent_id, program_id, academic_term_id, application_type, status, assigned_reviewer_id, submitted_at, date_created, date_updated)
	VALUES
		(:application_id, :constituent_id, :program_id, :academic_term_id, :application_type, :status, :assigned_reviewer_id, :submitted_at, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBApplication(app)); err != nil {
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
		assigned_reviewer_id = :assigned_reviewer_id,
		submitted_at = :submitted_at,
		date_updated = :date_updated
	WHERE
		application_id = :application_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBApplication(app)); err != nil {
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
		application_id, constituent_id, program_id, academic_term_id, application_type, status, assigned_reviewer_id, submitted_at, date_created, date_updated
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

	return toBusApplications(dbApplications), nil
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
		application_id, constituent_id, program_id, academic_term_id, application_type, status, assigned_reviewer_id, submitted_at, date_created, date_updated
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

	return toBusApplication(dbApplication), nil
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
