// Package admissionsbus provides business access to the admissions domain.
package admissionsbus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/sdk/delegate"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Set of error variables for admissions reference data operations.
var (
	ErrConstituentNotFound          = errors.New("constituent not found")
	ErrFirstNameRequired            = errors.New("first name required")
	ErrLastNameRequired             = errors.New("last name required")
	ErrDateOfBirthRequired          = errors.New("date of birth required")
	ErrDateOfBirthInFuture          = errors.New("date of birth cannot be in the future")
	ErrPrimaryPhoneRequired         = errors.New("primary phone required")
	ErrInvalidLifecycleStage        = errors.New("invalid lifecycle stage")
	ErrInvalidDuplicateStatus       = errors.New("invalid duplicate status")
	ErrInvalidDuplicateLink         = errors.New("duplicate status does not match duplicate link")
	ErrInvalidLifecycleChange       = errors.New("invalid lifecycle stage change")
	ErrInquiryNotFound              = errors.New("inquiry not found")
	ErrInquirySourceRequired        = errors.New("inquiry source required")
	ErrInvalidInquiryStatus         = errors.New("invalid inquiry status")
	ErrProgramNotFound              = errors.New("program not found")
	ErrAcademicTermNotFound         = errors.New("academic term not found")
	ErrInvalidTermDateRange         = errors.New("term start date must be before end date")
	ErrInvalidApplicationWindow     = errors.New("application deadline must be on or after application start date")
	ErrDuplicateReviewNotFound      = errors.New("duplicate review not found")
	ErrInvalidDuplicateReview       = errors.New("invalid duplicate review")
	ErrInvalidMatchType             = errors.New("invalid duplicate match type")
	ErrInvalidMatchScore            = errors.New("duplicate match score must be between 0 and 100")
	ErrMatchReasonRequired          = errors.New("duplicate match reason required")
	ErrInvalidReviewStatus          = errors.New("invalid duplicate review status")
	ErrInvalidResolution            = errors.New("invalid duplicate review resolution")
	ErrDuplicateReviewResolved      = errors.New("duplicate review already resolved")
	ErrResolutionActorRequired      = errors.New("resolution actor required")
	ErrApplicationNotFound          = errors.New("application not found")
	ErrInvalidApplicationType       = errors.New("invalid application type")
	ErrInvalidApplicationStatus     = errors.New("invalid application status")
	ErrDuplicateApplication         = errors.New("active application already exists for constituent term and program")
	ErrConstituentIDRequired        = errors.New("constituent id required")
	ErrProgramIDRequired            = errors.New("program id required")
	ErrAcademicTermIDRequired       = errors.New("academic term id required")
	ErrInactiveProgram              = errors.New("program is inactive")
	ErrInactiveAcademicTerm         = errors.New("academic term is inactive")
	ErrInvalidApplicationTransition = errors.New("invalid application status transition")
	ErrApplicationActorRequired     = errors.New("application transition actor required")
	ErrFormTemplateNotFound         = errors.New("application form template not found")
	ErrFormTemplateNameRequired     = errors.New("application form template name required")
	ErrFormTemplateFieldsRequired   = errors.New("application form template required fields required")
	ErrFormTemplateFieldInvalid     = errors.New("application form template field invalid")
	ErrFormTemplateChecklistInvalid = errors.New("application form template checklist item invalid")
	ErrFormTemplatePriorityInvalid  = errors.New("application form template priority must be greater than or equal to zero")
	ErrStaffProfileNotFound         = errors.New("staff profile not found")
	ErrStaffProfileUserRequired     = errors.New("staff profile user id required")
	ErrStaffProfileRoleRequired     = errors.New("staff profile role required")
	ErrApplicantProfileNotFound     = errors.New("applicant profile not found")
	ErrApplicantProfileUserRequired = errors.New("applicant profile user id required")
	ErrInvalidAdmissionsRole        = errors.New("invalid admissions role")
	ErrLeadScoreRuleNotFound        = errors.New("lead score rule not found")
	ErrLeadScoreNotFound            = errors.New("lead score not found")
	ErrLeadScoreRuleNameRequired    = errors.New("lead score rule name required")
	ErrLeadScoreCriteriaRequired    = errors.New("lead score criteria required")
	ErrInvalidLeadScoreCriterion    = errors.New("invalid lead score criterion")
	ErrInvalidLeadScorePoints       = errors.New("lead score points must be greater than or equal to zero")
	ErrInvalidLeadScorePriority     = errors.New("lead score priority must be greater than or equal to zero")
	ErrInvalidLeadScoreBand         = errors.New("invalid lead score band")
	ErrChecklistItemNotFound        = errors.New("checklist item not found")
	ErrChecklistItemKeyRequired     = errors.New("checklist item key required")
	ErrChecklistItemNameRequired    = errors.New("checklist item document name required")
	ErrChecklistItemOrderInvalid    = errors.New("checklist item display order must be greater than or equal to zero")
	ErrDocumentNotFound             = errors.New("document not found")
	ErrInvalidDocumentStatus        = errors.New("invalid document status")
	ErrDocumentFileNameRequired     = errors.New("document file name required")
	ErrDocumentContentTypeRequired  = errors.New("document content type required")
	ErrDocumentSizeInvalid          = errors.New("document size must be greater than zero")
	ErrDocumentStorageKeyRequired   = errors.New("document storage key required")
	ErrDocumentUploaderRequired     = errors.New("document uploader required")
	ErrDocumentReviewerRequired     = errors.New("document reviewer required")
	ErrDocumentStatusNotReviewable  = errors.New("document status is not a review action")
	ErrSyncJobNotFound              = errors.New("sync job not found")
	ErrSyncJobNameRequired          = errors.New("sync job name required")
	ErrInvalidSyncJobStatus         = errors.New("invalid sync job status")
	ErrInvalidSyncEventStatus       = errors.New("invalid sync event status")
	ErrInvalidSyncDirection         = errors.New("invalid sync direction")
	ErrInvalidSyncEventType         = errors.New("invalid sync event type")
	ErrSyncEventNotFound            = errors.New("sync event not found")
	ErrSyncEventResourceRequired    = errors.New("sync event resource required")
	ErrSyncEventPayloadHashRequired = errors.New("sync event payload hash required")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve admissions data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Health(ctx context.Context) (Health, error)
	CreateStaffProfile(ctx context.Context, profile StaffProfile) error
	UpdateStaffProfile(ctx context.Context, profile StaffProfile) error
	QueryStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter, orderBy order.By, page page.Page) ([]StaffProfile, error)
	CountStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter) (int, error)
	QueryStaffProfileByID(ctx context.Context, profileID uuid.UUID) (StaffProfile, error)
	QueryStaffProfileByUserID(ctx context.Context, userID uuid.UUID) (StaffProfile, error)
	CreateApplicantProfile(ctx context.Context, profile ApplicantProfile) error
	UpdateApplicantProfile(ctx context.Context, profile ApplicantProfile) error
	QueryApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter, orderBy order.By, page page.Page) ([]ApplicantProfile, error)
	CountApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter) (int, error)
	QueryApplicantProfileByID(ctx context.Context, profileID uuid.UUID) (ApplicantProfile, error)
	QueryApplicantProfileByUserID(ctx context.Context, userID uuid.UUID) (ApplicantProfile, error)
	QueryApplicantProfileByConstituentID(ctx context.Context, constituentID uuid.UUID) (ApplicantProfile, error)
	CreateInquiry(ctx context.Context, inquiry Inquiry) error
	UpdateInquiry(ctx context.Context, inquiry Inquiry) error
	QueryInquiries(ctx context.Context, filter InquiryQueryFilter, orderBy order.By, page page.Page) ([]Inquiry, error)
	CountInquiries(ctx context.Context, filter InquiryQueryFilter) (int, error)
	QueryInquiryByID(ctx context.Context, inquiryID uuid.UUID) (Inquiry, error)
	CreateLeadScoreRule(ctx context.Context, rule LeadScoreRule) error
	UpdateLeadScoreRule(ctx context.Context, rule LeadScoreRule) error
	QueryLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter, orderBy order.By, page page.Page) ([]LeadScoreRule, error)
	CountLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter) (int, error)
	QueryLeadScoreRuleByID(ctx context.Context, ruleID uuid.UUID) (LeadScoreRule, error)
	UpsertLeadScore(ctx context.Context, score LeadScore) error
	QueryLeadScores(ctx context.Context, filter LeadScoreQueryFilter, orderBy order.By, page page.Page) ([]LeadScore, error)
	CountLeadScores(ctx context.Context, filter LeadScoreQueryFilter) (int, error)
	QueryLeadScoreByID(ctx context.Context, scoreID uuid.UUID) (LeadScore, error)
	QueryLeadScoreByConstituentID(ctx context.Context, constituentID uuid.UUID) (LeadScore, error)
	CreateConstituent(ctx context.Context, cst Constituent) error
	UpdateConstituent(ctx context.Context, cst Constituent) error
	QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error)
	CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error)
	QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error)
	QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error)
	QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error)
	UpsertProgram(ctx context.Context, prg Program) error
	QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error)
	CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error)
	QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error)
	QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error)
	UpsertAcademicTerm(ctx context.Context, term AcademicTerm) error
	QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error)
	CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error)
	QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error)
	QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error)
	CreateDuplicateReview(ctx context.Context, review DuplicateReview) error
	UpdateDuplicateReview(ctx context.Context, review DuplicateReview) error
	QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error)
	CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error)
	QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error)
	CreateApplication(ctx context.Context, app Application) error
	QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error)
	CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error)
	QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error)
	QueryActiveApplicationByTuple(ctx context.Context, constituentID uuid.UUID, academicTermID uuid.UUID, programID uuid.UUID) (Application, error)
	UpdateApplication(ctx context.Context, app Application) error
	CreateApplicationFormTemplate(ctx context.Context, template ApplicationFormTemplate) error
	UpdateApplicationFormTemplate(ctx context.Context, template ApplicationFormTemplate) error
	QueryApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter, orderBy order.By, page page.Page) ([]ApplicationFormTemplate, error)
	CountApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter) (int, error)
	QueryApplicationFormTemplateByID(ctx context.Context, templateID uuid.UUID) (ApplicationFormTemplate, error)
	CreateApplicationTransition(ctx context.Context, transition ApplicationTransition) error
	QueryApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter, orderBy order.By, page page.Page) ([]ApplicationTransition, error)
	CountApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter) (int, error)
	CreateChecklistItem(ctx context.Context, item ChecklistItem) error
	UpdateChecklistItem(ctx context.Context, item ChecklistItem) error
	QueryChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter, orderBy order.By, page page.Page) ([]ChecklistItem, error)
	CountChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter) (int, error)
	QueryChecklistItemByID(ctx context.Context, itemID uuid.UUID) (ChecklistItem, error)
	CreateDocument(ctx context.Context, document Document) error
	UpdateDocument(ctx context.Context, document Document) error
	QueryDocuments(ctx context.Context, filter DocumentQueryFilter, orderBy order.By, page page.Page) ([]Document, error)
	CountDocuments(ctx context.Context, filter DocumentQueryFilter) (int, error)
	QueryDocumentByID(ctx context.Context, documentID uuid.UUID) (Document, error)
	CreateSyncJob(ctx context.Context, job SyncJob) error
	UpdateSyncJob(ctx context.Context, job SyncJob) error
	QuerySyncJobs(ctx context.Context, filter SyncJobQueryFilter, orderBy order.By, page page.Page) ([]SyncJob, error)
	CountSyncJobs(ctx context.Context, filter SyncJobQueryFilter) (int, error)
	QuerySyncJobByID(ctx context.Context, jobID uuid.UUID) (SyncJob, error)
	CreateSyncEvent(ctx context.Context, event SyncEvent) error
	UpdateSyncEvent(ctx context.Context, event SyncEvent) error
	QuerySyncEvents(ctx context.Context, filter SyncEventQueryFilter, orderBy order.By, page page.Page) ([]SyncEvent, error)
	CountSyncEvents(ctx context.Context, filter SyncEventQueryFilter) (int, error)
	QuerySyncEventByID(ctx context.Context, eventID uuid.UUID) (SyncEvent, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Health(ctx context.Context) (Health, error)
	CreateStaffProfile(ctx context.Context, np NewStaffProfile) (StaffProfile, error)
	UpdateStaffProfile(ctx context.Context, profile StaffProfile, np NewStaffProfile) (StaffProfile, error)
	QueryStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter, orderBy order.By, page page.Page) ([]StaffProfile, error)
	CountStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter) (int, error)
	QueryStaffProfileByID(ctx context.Context, profileID uuid.UUID) (StaffProfile, error)
	QueryStaffProfileByUserID(ctx context.Context, userID uuid.UUID) (StaffProfile, error)
	CreateApplicantProfile(ctx context.Context, np NewApplicantProfile) (ApplicantProfile, error)
	UpdateApplicantProfile(ctx context.Context, profile ApplicantProfile, np NewApplicantProfile) (ApplicantProfile, error)
	QueryApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter, orderBy order.By, page page.Page) ([]ApplicantProfile, error)
	CountApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter) (int, error)
	QueryApplicantProfileByID(ctx context.Context, profileID uuid.UUID) (ApplicantProfile, error)
	QueryApplicantProfileByUserID(ctx context.Context, userID uuid.UUID) (ApplicantProfile, error)
	QueryApplicantProfileByConstituentID(ctx context.Context, constituentID uuid.UUID) (ApplicantProfile, error)
	CreateInquiry(ctx context.Context, ni NewInquiry) (Inquiry, error)
	QueryInquiries(ctx context.Context, filter InquiryQueryFilter, orderBy order.By, page page.Page) ([]Inquiry, error)
	CountInquiries(ctx context.Context, filter InquiryQueryFilter) (int, error)
	QueryInquiryByID(ctx context.Context, inquiryID uuid.UUID) (Inquiry, error)
	CreateLeadScoreRule(ctx context.Context, nr NewLeadScoreRule) (LeadScoreRule, error)
	UpdateLeadScoreRule(ctx context.Context, rule LeadScoreRule, nr NewLeadScoreRule) (LeadScoreRule, error)
	QueryLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter, orderBy order.By, page page.Page) ([]LeadScoreRule, error)
	CountLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter) (int, error)
	QueryLeadScoreRuleByID(ctx context.Context, ruleID uuid.UUID) (LeadScoreRule, error)
	RecalculateLeadScoreForConstituent(ctx context.Context, constituentID uuid.UUID) (LeadScore, error)
	QueryLeadScores(ctx context.Context, filter LeadScoreQueryFilter, orderBy order.By, page page.Page) ([]LeadScore, error)
	CountLeadScores(ctx context.Context, filter LeadScoreQueryFilter) (int, error)
	QueryLeadScoreByID(ctx context.Context, scoreID uuid.UUID) (LeadScore, error)
	QueryLeadScoreByConstituentID(ctx context.Context, constituentID uuid.UUID) (LeadScore, error)
	CreateConstituent(ctx context.Context, nc NewConstituent) (Constituent, error)
	UpdateConstituent(ctx context.Context, cst Constituent, uc UpdateConstituent) (Constituent, error)
	QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error)
	CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error)
	QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error)
	QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error)
	QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error)
	UpsertProgram(ctx context.Context, up UpsertProgram) (Program, error)
	QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error)
	CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error)
	QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error)
	QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error)
	UpsertAcademicTerm(ctx context.Context, up UpsertAcademicTerm) (AcademicTerm, error)
	QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error)
	CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error)
	QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error)
	QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error)
	CreateDuplicateReview(ctx context.Context, nr NewDuplicateReview) (DuplicateReview, error)
	ResolveDuplicateReview(ctx context.Context, review DuplicateReview, rr ResolveDuplicateReview) (DuplicateReview, error)
	QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error)
	CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error)
	QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error)
	CreateApplication(ctx context.Context, na NewApplication) (Application, error)
	QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error)
	CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error)
	QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error)
	TransitionApplicationStatus(ctx context.Context, app Application, nt NewApplicationTransition) (Application, ApplicationTransition, error)
	CreateApplicationFormTemplate(ctx context.Context, nt NewApplicationFormTemplate) (ApplicationFormTemplate, error)
	UpdateApplicationFormTemplate(ctx context.Context, template ApplicationFormTemplate, nt NewApplicationFormTemplate) (ApplicationFormTemplate, error)
	QueryApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter, orderBy order.By, page page.Page) ([]ApplicationFormTemplate, error)
	CountApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter) (int, error)
	QueryApplicationFormTemplateByID(ctx context.Context, templateID uuid.UUID) (ApplicationFormTemplate, error)
	QueryApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter, orderBy order.By, page page.Page) ([]ApplicationTransition, error)
	CountApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter) (int, error)
	CreateChecklistItem(ctx context.Context, ni NewChecklistItem) (ChecklistItem, error)
	UpdateChecklistItem(ctx context.Context, item ChecklistItem, ni NewChecklistItem) (ChecklistItem, error)
	QueryChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter, orderBy order.By, page page.Page) ([]ChecklistItem, error)
	CountChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter) (int, error)
	QueryChecklistItemByID(ctx context.Context, itemID uuid.UUID) (ChecklistItem, error)
	CreateDocument(ctx context.Context, nd NewDocument) (Document, error)
	VerifyDocument(ctx context.Context, document Document, nv NewDocumentVerification) (Document, error)
	QueryDocuments(ctx context.Context, filter DocumentQueryFilter, orderBy order.By, page page.Page) ([]Document, error)
	CountDocuments(ctx context.Context, filter DocumentQueryFilter) (int, error)
	QueryDocumentByID(ctx context.Context, documentID uuid.UUID) (Document, error)
	CreateSyncJob(ctx context.Context, nj NewSyncJob) (SyncJob, error)
	UpdateSyncJob(ctx context.Context, job SyncJob, uj UpdateSyncJob) (SyncJob, error)
	QuerySyncJobs(ctx context.Context, filter SyncJobQueryFilter, orderBy order.By, page page.Page) ([]SyncJob, error)
	CountSyncJobs(ctx context.Context, filter SyncJobQueryFilter) (int, error)
	QuerySyncJobByID(ctx context.Context, jobID uuid.UUID) (SyncJob, error)
	EnqueueSyncEvent(ctx context.Context, ne NewSyncEvent) (SyncEvent, error)
	UpdateSyncEvent(ctx context.Context, event SyncEvent, ue UpdateSyncEvent) (SyncEvent, error)
	QuerySyncEvents(ctx context.Context, filter SyncEventQueryFilter, orderBy order.By, page page.Page) ([]SyncEvent, error)
	CountSyncEvents(ctx context.Context, filter SyncEventQueryFilter) (int, error)
	QuerySyncEventByID(ctx context.Context, eventID uuid.UUID) (SyncEvent, error)
}

