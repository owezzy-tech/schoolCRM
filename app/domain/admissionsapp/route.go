package admissionsapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log           *logger.Logger
	DB            *sqlx.DB
	AdmissionsBus admissionsbus.ExtBusiness
	AuditBus      auditbus.ExtBusiness
	AuthClient    authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleRead := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsRead)
	ruleManageConstituents := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsManageConstituents)
	ruleManageApplications := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsManageApplications)
	ruleManageEvents := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsManageEvents)
	ruleReviewApplications := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsReviewApplications)
	ruleResolveDuplicates := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsResolveDuplicates)
	ruleManageStaff := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsManageStaff)
	ruleManageLeadScoring := authorizeAdmissions(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsManageLeadScoring)
	ruleApplicantRead := authorizeApplicant(cfg.AuthClient, cfg.AdmissionsBus, auth.RuleAdmissionsRead)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	api := newApp(cfg.AdmissionsBus, cfg.AuditBus)

	app.HandlerFunc(http.MethodPost, version, "/admissions/inquiries", api.createInquiry, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/health", api.health, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/inquiries", api.queryInquiries, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/inquiries/{inquiry_id}", api.queryInquiryByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/staff-profiles", api.queryStaffProfiles, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodGet, version, "/admissions/staff-profiles/{staff_profile_id}", api.queryStaffProfileByID, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodPost, version, "/admissions/staff-profiles", api.createStaffProfile, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodPut, version, "/admissions/staff-profiles/{staff_profile_id}", api.updateStaffProfile, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant-profiles", api.queryApplicantProfiles, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant-profiles/{applicant_profile_id}", api.queryApplicantProfileByID, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applicant-profiles", api.createApplicantProfile, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodPut, version, "/admissions/applicant-profiles/{applicant_profile_id}", api.updateApplicantProfile, authen, ruleManageStaff)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/profile", api.queryCurrentApplicantProfile, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/programs", api.queryApplicantPrograms, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/programs/{program_id}", api.queryApplicantProgramByID, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/academic-terms", api.queryApplicantAcademicTerms, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/academic-terms/{academic_term_id}", api.queryApplicantAcademicTermByID, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/application-form-templates", api.queryApplicantApplicationFormTemplates, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/application-form-templates/{form_template_id}", api.queryApplicantApplicationFormTemplateByID, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/events", api.queryApplicantEvents, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/events/{event_id}", api.queryApplicantEventByID, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applicant/events/{event_id}/registrations", api.registerForApplicantEvent, authen, ruleApplicantRead, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/applications", api.queryApplicantApplications, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/applications/{application_id}", api.queryApplicantApplicationByID, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applicant/applications", api.createApplicantApplication, authen, ruleApplicantRead, transaction)
	app.HandlerFunc(http.MethodPut, version, "/admissions/applicant/applications/{application_id}", api.updateApplicantApplication, authen, ruleApplicantRead, transaction)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applicant/applications/{application_id}/transitions", api.transitionApplicantApplicationStatus, authen, ruleApplicantRead, transaction)
	app.HandlerFunc(http.MethodPut, version, "/admissions/applicant/custom-field-values", api.setApplicantCustomFieldValue, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/applications/{application_id}/checklist-items", api.queryApplicantChecklistItems, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applicant/applications/{application_id}/documents", api.queryApplicantDocuments, authen, ruleApplicantRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applicant/applications/{application_id}/documents", api.createApplicantDocument, authen, ruleApplicantRead, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/lead-score-rules", api.queryLeadScoreRules, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/lead-score-rules/{lead_score_rule_id}", api.queryLeadScoreRuleByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/lead-score-rules", api.createLeadScoreRule, authen, ruleManageLeadScoring)
	app.HandlerFunc(http.MethodPut, version, "/admissions/lead-score-rules/{lead_score_rule_id}", api.updateLeadScoreRule, authen, ruleManageLeadScoring)
	app.HandlerFunc(http.MethodGet, version, "/admissions/lead-scores", api.queryLeadScores, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/lead-scores/{lead_score_id}", api.queryLeadScoreByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents", api.queryConstituents, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents/{constituent_id}", api.queryConstituentByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents/{constituent_id}/lead-score", api.queryLeadScoreByConstituentID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/constituents/{constituent_id}/lead-score/recalculate", api.recalculateLeadScoreForConstituent, authen, ruleManageLeadScoring, transaction)
	app.HandlerFunc(http.MethodPost, version, "/admissions/constituents", api.createConstituent, authen, ruleManageConstituents)
	app.HandlerFunc(http.MethodPut, version, "/admissions/constituents/{constituent_id}", api.updateConstituent, authen, ruleManageConstituents)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs", api.queryPrograms, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs/{program_id}", api.queryProgramByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms", api.queryAcademicTerms, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms/{academic_term_id}", api.queryAcademicTermByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications", api.queryApplications, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}", api.queryApplicationByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications", api.createApplication, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodGet, version, "/admissions/events", api.queryEvents, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/events/{event_id}", api.queryEventByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/events/{event_id}/registrations", api.queryEventRegistrations, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/events/{event_id}/registrations", api.registerForEvent, authen, ruleManageEvents, transaction)
	app.HandlerFunc(http.MethodPost, version, "/admissions/event-registrations/{event_registration_id}/check-in", api.checkInEventRegistration, authen, ruleManageEvents, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/application-form-templates", api.queryApplicationFormTemplates, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/application-form-templates/{form_template_id}", api.queryApplicationFormTemplateByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/application-form-templates", api.createApplicationFormTemplate, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodPut, version, "/admissions/application-form-templates/{form_template_id}", api.updateApplicationFormTemplate, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodGet, version, "/admissions/custom-field-definitions", api.queryCustomFieldDefinitions, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/custom-field-definitions/{custom_field_definition_id}", api.queryCustomFieldDefinitionByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/custom-field-definitions", api.createCustomFieldDefinition, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodPut, version, "/admissions/custom-field-definitions/{custom_field_definition_id}", api.updateCustomFieldDefinition, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodGet, version, "/admissions/custom-field-values", api.queryCustomFieldValues, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/custom-field-values/{custom_field_value_id}", api.queryCustomFieldValueByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPut, version, "/admissions/custom-field-values", api.setCustomFieldValue, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}/transitions", api.queryApplicationTransitions, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications/{application_id}/transitions", api.transitionApplicationStatus, authen, ruleReviewApplications, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}/checklist-items", api.queryChecklistItems, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}/checklist-items/{checklist_item_id}", api.queryChecklistItemByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications/{application_id}/checklist-items", api.createChecklistItem, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodPut, version, "/admissions/applications/{application_id}/checklist-items/{checklist_item_id}", api.updateChecklistItem, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}/documents", api.queryDocuments, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications/{application_id}/documents", api.createDocument, authen, ruleManageApplications, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/documents/{document_id}", api.queryDocumentByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/documents/{document_id}/download", api.downloadDocument, authen, ruleRead)
	app.HandlerFunc(http.MethodPut, version, "/admissions/documents/{document_id}/verification", api.verifyDocument, authen, ruleReviewApplications, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/import-batches", api.queryImportBatches, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/import-batches/{import_batch_id}", api.queryImportBatchByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/import-batches", api.createImportBatch, authen, ruleManageApplications, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/import-batches/{import_batch_id}/invalid-rows", api.queryImportInvalidRows, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/import-batches/{import_batch_id}/invalid-rows", api.createImportInvalidRows, authen, ruleManageApplications)
	app.HandlerFunc(http.MethodPost, version, "/admissions/import-batches/{import_batch_id}/invalid-rows/download", api.downloadImportInvalidRows, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/import-invalid-rows/{import_invalid_row_id}", api.queryImportInvalidRowByID, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/duplicate-reviews", api.queryDuplicateReviews, authen, ruleRead)
	app.HandlerFunc(http.MethodGet, version, "/admissions/duplicate-reviews/{duplicate_review_id}", api.queryDuplicateReviewByID, authen, ruleRead)
	app.HandlerFunc(http.MethodPost, version, "/admissions/duplicate-reviews", api.createDuplicateReview, authen, ruleResolveDuplicates)
	app.HandlerFunc(http.MethodPut, version, "/admissions/duplicate-reviews/{duplicate_review_id}/resolution", api.resolveDuplicateReview, authen, ruleResolveDuplicates)
}
