package admissionsdb

import (
	"encoding/json"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb/dbarray"
)

type staffProfileDB struct {
	ID          uuid.UUID      `db:"staff_profile_id"`
	UserID      uuid.UUID      `db:"user_id"`
	Roles       dbarray.String `db:"admissions_roles"`
	Active      bool           `db:"is_active"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

type applicantProfileDB struct {
	ID            uuid.UUID `db:"applicant_profile_id"`
	UserID        uuid.UUID `db:"user_id"`
	ConstituentID uuid.UUID `db:"constituent_id"`
	Active        bool      `db:"is_active"`
	DateCreated   time.Time `db:"date_created"`
	DateUpdated   time.Time `db:"date_updated"`
}

type leadScoreRuleDB struct {
	ID          uuid.UUID       `db:"lead_score_rule_id"`
	Name        string          `db:"name"`
	Description *string         `db:"description"`
	Criteria    json.RawMessage `db:"criteria"`
	Points      int             `db:"points"`
	Active      bool            `db:"is_active"`
	Priority    int             `db:"priority"`
	DateCreated time.Time       `db:"date_created"`
	DateUpdated time.Time       `db:"date_updated"`
}

type leadScoreDB struct {
	ID             uuid.UUID       `db:"lead_score_id"`
	ConstituentID  uuid.UUID       `db:"constituent_id"`
	TotalScore     int             `db:"total_score"`
	Band           string          `db:"score_band"`
	Breakdown      json.RawMessage `db:"breakdown"`
	RecalculatedAt time.Time       `db:"recalculated_at"`
	DateCreated    time.Time       `db:"date_created"`
	DateUpdated    time.Time       `db:"date_updated"`
}

type constituentDB struct {
	ID                          uuid.UUID  `db:"constituent_id"`
	FirstName                   string     `db:"first_name"`
	LastName                    string     `db:"last_name"`
	PreferredName               *string    `db:"preferred_name"`
	MiddleName                  *string    `db:"middle_name"`
	Suffix                      *string    `db:"suffix"`
	DateOfBirth                 time.Time  `db:"date_of_birth"`
	PrimaryEmail                string     `db:"primary_email"`
	PrimaryPhone                string     `db:"primary_phone"`
	ExternalSISID               *string    `db:"external_sis_id"`
	NationalID                  *string    `db:"national_id"`
	NationalIDVerifiedAt        *time.Time `db:"national_id_verified_at"`
	NationalIDVerifiedByAdapter *string    `db:"national_id_verified_by_adapter"`
	UPI                         *string    `db:"upi"`
	UPIVerifiedAt               *time.Time `db:"upi_verified_at"`
	UPIVerifiedByAdapter        *string    `db:"upi_verified_by_adapter"`
	KCSEIndexNumber             *string    `db:"kcse_index_number"`
	KCSEIndexVerifiedAt         *time.Time `db:"kcse_index_verified_at"`
	KCSEIndexVerifiedByAdapter  *string    `db:"kcse_index_verified_by_adapter"`
	LifecycleStage              string     `db:"lifecycle_stage"`
	DuplicateStatus             string     `db:"duplicate_status"`
	DuplicateOfID               *uuid.UUID `db:"duplicate_of_id"`
	SISSyncedAt                 *time.Time `db:"sis_synced_at"`
	DateCreated                 time.Time  `db:"date_created"`
	DateUpdated                 time.Time  `db:"date_updated"`
}

type inquiryDB struct {
	ID                uuid.UUID  `db:"inquiry_id"`
	ConstituentID     uuid.UUID  `db:"constituent_id"`
	FirstName         string     `db:"first_name"`
	LastName          string     `db:"last_name"`
	DateOfBirth       time.Time  `db:"date_of_birth"`
	PrimaryEmail      string     `db:"primary_email"`
	PrimaryPhone      string     `db:"primary_phone"`
	ProgramOfInterest *uuid.UUID `db:"program_of_interest"`
	TermOfInterest    *uuid.UUID `db:"term_of_interest"`
	Source            string     `db:"source"`
	UTMSource         *string    `db:"utm_source"`
	UTMMedium         *string    `db:"utm_medium"`
	UTMCampaign       *string    `db:"utm_campaign"`
	Message           *string    `db:"message"`
	Status            string     `db:"status"`
	DateCreated       time.Time  `db:"date_created"`
	DateUpdated       time.Time  `db:"date_updated"`
}

type programDB struct {
	ID            uuid.UUID  `db:"program_id"`
	ExternalSISID string     `db:"external_sis_id"`
	Name          string     `db:"name"`
	Code          string     `db:"code"`
	Description   *string    `db:"description"`
	DegreeLevel   *string    `db:"degree_level"`
	Active        bool       `db:"is_active"`
	SyncedAt      *time.Time `db:"synced_at"`
	DateCreated   time.Time  `db:"date_created"`
	DateUpdated   time.Time  `db:"date_updated"`
}

type duplicateReviewDB struct {
	ID                     uuid.UUID  `db:"duplicate_review_id"`
	SourceConstituentID    uuid.UUID  `db:"source_constituent_id"`
	CandidateConstituentID uuid.UUID  `db:"candidate_constituent_id"`
	MatchType              string     `db:"match_type"`
	MatchScore             int        `db:"match_score"`
	MatchReason            string     `db:"match_reason"`
	Status                 string     `db:"status"`
	ResolvedBy             *uuid.UUID `db:"resolved_by"`
	ResolvedAt             *time.Time `db:"resolved_at"`
	ResolutionNote         *string    `db:"resolution_note"`
	DateCreated            time.Time  `db:"date_created"`
	DateUpdated            time.Time  `db:"date_updated"`
}

type applicationDB struct {
	ID                 uuid.UUID       `db:"application_id"`
	ConstituentID      uuid.UUID       `db:"constituent_id"`
	ProgramID          uuid.UUID       `db:"program_id"`
	AcademicTermID     uuid.UUID       `db:"academic_term_id"`
	ApplicationType    string          `db:"application_type"`
	Status             string          `db:"status"`
	KUCCPSPlacement    json.RawMessage `db:"kuccps_placement"`
	KCSEResult         json.RawMessage `db:"kcse_result"`
	AssignedReviewerID *uuid.UUID      `db:"assigned_reviewer_id"`
	SubmittedAt        *time.Time      `db:"submitted_at"`
	DateCreated        time.Time       `db:"date_created"`
	DateUpdated        time.Time       `db:"date_updated"`
}

type applicationFormTemplateDB struct {
	ID              uuid.UUID       `db:"form_template_id"`
	ProgramID       uuid.UUID       `db:"program_id"`
	AcademicTermID  uuid.UUID       `db:"academic_term_id"`
	ApplicationType string          `db:"application_type"`
	Name            string          `db:"name"`
	Description     *string         `db:"description"`
	Version         int             `db:"version"`
	RequiredFields  json.RawMessage `db:"required_fields"`
	ChecklistItems  json.RawMessage `db:"checklist_items"`
	Active          bool            `db:"is_active"`
	Priority        int             `db:"priority"`
	DateCreated     time.Time       `db:"date_created"`
	DateUpdated     time.Time       `db:"date_updated"`
}

type customFieldDefinitionDB struct {
	ID           uuid.UUID      `db:"custom_field_definition_id"`
	Owner        string         `db:"owner"`
	FieldKey     string         `db:"field_key"`
	Label        string         `db:"label"`
	Description  *string        `db:"description"`
	DataType     string         `db:"data_type"`
	Required     bool           `db:"is_required"`
	Options      dbarray.String `db:"options"`
	Validation   *string        `db:"validation"`
	Searchable   bool           `db:"is_searchable"`
	Reportable   bool           `db:"is_reportable"`
	Importable   bool           `db:"is_importable"`
	Exportable   bool           `db:"is_exportable"`
	DisplayOrder int            `db:"display_order"`
	Active       bool           `db:"is_active"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

type customFieldValueDB struct {
	ID           uuid.UUID `db:"custom_field_value_id"`
	DefinitionID uuid.UUID `db:"custom_field_definition_id"`
	Owner        string    `db:"owner"`
	OwnerID      uuid.UUID `db:"owner_id"`
	Value        string    `db:"value"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

type importBatchDB struct {
	ID                uuid.UUID       `db:"import_batch_id"`
	Source            string          `db:"source"`
	FileType          string          `db:"file_type"`
	Target            string          `db:"target"`
	Status            string          `db:"status"`
	FileName          string          `db:"file_name"`
	StorageKey        *string         `db:"storage_key"`
	UploadedByID      uuid.UUID       `db:"uploaded_by_id"`
	TotalRows         int             `db:"total_rows"`
	ValidRows         int             `db:"valid_rows"`
	InvalidRows       int             `db:"invalid_rows"`
	DuplicateRows     int             `db:"duplicate_rows"`
	FieldMapping      json.RawMessage `db:"field_mapping"`
	InvalidReportKey  *string         `db:"invalid_report_key"`
	ValidationSummary *string         `db:"validation_summary"`
	CommittedAt       *time.Time      `db:"committed_at"`
	DateCreated       time.Time       `db:"date_created"`
	DateUpdated       time.Time       `db:"date_updated"`
}

type importInvalidRowDB struct {
	ID          uuid.UUID       `db:"import_invalid_row_id"`
	BatchID     uuid.UUID       `db:"import_batch_id"`
	RowNumber   int             `db:"row_number"`
	FieldName   *string         `db:"field_name"`
	RawData     json.RawMessage `db:"raw_data"`
	ErrorCode   string          `db:"error_code"`
	ErrorDetail string          `db:"error_detail"`
	DateCreated time.Time       `db:"date_created"`
}

type applicationTransitionDB struct {
	ID            uuid.UUID       `db:"application_transition_id"`
	ApplicationID uuid.UUID       `db:"application_id"`
	FromStatus    string          `db:"from_status"`
	ToStatus      string          `db:"to_status"`
	ActorID       uuid.UUID       `db:"actor_id"`
	Reason        *string         `db:"reason"`
	Note          *string         `db:"note"`
	Metadata      json.RawMessage `db:"metadata"`
	DateCreated   time.Time       `db:"date_created"`
}

type checklistItemDB struct {
	ID            uuid.UUID `db:"checklist_item_id"`
	ApplicationID uuid.UUID `db:"application_id"`
	ItemKey       string    `db:"item_key"`
	DocumentName  string    `db:"document_name"`
	Description   *string   `db:"description"`
	Required      bool      `db:"is_required"`
	Status        string    `db:"status"`
	DisplayOrder  int       `db:"display_order"`
	DateCreated   time.Time `db:"date_created"`
	DateUpdated   time.Time `db:"date_updated"`
}

type documentDB struct {
	ID              uuid.UUID  `db:"document_id"`
	ApplicationID   uuid.UUID  `db:"application_id"`
	ChecklistItemID uuid.UUID  `db:"checklist_item_id"`
	FileName        string     `db:"file_name"`
	ContentType     string     `db:"content_type"`
	SizeBytes       int64      `db:"size_bytes"`
	StorageKey      string     `db:"storage_key"`
	Status          string     `db:"status"`
	ReviewerID      *uuid.UUID `db:"reviewer_id"`
	ReviewerNotes   *string    `db:"reviewer_notes"`
	UploadedByID    uuid.UUID  `db:"uploaded_by_id"`
	UploadedAt      time.Time  `db:"uploaded_at"`
	ReviewedAt      *time.Time `db:"reviewed_at"`
	DateCreated     time.Time  `db:"date_created"`
	DateUpdated     time.Time  `db:"date_updated"`
}

func toDBStaffProfile(bus admissionsbus.StaffProfile) staffProfileDB {
	return staffProfileDB{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Roles:       admissionsbus.AdmissionsRolesToStrings(bus.Roles),
		Active:      bus.Active,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusStaffProfile(db staffProfileDB) (admissionsbus.StaffProfile, error) {
	roles, err := admissionsbus.ParseAdmissionsRoles(db.Roles)
	if err != nil {
		return admissionsbus.StaffProfile{}, err
	}

	return admissionsbus.StaffProfile{
		ID:          db.ID,
		UserID:      db.UserID,
		Roles:       roles,
		Active:      db.Active,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}, nil
}

func toBusStaffProfiles(dbs []staffProfileDB) ([]admissionsbus.StaffProfile, error) {
	bus := make([]admissionsbus.StaffProfile, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusStaffProfile(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBApplicantProfile(bus admissionsbus.ApplicantProfile) applicantProfileDB {
	return applicantProfileDB{
		ID:            bus.ID,
		UserID:        bus.UserID,
		ConstituentID: bus.ConstituentID,
		Active:        bus.Active,
		DateCreated:   bus.DateCreated.UTC(),
		DateUpdated:   bus.DateUpdated.UTC(),
	}
}

func toBusApplicantProfile(db applicantProfileDB) admissionsbus.ApplicantProfile {
	return admissionsbus.ApplicantProfile{
		ID:            db.ID,
		UserID:        db.UserID,
		ConstituentID: db.ConstituentID,
		Active:        db.Active,
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}
}

func toBusApplicantProfiles(dbs []applicantProfileDB) []admissionsbus.ApplicantProfile {
	bus := make([]admissionsbus.ApplicantProfile, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusApplicantProfile(db)
	}

	return bus
}

func toDBLeadScoreRule(bus admissionsbus.LeadScoreRule) (leadScoreRuleDB, error) {
	criteria, err := json.Marshal(bus.Criteria)
	if err != nil {
		return leadScoreRuleDB{}, err
	}

	return leadScoreRuleDB{
		ID:          bus.ID,
		Name:        bus.Name,
		Description: bus.Description,
		Criteria:    criteria,
		Points:      bus.Points,
		Active:      bus.Active,
		Priority:    bus.Priority,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}, nil
}

func toBusLeadScoreRule(db leadScoreRuleDB) (admissionsbus.LeadScoreRule, error) {
	var criteria []admissionsbus.LeadScoreCriterion
	if err := json.Unmarshal(db.Criteria, &criteria); err != nil {
		return admissionsbus.LeadScoreRule{}, err
	}

	return admissionsbus.LeadScoreRule{
		ID:          db.ID,
		Name:        db.Name,
		Description: db.Description,
		Criteria:    criteria,
		Points:      db.Points,
		Active:      db.Active,
		Priority:    db.Priority,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}, nil
}

func toBusLeadScoreRules(dbs []leadScoreRuleDB) ([]admissionsbus.LeadScoreRule, error) {
	bus := make([]admissionsbus.LeadScoreRule, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusLeadScoreRule(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBLeadScore(bus admissionsbus.LeadScore) (leadScoreDB, error) {
	breakdown, err := json.Marshal(bus.Breakdown)
	if err != nil {
		return leadScoreDB{}, err
	}

	return leadScoreDB{
		ID:             bus.ID,
		ConstituentID:  bus.ConstituentID,
		TotalScore:     bus.TotalScore,
		Band:           bus.Band.String(),
		Breakdown:      breakdown,
		RecalculatedAt: bus.RecalculatedAt.UTC(),
		DateCreated:    bus.DateCreated.UTC(),
		DateUpdated:    bus.DateUpdated.UTC(),
	}, nil
}

func toBusLeadScore(db leadScoreDB) (admissionsbus.LeadScore, error) {
	var breakdown []admissionsbus.LeadScoreRuleResult
	if err := json.Unmarshal(db.Breakdown, &breakdown); err != nil {
		return admissionsbus.LeadScore{}, err
	}

	return admissionsbus.LeadScore{
		ID:             db.ID,
		ConstituentID:  db.ConstituentID,
		TotalScore:     db.TotalScore,
		Band:           admissionsbus.LeadScoreBand(db.Band),
		Breakdown:      breakdown,
		RecalculatedAt: db.RecalculatedAt.In(time.Local),
		DateCreated:    db.DateCreated.In(time.Local),
		DateUpdated:    db.DateUpdated.In(time.Local),
	}, nil
}

func toBusLeadScores(dbs []leadScoreDB) ([]admissionsbus.LeadScore, error) {
	bus := make([]admissionsbus.LeadScore, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusLeadScore(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBConstituent(bus admissionsbus.Constituent) constituentDB {
	return constituentDB{
		ID:                          bus.ID,
		FirstName:                   bus.FirstName,
		LastName:                    bus.LastName,
		PreferredName:               bus.PreferredName,
		MiddleName:                  bus.MiddleName,
		Suffix:                      bus.Suffix,
		DateOfBirth:                 bus.DateOfBirth.UTC(),
		PrimaryEmail:                bus.PrimaryEmail.String(),
		PrimaryPhone:                bus.PrimaryPhone,
		ExternalSISID:               bus.ExternalSISID,
		NationalID:                  bus.NationalID,
		NationalIDVerifiedAt:        utcTimePtr(bus.NationalIDVerifiedAt),
		NationalIDVerifiedByAdapter: bus.NationalIDVerifiedByAdapter,
		UPI:                         bus.UPI,
		UPIVerifiedAt:               utcTimePtr(bus.UPIVerifiedAt),
		UPIVerifiedByAdapter:        bus.UPIVerifiedByAdapter,
		KCSEIndexNumber:             bus.KCSEIndexNumber,
		KCSEIndexVerifiedAt:         utcTimePtr(bus.KCSEIndexVerifiedAt),
		KCSEIndexVerifiedByAdapter:  bus.KCSEIndexVerifiedByAdapter,
		LifecycleStage:              bus.LifecycleStage.String(),
		DuplicateStatus:             bus.DuplicateStatus.String(),
		DuplicateOfID:               bus.DuplicateOfID,
		SISSyncedAt:                 utcTimePtr(bus.SISSyncedAt),
		DateCreated:                 bus.DateCreated.UTC(),
		DateUpdated:                 bus.DateUpdated.UTC(),
	}
}

func toBusConstituent(db constituentDB) (admissionsbus.Constituent, error) {
	email, err := mail.ParseAddress(db.PrimaryEmail)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return admissionsbus.Constituent{
		ID:                          db.ID,
		FirstName:                   db.FirstName,
		LastName:                    db.LastName,
		PreferredName:               db.PreferredName,
		MiddleName:                  db.MiddleName,
		Suffix:                      db.Suffix,
		DateOfBirth:                 db.DateOfBirth.In(time.Local),
		PrimaryEmail:                *email,
		PrimaryPhone:                db.PrimaryPhone,
		ExternalSISID:               db.ExternalSISID,
		NationalID:                  db.NationalID,
		NationalIDVerifiedAt:        localTimePtr(db.NationalIDVerifiedAt),
		NationalIDVerifiedByAdapter: db.NationalIDVerifiedByAdapter,
		UPI:                         db.UPI,
		UPIVerifiedAt:               localTimePtr(db.UPIVerifiedAt),
		UPIVerifiedByAdapter:        db.UPIVerifiedByAdapter,
		KCSEIndexNumber:             db.KCSEIndexNumber,
		KCSEIndexVerifiedAt:         localTimePtr(db.KCSEIndexVerifiedAt),
		KCSEIndexVerifiedByAdapter:  db.KCSEIndexVerifiedByAdapter,
		LifecycleStage:              admissionsbus.LifecycleStage(db.LifecycleStage),
		DuplicateStatus:             admissionsbus.DuplicateStatus(db.DuplicateStatus),
		DuplicateOfID:               db.DuplicateOfID,
		SISSyncedAt:                 localTimePtr(db.SISSyncedAt),
		DateCreated:                 db.DateCreated.In(time.Local),
		DateUpdated:                 db.DateUpdated.In(time.Local),
	}, nil
}

func toBusConstituents(dbs []constituentDB) ([]admissionsbus.Constituent, error) {
	bus := make([]admissionsbus.Constituent, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusConstituent(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBInquiry(bus admissionsbus.Inquiry) inquiryDB {
	return inquiryDB{
		ID:                bus.ID,
		ConstituentID:     bus.ConstituentID,
		FirstName:         bus.FirstName,
		LastName:          bus.LastName,
		DateOfBirth:       bus.DateOfBirth.UTC(),
		PrimaryEmail:      bus.PrimaryEmail.String(),
		PrimaryPhone:      bus.PrimaryPhone,
		ProgramOfInterest: bus.ProgramOfInterest,
		TermOfInterest:    bus.TermOfInterest,
		Source:            bus.Source,
		UTMSource:         bus.UTMSource,
		UTMMedium:         bus.UTMMedium,
		UTMCampaign:       bus.UTMCampaign,
		Message:           bus.Message,
		Status:            bus.Status.String(),
		DateCreated:       bus.DateCreated.UTC(),
		DateUpdated:       bus.DateUpdated.UTC(),
	}
}

func toBusInquiry(db inquiryDB) (admissionsbus.Inquiry, error) {
	email, err := mail.ParseAddress(db.PrimaryEmail)
	if err != nil {
		return admissionsbus.Inquiry{}, err
	}

	return admissionsbus.Inquiry{
		ID:                db.ID,
		ConstituentID:     db.ConstituentID,
		FirstName:         db.FirstName,
		LastName:          db.LastName,
		DateOfBirth:       db.DateOfBirth.In(time.Local),
		PrimaryEmail:      *email,
		PrimaryPhone:      db.PrimaryPhone,
		ProgramOfInterest: db.ProgramOfInterest,
		TermOfInterest:    db.TermOfInterest,
		Source:            db.Source,
		UTMSource:         db.UTMSource,
		UTMMedium:         db.UTMMedium,
		UTMCampaign:       db.UTMCampaign,
		Message:           db.Message,
		Status:            admissionsbus.InquiryStatus(db.Status),
		DateCreated:       db.DateCreated.In(time.Local),
		DateUpdated:       db.DateUpdated.In(time.Local),
	}, nil
}

func toBusInquiries(dbs []inquiryDB) ([]admissionsbus.Inquiry, error) {
	bus := make([]admissionsbus.Inquiry, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusInquiry(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

type academicTermDB struct {
	ID                   uuid.UUID  `db:"academic_term_id"`
	ExternalSISID        string     `db:"external_sis_id"`
	Name                 string     `db:"name"`
	Code                 string     `db:"code"`
	TermType             *string    `db:"term_type"`
	StartDate            time.Time  `db:"start_date"`
	EndDate              time.Time  `db:"end_date"`
	ApplicationStartDate *time.Time `db:"application_start_date"`
	ApplicationDeadline  *time.Time `db:"application_deadline"`
	Active               bool       `db:"is_active"`
	SyncedAt             *time.Time `db:"synced_at"`
	DateCreated          time.Time  `db:"date_created"`
	DateUpdated          time.Time  `db:"date_updated"`
}

func toDBProgram(bus admissionsbus.Program) programDB {
	return programDB{
		ID:            bus.ID,
		ExternalSISID: bus.ExternalSISID,
		Name:          bus.Name,
		Code:          bus.Code,
		Description:   bus.Description,
		DegreeLevel:   bus.DegreeLevel,
		Active:        bus.Active,
		SyncedAt:      utcTimePtr(bus.SyncedAt),
		DateCreated:   bus.DateCreated.UTC(),
		DateUpdated:   bus.DateUpdated.UTC(),
	}
}

func toBusProgram(db programDB) admissionsbus.Program {
	return admissionsbus.Program{
		ID:            db.ID,
		ExternalSISID: db.ExternalSISID,
		Name:          db.Name,
		Code:          db.Code,
		Description:   db.Description,
		DegreeLevel:   db.DegreeLevel,
		Active:        db.Active,
		SyncedAt:      localTimePtr(db.SyncedAt),
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}
}

func toBusPrograms(dbs []programDB) []admissionsbus.Program {
	bus := make([]admissionsbus.Program, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusProgram(db)
	}

	return bus
}

func toDBAcademicTerm(bus admissionsbus.AcademicTerm) academicTermDB {
	return academicTermDB{
		ID:                   bus.ID,
		ExternalSISID:        bus.ExternalSISID,
		Name:                 bus.Name,
		Code:                 bus.Code,
		TermType:             bus.TermType,
		StartDate:            bus.StartDate.UTC(),
		EndDate:              bus.EndDate.UTC(),
		ApplicationStartDate: utcTimePtr(bus.ApplicationStartDate),
		ApplicationDeadline:  utcTimePtr(bus.ApplicationDeadline),
		Active:               bus.Active,
		SyncedAt:             utcTimePtr(bus.SyncedAt),
		DateCreated:          bus.DateCreated.UTC(),
		DateUpdated:          bus.DateUpdated.UTC(),
	}
}

func toBusAcademicTerm(db academicTermDB) admissionsbus.AcademicTerm {
	return admissionsbus.AcademicTerm{
		ID:                   db.ID,
		ExternalSISID:        db.ExternalSISID,
		Name:                 db.Name,
		Code:                 db.Code,
		TermType:             db.TermType,
		StartDate:            db.StartDate.In(time.Local),
		EndDate:              db.EndDate.In(time.Local),
		ApplicationStartDate: localTimePtr(db.ApplicationStartDate),
		ApplicationDeadline:  localTimePtr(db.ApplicationDeadline),
		Active:               db.Active,
		SyncedAt:             localTimePtr(db.SyncedAt),
		DateCreated:          db.DateCreated.In(time.Local),
		DateUpdated:          db.DateUpdated.In(time.Local),
	}
}

func toBusAcademicTerms(dbs []academicTermDB) []admissionsbus.AcademicTerm {
	bus := make([]admissionsbus.AcademicTerm, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusAcademicTerm(db)
	}

	return bus
}

func toDBDuplicateReview(bus admissionsbus.DuplicateReview) duplicateReviewDB {
	return duplicateReviewDB{
		ID:                     bus.ID,
		SourceConstituentID:    bus.SourceConstituentID,
		CandidateConstituentID: bus.CandidateConstituentID,
		MatchType:              bus.MatchType.String(),
		MatchScore:             bus.MatchScore,
		MatchReason:            bus.MatchReason,
		Status:                 bus.Status.String(),
		ResolvedBy:             bus.ResolvedBy,
		ResolvedAt:             utcTimePtr(bus.ResolvedAt),
		ResolutionNote:         bus.ResolutionNote,
		DateCreated:            bus.DateCreated.UTC(),
		DateUpdated:            bus.DateUpdated.UTC(),
	}
}

func toBusDuplicateReview(db duplicateReviewDB) admissionsbus.DuplicateReview {
	return admissionsbus.DuplicateReview{
		ID:                     db.ID,
		SourceConstituentID:    db.SourceConstituentID,
		CandidateConstituentID: db.CandidateConstituentID,
		MatchType:              admissionsbus.DuplicateReviewMatchType(db.MatchType),
		MatchScore:             db.MatchScore,
		MatchReason:            db.MatchReason,
		Status:                 admissionsbus.DuplicateReviewStatus(db.Status),
		ResolvedBy:             db.ResolvedBy,
		ResolvedAt:             localTimePtr(db.ResolvedAt),
		ResolutionNote:         db.ResolutionNote,
		DateCreated:            db.DateCreated.In(time.Local),
		DateUpdated:            db.DateUpdated.In(time.Local),
	}
}

func toBusDuplicateReviews(dbs []duplicateReviewDB) []admissionsbus.DuplicateReview {
	bus := make([]admissionsbus.DuplicateReview, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusDuplicateReview(db)
	}

	return bus
}

func toDBApplication(bus admissionsbus.Application) (applicationDB, error) {
	kuccpsPlacement, err := json.Marshal(bus.KUCCPSPlacement)
	if err != nil {
		return applicationDB{}, err
	}

	kcseResult, err := json.Marshal(bus.KCSEResult)
	if err != nil {
		return applicationDB{}, err
	}

	return applicationDB{
		ID:                 bus.ID,
		ConstituentID:      bus.ConstituentID,
		ProgramID:          bus.ProgramID,
		AcademicTermID:     bus.AcademicTermID,
		ApplicationType:    bus.ApplicationType.String(),
		Status:             bus.Status.String(),
		KUCCPSPlacement:    json.RawMessage(kuccpsPlacement),
		KCSEResult:         json.RawMessage(kcseResult),
		AssignedReviewerID: bus.AssignedReviewerID,
		SubmittedAt:        utcTimePtr(bus.SubmittedAt),
		DateCreated:        bus.DateCreated.UTC(),
		DateUpdated:        bus.DateUpdated.UTC(),
	}, nil
}

func toBusApplication(db applicationDB) (admissionsbus.Application, error) {
	var kuccpsPlacement *admissionsbus.KUCCPSPlacement
	if len(db.KUCCPSPlacement) > 0 {
		if err := json.Unmarshal(db.KUCCPSPlacement, &kuccpsPlacement); err != nil {
			return admissionsbus.Application{}, err
		}
	}

	var kcseResult *admissionsbus.ApplicationKCSEResult
	if len(db.KCSEResult) > 0 {
		if err := json.Unmarshal(db.KCSEResult, &kcseResult); err != nil {
			return admissionsbus.Application{}, err
		}
	}

	return admissionsbus.Application{
		ID:                 db.ID,
		ConstituentID:      db.ConstituentID,
		ProgramID:          db.ProgramID,
		AcademicTermID:     db.AcademicTermID,
		ApplicationType:    admissionsbus.ApplicationType(db.ApplicationType),
		Status:             admissionsbus.ApplicationStatus(db.Status),
		KUCCPSPlacement:    kuccpsPlacement,
		KCSEResult:         kcseResult,
		AssignedReviewerID: db.AssignedReviewerID,
		SubmittedAt:        localTimePtr(db.SubmittedAt),
		DateCreated:        db.DateCreated.In(time.Local),
		DateUpdated:        db.DateUpdated.In(time.Local),
	}, nil
}

func toBusApplications(dbs []applicationDB) ([]admissionsbus.Application, error) {
	bus := make([]admissionsbus.Application, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusApplication(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBApplicationFormTemplate(bus admissionsbus.ApplicationFormTemplate) (applicationFormTemplateDB, error) {
	requiredFields, err := json.Marshal(bus.RequiredFields)
	if err != nil {
		return applicationFormTemplateDB{}, err
	}

	checklistItems, err := json.Marshal(bus.ChecklistItems)
	if err != nil {
		return applicationFormTemplateDB{}, err
	}

	return applicationFormTemplateDB{
		ID:              bus.ID,
		ProgramID:       bus.ProgramID,
		AcademicTermID:  bus.AcademicTermID,
		ApplicationType: bus.ApplicationType.String(),
		Name:            bus.Name,
		Description:     bus.Description,
		Version:         bus.Version,
		RequiredFields:  requiredFields,
		ChecklistItems:  checklistItems,
		Active:          bus.Active,
		Priority:        bus.Priority,
		DateCreated:     bus.DateCreated.UTC(),
		DateUpdated:     bus.DateUpdated.UTC(),
	}, nil
}

func toBusApplicationFormTemplate(db applicationFormTemplateDB) (admissionsbus.ApplicationFormTemplate, error) {
	var requiredFields []admissionsbus.ApplicationFormField
	if err := json.Unmarshal(db.RequiredFields, &requiredFields); err != nil {
		return admissionsbus.ApplicationFormTemplate{}, err
	}

	var checklistItems []admissionsbus.ApplicationChecklistTemplateItem
	if err := json.Unmarshal(db.ChecklistItems, &checklistItems); err != nil {
		return admissionsbus.ApplicationFormTemplate{}, err
	}

	return admissionsbus.ApplicationFormTemplate{
		ID:              db.ID,
		ProgramID:       db.ProgramID,
		AcademicTermID:  db.AcademicTermID,
		ApplicationType: admissionsbus.ApplicationType(db.ApplicationType),
		Name:            db.Name,
		Description:     db.Description,
		Version:         db.Version,
		RequiredFields:  requiredFields,
		ChecklistItems:  checklistItems,
		Active:          db.Active,
		Priority:        db.Priority,
		DateCreated:     db.DateCreated.In(time.Local),
		DateUpdated:     db.DateUpdated.In(time.Local),
	}, nil
}

func toBusApplicationFormTemplates(dbs []applicationFormTemplateDB) ([]admissionsbus.ApplicationFormTemplate, error) {
	bus := make([]admissionsbus.ApplicationFormTemplate, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusApplicationFormTemplate(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBCustomFieldDefinition(bus admissionsbus.CustomFieldDefinition) customFieldDefinitionDB {
	return customFieldDefinitionDB{
		ID:           bus.ID,
		Owner:        bus.Owner.String(),
		FieldKey:     bus.FieldKey,
		Label:        bus.Label,
		Description:  bus.Description,
		DataType:     bus.DataType.String(),
		Required:     bus.Required,
		Options:      bus.Options,
		Validation:   bus.Validation,
		Searchable:   bus.Searchable,
		Reportable:   bus.Reportable,
		Importable:   bus.Importable,
		Exportable:   bus.Exportable,
		DisplayOrder: bus.DisplayOrder,
		Active:       bus.Active,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusCustomFieldDefinition(db customFieldDefinitionDB) admissionsbus.CustomFieldDefinition {
	return admissionsbus.CustomFieldDefinition{
		ID:           db.ID,
		Owner:        admissionsbus.CustomFieldOwner(db.Owner),
		FieldKey:     db.FieldKey,
		Label:        db.Label,
		Description:  db.Description,
		DataType:     admissionsbus.CustomFieldDataType(db.DataType),
		Required:     db.Required,
		Options:      []string(db.Options),
		Validation:   db.Validation,
		Searchable:   db.Searchable,
		Reportable:   db.Reportable,
		Importable:   db.Importable,
		Exportable:   db.Exportable,
		DisplayOrder: db.DisplayOrder,
		Active:       db.Active,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}
}

func toBusCustomFieldDefinitions(dbs []customFieldDefinitionDB) []admissionsbus.CustomFieldDefinition {
	bus := make([]admissionsbus.CustomFieldDefinition, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusCustomFieldDefinition(db)
	}

	return bus
}

func toDBCustomFieldValue(bus admissionsbus.CustomFieldValue) customFieldValueDB {
	return customFieldValueDB{
		ID:           bus.ID,
		DefinitionID: bus.DefinitionID,
		Owner:        bus.Owner.String(),
		OwnerID:      bus.OwnerID,
		Value:        bus.Value,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusCustomFieldValue(db customFieldValueDB) admissionsbus.CustomFieldValue {
	return admissionsbus.CustomFieldValue{
		ID:           db.ID,
		DefinitionID: db.DefinitionID,
		Owner:        admissionsbus.CustomFieldOwner(db.Owner),
		OwnerID:      db.OwnerID,
		Value:        db.Value,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}
}

func toBusCustomFieldValues(dbs []customFieldValueDB) []admissionsbus.CustomFieldValue {
	bus := make([]admissionsbus.CustomFieldValue, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusCustomFieldValue(db)
	}

	return bus
}

func toDBImportBatch(bus admissionsbus.ImportBatch) (importBatchDB, error) {
	fieldMapping, err := json.Marshal(bus.FieldMapping)
	if err != nil {
		return importBatchDB{}, err
	}

	return importBatchDB{
		ID:                bus.ID,
		Source:            bus.Source.String(),
		FileType:          bus.FileType.String(),
		Target:            bus.Target.String(),
		Status:            bus.Status.String(),
		FileName:          bus.FileName,
		StorageKey:        bus.StorageKey,
		UploadedByID:      bus.UploadedByID,
		TotalRows:         bus.TotalRows,
		ValidRows:         bus.ValidRows,
		InvalidRows:       bus.InvalidRows,
		DuplicateRows:     bus.DuplicateRows,
		FieldMapping:      fieldMapping,
		InvalidReportKey:  bus.InvalidReportKey,
		ValidationSummary: bus.ValidationSummary,
		CommittedAt:       utcTimePtr(bus.CommittedAt),
		DateCreated:       bus.DateCreated.UTC(),
		DateUpdated:       bus.DateUpdated.UTC(),
	}, nil
}

func toBusImportBatch(db importBatchDB) (admissionsbus.ImportBatch, error) {
	var fieldMapping map[string]string
	if err := json.Unmarshal(db.FieldMapping, &fieldMapping); err != nil {
		return admissionsbus.ImportBatch{}, err
	}

	return admissionsbus.ImportBatch{
		ID:                db.ID,
		Source:            admissionsbus.ImportSource(db.Source),
		FileType:          admissionsbus.ImportFileType(db.FileType),
		Target:            admissionsbus.ImportTarget(db.Target),
		Status:            admissionsbus.ImportBatchStatus(db.Status),
		FileName:          db.FileName,
		StorageKey:        db.StorageKey,
		UploadedByID:      db.UploadedByID,
		TotalRows:         db.TotalRows,
		ValidRows:         db.ValidRows,
		InvalidRows:       db.InvalidRows,
		DuplicateRows:     db.DuplicateRows,
		FieldMapping:      fieldMapping,
		InvalidReportKey:  db.InvalidReportKey,
		ValidationSummary: db.ValidationSummary,
		CommittedAt:       localTimePtr(db.CommittedAt),
		DateCreated:       db.DateCreated.In(time.Local),
		DateUpdated:       db.DateUpdated.In(time.Local),
	}, nil
}

func toBusImportBatches(dbs []importBatchDB) ([]admissionsbus.ImportBatch, error) {
	bus := make([]admissionsbus.ImportBatch, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusImportBatch(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBImportInvalidRow(bus admissionsbus.ImportInvalidRow) (importInvalidRowDB, error) {
	rawData, err := json.Marshal(bus.RawData)
	if err != nil {
		return importInvalidRowDB{}, err
	}

	return importInvalidRowDB{
		ID:          bus.ID,
		BatchID:     bus.BatchID,
		RowNumber:   bus.RowNumber,
		FieldName:   bus.FieldName,
		RawData:     rawData,
		ErrorCode:   bus.ErrorCode,
		ErrorDetail: bus.ErrorDetail,
		DateCreated: bus.DateCreated.UTC(),
	}, nil
}

func toDBImportInvalidRows(bus []admissionsbus.ImportInvalidRow) ([]importInvalidRowDB, error) {
	dbs := make([]importInvalidRowDB, len(bus))
	for i, row := range bus {
		var err error
		dbs[i], err = toDBImportInvalidRow(row)
		if err != nil {
			return nil, err
		}
	}

	return dbs, nil
}

func toBusImportInvalidRow(db importInvalidRowDB) (admissionsbus.ImportInvalidRow, error) {
	var rawData map[string]string
	if err := json.Unmarshal(db.RawData, &rawData); err != nil {
		return admissionsbus.ImportInvalidRow{}, err
	}

	return admissionsbus.ImportInvalidRow{
		ID:          db.ID,
		BatchID:     db.BatchID,
		RowNumber:   db.RowNumber,
		FieldName:   db.FieldName,
		RawData:     rawData,
		ErrorCode:   db.ErrorCode,
		ErrorDetail: db.ErrorDetail,
		DateCreated: db.DateCreated.In(time.Local),
	}, nil
}

func toBusImportInvalidRows(dbs []importInvalidRowDB) ([]admissionsbus.ImportInvalidRow, error) {
	bus := make([]admissionsbus.ImportInvalidRow, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusImportInvalidRow(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBApplicationTransition(bus admissionsbus.ApplicationTransition) applicationTransitionDB {
	return applicationTransitionDB{
		ID:            bus.ID,
		ApplicationID: bus.ApplicationID,
		FromStatus:    bus.FromStatus.String(),
		ToStatus:      bus.ToStatus.String(),
		ActorID:       bus.ActorID,
		Reason:        bus.Reason,
		Note:          bus.Note,
		Metadata:      json.RawMessage(bus.Metadata),
		DateCreated:   bus.DateCreated.UTC(),
	}
}

func toBusApplicationTransition(db applicationTransitionDB) admissionsbus.ApplicationTransition {
	return admissionsbus.ApplicationTransition{
		ID:            db.ID,
		ApplicationID: db.ApplicationID,
		FromStatus:    admissionsbus.ApplicationStatus(db.FromStatus),
		ToStatus:      admissionsbus.ApplicationStatus(db.ToStatus),
		ActorID:       db.ActorID,
		Reason:        db.Reason,
		Note:          db.Note,
		Metadata:      []byte(db.Metadata),
		DateCreated:   db.DateCreated.In(time.Local),
	}
}

func toBusApplicationTransitions(dbs []applicationTransitionDB) []admissionsbus.ApplicationTransition {
	bus := make([]admissionsbus.ApplicationTransition, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusApplicationTransition(db)
	}

	return bus
}

func toDBChecklistItem(bus admissionsbus.ChecklistItem) checklistItemDB {
	return checklistItemDB{
		ID:            bus.ID,
		ApplicationID: bus.ApplicationID,
		ItemKey:       bus.ItemKey,
		DocumentName:  bus.DocumentName,
		Description:   bus.Description,
		Required:      bus.Required,
		Status:        bus.Status.String(),
		DisplayOrder:  bus.DisplayOrder,
		DateCreated:   bus.DateCreated.UTC(),
		DateUpdated:   bus.DateUpdated.UTC(),
	}
}

func toBusChecklistItem(db checklistItemDB) admissionsbus.ChecklistItem {
	return admissionsbus.ChecklistItem{
		ID:            db.ID,
		ApplicationID: db.ApplicationID,
		ItemKey:       db.ItemKey,
		DocumentName:  db.DocumentName,
		Description:   db.Description,
		Required:      db.Required,
		Status:        admissionsbus.DocumentStatus(db.Status),
		DisplayOrder:  db.DisplayOrder,
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}
}

func toBusChecklistItems(dbs []checklistItemDB) []admissionsbus.ChecklistItem {
	bus := make([]admissionsbus.ChecklistItem, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusChecklistItem(db)
	}

	return bus
}

func toDBDocument(bus admissionsbus.Document) documentDB {
	return documentDB{
		ID:              bus.ID,
		ApplicationID:   bus.ApplicationID,
		ChecklistItemID: bus.ChecklistItemID,
		FileName:        bus.FileName,
		ContentType:     bus.ContentType,
		SizeBytes:       bus.SizeBytes,
		StorageKey:      bus.StorageKey,
		Status:          bus.Status.String(),
		ReviewerID:      bus.ReviewerID,
		ReviewerNotes:   bus.ReviewerNotes,
		UploadedByID:    bus.UploadedByID,
		UploadedAt:      bus.UploadedAt.UTC(),
		ReviewedAt:      utcTimePtr(bus.ReviewedAt),
		DateCreated:     bus.DateCreated.UTC(),
		DateUpdated:     bus.DateUpdated.UTC(),
	}
}

func toBusDocument(db documentDB) admissionsbus.Document {
	return admissionsbus.Document{
		ID:              db.ID,
		ApplicationID:   db.ApplicationID,
		ChecklistItemID: db.ChecklistItemID,
		FileName:        db.FileName,
		ContentType:     db.ContentType,
		SizeBytes:       db.SizeBytes,
		StorageKey:      db.StorageKey,
		Status:          admissionsbus.DocumentStatus(db.Status),
		ReviewerID:      db.ReviewerID,
		ReviewerNotes:   db.ReviewerNotes,
		UploadedByID:    db.UploadedByID,
		UploadedAt:      db.UploadedAt.In(time.Local),
		ReviewedAt:      localTimePtr(db.ReviewedAt),
		DateCreated:     db.DateCreated.In(time.Local),
		DateUpdated:     db.DateUpdated.In(time.Local),
	}
}

func toBusDocuments(dbs []documentDB) []admissionsbus.Document {
	bus := make([]admissionsbus.Document, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusDocument(db)
	}

	return bus
}

func utcTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	utc := t.UTC()
	return &utc
}

func localTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	local := t.In(time.Local)
	return &local
}
