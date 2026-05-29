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