// CreateStaffProfile adds a context-specific admissions staff profile for an identity user.
func (b *Business) CreateStaffProfile(ctx context.Context, np NewStaffProfile) (StaffProfile, error) {
	if err := validateNewStaffProfile(np); err != nil {
		return StaffProfile{}, err
	}

	now := time.Now()
	profile := StaffProfile{
		ID:          uuid.New(),
		UserID:      np.UserID,
		Roles:       np.Roles,
		Active:      np.Active,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.CreateStaffProfile(ctx, profile); err != nil {
		return StaffProfile{}, fmt.Errorf("create staff profile: %w", err)
	}

	return profile, nil
}

// UpdateStaffProfile replaces mutable admissions staff profile data.
func (b *Business) UpdateStaffProfile(ctx context.Context, profile StaffProfile, np NewStaffProfile) (StaffProfile, error) {
	if err := validateNewStaffProfile(np); err != nil {
		return StaffProfile{}, err
	}

	profile.UserID = np.UserID
	profile.Roles = np.Roles
	profile.Active = np.Active
	profile.DateUpdated = time.Now()

	if err := b.storer.UpdateStaffProfile(ctx, profile); err != nil {
		return StaffProfile{}, fmt.Errorf("update staff profile: %w", err)
	}

	return profile, nil
}

// QueryStaffProfiles retrieves a list of admissions staff profiles.
func (b *Business) QueryStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter, orderBy order.By, page page.Page) ([]StaffProfile, error) {
	profiles, err := b.storer.QueryStaffProfiles(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query staff profiles: %w", err)
	}

	return profiles, nil
}

// CountStaffProfiles returns the total number of admissions staff profiles.
func (b *Business) CountStaffProfiles(ctx context.Context, filter StaffProfileQueryFilter) (int, error) {
	return b.storer.CountStaffProfiles(ctx, filter)
}

// QueryStaffProfileByID finds an admissions staff profile by ID.
func (b *Business) QueryStaffProfileByID(ctx context.Context, profileID uuid.UUID) (StaffProfile, error) {
	profile, err := b.storer.QueryStaffProfileByID(ctx, profileID)
	if err != nil {
		return StaffProfile{}, fmt.Errorf("query staff profile: profileID[%s]: %w", profileID, err)
	}

	return profile, nil
}

// QueryStaffProfileByUserID finds an admissions staff profile by identity user ID.
func (b *Business) QueryStaffProfileByUserID(ctx context.Context, userID uuid.UUID) (StaffProfile, error) {
	profile, err := b.storer.QueryStaffProfileByUserID(ctx, userID)
	if err != nil {
		return StaffProfile{}, fmt.Errorf("query staff profile: userID[%s]: %w", userID, err)
	}

	return profile, nil
}

// CreateApplicantProfile adds a portal applicant context profile for an identity user.
func (b *Business) CreateApplicantProfile(ctx context.Context, np NewApplicantProfile) (ApplicantProfile, error) {
	if err := validateNewApplicantProfile(np); err != nil {
		return ApplicantProfile{}, err
	}

	if _, err := b.QueryConstituentByID(ctx, np.ConstituentID); err != nil {
		return ApplicantProfile{}, err
	}

	now := time.Now()
	profile := ApplicantProfile{
		ID:            uuid.New(),
		UserID:        np.UserID,
		ConstituentID: np.ConstituentID,
		Active:        np.Active,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := b.storer.CreateApplicantProfile(ctx, profile); err != nil {
		return ApplicantProfile{}, fmt.Errorf("create applicant profile: %w", err)
	}

	return profile, nil
}

// UpdateApplicantProfile replaces mutable admissions applicant profile data.
func (b *Business) UpdateApplicantProfile(ctx context.Context, profile ApplicantProfile, np NewApplicantProfile) (ApplicantProfile, error) {
	if err := validateNewApplicantProfile(np); err != nil {
		return ApplicantProfile{}, err
	}

	if _, err := b.QueryConstituentByID(ctx, np.ConstituentID); err != nil {
		return ApplicantProfile{}, err
	}

	profile.UserID = np.UserID
	profile.ConstituentID = np.ConstituentID
	profile.Active = np.Active
	profile.DateUpdated = time.Now()

	if err := b.storer.UpdateApplicantProfile(ctx, profile); err != nil {
		return ApplicantProfile{}, fmt.Errorf("update applicant profile: %w", err)
	}

	return profile, nil
}

// QueryApplicantProfiles retrieves a list of admissions applicant profiles.
func (b *Business) QueryApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter, orderBy order.By, page page.Page) ([]ApplicantProfile, error) {
	profiles, err := b.storer.QueryApplicantProfiles(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query applicant profiles: %w", err)
	}

	return profiles, nil
}

