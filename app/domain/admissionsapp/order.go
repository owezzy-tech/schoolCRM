package admissionsapp

import "github.com/owezzy/schoolCRM/business/domain/admissionsbus"

var programOrderByFields = map[string]string{
	"program_id": admissionsbus.OrderByProgramID,
	"name":       admissionsbus.OrderByProgramName,
	"code":       admissionsbus.OrderByProgramCode,
	"active":     admissionsbus.OrderByProgramActive,
}

var staffProfileOrderByFields = map[string]string{
	"staff_profile_id": admissionsbus.OrderByStaffProfileID,
	"user_id":          admissionsbus.OrderByStaffProfileUser,
	"date_created":     admissionsbus.OrderByStaffProfileDateCreated,
}

var applicantProfileOrderByFields = map[string]string{
	"applicant_profile_id": admissionsbus.OrderByApplicantProfileID,
	"user_id":              admissionsbus.OrderByApplicantProfileUser,
	"constituent_id":       admissionsbus.OrderByApplicantProfileConstituent,
	"date_created":         admissionsbus.OrderByApplicantProfileDateCreated,
}

var leadScoreRuleOrderByFields = map[string]string{
	"lead_score_rule_id": admissionsbus.OrderByLeadScoreRuleID,
	"name":               admissionsbus.OrderByLeadScoreRuleName,
	"priority":           admissionsbus.OrderByLeadScoreRulePriority,
	"date_created":       admissionsbus.OrderByLeadScoreRuleDateCreated,
}

var leadScoreOrderByFields = map[string]string{
	"lead_score_id":   admissionsbus.OrderByLeadScoreID,
	"total_score":     admissionsbus.OrderByLeadScoreTotalScore,
	"band":            admissionsbus.OrderByLeadScoreBand,
	"recalculated_at": admissionsbus.OrderByLeadScoreRecalculatedAt,
}

var constituentOrderByFields = map[string]string{
	"constituent_id":  admissionsbus.OrderByConstituentID,
	"last_name":       admissionsbus.OrderByConstituentLastName,
	"primary_email":   admissionsbus.OrderByConstituentPrimaryEmail,
	"lifecycle_stage": admissionsbus.OrderByConstituentLifecycleStage,
}

var inquiryOrderByFields = map[string]string{
	"inquiry_id":    admissionsbus.OrderByInquiryID,
	"primary_email": admissionsbus.OrderByInquiryEmail,
	"source":        admissionsbus.OrderByInquirySource,
	"status":        admissionsbus.OrderByInquiryStatus,
	"date_created":  admissionsbus.OrderByInquiryDateCreated,
}

var academicTermOrderByFields = map[string]string{
	"academic_term_id": admissionsbus.OrderByAcademicTermID,
	"name":             admissionsbus.OrderByAcademicTermName,
	"code":             admissionsbus.OrderByAcademicTermCode,
	"start_date":       admissionsbus.OrderByAcademicTermStartDate,
	"active":           admissionsbus.OrderByAcademicTermActive,
}

var duplicateReviewOrderByFields = map[string]string{
	"duplicate_review_id": admissionsbus.OrderByDuplicateReviewID,
	"status":              admissionsbus.OrderByDuplicateReviewStatus,
	"match_type":          admissionsbus.OrderByDuplicateReviewMatchType,
	"match_score":         admissionsbus.OrderByDuplicateReviewMatchScore,
	"date_created":        admissionsbus.OrderByDuplicateReviewDateCreated,
}

var applicationOrderByFields = map[string]string{
	"application_id":   admissionsbus.OrderByApplicationID,
	"status":           admissionsbus.OrderByApplicationStatus,
	"application_type": admissionsbus.OrderByApplicationType,
	"submitted_at":     admissionsbus.OrderByApplicationSubmittedAt,
	"date_created":     admissionsbus.OrderByApplicationDateCreated,
}

var applicationFormTemplateOrderByFields = map[string]string{
	"form_template_id": admissionsbus.OrderByApplicationFormTemplateID,
	"name":             admissionsbus.OrderByApplicationFormTemplateName,
	"application_type": admissionsbus.OrderByApplicationFormTemplateType,
	"version":          admissionsbus.OrderByApplicationFormTemplateVersion,
	"priority":         admissionsbus.OrderByApplicationFormTemplatePriority,
	"date_created":     admissionsbus.OrderByApplicationFormTemplateDateCreated,
}

