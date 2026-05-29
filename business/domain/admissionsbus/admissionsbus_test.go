package admissionsbus

import (
	"context"
	"errors"
	"net/mail"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

func TestUpsertAcademicTermValidatesTermDateRange(t *testing.T) {
	t.Parallel()

	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
	now := time.Now()

	_, err := bus.UpsertAcademicTerm(context.Background(), UpsertAcademicTerm{
		ExternalSISID: "TERM-2026-FALL",
		Name:          "Fall 2026",
		Code:          "202609",
		StartDate:     now,
		EndDate:       now,
		Active:        true,
	})

	if !errors.Is(err, ErrInvalidTermDateRange) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidTermDateRange)
	}
}

func TestUpsertAcademicTermValidatesApplicationWindow(t *testing.T) {
	t.Parallel()

	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
	start := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	appStart := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	deadline := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	_, err := bus.UpsertAcademicTerm(context.Background(), UpsertAcademicTerm{
		ExternalSISID:        "TERM-2026-SPRING",
		Name:                 "Spring 2026",
		Code:                 "202601",
		StartDate:            start,
		EndDate:              end,
		ApplicationStartDate: &appStart,
		ApplicationDeadline:  &deadline,
		Active:               true,
	})

	if !errors.Is(err, ErrInvalidApplicationWindow) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidApplicationWindow)
	}
}

func TestCreateConstituentRequiresIdentityFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "applicant@example.com"}

	tests := []struct {
		name string
		nc   NewConstituent
		want error
	}{
		{
			name: "first name",
			nc: NewConstituent{
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrFirstNameRequired,
		},
		{
			name: "last name",
			nc: NewConstituent{
				FirstName:    "Ada",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrLastNameRequired,
		},
		{
			name: "date of birth",
			nc: NewConstituent{
				FirstName:    "Ada",
				LastName:     "Applicant",
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrDateOfBirthRequired,
		},
		{
			name: "primary phone",
			nc: NewConstituent{
				FirstName:    "Ada",
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
			},
			want: ErrPrimaryPhoneRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateConstituent(context.Background(), tt.nc)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateConstituentAutoLinksExactEmailDuplicate(t *testing.T) {
	t.Parallel()

	email := mail.Address{Address: "applicant@example.com"}
	matchID := uuid.New()
	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{},
		constituentByEmail: map[string]Constituent{
			email.String(): {ID: matchID, PrimaryEmail: email},
		},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	created, err := bus.CreateConstituent(context.Background(), NewConstituent{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
	})
	if err != nil {
		t.Fatalf("CreateConstituent returned error: %v", err)
	}

	if created.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("DuplicateStatus = %s, want %s", created.DuplicateStatus, DuplicateStatusDuplicateOf)
	}

	if created.DuplicateOfID == nil || *created.DuplicateOfID != matchID {
		t.Fatalf("DuplicateOfID = %v, want %s", created.DuplicateOfID, matchID)
	}

	stored := store.constituents[created.ID]
	if stored.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("stored DuplicateStatus = %s, want %s", stored.DuplicateStatus, DuplicateStatusDuplicateOf)
	}
}

func TestCreateConstituentIgnoresMissingExactDuplicate(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "applicant@example.com"}

	_, err := bus.CreateConstituent(context.Background(), NewConstituent{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
	})
	if err != nil {
		t.Fatalf("CreateConstituent returned error: %v", err)
	}
}

func TestUpdateConstituentLifecycleTransitions(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	cst := Constituent{
		ID:              uuid.New(),
		LifecycleStage:  LifecycleStageProspect,
		DuplicateStatus: DuplicateStatusActive,
	}

	next := LifecycleStageInquiry
	updated, err := bus.UpdateConstituent(context.Background(), cst, UpdateConstituent{LifecycleStage: &next})
	if err != nil {
		t.Fatalf("UpdateConstituent returned error: %v", err)
	}

	if updated.LifecycleStage != LifecycleStageInquiry {
		t.Fatalf("LifecycleStage = %s, want %s", updated.LifecycleStage, LifecycleStageInquiry)
	}

	invalid := LifecycleStageEnrolled
	_, err = bus.UpdateConstituent(context.Background(), cst, UpdateConstituent{LifecycleStage: &invalid})
	if !errors.Is(err, ErrInvalidLifecycleChange) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidLifecycleChange)
	}
}

func TestCreateDuplicateReviewValidatesPair(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	id := uuid.New()

	_, err := bus.CreateDuplicateReview(context.Background(), NewDuplicateReview{
		SourceConstituentID:    id,
		CandidateConstituentID: id,
		MatchType:              DuplicateReviewMatchTypeExact,
		MatchScore:             100,
		MatchReason:            "same email",
	})

	if !errors.Is(err, ErrInvalidDuplicateReview) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidDuplicateReview)
	}
}