// CountApplicantProfiles returns the total number of admissions applicant profiles.
func (b *Business) CountApplicantProfiles(ctx context.Context, filter ApplicantProfileQueryFilter) (int, error) {
	return b.storer.CountApplicantProfiles(ctx, filter)
}

// QueryApplicantProfileByID finds an admissions applicant profile by ID.
func (b *Business) QueryApplicantProfileByID(ctx context.Context, profileID uuid.UUID) (ApplicantProfile, error) {
	profile, err := b.storer.QueryApplicantProfileByID(ctx, profileID)
	if err != nil {
		return ApplicantProfile{}, fmt.Errorf("query applicant profile: profileID[%s]: %w", profileID, err)
	}

	return profile, nil
}

// QueryApplicantProfileByUserID finds an admissions applicant profile by identity user ID.
func (b *Business) QueryApplicantProfileByUserID(ctx context.Context, userID uuid.UUID) (ApplicantProfile, error) {
	profile, err := b.storer.QueryApplicantProfileByUserID(ctx, userID)
	if err != nil {
		return ApplicantProfile{}, fmt.Errorf("query applicant profile: userID[%s]: %w", userID, err)
	}

	return profile, nil
}

// QueryApplicantProfileByConstituentID finds an admissions applicant profile by constituent ID.
func (b *Business) QueryApplicantProfileByConstituentID(ctx context.Context, constituentID uuid.UUID) (ApplicantProfile, error) {
	profile, err := b.storer.QueryApplicantProfileByConstituentID(ctx, constituentID)
	if err != nil {
		return ApplicantProfile{}, fmt.Errorf("query applicant profile: constituentID[%s]: %w", constituentID, err)
	}

	return profile, nil
}

// CreateInquiry records an anonymous inquiry and links it to a matched or new constituent.
func (b *Business) CreateInquiry(ctx context.Context, ni NewInquiry) (Inquiry, error) {
	if err := validateNewInquiry(ni); err != nil {
		return Inquiry{}, err
	}

	constituent, err := b.createOrMatchInquiryConstituent(ctx, ni)
	if err != nil {
		return Inquiry{}, err
	}

	now := time.Now()
	inquiry := Inquiry{
		ID:                uuid.New(),
		ConstituentID:     constituent.ID,
		FirstName:         strings.TrimSpace(ni.FirstName),
		LastName:          strings.TrimSpace(ni.LastName),
		DateOfBirth:       ni.DateOfBirth,
		PrimaryEmail:      ni.PrimaryEmail,
		PrimaryPhone:      strings.TrimSpace(ni.PrimaryPhone),
		ProgramOfInterest: ni.ProgramOfInterest,
		TermOfInterest:    ni.TermOfInterest,
		Source:            strings.TrimSpace(ni.Source),
		UTMSource:         trimStringPtr(ni.UTMSource),
		UTMMedium:         trimStringPtr(ni.UTMMedium),
		UTMCampaign:       trimStringPtr(ni.UTMCampaign),
		Message:           trimStringPtr(ni.Message),
		Status:            InquiryStatusNew,
		DateCreated:       now,
		DateUpdated:       now,
	}

	if err := b.storer.CreateInquiry(ctx, inquiry); err != nil {
		return Inquiry{}, fmt.Errorf("create inquiry: %w", err)
	}

	return inquiry, nil
}

// QueryInquiries retrieves admissions inquiries.
func (b *Business) QueryInquiries(ctx context.Context, filter InquiryQueryFilter, orderBy order.By, page page.Page) ([]Inquiry, error) {
	inquiries, err := b.storer.QueryInquiries(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query inquiries: %w", err)
	}

	return inquiries, nil
}

// CountInquiries returns the total number of admissions inquiries.
func (b *Business) CountInquiries(ctx context.Context, filter InquiryQueryFilter) (int, error) {
	return b.storer.CountInquiries(ctx, filter)
}

// QueryInquiryByID finds an admissions inquiry by ID.
func (b *Business) QueryInquiryByID(ctx context.Context, inquiryID uuid.UUID) (Inquiry, error) {
	inquiry, err := b.storer.QueryInquiryByID(ctx, inquiryID)
	if err != nil {
		return Inquiry{}, fmt.Errorf("query inquiry: inquiryID[%s]: %w", inquiryID, err)
	}

	return inquiry, nil
}

func (b *Business) createOrMatchInquiryConstituent(ctx context.Context, ni NewInquiry) (Constituent, error) {
	matched, err := b.storer.QueryConstituentByPrimaryEmail(ctx, ni.PrimaryEmail.String())
	if err != nil && !errors.Is(err, ErrConstituentNotFound) {
		return Constituent{}, fmt.Errorf("query inquiry constituent match: %w", err)
	}

	if err == nil {
		return matched, nil
	}

	constituent, err := b.CreateConstituent(ctx, NewConstituent{
		FirstName:       ni.FirstName,
		LastName:        ni.LastName,
		DateOfBirth:     ni.DateOfBirth,
		PrimaryEmail:    ni.PrimaryEmail,
		PrimaryPhone:    ni.PrimaryPhone,
		LifecycleStage:  LifecycleStageInquiry,
		DuplicateStatus: DuplicateStatusActive,
	})
	if err != nil {
		return Constituent{}, fmt.Errorf("create inquiry constituent: %w", err)
	}

	return constituent, nil
}

// CreateLeadScoreRule adds an explainable lead scoring rule.
func (b *Business) CreateLeadScoreRule(ctx context.Context, nr NewLeadScoreRule) (LeadScoreRule, error) {
	if err := validateNewLeadScoreRule(nr); err != nil {
		return LeadScoreRule{}, err
	}

	now := time.Now()
	rule := LeadScoreRule{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(nr.Name),
		Description: trimStringPtr(nr.Description),
		Criteria:    nr.Criteria,
		Points:      nr.Points,
		Active:      nr.Active,
		Priority:    nr.Priority,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.CreateLeadScoreRule(ctx, rule); err != nil {
		return LeadScoreRule{}, fmt.Errorf("create lead score rule: %w", err)
	}

	return rule, nil
}

// UpdateLeadScoreRule replaces mutable lead scoring rule data.
func (b *Business) UpdateLeadScoreRule(ctx context.Context, rule LeadScoreRule, nr NewLeadScoreRule) (LeadScoreRule, error) {
	if err := validateNewLeadScoreRule(nr); err != nil {
		return LeadScoreRule{}, err
	}

	rule.Name = strings.TrimSpace(nr.Name)
	rule.Description = trimStringPtr(nr.Description)
	rule.Criteria = nr.Criteria
	rule.Points = nr.Points
	rule.Active = nr.Active
	rule.Priority = nr.Priority
	rule.DateUpdated = time.Now()

	if err := b.storer.UpdateLeadScoreRule(ctx, rule); err != nil {
		return LeadScoreRule{}, fmt.Errorf("update lead score rule: %w", err)
	}

	return rule, nil
}

// QueryLeadScoreRules retrieves lead scoring rules.
func (b *Business) QueryLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter, orderBy order.By, page page.Page) ([]LeadScoreRule, error) {
	rules, err := b.storer.QueryLeadScoreRules(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query lead score rules: %w", err)
	}

	return rules, nil
}

// CountLeadScoreRules returns the total number of lead score rules.
func (b *Business) CountLeadScoreRules(ctx context.Context, filter LeadScoreRuleQueryFilter) (int, error) {
	return b.storer.CountLeadScoreRules(ctx, filter)
}

// QueryLeadScoreRuleByID finds a lead scoring rule by ID.
func (b *Business) QueryLeadScoreRuleByID(ctx context.Context, ruleID uuid.UUID) (LeadScoreRule, error) {
	rule, err := b.storer.QueryLeadScoreRuleByID(ctx, ruleID)
	if err != nil {
		return LeadScoreRule{}, fmt.Errorf("query lead score rule: ruleID[%s]: %w", ruleID, err)
	}

	return rule, nil
}

// RecalculateLeadScoreForConstituent recalculates and stores a constituent's latest lead score.
func (b *Business) RecalculateLeadScoreForConstituent(ctx context.Context, constituentID uuid.UUID) (LeadScore, error) {
	constituent, err := b.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return LeadScore{}, err
	}

	active := true
	rules, err := b.QueryLeadScoreRules(ctx, LeadScoreRuleQueryFilter{Active: &active}, DefaultLeadScoreRuleOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return LeadScore{}, err
	}

	applications, err := b.QueryApplications(ctx, ApplicationQueryFilter{ConstituentID: &constituentID}, DefaultApplicationOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return LeadScore{}, err
	}

	existing, err := b.storer.QueryLeadScoreByConstituentID(ctx, constituentID)
	if err != nil && !errors.Is(err, ErrLeadScoreNotFound) {
		return LeadScore{}, fmt.Errorf("query lead score: %w", err)
	}

	now := time.Now()
	score := LeadScore{
		ID:             existing.ID,
		ConstituentID:  constituentID,
		RecalculatedAt: now,
		DateCreated:    existing.DateCreated,
		DateUpdated:    now,
	}
	if score.ID == uuid.Nil {
		score.ID = uuid.New()
		score.DateCreated = now
	}

	for _, rule := range rules {
		result := evaluateLeadScoreRule(rule, constituent, applications)
		score.Breakdown = append(score.Breakdown, result)
		if result.Matched {
			score.TotalScore += result.Points
		}
	}
	score.Band = LeadScoreBandForTotal(score.TotalScore)

	if err := b.storer.UpsertLeadScore(ctx, score); err != nil {
		return LeadScore{}, fmt.Errorf("upsert lead score: %w", err)
	}

	return score, nil
}

// QueryLeadScores retrieves lead scores.
func (b *Business) QueryLeadScores(ctx context.Context, filter LeadScoreQueryFilter, orderBy order.By, page page.Page) ([]LeadScore, error) {
	scores, err := b.storer.QueryLeadScores(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query lead scores: %w", err)
	}

	return scores, nil
}

// CountLeadScores returns the total number of lead scores.
func (b *Business) CountLeadScores(ctx context.Context, filter LeadScoreQueryFilter) (int, error) {
	return b.storer.CountLeadScores(ctx, filter)
}

// QueryLeadScoreByID finds a lead score by ID.
func (b *Business) QueryLeadScoreByID(ctx context.Context, scoreID uuid.UUID) (LeadScore, error) {
	score, err := b.storer.QueryLeadScoreByID(ctx, scoreID)
	if err != nil {
		return LeadScore{}, fmt.Errorf("query lead score: scoreID[%s]: %w", scoreID, err)
	}

	return score, nil
}

// QueryLeadScoreByConstituentID finds a lead score by constituent ID.
func (b *Business) QueryLeadScoreByConstituentID(ctx context.Context, constituentID uuid.UUID) (LeadScore, error) {
	score, err := b.storer.QueryLeadScoreByConstituentID(ctx, constituentID)
	if err != nil {
		return LeadScore{}, fmt.Errorf("query lead score: constituentID[%s]: %w", constituentID, err)
	}

	return score, nil
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for admissions access.
type Business struct {
	log        *logger.Logger
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs an admissions business API for use.
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		delegate:   delegate,
		storer:     storer,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, b.delegate, storer, b.extensions...)

	return nb, nil
}

// Health returns the current admissions context scaffold metadata.
func (b *Business) Health(ctx context.Context) (Health, error) {
	health, err := b.storer.Health(ctx)
	if err != nil {
		return Health{}, fmt.Errorf("health: %w", err)
	}

	return health, nil
}

