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

// EventRegistration represents one admissions event registration for app responses.
type EventRegistration struct {
	ID              string  `json:"id"`
	ConstituentID   *string `json:"constituentId,omitempty"`
	ConstituentName string  `json:"constituentName"`
	Email           string  `json:"email"`
	Phone           *string `json:"phone,omitempty"`
	Status          string  `json:"status"`
	RegisteredAt    string  `json:"registeredAt"`
	MatchStatus     string  `json:"matchStatus"`
	Source          string  `json:"source"`
	CheckedInAt     *string `json:"checkedInAt,omitempty"`
	CheckedInByID   *string `json:"checkedInById,omitempty"`
	FirstName       string  `json:"-"`
	LastName        string  `json:"-"`
}

// Encode implements the encoder interface.
func (app EventRegistration) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// Event represents an admissions engagement event for app responses.
type Event struct {
	ID                      string              `json:"id"`
	Title                   string              `json:"title"`
	Type                    string              `json:"type"`
	Status                  string              `json:"status"`
	Description             string              `json:"description"`
	Start                   string              `json:"start"`
	End                     string              `json:"end"`
	Location                string              `json:"location"`
	IsVirtual               bool                `json:"isVirtual"`
	Capacity                int                 `json:"capacity"`
	RegisteredCount         int                 `json:"registeredCount"`
	CheckedInCount          int                 `json:"checkedInCount"`
	RegistrationDeadline    *string             `json:"registrationDeadline,omitempty"`
	AutoConfirmationEnabled bool                `json:"autoConfirmationEnabled"`
	AutoReminderEnabled     bool                `json:"autoReminderEnabled"`
	Registrations           []EventRegistration `json:"registrations"`
	DateCreated             string              `json:"dateCreated"`
	DateUpdated             string              `json:"dateUpdated"`
}

// Decode implements the decoder interface.
func (app *Event) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Encode implements the encoder interface.
func (app Event) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

