package admissionsbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// AdmissionsRole represents a context-specific role inside the admissions CRM.
type AdmissionsRole string

// Set of valid admissions CRM roles.
const (
	AdmissionsRoleAdmin               AdmissionsRole = "ADMISSIONS_ADMIN"
	AdmissionsRoleRecruiter           AdmissionsRole = "RECRUITER"
	AdmissionsRoleApplicationReviewer AdmissionsRole = "APPLICATION_REVIEWER"
	AdmissionsRoleMarketingManager    AdmissionsRole = "MARKETING_MANAGER"
	AdmissionsRoleEventManager        AdmissionsRole = "EVENT_MANAGER"
	AdmissionsRoleReportViewer        AdmissionsRole = "REPORT_VIEWER"
	AdmissionsRoleApplicant           AdmissionsRole = "APPLICANT"
)

// String returns the admissions role as a string.
func (role AdmissionsRole) String() string {
	return string(role)
}

// AdmissionsPermission represents an action permission inside the admissions CRM.
type AdmissionsPermission string

// Set of action-based admissions permissions.
const (
	AdmissionsPermissionRead               AdmissionsPermission = "admissions:read"
	AdmissionsPermissionManageConstituents AdmissionsPermission = "admissions:manage_constituents"
	AdmissionsPermissionManageApplications AdmissionsPermission = "admissions:manage_applications"
	AdmissionsPermissionManageEvents       AdmissionsPermission = "admissions:manage_events"
	AdmissionsPermissionReviewApplications AdmissionsPermission = "admissions:review_applications"
	AdmissionsPermissionResolveDuplicates  AdmissionsPermission = "admissions:resolve_duplicates"
	AdmissionsPermissionManageReferences   AdmissionsPermission = "admissions:manage_references"
	AdmissionsPermissionManageStaff        AdmissionsPermission = "admissions:manage_staff"
	AdmissionsPermissionManageLeadScoring  AdmissionsPermission = "admissions:manage_lead_scoring"
)

// String returns the admissions permission as a string.
func (permission AdmissionsPermission) String() string {
	return string(permission)
}

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

// NotificationChannel identifies a constituent communication channel.
type NotificationChannel string

// Set of supported constituent notification channels.
const (
	NotificationChannelSMS      NotificationChannel = "SMS"
	NotificationChannelWhatsApp NotificationChannel = "WHATSAPP"
	NotificationChannelEmail    NotificationChannel = "EMAIL"
)

// String returns the notification channel as a string.
func (channel NotificationChannel) String() string {
	return string(channel)
}

// NotificationPreferences captures constituent channel consent and priority.
type NotificationPreferences struct {
	SMSOptIn      bool
	WhatsAppOptIn bool
	EmailOptIn    bool
	Priority      []NotificationChannel
}

// DuplicateReviewStatus represents staff workflow state for a potential duplicate.
type DuplicateReviewStatus string

// Set of valid duplicate review statuses.
const (
	DuplicateReviewStatusPending  DuplicateReviewStatus = "PENDING"
	DuplicateReviewStatusLinked   DuplicateReviewStatus = "LINKED"
	DuplicateReviewStatusMerged   DuplicateReviewStatus = "MERGED"
	DuplicateReviewStatusRejected DuplicateReviewStatus = "REJECTED"
	DuplicateReviewStatusDeferred DuplicateReviewStatus = "DEFERRED"
)

// String returns the duplicate review status as a string.
func (status DuplicateReviewStatus) String() string {
	return string(status)
}

// DuplicateReviewMatchType represents the confidence class for a duplicate match.
type DuplicateReviewMatchType string

// Set of valid duplicate review match types.
const (
	DuplicateReviewMatchTypeExact DuplicateReviewMatchType = "EXACT"
	DuplicateReviewMatchTypeFuzzy DuplicateReviewMatchType = "FUZZY"
)

// String returns the duplicate review match type as a string.
func (matchType DuplicateReviewMatchType) String() string {
	return string(matchType)
}

// DuplicateReviewResolution represents a staff action taken on a duplicate review.
type DuplicateReviewResolution string

// Set of valid duplicate review resolution actions.
const (
	DuplicateReviewResolutionLink   DuplicateReviewResolution = "LINK"
	DuplicateReviewResolutionMerge  DuplicateReviewResolution = "MERGE"
	DuplicateReviewResolutionReject DuplicateReviewResolution = "REJECT"
	DuplicateReviewResolutionDefer  DuplicateReviewResolution = "DEFER"
)

// String returns the duplicate review resolution as a string.
func (resolution DuplicateReviewResolution) String() string {
	return string(resolution)
}

// Constituent is the durable person identity root for admissions workflows.
type Constituent struct {
	ID                          uuid.UUID
	FirstName                   string
	LastName                    string
	PreferredName               *string
	MiddleName                  *string
	Suffix                      *string
	DateOfBirth                 time.Time
	PrimaryEmail                mail.Address
	PrimaryPhone                string
	ExternalSISID               *string
	NationalID                  *string
	NationalIDVerifiedAt        *time.Time
	NationalIDVerifiedByAdapter *string
	UPI                         *string
	UPIVerifiedAt               *time.Time
	UPIVerifiedByAdapter        *string
	KCSEIndexNumber             *string
	KCSEIndexVerifiedAt         *time.Time
	KCSEIndexVerifiedByAdapter  *string
	LifecycleStage              LifecycleStage
	DuplicateStatus             DuplicateStatus
	DuplicateOfID               *uuid.UUID
	NotificationPreferences     NotificationPreferences
	SISSyncedAt                 *time.Time
	DateCreated                 time.Time
	DateUpdated                 time.Time
}