var customFieldDefinitionOrderByFields = map[string]string{
	"custom_field_definition_id": admissionsbus.OrderByCustomFieldDefinitionID,
	"owner":                      admissionsbus.OrderByCustomFieldDefinitionOwner,
	"field_key":                  admissionsbus.OrderByCustomFieldDefinitionFieldKey,
	"display_order":              admissionsbus.OrderByCustomFieldDefinitionDisplayOrder,
	"date_created":               admissionsbus.OrderByCustomFieldDefinitionDateCreated,
}

var customFieldValueOrderByFields = map[string]string{
	"custom_field_value_id":      admissionsbus.OrderByCustomFieldValueID,
	"custom_field_definition_id": admissionsbus.OrderByCustomFieldValueDefinition,
	"owner":                      admissionsbus.OrderByCustomFieldValueOwner,
	"owner_id":                   admissionsbus.OrderByCustomFieldValueOwnerID,
	"date_created":               admissionsbus.OrderByCustomFieldValueDateCreated,
}

var applicationTransitionOrderByFields = map[string]string{
	"application_transition_id": admissionsbus.OrderByApplicationTransitionID,
	"application_id":            admissionsbus.OrderByApplicationTransitionApplication,
	"actor_id":                  admissionsbus.OrderByApplicationTransitionActor,
	"date_created":              admissionsbus.OrderByApplicationTransitionDateCreated,
}

var checklistItemOrderByFields = map[string]string{
	"checklist_item_id": admissionsbus.OrderByChecklistItemID,
	"application_id":    admissionsbus.OrderByChecklistItemApplication,
	"status":            admissionsbus.OrderByChecklistItemStatus,
	"display_order":     admissionsbus.OrderByChecklistItemDisplayOrder,
	"date_created":      admissionsbus.OrderByChecklistItemDateCreated,
}

var documentOrderByFields = map[string]string{
	"document_id":       admissionsbus.OrderByDocumentID,
	"application_id":    admissionsbus.OrderByDocumentApplication,
	"checklist_item_id": admissionsbus.OrderByDocumentChecklistItem,
	"status":            admissionsbus.OrderByDocumentStatus,
	"uploaded_at":       admissionsbus.OrderByDocumentUploadedAt,
	"reviewed_at":       admissionsbus.OrderByDocumentReviewedAt,
}

var importBatchOrderByFields = map[string]string{
	"import_batch_id": admissionsbus.OrderByImportBatchID,
	"target":          admissionsbus.OrderByImportBatchTarget,
	"status":          admissionsbus.OrderByImportBatchStatus,
	"uploaded_by_id":  admissionsbus.OrderByImportBatchUploadedBy,
	"date_created":    admissionsbus.OrderByImportBatchDateCreated,
}

var importInvalidRowOrderByFields = map[string]string{
	"import_invalid_row_id": admissionsbus.OrderByImportInvalidRowID,
	"import_batch_id":       admissionsbus.OrderByImportInvalidRowBatch,
	"row_number":            admissionsbus.OrderByImportInvalidRowNumber,
	"date_created":          admissionsbus.OrderByImportInvalidRowDateCreated,
}

var syncJobOrderByFields = map[string]string{
	"sync_job_id":  admissionsbus.OrderBySyncJobID,
	"status":       admissionsbus.OrderBySyncJobStatus,
	"direction":    admissionsbus.OrderBySyncJobDirection,
	"started_at":   admissionsbus.OrderBySyncJobStartedAt,
	"completed_at": admissionsbus.OrderBySyncJobCompletedAt,
	"date_created": admissionsbus.OrderBySyncJobDateCreated,
}

var syncEventOrderByFields = map[string]string{
	"sync_event_id": admissionsbus.OrderBySyncEventID,
	"sync_job_id":   admissionsbus.OrderBySyncEventJob,
	"event_type":    admissionsbus.OrderBySyncEventType,
	"status":        admissionsbus.OrderBySyncEventStatus,
	"direction":     admissionsbus.OrderBySyncEventDirection,
	"resource_id":   admissionsbus.OrderBySyncEventResource,
	"date_created":  admissionsbus.OrderBySyncEventDateCreated,
}
