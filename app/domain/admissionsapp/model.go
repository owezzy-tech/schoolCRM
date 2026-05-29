package admissionsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
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

// Constituent represents durable admissions identity data.
type Constituent struct {
	ID              string  `json:"id"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	PreferredName   *string `json:"preferredName,omitempty"`
	MiddleName      *string `json:"middleName,omitempty"`
	Suffix          *string `json:"suffix,omitempty"`
	DateOfBirth     string  `json:"dateOfBirth"`
	PrimaryEmail    string  `json:"primaryEmail"`
	PrimaryPhone    string  `json:"primaryPhone"`
	ExternalSISID   *string `json:"externalSISID,omitempty"`
	LifecycleStage  string  `json:"lifecycleStage"`
	DuplicateStatus string  `json:"duplicateStatus"`
	DuplicateOfID   *string `json:"duplicateOfID,omitempty"`
	SISSyncedAt     *string `json:"sisSyncedAt,omitempty"`
	DateCreated     string  `json:"dateCreated"`
	DateUpdated     string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Constituent) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppConstituent(cst admissionsbus.Constituent) Constituent {
	return Constituent{
		ID:              cst.ID.String(),
		FirstName:       cst.FirstName,
		LastName:        cst.LastName,
		PreferredName:   cst.PreferredName,
		MiddleName:      cst.MiddleName,
		Suffix:          cst.Suffix,
		DateOfBirth:     cst.DateOfBirth.Format(time.RFC3339),
		PrimaryEmail:    cst.PrimaryEmail.String(),
		PrimaryPhone:    cst.PrimaryPhone,
		ExternalSISID:   cst.ExternalSISID,
		LifecycleStage:  cst.LifecycleStage.String(),
		DuplicateStatus: cst.DuplicateStatus.String(),
		DuplicateOfID:   uuidStringPtr(cst.DuplicateOfID),
		SISSyncedAt:     formatTimePtr(cst.SISSyncedAt),
		DateCreated:     cst.DateCreated.Format(time.RFC3339),
		DateUpdated:     cst.DateUpdated.Format(time.RFC3339),
	}
}

func toAppConstituents(constituents []admissionsbus.Constituent) []Constituent {
	app := make([]Constituent, len(constituents))
	for i, cst := range constituents {
		app[i] = toAppConstituent(cst)
	}

	return app
}

// NewConstituent defines the data needed to add a new constituent.
type NewConstituent struct {
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	PreferredName  *string `json:"preferredName"`
	MiddleName     *string `json:"middleName"`
	Suffix         *string `json:"suffix"`
	DateOfBirth    string  `json:"dateOfBirth"`
	PrimaryEmail   string  `json:"primaryEmail"`
	PrimaryPhone   string  `json:"primaryPhone"`
	ExternalSISID  *string `json:"externalSISID"`
	LifecycleStage string  `json:"lifecycleStage"`
}

// Decode implements the decoder interface.
func (app *NewConstituent) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewConstituent(_ context.Context, app NewConstituent) (admissionsbus.NewConstituent, error) {
	var fieldErrors errs.FieldErrors

	dob, err := time.Parse(time.RFC3339, app.DateOfBirth)
	if err != nil {
		fieldErrors.Add("dateOfBirth", err)
	}

	email, err := mail.ParseAddress(app.PrimaryEmail)
	if err != nil {
		fieldErrors.Add("primaryEmail", err)
	}

	stage := admissionsbus.LifecycleStage(app.LifecycleStage)

	if len(fieldErrors) > 0 {
		return admissionsbus.NewConstituent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewConstituent{
		FirstName:       app.FirstName,
		LastName:        app.LastName,
		PreferredName:   app.PreferredName,
		MiddleName:      app.MiddleName,
		Suffix:          app.Suffix,
		DateOfBirth:     dob,
		PrimaryEmail:    *email,
		PrimaryPhone:    app.PrimaryPhone,
		ExternalSISID:   app.ExternalSISID,
		LifecycleStage:  stage,
		DuplicateStatus: admissionsbus.DuplicateStatusActive,
	}, nil
}

// UpdateConstituent defines the data needed to update a constituent.
type UpdateConstituent struct {
	PreferredName  *string `json:"preferredName"`
	MiddleName     *string `json:"middleName"`
	Suffix         *string `json:"suffix"`
	PrimaryEmail   *string `json:"primaryEmail"`
	PrimaryPhone   *string `json:"primaryPhone"`
	LifecycleStage *string `json:"lifecycleStage"`
}

// Decode implements the decoder interface.
func (app *UpdateConstituent) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusUpdateConstituent(app UpdateConstituent) (admissionsbus.UpdateConstituent, error) {
	var fieldErrors errs.FieldErrors
	var email *mail.Address

	if app.PrimaryEmail != nil {
		parsed, err := mail.ParseAddress(*app.PrimaryEmail)
		if err != nil {
			fieldErrors.Add("primaryEmail", err)
		} else {
			email = parsed
		}
	}

	var stage *admissionsbus.LifecycleStage
	if app.LifecycleStage != nil {
		parsed := admissionsbus.LifecycleStage(*app.LifecycleStage)
		stage = &parsed
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.UpdateConstituent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.UpdateConstituent{
		PreferredName:  app.PreferredName,
		MiddleName:     app.MiddleName,
		Suffix:         app.Suffix,
		PrimaryEmail:   email,
		PrimaryPhone:   app.PrimaryPhone,
		LifecycleStage: stage,
	}, nil
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

// DuplicateReview represents a potential constituent duplicate requiring staff resolution.
type DuplicateReview struct {
	ID                     string  `json:"id"`
	SourceConstituentID    string  `json:"sourceConstituentID"`
	CandidateConstituentID string  `json:"candidateConstituentID"`
	MatchType              string  `json:"matchType"`
	MatchScore             int     `json:"matchScore"`
	MatchReason            string  `json:"matchReason"`
	Status                 string  `json:"status"`
	ResolvedBy             *string `json:"resolvedBy,omitempty"`
	ResolvedAt             *string `json:"resolvedAt,omitempty"`
	ResolutionNote         *string `json:"resolutionNote,omitempty"`
	DateCreated            string  `json:"dateCreated"`
	DateUpdated            string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app DuplicateReview) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppDuplicateReview(review admissionsbus.DuplicateReview) DuplicateReview {
	return DuplicateReview{
		ID:                     review.ID.String(),
		SourceConstituentID:    review.SourceConstituentID.String(),
		CandidateConstituentID: review.CandidateConstituentID.String(),
		MatchType:              review.MatchType.String(),
		MatchScore:             review.MatchScore,
		MatchReason:            review.MatchReason,
		Status:                 review.Status.String(),
		ResolvedBy:             uuidStringPtr(review.ResolvedBy),
		ResolvedAt:             formatTimePtr(review.ResolvedAt),
		ResolutionNote:         review.ResolutionNote,
		DateCreated:            review.DateCreated.Format(time.RFC3339),
		DateUpdated:            review.DateUpdated.Format(time.RFC3339),
	}
}

func toAppDuplicateReviews(reviews []admissionsbus.DuplicateReview) []DuplicateReview {
	app := make([]DuplicateReview, len(reviews))
	for i, review := range reviews {
		app[i] = toAppDuplicateReview(review)
	}

	return app
}

// NewDuplicateReview defines the data needed to enqueue a duplicate review.
type NewDuplicateReview struct {
	SourceConstituentID    string `json:"sourceConstituentID"`
	CandidateConstituentID string `json:"candidateConstituentID"`
	MatchType              string `json:"matchType"`
	MatchScore             int    `json:"matchScore"`
	MatchReason            string `json:"matchReason"`
}

// Decode implements the decoder interface.
func (app *NewDuplicateReview) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewDuplicateReview(app NewDuplicateReview) (admissionsbus.NewDuplicateReview, error) {
	var fieldErrors errs.FieldErrors

	sourceID, err := uuid.Parse(app.SourceConstituentID)
	if err != nil {
		fieldErrors.Add("sourceConstituentID", err)
	}

	candidateID, err := uuid.Parse(app.CandidateConstituentID)
	if err != nil {
		fieldErrors.Add("candidateConstituentID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewDuplicateReview{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewDuplicateReview{
		SourceConstituentID:    sourceID,
		CandidateConstituentID: candidateID,
		MatchType:              admissionsbus.DuplicateReviewMatchType(app.MatchType),
		MatchScore:             app.MatchScore,
		MatchReason:            app.MatchReason,
	}, nil
}

// ResolveDuplicateReview defines a staff duplicate resolution action.
type ResolveDuplicateReview struct {
	Resolution string  `json:"resolution"`
	Note       *string `json:"note"`
}

// Decode implements the decoder interface.
func (app *ResolveDuplicateReview) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusResolveDuplicateReview(app ResolveDuplicateReview, actorID uuid.UUID) admissionsbus.ResolveDuplicateReview {
	return admissionsbus.ResolveDuplicateReview{
		Resolution: admissionsbus.DuplicateReviewResolution(app.Resolution),
		ActorID:    actorID,
		Note:       app.Note,
	}
}

// Application represents a constituent's application for a program and academic term.
type Application struct {
	ID                 string  `json:"id"`
	ConstituentID      string  `json:"constituentID"`
	ProgramID          string  `json:"programID"`
	AcademicTermID     string  `json:"academicTermID"`
	ApplicationType    string  `json:"applicationType"`
	Status             string  `json:"status"`
	AssignedReviewerID *string `json:"assignedReviewerID,omitempty"`
	SubmittedAt        *string `json:"submittedAt,omitempty"`
	DateCreated        string  `json:"dateCreated"`
	DateUpdated        string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Application) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppApplication(application admissionsbus.Application) Application {
	return Application{
		ID:                 application.ID.String(),
		ConstituentID:      application.ConstituentID.String(),
		ProgramID:          application.ProgramID.String(),
		AcademicTermID:     application.AcademicTermID.String(),
		ApplicationType:    application.ApplicationType.String(),
		Status:             application.Status.String(),
		AssignedReviewerID: uuidStringPtr(application.AssignedReviewerID),
		SubmittedAt:        formatTimePtr(application.SubmittedAt),
		DateCreated:        application.DateCreated.Format(time.RFC3339),
		DateUpdated:        application.DateUpdated.Format(time.RFC3339),
	}
}

func toAppApplications(applications []admissionsbus.Application) []Application {
	app := make([]Application, len(applications))
	for i, application := range applications {
		app[i] = toAppApplication(application)
	}

	return app
}

// NewApplication defines the data needed to create a draft application.
type NewApplication struct {
	ConstituentID      string  `json:"constituentID"`
	ProgramID          string  `json:"programID"`
	AcademicTermID     string  `json:"academicTermID"`
	ApplicationType    string  `json:"applicationType"`
	AssignedReviewerID *string `json:"assignedReviewerID"`
}

// Decode implements the decoder interface.
func (app *NewApplication) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewApplication(app NewApplication) (admissionsbus.NewApplication, error) {
	var fieldErrors errs.FieldErrors

	constituentID, err := uuid.Parse(app.ConstituentID)
	if err != nil {
		fieldErrors.Add("constituentID", err)
	}

	programID, err := uuid.Parse(app.ProgramID)
	if err != nil {
		fieldErrors.Add("programID", err)
	}

	academicTermID, err := uuid.Parse(app.AcademicTermID)
	if err != nil {
		fieldErrors.Add("academicTermID", err)
	}

	var assignedReviewerID *uuid.UUID
	if app.AssignedReviewerID != nil {
		parsed, err := uuid.Parse(*app.AssignedReviewerID)
		if err != nil {
			fieldErrors.Add("assignedReviewerID", err)
		} else {
			assignedReviewerID = &parsed
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewApplication{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewApplication{
		ConstituentID:      constituentID,
		ProgramID:          programID,
		AcademicTermID:     academicTermID,
		ApplicationType:    admissionsbus.ApplicationType(app.ApplicationType),
		AssignedReviewerID: assignedReviewerID,
	}, nil
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	formatted := t.Format(time.RFC3339)
	return &formatted
}

func uuidStringPtr(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	formatted := id.String()
	return &formatted
}