// CreateConstituent adds a new Constituent to the admissions context.
func (b *Business) CreateConstituent(ctx context.Context, nc NewConstituent) (Constituent, error) {
	stage := nc.LifecycleStage
	if stage == "" {
		stage = LifecycleStageProspect
	}

	status := nc.DuplicateStatus
	if status == "" {
		status = DuplicateStatusActive
	}

	if err := validateRequiredConstituentFields(nc.FirstName, nc.LastName, nc.DateOfBirth, nc.PrimaryPhone); err != nil {
		return Constituent{}, err
	}

	if err := validateLifecycleStage(stage); err != nil {
		return Constituent{}, err
	}

	if err := validateDuplicateStatus(status, nc.DuplicateOfID); err != nil {
		return Constituent{}, err
	}

	now := time.Now()
	cst := Constituent{
		ID:              uuid.New(),
		FirstName:       strings.TrimSpace(nc.FirstName),
		LastName:        strings.TrimSpace(nc.LastName),
		PreferredName:   trimStringPtr(nc.PreferredName),
		MiddleName:      trimStringPtr(nc.MiddleName),
		Suffix:          trimStringPtr(nc.Suffix),
		DateOfBirth:     nc.DateOfBirth,
		PrimaryEmail:    nc.PrimaryEmail,
		PrimaryPhone:    strings.TrimSpace(nc.PrimaryPhone),
		ExternalSISID:   trimStringPtr(nc.ExternalSISID),
		LifecycleStage:  stage,
		DuplicateStatus: status,
		DuplicateOfID:   nc.DuplicateOfID,
		SISSyncedAt:     nc.SISSyncedAt,
		DateCreated:     now,
		DateUpdated:     now,
	}

	match, err := b.queryTrustedExactDuplicate(ctx, cst)
	if err != nil {
		return Constituent{}, err
	}

	if match != nil {
		if cst.ExternalSISID != nil && match.ExternalSISID != nil && *cst.ExternalSISID == *match.ExternalSISID {
			return *match, nil
		}

		cst.DuplicateStatus = DuplicateStatusDuplicateOf
		cst.DuplicateOfID = &match.ID
	}

	if err := b.storer.CreateConstituent(ctx, cst); err != nil {
		return Constituent{}, fmt.Errorf("create constituent: %w", err)
	}

	return cst, nil
}

// UpdateConstituent modifies mutable information for a Constituent.
func (b *Business) UpdateConstituent(ctx context.Context, cst Constituent, uc UpdateConstituent) (Constituent, error) {
	if uc.PreferredName != nil {
		cst.PreferredName = trimStringPtr(uc.PreferredName)
	}

	if uc.MiddleName != nil {
		cst.MiddleName = trimStringPtr(uc.MiddleName)
	}

	if uc.Suffix != nil {
		cst.Suffix = trimStringPtr(uc.Suffix)
	}

	if uc.PrimaryEmail != nil {
		cst.PrimaryEmail = *uc.PrimaryEmail
	}

	if uc.PrimaryPhone != nil {
		phone := strings.TrimSpace(*uc.PrimaryPhone)
		if phone == "" {
			return Constituent{}, ErrPrimaryPhoneRequired
		}
		cst.PrimaryPhone = phone
	}

	if uc.LifecycleStage != nil {
		if err := validateLifecycleStage(*uc.LifecycleStage); err != nil {
			return Constituent{}, err
		}

		if !canChangeLifecycleStage(cst.LifecycleStage, *uc.LifecycleStage) {
			return Constituent{}, ErrInvalidLifecycleChange
		}

		cst.LifecycleStage = *uc.LifecycleStage
	}

	if uc.DuplicateStatus != nil || uc.DuplicateOfID != nil {
		status := cst.DuplicateStatus
		if uc.DuplicateStatus != nil {
			status = *uc.DuplicateStatus
		}

		duplicateOfID := cst.DuplicateOfID
		if uc.DuplicateOfID != nil {
			duplicateOfID = uc.DuplicateOfID
		}

		if err := validateDuplicateStatus(status, duplicateOfID); err != nil {
			return Constituent{}, err
		}

		cst.DuplicateStatus = status
		cst.DuplicateOfID = duplicateOfID
	}

	if uc.SISSyncedAt != nil {
		cst.SISSyncedAt = uc.SISSyncedAt
	}

	cst.DateUpdated = time.Now()

	if err := b.storer.UpdateConstituent(ctx, cst); err != nil {
		return Constituent{}, fmt.Errorf("update constituent: %w", err)
	}

	return cst, nil
}

// QueryConstituents retrieves a list of existing constituents.
func (b *Business) QueryConstituents(ctx context.Context, filter ConstituentQueryFilter, orderBy order.By, page page.Page) ([]Constituent, error) {
	constituents, err := b.storer.QueryConstituents(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query constituents: %w", err)
	}

	return constituents, nil
}

// CountConstituents returns the total number of constituents.
func (b *Business) CountConstituents(ctx context.Context, filter ConstituentQueryFilter) (int, error) {
	return b.storer.CountConstituents(ctx, filter)
}

// QueryConstituentByID finds a Constituent by ID.
func (b *Business) QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: constituentID[%s]: %w", constituentID, err)
	}

	return cst, nil
}

// QueryConstituentByPrimaryEmail finds a Constituent by primary email.
func (b *Business) QueryConstituentByPrimaryEmail(ctx context.Context, email string) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByPrimaryEmail(ctx, email)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: primaryEmail[%s]: %w", email, err)
	}

	return cst, nil
}

// QueryConstituentByExternalSISID finds a Constituent by SIS ID.
func (b *Business) QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (Constituent, error) {
	cst, err := b.storer.QueryConstituentByExternalSISID(ctx, externalSISID)
	if err != nil {
		return Constituent{}, fmt.Errorf("query constituent: externalSISID[%s]: %w", externalSISID, err)
	}

	return cst, nil
}

func (b *Business) queryTrustedExactDuplicate(ctx context.Context, cst Constituent) (*Constituent, error) {
	if cst.ExternalSISID != nil {
		match, err := b.storer.QueryConstituentByExternalSISID(ctx, *cst.ExternalSISID)
		if err != nil && !errors.Is(err, ErrConstituentNotFound) {
			return nil, fmt.Errorf("query external sis duplicate: %w", err)
		}

		if err == nil && match.ID != cst.ID {
			return &match, nil
		}
	}

	match, err := b.storer.QueryConstituentByPrimaryEmail(ctx, cst.PrimaryEmail.String())
	if err != nil && !errors.Is(err, ErrConstituentNotFound) {
		return nil, fmt.Errorf("query email duplicate: %w", err)
	}

	if err == nil && match.ID != cst.ID {
		return &match, nil
	}

	return nil, nil
}

// UpsertProgram creates or updates SIS-owned Program reference data for sync/import paths.
func (b *Business) UpsertProgram(ctx context.Context, up UpsertProgram) (Program, error) {
	now := time.Now()
	id := uuid.New()
	if up.ID != nil {
		id = *up.ID
	}

	prg := Program{
		ID:            id,
		ExternalSISID: up.ExternalSISID,
		Name:          up.Name,
		Code:          up.Code,
		Description:   up.Description,
		DegreeLevel:   up.DegreeLevel,
		Active:        up.Active,
		SyncedAt:      up.SyncedAt,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := b.storer.UpsertProgram(ctx, prg); err != nil {
		return Program{}, fmt.Errorf("upsert program: %w", err)
	}

	return b.QueryProgramByExternalSISID(ctx, up.ExternalSISID)
}

// QueryPrograms retrieves a list of existing Program reference records.
func (b *Business) QueryPrograms(ctx context.Context, filter ProgramQueryFilter, orderBy order.By, page page.Page) ([]Program, error) {
	programs, err := b.storer.QueryPrograms(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query programs: %w", err)
	}

	return programs, nil
}

// CountPrograms returns the total number of Program reference records.
func (b *Business) CountPrograms(ctx context.Context, filter ProgramQueryFilter) (int, error) {
	return b.storer.CountPrograms(ctx, filter)
}

// QueryProgramByID finds a Program by ID.
func (b *Business) QueryProgramByID(ctx context.Context, programID uuid.UUID) (Program, error) {
	program, err := b.storer.QueryProgramByID(ctx, programID)
	if err != nil {
		return Program{}, fmt.Errorf("query program: programID[%s]: %w", programID, err)
	}

	return program, nil
}

// QueryProgramByExternalSISID finds a Program by immutable SIS ID.
func (b *Business) QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (Program, error) {
	program, err := b.storer.QueryProgramByExternalSISID(ctx, externalSISID)
	if err != nil {
		return Program{}, fmt.Errorf("query program: externalSISID[%s]: %w", externalSISID, err)
	}

	return program, nil
}

// UpsertAcademicTerm creates or updates SIS-owned AcademicTerm reference data for sync/import paths.
func (b *Business) UpsertAcademicTerm(ctx context.Context, up UpsertAcademicTerm) (AcademicTerm, error) {
	if !up.StartDate.Before(up.EndDate) {
		return AcademicTerm{}, ErrInvalidTermDateRange
	}

	if up.ApplicationStartDate != nil && up.ApplicationDeadline != nil && up.ApplicationDeadline.Before(*up.ApplicationStartDate) {
		return AcademicTerm{}, ErrInvalidApplicationWindow
	}

	now := time.Now()
	id := uuid.New()
	if up.ID != nil {
		id = *up.ID
	}

	term := AcademicTerm{
		ID:                   id,
		ExternalSISID:        up.ExternalSISID,
		Name:                 up.Name,
		Code:                 up.Code,
		TermType:             up.TermType,
		StartDate:            up.StartDate,
		EndDate:              up.EndDate,
		ApplicationStartDate: up.ApplicationStartDate,
		ApplicationDeadline:  up.ApplicationDeadline,
		Active:               up.Active,
		SyncedAt:             up.SyncedAt,
		DateCreated:          now,
		DateUpdated:          now,
	}

	if err := b.storer.UpsertAcademicTerm(ctx, term); err != nil {
		return AcademicTerm{}, fmt.Errorf("upsert academic term: %w", err)
	}

	return b.QueryAcademicTermByExternalSISID(ctx, up.ExternalSISID)
}

// QueryAcademicTerms retrieves a list of existing AcademicTerm reference records.
func (b *Business) QueryAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]AcademicTerm, error) {
	terms, err := b.storer.QueryAcademicTerms(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query academic terms: %w", err)
	}

	return terms, nil
}

// CountAcademicTerms returns the total number of AcademicTerm reference records.
func (b *Business) CountAcademicTerms(ctx context.Context, filter AcademicTermQueryFilter) (int, error) {
	return b.storer.CountAcademicTerms(ctx, filter)
}

// QueryAcademicTermByID finds an AcademicTerm by ID.
func (b *Business) QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (AcademicTerm, error) {
	term, err := b.storer.QueryAcademicTermByID(ctx, termID)
	if err != nil {
		return AcademicTerm{}, fmt.Errorf("query academic term: termID[%s]: %w", termID, err)
	}

	return term, nil
}

// QueryAcademicTermByExternalSISID finds an AcademicTerm by immutable SIS ID.
func (b *Business) QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (AcademicTerm, error) {
	term, err := b.storer.QueryAcademicTermByExternalSISID(ctx, externalSISID)
	if err != nil {
		return AcademicTerm{}, fmt.Errorf("query academic term: externalSISID[%s]: %w", externalSISID, err)
	}

	return term, nil
}

// CreateDuplicateReview adds a possible duplicate pair to the staff review queue.
func (b *Business) CreateDuplicateReview(ctx context.Context, nr NewDuplicateReview) (DuplicateReview, error) {
	if err := validateNewDuplicateReview(nr); err != nil {
		return DuplicateReview{}, err
	}

	now := time.Now()
	review := DuplicateReview{
		ID:                     uuid.New(),
		SourceConstituentID:    nr.SourceConstituentID,
		CandidateConstituentID: nr.CandidateConstituentID,
		MatchType:              nr.MatchType,
		MatchScore:             nr.MatchScore,
		MatchReason:            strings.TrimSpace(nr.MatchReason),
		Status:                 DuplicateReviewStatusPending,
		DateCreated:            now,
		DateUpdated:            now,
	}

	if err := b.storer.CreateDuplicateReview(ctx, review); err != nil {
		return DuplicateReview{}, fmt.Errorf("create duplicate review: %w", err)
	}

	return review, nil
}