type NewEventRegistration struct {
	EventID       string  `json:"eventId"`
	ConstituentID *string `json:"constituentId"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	Email         string  `json:"email"`
	Phone         *string `json:"phone"`
	Source        string  `json:"source"`
	MatchStatus   string  `json:"matchStatus"`
}

type NewEvent struct {
	Title                   string  `json:"title"`
	Type                    string  `json:"type"`
	Status                  string  `json:"status"`
	Description             string  `json:"description"`
	Start                   string  `json:"start"`
	End                     string  `json:"end"`
	Location                string  `json:"location"`
	IsVirtual               bool    `json:"isVirtual"`
	Capacity                int     `json:"capacity"`
	RegistrationDeadline    *string `json:"registrationDeadline"`
	AutoConfirmationEnabled bool    `json:"autoConfirmationEnabled"`
	AutoReminderEnabled     bool    `json:"autoReminderEnabled"`
}

// Decode implements the decoder interface.
func (app *NewEventRegistration) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Decode implements the decoder interface.
func (app *NewEvent) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

type NewEventCheckIn struct {
	RegistrationID string `json:"registrationId"`
}

// Decode implements the decoder interface.
func (app *NewEventCheckIn) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toAppEventRegistration(registration admissionsbus.EventRegistration) EventRegistration {
	fullName := registration.FirstName
	if registration.LastName != "" {
		fullName = registration.FirstName + " " + registration.LastName
	}

	checkedInByID := uuidStringPtr(registration.CheckedInByID)

	return EventRegistration{
		ID:              registration.ID.String(),
		ConstituentID:   uuidStringPtr(registration.ConstituentID),
		ConstituentName: fullName,
		Email:           registration.Email,
		Phone:           registration.Phone,
		Status:          registration.Status.String(),
		RegisteredAt:    registration.RegisteredAt.Format(time.RFC3339),
		MatchStatus:     registration.MatchStatus.String(),
		Source:          registration.Source.String(),
		CheckedInAt:     formatTimePtr(registration.CheckedInAt),
		CheckedInByID:   checkedInByID,
		FirstName:       registration.FirstName,
		LastName:        registration.LastName,
	}
}

func toAppEventRegistrations(registrations []admissionsbus.EventRegistration) []EventRegistration {
	app := make([]EventRegistration, len(registrations))
	for i, registration := range registrations {
		app[i] = toAppEventRegistration(registration)
	}

	return app
}

func toAppEvent(event admissionsbus.Event, registrations []admissionsbus.EventRegistration) Event {
	registeredCount := len(registrations)
	checkedInCount := 0
	for _, registration := range registrations {
		if registration.Status == admissionsbus.EventRegistrationStatusCheckedIn {
			checkedInCount++
		}
	}

	return Event{
		ID:                      event.ID.String(),
		Title:                   event.Title,
		Type:                    event.Type.String(),
		Status:                  event.Status.String(),
		Description:             event.Description,
		Start:                   event.StartTime.Format(time.RFC3339),
		End:                     event.EndTime.Format(time.RFC3339),
		Location:                event.Location,
		IsVirtual:               event.IsVirtual,
		Capacity:                event.Capacity,
		RegisteredCount:         registeredCount,
		CheckedInCount:          checkedInCount,
		RegistrationDeadline:    formatTimePtr(event.RegistrationDeadline),
		AutoConfirmationEnabled: event.AutoConfirmationEnabled,
		AutoReminderEnabled:     event.AutoReminderEnabled,
		Registrations:           toAppEventRegistrations(registrations),
		DateCreated:             event.DateCreated.Format(time.RFC3339),
		DateUpdated:             event.DateUpdated.Format(time.RFC3339),
	}
}

func toBusNewEventRegistration(app NewEventRegistration) (admissionsbus.NewEventRegistration, error) {
	var fieldErrors errs.FieldErrors

	eventID, err := uuid.Parse(app.EventID)
	if err != nil {
		fieldErrors.Add("eventId", err)
	}

	constituentID, err := parseUUIDPtr(app.ConstituentID)
	if err != nil {
		fieldErrors.Add("constituentId", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewEventRegistration{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewEventRegistration{
		EventID:       eventID,
		ConstituentID: constituentID,
		FirstName:     app.FirstName,
		LastName:      app.LastName,
		Email:         app.Email,
		Phone:         app.Phone,
		Source:        admissionsbus.EventRegistrationSource(app.Source),
		MatchStatus:   admissionsbus.EventRegistrationMatchStatus(app.MatchStatus),
	}, nil
}

func toBusNewEvent(app NewEvent) (admissionsbus.NewEvent, error) {
	var fieldErrors errs.FieldErrors

	startTime, err := time.Parse(time.RFC3339, app.Start)
	if err != nil {
		fieldErrors.Add("start", err)
	}

	endTime, err := time.Parse(time.RFC3339, app.End)
	if err != nil {
		fieldErrors.Add("end", err)
	}

	registrationDeadline, err := parseTimePtr(app.RegistrationDeadline)
	if err != nil {
		fieldErrors.Add("registrationDeadline", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewEvent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewEvent{
		Title:                   app.Title,
		Type:                    admissionsbus.EventType(app.Type),
		Status:                  admissionsbus.EventStatus(app.Status),
		Description:             app.Description,
		StartTime:               startTime,
		EndTime:                 endTime,
		Location:                app.Location,
		IsVirtual:               app.IsVirtual,
		Capacity:                app.Capacity,
		RegistrationDeadline:    registrationDeadline,
		AutoConfirmationEnabled: app.AutoConfirmationEnabled,
		AutoReminderEnabled:     app.AutoReminderEnabled,
	}, nil
}

func toBusNewEventCheckIn(app NewEventCheckIn, checkedInByID uuid.UUID) (admissionsbus.NewEventCheckIn, error) {
	var fieldErrors errs.FieldErrors

	registrationID, err := uuid.Parse(app.RegistrationID)
	if err != nil {
		fieldErrors.Add("registrationId", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewEventCheckIn{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewEventCheckIn{
		RegistrationID: registrationID,
		CheckedInByID:  checkedInByID,
	}, nil
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
	ID                          string                  `json:"id"`
	FirstName                   string                  `json:"firstName"`
	LastName                    string                  `json:"lastName"`
	PreferredName               *string                 `json:"preferredName,omitempty"`
	MiddleName                  *string                 `json:"middleName,omitempty"`
	Suffix                      *string                 `json:"suffix,omitempty"`
	DateOfBirth                 string                  `json:"dateOfBirth"`
	PrimaryEmail                string                  `json:"primaryEmail"`
	PrimaryPhone                string                  `json:"primaryPhone"`
	ExternalSISID               *string                 `json:"externalSISID,omitempty"`
	NationalID                  *string                 `json:"nationalID,omitempty"`
	NationalIDVerifiedAt        *string                 `json:"nationalIDVerifiedAt,omitempty"`
	NationalIDVerifiedByAdapter *string                 `json:"nationalIDVerifiedByAdapter,omitempty"`
	UPI                         *string                 `json:"upi,omitempty"`
	UPIVerifiedAt               *string                 `json:"upiVerifiedAt,omitempty"`
	UPIVerifiedByAdapter        *string                 `json:"upiVerifiedByAdapter,omitempty"`
	KCSEIndexNumber             *string                 `json:"kcseIndexNumber,omitempty"`
	KCSEIndexVerifiedAt         *string                 `json:"kcseIndexVerifiedAt,omitempty"`
	KCSEIndexVerifiedByAdapter  *string                 `json:"kcseIndexVerifiedByAdapter,omitempty"`
	LifecycleStage              string                  `json:"lifecycleStage"`
	DuplicateStatus             string                  `json:"duplicateStatus"`
	DuplicateOfID               *string                 `json:"duplicateOfID,omitempty"`
	NotificationPreferences     NotificationPreferences `json:"notificationPreferences"`
	SISSyncedAt                 *string                 `json:"sisSyncedAt,omitempty"`
	DateCreated                 string                  `json:"dateCreated"`
	DateUpdated                 string                  `json:"dateUpdated"`
}

// NotificationPreferences represents constituent notification opt-ins and priority order.
type NotificationPreferences struct {
	SMSOptIn      bool     `json:"smsOptIn"`
	WhatsAppOptIn bool     `json:"whatsAppOptIn"`
	EmailOptIn    bool     `json:"emailOptIn"`
	Priority      []string `json:"priority"`
}

// Encode implements the encoder interface.
func (app Constituent) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppConstituent(cst admissionsbus.Constituent) Constituent {
	return Constituent{
		ID:                          cst.ID.String(),
		FirstName:                   cst.FirstName,
		LastName:                    cst.LastName,
		PreferredName:               cst.PreferredName,
		MiddleName:                  cst.MiddleName,
		Suffix:                      cst.Suffix,
		DateOfBirth:                 cst.DateOfBirth.Format(time.RFC3339),
		PrimaryEmail:                cst.PrimaryEmail.String(),
		PrimaryPhone:                cst.PrimaryPhone,
		ExternalSISID:               cst.ExternalSISID,
		NationalID:                  cst.NationalID,
		NationalIDVerifiedAt:        formatTimePtr(cst.NationalIDVerifiedAt),
		NationalIDVerifiedByAdapter: cst.NationalIDVerifiedByAdapter,
		UPI:                         cst.UPI,
		UPIVerifiedAt:               formatTimePtr(cst.UPIVerifiedAt),
		UPIVerifiedByAdapter:        cst.UPIVerifiedByAdapter,
		KCSEIndexNumber:             cst.KCSEIndexNumber,
		KCSEIndexVerifiedAt:         formatTimePtr(cst.KCSEIndexVerifiedAt),
		KCSEIndexVerifiedByAdapter:  cst.KCSEIndexVerifiedByAdapter,
		LifecycleStage:              cst.LifecycleStage.String(),
		DuplicateStatus:             cst.DuplicateStatus.String(),
		DuplicateOfID:               uuidStringPtr(cst.DuplicateOfID),
		NotificationPreferences:     toAppNotificationPreferences(cst.NotificationPreferences),
		SISSyncedAt:                 formatTimePtr(cst.SISSyncedAt),
		DateCreated:                 cst.DateCreated.Format(time.RFC3339),
		DateUpdated:                 cst.DateUpdated.Format(time.RFC3339),
	}
}

func toAppNotificationPreferences(preferences admissionsbus.NotificationPreferences) NotificationPreferences {
	return NotificationPreferences{
		SMSOptIn:      preferences.SMSOptIn,
		WhatsAppOptIn: preferences.WhatsAppOptIn,
		EmailOptIn:    preferences.EmailOptIn,
		Priority:      admissionsbus.NotificationChannelsToStrings(preferences.Priority),
	}
}

func toBusNotificationPreferences(app NotificationPreferences) (admissionsbus.NotificationPreferences, error) {
	channels, err := admissionsbus.ParseNotificationChannels(app.Priority)
	if err != nil {
		return admissionsbus.NotificationPreferences{}, err
	}

	return admissionsbus.NormalizeNotificationPreferences(admissionsbus.NotificationPreferences{
		SMSOptIn:      app.SMSOptIn,
		WhatsAppOptIn: app.WhatsAppOptIn,
		EmailOptIn:    app.EmailOptIn,
		Priority:      channels,
	})
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
	FirstName                   string                   `json:"firstName"`
	LastName                    string                   `json:"lastName"`
	PreferredName               *string                  `json:"preferredName"`
	MiddleName                  *string                  `json:"middleName"`
	Suffix                      *string                  `json:"suffix"`
	DateOfBirth                 string                   `json:"dateOfBirth"`
	PrimaryEmail                string                   `json:"primaryEmail"`
	PrimaryPhone                string                   `json:"primaryPhone"`
	ExternalSISID               *string                  `json:"externalSISID"`
	NationalID                  *string                  `json:"nationalID"`
	NationalIDVerifiedAt        *string                  `json:"nationalIDVerifiedAt"`
	NationalIDVerifiedByAdapter *string                  `json:"nationalIDVerifiedByAdapter"`
	UPI                         *string                  `json:"upi"`
	UPIVerifiedAt               *string                  `json:"upiVerifiedAt"`
	UPIVerifiedByAdapter        *string                  `json:"upiVerifiedByAdapter"`
	KCSEIndexNumber             *string                  `json:"kcseIndexNumber"`
	KCSEIndexVerifiedAt         *string                  `json:"kcseIndexVerifiedAt"`
	KCSEIndexVerifiedByAdapter  *string                  `json:"kcseIndexVerifiedByAdapter"`
	LifecycleStage              string                   `json:"lifecycleStage"`
	NotificationPreferences     *NotificationPreferences `json:"notificationPreferences"`
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

	nationalID, err := parseKenyaNationalIDPtr(app.NationalID)
	if err != nil {
		fieldErrors.Add("nationalID", err)
	}

	upi, err := parseKenyaUPIPtr(app.UPI)
	if err != nil {
		fieldErrors.Add("upi", err)
	}

	kcseIndexNumber, err := parseKenyaKCSEIndexNumberPtr(app.KCSEIndexNumber)
	if err != nil {
		fieldErrors.Add("kcseIndexNumber", err)
	}

	nationalIDVerifiedAt, err := parseTimePtr(app.NationalIDVerifiedAt)
	if err != nil {
		fieldErrors.Add("nationalIDVerifiedAt", err)
	}

	upiVerifiedAt, err := parseTimePtr(app.UPIVerifiedAt)
	if err != nil {
		fieldErrors.Add("upiVerifiedAt", err)
	}

	kcseIndexVerifiedAt, err := parseTimePtr(app.KCSEIndexVerifiedAt)
	if err != nil {
		fieldErrors.Add("kcseIndexVerifiedAt", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewConstituent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	var notificationPreferences *admissionsbus.NotificationPreferences
	if app.NotificationPreferences != nil {
		preferences, err := toBusNotificationPreferences(*app.NotificationPreferences)
		if err != nil {
			fieldErrors.Add("notificationPreferences", err)
		} else {
			notificationPreferences = &preferences
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewConstituent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewConstituent{
		FirstName:                   app.FirstName,
		LastName:                    app.LastName,
		PreferredName:               app.PreferredName,
		MiddleName:                  app.MiddleName,
		Suffix:                      app.Suffix,
		DateOfBirth:                 dob,
		PrimaryEmail:                *email,
		PrimaryPhone:                app.PrimaryPhone,
		ExternalSISID:               app.ExternalSISID,
		NationalID:                  nationalID,
		NationalIDVerifiedAt:        nationalIDVerifiedAt,
		NationalIDVerifiedByAdapter: app.NationalIDVerifiedByAdapter,
		UPI:                         upi,
		UPIVerifiedAt:               upiVerifiedAt,
		UPIVerifiedByAdapter:        app.UPIVerifiedByAdapter,
		KCSEIndexNumber:             kcseIndexNumber,
		KCSEIndexVerifiedAt:         kcseIndexVerifiedAt,
		KCSEIndexVerifiedByAdapter:  app.KCSEIndexVerifiedByAdapter,
		LifecycleStage:              stage,
		DuplicateStatus:             admissionsbus.DuplicateStatusActive,
		NotificationPreferences:     notificationPreferences,
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
	PreferredName               *string                  `json:"preferredName"`
	MiddleName                  *string                  `json:"middleName"`
	Suffix                      *string                  `json:"suffix"`
	PrimaryEmail                *string                  `json:"primaryEmail"`
	PrimaryPhone                *string                  `json:"primaryPhone"`
	NationalID                  *string                  `json:"nationalID"`
	NationalIDVerifiedAt        *string                  `json:"nationalIDVerifiedAt"`
	NationalIDVerifiedByAdapter *string                  `json:"nationalIDVerifiedByAdapter"`
	UPI                         *string                  `json:"upi"`
	UPIVerifiedAt               *string                  `json:"upiVerifiedAt"`
	UPIVerifiedByAdapter        *string                  `json:"upiVerifiedByAdapter"`
	KCSEIndexNumber             *string                  `json:"kcseIndexNumber"`
	KCSEIndexVerifiedAt         *string                  `json:"kcseIndexVerifiedAt"`
	KCSEIndexVerifiedByAdapter  *string                  `json:"kcseIndexVerifiedByAdapter"`
	LifecycleStage              *string                  `json:"lifecycleStage"`
	NotificationPreferences     *NotificationPreferences `json:"notificationPreferences"`
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

	nationalID, err := parseKenyaNationalIDPtr(app.NationalID)
	if err != nil {
		fieldErrors.Add("nationalID", err)
	}

	upi, err := parseKenyaUPIPtr(app.UPI)
	if err != nil {
		fieldErrors.Add("upi", err)
	}

	kcseIndexNumber, err := parseKenyaKCSEIndexNumberPtr(app.KCSEIndexNumber)
	if err != nil {
		fieldErrors.Add("kcseIndexNumber", err)
	}

	nationalIDVerifiedAt, err := parseTimePtr(app.NationalIDVerifiedAt)
	if err != nil {
		fieldErrors.Add("nationalIDVerifiedAt", err)
	}

	upiVerifiedAt, err := parseTimePtr(app.UPIVerifiedAt)
	if err != nil {
		fieldErrors.Add("upiVerifiedAt", err)
	}

	kcseIndexVerifiedAt, err := parseTimePtr(app.KCSEIndexVerifiedAt)
	if err != nil {
		fieldErrors.Add("kcseIndexVerifiedAt", err)
	}

	var notificationPreferences *admissionsbus.NotificationPreferences
	if app.NotificationPreferences != nil {
		preferences, err := toBusNotificationPreferences(*app.NotificationPreferences)
		if err != nil {
			fieldErrors.Add("notificationPreferences", err)
		} else {
			notificationPreferences = &preferences
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.UpdateConstituent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.UpdateConstituent{
		PreferredName:               app.PreferredName,
		MiddleName:                  app.MiddleName,
		Suffix:                      app.Suffix,
		PrimaryEmail:                email,
		PrimaryPhone:                app.PrimaryPhone,
		NationalID:                  nationalID,
		NationalIDVerifiedAt:        nationalIDVerifiedAt,
		NationalIDVerifiedByAdapter: app.NationalIDVerifiedByAdapter,
		UPI:                         upi,
		UPIVerifiedAt:               upiVerifiedAt,
		UPIVerifiedByAdapter:        app.UPIVerifiedByAdapter,
		KCSEIndexNumber:             kcseIndexNumber,
		KCSEIndexVerifiedAt:         kcseIndexVerifiedAt,
		KCSEIndexVerifiedByAdapter:  app.KCSEIndexVerifiedByAdapter,
		LifecycleStage:              stage,
		NotificationPreferences:     notificationPreferences,
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
	ID                 string                 `json:"id"`
	ConstituentID      string                 `json:"constituentID"`
	ProgramID          string                 `json:"programID"`
	AcademicTermID     string                 `json:"academicTermID"`
	ApplicationType    string                 `json:"applicationType"`
	Status             string                 `json:"status"`
	KUCCPSPlacement    *KUCCPSPlacement       `json:"kuccpsPlacement,omitempty"`
	KCSEResult         *ApplicationKCSEResult `json:"kcseResult,omitempty"`
	AssignedReviewerID *string                `json:"assignedReviewerID,omitempty"`
	SubmittedAt        *string                `json:"submittedAt,omitempty"`
	DateCreated        string                 `json:"dateCreated"`
	DateUpdated        string                 `json:"dateUpdated"`
}

// KUCCPSPlacement captures a normalized KUCCPS placement snapshot on an application.
type KUCCPSPlacement struct {
	PlacementID        string   `json:"placementID"`
	AdmissionNumber    *string  `json:"admissionNumber,omitempty"`
	InstitutionCode    string   `json:"institutionCode"`
	ProgrammeCode      string   `json:"programmeCode"`
	ProgrammeName      string   `json:"programmeName"`
	PlacementYear      int      `json:"placementYear"`
	ClusterCode        *string  `json:"clusterCode,omitempty"`
	ClusterPoints      *float64 `json:"clusterPoints,omitempty"`
	WeightedPointsNote *string  `json:"weightedPointsNote,omitempty"`
}

// ApplicationKCSESubject stores one KCSE subject grade snapshot on an application.
type ApplicationKCSESubject struct {
	SubjectCode string `json:"subjectCode"`
	Grade       string `json:"grade"`
	Points      int    `json:"points"`
}

// ApplicationKCSEResult stores the KCSE result snapshot submitted with an application.
type ApplicationKCSEResult struct {
	IndexNumber string                   `json:"indexNumber"`
	ExamYear    int                      `json:"examYear"`
	Subjects    []ApplicationKCSESubject `json:"subjects"`
	MeanGrade   string                   `json:"meanGrade"`
	MeanPoints  int                      `json:"meanPoints"`
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

// CustomFieldDefinition represents a configurable admissions field owned by a constituent or application.
type CustomFieldDefinition struct {
	ID           string   `json:"id"`
	Owner        string   `json:"owner"`
	FieldKey     string   `json:"fieldKey"`
	Label        string   `json:"label"`
	Description  *string  `json:"description,omitempty"`
	DataType     string   `json:"dataType"`
	Required     bool     `json:"required"`
	Options      []string `json:"options"`
	Validation   *string  `json:"validation,omitempty"`
	Searchable   bool     `json:"searchable"`
	Reportable   bool     `json:"reportable"`
	Importable   bool     `json:"importable"`
	Exportable   bool     `json:"exportable"`
	DisplayOrder int      `json:"displayOrder"`
	Active       bool     `json:"active"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app CustomFieldDefinition) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCustomFieldDefinition(definition admissionsbus.CustomFieldDefinition) CustomFieldDefinition {
	return CustomFieldDefinition{
		ID:           definition.ID.String(),
		Owner:        definition.Owner.String(),
		FieldKey:     definition.FieldKey,
		Label:        definition.Label,
		Description:  definition.Description,
		DataType:     definition.DataType.String(),
		Required:     definition.Required,
		Options:      definition.Options,
		Validation:   definition.Validation,
		Searchable:   definition.Searchable,
		Reportable:   definition.Reportable,
		Importable:   definition.Importable,
		Exportable:   definition.Exportable,
		DisplayOrder: definition.DisplayOrder,
		Active:       definition.Active,
		DateCreated:  definition.DateCreated.Format(time.RFC3339),
		DateUpdated:  definition.DateUpdated.Format(time.RFC3339),
	}
}

