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
	ID              uuid.UUID
	FirstName       string
	LastName        string
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	DateOfBirth     time.Time
	PrimaryEmail    mail.Address
	PrimaryPhone    string
	ExternalSISID   *string
	LifecycleStage  LifecycleStage
	DuplicateStatus DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
	DateCreated     time.Time
	DateUpdated     time.Time
}

// NewConstituent is what we require from clients when adding a Constituent.
type NewConstituent struct {
	FirstName       string
	LastName        string
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	DateOfBirth     time.Time
	PrimaryEmail    mail.Address
	PrimaryPhone    string
	ExternalSISID   *string
	LifecycleStage  LifecycleStage
	DuplicateStatus DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
}

// UpdateConstituent defines what information may be provided to modify a Constituent.
type UpdateConstituent struct {
	PreferredName   *string
	MiddleName      *string
	Suffix          *string
	PrimaryEmail    *mail.Address
	PrimaryPhone    *string
	LifecycleStage  *LifecycleStage
	DuplicateStatus *DuplicateStatus
	DuplicateOfID   *uuid.UUID
	SISSyncedAt     *time.Time
}

// Inquiry captures pre-application interest in the school.
type Inquiry struct{}

// ApplicationType represents the admissions category for an application.
type ApplicationType string

// Set of valid application types.
const (
	ApplicationTypeFreshman ApplicationType = "FRESHMAN"
	ApplicationTypeTransfer ApplicationType = "TRANSFER"
	ApplicationTypeGraduate ApplicationType = "GRADUATE"
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

// Application represents a constituent's program application for a term.
type Application struct {
	ID                 uuid.UUID
	ConstituentID      uuid.UUID
	ProgramID          uuid.UUID
	AcademicTermID     uuid.UUID
	ApplicationType    ApplicationType
	Status             ApplicationStatus
	AssignedReviewerID *uuid.UUID
	SubmittedAt        *time.Time
	DateCreated        time.Time
	DateUpdated        time.Time
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
	AssignedReviewerID *uuid.UUID
}

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
		"leadScoreRule",
		"leadScore",
		"checklist",
		"document",
		"decision",
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
	LifecycleStage  *LifecycleStage
	DuplicateStatus *DuplicateStatus
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