// NewConstituent is what we require from clients when adding a Constituent.
type NewConstituent struct {
	FirstName                   string
	LastName                    string
	PreferredName               *string
	MiddleName                  *string
	Suffix                      *string
	DateOfBirth                 time.Time
	PrimaryEmail                mail.Address
	PrimaryPhone                string
	ExternalSISID               *string
	NationalID                  *string
	NationalIDVerifiedAt        *time.Time
	NationalIDVerifiedByAdapter *string
	UPI                         *string
	UPIVerifiedAt               *time.Time
	UPIVerifiedByAdapter        *string
	KCSEIndexNumber             *string
	KCSEIndexVerifiedAt         *time.Time
	KCSEIndexVerifiedByAdapter  *string
	LifecycleStage              LifecycleStage
	DuplicateStatus             DuplicateStatus
	DuplicateOfID               *uuid.UUID
	NotificationPreferences     *NotificationPreferences
	SISSyncedAt                 *time.Time
}

// UpdateConstituent defines what information may be provided to modify a Constituent.
type UpdateConstituent struct {
	PreferredName               *string
	MiddleName                  *string
	Suffix                      *string
	PrimaryEmail                *mail.Address
	PrimaryPhone                *string
	NationalID                  *string
	NationalIDVerifiedAt        *time.Time
	NationalIDVerifiedByAdapter *string
	UPI                         *string
	UPIVerifiedAt               *time.Time
	UPIVerifiedByAdapter        *string
	KCSEIndexNumber             *string
	KCSEIndexVerifiedAt         *time.Time
	KCSEIndexVerifiedByAdapter  *string
	LifecycleStage              *LifecycleStage
	DuplicateStatus             *DuplicateStatus
	DuplicateOfID               *uuid.UUID
	NotificationPreferences     *NotificationPreferences
	SISSyncedAt                 *time.Time
}

// InquiryStatus represents staff follow-up state for an inquiry.
type InquiryStatus string

// Set of valid inquiry statuses.
const (
	InquiryStatusNew       InquiryStatus = "NEW"
	InquiryStatusContacted InquiryStatus = "CONTACTED"
	InquiryStatusConverted InquiryStatus = "CONVERTED"
	InquiryStatusClosed    InquiryStatus = "CLOSED"
)

// String returns the inquiry status as a string.
func (status InquiryStatus) String() string {
	return string(status)
}

// Inquiry captures pre-application interest in the school.
type Inquiry struct {
	ID                uuid.UUID
	ConstituentID     uuid.UUID
	FirstName         string
	LastName          string
	DateOfBirth       time.Time
	PrimaryEmail      mail.Address
	PrimaryPhone      string
	ProgramOfInterest *uuid.UUID
	TermOfInterest    *uuid.UUID
	Source            string
	UTMSource         *string
	UTMMedium         *string
	UTMCampaign       *string
	Message           *string
	Status            InquiryStatus
	DateCreated       time.Time
	DateUpdated       time.Time
}

// NewInquiry is what we require from anonymous prospects when submitting an inquiry.
type NewInquiry struct {
	FirstName         string
	LastName          string
	DateOfBirth       time.Time
	PrimaryEmail      mail.Address
	PrimaryPhone      string
	ProgramOfInterest *uuid.UUID
	TermOfInterest    *uuid.UUID
	Source            string
	UTMSource         *string
	UTMMedium         *string
	UTMCampaign       *string
	Message           *string
}

// ApplicationType represents the admissions category for an application.
type ApplicationType string

// Set of valid application types.
const (
	ApplicationTypeKUCCPSPlacement        ApplicationType = "KUCCPS_PLACEMENT"
	ApplicationTypeSelfSponsoredUndergrad ApplicationType = "SELF_SPONSORED_UNDERGRAD"
	ApplicationTypeDiploma                ApplicationType = "DIPLOMA"
	ApplicationTypeMasters                ApplicationType = "MASTERS"
	ApplicationTypePhD                    ApplicationType = "PHD"
	ApplicationTypeTVET                   ApplicationType = "TVET"
	ApplicationTypeBridging               ApplicationType = "BRIDGING"
	ApplicationTypeCertificate            ApplicationType = "CERTIFICATE"
)

// String returns the application type as a string.
func (applicationType ApplicationType) String() string {
	return string(applicationType)
}

// ApplicationStatus represents the workflow state of an admissions application.
type ApplicationStatus string

// Set of valid application statuses.
const (
	ApplicationStatusDraft             ApplicationStatus = "DRAFT"
	ApplicationStatusSubmitted         ApplicationStatus = "SUBMITTED"
	ApplicationStatusAwaitingDocuments ApplicationStatus = "AWAITING_DOCUMENTS"
	ApplicationStatusReadyForReview    ApplicationStatus = "READY_FOR_REVIEW"
	ApplicationStatusInReview          ApplicationStatus = "IN_REVIEW"
	ApplicationStatusDecisionPending   ApplicationStatus = "DECISION_PENDING"
	ApplicationStatusAdmitted          ApplicationStatus = "ADMITTED"
	ApplicationStatusDenied            ApplicationStatus = "DENIED"
	ApplicationStatusWaitlisted        ApplicationStatus = "WAITLISTED"
	ApplicationStatusDeferred          ApplicationStatus = "DEFERRED"
	ApplicationStatusWithdrawn         ApplicationStatus = "WITHDRAWN"
	ApplicationStatusEnrolled          ApplicationStatus = "ENROLLED"
)

// String returns the application status as a string.
func (status ApplicationStatus) String() string {
	return string(status)
}

// EventType represents the kind of admissions engagement event.
type EventType string

// Set of valid admissions event types.
const (
	EventTypeOpenDay    EventType = "open-day"
	EventTypeWebinar    EventType = "webinar"
	EventTypeInfoSession EventType = "info-session"
	EventTypeCampusTour EventType = "campus-tour"
	EventTypeFair       EventType = "fair"
)

// String returns the event type as a string.
func (eventType EventType) String() string {
	return string(eventType)
}

