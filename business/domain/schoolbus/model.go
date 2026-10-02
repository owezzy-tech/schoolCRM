// Package schoolbus owns school-scoped membership and administration rules.
package schoolbus

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrForbidden = errors.New("school action is not permitted")
	ErrNotFound  = errors.New("school resource not found")
	ErrInvalid   = errors.New("invalid school request")
)

// Capability grants authority within a school or one of its departments.
type Capability string

const (
	ManageMembers  Capability = "manage_members"
	Teach          Capability = "teach"
	ReviewLessons  Capability = "review_lessons"
	ApproveLessons Capability = "approve_lessons"
)

// Scope identifies a school and optionally a department within it.
type Scope struct {
	SchoolID     uuid.UUID  `json:"schoolID"`
	DepartmentID *uuid.UUID `json:"departmentID,omitempty"`
}

type School struct {
	ID          uuid.UUID `db:"school_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	DateCreated time.Time `db:"date_created" json:"dateCreated"`
}

type Department struct {
	ID          uuid.UUID `db:"department_id" json:"id"`
	SchoolID    uuid.UUID `db:"school_id" json:"schoolID"`
	Name        string    `db:"name" json:"name"`
	DateCreated time.Time `db:"date_created" json:"dateCreated"`
}

type Membership struct {
	ID           uuid.UUID  `db:"membership_id" json:"id"`
	SchoolID     uuid.UUID  `db:"school_id" json:"schoolID"`
	DepartmentID *uuid.UUID `db:"department_id" json:"departmentID,omitempty"`
	UserID       uuid.UUID  `db:"user_id" json:"userID"`
	Capability   Capability `db:"capability" json:"capability"`
	Active       bool       `db:"active" json:"active"`
	DateUpdated  time.Time  `db:"date_updated" json:"dateUpdated"`
}

// Actor is loaded from the current user row, never cached token roles.
type Actor struct {
	Enabled     bool `db:"enabled"`
	SuperAdmin  bool `db:"super_admin"`
	SchoolAdmin bool `db:"school_admin"`
}
