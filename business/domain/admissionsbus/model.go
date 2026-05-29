package admissionsbus

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
type Program struct{}

// AcademicTerm is SIS-owned reference data for application cycles.
type AcademicTerm struct{}

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
