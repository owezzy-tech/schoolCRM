package admissionsbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// Health describes the currently available admissions bounded-context seams.
type Health struct {
	Context    string
	Status     string
	Aggregates []string
}

// LifecycleStage represents a constituent's overall admissions journey.
type LifecycleStage string

// Set of valid constituent lifecycle stages.
const (
	LifecycleStageProspect  LifecycleStage = "PROSPECT"
	LifecycleStageInquiry   LifecycleStage = "INQUIRY"
	LifecycleStageApplicant LifecycleStage = "APPLICANT"
	LifecycleStageAdmitted  LifecycleStage = "ADMITTED"
	LifecycleStageEnrolled  LifecycleStage = "ENROLLED"
	LifecycleStageAlumni    LifecycleStage = "ALUMNI"
)

// String returns the lifecycle stage as a string.
func (stage LifecycleStage) String() string {
	return string(stage)
}

// DuplicateStatus represents whether a constituent is canonical or duplicate-linked.
type DuplicateStatus string

// Set of valid constituent duplicate statuses.
const (
	DuplicateStatusActive      DuplicateStatus = "ACTIVE"
	DuplicateStatusMerged      DuplicateStatus = "MERGED"
	DuplicateStatusDuplicateOf DuplicateStatus = "DUPLICATE_OF"
)

// String returns the duplicate status as a string.
func (status DuplicateStatus) String() string {
	return string(status)
}

// Constituent is the durable person identity root for admissions workflows.
type Constituent struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	DateOfBirth     time.Time
	PrimaryEmail    mail.Address
	PrimaryPhone    string
	ExternalSISID   *string
	LifecycleStage  LifecycleStage
	DuplicateStatus DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
	DateCreated     time.Time
	DateUpdated     time.Time
}

// NewConstituent is what we require from clients when adding a Constituent.
type NewConstituent struct {
	FirstName       string
	LastName        string
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	DateOfBirth     time.Time
	PrimaryEmail    mail.Address
	PrimaryPhone    string
	ExternalSISID   *string
	LifecycleStage  LifecycleStage
	DuplicateStatus DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
}

// UpdateConstituent defines what information may be provided to modify a Constituent.
type UpdateConstituent struct {
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	PrimaryEmail    *mail.Address
	PrimaryPhone    *string
	LifecycleStage  *LifecycleStage
	DuplicateStatus *DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
}

// Inquiry captures pre-application interest in the school.
type Inquiry struct{}

// Application represents a constituent's program application for a term.
type Application struct{}

// Checklist groups required admissions items for an application.
type Checklist struct{}

// Document represents applicant-submitted evidence for checklist items.
type Document struct{}

// Decision represents the outcome of an application review.
type Decision struct{}

// Program is SIS-owned reference data used by admissions applications.
type Program struct {
	ID            uuid.UUID
	ExternalSISID string
	Name          string
	Code          string
	Description   *string
	DegreeLevel   *string
	Active        bool
	SyncedAt      *time.Time
	DateCreated   time.Time
	DateUpdated   time.Time
}

// UpsertProgram is the sync/import-owned data needed to create or update a Program.
type UpsertProgram struct {
	ID            *uuid.UUID
	ExternalSISID string
	Name          string
	Code          string
	Description   *string
	DegreeLevel   *string
	Active        bool
	SyncedAt      *time.Time
}

// AcademicTerm is SIS-owned reference data for application cycles.
type AcademicTerm struct {
	ID                   uuid.UUID
	ExternalSISID        string
	Name                 string
	Code                 string
	TermType             *string
	StartDate            time.Time
	EndDate              time.Time
	ApplicationStartDate *time.Time
	ApplicationDeadline  *time.Time
	Active               bool
	SyncedAt             *time.Time
	DateCreated          time.Time
	DateUpdated          time.Time
}

// UpsertAcademicTerm is the sync/import-owned data needed to create or update an AcademicTerm.
type UpsertAcademicTerm struct {
	ID                   *uuid.UUID
	ExternalSISID        string
	Name                 string
	Code                 string
	TermType             *string
	StartDate            time.Time
	EndDate              time.Time
	ApplicationStartDate *time.Time
	ApplicationDeadline  *time.Time
	Active               bool
	SyncedAt             *time.Time
}

// DuplicateReview represents a potential constituent duplicate requiring resolution.
type DuplicateReview struct{}

// AggregateNames returns the scaffolded admissions aggregate names.
func AggregateNames() []string {
	return []string{
		"constituent",
		"inquiry",
		"application",
		"checklist",
		"document",
		"decision",
		"program",
		"academicTerm",
		"duplicateReview",
	}
}

// ProgramQueryFilter holds the available fields a program query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ProgramQueryFilter struct {
	ID            *uuid.UUID
	ExternalSISID *string
	Code          *string
	Active        *bool
}

// AcademicTermQueryFilter holds the available fields an academic term query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type AcademicTermQueryFilter struct {
	ID            *uuid.UUID
	ExternalSISID *string
	Code          *string
	Active        *bool
}

// ConstituentQueryFilter holds the available fields a constituent query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ConstituentQueryFilter struct {
	ID              *uuid.UUID
	PrimaryEmail    *mail.Address
	ExternalSISID   *string
	LifecycleStage  *LifecycleStage
	DuplicateStatus *DuplicateStatus
}