// EventStatus represents the workflow state of an admissions event.
type EventStatus string

// Set of valid admissions event statuses.
const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusUpcoming  EventStatus = "upcoming"
	EventStatusLive      EventStatus = "live"
	EventStatusCompleted EventStatus = "completed"
	EventStatusCancelled EventStatus = "cancelled"
)

// String returns the event status as a string.
func (status EventStatus) String() string {
	return string(status)
}

// EventRegistrationStatus represents the state of one event registration.
type EventRegistrationStatus string

const (
	EventRegistrationStatusRegistered EventRegistrationStatus = "registered"
	EventRegistrationStatusCheckedIn  EventRegistrationStatus = "checked-in"
	EventRegistrationStatusCancelled  EventRegistrationStatus = "cancelled"
)

func (status EventRegistrationStatus) String() string {
	return string(status)
}

// EventRegistrationMatchStatus represents CRM matching confidence for a registration.
type EventRegistrationMatchStatus string

const (
	EventRegistrationMatchStatusMatched      EventRegistrationMatchStatus = "matched"
	EventRegistrationMatchStatusNewProspect  EventRegistrationMatchStatus = "new-prospect"
	EventRegistrationMatchStatusNeedsReview  EventRegistrationMatchStatus = "needs-review"
)

func (status EventRegistrationMatchStatus) String() string {
	return string(status)
}

// EventRegistrationSource represents how a prospect registered for an event.
type EventRegistrationSource string

const (
	EventRegistrationSourcePortal   EventRegistrationSource = "portal"
	EventRegistrationSourceStaff    EventRegistrationSource = "staff"
	EventRegistrationSourceCampaign EventRegistrationSource = "campaign"
)

func (source EventRegistrationSource) String() string {
	return string(source)
}

// Event represents an admissions engagement event.
type Event struct {
	ID                        uuid.UUID
	Title                     string
	Type                      EventType
	Status                    EventStatus
	Description               string
	StartTime                 time.Time
	EndTime                   time.Time
	Location                  string
	IsVirtual                 bool
	Capacity                  int
	RegistrationDeadline      *time.Time
	AutoConfirmationEnabled   bool
	AutoReminderEnabled       bool
	DateCreated               time.Time
	DateUpdated               time.Time
}