// ResolveDuplicateReview records a staff decision for a pending duplicate review.
func (b *Business) ResolveDuplicateReview(ctx context.Context, review DuplicateReview, rr ResolveDuplicateReview) (DuplicateReview, error) {
	if review.Status != DuplicateReviewStatusPending && review.Status != DuplicateReviewStatusDeferred {
		return DuplicateReview{}, ErrDuplicateReviewResolved
	}

	if rr.ActorID == uuid.Nil {
		return DuplicateReview{}, ErrResolutionActorRequired
	}

	status, err := statusForResolution(rr.Resolution)
	if err != nil {
		return DuplicateReview{}, err
	}

	now := time.Now()
	review.Status = status
	review.ResolvedBy = &rr.ActorID
	review.ResolvedAt = &now
	review.ResolutionNote = trimStringPtr(rr.Note)
	review.DateUpdated = now

	if err := b.storer.UpdateDuplicateReview(ctx, review); err != nil {
		return DuplicateReview{}, fmt.Errorf("update duplicate review: %w", err)
	}

	if rr.Resolution == DuplicateReviewResolutionLink || rr.Resolution == DuplicateReviewResolutionMerge {
		source, err := b.storer.QueryConstituentByID(ctx, review.SourceConstituentID)
		if err != nil {
			return DuplicateReview{}, fmt.Errorf("query source constituent: %w", err)
		}

		candidate, err := b.storer.QueryConstituentByID(ctx, review.CandidateConstituentID)
		if err != nil {
			return DuplicateReview{}, fmt.Errorf("query candidate constituent: %w", err)
		}

		duplicateStatus := DuplicateStatusDuplicateOf
		if rr.Resolution == DuplicateReviewResolutionMerge {
			duplicateStatus = DuplicateStatusMerged
		}

		source.DuplicateStatus = duplicateStatus
		source.DuplicateOfID = &candidate.ID
		source.DateUpdated = now

		if err := b.storer.UpdateConstituent(ctx, source); err != nil {
			return DuplicateReview{}, fmt.Errorf("update source duplicate link: %w", err)
		}
	}

	return review, nil
}

// QueryDuplicateReviews retrieves a list of existing duplicate reviews.
func (b *Business) QueryDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]DuplicateReview, error) {
	reviews, err := b.storer.QueryDuplicateReviews(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query duplicate reviews: %w", err)
	}

	return reviews, nil
}

// CountDuplicateReviews returns the total number of duplicate reviews.
func (b *Business) CountDuplicateReviews(ctx context.Context, filter DuplicateReviewQueryFilter) (int, error) {
	return b.storer.CountDuplicateReviews(ctx, filter)
}

// QueryDuplicateReviewByID finds a DuplicateReview by ID.
func (b *Business) QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (DuplicateReview, error) {
	review, err := b.storer.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return DuplicateReview{}, fmt.Errorf("query duplicate review: reviewID[%s]: %w", reviewID, err)
	}

	return review, nil
}

// CreateApplication adds a draft application while enforcing active application uniqueness.
func (b *Business) CreateApplication(ctx context.Context, na NewApplication) (Application, error) {
	if err := validateNewApplication(na); err != nil {
		return Application{}, err
	}

	if _, err := b.storer.QueryConstituentByID(ctx, na.ConstituentID); err != nil {
		return Application{}, fmt.Errorf("query constituent: %w", err)
	}

	program, err := b.storer.QueryProgramByID(ctx, na.ProgramID)
	if err != nil {
		return Application{}, fmt.Errorf("query program: %w", err)
	}
	if !program.Active {
		return Application{}, ErrInactiveProgram
	}

	term, err := b.storer.QueryAcademicTermByID(ctx, na.AcademicTermID)
	if err != nil {
		return Application{}, fmt.Errorf("query academic term: %w", err)
	}
	if !term.Active {
		return Application{}, ErrInactiveAcademicTerm
	}

	if _, err := b.storer.QueryActiveApplicationByTuple(ctx, na.ConstituentID, na.AcademicTermID, na.ProgramID); err == nil {
		return Application{}, ErrDuplicateApplication
	} else if !errors.Is(err, ErrApplicationNotFound) {
		return Application{}, fmt.Errorf("query active application: %w", err)
	}

	now := time.Now()
	app := Application{
		ID:                 uuid.New(),
		ConstituentID:      na.ConstituentID,
		ProgramID:          na.ProgramID,
		AcademicTermID:     na.AcademicTermID,
		ApplicationType:    na.ApplicationType,
		Status:             ApplicationStatusDraft,
		AssignedReviewerID: na.AssignedReviewerID,
		DateCreated:        now,
		DateUpdated:        now,
	}

	if err := b.storer.CreateApplication(ctx, app); err != nil {
		return Application{}, fmt.Errorf("create application: %w", err)
	}

	return app, nil
}

// QueryApplications retrieves a list of existing applications.
func (b *Business) QueryApplications(ctx context.Context, filter ApplicationQueryFilter, orderBy order.By, page page.Page) ([]Application, error) {
	applications, err := b.storer.QueryApplications(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query applications: %w", err)
	}

	return applications, nil
}

// CountApplications returns the total number of applications.
func (b *Business) CountApplications(ctx context.Context, filter ApplicationQueryFilter) (int, error) {
	return b.storer.CountApplications(ctx, filter)
}

// QueryApplicationByID finds an Application by ID.
func (b *Business) QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (Application, error) {
	app, err := b.storer.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return Application{}, fmt.Errorf("query application: applicationID[%s]: %w", applicationID, err)
	}

	return app, nil
}

// CreateApplicationFormTemplate adds a configurable application form template.
func (b *Business) CreateApplicationFormTemplate(ctx context.Context, nt NewApplicationFormTemplate) (ApplicationFormTemplate, error) {
	if err := validateNewApplicationFormTemplate(nt); err != nil {
		return ApplicationFormTemplate{}, err
	}

	program, err := b.storer.QueryProgramByID(ctx, nt.ProgramID)
	if err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("query program: %w", err)
	}
	if !program.Active {
		return ApplicationFormTemplate{}, ErrInactiveProgram
	}

	term, err := b.storer.QueryAcademicTermByID(ctx, nt.AcademicTermID)
	if err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("query academic term: %w", err)
	}
	if !term.Active {
		return ApplicationFormTemplate{}, ErrInactiveAcademicTerm
	}

	now := time.Now()
	template := ApplicationFormTemplate{
		ID:              uuid.New(),
		ProgramID:       nt.ProgramID,
		AcademicTermID:  nt.AcademicTermID,
		ApplicationType: nt.ApplicationType,
		Name:            strings.TrimSpace(nt.Name),
		Description:     trimStringPtr(nt.Description),
		Version:         1,
		RequiredFields:  normalizeApplicationFormFields(nt.RequiredFields),
		ChecklistItems:  normalizeChecklistTemplateItems(nt.ChecklistItems),
		Active:          nt.Active,
		Priority:        nt.Priority,
		DateCreated:     now,
		DateUpdated:     now,
	}

	if err := b.storer.CreateApplicationFormTemplate(ctx, template); err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("create application form template: %w", err)
	}

	return template, nil
}

// UpdateApplicationFormTemplate creates a new template version and updates the mutable template configuration.
func (b *Business) UpdateApplicationFormTemplate(ctx context.Context, template ApplicationFormTemplate, nt NewApplicationFormTemplate) (ApplicationFormTemplate, error) {
	if err := validateNewApplicationFormTemplate(nt); err != nil {
		return ApplicationFormTemplate{}, err
	}

	program, err := b.storer.QueryProgramByID(ctx, nt.ProgramID)
	if err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("query program: %w", err)
	}
	if !program.Active {
		return ApplicationFormTemplate{}, ErrInactiveProgram
	}

	term, err := b.storer.QueryAcademicTermByID(ctx, nt.AcademicTermID)
	if err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("query academic term: %w", err)
	}
	if !term.Active {
		return ApplicationFormTemplate{}, ErrInactiveAcademicTerm
	}

	template.ProgramID = nt.ProgramID
	template.AcademicTermID = nt.AcademicTermID
	template.ApplicationType = nt.ApplicationType
	template.Name = strings.TrimSpace(nt.Name)
	template.Description = trimStringPtr(nt.Description)
	template.Version++
	template.RequiredFields = normalizeApplicationFormFields(nt.RequiredFields)
	template.ChecklistItems = normalizeChecklistTemplateItems(nt.ChecklistItems)
	template.Active = nt.Active
	template.Priority = nt.Priority
	template.DateUpdated = time.Now()

	if err := b.storer.UpdateApplicationFormTemplate(ctx, template); err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("update application form template: %w", err)
	}

	return template, nil
}

// QueryApplicationFormTemplates retrieves configurable application form templates.
func (b *Business) QueryApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter, orderBy order.By, page page.Page) ([]ApplicationFormTemplate, error) {
	templates, err := b.storer.QueryApplicationFormTemplates(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query application form templates: %w", err)
	}

	return templates, nil
}

// CountApplicationFormTemplates returns the total number of application form templates.
func (b *Business) CountApplicationFormTemplates(ctx context.Context, filter ApplicationFormTemplateQueryFilter) (int, error) {
	return b.storer.CountApplicationFormTemplates(ctx, filter)
}

// QueryApplicationFormTemplateByID finds an application form template by ID.
func (b *Business) QueryApplicationFormTemplateByID(ctx context.Context, templateID uuid.UUID) (ApplicationFormTemplate, error) {
	template, err := b.storer.QueryApplicationFormTemplateByID(ctx, templateID)
	if err != nil {
		return ApplicationFormTemplate{}, fmt.Errorf("query application form template: templateID[%s]: %w", templateID, err)
	}

	return template, nil
}

// TransitionApplicationStatus changes an Application status and records immutable transition history.
func (b *Business) TransitionApplicationStatus(ctx context.Context, app Application, nt NewApplicationTransition) (Application, ApplicationTransition, error) {
	if nt.ActorID == uuid.Nil {
		return Application{}, ApplicationTransition{}, ErrApplicationActorRequired
	}

	if err := validateApplicationStatus(nt.ToStatus); err != nil {
		return Application{}, ApplicationTransition{}, err
	}

	if !canTransitionApplicationStatus(app.Status, nt.ToStatus) {
		return Application{}, ApplicationTransition{}, ErrInvalidApplicationTransition
	}

	now := time.Now()
	transition := ApplicationTransition{
		ID:            uuid.New(),
		ApplicationID: app.ID,
		FromStatus:    app.Status,
		ToStatus:      nt.ToStatus,
		ActorID:       nt.ActorID,
		Reason:        trimStringPtr(nt.Reason),
		Note:          trimStringPtr(nt.Note),
		Metadata:      nt.Metadata,
		DateCreated:   now,
	}

	app.Status = nt.ToStatus
	app.DateUpdated = now
	if nt.ToStatus == ApplicationStatusSubmitted && app.SubmittedAt == nil {
		app.SubmittedAt = &now
	}

	if err := b.storer.UpdateApplication(ctx, app); err != nil {
		return Application{}, ApplicationTransition{}, fmt.Errorf("update application: %w", err)
	}

	if err := b.storer.CreateApplicationTransition(ctx, transition); err != nil {
		return Application{}, ApplicationTransition{}, fmt.Errorf("create application transition: %w", err)
	}

	return app, transition, nil
}