func toAppCustomFieldDefinitions(definitions []admissionsbus.CustomFieldDefinition) []CustomFieldDefinition {
	app := make([]CustomFieldDefinition, len(definitions))
	for i, definition := range definitions {
		app[i] = toAppCustomFieldDefinition(definition)
	}

	return app
}

// NewCustomFieldDefinition defines the data needed to create or update an admissions custom field definition.
type NewCustomFieldDefinition struct {
	Owner        string   `json:"owner"`
	FieldKey     string   `json:"fieldKey"`
	Label        string   `json:"label"`
	Description  *string  `json:"description"`
	DataType     string   `json:"dataType"`
	Required     bool     `json:"required"`
	Options      []string `json:"options"`
	Validation   *string  `json:"validation"`
	Searchable   bool     `json:"searchable"`
	Reportable   bool     `json:"reportable"`
	Importable   bool     `json:"importable"`
	Exportable   bool     `json:"exportable"`
	DisplayOrder int      `json:"displayOrder"`
	Active       bool     `json:"active"`
}

// Decode implements the decoder interface.
func (app *NewCustomFieldDefinition) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewCustomFieldDefinition(app NewCustomFieldDefinition) admissionsbus.NewCustomFieldDefinition {
	return admissionsbus.NewCustomFieldDefinition{
		Owner:        admissionsbus.CustomFieldOwner(app.Owner),
		FieldKey:     app.FieldKey,
		Label:        app.Label,
		Description:  app.Description,
		DataType:     admissionsbus.CustomFieldDataType(app.DataType),
		Required:     app.Required,
		Options:      app.Options,
		Validation:   app.Validation,
		Searchable:   app.Searchable,
		Reportable:   app.Reportable,
		Importable:   app.Importable,
		Exportable:   app.Exportable,
		DisplayOrder: app.DisplayOrder,
		Active:       app.Active,
	}
}

