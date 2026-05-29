// Package admissionsdb contains admissions related persistence functionality.
package admissionsdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

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
