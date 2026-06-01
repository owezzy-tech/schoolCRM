// Package admissionsotel provides an extension for admissionsbus that adds
// otel tracking.
package admissionsotel

import (
	"context"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/otel"
)

// Extension provides a wrapper for otel functionality around the admissionsbus.
type Extension struct {
	bus admissionsbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the admissionsbus with otel.
func NewExtension() admissionsbus.Extension {
	return func(bus admissionsbus.ExtBusiness) admissionsbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

// NewWithTx does not apply otel.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (admissionsbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Health applies otel to the admissions scaffold health check.
func (ext *Extension) Health(ctx context.Context) (admissionsbus.Health, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.health")
	defer span.End()

	return ext.bus.Health(ctx)
}

// CreateStaffProfile applies otel to admissions staff profile creation.
func (ext *Extension) CreateStaffProfile(ctx context.Context, np admissionsbus.NewStaffProfile) (admissionsbus.StaffProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createstaffprofile")
	defer span.End()

	return ext.bus.CreateStaffProfile(ctx, np)
}

// UpdateStaffProfile applies otel to admissions staff profile updates.
func (ext *Extension) UpdateStaffProfile(ctx context.Context, profile admissionsbus.StaffProfile, np admissionsbus.NewStaffProfile) (admissionsbus.StaffProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updatestaffprofile")
	defer span.End()

	return ext.bus.UpdateStaffProfile(ctx, profile, np)
}

// QueryStaffProfiles applies otel to admissions staff profile queries.
func (ext *Extension) QueryStaffProfiles(ctx context.Context, filter admissionsbus.StaffProfileQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.StaffProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querystaffprofiles")
	defer span.End()

	return ext.bus.QueryStaffProfiles(ctx, filter, orderBy, page)
}

// CountStaffProfiles applies otel to admissions staff profile counts.
func (ext *Extension) CountStaffProfiles(ctx context.Context, filter admissionsbus.StaffProfileQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countstaffprofiles")
	defer span.End()

	return ext.bus.CountStaffProfiles(ctx, filter)
}

// QueryStaffProfileByID applies otel to admissions staff profile ID lookups.
func (ext *Extension) QueryStaffProfileByID(ctx context.Context, profileID uuid.UUID) (admissionsbus.StaffProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querystaffprofilebyid")
	defer span.End()

	return ext.bus.QueryStaffProfileByID(ctx, profileID)
}

// QueryStaffProfileByUserID applies otel to admissions staff profile user lookups.
func (ext *Extension) QueryStaffProfileByUserID(ctx context.Context, userID uuid.UUID) (admissionsbus.StaffProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querystaffprofilebyuserid")
	defer span.End()

	return ext.bus.QueryStaffProfileByUserID(ctx, userID)
}

// CreateApplicantProfile applies otel to admissions applicant profile creation.
func (ext *Extension) CreateApplicantProfile(ctx context.Context, np admissionsbus.NewApplicantProfile) (admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createapplicantprofile")
	defer span.End()

	return ext.bus.CreateApplicantProfile(ctx, np)
}

// UpdateApplicantProfile applies otel to admissions applicant profile updates.
func (ext *Extension) UpdateApplicantProfile(ctx context.Context, profile admissionsbus.ApplicantProfile, np admissionsbus.NewApplicantProfile) (admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateapplicantprofile")
	defer span.End()

	return ext.bus.UpdateApplicantProfile(ctx, profile, np)
}

// QueryApplicantProfiles applies otel to admissions applicant profile queries.
func (ext *Extension) QueryApplicantProfiles(ctx context.Context, filter admissionsbus.ApplicantProfileQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicantprofiles")
	defer span.End()

	return ext.bus.QueryApplicantProfiles(ctx, filter, orderBy, page)
}

// CountApplicantProfiles applies otel to admissions applicant profile counts.
func (ext *Extension) CountApplicantProfiles(ctx context.Context, filter admissionsbus.ApplicantProfileQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countapplicantprofiles")
	defer span.End()

	return ext.bus.CountApplicantProfiles(ctx, filter)
}

// QueryApplicantProfileByID applies otel to admissions applicant profile ID lookups.
func (ext *Extension) QueryApplicantProfileByID(ctx context.Context, profileID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicantprofilebyid")
	defer span.End()

	return ext.bus.QueryApplicantProfileByID(ctx, profileID)
}

// QueryApplicantProfileByUserID applies otel to admissions applicant profile user lookups.
func (ext *Extension) QueryApplicantProfileByUserID(ctx context.Context, userID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicantprofilebyuserid")
	defer span.End()

	return ext.bus.QueryApplicantProfileByUserID(ctx, userID)
}

// QueryApplicantProfileByConstituentID applies otel to admissions applicant profile constituent lookups.
func (ext *Extension) QueryApplicantProfileByConstituentID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.ApplicantProfile, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicantprofilebyconstituentid")
	defer span.End()

	return ext.bus.QueryApplicantProfileByConstituentID(ctx, constituentID)
}

// CreateInquiry applies otel to admissions inquiry creation.
func (ext *Extension) CreateInquiry(ctx context.Context, ni admissionsbus.NewInquiry) (admissionsbus.Inquiry, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createinquiry")
	defer span.End()

	return ext.bus.CreateInquiry(ctx, ni)
}

// QueryInquiries applies otel to admissions inquiry queries.
func (ext *Extension) QueryInquiries(ctx context.Context, filter admissionsbus.InquiryQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Inquiry, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryinquiries")
	defer span.End()

	return ext.bus.QueryInquiries(ctx, filter, orderBy, page)
}

// CountInquiries applies otel to admissions inquiry counts.
func (ext *Extension) CountInquiries(ctx context.Context, filter admissionsbus.InquiryQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countinquiries")
	defer span.End()

	return ext.bus.CountInquiries(ctx, filter)
}

// QueryInquiryByID applies otel to admissions inquiry ID lookups.
func (ext *Extension) QueryInquiryByID(ctx context.Context, inquiryID uuid.UUID) (admissionsbus.Inquiry, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryinquirybyid")
	defer span.End()

	return ext.bus.QueryInquiryByID(ctx, inquiryID)
}

// CreateLeadScoreRule applies otel to lead score rule creation.
func (ext *Extension) CreateLeadScoreRule(ctx context.Context, nr admissionsbus.NewLeadScoreRule) (admissionsbus.LeadScoreRule, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createleadscorerule")
	defer span.End()

	return ext.bus.CreateLeadScoreRule(ctx, nr)
}

// UpdateLeadScoreRule applies otel to lead score rule updates.
func (ext *Extension) UpdateLeadScoreRule(ctx context.Context, rule admissionsbus.LeadScoreRule, nr admissionsbus.NewLeadScoreRule) (admissionsbus.LeadScoreRule, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateleadscorerule")
	defer span.End()

	return ext.bus.UpdateLeadScoreRule(ctx, rule, nr)
}

// QueryLeadScoreRules applies otel to lead score rule queries.
func (ext *Extension) QueryLeadScoreRules(ctx context.Context, filter admissionsbus.LeadScoreRuleQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.LeadScoreRule, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryleadscorerules")
	defer span.End()

	return ext.bus.QueryLeadScoreRules(ctx, filter, orderBy, page)
}

// CountLeadScoreRules applies otel to lead score rule counts.
func (ext *Extension) CountLeadScoreRules(ctx context.Context, filter admissionsbus.LeadScoreRuleQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countleadscorerules")
	defer span.End()

	return ext.bus.CountLeadScoreRules(ctx, filter)
}

// QueryLeadScoreRuleByID applies otel to lead score rule ID lookups.
func (ext *Extension) QueryLeadScoreRuleByID(ctx context.Context, ruleID uuid.UUID) (admissionsbus.LeadScoreRule, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryleadscorerulebyid")
	defer span.End()

	return ext.bus.QueryLeadScoreRuleByID(ctx, ruleID)
}

// RecalculateLeadScoreForConstituent applies otel to lead score recalculation.
func (ext *Extension) RecalculateLeadScoreForConstituent(ctx context.Context, constituentID uuid.UUID) (admissionsbus.LeadScore, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.recalculateleadscoreforconstituent")
	defer span.End()

	return ext.bus.RecalculateLeadScoreForConstituent(ctx, constituentID)
}

// QueryLeadScores applies otel to lead score queries.
func (ext *Extension) QueryLeadScores(ctx context.Context, filter admissionsbus.LeadScoreQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.LeadScore, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryleadscores")
	defer span.End()

	return ext.bus.QueryLeadScores(ctx, filter, orderBy, page)
}

// CountLeadScores applies otel to lead score counts.
func (ext *Extension) CountLeadScores(ctx context.Context, filter admissionsbus.LeadScoreQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countleadscores")
	defer span.End()

	return ext.bus.CountLeadScores(ctx, filter)
}

// QueryLeadScoreByID applies otel to lead score ID lookups.
func (ext *Extension) QueryLeadScoreByID(ctx context.Context, scoreID uuid.UUID) (admissionsbus.LeadScore, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryleadscorebyid")
	defer span.End()

	return ext.bus.QueryLeadScoreByID(ctx, scoreID)
}

// QueryLeadScoreByConstituentID applies otel to lead score constituent lookups.
func (ext *Extension) QueryLeadScoreByConstituentID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.LeadScore, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryleadscorebyconstituentid")
	defer span.End()

	return ext.bus.QueryLeadScoreByConstituentID(ctx, constituentID)
}