func TestCreateDuplicateReviewValidatesMatchData(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nr   NewDuplicateReview
		want error
	}{
		{
			name: "match type",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchType("UNKNOWN"),
				MatchScore:             100,
				MatchReason:            "same email",
			},
			want: ErrInvalidMatchType,
		},
		{
			name: "match score",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchTypeExact,
				MatchScore:             101,
				MatchReason:            "same email",
			},
			want: ErrInvalidMatchScore,
		},
		{
			name: "match reason",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchTypeExact,
				MatchScore:             100,
			},
			want: ErrMatchReasonRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateDuplicateReview(context.Background(), tt.nr)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestResolveDuplicateReviewLinksSourceConstituent(t *testing.T) {
	t.Parallel()

	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	sourceID := uuid.New()
	candidateID := uuid.New()
	store.constituents[sourceID] = Constituent{ID: sourceID, DuplicateStatus: DuplicateStatusActive}
	store.constituents[candidateID] = Constituent{ID: candidateID, DuplicateStatus: DuplicateStatusActive}

	review := DuplicateReview{
		ID:                     uuid.New(),
		SourceConstituentID:    sourceID,
		CandidateConstituentID: candidateID,
		MatchType:              DuplicateReviewMatchTypeExact,
		MatchScore:             100,
		MatchReason:            "same email",
		Status:                 DuplicateReviewStatusPending,
	}

	resolved, err := bus.ResolveDuplicateReview(context.Background(), review, ResolveDuplicateReview{
		Resolution: DuplicateReviewResolutionLink,
		ActorID:    uuid.New(),
	})
	if err != nil {
		t.Fatalf("ResolveDuplicateReview returned error: %v", err)
	}

	if resolved.Status != DuplicateReviewStatusLinked {
		t.Fatalf("Status = %s, want %s", resolved.Status, DuplicateReviewStatusLinked)
	}

	source := store.constituents[sourceID]
	if source.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("DuplicateStatus = %s, want %s", source.DuplicateStatus, DuplicateStatusDuplicateOf)
	}

	if source.DuplicateOfID == nil || *source.DuplicateOfID != candidateID {
		t.Fatalf("DuplicateOfID = %v, want %s", source.DuplicateOfID, candidateID)
	}
}

func TestCreateApplicationPreventsDuplicateActiveApplication(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	store.applications = append(store.applications, Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
		Status:          ApplicationStatusDraft,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
	})

	if !errors.Is(err, ErrDuplicateApplication) {
		t.Fatalf("err = %v, want %v", err, ErrDuplicateApplication)
	}
}

func TestCreateApplicationAllowsClosedPriorApplication(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	store.applications = append(store.applications, Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
		Status:          ApplicationStatusWithdrawn,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	created, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeTransfer,
	})
	if err != nil {
		t.Fatalf("CreateApplication returned error: %v", err)
	}

	if created.Status != ApplicationStatusDraft {
		t.Fatalf("Status = %s, want %s", created.Status, ApplicationStatusDraft)
	}

	if len(store.applications) != 2 {
		t.Fatalf("application count = %d, want 2", len(store.applications))
	}
}

func TestCreateApplicationValidatesApplicationType(t *testing.T) {
	t.Parallel()

	store := newApplicationStubStore(uuid.New(), uuid.New(), uuid.New())
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   uuid.New(),
		ProgramID:       uuid.New(),
		AcademicTermID:  uuid.New(),
		ApplicationType: ApplicationType("UNKNOWN"),
	})

	if !errors.Is(err, ErrInvalidApplicationType) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidApplicationType)
	}
}

func newTestBusiness() ExtBusiness {
	return NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
}

