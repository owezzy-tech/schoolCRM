package admissionsdb

import (
	"bytes"
	"strings"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

func (s *Store) applyConstituentFilter(filter admissionsbus.ConstituentQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["constituent_id"] = filter.ID
		wc = append(wc, "constituent_id = :constituent_id")
	}

	if filter.PrimaryEmail != nil {
		data["primary_email"] = filter.PrimaryEmail.String()
		wc = append(wc, "primary_email = :primary_email")
	}

	if filter.ExternalSISID != nil {
		data["external_sis_id"] = filter.ExternalSISID
		wc = append(wc, "external_sis_id = :external_sis_id")
	}

	if filter.NationalID != nil {
		data["national_id"] = filter.NationalID
		wc = append(wc, "national_id = :national_id")
	}

	if filter.UPI != nil {
		data["upi"] = filter.UPI
		wc = append(wc, "upi = :upi")
	}

	if filter.KCSEIndexNumber != nil {
		data["kcse_index_number"] = filter.KCSEIndexNumber
		wc = append(wc, "kcse_index_number = :kcse_index_number")
	}

	if filter.LifecycleStage != nil {
		data["lifecycle_stage"] = filter.LifecycleStage.String()
		wc = append(wc, "lifecycle_stage = :lifecycle_stage")
	}

	if filter.DuplicateStatus != nil {
		data["duplicate_status"] = filter.DuplicateStatus.String()
		wc = append(wc, "duplicate_status = :duplicate_status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyInquiryFilter(filter admissionsbus.InquiryQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["inquiry_id"] = filter.ID
		wc = append(wc, "inquiry_id = :inquiry_id")
	}

	if filter.ConstituentID != nil {
		data["constituent_id"] = filter.ConstituentID
		wc = append(wc, "constituent_id = :constituent_id")
	}

	if filter.PrimaryEmail != nil {
		data["primary_email"] = filter.PrimaryEmail.String()
		wc = append(wc, "primary_email = :primary_email")
	}

	if filter.Source != nil {
		data["source"] = filter.Source
		wc = append(wc, "source = :source")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyStaffProfileFilter(filter admissionsbus.StaffProfileQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["staff_profile_id"] = filter.ID
		wc = append(wc, "staff_profile_id = :staff_profile_id")
	}

	if filter.UserID != nil {
		data["user_id"] = filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Role != nil {
		data["admissions_role"] = filter.Role.String()
		wc = append(wc, ":admissions_role = ANY(admissions_roles)")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyApplicantProfileFilter(filter admissionsbus.ApplicantProfileQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["applicant_profile_id"] = filter.ID
		wc = append(wc, "applicant_profile_id = :applicant_profile_id")
	}

	if filter.UserID != nil {
		data["user_id"] = filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.ConstituentID != nil {
		data["constituent_id"] = filter.ConstituentID
		wc = append(wc, "constituent_id = :constituent_id")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyLeadScoreRuleFilter(filter admissionsbus.LeadScoreRuleQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["lead_score_rule_id"] = filter.ID
		wc = append(wc, "lead_score_rule_id = :lead_score_rule_id")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyLeadScoreFilter(filter admissionsbus.LeadScoreQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["lead_score_id"] = filter.ID
		wc = append(wc, "lead_score_id = :lead_score_id")
	}

	if filter.ConstituentID != nil {
		data["constituent_id"] = filter.ConstituentID
		wc = append(wc, "constituent_id = :constituent_id")
	}

	if filter.Band != nil {
		data["score_band"] = filter.Band.String()
		wc = append(wc, "score_band = :score_band")
	}

	if filter.MinScore != nil {
		data["min_score"] = filter.MinScore
		wc = append(wc, "total_score >= :min_score")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyApplicationFormTemplateFilter(filter admissionsbus.ApplicationFormTemplateQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["form_template_id"] = filter.ID
		wc = append(wc, "form_template_id = :form_template_id")
	}

	if filter.ProgramID != nil {
		data["program_id"] = filter.ProgramID
		wc = append(wc, "program_id = :program_id")
	}

	if filter.AcademicTermID != nil {
		data["academic_term_id"] = filter.AcademicTermID
		wc = append(wc, "academic_term_id = :academic_term_id")
	}

	if filter.ApplicationType != nil {
		data["application_type"] = filter.ApplicationType.String()
		wc = append(wc, "application_type = :application_type")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if filter.Version != nil {
		data["version"] = filter.Version
		wc = append(wc, "version = :version")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyCustomFieldDefinitionFilter(filter admissionsbus.CustomFieldDefinitionQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["custom_field_definition_id"] = filter.ID
		wc = append(wc, "custom_field_definition_id = :custom_field_definition_id")
	}

	if filter.Owner != nil {
		data["owner"] = filter.Owner.String()
		wc = append(wc, "owner = :owner")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyCustomFieldValueFilter(filter admissionsbus.CustomFieldValueQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["custom_field_value_id"] = filter.ID
		wc = append(wc, "custom_field_value_id = :custom_field_value_id")
	}

	if filter.DefinitionID != nil {
		data["custom_field_definition_id"] = filter.DefinitionID
		wc = append(wc, "custom_field_definition_id = :custom_field_definition_id")
	}

	if filter.Owner != nil {
		data["owner"] = filter.Owner.String()
		wc = append(wc, "owner = :owner")
	}

	if filter.OwnerID != nil {
		data["owner_id"] = filter.OwnerID
		wc = append(wc, "owner_id = :owner_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyProgramFilter(filter admissionsbus.ProgramQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["program_id"] = filter.ID
		wc = append(wc, "program_id = :program_id")
	}

	if filter.ExternalSISID != nil {
		data["external_sis_id"] = filter.ExternalSISID
		wc = append(wc, "external_sis_id = :external_sis_id")
	}

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyAcademicTermFilter(filter admissionsbus.AcademicTermQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["academic_term_id"] = filter.ID
		wc = append(wc, "academic_term_id = :academic_term_id")
	}

	if filter.ExternalSISID != nil {
		data["external_sis_id"] = filter.ExternalSISID
		wc = append(wc, "external_sis_id = :external_sis_id")
	}

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Active != nil {
		data["is_active"] = filter.Active
		wc = append(wc, "is_active = :is_active")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyDuplicateReviewFilter(filter admissionsbus.DuplicateReviewQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["duplicate_review_id"] = filter.ID
		wc = append(wc, "duplicate_review_id = :duplicate_review_id")
	}

	if filter.SourceConstituentID != nil {
		data["source_constituent_id"] = filter.SourceConstituentID
		wc = append(wc, "source_constituent_id = :source_constituent_id")
	}

	if filter.CandidateConstituentID != nil {
		data["candidate_constituent_id"] = filter.CandidateConstituentID
		wc = append(wc, "candidate_constituent_id = :candidate_constituent_id")
	}

	if filter.MatchType != nil {
		data["match_type"] = filter.MatchType.String()
		wc = append(wc, "match_type = :match_type")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyApplicationFilter(filter admissionsbus.ApplicationQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["application_id"] = filter.ID
		wc = append(wc, "application_id = :application_id")
	}

	if filter.ConstituentID != nil {
		data["constituent_id"] = filter.ConstituentID
		wc = append(wc, "constituent_id = :constituent_id")
	}

	if filter.ProgramID != nil {
		data["program_id"] = filter.ProgramID
		wc = append(wc, "program_id = :program_id")
	}

	if filter.AcademicTermID != nil {
		data["academic_term_id"] = filter.AcademicTermID
		wc = append(wc, "academic_term_id = :academic_term_id")
	}

	if filter.ApplicationType != nil {
		data["application_type"] = filter.ApplicationType.String()
		wc = append(wc, "application_type = :application_type")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if filter.ActiveOnly != nil && *filter.ActiveOnly {
		// Keep this predicate aligned with admissionsbus.isApplicationActive.
		wc = append(wc, "status NOT IN ('DENIED', 'WITHDRAWN', 'ENROLLED')")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyApplicationTransitionFilter(filter admissionsbus.ApplicationTransitionQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["application_transition_id"] = filter.ID
		wc = append(wc, "application_transition_id = :application_transition_id")
	}

	if filter.ApplicationID != nil {
		data["application_id"] = filter.ApplicationID
		wc = append(wc, "application_id = :application_id")
	}

	if filter.ActorID != nil {
		data["actor_id"] = filter.ActorID
		wc = append(wc, "actor_id = :actor_id")
	}

	if filter.FromStatus != nil {
		data["from_status"] = filter.FromStatus.String()
		wc = append(wc, "from_status = :from_status")
	}

	if filter.ToStatus != nil {
		data["to_status"] = filter.ToStatus.String()
		wc = append(wc, "to_status = :to_status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyChecklistItemFilter(filter admissionsbus.ChecklistItemQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["checklist_item_id"] = filter.ID
		wc = append(wc, "checklist_item_id = :checklist_item_id")
	}

	if filter.ApplicationID != nil {
		data["application_id"] = filter.ApplicationID
		wc = append(wc, "application_id = :application_id")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if filter.Required != nil {
		data["is_required"] = filter.Required
		wc = append(wc, "is_required = :is_required")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyDocumentFilter(filter admissionsbus.DocumentQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["document_id"] = filter.ID
		wc = append(wc, "document_id = :document_id")
	}

	if filter.ApplicationID != nil {
		data["application_id"] = filter.ApplicationID
		wc = append(wc, "application_id = :application_id")
	}

	if filter.ChecklistItemID != nil {
		data["checklist_item_id"] = filter.ChecklistItemID
		wc = append(wc, "checklist_item_id = :checklist_item_id")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if filter.UploadedByID != nil {
		data["uploaded_by_id"] = filter.UploadedByID
		wc = append(wc, "uploaded_by_id = :uploaded_by_id")
	}

	if filter.ReviewerID != nil {
		data["reviewer_id"] = filter.ReviewerID
		wc = append(wc, "reviewer_id = :reviewer_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyImportBatchFilter(filter admissionsbus.ImportBatchQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["import_batch_id"] = filter.ID
		wc = append(wc, "import_batch_id = :import_batch_id")
	}

	if filter.Source != nil {
		data["source"] = filter.Source.String()
		wc = append(wc, "source = :source")
	}

	if filter.FileType != nil {
		data["file_type"] = filter.FileType.String()
		wc = append(wc, "file_type = :file_type")
	}

	if filter.Target != nil {
		data["target"] = filter.Target.String()
		wc = append(wc, "target = :target")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if filter.UploadedByID != nil {
		data["uploaded_by_id"] = filter.UploadedByID
		wc = append(wc, "uploaded_by_id = :uploaded_by_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyImportInvalidRowFilter(filter admissionsbus.ImportInvalidRowQueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["import_invalid_row_id"] = filter.ID
		wc = append(wc, "import_invalid_row_id = :import_invalid_row_id")
	}

	if filter.BatchID != nil {
		data["import_batch_id"] = filter.BatchID
		wc = append(wc, "import_batch_id = :import_batch_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
