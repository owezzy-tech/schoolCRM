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
	admissionsbus.OrderByConstituentDateCreated:    "date_created",
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
	admissionsbus.OrderByApplicationID:           "application_id",
	admissionsbus.OrderByApplicationStatus:       "status",
	admissionsbus.OrderByApplicationType:         "application_type",
	admissionsbus.OrderByApplicationSubmittedAt:  "submitted_at",
	admissionsbus.OrderByApplicationDateCreated:  "date_created",
	admissionsbus.OrderByApplicationDateUpdated:  "date_updated",
}

var eventOrderByFields = map[string]string{
	admissionsbus.OrderByEventID:          "event_id",
	admissionsbus.OrderByEventType:        "event_type",
	admissionsbus.OrderByEventStatus:      "status",
	admissionsbus.OrderByEventStartTime:   "start_time",
	admissionsbus.OrderByEventDateCreated: "date_created",
}

var eventRegistrationOrderByFields = map[string]string{
	admissionsbus.OrderByEventRegistrationID:           "event_registration_id",
	admissionsbus.OrderByEventRegistrationEvent:        "event_id",
	admissionsbus.OrderByEventRegistrationStatus:       "status",
	admissionsbus.OrderByEventRegistrationRegisteredAt: "registered_at",
	admissionsbus.OrderByEventRegistrationCheckedInAt:  "checked_in_at",
}

var applicationFormTemplateOrderByFields = map[string]string{
	admissionsbus.OrderByApplicationFormTemplateID:          "form_template_id",
	admissionsbus.OrderByApplicationFormTemplateName:        "name",
	admissionsbus.OrderByApplicationFormTemplateType:        "application_type",
	admissionsbus.OrderByApplicationFormTemplateVersion:     "version",
	admissionsbus.OrderByApplicationFormTemplatePriority:    "priority",
	admissionsbus.OrderByApplicationFormTemplateDateCreated: "date_created",
}

var customFieldDefinitionOrderByFields = map[string]string{
	admissionsbus.OrderByCustomFieldDefinitionID:           "custom_field_definition_id",
	admissionsbus.OrderByCustomFieldDefinitionOwner:        "owner",
	admissionsbus.OrderByCustomFieldDefinitionFieldKey:     "field_key",
	admissionsbus.OrderByCustomFieldDefinitionDisplayOrder: "display_order",
	admissionsbus.OrderByCustomFieldDefinitionDateCreated:  "date_created",
}

var customFieldValueOrderByFields = map[string]string{
	admissionsbus.OrderByCustomFieldValueID:          "custom_field_value_id",
	admissionsbus.OrderByCustomFieldValueDefinition:  "custom_field_definition_id",
	admissionsbus.OrderByCustomFieldValueOwner:       "owner",
	admissionsbus.OrderByCustomFieldValueOwnerID:     "owner_id",
	admissionsbus.OrderByCustomFieldValueDateCreated: "date_created",
}

var applicationTransitionOrderByFields = map[string]string{
	admissionsbus.OrderByApplicationTransitionID:          "application_transition_id",
	admissionsbus.OrderByApplicationTransitionApplication: "application_id",
	admissionsbus.OrderByApplicationTransitionActor:       "actor_id",
	admissionsbus.OrderByApplicationTransitionDateCreated: "date_created",
}

var checklistItemOrderByFields = map[string]string{
	admissionsbus.OrderByChecklistItemID:           "checklist_item_id",
	admissionsbus.OrderByChecklistItemApplication:  "application_id",
	admissionsbus.OrderByChecklistItemStatus:       "status",
	admissionsbus.OrderByChecklistItemDisplayOrder: "display_order",
	admissionsbus.OrderByChecklistItemDateCreated:  "date_created",
}

var documentOrderByFields = map[string]string{
	admissionsbus.OrderByDocumentID:            "document_id",
	admissionsbus.OrderByDocumentApplication:   "application_id",
	admissionsbus.OrderByDocumentChecklistItem: "checklist_item_id",
	admissionsbus.OrderByDocumentStatus:        "status",
	admissionsbus.OrderByDocumentUploadedAt:    "uploaded_at",
	admissionsbus.OrderByDocumentReviewedAt:    "reviewed_at",
}