// CreateConstituent applies otel to Constituent creation.
func (ext *Extension) CreateConstituent(ctx context.Context, nc admissionsbus.NewConstituent) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createconstituent")
	defer span.End()

	return ext.bus.CreateConstituent(ctx, nc)
}

// UpdateConstituent applies otel to Constituent updates.
func (ext *Extension) UpdateConstituent(ctx context.Context, cst admissionsbus.Constituent, uc admissionsbus.UpdateConstituent) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateconstituent")
	defer span.End()

	return ext.bus.UpdateConstituent(ctx, cst, uc)
}

// QueryConstituents applies otel to Constituent queries.
func (ext *Extension) QueryConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituents")
	defer span.End()

	return ext.bus.QueryConstituents(ctx, filter, orderBy, page)
}

// CountConstituents applies otel to Constituent counts.
func (ext *Extension) CountConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countconstituents")
	defer span.End()

	return ext.bus.CountConstituents(ctx, filter)
}

// QueryConstituentByID applies otel to Constituent ID lookups.
func (ext *Extension) QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyid")
	defer span.End()

	return ext.bus.QueryConstituentByID(ctx, constituentID)
}

// QueryConstituentByPrimaryEmail applies otel to Constituent email lookups.
func (ext *Extension) QueryConstituentByPrimaryEmail(ctx context.Context, email string) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyprimaryemail")
	defer span.End()

	return ext.bus.QueryConstituentByPrimaryEmail(ctx, email)
}