// EventRegistration represents one registration row for an admissions event.
type EventRegistration struct {
	ID           uuid.UUID
	EventID      uuid.UUID
	ConstituentID *uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Phone        *string
	Status       EventRegistrationStatus
	MatchStatus  EventRegistrationMatchStatus
	Source       EventRegistrationSource
	RegisteredAt time.Time
	CheckedInAt  *time.Time
	CheckedInByID *uuid.UUID
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewEventRegistration is what we require to register an attendee for an event.
type NewEventRegistration struct {
	EventID       uuid.UUID
	ConstituentID *uuid.UUID
	FirstName     string
	LastName      string
	Email         string
	Phone         *string
	Source        EventRegistrationSource
	MatchStatus   EventRegistrationMatchStatus
}

// NewEventCheckIn is what we require to record a staff check-in.
type NewEventCheckIn struct {
	RegistrationID uuid.UUID
	CheckedInByID  uuid.UUID
}

// Application represents a constituent's program application for a term.
type Application struct {
	ID                 uuid.UUID
	ConstituentID      uuid.UUID
	ProgramID          uuid.UUID
	AcademicTermID     uuid.UUID
	ApplicationType    ApplicationType
	Status             ApplicationStatus
	KUCCPSPlacement    *KUCCPSPlacement
	KCSEResult         *ApplicationKCSEResult
	AssignedReviewerID *uuid.UUID
	SubmittedAt        *time.Time
	DateCreated        time.Time
	DateUpdated        time.Time
}

// KUCCPSPlacement captures a normalized KUCCPS placement snapshot on an application.
type KUCCPSPlacement struct {
	PlacementID        string
	AdmissionNumber    *string
	InstitutionCode    string
	ProgrammeCode      string
	ProgrammeName      string
	PlacementYear      int
	ClusterCode        *string
	ClusterPoints      *float64
	WeightedPointsNote *string
}

// ApplicationKCSESubject stores one KCSE subject grade snapshot on an application.
type ApplicationKCSESubject struct {
	SubjectCode string
	Grade       string
	Points      int
}

// ApplicationKCSEResult stores the KCSE result snapshot submitted with an application.
type ApplicationKCSEResult struct {
	IndexNumber string
	ExamYear    int
	Subjects    []ApplicationKCSESubject
	MeanGrade   string
	MeanPoints  int
}

// ApplicationFormField defines a configurable, non-core application form field.
type ApplicationFormField struct {
	FieldName    string
	FieldType    string
	Required     bool
	DisplayOrder int
	Validation   *string
}

// CustomFieldOwner represents the admissions aggregate a custom field belongs to.
type CustomFieldOwner string

// Set of aggregates that may own custom fields.
const (
	CustomFieldOwnerConstituent CustomFieldOwner = "CONSTITUENT"
	CustomFieldOwnerApplication CustomFieldOwner = "APPLICATION"
)

// String returns the custom field owner as a string.
func (owner CustomFieldOwner) String() string {
	return string(owner)
}

// CustomFieldDataType represents the validated storage type for custom field values.
type CustomFieldDataType string

// Set of supported custom field data types.
const (
	CustomFieldDataTypeText     CustomFieldDataType = "TEXT"
	CustomFieldDataTypeTextarea CustomFieldDataType = "TEXTAREA"
	CustomFieldDataTypeNumber   CustomFieldDataType = "NUMBER"
	CustomFieldDataTypeDate     CustomFieldDataType = "DATE"
	CustomFieldDataTypeSelect   CustomFieldDataType = "SELECT"
	CustomFieldDataTypeBoolean  CustomFieldDataType = "BOOLEAN"
)

// String returns the custom field data type as a string.
func (dataType CustomFieldDataType) String() string {
	return string(dataType)
}

// CustomFieldDefinition describes one user-defined admissions field without making core fields dynamic.
type CustomFieldDefinition struct {
	ID           uuid.UUID
	Owner        CustomFieldOwner
	FieldKey     string
	Label        string
	Description  *string
	DataType     CustomFieldDataType
	Required     bool
	Options      []string
	Validation   *string
	Searchable   bool
	Reportable   bool
	Importable   bool
	Exportable   bool
	DisplayOrder int
	Active       bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewCustomFieldDefinition is what we require to create or update a custom field definition.
type NewCustomFieldDefinition struct {
	Owner        CustomFieldOwner
	FieldKey     string
	Label        string
	Description  *string
	DataType     CustomFieldDataType
	Required     bool
	Options      []string
	Validation   *string
	Searchable   bool
	Reportable   bool
	Importable   bool
	Exportable   bool
	DisplayOrder int
	Active       bool
}

// CustomFieldValue stores one validated custom field value for a constituent or application.
type CustomFieldValue struct {
	ID           uuid.UUID
	DefinitionID uuid.UUID
	Owner        CustomFieldOwner
	OwnerID      uuid.UUID
	Value        string
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewCustomFieldValue is what we require to set one custom field value.
type NewCustomFieldValue struct {
	DefinitionID uuid.UUID
	Owner        CustomFieldOwner
	OwnerID      uuid.UUID
	Value        string
}

// ApplicationChecklistTemplateItem defines a document/checklist requirement attached to a form template.
type ApplicationChecklistTemplateItem struct {
	ItemKey      string
	DocumentName string
	Description  *string
	Required     bool
	DisplayOrder int
}

// ApplicationFormTemplate defines configurable form requirements for a program, term, and application type.
type ApplicationFormTemplate struct {
	ID              uuid.UUID
	ProgramID       uuid.UUID
	AcademicTermID  uuid.UUID
	ApplicationType ApplicationType
	Name            string
	Description     *string
	Version         int
	RequiredFields  []ApplicationFormField
	ChecklistItems  []ApplicationChecklistTemplateItem
	Active          bool
	Priority        int
	DateCreated     time.Time
	DateUpdated     time.Time
}

// NewApplicationFormTemplate is what we require to create or update a form template.
type NewApplicationFormTemplate struct {
	ProgramID       uuid.UUID
	AcademicTermID  uuid.UUID
	ApplicationType ApplicationType
	Name            string
	Description     *string
	RequiredFields  []ApplicationFormField
	ChecklistItems  []ApplicationChecklistTemplateItem
	Active          bool
	Priority        int
}

// ApplicationTransition records an immutable application status transition.
type ApplicationTransition struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	FromStatus    ApplicationStatus
	ToStatus      ApplicationStatus
	ActorID       uuid.UUID
	Reason        *string
	Note          *string
	Metadata      []byte
	DateCreated   time.Time
}

// LeadScoreBand represents an explainable score tier for a constituent.
type LeadScoreBand string

// Set of derived lead score bands.
const (
	LeadScoreBandCold         LeadScoreBand = "COLD"
	LeadScoreBandWarm         LeadScoreBand = "WARM"
	LeadScoreBandHot          LeadScoreBand = "HOT"
	LeadScoreBandReadyToApply LeadScoreBand = "READY_TO_APPLY"
)

// String returns the lead score band as a string.
func (band LeadScoreBand) String() string {
	return string(band)
}

// LeadScoreCriterionField represents a supported rule criterion field.
type LeadScoreCriterionField string

// Set of supported lead score criterion fields.
const (
	LeadScoreCriterionFieldLifecycleStage    LeadScoreCriterionField = "lifecycle_stage"
	LeadScoreCriterionFieldApplicationType   LeadScoreCriterionField = "application_type"
	LeadScoreCriterionFieldApplicationStatus LeadScoreCriterionField = "application_status"
	LeadScoreCriterionFieldProgramID         LeadScoreCriterionField = "program_id"
	LeadScoreCriterionFieldAcademicTermID    LeadScoreCriterionField = "academic_term_id"
)

// String returns the criterion field as a string.
func (field LeadScoreCriterionField) String() string {
	return string(field)
}

// LeadScoreCriterionOperator represents a supported rule comparison operator.
type LeadScoreCriterionOperator string

// Set of supported lead score criterion operators.
const (
	LeadScoreCriterionOperatorEquals LeadScoreCriterionOperator = "EQ"
	LeadScoreCriterionOperatorIn     LeadScoreCriterionOperator = "IN"
)

// String returns the criterion operator as a string.
func (operator LeadScoreCriterionOperator) String() string {
	return string(operator)
}

// LeadScoreCriterion is one condition in a lead score rule.
type LeadScoreCriterion struct {
	Field    LeadScoreCriterionField
	Operator LeadScoreCriterionOperator
	Values   []string
}

// LeadScoreRule defines an enabled/disabled rule that contributes points to a constituent score.
type LeadScoreRule struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Criteria    []LeadScoreCriterion
	Points      int
	Active      bool
	Priority    int
	DateCreated time.Time
	DateUpdated time.Time
}

// NewLeadScoreRule is what we require to create or update a lead score rule.
type NewLeadScoreRule struct {
	Name        string
	Description *string
	Criteria    []LeadScoreCriterion
	Points      int
	Active      bool
	Priority    int
}

// LeadScoreRuleResult explains how one rule contributed to a lead score.
type LeadScoreRuleResult struct {
	RuleID  uuid.UUID
	Name    string
	Points  int
	Matched bool
	Reason  string
}

// LeadScore records the latest explainable score for a constituent.
type LeadScore struct {
	ID             uuid.UUID
	ConstituentID  uuid.UUID
	TotalScore     int
	Band           LeadScoreBand
	Breakdown      []LeadScoreRuleResult
	RecalculatedAt time.Time
	DateCreated    time.Time
	DateUpdated    time.Time
}

// StaffProfile connects an identity user to admissions-specific staff roles.
type StaffProfile struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Roles       []AdmissionsRole
	Active      bool
	DateCreated time.Time
	DateUpdated time.Time
}