var importBatchOrderByFields = map[string]string{
	admissionsbus.OrderByImportBatchID:          "import_batch_id",
	admissionsbus.OrderByImportBatchTarget:      "target",
	admissionsbus.OrderByImportBatchStatus:      "status",
	admissionsbus.OrderByImportBatchUploadedBy:  "uploaded_by_id",
	admissionsbus.OrderByImportBatchDateCreated: "date_created",
}

var importInvalidRowOrderByFields = map[string]string{
	admissionsbus.OrderByImportInvalidRowID:          "import_invalid_row_id",
	admissionsbus.OrderByImportInvalidRowBatch:       "import_batch_id",
	admissionsbus.OrderByImportInvalidRowNumber:      "row_number",
	admissionsbus.OrderByImportInvalidRowDateCreated: "date_created",
}

var syncJobOrderByFields = map[string]string{
	admissionsbus.OrderBySyncJobID:          "sync_job_id",
	admissionsbus.OrderBySyncJobStatus:      "status",
	admissionsbus.OrderBySyncJobDirection:   "direction",
	admissionsbus.OrderBySyncJobStartedAt:   "started_at",
	admissionsbus.OrderBySyncJobCompletedAt: "completed_at",
	admissionsbus.OrderBySyncJobDateCreated: "date_created",
}

var syncEventOrderByFields = map[string]string{
	admissionsbus.OrderBySyncEventID:          "sync_event_id",
	admissionsbus.OrderBySyncEventJob:         "sync_job_id",
	admissionsbus.OrderBySyncEventType:        "event_type",
	admissionsbus.OrderBySyncEventStatus:      "status",
	admissionsbus.OrderBySyncEventDirection:   "direction",
	admissionsbus.OrderBySyncEventResource:    "resource_type",
	admissionsbus.OrderBySyncEventDateCreated: "date_created",
}

var campaignOrderByFields = map[string]string{
	admissionsbus.OrderByCampaignID:          "campaign_id",
	admissionsbus.OrderByCampaignName:        "name",
	admissionsbus.OrderByCampaignStatus:      "status",
	admissionsbus.OrderByCampaignChannel:     "channel",
	admissionsbus.OrderByCampaignStartsAt:    "starts_at",
	admissionsbus.OrderByCampaignDateCreated: "date_created",
}

var campaignAuditEventOrderByFields = map[string]string{
	admissionsbus.OrderByCampaignAuditEventID:         "campaign_audit_event_id",
	admissionsbus.OrderByCampaignAuditEventCampaign:   "campaign_id",
	admissionsbus.OrderByCampaignAuditEventOccurredAt: "occurred_at",
}

var communicationOrderByFields = map[string]string{
	admissionsbus.OrderByCommunicationID:          "communication_id",
	admissionsbus.OrderByCommunicationChannel:     "channel",
	admissionsbus.OrderByCommunicationDirection:   "direction",
	admissionsbus.OrderByCommunicationStatus:      "status",
	admissionsbus.OrderByCommunicationOccurredAt:  "occurred_at",
	admissionsbus.OrderByCommunicationDateCreated: "date_created",
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

func eventOrderByClause(orderBy order.By) (string, error) {
	by, exists := eventOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func eventRegistrationOrderByClause(orderBy order.By) (string, error) {
	by, exists := eventRegistrationOrderByFields[orderBy.Field]
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

func customFieldDefinitionOrderByClause(orderBy order.By) (string, error) {
	by, exists := customFieldDefinitionOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func customFieldValueOrderByClause(orderBy order.By) (string, error) {
	by, exists := customFieldValueOrderByFields[orderBy.Field]
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

func checklistItemOrderByClause(orderBy order.By) (string, error) {
	by, exists := checklistItemOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func documentOrderByClause(orderBy order.By) (string, error) {
	by, exists := documentOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func importBatchOrderByClause(orderBy order.By) (string, error) {
	by, exists := importBatchOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func importInvalidRowOrderByClause(orderBy order.By) (string, error) {
	by, exists := importInvalidRowOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func syncJobOrderByClause(orderBy order.By) (string, error) {
	by, exists := syncJobOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func syncEventOrderByClause(orderBy order.By) (string, error) {
	by, exists := syncEventOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func campaignOrderByClause(orderBy order.By) (string, error) {
	by, exists := campaignOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func campaignAuditEventOrderByClause(orderBy order.By) (string, error) {
	by, exists := campaignAuditEventOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func communicationOrderByClause(orderBy order.By) (string, error) {
	by, exists := communicationOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