// QueryApplicationTransitions retrieves a list of application transition history records.
func (b *Business) QueryApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter, orderBy order.By, page page.Page) ([]ApplicationTransition, error) {
	transitions, err := b.storer.QueryApplicationTransitions(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query application transitions: %w", err)
	}

	return transitions, nil
}

// CountApplicationTransitions returns the total number of application transition records.
func (b *Business) CountApplicationTransitions(ctx context.Context, filter ApplicationTransitionQueryFilter) (int, error) {
	return b.storer.CountApplicationTransitions(ctx, filter)
}

// CreateChecklistItem adds a document requirement to an application checklist.
func (b *Business) CreateChecklistItem(ctx context.Context, ni NewChecklistItem) (ChecklistItem, error) {
	if err := validateNewChecklistItem(ni); err != nil {
		return ChecklistItem{}, err
	}

	if _, err := b.storer.QueryApplicationByID(ctx, ni.ApplicationID); err != nil {
		return ChecklistItem{}, fmt.Errorf("query application: %w", err)
	}

	now := time.Now()
	item := ChecklistItem{
		ID:            uuid.New(),
		ApplicationID: ni.ApplicationID,
		ItemKey:       strings.TrimSpace(ni.ItemKey),
		DocumentName:  strings.TrimSpace(ni.DocumentName),
		Description:   trimStringPtr(ni.Description),
		Required:      ni.Required,
		Status:        DocumentStatusPendingReview,
		DisplayOrder:  ni.DisplayOrder,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := b.storer.CreateChecklistItem(ctx, item); err != nil {
		return ChecklistItem{}, fmt.Errorf("create checklist item: %w", err)
	}

	return item, nil
}

// UpdateChecklistItem updates a document requirement while preserving its review status.
func (b *Business) UpdateChecklistItem(ctx context.Context, item ChecklistItem, ni NewChecklistItem) (ChecklistItem, error) {
	if err := validateNewChecklistItem(ni); err != nil {
		return ChecklistItem{}, err
	}

	if _, err := b.storer.QueryApplicationByID(ctx, ni.ApplicationID); err != nil {
		return ChecklistItem{}, fmt.Errorf("query application: %w", err)
	}

	item.ApplicationID = ni.ApplicationID
	item.ItemKey = strings.TrimSpace(ni.ItemKey)
	item.DocumentName = strings.TrimSpace(ni.DocumentName)
	item.Description = trimStringPtr(ni.Description)
	item.Required = ni.Required
	item.DisplayOrder = ni.DisplayOrder
	item.DateUpdated = time.Now()

	if err := b.storer.UpdateChecklistItem(ctx, item); err != nil {
		return ChecklistItem{}, fmt.Errorf("update checklist item: %w", err)
	}

	return item, nil
}

// QueryChecklistItems retrieves application checklist requirements.
func (b *Business) QueryChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter, orderBy order.By, page page.Page) ([]ChecklistItem, error) {
	items, err := b.storer.QueryChecklistItems(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query checklist items: %w", err)
	}

	return items, nil
}

// CountChecklistItems returns the total number of checklist requirements.
func (b *Business) CountChecklistItems(ctx context.Context, filter ChecklistItemQueryFilter) (int, error) {
	return b.storer.CountChecklistItems(ctx, filter)
}

// QueryChecklistItemByID finds a checklist requirement by ID.
func (b *Business) QueryChecklistItemByID(ctx context.Context, itemID uuid.UUID) (ChecklistItem, error) {
	item, err := b.storer.QueryChecklistItemByID(ctx, itemID)
	if err != nil {
		return ChecklistItem{}, fmt.Errorf("query checklist item: itemID[%s]: %w", itemID, err)
	}

	return item, nil
}

// CreateDocument records uploaded document metadata for a checklist item.
func (b *Business) CreateDocument(ctx context.Context, nd NewDocument) (Document, error) {
	if err := validateNewDocument(nd); err != nil {
		return Document{}, err
	}

	item, err := b.storer.QueryChecklistItemByID(ctx, nd.ChecklistItemID)
	if err != nil {
		return Document{}, fmt.Errorf("query checklist item: %w", err)
	}
	if item.ApplicationID != nd.ApplicationID {
		return Document{}, ErrApplicationNotFound
	}

	now := time.Now()
	document := Document{
		ID:              uuid.New(),
		ApplicationID:   nd.ApplicationID,
		ChecklistItemID: nd.ChecklistItemID,
		FileName:        strings.TrimSpace(nd.FileName),
		ContentType:     strings.TrimSpace(nd.ContentType),
		SizeBytes:       nd.SizeBytes,
		StorageKey:      strings.TrimSpace(nd.StorageKey),
		Status:          DocumentStatusPendingReview,
		UploadedByID:    nd.UploadedByID,
		UploadedAt:      now,
		DateCreated:     now,
		DateUpdated:     now,
	}

	item.Status = DocumentStatusPendingReview
	item.DateUpdated = now

	if err := b.storer.CreateDocument(ctx, document); err != nil {
		return Document{}, fmt.Errorf("create document: %w", err)
	}

	if err := b.storer.UpdateChecklistItem(ctx, item); err != nil {
		return Document{}, fmt.Errorf("update checklist item: %w", err)
	}

	return document, nil
}

// VerifyDocument records reviewer action for uploaded document metadata.
func (b *Business) VerifyDocument(ctx context.Context, document Document, nv NewDocumentVerification) (Document, error) {
	if err := validateNewDocumentVerification(nv); err != nil {
		return Document{}, err
	}

	item, err := b.storer.QueryChecklistItemByID(ctx, document.ChecklistItemID)
	if err != nil {
		return Document{}, fmt.Errorf("query checklist item: %w", err)
	}

	now := time.Now()
	document.Status = nv.Status
	document.ReviewerID = &nv.ReviewerID
	document.ReviewerNotes = trimStringPtr(nv.ReviewerNotes)
	document.ReviewedAt = &now
	document.DateUpdated = now

	item.Status = nv.Status
	item.DateUpdated = now

	if err := b.storer.UpdateDocument(ctx, document); err != nil {
		return Document{}, fmt.Errorf("update document: %w", err)
	}

	if err := b.storer.UpdateChecklistItem(ctx, item); err != nil {
		return Document{}, fmt.Errorf("update checklist item: %w", err)
	}

	return document, nil
}

// QueryDocuments retrieves uploaded document metadata.
func (b *Business) QueryDocuments(ctx context.Context, filter DocumentQueryFilter, orderBy order.By, page page.Page) ([]Document, error) {
	documents, err := b.storer.QueryDocuments(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query documents: %w", err)
	}

	return documents, nil
}

// CountDocuments returns the total number of uploaded document metadata records.
func (b *Business) CountDocuments(ctx context.Context, filter DocumentQueryFilter) (int, error) {
	return b.storer.CountDocuments(ctx, filter)
}

// QueryDocumentByID finds uploaded document metadata by ID.
func (b *Business) QueryDocumentByID(ctx context.Context, documentID uuid.UUID) (Document, error) {
	document, err := b.storer.QueryDocumentByID(ctx, documentID)
	if err != nil {
		return Document{}, fmt.Errorf("query document: documentID[%s]: %w", documentID, err)
	}

	return document, nil
}