// NewStaffProfile is what we require to create an admissions staff profile.
type NewStaffProfile struct {
	UserID uuid.UUID
	Roles  []AdmissionsRole
	Active bool
}

// ApplicantProfile connects an authenticated identity user to a constituent record.
type ApplicantProfile struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ConstituentID uuid.UUID
	Active        bool
	DateCreated   time.Time
	DateUpdated   time.Time
}

// NewApplicantProfile is what we require to create an admissions applicant profile.
type NewApplicantProfile struct {
	UserID        uuid.UUID
	ConstituentID uuid.UUID
	Active        bool
}

// NewApplicationTransition is what we require to change an application status.
type NewApplicationTransition struct {
	ToStatus ApplicationStatus
	ActorID  uuid.UUID
	Reason   *string
	Note     *string
	Metadata []byte
}

// NewApplication is what we require from clients when adding an Application.
type NewApplication struct {
	ConstituentID      uuid.UUID
	ProgramID          uuid.UUID
	AcademicTermID     uuid.UUID
	ApplicationType    ApplicationType
	KUCCPSPlacement    *KUCCPSPlacement
	KCSEResult         *ApplicationKCSEResult
	AssignedReviewerID *uuid.UUID
}

// DocumentStatus represents the review state of applicant-submitted evidence.
type DocumentStatus string

// Set of valid document statuses.
const (
	DocumentStatusUploaded      DocumentStatus = "UPLOADED"
	DocumentStatusPendingReview DocumentStatus = "PENDING_REVIEW"
	DocumentStatusAccepted      DocumentStatus = "ACCEPTED"
	DocumentStatusRejected      DocumentStatus = "REJECTED"
	DocumentStatusWaived        DocumentStatus = "WAIVED"
	DocumentStatusExpired       DocumentStatus = "EXPIRED"
	DocumentStatusSyncedToSIS   DocumentStatus = "SYNCED_TO_SIS"
)

// String returns the document status as a string.
func (status DocumentStatus) String() string {
	return string(status)
}

// ImportBatchStatus represents the operational state of an admissions import.
type ImportBatchStatus string

// Set of valid admissions import batch statuses.
const (
	ImportBatchStatusPreviewed        ImportBatchStatus = "PREVIEWED"
	ImportBatchStatusValidationFailed ImportBatchStatus = "VALIDATION_FAILED"
	ImportBatchStatusQueued           ImportBatchStatus = "QUEUED"
	ImportBatchStatusProcessing       ImportBatchStatus = "PROCESSING"
	ImportBatchStatusCompleted        ImportBatchStatus = "COMPLETED"
	ImportBatchStatusFailed           ImportBatchStatus = "FAILED"
)

// String returns the import batch status as a string.
func (status ImportBatchStatus) String() string {
	return string(status)
}

// ImportSource represents where an import file originated.
type ImportSource string

// Set of valid admissions import sources.
const (
	ImportSourceManualUpload ImportSource = "MANUAL_UPLOAD"
	ImportSourceSISExport    ImportSource = "SIS_EXPORT"
)

// String returns the import source as a string.
func (source ImportSource) String() string {
	return string(source)
}

// ImportFileType represents the supported import file formats.
type ImportFileType string

// Set of valid admissions import file types.
const (
	ImportFileTypeCSV  ImportFileType = "CSV"
	ImportFileTypeXLSX ImportFileType = "XLSX"
)

// String returns the import file type as a string.
func (fileType ImportFileType) String() string {
	return string(fileType)
}

// ImportTarget represents the admissions aggregate receiving imported rows.
type ImportTarget string

// Set of supported import targets.
const (
	ImportTargetConstituents ImportTarget = "CONSTITUENTS"
	ImportTargetApplications ImportTarget = "APPLICATIONS"
)

// String returns the import target as a string.
func (target ImportTarget) String() string {
	return string(target)
}

// SyncJobStatus represents the operational state of a SIS batch reconciliation run.
type SyncJobStatus string

// Set of valid SIS sync job statuses.
const (
	SyncJobStatusQueued     SyncJobStatus = "QUEUED"
	SyncJobStatusRunning    SyncJobStatus = "RUNNING"
	SyncJobStatusSucceeded  SyncJobStatus = "SUCCEEDED"
	SyncJobStatusFailed     SyncJobStatus = "FAILED"
	SyncJobStatusRetryReady SyncJobStatus = "RETRY_READY"
)

// String returns the sync job status as a string.
func (status SyncJobStatus) String() string {
	return string(status)
}

// IntegrationAdapter identifies the external Kenya integration provider for sync tracking.
// It intentionally lives in the root admissions domain to avoid importing adapters/* back
// into the package that adapter implementations already depend on.
type IntegrationAdapter string

// Set of supported Kenya integration adapters.
const (
	IntegrationAdapterKUCCPS        IntegrationAdapter = "kuccps"
	IntegrationAdapterKNEC          IntegrationAdapter = "knec"
	IntegrationAdapterIPRS          IntegrationAdapter = "iprs"
	IntegrationAdapterMPesaDaraja   IntegrationAdapter = "mpesa_daraja"
	IntegrationAdapterCelcomAfrica  IntegrationAdapter = "celcom_africa"
	IntegrationAdapterWhatsAppCloud IntegrationAdapter = "whatsapp_cloud"
)

// String returns the integration adapter as a string.
func (adapter IntegrationAdapter) String() string {
	return string(adapter)
}

