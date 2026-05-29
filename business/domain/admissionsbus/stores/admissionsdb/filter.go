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
		wc = append(wc, "status NOT IN ('DENIED', 'WITHDRAWN', 'ENROLLED')")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