// QueryConstituentByExternalSISID applies otel to Constituent SIS ID lookups.
func (ext *Extension) QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyexternalsisid")
	defer span.End()

	return ext.bus.QueryConstituentByExternalSISID(ctx, externalSISID)
}

// UpsertProgram applies otel to Program sync/import upserts.
func (ext *Extension) UpsertProgram(ctx context.Context, up admissionsbus.UpsertProgram) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.upsertprogram")
	defer span.End()

	return ext.bus.UpsertProgram(ctx, up)
}

// QueryPrograms applies otel to Program queries.
func (ext *Extension) QueryPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprograms")
	defer span.End()

	return ext.bus.QueryPrograms(ctx, filter, orderBy, page)
}

// CountPrograms applies otel to Program counts.
func (ext *Extension) CountPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countprograms")
	defer span.End()

	return ext.bus.CountPrograms(ctx, filter)
}

// QueryProgramByID applies otel to Program ID lookups.
func (ext *Extension) QueryProgramByID(ctx context.Context, programID uuid.UUID) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprogrambyid")
	defer span.End()

	return ext.bus.QueryProgramByID(ctx, programID)
}

// QueryProgramByExternalSISID applies otel to Program SIS ID lookups.
func (ext *Extension) QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprogrambyexternalsisid")
	defer span.End()

	return ext.bus.QueryProgramByExternalSISID(ctx, externalSISID)
}