// Valid reports whether the adapter is one of the supported Kenya providers.
func (adapter IntegrationAdapter) Valid() bool {
	switch adapter {
	case IntegrationAdapterKUCCPS,
		IntegrationAdapterKNEC,
		IntegrationAdapterIPRS,
		IntegrationAdapterMPesaDaraja,
		IntegrationAdapterCelcomAfrica,
		IntegrationAdapterWhatsAppCloud:
		return true
	default:
		return false
	}
}

// SyncEventStatus represents the queue state for a real-time SIS sync event.
type SyncEventStatus string

// Set of valid SIS sync event statuses.
const (
	SyncEventStatusQueued     SyncEventStatus = "QUEUED"
	SyncEventStatusProcessing SyncEventStatus = "PROCESSING"
	SyncEventStatusSucceeded  SyncEventStatus = "SUCCEEDED"
	SyncEventStatusFailed     SyncEventStatus = "FAILED"
	SyncEventStatusRetryReady SyncEventStatus = "RETRY_READY"
)

// String returns the sync event status as a string.
func (status SyncEventStatus) String() string {
	return string(status)
}

// SyncDirection represents whether CRM is pulling canonical SIS state or pushing CRM events to SIS.
type SyncDirection string

// Set of valid SIS sync directions.
const (
	SyncDirectionInbound  SyncDirection = "INBOUND"
	SyncDirectionOutbound SyncDirection = "OUTBOUND"
)

// String returns the sync direction as a string.
func (direction SyncDirection) String() string {
	return string(direction)
}

// SyncEventType represents the approved field-set actions exchanged with SIS.
type SyncEventType string

// Set of valid SIS sync event types.
const (
	SyncEventTypeBatchTermsPull         SyncEventType = "BATCH_TERMS_PULL"
	SyncEventTypeBatchProgramsPull      SyncEventType = "BATCH_PROGRAMS_PULL"
	SyncEventTypeBatchPersonMatchesPull SyncEventType = "BATCH_PERSON_MATCHES_PULL"
	SyncEventTypeBatchEnrollmentPull    SyncEventType = "BATCH_ENROLLMENT_PULL"
	SyncEventTypeApplicationSubmission  SyncEventType = "APPLICATION_SUBMISSION"
	SyncEventTypeApplicationDecision    SyncEventType = "APPLICATION_DECISION"
	SyncEventTypeDocumentStatus         SyncEventType = "DOCUMENT_STATUS"
	SyncEventTypeEnrollmentIntent       SyncEventType = "ENROLLMENT_INTENT"
	SyncEventTypeKUCCPSPlacementPull    SyncEventType = "KUCCPS_PLACEMENT_PULL"
	SyncEventTypeKUCCPSPlacementConfirm SyncEventType = "KUCCPS_PLACEMENT_CONFIRM"
	SyncEventTypeKNECResultVerification SyncEventType = "KNEC_RESULT_VERIFICATION"
	SyncEventTypeIPRSIdentityVerify     SyncEventType = "IPRS_IDENTITY_VERIFICATION"
	SyncEventTypeMPesaSTKPush           SyncEventType = "MPESA_STK_PUSH"
	SyncEventTypeMPesaC2BCallback       SyncEventType = "MPESA_C2B_CALLBACK"
	SyncEventTypeMPesaTransactionQuery  SyncEventType = "MPESA_TRANSACTION_QUERY"
	SyncEventTypeSMSOutbound            SyncEventType = "SMS_OUTBOUND"
	SyncEventTypeSMSDeliveryReport      SyncEventType = "SMS_DELIVERY_REPORT"
	SyncEventTypeWhatsAppMessageSend    SyncEventType = "WHATSAPP_MESSAGE_SEND"
	SyncEventTypeWhatsAppWebhookInbound SyncEventType = "WHATSAPP_WEBHOOK_INBOUND"
)

// String returns the sync event type as a string.
func (eventType SyncEventType) String() string {
	return string(eventType)
}

// ChecklistItem represents one document requirement for an application.
type ChecklistItem struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	ItemKey       string
	DocumentName  string
	Description   *string
	Required      bool
	Status        DocumentStatus
	DisplayOrder  int
	DateCreated   time.Time
	DateUpdated   time.Time
}

// NewChecklistItem is what we require to create a checklist item.
type NewChecklistItem struct {
	ApplicationID uuid.UUID
	ItemKey       string
	DocumentName  string
	Description   *string
	Required      bool
	DisplayOrder  int
}

// Document represents applicant-submitted evidence for checklist items.
type Document struct {
	ID              uuid.UUID
	ApplicationID   uuid.UUID
	ChecklistItemID uuid.UUID
	FileName        string
	ContentType     string
	SizeBytes       int64
	StorageKey      string
	Status          DocumentStatus
	ReviewerID      *uuid.UUID
	ReviewerNotes   *string
	UploadedByID    uuid.UUID
	UploadedAt      time.Time
	ReviewedAt      *time.Time
	DateCreated     time.Time
	DateUpdated     time.Time
}

// NewDocument is what we require to record uploaded document metadata.
type NewDocument struct {
	ApplicationID   uuid.UUID
	ChecklistItemID uuid.UUID
	FileName        string
	ContentType     string
	SizeBytes       int64
	StorageKey      string
	UploadedByID    uuid.UUID
}

// NewDocumentVerification is what we require to review document metadata.
type NewDocumentVerification struct {
	Status        DocumentStatus
	ReviewerID    uuid.UUID
	ReviewerNotes *string
}

// ImportBatch records one import preview or commit run.
type ImportBatch struct {
	ID                uuid.UUID
	Source            ImportSource
	FileType          ImportFileType
	Target            ImportTarget
	Status            ImportBatchStatus
	FileName          string
	StorageKey        *string
	UploadedByID      uuid.UUID
	TotalRows         int
	ValidRows         int
	InvalidRows       int
	DuplicateRows     int
	FieldMapping      map[string]string
	InvalidReportKey  *string
	ValidationSummary *string
	CommittedAt       *time.Time
	DateCreated       time.Time
	DateUpdated       time.Time
}

