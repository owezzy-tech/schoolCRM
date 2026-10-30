// Package lessonbus owns versioned lesson plans and their publication workflow.
package lessonbus

import (
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
)

var (
	ErrForbidden = errors.New("lesson action is not permitted")
	ErrNotFound  = errors.New("lesson resource not found")
	ErrInvalid   = errors.New("invalid lesson request")
	ErrConflict  = errors.New("lesson version is stale or not in the required state")
)

// Status is the workflow state of one immutable lesson version.
type Status string

const (
	Draft            Status = "draft"
	HODReview        Status = "hod_review"
	DeanApproval     Status = "dean_approval"
	ChangesRequested Status = "changes_requested"
	Approved         Status = "approved"
	Published        Status = "published"
)

// Decision is an HOD or dean outcome for one version.
type Decision string

const (
	Accept         Decision = "approve"
	RequestChanges Decision = "request_changes"
)

// Plan points at its latest version and, separately, its live publication.
type Plan struct {
	ID               uuid.UUID `db:"plan_id" json:"id"`
	SchoolID         uuid.UUID `db:"school_id" json:"schoolID"`
	DepartmentID     uuid.UUID `db:"department_id" json:"departmentID"`
	AuthorID         uuid.UUID `db:"author_id" json:"authorID"`
	Title            string    `db:"title" json:"title"`
	Status           Status    `db:"status" json:"status"`
	CurrentVersion   int       `db:"current_version" json:"currentVersion"`
	PublishedVersion *int      `db:"published_version" json:"publishedVersion"`
	DateCreated      time.Time `db:"date_created" json:"dateCreated"`
	DateUpdated      time.Time `db:"date_updated" json:"dateUpdated"`
}

// Version is an immutable lesson snapshot plus the decisions bound to it.
type Version struct {
	PlanID           uuid.UUID       `db:"plan_id" json:"planID"`
	Number           int             `db:"version_number" json:"version"`
	Title            string          `db:"title" json:"title"`
	Content          json.RawMessage `db:"content" json:"content,omitempty"`
	ChangeSummary    string          `db:"change_summary" json:"changeSummary"`
	AuthorID         uuid.UUID       `db:"author_id" json:"authorID"`
	Status           Status          `db:"status" json:"status"`
	ReviewerID       *uuid.UUID      `db:"reviewer_id" json:"reviewerID"`
	ReviewedAt       *time.Time      `db:"reviewed_at" json:"reviewedAt"`
	ApproverID       *uuid.UUID      `db:"approver_id" json:"approverID"`
	ApprovedAt       *time.Time      `db:"approved_at" json:"approvedAt"`
	PublishedAt      *time.Time      `db:"published_at" json:"publishedAt"`
	ReviewFeedback   string          `db:"review_feedback" json:"reviewFeedback"`
	ApprovalFeedback string          `db:"approval_feedback" json:"approvalFeedback"`
	DateCreated      time.Time       `db:"date_created" json:"dateCreated"`
}

// Revision is the author-supplied content of a new version.
type Revision struct {
	Title         string
	Content       json.RawMessage
	ChangeSummary string
	RequestID     uuid.UUID
}

// Creation retains the original response for one actor-owned creation command.
type Creation struct {
	ActorID   uuid.UUID
	RequestID uuid.UUID
	InputHash string
	Result    Plan
}

// Authority is the actor's current standing in one department, read inside the
// command's transaction rather than from token claims.
type Authority struct {
	Enabled      bool
	Capabilities []schoolbus.Capability
}

func (a Authority) has(capabilities ...schoolbus.Capability) bool {
	return a.Enabled && slices.ContainsFunc(capabilities, func(c schoolbus.Capability) bool {
		return slices.Contains(a.Capabilities, c)
	})
}
