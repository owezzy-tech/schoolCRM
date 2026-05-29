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

func newTestBusiness() ExtBusiness {
	return NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

type stubStore struct{}

func (s *stubStore) NewWithTx(sqldb.CommitRollbacker) (Storer, error) {
	return s, nil
}

func (s *stubStore) Health(context.Context) (Health, error) {
	return Health{}, nil
}

func (s *stubStore) CreateConstituent(context.Context, Constituent) error {
	return nil
}

func (s *stubStore) UpdateConstituent(context.Context, Constituent) error {
	return nil
}

func (s *stubStore) QueryConstituents(context.Context, ConstituentQueryFilter, order.By, page.Page) ([]Constituent, error) {
	return nil, nil
}

func (s *stubStore) CountConstituents(context.Context, ConstituentQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryConstituentByID(context.Context, uuid.UUID) (Constituent, error) {
	return Constituent{}, nil
}

func (s *stubStore) QueryConstituentByPrimaryEmail(context.Context, string) (Constituent, error) {
	return Constituent{}, nil
}

func (s *stubStore) QueryConstituentByExternalSISID(context.Context, string) (Constituent, error) {
	return Constituent{}, nil
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

func (s *stubStore) QueryProgramByID(context.Context, uuid.UUID) (Program, error) {
	return Program{}, nil
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

func (s *stubStore) QueryAcademicTermByID(context.Context, uuid.UUID) (AcademicTerm, error) {
	return AcademicTerm{}, nil
}

func (s *stubStore) QueryAcademicTermByExternalSISID(context.Context, string) (AcademicTerm, error) {
	return AcademicTerm{}, nil
}
