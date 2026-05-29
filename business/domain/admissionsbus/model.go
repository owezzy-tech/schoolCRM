package admissionsbus

import (
	"time"

	"github.com/google/uuid"
)

// Health describes the currently available admissions bounded-context seams.
type Health struct {
	Context    string
	Status     string
	Aggregates []string
}

// Constituent is the durable person identity root for admissions workflows.
type Constituent struct{}

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
