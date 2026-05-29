// Package admissionsbus provides business access to the admissions domain.
package admissionsbus

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	ErrConstituentNotFound      = errors.New("constituent not found")
	ErrFirstNameRequired        = errors.New("first name required")
	ErrLastNameRequired         = errors.New("last name required")
	ErrDateOfBirthRequired      = errors.New("date of birth required")
	ErrDateOfBirthInFuture      = errors.New("date of birth cannot be in the future")
	ErrPrimaryPhoneRequired     = errors.New("primary phone required")
	ErrInvalidLifecycleStage    = errors.New("invalid lifecycle stage")
	ErrInvalidDuplicateStatus   = errors.New("invalid duplicate status")
	ErrInvalidDuplicateLink     = errors.New("duplicate status does not match duplicate link")
	ErrInvalidLifecycleChange   = errors.New("invalid lifecycle stage change")
	ErrProgramNotFound          = errors.New("program not found")
	ErrAcademicTermNotFound     = errors.New("academic term not found")
	ErrInvalidTermDateRange     = errors.New("term start date must be before end date")
	ErrInvalidApplicationWindow = errors.New("application deadline must be on or after application start date")
	ErrDuplicateReviewNotFound  = errors.New("duplicate review not found")
	ErrInvalidDuplicateReview   = errors.New("invalid duplicate review")
	ErrInvalidMatchType         = errors.New("invalid duplicate match type")
	ErrInvalidMatchScore        = errors.New("duplicate match score must be between 0 and 100")
	ErrMatchReasonRequired      = errors.New("duplicate match reason required")
	ErrInvalidReviewStatus      = errors.New("invalid duplicate review status")
	ErrInvalidResolution        = errors.New("invalid duplicate review resolution")
	ErrDuplicateReviewResolved  = errors.New("duplicate review already resolved")
	ErrResolutionActorRequired  = errors.New("resolution actor required")
	ErrApplicationNotFound      = errors.New("application not found")
	ErrInvalidApplicationType   = errors.New("invalid application type")
	ErrInvalidApplicationStatus = errors.New("invalid application status")
	ErrDuplicateApplication     = errors.New("active application already exists for constituent term and program")
	ErrConstituentIDRequired    = errors.New("constituent id required")
	ErrProgramIDRequired        = errors.New("program id required")
	ErrAcademicTermIDRequired   = errors.New("academic term id required")
	ErrInactiveProgram          = errors.New("program is inactive")
	ErrInactiveAcademicTerm     = errors.New("academic term is inactive")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve admissions data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Health(ctx context.Context) (Health, error)
	CreateConstituent(ctx context.Context, cst Constituent) error
	UpdateConstituent(ctx context.Context, cst Constituent) error
	QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error)
	CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error)
	QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error)
	QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error)
	QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error)
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
	CreateDuplicateReview(ctx context.Context, review DuplicateReview) error
	UpdateDuplicateReview(ctx context.Context, review DuplicateReview) error
	QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error)
	CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error)
	QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error)
	CreateApplication(ctx context.Context, app Application) error
	QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error)
	CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error)
	QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error)
	QueryActiveApplicationByTuple(ctx context.Context, constituentID uuid.UUID, academicTermID uuid.UUID, programID uuid.UUID) (Application, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Health(ctx context.Context) (Health, error)
	CreateConstituent(ctx context.Context, nc NewConstituent) (Constituent, error)
	UpdateConstituent(ctx context.Context, cst Constituent, uc UpdateConstituent) (Constituent, error)
	QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error)
	CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error)
	QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error)
	QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error)
	QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error)
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
	CreateDuplicateReview(ctx context.Context, nr NewDuplicateReview) (DuplicateReview, error)
	ResolveDuplicateReview(ctx context.Context, review DuplicateReview, rr ResolveDuplicateReview) (DuplicateReview, error)
	QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error)
	CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error)
	QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error)
	CreateApplication(ctx context.Context, na NewApplication) (Application, error)
	QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error)
	CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error)
	QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error)
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