func toBusNewCustomFieldValue(app NewCustomFieldValue) (admissionsbus.NewCustomFieldValue, error) {
	var fieldErrors errs.FieldErrors

	definitionID, err := uuid.Parse(app.DefinitionID)
	if err != nil {
		fieldErrors.Add("definitionID", err)
	}

	ownerID, err := uuid.Parse(app.OwnerID)
	if err != nil {
		fieldErrors.Add("ownerID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewCustomFieldValue{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewCustomFieldValue{
		DefinitionID: definitionID,
		Owner:        admissionsbus.CustomFieldOwner(app.Owner),
		OwnerID:      ownerID,
		Value:        app.Value,
	}, nil
}

// CustomFieldValue represents one custom field value for a constituent or application.
type CustomFieldValue struct {
	ID           string `json:"id"`
	DefinitionID string `json:"definitionID"`
	Owner        string `json:"owner"`
	OwnerID      string `json:"ownerID"`
	Value        string `json:"value"`
	DateCreated  string `json:"dateCreated"`
	DateUpdated  string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app CustomFieldValue) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCustomFieldValue(value admissionsbus.CustomFieldValue) CustomFieldValue {
	return CustomFieldValue{
		ID:           value.ID.String(),
		DefinitionID: value.DefinitionID.String(),
		Owner:        value.Owner.String(),
		OwnerID:      value.OwnerID.String(),
		Value:        value.Value,
		DateCreated:  value.DateCreated.Format(time.RFC3339),
		DateUpdated:  value.DateUpdated.Format(time.RFC3339),
	}
}

func toAppCustomFieldValues(values []admissionsbus.CustomFieldValue) []CustomFieldValue {
	app := make([]CustomFieldValue, len(values))
	for i, value := range values {
		app[i] = toAppCustomFieldValue(value)
	}

	return app
}

// NewCustomFieldValue defines the data needed to set one custom field value.
type NewCustomFieldValue struct {
	DefinitionID string `json:"definitionID"`
	Owner        string `json:"owner"`
	OwnerID      string `json:"ownerID"`
	Value        string `json:"value"`
}

// Decode implements the decoder interface.
func (app *NewCustomFieldValue) Decode(data []byte) error {
	return json.Unmarshal(data, app)
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

// ChecklistItem represents one document requirement for an application.
type ChecklistItem struct {
	ID            string  `json:"id"`
	ApplicationID string  `json:"applicationID"`
	ItemKey       string  `json:"itemKey"`
	DocumentName  string  `json:"documentName"`
	Description   *string `json:"description,omitempty"`
	Required      bool    `json:"required"`
	Status        string  `json:"status"`
	DisplayOrder  int     `json:"displayOrder"`
	DateCreated   string  `json:"dateCreated"`
	DateUpdated   string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ChecklistItem) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppChecklistItem(item admissionsbus.ChecklistItem) ChecklistItem {
	return ChecklistItem{
		ID:            item.ID.String(),
		ApplicationID: item.ApplicationID.String(),
		ItemKey:       item.ItemKey,
		DocumentName:  item.DocumentName,
		Description:   item.Description,
		Required:      item.Required,
		Status:        item.Status.String(),
		DisplayOrder:  item.DisplayOrder,
		DateCreated:   item.DateCreated.Format(time.RFC3339),
		DateUpdated:   item.DateUpdated.Format(time.RFC3339),
	}
}

func toAppChecklistItems(items []admissionsbus.ChecklistItem) []ChecklistItem {
	app := make([]ChecklistItem, len(items))
	for i, item := range items {
		app[i] = toAppChecklistItem(item)
	}

	return app
}

// Document represents uploaded admissions document metadata.
type Document struct {
	ID              string  `json:"id"`
	ApplicationID   string  `json:"applicationID"`
	ChecklistItemID string  `json:"checklistItemID"`
	FileName        string  `json:"fileName"`
	ContentType     string  `json:"contentType"`
	SizeBytes       int64   `json:"sizeBytes"`
	StorageKey      string  `json:"storageKey"`
	Status          string  `json:"status"`
	ReviewerID      *string `json:"reviewerID,omitempty"`
	ReviewerNotes   *string `json:"reviewerNotes,omitempty"`
	UploadedByID    string  `json:"uploadedByID"`
	UploadedAt      string  `json:"uploadedAt"`
	ReviewedAt      *string `json:"reviewedAt,omitempty"`
	DateCreated     string  `json:"dateCreated"`
	DateUpdated     string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Document) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppDocument(document admissionsbus.Document) Document {
	return Document{
		ID:              document.ID.String(),
		ApplicationID:   document.ApplicationID.String(),
		ChecklistItemID: document.ChecklistItemID.String(),
		FileName:        document.FileName,
		ContentType:     document.ContentType,
		SizeBytes:       document.SizeBytes,
		StorageKey:      document.StorageKey,
		Status:          document.Status.String(),
		ReviewerID:      uuidStringPtr(document.ReviewerID),
		ReviewerNotes:   document.ReviewerNotes,
		UploadedByID:    document.UploadedByID.String(),
		UploadedAt:      document.UploadedAt.Format(time.RFC3339),
		ReviewedAt:      formatTimePtr(document.ReviewedAt),
		DateCreated:     document.DateCreated.Format(time.RFC3339),
		DateUpdated:     document.DateUpdated.Format(time.RFC3339),
	}
}

func toAppDocuments(documents []admissionsbus.Document) []Document {
	app := make([]Document, len(documents))
	for i, document := range documents {
		app[i] = toAppDocument(document)
	}

	return app
}

// ImportBatch represents an admissions CSV or Excel import batch.
type ImportBatch struct {
	ID                string            `json:"id"`
	Source            string            `json:"source"`
	FileType          string            `json:"fileType"`
	Target            string            `json:"target"`
	Status            string            `json:"status"`
	FileName          string            `json:"fileName"`
	StorageKey        *string           `json:"storageKey,omitempty"`
	UploadedByID      string            `json:"uploadedByID"`
	TotalRows         int               `json:"totalRows"`
	ValidRows         int               `json:"validRows"`
	InvalidRows       int               `json:"invalidRows"`
	DuplicateRows     int               `json:"duplicateRows"`
	FieldMapping      map[string]string `json:"fieldMapping"`
	InvalidReportKey  *string           `json:"invalidReportKey,omitempty"`
	ValidationSummary *string           `json:"validationSummary,omitempty"`
	CommittedAt       *string           `json:"committedAt,omitempty"`
	DateCreated       string            `json:"dateCreated"`
	DateUpdated       string            `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ImportBatch) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppImportBatch(batch admissionsbus.ImportBatch) ImportBatch {
	return ImportBatch{
		ID:                batch.ID.String(),
		Source:            batch.Source.String(),
		FileType:          batch.FileType.String(),
		Target:            batch.Target.String(),
		Status:            batch.Status.String(),
		FileName:          batch.FileName,
		StorageKey:        batch.StorageKey,
		UploadedByID:      batch.UploadedByID.String(),
		TotalRows:         batch.TotalRows,
		ValidRows:         batch.ValidRows,
		InvalidRows:       batch.InvalidRows,
		DuplicateRows:     batch.DuplicateRows,
		FieldMapping:      batch.FieldMapping,
		InvalidReportKey:  batch.InvalidReportKey,
		ValidationSummary: batch.ValidationSummary,
		CommittedAt:       formatTimePtr(batch.CommittedAt),
		DateCreated:       batch.DateCreated.Format(time.RFC3339),
		DateUpdated:       batch.DateUpdated.Format(time.RFC3339),
	}
}

func toAppImportBatches(batches []admissionsbus.ImportBatch) []ImportBatch {
	app := make([]ImportBatch, len(batches))
	for i, batch := range batches {
		app[i] = toAppImportBatch(batch)
	}

	return app
}

// ImportInvalidRow represents one invalid import row available for correction download.
type ImportInvalidRow struct {
	ID          string            `json:"id"`
	BatchID     string            `json:"batchID"`
	RowNumber   int               `json:"rowNumber"`
	FieldName   *string           `json:"fieldName,omitempty"`
	RawData     map[string]string `json:"rawData"`
	ErrorCode   string            `json:"errorCode"`
	ErrorDetail string            `json:"errorDetail"`
	DateCreated string            `json:"dateCreated"`
}

// Encode implements the encoder interface.
func (app ImportInvalidRow) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppImportInvalidRow(row admissionsbus.ImportInvalidRow) ImportInvalidRow {
	return ImportInvalidRow{
		ID:          row.ID.String(),
		BatchID:     row.BatchID.String(),
		RowNumber:   row.RowNumber,
		FieldName:   row.FieldName,
		RawData:     row.RawData,
		ErrorCode:   row.ErrorCode,
		ErrorDetail: row.ErrorDetail,
		DateCreated: row.DateCreated.Format(time.RFC3339),
	}
}

func toAppImportInvalidRows(rows []admissionsbus.ImportInvalidRow) []ImportInvalidRow {
	app := make([]ImportInvalidRow, len(rows))
	for i, row := range rows {
		app[i] = toAppImportInvalidRow(row)
	}

	return app
}

// NewImportBatch defines the data needed to record an import preview or commit.
type NewImportBatch struct {
	Source            string            `json:"source"`
	FileType          string            `json:"fileType"`
	Target            string            `json:"target"`
	Status            string            `json:"status"`
	FileName          string            `json:"fileName"`
	StorageKey        *string           `json:"storageKey"`
	TotalRows         int               `json:"totalRows"`
	ValidRows         int               `json:"validRows"`
	InvalidRows       int               `json:"invalidRows"`
	DuplicateRows     int               `json:"duplicateRows"`
	FieldMapping      map[string]string `json:"fieldMapping"`
	InvalidReportKey  *string           `json:"invalidReportKey"`
	ValidationSummary *string           `json:"validationSummary"`
}

// Decode implements the decoder interface.
func (app *NewImportBatch) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewImportBatch(app NewImportBatch, uploadedByID uuid.UUID) admissionsbus.NewImportBatch {
	return admissionsbus.NewImportBatch{
		Source:            admissionsbus.ImportSource(app.Source),
		FileType:          admissionsbus.ImportFileType(app.FileType),
		Target:            admissionsbus.ImportTarget(app.Target),
		Status:            admissionsbus.ImportBatchStatus(app.Status),
		FileName:          app.FileName,
		StorageKey:        app.StorageKey,
		UploadedByID:      uploadedByID,
		TotalRows:         app.TotalRows,
		ValidRows:         app.ValidRows,
		InvalidRows:       app.InvalidRows,
		DuplicateRows:     app.DuplicateRows,
		FieldMapping:      app.FieldMapping,
		InvalidReportKey:  app.InvalidReportKey,
		ValidationSummary: app.ValidationSummary,
	}
}

// NewImportInvalidRow defines invalid row details attached to an import batch.
type NewImportInvalidRow struct {
	RowNumber   int               `json:"rowNumber"`
	FieldName   *string           `json:"fieldName"`
	RawData     map[string]string `json:"rawData"`
	ErrorCode   string            `json:"errorCode"`
	ErrorDetail string            `json:"errorDetail"`
}

// NewImportInvalidRows defines the rows attached to an import batch.
type NewImportInvalidRows struct {
	Rows []NewImportInvalidRow `json:"rows"`
}

// Decode implements the decoder interface.
func (app *NewImportInvalidRows) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewImportInvalidRows(app NewImportInvalidRows, batchID uuid.UUID) []admissionsbus.NewImportInvalidRow {
	rows := make([]admissionsbus.NewImportInvalidRow, len(app.Rows))
	for i, row := range app.Rows {
		rows[i] = admissionsbus.NewImportInvalidRow{
			BatchID:     batchID,
			RowNumber:   row.RowNumber,
			FieldName:   row.FieldName,
			RawData:     row.RawData,
			ErrorCode:   row.ErrorCode,
			ErrorDetail: row.ErrorDetail,
		}
	}

	return rows
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
		KUCCPSPlacement:    toAppKUCCPSPlacement(application.KUCCPSPlacement),
		KCSEResult:         toAppApplicationKCSEResult(application.KCSEResult),
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
	ConstituentID      string                 `json:"constituentID"`
	ProgramID          string                 `json:"programID"`
	AcademicTermID     string                 `json:"academicTermID"`
	ApplicationType    string                 `json:"applicationType"`
	KUCCPSPlacement    *KUCCPSPlacement       `json:"kuccpsPlacement"`
	KCSEResult         *ApplicationKCSEResult `json:"kcseResult"`
	AssignedReviewerID *string                `json:"assignedReviewerID"`
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

// NewChecklistItem defines the data needed to create or update an application checklist item.
type NewChecklistItem struct {
	ItemKey      string  `json:"itemKey"`
	DocumentName string  `json:"documentName"`
	Description  *string `json:"description"`
	Required     bool    `json:"required"`
	DisplayOrder int     `json:"displayOrder"`
}

// Decode implements the decoder interface.
func (app *NewChecklistItem) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewChecklistItem(app NewChecklistItem, applicationID uuid.UUID) admissionsbus.NewChecklistItem {
	return admissionsbus.NewChecklistItem{
		ApplicationID: applicationID,
		ItemKey:       app.ItemKey,
		DocumentName:  app.DocumentName,
		Description:   app.Description,
		Required:      app.Required,
		DisplayOrder:  app.DisplayOrder,
	}
}

// NewDocument defines uploaded document metadata. File bytes are stored outside the CRM database.
type NewDocument struct {
	ChecklistItemID string `json:"checklistItemID"`
	FileName        string `json:"fileName"`
	ContentType     string `json:"contentType"`
	SizeBytes       int64  `json:"sizeBytes"`
	StorageKey      string `json:"storageKey"`
}

// Decode implements the decoder interface.
func (app *NewDocument) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewDocument(app NewDocument, applicationID uuid.UUID, uploadedByID uuid.UUID) (admissionsbus.NewDocument, error) {
	var fieldErrors errs.FieldErrors

	checklistItemID, err := uuid.Parse(app.ChecklistItemID)
	if err != nil {
		fieldErrors.Add("checklistItemID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewDocument{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewDocument{
		ApplicationID:   applicationID,
		ChecklistItemID: checklistItemID,
		FileName:        app.FileName,
		ContentType:     app.ContentType,
		SizeBytes:       app.SizeBytes,
		StorageKey:      app.StorageKey,
		UploadedByID:    uploadedByID,
	}, nil
}

// NewDocumentVerification defines reviewer action for uploaded document metadata.
type NewDocumentVerification struct {
	Status        string  `json:"status"`
	ReviewerNotes *string `json:"reviewerNotes"`
}

// Decode implements the decoder interface.
func (app *NewDocumentVerification) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewDocumentVerification(app NewDocumentVerification, reviewerID uuid.UUID) admissionsbus.NewDocumentVerification {
	return admissionsbus.NewDocumentVerification{
		Status:        admissionsbus.DocumentStatus(app.Status),
		ReviewerID:    reviewerID,
		ReviewerNotes: app.ReviewerNotes,
	}
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
		KUCCPSPlacement:    toBusKUCCPSPlacement(app.KUCCPSPlacement),
		KCSEResult:         toBusApplicationKCSEResult(app.KCSEResult),
		AssignedReviewerID: assignedReviewerID,
	}, nil
}

func toAppKUCCPSPlacement(placement *admissionsbus.KUCCPSPlacement) *KUCCPSPlacement {
	if placement == nil {
		return nil
	}

	return &KUCCPSPlacement{
		PlacementID:        placement.PlacementID,
		AdmissionNumber:    placement.AdmissionNumber,
		InstitutionCode:    placement.InstitutionCode,
		ProgrammeCode:      placement.ProgrammeCode,
		ProgrammeName:      placement.ProgrammeName,
		PlacementYear:      placement.PlacementYear,
		ClusterCode:        placement.ClusterCode,
		ClusterPoints:      placement.ClusterPoints,
		WeightedPointsNote: placement.WeightedPointsNote,
	}
}

func toBusKUCCPSPlacement(placement *KUCCPSPlacement) *admissionsbus.KUCCPSPlacement {
	if placement == nil {
		return nil
	}

	return &admissionsbus.KUCCPSPlacement{
		PlacementID:        placement.PlacementID,
		AdmissionNumber:    placement.AdmissionNumber,
		InstitutionCode:    placement.InstitutionCode,
		ProgrammeCode:      placement.ProgrammeCode,
		ProgrammeName:      placement.ProgrammeName,
		PlacementYear:      placement.PlacementYear,
		ClusterCode:        placement.ClusterCode,
		ClusterPoints:      placement.ClusterPoints,
		WeightedPointsNote: placement.WeightedPointsNote,
	}
}

func toAppApplicationKCSEResult(result *admissionsbus.ApplicationKCSEResult) *ApplicationKCSEResult {
	if result == nil {
		return nil
	}

	return &ApplicationKCSEResult{
		IndexNumber: result.IndexNumber,
		ExamYear:    result.ExamYear,
		Subjects:    toAppApplicationKCSESubjects(result.Subjects),
		MeanGrade:   result.MeanGrade,
		MeanPoints:  result.MeanPoints,
	}
}

func toBusApplicationKCSEResult(result *ApplicationKCSEResult) *admissionsbus.ApplicationKCSEResult {
	if result == nil {
		return nil
	}

	return &admissionsbus.ApplicationKCSEResult{
		IndexNumber: result.IndexNumber,
		ExamYear:    result.ExamYear,
		Subjects:    toBusApplicationKCSESubjects(result.Subjects),
		MeanGrade:   result.MeanGrade,
		MeanPoints:  result.MeanPoints,
	}
}

func toAppApplicationKCSESubjects(subjects []admissionsbus.ApplicationKCSESubject) []ApplicationKCSESubject {
	app := make([]ApplicationKCSESubject, len(subjects))
	for i, subject := range subjects {
		app[i] = ApplicationKCSESubject{
			SubjectCode: subject.SubjectCode,
			Grade:       subject.Grade,
			Points:      subject.Points,
		}
	}

	return app
}

func toBusApplicationKCSESubjects(subjects []ApplicationKCSESubject) []admissionsbus.ApplicationKCSESubject {
	bus := make([]admissionsbus.ApplicationKCSESubject, len(subjects))
	for i, subject := range subjects {
		bus[i] = admissionsbus.ApplicationKCSESubject{
			SubjectCode: subject.SubjectCode,
			Grade:       subject.Grade,
			Points:      subject.Points,
		}
	}

	return bus
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

// SyncJob represents a SIS batch reconciliation run.
type SyncJob struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	Direction      string  `json:"direction"`
	StartedAt      *string `json:"startedAt,omitempty"`
	CompletedAt    *string `json:"completedAt,omitempty"`
	RecordsPulled  int     `json:"recordsPulled"`
	RecordsPushed  int     `json:"recordsPushed"`
	EventsRequeued int     `json:"eventsRequeued"`
	FailureReason  *string `json:"failureReason,omitempty"`
	Retryable      bool    `json:"retryable"`
	CreatedByID    *string `json:"createdByID,omitempty"`
	DateCreated    string  `json:"dateCreated"`
	DateUpdated    string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app SyncJob) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppSyncJob(job admissionsbus.SyncJob) SyncJob {
	return SyncJob{
		ID:             job.ID.String(),
		Name:           job.Name,
		Status:         job.Status.String(),
		Direction:      job.Direction.String(),
		StartedAt:      formatTimePtr(job.StartedAt),
		CompletedAt:    formatTimePtr(job.CompletedAt),
		RecordsPulled:  job.RecordsPulled,
		RecordsPushed:  job.RecordsPushed,
		EventsRequeued: job.EventsRequeued,
		FailureReason:  job.FailureReason,
		Retryable:      job.Retryable,
		CreatedByID:    uuidStringPtr(job.CreatedByID),
		DateCreated:    job.DateCreated.Format(time.RFC3339),
		DateUpdated:    job.DateUpdated.Format(time.RFC3339),
	}
}

func toAppSyncJobs(jobs []admissionsbus.SyncJob) []SyncJob {
	app := make([]SyncJob, len(jobs))
	for i, job := range jobs {
		app[i] = toAppSyncJob(job)
	}

	return app
}

// NewSyncJob defines the data needed to create a SIS batch reconciliation run.
type NewSyncJob struct {
	Name        string  `json:"name"`
	Direction   string  `json:"direction"`
	Status      string  `json:"status"`
	StartedAt   *string `json:"startedAt"`
	CreatedByID *string `json:"createdByID"`
}

// Decode implements the decoder interface.
func (app *NewSyncJob) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewSyncJob(app NewSyncJob) (admissionsbus.NewSyncJob, error) {
	var fieldErrors errs.FieldErrors

	startedAt, err := parseTimePtr(app.StartedAt)
	if err != nil {
		fieldErrors.Add("startedAt", err)
	}

	createdByID, err := parseUUIDPtr(app.CreatedByID)
	if err != nil {
		fieldErrors.Add("createdByID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewSyncJob{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewSyncJob{
		Name:        app.Name,
		Direction:   admissionsbus.SyncDirection(app.Direction),
		Status:      admissionsbus.SyncJobStatus(app.Status),
		StartedAt:   startedAt,
		CreatedByID: createdByID,
	}, nil
}

// SyncEvent represents a selected real-time SIS sync event.
type SyncEvent struct {
	ID            string  `json:"id"`
	JobID         *string `json:"jobID,omitempty"`
	EventType     string  `json:"eventType"`
	Status        string  `json:"status"`
	Direction     string  `json:"direction"`
	ResourceType  string  `json:"resourceType"`
	ResourceID    string  `json:"resourceID"`
	PayloadHash   string  `json:"payloadHash"`
	Attempts      int     `json:"attempts"`
	NextRetryAt   *string `json:"nextRetryAt,omitempty"`
	FailureReason *string `json:"failureReason,omitempty"`
	AuditMessage  string  `json:"auditMessage"`
	DateCreated   string  `json:"dateCreated"`
	DateUpdated   string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app SyncEvent) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppSyncEvent(event admissionsbus.SyncEvent) SyncEvent {
	return SyncEvent{
		ID:            event.ID.String(),
		JobID:         uuidStringPtr(event.JobID),
		EventType:     event.EventType.String(),
		Status:        event.Status.String(),
		Direction:     event.Direction.String(),
		ResourceType:  event.ResourceType,
		ResourceID:    event.ResourceID.String(),
		PayloadHash:   event.PayloadHash,
		Attempts:      event.Attempts,
		NextRetryAt:   formatTimePtr(event.NextRetryAt),
		FailureReason: event.FailureReason,
		AuditMessage:  event.AuditMessage,
		DateCreated:   event.DateCreated.Format(time.RFC3339),
		DateUpdated:   event.DateUpdated.Format(time.RFC3339),
	}
}

func toAppSyncEvents(events []admissionsbus.SyncEvent) []SyncEvent {
	app := make([]SyncEvent, len(events))
	for i, event := range events {
		app[i] = toAppSyncEvent(event)
	}

	return app
}

// NewSyncEvent defines the data needed to enqueue a selected real-time SIS sync event.
type NewSyncEvent struct {
	JobID        *string `json:"jobID"`
	EventType    string  `json:"eventType"`
	Direction    string  `json:"direction"`
	ResourceType string  `json:"resourceType"`
	ResourceID   string  `json:"resourceID"`
	PayloadHash  string  `json:"payloadHash"`
	AuditMessage string  `json:"auditMessage"`
}

// Decode implements the decoder interface.
func (app *NewSyncEvent) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewSyncEvent(app NewSyncEvent) (admissionsbus.NewSyncEvent, error) {
	var fieldErrors errs.FieldErrors

	jobID, err := parseUUIDPtr(app.JobID)
	if err != nil {
		fieldErrors.Add("jobID", err)
	}

	resourceID, err := uuid.Parse(app.ResourceID)
	if err != nil {
		fieldErrors.Add("resourceID", err)
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.NewSyncEvent{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	return admissionsbus.NewSyncEvent{
		JobID:        jobID,
		EventType:    admissionsbus.SyncEventType(app.EventType),
		Direction:    admissionsbus.SyncDirection(app.Direction),
		ResourceType: app.ResourceType,
		ResourceID:   resourceID,
		PayloadHash:  app.PayloadHash,
		AuditMessage: app.AuditMessage,
	}, nil
}

// CampaignAuditEvent represents one lifecycle action for an admissions campaign.
type CampaignAuditEvent struct {
	ID          string `json:"id"`
	CampaignID  string `json:"campaignID"`
	ActorName   string `json:"actorName"`
	Action      string `json:"action"`
	OccurredAt  string `json:"occurredAt"`
	DateCreated string `json:"dateCreated"`
}

// Campaign represents an admissions marketing or operational campaign.
type Campaign struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Status         string               `json:"status"`
	Channel        string               `json:"channel"`
	AudienceName   string               `json:"audienceName"`
	TemplateName   string               `json:"templateName"`
	MessagePreview string               `json:"messagePreview"`
	Segment        json.RawMessage      `json:"segment"`
	Metrics        json.RawMessage      `json:"metrics"`
	StartsAt       *string              `json:"startsAt,omitempty"`
	EndsAt         *string              `json:"endsAt,omitempty"`
	CreatedByID    *string              `json:"createdByID,omitempty"`
	AuditTrail     []CampaignAuditEvent `json:"auditTrail"`
	DateCreated    string               `json:"dateCreated"`
	DateUpdated    string               `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Campaign) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCampaign(campaign admissionsbus.Campaign, auditTrail []admissionsbus.CampaignAuditEvent) Campaign {
	return Campaign{
		ID:             campaign.ID.String(),
		Name:           campaign.Name,
		Status:         campaign.Status.String(),
		Channel:        campaign.Channel.String(),
		AudienceName:   campaign.AudienceName,
		TemplateName:   campaign.TemplateName,
		MessagePreview: campaign.MessagePreview,
		Segment:        campaign.Segment,
		Metrics:        campaign.Metrics,
		StartsAt:       formatTimePtr(campaign.StartsAt),
		EndsAt:         formatTimePtr(campaign.EndsAt),
		CreatedByID:    uuidStringPtr(campaign.CreatedByID),
		AuditTrail:     toAppCampaignAuditEvents(auditTrail),
		DateCreated:    campaign.DateCreated.Format(time.RFC3339),
		DateUpdated:    campaign.DateUpdated.Format(time.RFC3339),
	}
}

func toAppCampaigns(campaigns []admissionsbus.Campaign, audits map[uuid.UUID][]admissionsbus.CampaignAuditEvent) []Campaign {
	app := make([]Campaign, len(campaigns))
	for i, campaign := range campaigns {
		app[i] = toAppCampaign(campaign, audits[campaign.ID])
	}

	return app
}

func toAppCampaignAuditEvent(event admissionsbus.CampaignAuditEvent) CampaignAuditEvent {
	return CampaignAuditEvent{
		ID:          event.ID.String(),
		CampaignID:  event.CampaignID.String(),
		ActorName:   event.ActorName,
		Action:      event.Action,
		OccurredAt:  event.OccurredAt.Format(time.RFC3339),
		DateCreated: event.DateCreated.Format(time.RFC3339),
	}
}

func toAppCampaignAuditEvents(events []admissionsbus.CampaignAuditEvent) []CampaignAuditEvent {
	app := make([]CampaignAuditEvent, len(events))
	for i, event := range events {
		app[i] = toAppCampaignAuditEvent(event)
	}

	return app
}

// Communication represents one inbound, outbound, provider-tracked, or manually logged touchpoint.
type Communication struct {
	ID                string          `json:"id"`
	ExternalMessageID string          `json:"externalMessageID"`
	Channel           string          `json:"channel"`
	Direction         string          `json:"direction"`
	ConstituentID     string          `json:"constituentID"`
	ApplicationID     *string         `json:"applicationID,omitempty"`
	CampaignID        *string         `json:"campaignID,omitempty"`
	RecipientSender   string          `json:"recipientSender"`
	RecipientInitials string          `json:"recipientInitials"`
	Subject           string          `json:"subject"`
	Preview           string          `json:"preview"`
	Status            string          `json:"status"`
	Provider          *string         `json:"provider,omitempty"`
	OwnerName         string          `json:"ownerName"`
	Outcome           *string         `json:"outcome,omitempty"`
	Duration          *string         `json:"duration,omitempty"`
	OccurredAt        string          `json:"occurredAt"`
	ProviderPayload   json.RawMessage `json:"providerPayload,omitempty"`
	DateCreated       string          `json:"dateCreated"`
	DateUpdated       string          `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Communication) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCommunication(communication admissionsbus.Communication) Communication {
	return Communication{
		ID:                communication.ID.String(),
		ExternalMessageID: communication.ExternalMessageID,
		Channel:           communication.Channel.String(),
		Direction:         communication.Direction.String(),
		ConstituentID:     communication.ConstituentID.String(),
		ApplicationID:     uuidStringPtr(communication.ApplicationID),
		CampaignID:        uuidStringPtr(communication.CampaignID),
		RecipientSender:   communication.RecipientSender,
		RecipientInitials: communication.RecipientInitials,
		Subject:           communication.Subject,
		Preview:           communication.Preview,
		Status:            communication.Status.String(),
		Provider:          communication.Provider,
		OwnerName:         communication.OwnerName,
		Outcome:           communication.Outcome,
		Duration:          communication.Duration,
		OccurredAt:        communication.OccurredAt.Format(time.RFC3339),
		ProviderPayload:   communication.ProviderPayload,
		DateCreated:       communication.DateCreated.Format(time.RFC3339),
		DateUpdated:       communication.DateUpdated.Format(time.RFC3339),
	}
}

func toAppCommunications(communications []admissionsbus.Communication) []Communication {
	app := make([]Communication, len(communications))
	for i, communication := range communications {
		app[i] = toAppCommunication(communication)
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

func uuidStringPtr(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	formatted := id.String()
	return &formatted
}

func parseTimePtr(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseKenyaNationalIDPtr(value *string) (*string, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	id, err := admissionsbus.ParseKenyaNationalID(*value)
	if err != nil {
		return nil, err
	}

	normalized := id.String()
	return &normalized, nil
}

func parseKenyaUPIPtr(value *string) (*string, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	upi, err := admissionsbus.ParseKenyaUPI(*value)
	if err != nil {
		return nil, err
	}

	normalized := upi.String()
	return &normalized, nil
}

func parseKenyaKCSEIndexNumberPtr(value *string) (*string, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	indexNumber, err := admissionsbus.ParseKenyaKCSEIndexNumber(*value)
	if err != nil {
		return nil, err
	}

	normalized := indexNumber.String()
	return &normalized, nil
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