// UpsertAcademicTerm applies otel to AcademicTerm sync/import upserts.
func (ext *Extension) UpsertAcademicTerm(ctx context.Context, up admissionsbus.UpsertAcademicTerm) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.upsertacademicterm")
	defer span.End()

	return ext.bus.UpsertAcademicTerm(ctx, up)
}

// QueryAcademicTerms applies otel to AcademicTerm queries.
func (ext *Extension) QueryAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademicterms")
	defer span.End()

	return ext.bus.QueryAcademicTerms(ctx, filter, orderBy, page)
}

// CountAcademicTerms applies otel to AcademicTerm counts.
func (ext *Extension) CountAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countacademicterms")
	defer span.End()

	return ext.bus.CountAcademicTerms(ctx, filter)
}

// QueryAcademicTermByID applies otel to AcademicTerm ID lookups.
func (ext *Extension) QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademictermbyid")
	defer span.End()

	return ext.bus.QueryAcademicTermByID(ctx, termID)
}

// QueryAcademicTermByExternalSISID applies otel to AcademicTerm SIS ID lookups.
func (ext *Extension) QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademictermbyexternalsisid")
	defer span.End()

	return ext.bus.QueryAcademicTermByExternalSISID(ctx, externalSISID)
}

// CreateDuplicateReview applies otel to duplicate review creation.
func (ext *Extension) CreateDuplicateReview(ctx context.Context, nr admissionsbus.NewDuplicateReview) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createduplicatereview")
	defer span.End()

	return ext.bus.CreateDuplicateReview(ctx, nr)
}

// ResolveDuplicateReview applies otel to duplicate review resolution.
func (ext *Extension) ResolveDuplicateReview(ctx context.Context, review admissionsbus.DuplicateReview, rr admissionsbus.ResolveDuplicateReview) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.resolveduplicatereview")
	defer span.End()

	return ext.bus.ResolveDuplicateReview(ctx, review, rr)
}

// QueryDuplicateReviews applies otel to duplicate review queries.
func (ext *Extension) QueryDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryduplicatereviews")
	defer span.End()

	return ext.bus.QueryDuplicateReviews(ctx, filter, orderBy, page)
}

// CountDuplicateReviews applies otel to duplicate review counts.
func (ext *Extension) CountDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countduplicatereviews")
	defer span.End()

	return ext.bus.CountDuplicateReviews(ctx, filter)
}

// QueryDuplicateReviewByID applies otel to duplicate review ID lookups.
func (ext *Extension) QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryduplicatereviewbyid")
	defer span.End()

	return ext.bus.QueryDuplicateReviewByID(ctx, reviewID)
}

// CreateApplication applies otel to Application creation.
func (ext *Extension) CreateApplication(ctx context.Context, na admissionsbus.NewApplication) (admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createapplication")
	defer span.End()

	return ext.bus.CreateApplication(ctx, na)
}

// QueryApplications applies otel to Application queries.
func (ext *Extension) QueryApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplications")
	defer span.End()

	return ext.bus.QueryApplications(ctx, filter, orderBy, page)
}

// CountApplications applies otel to Application counts.
func (ext *Extension) CountApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countapplications")
	defer span.End()

	return ext.bus.CountApplications(ctx, filter)
}

