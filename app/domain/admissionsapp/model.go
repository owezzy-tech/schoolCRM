package admissionsapp

import (
	"encoding/json"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// Health represents the admissions bounded-context scaffold status.
type Health struct {
	Context    string   `json:"context"`
	Status     string   `json:"status"`
	Aggregates []string `json:"aggregates"`
}

// Encode implements the encoder interface.
func (app Health) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppHealth(health admissionsbus.Health) Health {
	return Health{
		Context:    health.Context,
		Status:     health.Status,
		Aggregates: health.Aggregates,
	}
}

// Program represents SIS-owned program reference data.
type Program struct {
	ID            string  `json:"id"`
	ExternalSISID string  `json:"externalSISID"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Description   *string `json:"description,omitempty"`
	DegreeLevel   *string `json:"degreeLevel,omitempty"`
	Active        bool    `json:"active"`
	SyncedAt      *string `json:"syncedAt,omitempty"`
	DateCreated   string  `json:"dateCreated"`
	DateUpdated   string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Program) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppProgram(program admissionsbus.Program) Program {
	return Program{
		ID:            program.ID.String(),
		ExternalSISID: program.ExternalSISID,
		Name:          program.Name,
		Code:          program.Code,
		Description:   program.Description,
		DegreeLevel:   program.DegreeLevel,
		Active:        program.Active,
		SyncedAt:      formatTimePtr(program.SyncedAt),
		DateCreated:   program.DateCreated.Format(time.RFC3339),
		DateUpdated:   program.DateUpdated.Format(time.RFC3339),
	}
}

func toAppPrograms(programs []admissionsbus.Program) []Program {
	app := make([]Program, len(programs))
	for i, program := range programs {
		app[i] = toAppProgram(program)
	}

	return app
}

// AcademicTerm represents SIS-owned academic term reference data.
type AcademicTerm struct {
	ID                   string  `json:"id"`
	ExternalSISID        string  `json:"externalSISID"`
	Name                 string  `json:"name"`
	Code                 string  `json:"code"`
	TermType             *string `json:"termType,omitempty"`
	StartDate            string  `json:"startDate"`
	EndDate              string  `json:"endDate"`
	ApplicationStartDate *string `json:"applicationStartDate,omitempty"`
	ApplicationDeadline  *string `json:"applicationDeadline,omitempty"`
	Active               bool    `json:"active"`
	SyncedAt             *string `json:"syncedAt,omitempty"`
	DateCreated          string  `json:"dateCreated"`
	DateUpdated          string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app AcademicTerm) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppAcademicTerm(term admissionsbus.AcademicTerm) AcademicTerm {
	return AcademicTerm{
		ID:                   term.ID.String(),
		ExternalSISID:        term.ExternalSISID,
		Name:                 term.Name,
		Code:                 term.Code,
		TermType:             term.TermType,
		StartDate:            term.StartDate.Format(time.RFC3339),
		EndDate:              term.EndDate.Format(time.RFC3339),
		ApplicationStartDate: formatTimePtr(term.ApplicationStartDate),
		ApplicationDeadline:  formatTimePtr(term.ApplicationDeadline),
		Active:               term.Active,
		SyncedAt:             formatTimePtr(term.SyncedAt),
		DateCreated:          term.DateCreated.Format(time.RFC3339),
		DateUpdated:          term.DateUpdated.Format(time.RFC3339),
	}
}

func toAppAcademicTerms(terms []admissionsbus.AcademicTerm) []AcademicTerm {
	app := make([]AcademicTerm, len(terms))
	for i, term := range terms {
		app[i] = toAppAcademicTerm(term)
	}

	return app
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	formatted := t.Format(time.RFC3339)
	return &formatted
}