// CreateConstituent adds a new Constituent to the admissions context.
func (b *Business) CreateConstituent(ctx context.Context, nc NewConstituent) (Constituent, error) {
	stage := nc.LifecycleStage
	if stage == "" {
		stage = LifecycleStageProspect
	}

	status := nc.DuplicateStatus
	if status == "" {
		status = DuplicateStatusActive
	}

	if err := validateRequiredConstituentFields(nc.FirstName, nc.LastName, nc.DateOfBirth, nc.PrimaryPhone); err != nil {
		return Constituent{}, err
	}

	if err := validateLifecycleStage(stage); err != nil {
		return Constituent{}, err
	}

	if err := validateDuplicateStatus(status, nc.DuplicateOfID); err != nil {
		return Constituent{}, err
	}

	now := time.Now()
	cst := Constituent{
		ID:              uuid.New(),
		FirstName:       strings.TrimSpace(nc.FirstName),
		LastName:        strings.TrimSpace(nc.LastName),
		PreferredName:   trimStringPtr(nc.PreferredName),
		MiddleName:      trimStringPtr(nc.MiddleName),
		Suffix:          trimStringPtr(nc.Suffix),
		DateOfBirth:     nc.DateOfBirth,
		PrimaryEmail:    nc.PrimaryEmail,
		PrimaryPhone:    strings.TrimSpace(nc.PrimaryPhone),
		ExternalSISID:   trimStringPtr(nc.ExternalSISID),
		LifecycleStage:  stage,
		DuplicateStatus: status,
		DuplicateOfID:   nc.DuplicateOfID,
		SISSyncedAt:     nc.SISSyncedAt,
		DateCreated:     now,
		DateUpdated:     now,
	}

	match, err := b.queryTrustedExactDuplicate(ctx, cst)
	if err != nil {
		return Constituent{}, err
	}

	if match != nil {
		if cst.ExternalSISID != nil && match.ExternalSISID != nil && *cst.ExternalSISID == *match.ExternalSISID {
			return *match, nil
		}

		cst.DuplicateStatus = DuplicateStatusDuplicateOf
		cst.DuplicateOfID = &match.ID
	}

	if err := b.storer.CreateConstituent(ctx, cst); err != nil {
		return Constituent{}, fmt.Errorf("create constituent: %w", err)
	}

	return cst, nil
}

// UpdateConstituent modifies mutable information for a Constituent.
func (b *Business) UpdateConstituent(ctx context.Context, cst Constituent, uc UpdateConstituent) (Constituent, error) {
	if uc.PreferredName != nil {
		cst.PreferredName = trimStringPtr(uc.PreferredName)
	}

	if uc.MiddleName != nil {
		cst.MiddleName = trimStringPtr(uc.MiddleName)
	}

	if uc.Suffix != nil {
		cst.Suffix = trimStringPtr(uc.Suffix)
	}

	if uc.PrimaryEmail != nil {
		cst.PrimaryEmail = *uc.PrimaryEmail
	}

	if uc.PrimaryPhone != nil {
		phone := strings.TrimSpace(*uc.PrimaryPhone)
		if phone == "" {
			return Constituent{}, ErrPrimaryPhoneRequired
		}
		cst.PrimaryPhone = phone
	}

	if uc.LifecycleStage != nil {
		if err := validateLifecycleStage(*uc.LifecycleStage); err != nil {
			return Constituent{}, err
		}

		if !canChangeLifecycleStage(cst.LifecycleStage, *uc.LifecycleStage) {
			return Constituent{}, ErrInvalidLifecycleChange
		}

		cst.LifecycleStage = *uc.LifecycleStage
	}

	if uc.DuplicateStatus != nil || uc.DuplicateOfID != nil {
		status := cst.DuplicateStatus
		if uc.DuplicateStatus != nil {
			status = *uc.DuplicateStatus
		}

		duplicateOfID := cst.DuplicateOfID
		if uc.DuplicateOfID != nil {
			duplicateOfID = uc.DuplicateOfID
		}

		if err := validateDuplicateStatus(status, duplicateOfID); err != nil {
			return Constituent{}, err
		}

		cst.DuplicateStatus = status
		cst.DuplicateOfID = duplicateOfID
	}

	if uc.SISSyncedAt != nil {
		cst.SISSyncedAt = uc.SISSyncedAt
	}

	cst.DateUpdated = time.Now()

	if err := b.storer.UpdateConstituent(ctx, cst); err != nil {
		return Constituent{}, fmt.Errorf("update constituent: %w", err)
	}

	return cst, nil
}

// QueryConstituents retrieves a list of existing constituents.
func (b *Business) QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error) {
	constituents, err := b.storer.QueryConstituents(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query constituents: %w", err)
	}

	return constituents, nil
}