// QueryApplicationByID applies otel to Application ID lookups.
func (ext *Extension) QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicationbyid")
	defer span.End()

	return ext.bus.QueryApplicationByID(ctx, applicationID)
}

// CreateApplicationFormTemplate applies otel to Application form template creation.
func (ext *Extension) CreateApplicationFormTemplate(ctx context.Context, nt admissionsbus.NewApplicationFormTemplate) (admissionsbus.ApplicationFormTemplate, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createapplicationformtemplate")
	defer span.End()

	return ext.bus.CreateApplicationFormTemplate(ctx, nt)
}

// UpdateApplicationFormTemplate applies otel to Application form template updates.
func (ext *Extension) UpdateApplicationFormTemplate(ctx context.Context, template admissionsbus.ApplicationFormTemplate, nt admissionsbus.NewApplicationFormTemplate) (admissionsbus.ApplicationFormTemplate, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateapplicationformtemplate")
	defer span.End()

	return ext.bus.UpdateApplicationFormTemplate(ctx, template, nt)
}

// QueryApplicationFormTemplates applies otel to Application form template queries.
func (ext *Extension) QueryApplicationFormTemplates(ctx context.Context, filter admissionsbus.ApplicationFormTemplateQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicationFormTemplate, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicationformtemplates")
	defer span.End()

	return ext.bus.QueryApplicationFormTemplates(ctx, filter, orderBy, page)
}

// CountApplicationFormTemplates applies otel to Application form template counts.
func (ext *Extension) CountApplicationFormTemplates(ctx context.Context, filter admissionsbus.ApplicationFormTemplateQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countapplicationformtemplates")
	defer span.End()

	return ext.bus.CountApplicationFormTemplates(ctx, filter)
}

// QueryApplicationFormTemplateByID applies otel to Application form template ID lookups.
func (ext *Extension) QueryApplicationFormTemplateByID(ctx context.Context, templateID uuid.UUID) (admissionsbus.ApplicationFormTemplate, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicationformtemplatebyid")
	defer span.End()

	return ext.bus.QueryApplicationFormTemplateByID(ctx, templateID)
}

// CreateCustomFieldDefinition applies otel to custom field definition creation.
func (ext *Extension) CreateCustomFieldDefinition(ctx context.Context, nd admissionsbus.NewCustomFieldDefinition) (admissionsbus.CustomFieldDefinition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createcustomfielddefinition")
	defer span.End()

	return ext.bus.CreateCustomFieldDefinition(ctx, nd)
}

// UpdateCustomFieldDefinition applies otel to custom field definition updates.
func (ext *Extension) UpdateCustomFieldDefinition(ctx context.Context, definition admissionsbus.CustomFieldDefinition, nd admissionsbus.NewCustomFieldDefinition) (admissionsbus.CustomFieldDefinition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updatecustomfielddefinition")
	defer span.End()

	return ext.bus.UpdateCustomFieldDefinition(ctx, definition, nd)
}

// QueryCustomFieldDefinitions applies otel to custom field definition queries.
func (ext *Extension) QueryCustomFieldDefinitions(ctx context.Context, filter admissionsbus.CustomFieldDefinitionQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.CustomFieldDefinition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querycustomfielddefinitions")
	defer span.End()

	return ext.bus.QueryCustomFieldDefinitions(ctx, filter, orderBy, page)
}

// CountCustomFieldDefinitions applies otel to custom field definition counts.
func (ext *Extension) CountCustomFieldDefinitions(ctx context.Context, filter admissionsbus.CustomFieldDefinitionQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countcustomfielddefinitions")
	defer span.End()

	return ext.bus.CountCustomFieldDefinitions(ctx, filter)
}

// QueryCustomFieldDefinitionByID applies otel to custom field definition ID lookups.
func (ext *Extension) QueryCustomFieldDefinitionByID(ctx context.Context, definitionID uuid.UUID) (admissionsbus.CustomFieldDefinition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querycustomfielddefinitionbyid")
	defer span.End()

	return ext.bus.QueryCustomFieldDefinitionByID(ctx, definitionID)
}

