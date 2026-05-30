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