// CountConstituents returns the total number of constituents.
func (b *Business) CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error) {
	return b.storer.CountConstituents(ctx, filter)
}

// QueryConstituentByID finds a Constituent by ID.
func (b *Business) QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: constituentID[%s]: %w", constituentID, err)
	}

	return cst, nil
}

// QueryConstituentByPrimaryEmail finds a Constituent by primary email.
func (b *Business) QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByPrimaryEmail(ctx, email)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: primaryEmail[%s]: %w", email, err)
	}

	return cst, nil
}

// QueryConstituentByExternalSISID finds a Constituent by SIS ID.
func (b *Business) QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByExternalSISID(ctx, externalSISID)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: externalSISID[%s]: %w", externalSISID, err)
	}

	return cst, nil
}

func (b *Business) queryTrustedExactDuplicate(ctx context.Context, cst Constituent) (*Constituent, error) {
	if cst.ExternalSISID != nil {
		match, err := b.storer.QueryConstituentByExternalSISID(ctx, *cst.ExternalSISID)
		if err != nil && !errors.Is(err, ErrConstituentNotFound) {
			return nil, fmt.Errorf("query external sis duplicate: %w", err)
		}

		if err == nil && match.ID != cst.ID {
			return &match, nil
		}
	}

	match, err := b.storer.QueryConstituentByPrimaryEmail(ctx, cst.PrimaryEmail.String())
	if err != nil && !errors.Is(err, ErrConstituentNotFound) {
		return nil, fmt.Errorf("query email duplicate: %w", err)
	}

	if err == nil && match.ID != cst.ID {
		return &match, nil
	}

	return nil, nil
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

// CreateDuplicateReview adds a possible duplicate pair to the staff review queue.
func (b *Business) CreateDuplicateReview(ctx context.Context, nr NewDuplicateReview) (DuplicateReview, error) {
	if err := validateNewDuplicateReview(nr); err != nil {
		return DuplicateReview{}, err
	}

	now := time.Now()
	review := DuplicateReview{
		ID:                     uuid.New(),
		SourceConstituentID:    nr.SourceConstituentID,
		CandidateConstituentID: nr.CandidateConstituentID,
		MatchType:              nr.MatchType,
		MatchScore:             nr.MatchScore,
		MatchReason:            strings.TrimSpace(nr.MatchReason),
		Status:                 DuplicateReviewStatusPending,
		DateCreated:            now,
		DateUpdated:            now,
	}

	if err := b.storer.CreateDuplicateReview(ctx, review); err != nil {
		return DuplicateReview{}, fmt.Errorf("create duplicate review: %w", err)
	}

	return review, nil
}

// ResolveDuplicateReview records a staff decision for a pending duplicate review.
func (b *Business) ResolveDuplicateReview(ctx context.Context, review DuplicateReview, rr ResolveDuplicateReview) (DuplicateReview, error) {
	if review.Status != DuplicateReviewStatusPending && review.Status != DuplicateReviewStatusDeferred {
		return DuplicateReview{}, ErrDuplicateReviewResolved
	}

	if rr.ActorID == uuid.Nil {
		return DuplicateReview{}, ErrResolutionActorRequired
	}

	status, err := statusForResolution(rr.Resolution)
	if err != nil {
		return DuplicateReview{}, err
	}

	now := time.Now()
	review.Status = status
	review.ResolvedBy = &rr.ActorID
	review.ResolvedAt = &now
	review.ResolutionNote = trimStringPtr(rr.Note)
	review.DateUpdated = now

	if err := b.storer.UpdateDuplicateReview(ctx, review); err != nil {
		return DuplicateReview{}, fmt.Errorf("update duplicate review: %w", err)
	}

	if rr.Resolution == DuplicateReviewResolutionLink || rr.Resolution == DuplicateReviewResolutionMerge {
		source, err := b.storer.QueryConstituentByID(ctx, review.SourceConstituentID)
		if err != nil {
			return DuplicateReview{}, fmt.Errorf("query source constituent: %w", err)
		}

		candidate, err := b.storer.QueryConstituentByID(ctx, review.CandidateConstituentID)
		if err != nil {
			return DuplicateReview{}, fmt.Errorf("query candidate constituent: %w", err)
		}

		duplicateStatus := DuplicateStatusDuplicateOf
		if rr.Resolution == DuplicateReviewResolutionMerge {
			duplicateStatus = DuplicateStatusMerged
		}

		source.DuplicateStatus = duplicateStatus
		source.DuplicateOfID = &candidate.ID
		source.DateUpdated = now

		if err := b.storer.UpdateConstituent(ctx, source); err != nil {
			return DuplicateReview{}, fmt.Errorf("update source duplicate link: %w", err)
		}
	}

	return review, nil
}

// QueryDuplicateReviews retrieves a list of existing duplicate reviews.
func (b *Business) QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error) {
	reviews, err := b.storer.QueryDuplicateReviews(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query duplicate reviews: %w", err)
	}

	return reviews, nil
}

