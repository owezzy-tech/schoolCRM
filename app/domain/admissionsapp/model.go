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

// StaffProfile represents an admissions staff context profile.
type StaffProfile struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userID"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Active      bool     `json:"active"`
	DateCreated string   `json:"dateCreated"`
	DateUpdated string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app StaffProfile) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppStaffProfile(profile admissionsbus.StaffProfile) StaffProfile {
	permissions := admissionsbus.AdmissionsPermissionsForRoles(profile.Roles)

	return StaffProfile{
		ID:          profile.ID.String(),
		UserID:      profile.UserID.String(),
		Roles:       admissionsbus.AdmissionsRolesToStrings(profile.Roles),
		Permissions: admissionsbus.AdmissionsPermissionsToStrings(permissions),
		Active:      profile.Active,
		DateCreated: profile.DateCreated.Format(time.RFC3339),
		DateUpdated: profile.DateUpdated.Format(time.RFC3339),
	}
}

func toAppStaffProfiles(profiles []admissionsbus.StaffProfile) []StaffProfile {
	app := make([]StaffProfile, len(profiles))
	for i, profile := range profiles {
		app[i] = toAppStaffProfile(profile)
	}

	return app
}

// NewStaffProfile defines the data needed to create or update an admissions staff profile.
type NewStaffProfile struct {
	UserID string   `json:"userID"`
	Roles  []string `json:"roles"`
	Active bool     `json:"active"`
}

