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