// CountDuplicateReviews returns the total number of duplicate reviews.
func (b *Business) CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error) {
	return b.storer.CountDuplicateReviews(ctx, filter)
}

// QueryDuplicateReviewByID finds a DuplicateReview by ID.
func (b *Business) QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error) {
	review, err := b.storer.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return DuplicateReview{}, fmt.Errorf("query duplicate review: reviewID[%s]: %w", reviewID, err)
	}

	return review, nil
}

// CreateApplication adds a draft application while enforcing active application uniqueness.
func (b *Business) CreateApplication(ctx context.Context, na NewApplication) (Application, error) {
	if err := validateNewApplication(na); err != nil {
		return Application{}, err
	}

	if _, err := b.storer.QueryConstituentByID(ctx, na.ConstituentID); err != nil {
		return Application{}, fmt.Errorf("query constituent: %w", err)
	}

	program, err := b.storer.QueryProgramByID(ctx, na.ProgramID)
	if err != nil {
		return Application{}, fmt.Errorf("query program: %w", err)
	}
	if !program.Active {
		return Application{}, ErrInactiveProgram
	}

	term, err := b.storer.QueryAcademicTermByID(ctx, na.AcademicTermID)
	if err != nil {
		return Application{}, fmt.Errorf("query academic term: %w", err)
	}
	if !term.Active {
		return Application{}, ErrInactiveAcademicTerm
	}

	if _, err := b.storer.QueryActiveApplicationByTuple(ctx, na.ConstituentID, na.AcademicTermID, na.ProgramID); err == nil {
		return Application{}, ErrDuplicateApplication
	} else if !errors.Is(err, ErrApplicationNotFound) {
		return Application{}, fmt.Errorf("query active application: %w", err)
	}

	now := time.Now()
	app := Application{
		ID:                 uuid.New(),
		ConstituentID:      na.ConstituentID,
		ProgramID:          na.ProgramID,
		AcademicTermID:     na.AcademicTermID,
		ApplicationType:    na.ApplicationType,
		Status:             ApplicationStatusDraft,
		AssignedReviewerID: na.AssignedReviewerID,
		DateCreated:        now,
		DateUpdated:        now,
	}

	if err := b.storer.CreateApplication(ctx, app); err != nil {
		return Application{}, fmt.Errorf("create application: %w", err)
	}

	return app, nil
}

// QueryApplications retrieves a list of existing applications.
func (b *Business) QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error) {
	applications, err := b.storer.QueryApplications(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query applications: %w", err)
	}

	return applications, nil
}

// CountApplications returns the total number of applications.
func (b *Business) CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error) {
	return b.storer.CountApplications(ctx, filter)
}

// QueryApplicationByID finds an Application by ID.
func (b *Business) QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error) {
	app, err := b.storer.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return Application{}, fmt.Errorf("query application: applicationID[%s]: %w", applicationID, err)
	}

	return app, nil
}

func validateRequiredConstituentFields(firstName string, lastName string, dob time.Time, primaryPhone string) error {
	if strings.TrimSpace(firstName) == "" {
		return ErrFirstNameRequired
	}

	if strings.TrimSpace(lastName) == "" {
		return ErrLastNameRequired
	}

	if dob.IsZero() {
		return ErrDateOfBirthRequired
	}

	if dob.After(time.Now()) {
		return ErrDateOfBirthInFuture
	}

	if strings.TrimSpace(primaryPhone) == "" {
		return ErrPrimaryPhoneRequired
	}

	return nil
}

func validateLifecycleStage(stage LifecycleStage) error {
	switch stage {
	case LifecycleStageProspect,
		LifecycleStageInquiry,
		LifecycleStageApplicant,
		LifecycleStageAdmitted,
		LifecycleStageEnrolled,
		LifecycleStageAlumni:
		return nil
	default:
		return ErrInvalidLifecycleStage
	}
}