func newApplicationStubStore(constituentID uuid.UUID, programID uuid.UUID, termID uuid.UUID) *stubStore {
	return &stubStore{
		constituents: map[uuid.UUID]Constituent{
			constituentID: {ID: constituentID, LifecycleStage: LifecycleStageApplicant, DuplicateStatus: DuplicateStatusActive},
		},
		programs: map[uuid.UUID]Program{
			programID: {ID: programID, Active: true},
		},
		terms: map[uuid.UUID]AcademicTerm{
			termID: {ID: termID, Active: true},
		},
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

type stubStore struct {
	constituents       map[uuid.UUID]Constituent
	constituentByEmail map[string]Constituent
	duplicateReviews   []DuplicateReview
	programs           map[uuid.UUID]Program
	terms              map[uuid.UUID]AcademicTerm
	applications       []Application
}

func (s *stubStore) NewWithTx(sqldb.CommitRollbacker) (Storer, error) {
	return s, nil
}

func (s *stubStore) Health(context.Context) (Health, error) {
	return Health{}, nil
}

func (s *stubStore) CreateConstituent(_ context.Context, cst Constituent) error {
	if s.constituents != nil {
		s.constituents[cst.ID] = cst
	}

	if s.constituentByEmail != nil {
		if _, exists := s.constituentByEmail[cst.PrimaryEmail.String()]; !exists {
			s.constituentByEmail[cst.PrimaryEmail.String()] = cst
		}
	}

	return nil
}

func (s *stubStore) UpdateConstituent(_ context.Context, cst Constituent) error {
	if s.constituents != nil {
		s.constituents[cst.ID] = cst
	}
	return nil
}

func (s *stubStore) QueryConstituents(context.Context, ConstituentQueryFilter, order.By, page.Page) ([]Constituent, error) {
	return nil, nil
}

func (s *stubStore) CountConstituents(context.Context, ConstituentQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryConstituentByID(_ context.Context, id uuid.UUID) (Constituent, error) {
	if s.constituents != nil {
		cst, exists := s.constituents[id]
		if exists {
			return cst, nil
		}
	}
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByPrimaryEmail(_ context.Context, email string) (Constituent, error) {
	if s.constituentByEmail != nil {
		cst, exists := s.constituentByEmail[email]
		if exists {
			return cst, nil
		}
	}

	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByExternalSISID(context.Context, string) (Constituent, error) {
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) UpsertProgram(context.Context, Program) error {
	return nil
}

func (s *stubStore) QueryPrograms(context.Context, ProgramQueryFilter, order.By, page.Page) ([]Program, error) {
	return nil, nil
}

func (s *stubStore) CountPrograms(context.Context, ProgramQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryProgramByID(_ context.Context, id uuid.UUID) (Program, error) {
	if s.programs != nil {
		program, exists := s.programs[id]
		if exists {
			return program, nil
		}
	}
	return Program{}, ErrProgramNotFound
}

func (s *stubStore) QueryProgramByExternalSISID(context.Context, string) (Program, error) {
	return Program{}, nil
}

func (s *stubStore) UpsertAcademicTerm(context.Context, AcademicTerm) error {
	return nil
}

func (s *stubStore) QueryAcademicTerms(context.Context, AcademicTermQueryFilter, order.By, page.Page) ([]AcademicTerm, error) {
	return nil, nil
}

func (s *stubStore) CountAcademicTerms(context.Context, AcademicTermQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryAcademicTermByID(_ context.Context, id uuid.UUID) (AcademicTerm, error) {
	if s.terms != nil {
		term, exists := s.terms[id]
		if exists {
			return term, nil
		}
	}
	return AcademicTerm{}, ErrAcademicTermNotFound
}

func (s *stubStore) QueryAcademicTermByExternalSISID(context.Context, string) (AcademicTerm, error) {
	return AcademicTerm{}, nil
}

func (s *stubStore) CreateDuplicateReview(_ context.Context, review DuplicateReview) error {
	s.duplicateReviews = append(s.duplicateReviews, review)
	return nil
}

func (s *stubStore) UpdateDuplicateReview(context.Context, DuplicateReview) error {
	return nil
}

func (s *stubStore) QueryDuplicateReviews(context.Context, DuplicateReviewQueryFilter, order.By, page.Page) ([]DuplicateReview, error) {
	return nil, nil
}

func (s *stubStore) CountDuplicateReviews(context.Context, DuplicateReviewQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryDuplicateReviewByID(context.Context, uuid.UUID) (DuplicateReview, error) {
	return DuplicateReview{}, nil
}

func (s *stubStore) CreateApplication(_ context.Context, app Application) error {
	s.applications = append(s.applications, app)
	return nil
}

func (s *stubStore) QueryApplications(context.Context, ApplicationQueryFilter, order.By, page.Page) ([]Application, error) {
	return s.applications, nil
}

func (s *stubStore) CountApplications(context.Context, ApplicationQueryFilter) (int, error) {
	return len(s.applications), nil
}

func (s *stubStore) QueryApplicationByID(_ context.Context, id uuid.UUID) (Application, error) {
	for _, app := range s.applications {
		if app.ID == id {
			return app, nil
		}
	}
	return Application{}, ErrApplicationNotFound
}

func (s *stubStore) QueryActiveApplicationByTuple(_ context.Context, constituentID uuid.UUID, academicTermID uuid.UUID, programID uuid.UUID) (Application, error) {
	for _, app := range s.applications {
		if app.ConstituentID == constituentID && app.AcademicTermID == academicTermID && app.ProgramID == programID && isApplicationActive(app.Status) {
			return app, nil
		}
	}
	return Application{}, ErrApplicationNotFound
}