// SetCustomFieldValue applies otel to custom field value upserts.
func (ext *Extension) SetCustomFieldValue(ctx context.Context, nv admissionsbus.NewCustomFieldValue) (admissionsbus.CustomFieldValue, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.setcustomfieldvalue")
	defer span.End()

	return ext.bus.SetCustomFieldValue(ctx, nv)
}

// QueryCustomFieldValues applies otel to custom field value queries.
func (ext *Extension) QueryCustomFieldValues(ctx context.Context, filter admissionsbus.CustomFieldValueQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.CustomFieldValue, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querycustomfieldvalues")
	defer span.End()

	return ext.bus.QueryCustomFieldValues(ctx, filter, orderBy, page)
}

// CountCustomFieldValues applies otel to custom field value counts.
func (ext *Extension) CountCustomFieldValues(ctx context.Context, filter admissionsbus.CustomFieldValueQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countcustomfieldvalues")
	defer span.End()

	return ext.bus.CountCustomFieldValues(ctx, filter)
}

// QueryCustomFieldValueByID applies otel to custom field value ID lookups.
func (ext *Extension) QueryCustomFieldValueByID(ctx context.Context, valueID uuid.UUID) (admissionsbus.CustomFieldValue, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querycustomfieldvaluebyid")
	defer span.End()

	return ext.bus.QueryCustomFieldValueByID(ctx, valueID)
}

// TransitionApplicationStatus applies otel to Application status transitions.
func (ext *Extension) TransitionApplicationStatus(ctx context.Context, app admissionsbus.Application, nt admissionsbus.NewApplicationTransition) (admissionsbus.Application, admissionsbus.ApplicationTransition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.transitionapplicationstatus")
	defer span.End()

	return ext.bus.TransitionApplicationStatus(ctx, app, nt)
}

// QueryApplicationTransitions applies otel to Application transition queries.
func (ext *Extension) QueryApplicationTransitions(ctx context.Context, filter admissionsbus.ApplicationTransitionQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ApplicationTransition, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicationtransitions")
	defer span.End()

	return ext.bus.QueryApplicationTransitions(ctx, filter, orderBy, page)
}

// CountApplicationTransitions applies otel to Application transition counts.
func (ext *Extension) CountApplicationTransitions(ctx context.Context, filter admissionsbus.ApplicationTransitionQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countapplicationtransitions")
	defer span.End()

	return ext.bus.CountApplicationTransitions(ctx, filter)
}

// CreateChecklistItem applies otel to checklist item creation.
func (ext *Extension) CreateChecklistItem(ctx context.Context, ni admissionsbus.NewChecklistItem) (admissionsbus.ChecklistItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createchecklistitem")
	defer span.End()

	return ext.bus.CreateChecklistItem(ctx, ni)
}

// UpdateChecklistItem applies otel to checklist item updates.
func (ext *Extension) UpdateChecklistItem(ctx context.Context, item admissionsbus.ChecklistItem, ni admissionsbus.NewChecklistItem) (admissionsbus.ChecklistItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updatechecklistitem")
	defer span.End()

	return ext.bus.UpdateChecklistItem(ctx, item, ni)
}

// QueryChecklistItems applies otel to checklist item queries.
func (ext *Extension) QueryChecklistItems(ctx context.Context, filter admissionsbus.ChecklistItemQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ChecklistItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querychecklistitems")
	defer span.End()

	return ext.bus.QueryChecklistItems(ctx, filter, orderBy, page)
}

// CountChecklistItems applies otel to checklist item counts.
func (ext *Extension) CountChecklistItems(ctx context.Context, filter admissionsbus.ChecklistItemQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countchecklistitems")
	defer span.End()

	return ext.bus.CountChecklistItems(ctx, filter)
}

// QueryChecklistItemByID applies otel to checklist item ID lookups.
func (ext *Extension) QueryChecklistItemByID(ctx context.Context, itemID uuid.UUID) (admissionsbus.ChecklistItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querychecklistitembyid")
	defer span.End()

	return ext.bus.QueryChecklistItemByID(ctx, itemID)
}

