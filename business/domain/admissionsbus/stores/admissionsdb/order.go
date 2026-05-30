package admissionsdb

import (
	"fmt"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
)

var programOrderByFields = map[string]string{
	admissionsbus.OrderByProgramID:     "program_id",
	admissionsbus.OrderByProgramName:   "name",
	admissionsbus.OrderByProgramCode:   "code",
	admissionsbus.OrderByProgramActive: "is_active",
}

var staffProfileOrderByFields = map[string]string{
	admissionsbus.OrderByStaffProfileID:          "staff_profile_id",
	admissionsbus.OrderByStaffProfileUser:        "user_id",
	admissionsbus.OrderByStaffProfileDateCreated: "date_created",
}

var applicantProfileOrderByFields = map[string]string{
	admissionsbus.OrderByApplicantProfileID:          "applicant_profile_id",
	admissionsbus.OrderByApplicantProfileUser:        "user_id",
	admissionsbus.OrderByApplicantProfileConstituent: "constituent_id",
	admissionsbus.OrderByApplicantProfileDateCreated: "date_created",
}

var leadScoreRuleOrderByFields = map[string]string{
	admissionsbus.OrderByLeadScoreRuleID:          "lead_score_rule_id",
	admissionsbus.OrderByLeadScoreRuleName:        "name",
	admissionsbus.OrderByLeadScoreRulePriority:    "priority",
	admissionsbus.OrderByLeadScoreRuleDateCreated: "date_created",
}

var leadScoreOrderByFields = map[string]string{
	admissionsbus.OrderByLeadScoreID:             "lead_score_id",
	admissionsbus.OrderByLeadScoreTotalScore:     "total_score",
	admissionsbus.OrderByLeadScoreBand:           "score_band",
	admissionsbus.OrderByLeadScoreRecalculatedAt: "recalculated_at",
}

var constituentOrderByFields = map[string]string{
	admissionsbus.OrderByConstituentID:             "constituent_id",
	admissionsbus.OrderByConstituentLastName:       "last_name",
	admissionsbus.OrderByConstituentPrimaryEmail:   "primary_email",
	admissionsbus.OrderByConstituentLifecycleStage: "lifecycle_stage",
}

var inquiryOrderByFields = map[string]string{
	admissionsbus.OrderByInquiryID:          "inquiry_id",
	admissionsbus.OrderByInquiryEmail:       "primary_email",
	admissionsbus.OrderByInquirySource:      "source",
	admissionsbus.OrderByInquiryStatus:      "status",
	admissionsbus.OrderByInquiryDateCreated: "date_created",
}

var academicTermOrderByFields = map[string]string{
	admissionsbus.OrderByAcademicTermID:        "academic_term_id",
	admissionsbus.OrderByAcademicTermName:      "name",
	admissionsbus.OrderByAcademicTermCode:      "code",
	admissionsbus.OrderByAcademicTermStartDate: "start_date",
	admissionsbus.OrderByAcademicTermActive:    "is_active",
}

var duplicateReviewOrderByFields = map[string]string{
	admissionsbus.OrderByDuplicateReviewID:          "duplicate_review_id",
	admissionsbus.OrderByDuplicateReviewStatus:      "status",
	admissionsbus.OrderByDuplicateReviewMatchType:   "match_type",
	admissionsbus.OrderByDuplicateReviewMatchScore:  "match_score",
	admissionsbus.OrderByDuplicateReviewDateCreated: "date_created",
}

var applicationOrderByFields = map[string]string{
	admissionsbus.OrderByApplicationID:          "application_id",
	admissionsbus.OrderByApplicationStatus:      "status",
	admissionsbus.OrderByApplicationType:        "application_type",
	admissionsbus.OrderByApplicationSubmittedAt: "submitted_at",
	admissionsbus.OrderByApplicationDateCreated: "date_created",
}

var applicationFormTemplateOrderByFields = map[string]string{
	admissionsbus.OrderByApplicationFormTemplateID:          "form_template_id",
	admissionsbus.OrderByApplicationFormTemplateName:        "name",
	admissionsbus.OrderByApplicationFormTemplateType:        "application_type",
	admissionsbus.OrderByApplicationFormTemplateVersion:     "version",
	admissionsbus.OrderByApplicationFormTemplatePriority:    "priority",
	admissionsbus.OrderByApplicationFormTemplateDateCreated: "date_created",
}

var applicationTransitionOrderByFields = map[string]string{
	admissionsbus.OrderByApplicationTransitionID:          "application_transition_id",
	admissionsbus.OrderByApplicationTransitionApplication: "application_id",
	admissionsbus.OrderByApplicationTransitionActor:       "actor_id",
	admissionsbus.OrderByApplicationTransitionDateCreated: "date_created",
}

func programOrderByClause(orderBy order.By) (string, error) {
	by, exists := programOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func staffProfileOrderByClause(orderBy order.By) (string, error) {
	by, exists := staffProfileOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func applicantProfileOrderByClause(orderBy order.By) (string, error) {
	by, exists := applicantProfileOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func leadScoreRuleOrderByClause(orderBy order.By) (string, error) {
	by, exists := leadScoreRuleOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func leadScoreOrderByClause(orderBy order.By) (string, error) {
	by, exists := leadScoreOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func constituentOrderByClause(orderBy order.By) (string, error) {
	by, exists := constituentOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func inquiryOrderByClause(orderBy order.By) (string, error) {
	by, exists := inquiryOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func academicTermOrderByClause(orderBy order.By) (string, error) {
	by, exists := academicTermOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func duplicateReviewOrderByClause(orderBy order.By) (string, error) {
	by, exists := duplicateReviewOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func applicationOrderByClause(orderBy order.By) (string, error) {
	by, exists := applicationOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func applicationFormTemplateOrderByClause(orderBy order.By) (string, error) {
	by, exists := applicationFormTemplateOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func applicationTransitionOrderByClause(orderBy order.By) (string, error) {
	by, exists := applicationTransitionOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