// CreateSyncJob schedules or starts a SIS batch reconciliation run.
func (b *Business) CreateSyncJob(ctx context.Context, nj NewSyncJob) (SyncJob, error) {
	if err := validateNewSyncJob(nj); err != nil {
		return SyncJob{}, err
	}

	now := time.Now()
	job := SyncJob{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(nj.Name),
		Status:      nj.Status,
		Direction:   nj.Direction,
		StartedAt:   nj.StartedAt,
		CreatedByID: nj.CreatedByID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.CreateSyncJob(ctx, job); err != nil {
		return SyncJob{}, fmt.Errorf("create sync job: %w", err)
	}

	return job, nil
}

// UpdateSyncJob records the outcome and retry state for a SIS batch reconciliation run.
func (b *Business) UpdateSyncJob(ctx context.Context, job SyncJob, uj UpdateSyncJob) (SyncJob, error) {
	if err := validateSyncJobStatus(uj.Status); err != nil {
		return SyncJob{}, err
	}

	job.Status = uj.Status
	job.CompletedAt = uj.CompletedAt
	job.RecordsPulled = uj.RecordsPulled
	job.RecordsPushed = uj.RecordsPushed
	job.EventsRequeued = uj.EventsRequeued
	job.FailureReason = trimStringPtr(uj.FailureReason)
	job.Retryable = uj.Retryable
	job.DateUpdated = time.Now()

	if err := b.storer.UpdateSyncJob(ctx, job); err != nil {
		return SyncJob{}, fmt.Errorf("update sync job: %w", err)
	}

	return job, nil
}

// QuerySyncJobs retrieves SIS batch reconciliation runs.
func (b *Business) QuerySyncJobs(ctx context.Context, filter SyncJobQueryFilter, orderBy order.By, page page.Page) ([]SyncJob, error) {
	jobs, err := b.storer.QuerySyncJobs(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query sync jobs: %w", err)
	}

	return jobs, nil
}

// CountSyncJobs returns the total number of SIS batch reconciliation runs.
func (b *Business) CountSyncJobs(ctx context.Context, filter SyncJobQueryFilter) (int, error) {
	return b.storer.CountSyncJobs(ctx, filter)
}

// QuerySyncJobByID finds a SIS batch reconciliation run by ID.
func (b *Business) QuerySyncJobByID(ctx context.Context, jobID uuid.UUID) (SyncJob, error) {
	job, err := b.storer.QuerySyncJobByID(ctx, jobID)
	if err != nil {
		return SyncJob{}, fmt.Errorf("query sync job: jobID[%s]: %w", jobID, err)
	}

	return job, nil
}

// EnqueueSyncEvent records a selected real-time SIS sync event for queue processing.
func (b *Business) EnqueueSyncEvent(ctx context.Context, ne NewSyncEvent) (SyncEvent, error) {
	if err := validateNewSyncEvent(ne); err != nil {
		return SyncEvent{}, err
	}

	now := time.Now()
	event := SyncEvent{
		ID:           uuid.New(),
		JobID:        ne.JobID,
		EventType:    ne.EventType,
		Status:       SyncEventStatusQueued,
		Direction:    ne.Direction,
		ResourceType: strings.TrimSpace(ne.ResourceType),
		ResourceID:   ne.ResourceID,
		PayloadHash:  strings.TrimSpace(ne.PayloadHash),
		AuditMessage: strings.TrimSpace(ne.AuditMessage),
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.CreateSyncEvent(ctx, event); err != nil {
		return SyncEvent{}, fmt.Errorf("create sync event: %w", err)
	}

	return event, nil
}

// UpdateSyncEvent records queue processing status, retry scheduling, and failure visibility.
func (b *Business) UpdateSyncEvent(ctx context.Context, event SyncEvent, ue UpdateSyncEvent) (SyncEvent, error) {
	if err := validateSyncEventStatus(ue.Status); err != nil {
		return SyncEvent{}, err
	}

	event.Status = ue.Status
	event.Attempts = ue.Attempts
	event.NextRetryAt = ue.NextRetryAt
	event.FailureReason = trimStringPtr(ue.FailureReason)
	event.AuditMessage = strings.TrimSpace(ue.AuditMessage)
	event.DateUpdated = time.Now()

	if err := b.storer.UpdateSyncEvent(ctx, event); err != nil {
		return SyncEvent{}, fmt.Errorf("update sync event: %w", err)
	}

	return event, nil
}

// QuerySyncEvents retrieves selected real-time SIS sync events.
func (b *Business) QuerySyncEvents(ctx context.Context, filter SyncEventQueryFilter, orderBy order.By, page page.Page) ([]SyncEvent, error) {
	events, err := b.storer.QuerySyncEvents(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query sync events: %w", err)
	}

	return events, nil
}

// CountSyncEvents returns the total number of selected real-time SIS sync events.
func (b *Business) CountSyncEvents(ctx context.Context, filter SyncEventQueryFilter) (int, error) {
	return b.storer.CountSyncEvents(ctx, filter)
}

// QuerySyncEventByID finds a selected real-time SIS sync event by ID.
func (b *Business) QuerySyncEventByID(ctx context.Context, eventID uuid.UUID) (SyncEvent, error) {
	event, err := b.storer.QuerySyncEventByID(ctx, eventID)
	if err != nil {
		return SyncEvent{}, fmt.Errorf("query sync event: eventID[%s]: %w", eventID, err)
	}

	return event, nil
}

func validateRequiredConstituentFields(firstName string, lastName string, dob time.Time, primaryPhone string) error {
	if strings.TrimSpace(firstName) == "" {
		return ErrFirstNameRequired
	}

	if strings.TrimSpace(lastName) == "" {
		return ErrLastNameRequired
	}

	if dob.IsZero() {
		return ErrDateOfBirthRequired
	}

	if dob.After(time.Now()) {
		return ErrDateOfBirthInFuture
	}

	if strings.TrimSpace(primaryPhone) == "" {
		return ErrPrimaryPhoneRequired
	}

	return nil
}

func validateNewStaffProfile(np NewStaffProfile) error {
	if np.UserID == uuid.Nil {
		return ErrStaffProfileUserRequired
	}

	if len(np.Roles) == 0 {
		return ErrStaffProfileRoleRequired
	}

	for _, role := range np.Roles {
		if err := validateAdmissionsRole(role); err != nil {
			return err
		}
	}

	return nil
}

func validateNewApplicantProfile(np NewApplicantProfile) error {
	if np.UserID == uuid.Nil {
		return ErrApplicantProfileUserRequired
	}

	if np.ConstituentID == uuid.Nil {
		return ErrConstituentIDRequired
	}

	return nil
}

func validateNewInquiry(ni NewInquiry) error {
	if err := validateRequiredConstituentFields(ni.FirstName, ni.LastName, ni.DateOfBirth, ni.PrimaryPhone); err != nil {
		return err
	}

	if strings.TrimSpace(ni.Source) == "" {
		return ErrInquirySourceRequired
	}

	return nil
}

func validateInquiryStatus(status InquiryStatus) error {
	switch status {
	case InquiryStatusNew,
		InquiryStatusContacted,
		InquiryStatusConverted,
		InquiryStatusClosed:
		return nil
	default:
		return ErrInvalidInquiryStatus
	}
}

func validateNewApplicationFormTemplate(nt NewApplicationFormTemplate) error {
	if nt.ProgramID == uuid.Nil {
		return ErrProgramIDRequired
	}

	if nt.AcademicTermID == uuid.Nil {
		return ErrAcademicTermIDRequired
	}

	if err := validateApplicationType(nt.ApplicationType); err != nil {
		return err
	}

	if strings.TrimSpace(nt.Name) == "" {
		return ErrFormTemplateNameRequired
	}

	if len(nt.RequiredFields) == 0 {
		return ErrFormTemplateFieldsRequired
	}

	for _, field := range nt.RequiredFields {
		if err := validateApplicationFormField(field); err != nil {
			return err
		}
	}

	for _, item := range nt.ChecklistItems {
		if err := validateChecklistTemplateItem(item); err != nil {
			return err
		}
	}

	if nt.Priority < 0 {
		return ErrFormTemplatePriorityInvalid
	}

	return nil
}

func validateApplicationFormField(field ApplicationFormField) error {
	if strings.TrimSpace(field.FieldName) == "" || strings.TrimSpace(field.FieldType) == "" || field.DisplayOrder < 0 {
		return ErrFormTemplateFieldInvalid
	}

	return nil
}

func validateChecklistTemplateItem(item ApplicationChecklistTemplateItem) error {
	if strings.TrimSpace(item.ItemKey) == "" || strings.TrimSpace(item.DocumentName) == "" || item.DisplayOrder < 0 {
		return ErrFormTemplateChecklistInvalid
	}

	return nil
}

func validateNewChecklistItem(ni NewChecklistItem) error {
	if ni.ApplicationID == uuid.Nil {
		return ErrApplicationNotFound
	}

	if strings.TrimSpace(ni.ItemKey) == "" {
		return ErrChecklistItemKeyRequired
	}

	if strings.TrimSpace(ni.DocumentName) == "" {
		return ErrChecklistItemNameRequired
	}

	if ni.DisplayOrder < 0 {
		return ErrChecklistItemOrderInvalid
	}

	return nil
}

func validateNewDocument(nd NewDocument) error {
	if nd.ApplicationID == uuid.Nil {
		return ErrApplicationNotFound
	}

	if nd.ChecklistItemID == uuid.Nil {
		return ErrChecklistItemNotFound
	}

	if strings.TrimSpace(nd.FileName) == "" {
		return ErrDocumentFileNameRequired
	}

	if strings.TrimSpace(nd.ContentType) == "" {
		return ErrDocumentContentTypeRequired
	}

	if nd.SizeBytes <= 0 {
		return ErrDocumentSizeInvalid
	}

	if strings.TrimSpace(nd.StorageKey) == "" {
		return ErrDocumentStorageKeyRequired
	}

	if nd.UploadedByID == uuid.Nil {
		return ErrDocumentUploaderRequired
	}

	return nil
}

func validateNewDocumentVerification(nv NewDocumentVerification) error {
	if nv.ReviewerID == uuid.Nil {
		return ErrDocumentReviewerRequired
	}

	if !isReviewDocumentStatus(nv.Status) {
		return ErrDocumentStatusNotReviewable
	}

	return nil
}

func validateNewSyncJob(nj NewSyncJob) error {
	if strings.TrimSpace(nj.Name) == "" {
		return ErrSyncJobNameRequired
	}

	if err := validateSyncDirection(nj.Direction); err != nil {
		return err
	}

	return validateSyncJobStatus(nj.Status)
}

func validateNewSyncEvent(ne NewSyncEvent) error {
	if err := validateSyncEventType(ne.EventType); err != nil {
		return err
	}

	if err := validateSyncDirection(ne.Direction); err != nil {
		return err
	}

	if strings.TrimSpace(ne.ResourceType) == "" || ne.ResourceID == uuid.Nil {
		return ErrSyncEventResourceRequired
	}

	if strings.TrimSpace(ne.PayloadHash) == "" {
		return ErrSyncEventPayloadHashRequired
	}

	return nil
}

func validateSyncJobStatus(status SyncJobStatus) error {
	switch status {
	case SyncJobStatusQueued,
		SyncJobStatusRunning,
		SyncJobStatusSucceeded,
		SyncJobStatusFailed,
		SyncJobStatusRetryReady:
		return nil
	default:
		return ErrInvalidSyncJobStatus
	}
}

func validateSyncEventStatus(status SyncEventStatus) error {
	switch status {
	case SyncEventStatusQueued,
		SyncEventStatusProcessing,
		SyncEventStatusSucceeded,
		SyncEventStatusFailed,
		SyncEventStatusRetryReady:
		return nil
	default:
		return ErrInvalidSyncEventStatus
	}
}

func validateSyncDirection(direction SyncDirection) error {
	switch direction {
	case SyncDirectionInbound,
		SyncDirectionOutbound:
		return nil
	default:
		return ErrInvalidSyncDirection
	}
}

func validateSyncEventType(eventType SyncEventType) error {
	switch eventType {
	case SyncEventTypeBatchTermsPull,
		SyncEventTypeBatchProgramsPull,
		SyncEventTypeBatchPersonMatchesPull,
		SyncEventTypeBatchEnrollmentPull,
		SyncEventTypeApplicationSubmission,
		SyncEventTypeApplicationDecision,
		SyncEventTypeDocumentStatus,
		SyncEventTypeEnrollmentIntent:
		return nil
	default:
		return ErrInvalidSyncEventType
	}
}

func normalizeApplicationFormFields(fields []ApplicationFormField) []ApplicationFormField {
	normalized := make([]ApplicationFormField, len(fields))
	for i, field := range fields {
		normalized[i] = ApplicationFormField{
			FieldName:    strings.TrimSpace(field.FieldName),
			FieldType:    strings.TrimSpace(field.FieldType),
			Required:     field.Required,
			DisplayOrder: field.DisplayOrder,
			Validation:   trimStringPtr(field.Validation),
		}
	}

	return normalized
}

func normalizeChecklistTemplateItems(items []ApplicationChecklistTemplateItem) []ApplicationChecklistTemplateItem {
	normalized := make([]ApplicationChecklistTemplateItem, len(items))
	for i, item := range items {
		normalized[i] = ApplicationChecklistTemplateItem{
			ItemKey:      strings.TrimSpace(item.ItemKey),
			DocumentName: strings.TrimSpace(item.DocumentName),
			Description:  trimStringPtr(item.Description),
			Required:     item.Required,
			DisplayOrder: item.DisplayOrder,
		}
	}

	return normalized
}

func validateAdmissionsRole(role AdmissionsRole) error {
	switch role {
	case AdmissionsRoleAdmin,
		AdmissionsRoleRecruiter,
		AdmissionsRoleApplicationReviewer,
		AdmissionsRoleMarketingManager,
		AdmissionsRoleEventManager,
		AdmissionsRoleReportViewer,
		AdmissionsRoleApplicant:
		return nil
	default:
		return ErrInvalidAdmissionsRole
	}
}

// AdmissionsPermissionsForRoles returns the action permissions granted by a set
// of admissions context roles.
func AdmissionsPermissionsForRoles(roles []AdmissionsRole) []AdmissionsPermission {
	permissions := make(map[AdmissionsPermission]struct{})

	for _, role := range roles {
		for _, permission := range admissionsPermissionsForRole(role) {
			permissions[permission] = struct{}{}
		}
	}

	result := make([]AdmissionsPermission, 0, len(permissions))
	for permission := range permissions {
		result = append(result, permission)
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i] < result[j]
	})

	return result
}

func admissionsPermissionsForRole(role AdmissionsRole) []AdmissionsPermission {
	switch role {
	case AdmissionsRoleAdmin:
		return []AdmissionsPermission{
			AdmissionsPermissionRead,
			AdmissionsPermissionManageConstituents,
			AdmissionsPermissionManageApplications,
			AdmissionsPermissionReviewApplications,
			AdmissionsPermissionResolveDuplicates,
			AdmissionsPermissionManageReferences,
			AdmissionsPermissionManageStaff,
			AdmissionsPermissionManageLeadScoring,
		}
	case AdmissionsRoleRecruiter:
		return []AdmissionsPermission{
			AdmissionsPermissionRead,
			AdmissionsPermissionManageConstituents,
			AdmissionsPermissionManageApplications,
		}
	case AdmissionsRoleApplicationReviewer:
		return []AdmissionsPermission{
			AdmissionsPermissionRead,
			AdmissionsPermissionReviewApplications,
			AdmissionsPermissionResolveDuplicates,
		}
	case AdmissionsRoleMarketingManager,
		AdmissionsRoleEventManager:
		return []AdmissionsPermission{
			AdmissionsPermissionRead,
			AdmissionsPermissionManageLeadScoring,
		}
	case AdmissionsRoleReportViewer,
		AdmissionsRoleApplicant:
		return []AdmissionsPermission{AdmissionsPermissionRead}
	default:
		return nil
	}
}

// ParseAdmissionsRoles converts persisted role strings into admissions roles.
func ParseAdmissionsRoles(values []string) ([]AdmissionsRole, error) {
	roles := make([]AdmissionsRole, len(values))
	for i, value := range values {
		role := AdmissionsRole(value)
		if err := validateAdmissionsRole(role); err != nil {
			return nil, err
		}
		roles[i] = role
	}

	return roles, nil
}