// NewImportBatch is what we require to record an admissions import preview or commit.
type NewImportBatch struct {
	Source            ImportSource
	FileType          ImportFileType
	Target            ImportTarget
	Status            ImportBatchStatus
	FileName          string
	StorageKey        *string
	UploadedByID      uuid.UUID
	TotalRows         int
	ValidRows         int
	InvalidRows       int
	DuplicateRows     int
	FieldMapping      map[string]string
	InvalidReportKey  *string
	ValidationSummary *string
}

// ImportInvalidRow records one invalid row discovered during preview or commit validation.
type ImportInvalidRow struct {
	ID          uuid.UUID
	BatchID     uuid.UUID
	RowNumber   int
	FieldName   *string
	RawData     map[string]string
	ErrorCode   string
	ErrorDetail string
	DateCreated time.Time
}

// NewImportInvalidRow is what we require to attach an invalid row to an import batch.
type NewImportInvalidRow struct {
	BatchID     uuid.UUID
	RowNumber   int
	FieldName   *string
	RawData     map[string]string
	ErrorCode   string
	ErrorDetail string
}

// SyncJob tracks a SIS batch reconciliation run and its operational result.
type SyncJob struct {
	ID                uuid.UUID
	Name              string
	Adapter           IntegrationAdapter
	Operation         string
	Status            SyncJobStatus
	Direction         SyncDirection
	StartedAt         *time.Time
	CompletedAt       *time.Time
	RecordsPulled     int
	RecordsPushed     int
	EventsRequeued    int
	AttemptCount      int
	MaxAttempts       int
	NextRetryAt       *time.Time
	ExternalRef       *string
	ExternalReceiptID *string
	ErrorCode         *string
	ErrorDetail       *string
	LastErrorAt       *time.Time
	FailureReason     *string
	Retryable         bool
	CreatedByID       *uuid.UUID
	DateCreated       time.Time
	DateUpdated       time.Time
}

// NewSyncJob is what the sync runner requires when starting or scheduling a batch reconciliation run.
type NewSyncJob struct {
	Name        string
	Adapter     IntegrationAdapter
	Operation   string
	Direction   SyncDirection
	Status      SyncJobStatus
	StartedAt   *time.Time
	MaxAttempts int
	CreatedByID *uuid.UUID
}

// UpdateSyncJob is what the sync runner may change as a batch reconciliation run progresses.
type UpdateSyncJob struct {
	Status            SyncJobStatus
	CompletedAt       *time.Time
	RecordsPulled     int
	RecordsPushed     int
	EventsRequeued    int
	AttemptCount      int
	NextRetryAt       *time.Time
	ExternalRef       *string
	ExternalReceiptID *string
	ErrorCode         *string
	ErrorDetail       *string
	LastErrorAt       *time.Time
	FailureReason     *string
	Retryable         bool
}

// SyncEvent tracks a selected real-time SIS event, including retries and failure visibility.
type SyncEvent struct {
	ID                uuid.UUID
	JobID             *uuid.UUID
	Adapter           IntegrationAdapter
	Operation         string
	EventType         SyncEventType
	Status            SyncEventStatus
	Direction         SyncDirection
	ResourceType      string
	ResourceID        uuid.UUID
	ExternalRef       *string
	ExternalReceiptID *string
	PayloadHash       string
	Attempts          int
	MaxAttempts       int
	NextRetryAt       *time.Time
	ErrorCode         *string
	ErrorDetail       *string
	LastErrorAt       *time.Time
	FailureReason     *string
	AuditMessage      string
	DateCreated       time.Time
	DateUpdated       time.Time
}

// NewSyncEvent is what workflow hooks provide when enqueueing a selected real-time SIS event.
type NewSyncEvent struct {
	JobID             *uuid.UUID
	Adapter           IntegrationAdapter
	Operation         string
	EventType         SyncEventType
	Direction         SyncDirection
	ResourceType      string
	ResourceID        uuid.UUID
	ExternalRef       *string
	ExternalReceiptID *string
	PayloadHash       string
	MaxAttempts       int
	AuditMessage      string
}

// UpdateSyncEvent is what the queue runner may change as an event is processed or retried.
type UpdateSyncEvent struct {
	Status            SyncEventStatus
	Attempts          int
	NextRetryAt       *time.Time
	ExternalRef       *string
	ExternalReceiptID *string
	ErrorCode         *string
	ErrorDetail       *string
	LastErrorAt       *time.Time
	FailureReason     *string
	AuditMessage      string
}

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
type DuplicateReview struct {
	ID                     uuid.UUID
	SourceConstituentID    uuid.UUID
	CandidateConstituentID uuid.UUID
	MatchType              DuplicateReviewMatchType
	MatchScore             int
	MatchReason            string
	Status                 DuplicateReviewStatus
	ResolvedBy             *uuid.UUID
	ResolvedAt             *time.Time
	ResolutionNote         *string
	DateCreated            time.Time
	DateUpdated            time.Time
}

// NewDuplicateReview is what we require when enqueueing a possible duplicate.
type NewDuplicateReview struct {
	SourceConstituentID    uuid.UUID
	CandidateConstituentID uuid.UUID
	MatchType              DuplicateReviewMatchType
	MatchScore             int
	MatchReason            string
}

// ResolveDuplicateReview defines a staff decision for a duplicate review.
type ResolveDuplicateReview struct {
	Resolution DuplicateReviewResolution
	ActorID    uuid.UUID
	Note       *string
}

// AggregateNames returns the scaffolded admissions aggregate names.
func AggregateNames() []string {
	return []string{
		"staffProfile",
		"applicantProfile",
		"constituent",
		"inquiry",
		"application",
		"customFieldDefinition",
		"customFieldValue",
		"leadScoreRule",
		"leadScore",
		"checklist",
		"document",
		"importBatch",
		"importInvalidRow",
		"decision",
		"syncJob",
		"syncEvent",
		"program",
		"academicTerm",
		"duplicateReview",
	}
}