func canChangeLifecycleStage(from LifecycleStage, to LifecycleStage) bool {
	if from == to {
		return true
	}

	if to == LifecycleStageAlumni {
		return true
	}

	transitions := map[LifecycleStage]LifecycleStage{
		LifecycleStageProspect:  LifecycleStageInquiry,
		LifecycleStageInquiry:   LifecycleStageApplicant,
		LifecycleStageApplicant: LifecycleStageAdmitted,
		LifecycleStageAdmitted:  LifecycleStageEnrolled,
	}

	return transitions[from] == to
}

func validateDuplicateStatus(status DuplicateStatus, duplicateOfID *uuid.UUID) error {
	switch status {
	case DuplicateStatusActive:
		if duplicateOfID != nil {
			return ErrInvalidDuplicateLink
		}
		return nil
	case DuplicateStatusMerged, DuplicateStatusDuplicateOf:
		if duplicateOfID == nil {
			return ErrInvalidDuplicateLink
		}
		return nil
	default:
		return ErrInvalidDuplicateStatus
	}
}

func validateNewDuplicateReview(nr NewDuplicateReview) error {
	if nr.SourceConstituentID == uuid.Nil || nr.CandidateConstituentID == uuid.Nil || nr.SourceConstituentID == nr.CandidateConstituentID {
		return ErrInvalidDuplicateReview
	}

	if err := validateDuplicateReviewMatchType(nr.MatchType); err != nil {
		return err
	}

	if nr.MatchScore < 0 || nr.MatchScore > 100 {
		return ErrInvalidMatchScore
	}

	if strings.TrimSpace(nr.MatchReason) == "" {
		return ErrMatchReasonRequired
	}

	return nil
}

func validateDuplicateReviewMatchType(matchType DuplicateReviewMatchType) error {
	switch matchType {
	case DuplicateReviewMatchTypeExact, DuplicateReviewMatchTypeFuzzy:
		return nil
	default:
		return ErrInvalidMatchType
	}
}

func validateDuplicateReviewStatus(status DuplicateReviewStatus) error {
	switch status {
	case DuplicateReviewStatusPending,
		DuplicateReviewStatusLinked,
		DuplicateReviewStatusMerged,
		DuplicateReviewStatusRejected,
		DuplicateReviewStatusDeferred:
		return nil
	default:
		return ErrInvalidReviewStatus
	}
}

func statusForResolution(resolution DuplicateReviewResolution) (DuplicateReviewStatus, error) {
	switch resolution {
	case DuplicateReviewResolutionLink:
		return DuplicateReviewStatusLinked, nil
	case DuplicateReviewResolutionMerge:
		return DuplicateReviewStatusMerged, nil
	case DuplicateReviewResolutionReject:
		return DuplicateReviewStatusRejected, nil
	case DuplicateReviewResolutionDefer:
		return DuplicateReviewStatusDeferred, nil
	default:
		return "", ErrInvalidResolution
	}
}

func validateNewApplication(na NewApplication) error {
	if na.ConstituentID == uuid.Nil {
		return ErrConstituentIDRequired
	}

	if na.ProgramID == uuid.Nil {
		return ErrProgramIDRequired
	}

	if na.AcademicTermID == uuid.Nil {
		return ErrAcademicTermIDRequired
	}

	return validateApplicationType(na.ApplicationType)
}

func validateApplicationType(applicationType ApplicationType) error {
	switch applicationType {
	case ApplicationTypeFreshman,
		ApplicationTypeTransfer,
		ApplicationTypeGraduate:
		return nil
	default:
		return ErrInvalidApplicationType
	}
}

func validateApplicationStatus(status ApplicationStatus) error {
	switch status {
	case ApplicationStatusDraft,
		ApplicationStatusSubmitted,
		ApplicationStatusAwaitingDocuments,
		ApplicationStatusReadyForReview,
		ApplicationStatusInReview,
		ApplicationStatusDecisionPending,
		ApplicationStatusAdmitted,
		ApplicationStatusDenied,
		ApplicationStatusWaitlisted,
		ApplicationStatusDeferred,
		ApplicationStatusWithdrawn,
		ApplicationStatusEnrolled:
		return nil
	default:
		return ErrInvalidApplicationStatus
	}
}

func isApplicationActive(status ApplicationStatus) bool {
	switch status {
	case ApplicationStatusDenied,
		ApplicationStatusWithdrawn,
		ApplicationStatusEnrolled:
		return false
	default:
		return true
	}
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