// AdmissionsRolesToStrings converts admissions roles into strings for storage and claims.
func AdmissionsRolesToStrings(roles []AdmissionsRole) []string {
	values := make([]string, len(roles))
	for i, role := range roles {
		values[i] = role.String()
	}

	return values
}

// AdmissionsPermissionsToStrings converts admissions permissions into strings for clients and audits.
func AdmissionsPermissionsToStrings(permissions []AdmissionsPermission) []string {
	values := make([]string, len(permissions))
	for i, permission := range permissions {
		values[i] = permission.String()
	}

	return values
}

func validateNewLeadScoreRule(nr NewLeadScoreRule) error {
	if strings.TrimSpace(nr.Name) == "" {
		return ErrLeadScoreRuleNameRequired
	}

	if len(nr.Criteria) == 0 {
		return ErrLeadScoreCriteriaRequired
	}

	if nr.Points < 0 {
		return ErrInvalidLeadScorePoints
	}

	if nr.Priority < 0 {
		return ErrInvalidLeadScorePriority
	}

	for _, criterion := range nr.Criteria {
		if err := validateLeadScoreCriterion(criterion); err != nil {
			return err
		}
	}

	return nil
}

func validateLeadScoreCriterion(criterion LeadScoreCriterion) error {
	if len(criterion.Values) == 0 {
		return ErrInvalidLeadScoreCriterion
	}

	switch criterion.Field {
	case LeadScoreCriterionFieldLifecycleStage,
		LeadScoreCriterionFieldApplicationType,
		LeadScoreCriterionFieldApplicationStatus,
		LeadScoreCriterionFieldProgramID,
		LeadScoreCriterionFieldAcademicTermID:
	default:
		return ErrInvalidLeadScoreCriterion
	}

	switch criterion.Operator {
	case LeadScoreCriterionOperatorEquals:
		if len(criterion.Values) != 1 {
			return ErrInvalidLeadScoreCriterion
		}
	case LeadScoreCriterionOperatorIn:
	default:
		return ErrInvalidLeadScoreCriterion
	}

	for _, value := range criterion.Values {
		if strings.TrimSpace(value) == "" {
			return ErrInvalidLeadScoreCriterion
		}
	}

	return nil
}

// LeadScoreBandForTotal derives the score band from a total score.
func LeadScoreBandForTotal(total int) LeadScoreBand {
	switch {
	case total >= 76:
		return LeadScoreBandReadyToApply
	case total >= 51:
		return LeadScoreBandHot
	case total >= 26:
		return LeadScoreBandWarm
	default:
		return LeadScoreBandCold
	}
}

func validateLeadScoreBand(band LeadScoreBand) error {
	switch band {
	case LeadScoreBandCold,
		LeadScoreBandWarm,
		LeadScoreBandHot,
		LeadScoreBandReadyToApply:
		return nil
	default:
		return ErrInvalidLeadScoreBand
	}
}

func evaluateLeadScoreRule(rule LeadScoreRule, constituent Constituent, applications []Application) LeadScoreRuleResult {
	result := LeadScoreRuleResult{
		RuleID: rule.ID,
		Name:   rule.Name,
		Points: rule.Points,
	}

	for _, criterion := range rule.Criteria {
		matched, value := evaluateLeadScoreCriterion(criterion, constituent, applications)
		if !matched {
			result.Reason = fmt.Sprintf("%s did not match %s", criterion.Field, strings.Join(criterion.Values, ", "))
			return result
		}

		result.Reason = appendReason(result.Reason, fmt.Sprintf("%s matched %s", criterion.Field, value))
	}

	result.Matched = true
	return result
}

func evaluateLeadScoreCriterion(criterion LeadScoreCriterion, constituent Constituent, applications []Application) (bool, string) {
	values := valuesForLeadScoreCriterion(criterion.Field, constituent, applications)
	for _, candidate := range values {
		if criterionMatchesValue(criterion, candidate) {
			return true, candidate
		}
	}

	return false, ""
}

func valuesForLeadScoreCriterion(field LeadScoreCriterionField, constituent Constituent, applications []Application) []string {
	switch field {
	case LeadScoreCriterionFieldLifecycleStage:
		return []string{constituent.LifecycleStage.String()}
	case LeadScoreCriterionFieldApplicationType:
		values := make([]string, 0, len(applications))
		for _, app := range applications {
			values = append(values, app.ApplicationType.String())
		}
		return values
	case LeadScoreCriterionFieldApplicationStatus:
		values := make([]string, 0, len(applications))
		for _, app := range applications {
			values = append(values, app.Status.String())
		}
		return values
	case LeadScoreCriterionFieldProgramID:
		values := make([]string, 0, len(applications))
		for _, app := range applications {
			values = append(values, app.ProgramID.String())
		}
		return values
	case LeadScoreCriterionFieldAcademicTermID:
		values := make([]string, 0, len(applications))
		for _, app := range applications {
			values = append(values, app.AcademicTermID.String())
		}
		return values
	default:
		return nil
	}
}

func criterionMatchesValue(criterion LeadScoreCriterion, candidate string) bool {
	for _, value := range criterion.Values {
		if strings.EqualFold(candidate, strings.TrimSpace(value)) {
			return true
		}
	}

	return false
}

func appendReason(existing string, next string) string {
	if existing == "" {
		return next
	}

	return existing + "; " + next
}

func validateLifecycleStage(stage LifecycleStage) error {
	switch stage {
	case LifecycleStageProspect,
		LifecycleStageInquiry,
		LifecycleStageApplicant,
		LifecycleStageAdmitted,
		LifecycleStageEnrolled,
		LifecycleStageAlumni:
		return nil
	default:
		return ErrInvalidLifecycleStage
	}
}

func canChangeLifecycleStage(from LifecycleStage, to LifecycleStage) bool {
	if from == to {
		return true
	}

	if to == LifecycleStageAlumni {
		return true
	}

	transitions := map[LifecycleStage]LifecycleStage{
		LifecycleStageProspect:  LifecycleStageInquiry,
		LifecycleStageInquiry:   LifecycleStageApplicant,
		LifecycleStageApplicant: LifecycleStageAdmitted,
		LifecycleStageAdmitted:  LifecycleStageEnrolled,
	}

	return transitions[from] == to
}

func validateDuplicateStatus(status DuplicateStatus, duplicateOfID *uuid.UUID) error {
	switch status {
	case DuplicateStatusActive:
		if duplicateOfID != nil {
			return ErrInvalidDuplicateLink
		}
		return nil
	case DuplicateStatusMerged, DuplicateStatusDuplicateOf:
		if duplicateOfID == nil {
			return ErrInvalidDuplicateLink
		}
		return nil
	default:
		return ErrInvalidDuplicateStatus
	}
}

func validateNewDuplicateReview(nr NewDuplicateReview) error {
	if nr.SourceConstituentID == uuid.Nil || nr.CandidateConstituentID == uuid.Nil || nr.SourceConstituentID == nr.CandidateConstituentID {
		return ErrInvalidDuplicateReview
	}

	if err := validateDuplicateReviewMatchType(nr.MatchType); err != nil {
		return err
	}

	if nr.MatchScore < 0 || nr.MatchScore > 100 {
		return ErrInvalidMatchScore
	}

	if strings.TrimSpace(nr.MatchReason) == "" {
		return ErrMatchReasonRequired
	}

	return nil
}

func validateDuplicateReviewMatchType(matchType DuplicateReviewMatchType) error {
	switch matchType {
	case DuplicateReviewMatchTypeExact, DuplicateReviewMatchTypeFuzzy:
		return nil
	default:
		return ErrInvalidMatchType
	}
}

func validateDuplicateReviewStatus(status DuplicateReviewStatus) error {
	switch status {
	case DuplicateReviewStatusPending,
		DuplicateReviewStatusLinked,
		DuplicateReviewStatusMerged,
		DuplicateReviewStatusRejected,
		DuplicateReviewStatusDeferred:
		return nil
	default:
		return ErrInvalidReviewStatus
	}
}

func statusForResolution(resolution DuplicateReviewResolution) (DuplicateReviewStatus, error) {
	switch resolution {
	case DuplicateReviewResolutionLink:
		return DuplicateReviewStatusLinked, nil
	case DuplicateReviewResolutionMerge:
		return DuplicateReviewStatusMerged, nil
	case DuplicateReviewResolutionReject:
		return DuplicateReviewStatusRejected, nil
	case DuplicateReviewResolutionDefer:
		return DuplicateReviewStatusDeferred, nil
	default:
		return "", ErrInvalidResolution
	}
}

func validateNewApplication(na NewApplication) error {
	if na.ConstituentID == uuid.Nil {
		return ErrConstituentIDRequired
	}

	if na.ProgramID == uuid.Nil {
		return ErrProgramIDRequired
	}

	if na.AcademicTermID == uuid.Nil {
		return ErrAcademicTermIDRequired
	}

	return validateApplicationType(na.ApplicationType)
}

func validateApplicationType(applicationType ApplicationType) error {
	switch applicationType {
	case ApplicationTypeFreshman,
		ApplicationTypeTransfer,
		ApplicationTypeGraduate:
		return nil
	default:
		return ErrInvalidApplicationType
	}
}

func validateApplicationStatus(status ApplicationStatus) error {
	switch status {
	case ApplicationStatusDraft,
		ApplicationStatusSubmitted,
		ApplicationStatusAwaitingDocuments,
		ApplicationStatusReadyForReview,
		ApplicationStatusInReview,
		ApplicationStatusDecisionPending,
		ApplicationStatusAdmitted,
		ApplicationStatusDenied,
		ApplicationStatusWaitlisted,
		ApplicationStatusDeferred,
		ApplicationStatusWithdrawn,
		ApplicationStatusEnrolled:
		return nil
	default:
		return ErrInvalidApplicationStatus
	}
}

func validateDocumentStatus(status DocumentStatus) error {
	switch status {
	case DocumentStatusUploaded,
		DocumentStatusPendingReview,
		DocumentStatusAccepted,
		DocumentStatusRejected,
		DocumentStatusWaived,
		DocumentStatusExpired,
		DocumentStatusSyncedToSIS:
		return nil
	default:
		return ErrInvalidDocumentStatus
	}
}

func isReviewDocumentStatus(status DocumentStatus) bool {
	switch status {
	case DocumentStatusAccepted,
		DocumentStatusRejected,
		DocumentStatusWaived:
		return true
	default:
		return false
	}
}

func isApplicationActive(status ApplicationStatus) bool {
	switch status {
	case ApplicationStatusDenied,
		ApplicationStatusWithdrawn,
		ApplicationStatusEnrolled:
		return false
	default:
		return true
	}
}

func canTransitionApplicationStatus(from ApplicationStatus, to ApplicationStatus) bool {
	transitions := map[ApplicationStatus][]ApplicationStatus{
		ApplicationStatusDraft: {
			ApplicationStatusSubmitted,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusSubmitted: {
			ApplicationStatusAwaitingDocuments,
			ApplicationStatusReadyForReview,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusAwaitingDocuments: {
			ApplicationStatusReadyForReview,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusReadyForReview: {
			ApplicationStatusInReview,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusInReview: {
			ApplicationStatusDecisionPending,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusDecisionPending: {
			ApplicationStatusAdmitted,
			ApplicationStatusDenied,
			ApplicationStatusWaitlisted,
			ApplicationStatusDeferred,
		},
		ApplicationStatusAdmitted: {
			ApplicationStatusEnrolled,
		},
		ApplicationStatusWaitlisted: {
			ApplicationStatusAdmitted,
			ApplicationStatusDeferred,
		},
		ApplicationStatusDeferred: {
			ApplicationStatusSubmitted,
		},
	}

	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}

	return false
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