// Decode implements the decoder interface.
func (app *NewStaffProfile) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewStaffProfile(app NewStaffProfile) (admissionsbus.NewStaffProfile, error) {
	var fieldErrors errs.FieldErrors

	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		fieldErrors.Add("userID", err)
	}

	roles, err := admissionsbus.ParseAdmissionsRoles(app.Roles)
	if err != nil {
		fieldErrors.Add("roles", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewStaffProfile{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewStaffProfile{
		UserID: userID,
		Roles:  roles,
		Active: app.Active,
	}, nil
}

// ApplicantProfile represents an admissions applicant context profile.
type ApplicantProfile struct {
	ID            string `json:"id"`
	UserID        string `json:"userID"`
	ConstituentID string `json:"constituentID"`
	Active        bool   `json:"active"`
	DateCreated   string `json:"dateCreated"`
	DateUpdated   string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ApplicantProfile) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppApplicantProfile(profile admissionsbus.ApplicantProfile) ApplicantProfile {
	return ApplicantProfile{
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
		ConstituentID: profile.ConstituentID.String(),
		Active:        profile.Active,
		DateCreated:   profile.DateCreated.Format(time.RFC3339),
		DateUpdated:   profile.DateUpdated.Format(time.RFC3339),
	}
}

func toAppApplicantProfiles(profiles []admissionsbus.ApplicantProfile) []ApplicantProfile {
	app := make([]ApplicantProfile, len(profiles))
	for i, profile := range profiles {
		app[i] = toAppApplicantProfile(profile)
	}

	return app
}

// NewApplicantProfile defines the data needed to create or update an admissions applicant profile.
type NewApplicantProfile struct {
	UserID        string `json:"userID"`
	ConstituentID string `json:"constituentID"`
	Active        bool   `json:"active"`
}

// Decode implements the decoder interface.
func (app *NewApplicantProfile) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewApplicantProfile(app NewApplicantProfile) (admissionsbus.NewApplicantProfile, error) {
	var fieldErrors errs.FieldErrors

	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		fieldErrors.Add("userID", err)
	}

	constituentID, err := uuid.Parse(app.ConstituentID)
	if err != nil {
		fieldErrors.Add("constituentID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewApplicantProfile{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewApplicantProfile{
		UserID:        userID,
		ConstituentID: constituentID,
		Active:        app.Active,
	}, nil
}

// LeadScoreCriterion represents a single condition in a lead score rule.
type LeadScoreCriterion struct {
	Field    string   `json:"field"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

func toAppLeadScoreCriterion(criterion admissionsbus.LeadScoreCriterion) LeadScoreCriterion {
	return LeadScoreCriterion{
		Field:    criterion.Field.String(),
		Operator: criterion.Operator.String(),
		Values:   criterion.Values,
	}
}

func toAppLeadScoreCriteria(criteria []admissionsbus.LeadScoreCriterion) []LeadScoreCriterion {
	app := make([]LeadScoreCriterion, len(criteria))
	for i, criterion := range criteria {
		app[i] = toAppLeadScoreCriterion(criterion)
	}

	return app
}

func toBusLeadScoreCriterion(app LeadScoreCriterion) admissionsbus.LeadScoreCriterion {
	return admissionsbus.LeadScoreCriterion{
		Field:    admissionsbus.LeadScoreCriterionField(app.Field),
		Operator: admissionsbus.LeadScoreCriterionOperator(app.Operator),
		Values:   app.Values,
	}
}

func toBusLeadScoreCriteria(app []LeadScoreCriterion) []admissionsbus.LeadScoreCriterion {
	criteria := make([]admissionsbus.LeadScoreCriterion, len(app))
	for i, criterion := range app {
		criteria[i] = toBusLeadScoreCriterion(criterion)
	}

	return criteria
}

// LeadScoreRule represents an explainable rule contributing points to a lead score.
type LeadScoreRule struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description,omitempty"`
	Criteria    []LeadScoreCriterion `json:"criteria"`
	Points      int                  `json:"points"`
	Active      bool                 `json:"active"`
	Priority    int                  `json:"priority"`
	DateCreated string               `json:"dateCreated"`
	DateUpdated string               `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app LeadScoreRule) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppLeadScoreRule(rule admissionsbus.LeadScoreRule) LeadScoreRule {
	return LeadScoreRule{
		ID:          rule.ID.String(),
		Name:        rule.Name,
		Description: rule.Description,
		Criteria:    toAppLeadScoreCriteria(rule.Criteria),
		Points:      rule.Points,
		Active:      rule.Active,
		Priority:    rule.Priority,
		DateCreated: rule.DateCreated.Format(time.RFC3339),
		DateUpdated: rule.DateUpdated.Format(time.RFC3339),
	}
}

func toAppLeadScoreRules(rules []admissionsbus.LeadScoreRule) []LeadScoreRule {
	app := make([]LeadScoreRule, len(rules))
	for i, rule := range rules {
		app[i] = toAppLeadScoreRule(rule)
	}

	return app
}

// NewLeadScoreRule defines the data needed to create or update a lead score rule.
type NewLeadScoreRule struct {
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Criteria    []LeadScoreCriterion `json:"criteria"`
	Points      int                  `json:"points"`
	Active      bool                 `json:"active"`
	Priority    int                  `json:"priority"`
}

// Decode implements the decoder interface.
func (app *NewLeadScoreRule) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewLeadScoreRule(app NewLeadScoreRule) admissionsbus.NewLeadScoreRule {
	return admissionsbus.NewLeadScoreRule{
		Name:        app.Name,
		Description: app.Description,
		Criteria:    toBusLeadScoreCriteria(app.Criteria),
		Points:      app.Points,
		Active:      app.Active,
		Priority:    app.Priority,
	}
}

// LeadScoreRuleResult explains how one rule contributed to a score.
type LeadScoreRuleResult struct {
	RuleID  string `json:"ruleID"`
	Name    string `json:"name"`
	Points  int    `json:"points"`
	Matched bool   `json:"matched"`
	Reason  string `json:"reason"`
}

func toAppLeadScoreRuleResult(result admissionsbus.LeadScoreRuleResult) LeadScoreRuleResult {
	return LeadScoreRuleResult{
		RuleID:  result.RuleID.String(),
		Name:    result.Name,
		Points:  result.Points,
		Matched: result.Matched,
		Reason:  result.Reason,
	}
}

func toAppLeadScoreRuleResults(results []admissionsbus.LeadScoreRuleResult) []LeadScoreRuleResult {
	app := make([]LeadScoreRuleResult, len(results))
	for i, result := range results {
		app[i] = toAppLeadScoreRuleResult(result)
	}

	return app
}

// LeadScore represents the latest explainable score for a constituent.
type LeadScore struct {
	ID             string                `json:"id"`
	ConstituentID  string                `json:"constituentID"`
	TotalScore     int                   `json:"totalScore"`
	Band           string                `json:"band"`
	Breakdown      []LeadScoreRuleResult `json:"breakdown"`
	RecalculatedAt string                `json:"recalculatedAt"`
	DateCreated    string                `json:"dateCreated"`
	DateUpdated    string                `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app LeadScore) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppLeadScore(score admissionsbus.LeadScore) LeadScore {
	return LeadScore{
		ID:             score.ID.String(),
		ConstituentID:  score.ConstituentID.String(),
		TotalScore:     score.TotalScore,
		Band:           score.Band.String(),
		Breakdown:      toAppLeadScoreRuleResults(score.Breakdown),
		RecalculatedAt: score.RecalculatedAt.Format(time.RFC3339),
		DateCreated:    score.DateCreated.Format(time.RFC3339),
		DateUpdated:    score.DateUpdated.Format(time.RFC3339),
	}
}

func toAppLeadScores(scores []admissionsbus.LeadScore) []LeadScore {
	app := make([]LeadScore, len(scores))
	for i, score := range scores {
		app[i] = toAppLeadScore(score)
	}

	return app
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

// Inquiry represents a public admissions inquiry submission.
type Inquiry struct {
	ID                string  `json:"id"`
	ConstituentID     string  `json:"constituentID"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	DateOfBirth       string  `json:"dateOfBirth"`
	PrimaryEmail      string  `json:"primaryEmail"`
	PrimaryPhone      string  `json:"primaryPhone"`
	ProgramOfInterest *string `json:"programOfInterest,omitempty"`
	TermOfInterest    *string `json:"termOfInterest,omitempty"`
	Source            string  `json:"source"`
	UTMSource         *string `json:"utmSource,omitempty"`
	UTMMedium         *string `json:"utmMedium,omitempty"`
	UTMCampaign       *string `json:"utmCampaign,omitempty"`
	Message           *string `json:"message,omitempty"`
	Status            string  `json:"status"`
	DateCreated       string  `json:"dateCreated"`
	DateUpdated       string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Inquiry) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppInquiry(inquiry admissionsbus.Inquiry) Inquiry {
	return Inquiry{
		ID:                inquiry.ID.String(),
		ConstituentID:     inquiry.ConstituentID.String(),
		FirstName:         inquiry.FirstName,
		LastName:          inquiry.LastName,
		DateOfBirth:       inquiry.DateOfBirth.Format(time.RFC3339),
		PrimaryEmail:      inquiry.PrimaryEmail.String(),
		PrimaryPhone:      inquiry.PrimaryPhone,
		ProgramOfInterest: uuidStringPtr(inquiry.ProgramOfInterest),
		TermOfInterest:    uuidStringPtr(inquiry.TermOfInterest),
		Source:            inquiry.Source,
		UTMSource:         inquiry.UTMSource,
		UTMMedium:         inquiry.UTMMedium,
		UTMCampaign:       inquiry.UTMCampaign,
		Message:           inquiry.Message,
		Status:            inquiry.Status.String(),
		DateCreated:       inquiry.DateCreated.Format(time.RFC3339),
		DateUpdated:       inquiry.DateUpdated.Format(time.RFC3339),
	}
}

func toAppInquiries(inquiries []admissionsbus.Inquiry) []Inquiry {
	app := make([]Inquiry, len(inquiries))
	for i, inquiry := range inquiries {
		app[i] = toAppInquiry(inquiry)
	}

	return app
}

// NewInquiry defines the data needed to submit a public inquiry form.
type NewInquiry struct {
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	DateOfBirth       string  `json:"dateOfBirth"`
	PrimaryEmail      string  `json:"primaryEmail"`
	PrimaryPhone      string  `json:"primaryPhone"`
	ProgramOfInterest *string `json:"programOfInterest"`
	TermOfInterest    *string `json:"termOfInterest"`
	Source            string  `json:"source"`
	UTMSource         *string `json:"utmSource"`
	UTMMedium         *string `json:"utmMedium"`
	UTMCampaign       *string `json:"utmCampaign"`
	Message           *string `json:"message"`
}

// Decode implements the decoder interface.
func (app *NewInquiry) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewInquiry(app NewInquiry) (admissionsbus.NewInquiry, error) {
	var fieldErrors errs.FieldErrors

	dob, err := time.Parse(time.RFC3339, app.DateOfBirth)
	if err != nil {
		fieldErrors.Add("dateOfBirth", err)
	}

	email, err := mail.ParseAddress(app.PrimaryEmail)
	if err != nil {
		fieldErrors.Add("primaryEmail", err)
	}

	programID, err := parseUUIDPtr(app.ProgramOfInterest)
	if err != nil {
		fieldErrors.Add("programOfInterest", err)
	}

	termID, err := parseUUIDPtr(app.TermOfInterest)
	if err != nil {
		fieldErrors.Add("termOfInterest", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewInquiry{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewInquiry{
		FirstName:         app.FirstName,
		LastName:          app.LastName,
		DateOfBirth:       dob,
		PrimaryEmail:      *email,
		PrimaryPhone:      app.PrimaryPhone,
		ProgramOfInterest: programID,
		TermOfInterest:    termID,
		Source:            app.Source,
		UTMSource:         app.UTMSource,
		UTMMedium:         app.UTMMedium,
		UTMCampaign:       app.UTMCampaign,
		Message:           app.Message,
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

// ApplicationFormField represents a configurable, non-core application form field.
type ApplicationFormField struct {
	FieldName    string  `json:"fieldName"`
	FieldType    string  `json:"fieldType"`
	Required     bool    `json:"required"`
	DisplayOrder int     `json:"displayOrder"`
	Validation   *string `json:"validation,omitempty"`
}

// ApplicationChecklistTemplateItem represents a checklist/document requirement in a form template.
type ApplicationChecklistTemplateItem struct {
	ItemKey      string  `json:"itemKey"`
	DocumentName string  `json:"documentName"`
	Description  *string `json:"description,omitempty"`
	Required     bool    `json:"required"`
	DisplayOrder int     `json:"displayOrder"`
}

// ApplicationFormTemplate represents configurable application form requirements.
type ApplicationFormTemplate struct {
	ID              string                             `json:"id"`
	ProgramID       string                             `json:"programID"`
	AcademicTermID  string                             `json:"academicTermID"`
	ApplicationType string                             `json:"applicationType"`
	Name            string                             `json:"name"`
	Description     *string                            `json:"description,omitempty"`
	Version         int                                `json:"version"`
	RequiredFields  []ApplicationFormField             `json:"requiredFields"`
	ChecklistItems  []ApplicationChecklistTemplateItem `json:"checklistItems"`
	Active          bool                               `json:"active"`
	Priority        int                                `json:"priority"`
	DateCreated     string                             `json:"dateCreated"`
	DateUpdated     string                             `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ApplicationFormTemplate) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppApplicationFormField(field admissionsbus.ApplicationFormField) ApplicationFormField {
	return ApplicationFormField{
		FieldName:    field.FieldName,
		FieldType:    field.FieldType,
		Required:     field.Required,
		DisplayOrder: field.DisplayOrder,
		Validation:   field.Validation,
	}
}

func toAppApplicationFormFields(fields []admissionsbus.ApplicationFormField) []ApplicationFormField {
	app := make([]ApplicationFormField, len(fields))
	for i, field := range fields {
		app[i] = toAppApplicationFormField(field)
	}

	return app
}

func toAppChecklistTemplateItem(item admissionsbus.ApplicationChecklistTemplateItem) ApplicationChecklistTemplateItem {
	return ApplicationChecklistTemplateItem{
		ItemKey:      item.ItemKey,
		DocumentName: item.DocumentName,
		Description:  item.Description,
		Required:     item.Required,
		DisplayOrder: item.DisplayOrder,
	}
}

func toAppChecklistTemplateItems(items []admissionsbus.ApplicationChecklistTemplateItem) []ApplicationChecklistTemplateItem {
	app := make([]ApplicationChecklistTemplateItem, len(items))
	for i, item := range items {
		app[i] = toAppChecklistTemplateItem(item)
	}

	return app
}

func toAppApplicationFormTemplate(template admissionsbus.ApplicationFormTemplate) ApplicationFormTemplate {
	return ApplicationFormTemplate{
		ID:              template.ID.String(),
		ProgramID:       template.ProgramID.String(),
		AcademicTermID:  template.AcademicTermID.String(),
		ApplicationType: template.ApplicationType.String(),
		Name:            template.Name,
		Description:     template.Description,
		Version:         template.Version,
		RequiredFields:  toAppApplicationFormFields(template.RequiredFields),
		ChecklistItems:  toAppChecklistTemplateItems(template.ChecklistItems),
		Active:          template.Active,
		Priority:        template.Priority,
		DateCreated:     template.DateCreated.Format(time.RFC3339),
		DateUpdated:     template.DateUpdated.Format(time.RFC3339),
	}
}

func toAppApplicationFormTemplates(templates []admissionsbus.ApplicationFormTemplate) []ApplicationFormTemplate {
	app := make([]ApplicationFormTemplate, len(templates))
	for i, template := range templates {
		app[i] = toAppApplicationFormTemplate(template)
	}

	return app
}

// ApplicationTransition represents immutable application status transition history.
type ApplicationTransition struct {
	ID            string          `json:"id"`
	ApplicationID string          `json:"applicationID"`
	FromStatus    string          `json:"fromStatus"`
	ToStatus      string          `json:"toStatus"`
	ActorID       string          `json:"actorID"`
	Reason        *string         `json:"reason,omitempty"`
	Note          *string         `json:"note,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	DateCreated   string          `json:"dateCreated"`
}

// Encode implements the encoder interface.
func (app ApplicationTransition) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppApplicationTransition(transition admissionsbus.ApplicationTransition) ApplicationTransition {
	return ApplicationTransition{
		ID:            transition.ID.String(),
		ApplicationID: transition.ApplicationID.String(),
		FromStatus:    transition.FromStatus.String(),
		ToStatus:      transition.ToStatus.String(),
		ActorID:       transition.ActorID.String(),
		Reason:        transition.Reason,
		Note:          transition.Note,
		Metadata:      json.RawMessage(transition.Metadata),
		DateCreated:   transition.DateCreated.Format(time.RFC3339),
	}
}

func toAppApplicationTransitions(transitions []admissionsbus.ApplicationTransition) []ApplicationTransition {
	app := make([]ApplicationTransition, len(transitions))
	for i, transition := range transitions {
		app[i] = toAppApplicationTransition(transition)
	}

	return app
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

// NewApplicationFormTemplate defines the data needed to create or update a form template.
type NewApplicationFormTemplate struct {
	ProgramID       string                             `json:"programID"`
	AcademicTermID  string                             `json:"academicTermID"`
	ApplicationType string                             `json:"applicationType"`
	Name            string                             `json:"name"`
	Description     *string                            `json:"description"`
	RequiredFields  []ApplicationFormField             `json:"requiredFields"`
	ChecklistItems  []ApplicationChecklistTemplateItem `json:"checklistItems"`
	Active          bool                               `json:"active"`
	Priority        int                                `json:"priority"`
}

// Decode implements the decoder interface.
func (app *NewApplicationFormTemplate) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusApplicationFormField(app ApplicationFormField) admissionsbus.ApplicationFormField {
	return admissionsbus.ApplicationFormField{
		FieldName:    app.FieldName,
		FieldType:    app.FieldType,
		Required:     app.Required,
		DisplayOrder: app.DisplayOrder,
		Validation:   app.Validation,
	}
}

func toBusApplicationFormFields(app []ApplicationFormField) []admissionsbus.ApplicationFormField {
	fields := make([]admissionsbus.ApplicationFormField, len(app))
	for i, field := range app {
		fields[i] = toBusApplicationFormField(field)
	}

	return fields
}

func toBusChecklistTemplateItem(app ApplicationChecklistTemplateItem) admissionsbus.ApplicationChecklistTemplateItem {
	return admissionsbus.ApplicationChecklistTemplateItem{
		ItemKey:      app.ItemKey,
		DocumentName: app.DocumentName,
		Description:  app.Description,
		Required:     app.Required,
		DisplayOrder: app.DisplayOrder,
	}
}

func toBusChecklistTemplateItems(app []ApplicationChecklistTemplateItem) []admissionsbus.ApplicationChecklistTemplateItem {
	items := make([]admissionsbus.ApplicationChecklistTemplateItem, len(app))
	for i, item := range app {
		items[i] = toBusChecklistTemplateItem(item)
	}

	return items
}

func toBusNewApplicationFormTemplate(app NewApplicationFormTemplate) (admissionsbus.NewApplicationFormTemplate, error) {
	var fieldErrors errs.FieldErrors

	programID, err := uuid.Parse(app.ProgramID)
	if err != nil {
		fieldErrors.Add("programID", err)
	}

	academicTermID, err := uuid.Parse(app.AcademicTermID)
	if err != nil {
		fieldErrors.Add("academicTermID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewApplicationFormTemplate{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewApplicationFormTemplate{
		ProgramID:       programID,
		AcademicTermID:  academicTermID,
		ApplicationType: admissionsbus.ApplicationType(app.ApplicationType),
		Name:            app.Name,
		Description:     app.Description,
		RequiredFields:  toBusApplicationFormFields(app.RequiredFields),
		ChecklistItems:  toBusChecklistTemplateItems(app.ChecklistItems),
		Active:          app.Active,
		Priority:        app.Priority,
	}, nil
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

// NewApplicationTransition defines the data needed to change an application status.
type NewApplicationTransition struct {
	ToStatus string          `json:"toStatus"`
	Reason   *string         `json:"reason"`
	Note     *string         `json:"note"`
	Metadata json.RawMessage `json:"metadata"`
}

// Decode implements the decoder interface.
func (app *NewApplicationTransition) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewApplicationTransition(app NewApplicationTransition, actorID uuid.UUID) admissionsbus.NewApplicationTransition {
	return admissionsbus.NewApplicationTransition{
		ToStatus: admissionsbus.ApplicationStatus(app.ToStatus),
		ActorID:  actorID,
		Reason:   app.Reason,
		Note:     app.Note,
		Metadata: []byte(app.Metadata),
	}
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

func parseUUIDPtr(value *string) (*uuid.UUID, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	id, err := uuid.Parse(*value)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