// CreateDocument applies otel to document metadata creation.
func (ext *Extension) CreateDocument(ctx context.Context, nd admissionsbus.NewDocument) (admissionsbus.Document, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createdocument")
	defer span.End()

	return ext.bus.CreateDocument(ctx, nd)
}

// VerifyDocument applies otel to document verification.
func (ext *Extension) VerifyDocument(ctx context.Context, document admissionsbus.Document, nv admissionsbus.NewDocumentVerification) (admissionsbus.Document, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.verifydocument")
	defer span.End()

	return ext.bus.VerifyDocument(ctx, document, nv)
}

// QueryDocuments applies otel to document metadata queries.
func (ext *Extension) QueryDocuments(ctx context.Context, filter admissionsbus.DocumentQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Document, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querydocuments")
	defer span.End()

	return ext.bus.QueryDocuments(ctx, filter, orderBy, page)
}

// CountDocuments applies otel to document metadata counts.
func (ext *Extension) CountDocuments(ctx context.Context, filter admissionsbus.DocumentQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countdocuments")
	defer span.End()

	return ext.bus.CountDocuments(ctx, filter)
}

// QueryDocumentByID applies otel to document metadata ID lookups.
func (ext *Extension) QueryDocumentByID(ctx context.Context, documentID uuid.UUID) (admissionsbus.Document, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querydocumentbyid")
	defer span.End()

	return ext.bus.QueryDocumentByID(ctx, documentID)
}

// CreateImportBatch applies otel to admissions import batch creation.
func (ext *Extension) CreateImportBatch(ctx context.Context, nb admissionsbus.NewImportBatch) (admissionsbus.ImportBatch, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createimportbatch")
	defer span.End()

	return ext.bus.CreateImportBatch(ctx, nb)
}

// UpdateImportBatch applies otel to admissions import batch updates.
func (ext *Extension) UpdateImportBatch(ctx context.Context, batch admissionsbus.ImportBatch, nb admissionsbus.NewImportBatch) (admissionsbus.ImportBatch, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateimportbatch")
	defer span.End()

	return ext.bus.UpdateImportBatch(ctx, batch, nb)
}

// QueryImportBatches applies otel to admissions import batch queries.
func (ext *Extension) QueryImportBatches(ctx context.Context, filter admissionsbus.ImportBatchQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ImportBatch, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryimportbatches")
	defer span.End()

	return ext.bus.QueryImportBatches(ctx, filter, orderBy, page)
}

// CountImportBatches applies otel to admissions import batch counts.
func (ext *Extension) CountImportBatches(ctx context.Context, filter admissionsbus.ImportBatchQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countimportbatches")
	defer span.End()

	return ext.bus.CountImportBatches(ctx, filter)
}

// QueryImportBatchByID applies otel to admissions import batch ID lookups.
func (ext *Extension) QueryImportBatchByID(ctx context.Context, batchID uuid.UUID) (admissionsbus.ImportBatch, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryimportbatchbyid")
	defer span.End()

	return ext.bus.QueryImportBatchByID(ctx, batchID)
}

// CreateImportInvalidRows applies otel to admissions import invalid row creation.
func (ext *Extension) CreateImportInvalidRows(ctx context.Context, rows []admissionsbus.NewImportInvalidRow) ([]admissionsbus.ImportInvalidRow, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createimportinvalidrows")
	defer span.End()

	return ext.bus.CreateImportInvalidRows(ctx, rows)
}

// QueryImportInvalidRows applies otel to admissions import invalid row queries.
func (ext *Extension) QueryImportInvalidRows(ctx context.Context, filter admissionsbus.ImportInvalidRowQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.ImportInvalidRow, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryimportinvalidrows")
	defer span.End()

	return ext.bus.QueryImportInvalidRows(ctx, filter, orderBy, page)
}

