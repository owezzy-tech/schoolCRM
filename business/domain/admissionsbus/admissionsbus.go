// Package admissionsbus provides business access to the admissions domain.
package admissionsbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/sdk/delegate"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Set of error variables for admissions reference data operations.
var (
	ErrProgramNotFound          = errors.New("program not found")
	ErrAcademicTermNotFound     = errors.New("academic term not found")
	ErrInvalidTermDateRange     = errors.New("term start date must be before end date")
	ErrInvalidApplicationWindow = errors.New("application deadline must be on or after application start date")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve admissions data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Health(ctx context.Context) (Health, error)
	UpsertProgram(ctx context.Context, prg Program) error
	QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error)
	CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error)
	QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error)
	QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error)
	UpsertAcademicTerm(ctx context.Context, term AcademicTerm) error
	QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error)
	CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error)
	QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error)
	QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Health(ctx context.Context) (Health, error)
	UpsertProgram(ctx context.Context, up UpsertProgram) (Program, error)
	QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error)
	CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error)
	QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error)
	QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error)
	UpsertAcademicTerm(ctx context.Context, up UpsertAcademicTerm) (AcademicTerm, error)
	QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error)
	CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error)
	QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error)
	QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for admissions access.
type Business struct {
	log        *logger.Logger
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs an admissions business API for use.
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		delegate:   delegate,
		storer:     storer,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, b.delegate, storer, b.extensions...)

	return nb, nil
}

// Health returns the current admissions context scaffold metadata.
func (b *Business) Health(ctx context.Context) (Health, error) {
	health, err := b.storer.Health(ctx)
	if err != nil {
		return Health{}, fmt.Errorf("health: %w", err)
	}

	return health, nil
}

// UpsertProgram creates or updates SIS-owned Program reference data for sync/import paths.
func (b *Business) UpsertProgram(ctx context.Context, up UpsertProgram) (Program, error) {
	now := time.Now()
	id := uuid.New()
	if up.ID != nil {
		id = *up.ID
	}

	prg := Program{
		ID:            id,
		ExternalSISID: up.ExternalSISID,
		Name:          up.Name,
		Code:          up.Code,
		Description:   up.Description,
		DegreeLevel:   up.DegreeLevel,
		Active:        up.Active,
		SyncedAt:      up.SyncedAt,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := b.storer.UpsertProgram(ctx, prg); err != nil {
		return Program{}, fmt.Errorf("upsert program: %w", err)
	}

	return b.QueryProgramByExternalSISID(ctx, up.ExternalSISID)
}

// QueryPrograms retrieves a list of existing Program reference records.
func (b *Business) QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error) {
	programs, err := b.storer.QueryPrograms(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query programs: %w", err)
	}

	return programs, nil
}

// CountPrograms returns the total number of Program reference records.
func (b *Business) CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error) {
	return b.storer.CountPrograms(ctx, filter)
}

// QueryProgramByID finds a Program by ID.
func (b *Business) QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error) {
	program, err := b.storer.QueryProgramByID(ctx, programID)
	if err != nil {
		return Program{}, fmt.Errorf("query program: programID[%s]: %w", programID, err)
	}

	return program, nil
}

// QueryProgramByExternalSISID finds a Program by immutable SIS ID.
func (b *Business) QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error) {
	program, err := b.storer.QueryProgramByExternalSISID(ctx, externalSISID)
	if err != nil {
		return Program{}, fmt.Errorf("query program: externalSISID[%s]: %w", externalSISID, err)
	}

	return program, nil
}

// UpsertAcademicTerm creates or updates SIS-owned AcademicTerm reference data for sync/import paths.
func (b *Business) UpsertAcademicTerm(ctx context.Context, up UpsertAcademicTerm) (AcademicTerm, error) {
	if !up.StartDate.Before(up.EndDate) {
		return AcademicTerm{}, ErrInvalidTermDateRange
	}

	if up.ApplicationStartDate != nil && up.ApplicationDeadline != nil && up.ApplicationDeadline.Before(*up.ApplicationStartDate) {
		return AcademicTerm{}, ErrInvalidApplicationWindow
	}

	now := time.Now()
	id := uuid.New()
	if up.ID != nil {
		id = *up.ID
	}

	term := AcademicTerm{
		ID:                   id,
		ExternalSISID:        up.ExternalSISID,
		Name:                 up.Name,
		Code:                 up.Code,
		TermType:             up.TermType,
		StartDate:            up.StartDate,
		EndDate:              up.EndDate,
		ApplicationStartDate: up.ApplicationStartDate,
		ApplicationDeadline:  up.ApplicationDeadline,
		Active:               up.Active,
		SyncedAt:             up.SyncedAt,
		DateCreated:          now,
		DateUpdated:          now,
	}

	if err := b.storer.UpsertAcademicTerm(ctx, term); err != nil {
		return AcademicTerm{}, fmt.Errorf("upsert academic term: %w", err)
	}

	return b.QueryAcademicTermByExternalSISID(ctx, up.ExternalSISID)
}

// QueryAcademicTerms retrieves a list of existing AcademicTerm reference records.
func (b *Business) QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error) {
	terms, err := b.storer.QueryAcademicTerms(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query academic terms: %w", err)
	}

	return terms, nil
}

// CountAcademicTerms returns the total number of AcademicTerm reference records.
func (b *Business) CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error) {
	return b.storer.CountAcademicTerms(ctx, filter)
}

// QueryAcademicTermByID finds an AcademicTerm by ID.
func (b *Business) QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error) {
	term, err := b.storer.QueryAcademicTermByID(ctx, termID)
	if err != nil {
		return AcademicTerm{}, fmt.Errorf("query academic term: termID[%s]: %w", termID, err)
	}

	return term, nil
}

// QueryAcademicTermByExternalSISID finds an AcademicTerm by immutable SIS ID.
func (b *Business) QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error) {
	term, err := b.storer.QueryAcademicTermByExternalSISID(ctx, externalSISID)
	if err != nil {
		return AcademicTerm{}, fmt.Errorf("query academic term: externalSISID[%s]: %w", externalSISID, err)
	}

	return term, nil
}