// StaffProfileQueryFilter holds the available fields a staff profile query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type StaffProfileQueryFilter struct {
	ID     *uuid.UUID
	UserID *uuid.UUID
	Role   *AdmissionsRole
	Active *bool
}

// ApplicantProfileQueryFilter holds the available fields an applicant profile query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ApplicantProfileQueryFilter struct {
	ID            *uuid.UUID
	UserID        *uuid.UUID
	ConstituentID *uuid.UUID
	Active        *bool
}

// InquiryQueryFilter holds the available fields an inquiry query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type InquiryQueryFilter struct {
	ID            *uuid.UUID
	ConstituentID *uuid.UUID
	PrimaryEmail  *mail.Address
	Source        *string
	Status        *InquiryStatus
}

// LeadScoreRuleQueryFilter holds the available fields a lead score rule query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type LeadScoreRuleQueryFilter struct {
	ID     *uuid.UUID
	Active *bool
}

// LeadScoreQueryFilter holds the available fields a lead score query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type LeadScoreQueryFilter struct {
	ID            *uuid.UUID
	ConstituentID *uuid.UUID
	Band          *LeadScoreBand
	MinScore      *int
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
	NationalID      *string
	UPI             *string
	KCSEIndexNumber *string
	LifecycleStage  *LifecycleStage
	DuplicateStatus *DuplicateStatus
	CustomFields    map[string]string
}

// DuplicateReviewQueryFilter holds the available fields a duplicate review query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type DuplicateReviewQueryFilter struct {
	ID                     *uuid.UUID
	SourceConstituentID    *uuid.UUID
	CandidateConstituentID *uuid.UUID
	MatchType              *DuplicateReviewMatchType
	Status                 *DuplicateReviewStatus
}

// ApplicationQueryFilter holds the available fields an application query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ApplicationQueryFilter struct {
	ID              *uuid.UUID
	ConstituentID   *uuid.UUID
	ProgramID       *uuid.UUID
	AcademicTermID  *uuid.UUID
	ApplicationType *ApplicationType
	Status          *ApplicationStatus
	ActiveOnly      *bool
	CustomFields    map[string]string
}

// EventQueryFilter holds the available fields an event query can be filtered on.
type EventQueryFilter struct {
	ID       *uuid.UUID
	Type     *EventType
	Status   *EventStatus
	Virtual  *bool
}

// EventRegistrationQueryFilter holds the available fields an event registration query can be filtered on.
type EventRegistrationQueryFilter struct {
	ID            *uuid.UUID
	EventID       *uuid.UUID
	ConstituentID *uuid.UUID
	Status        *EventRegistrationStatus
	MatchStatus   *EventRegistrationMatchStatus
	Source        *EventRegistrationSource
}

// CustomFieldDefinitionQueryFilter holds fields a custom field definition query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type CustomFieldDefinitionQueryFilter struct {
	ID     *uuid.UUID
	Owner  *CustomFieldOwner
	Active *bool
}

// CustomFieldValueQueryFilter holds fields a custom field value query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type CustomFieldValueQueryFilter struct {
	ID           *uuid.UUID
	DefinitionID *uuid.UUID
	Owner        *CustomFieldOwner
	OwnerID      *uuid.UUID
}

// ApplicationFormTemplateQueryFilter holds fields an application form template query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ApplicationFormTemplateQueryFilter struct {
	ID              *uuid.UUID
	ProgramID       *uuid.UUID
	AcademicTermID  *uuid.UUID
	ApplicationType *ApplicationType
	Active          *bool
	Version         *int
}

// ApplicationTransitionQueryFilter holds the available fields an application transition query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ApplicationTransitionQueryFilter struct {
	ID            *uuid.UUID
	ApplicationID *uuid.UUID
	ActorID       *uuid.UUID
	FromStatus    *ApplicationStatus
	ToStatus      *ApplicationStatus
}

// ChecklistItemQueryFilter holds fields a checklist item query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ChecklistItemQueryFilter struct {
	ID            *uuid.UUID
	ApplicationID *uuid.UUID
	Status        *DocumentStatus
	Required      *bool
}

// DocumentQueryFilter holds fields a document query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type DocumentQueryFilter struct {
	ID              *uuid.UUID
	ApplicationID   *uuid.UUID
	ChecklistItemID *uuid.UUID
	Status          *DocumentStatus
	UploadedByID    *uuid.UUID
	ReviewerID      *uuid.UUID
}

// ImportBatchQueryFilter holds fields an import batch query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ImportBatchQueryFilter struct {
	ID           *uuid.UUID
	Source       *ImportSource
	FileType     *ImportFileType
	Target       *ImportTarget
	Status       *ImportBatchStatus
	UploadedByID *uuid.UUID
}

// ImportInvalidRowQueryFilter holds fields an import invalid row query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ImportInvalidRowQueryFilter struct {
	ID      *uuid.UUID
	BatchID *uuid.UUID
}

// SyncJobQueryFilter holds fields a SIS sync job query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type SyncJobQueryFilter struct {
	ID              *uuid.UUID
	Adapter         *IntegrationAdapter
	Status          *SyncJobStatus
	Direction       *SyncDirection
	Retryable       *bool
	NextRetryBefore *time.Time
}

// SyncEventQueryFilter holds fields a SIS sync event query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type SyncEventQueryFilter struct {
	ID           *uuid.UUID
	JobID        *uuid.UUID
	Adapter      *IntegrationAdapter
	EventType    *SyncEventType
	Status       *SyncEventStatus
	Direction    *SyncDirection
	ResourceType *string
	ResourceID   *uuid.UUID
}