// CountImportInvalidRows applies otel to admissions import invalid row counts.
func (ext *Extension) CountImportInvalidRows(ctx context.Context, filter admissionsbus.ImportInvalidRowQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countimportinvalidrows")
	defer span.End()

	return ext.bus.CountImportInvalidRows(ctx, filter)
}

// QueryImportInvalidRowByID applies otel to admissions import invalid row ID lookups.
func (ext *Extension) QueryImportInvalidRowByID(ctx context.Context, rowID uuid.UUID) (admissionsbus.ImportInvalidRow, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryimportinvalidrowbyid")
	defer span.End()

	return ext.bus.QueryImportInvalidRowByID(ctx, rowID)
}

// CreateSyncJob applies otel to SIS sync job creation.
func (ext *Extension) CreateSyncJob(ctx context.Context, nj admissionsbus.NewSyncJob) (admissionsbus.SyncJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createsyncjob")
	defer span.End()

	return ext.bus.CreateSyncJob(ctx, nj)
}

// UpdateSyncJob applies otel to SIS sync job updates.
func (ext *Extension) UpdateSyncJob(ctx context.Context, job admissionsbus.SyncJob, uj admissionsbus.UpdateSyncJob) (admissionsbus.SyncJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updatesyncjob")
	defer span.End()

	return ext.bus.UpdateSyncJob(ctx, job, uj)
}

// QuerySyncJobs applies otel to SIS sync job queries.
func (ext *Extension) QuerySyncJobs(ctx context.Context, filter admissionsbus.SyncJobQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.SyncJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querysyncjobs")
	defer span.End()

	return ext.bus.QuerySyncJobs(ctx, filter, orderBy, page)
}

// CountSyncJobs applies otel to SIS sync job counts.
func (ext *Extension) CountSyncJobs(ctx context.Context, filter admissionsbus.SyncJobQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countsyncjobs")
	defer span.End()

	return ext.bus.CountSyncJobs(ctx, filter)
}

// QuerySyncJobByID applies otel to SIS sync job ID lookups.
func (ext *Extension) QuerySyncJobByID(ctx context.Context, jobID uuid.UUID) (admissionsbus.SyncJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querysyncjobbyid")
	defer span.End()

	return ext.bus.QuerySyncJobByID(ctx, jobID)
}

// EnqueueSyncEvent applies otel to SIS sync event enqueueing.
func (ext *Extension) EnqueueSyncEvent(ctx context.Context, ne admissionsbus.NewSyncEvent) (admissionsbus.SyncEvent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.enqueuesyncevent")
	defer span.End()

	return ext.bus.EnqueueSyncEvent(ctx, ne)
}

// UpdateSyncEvent applies otel to SIS sync event updates.
func (ext *Extension) UpdateSyncEvent(ctx context.Context, event admissionsbus.SyncEvent, ue admissionsbus.UpdateSyncEvent) (admissionsbus.SyncEvent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updatesyncevent")
	defer span.End()

	return ext.bus.UpdateSyncEvent(ctx, event, ue)
}

// QuerySyncEvents applies otel to SIS sync event queries.
func (ext *Extension) QuerySyncEvents(ctx context.Context, filter admissionsbus.SyncEventQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.SyncEvent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querysyncevents")
	defer span.End()

	return ext.bus.QuerySyncEvents(ctx, filter, orderBy, page)
}

// CountSyncEvents applies otel to SIS sync event counts.
func (ext *Extension) CountSyncEvents(ctx context.Context, filter admissionsbus.SyncEventQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countsyncevents")
	defer span.End()

	return ext.bus.CountSyncEvents(ctx, filter)
}

// QuerySyncEventByID applies otel to SIS sync event ID lookups.
func (ext *Extension) QuerySyncEventByID(ctx context.Context, eventID uuid.UUID) (admissionsbus.SyncEvent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.querysynceventbyid")
	defer span.End()

	return ext.bus.QuerySyncEventByID(ctx, eventID)
}
