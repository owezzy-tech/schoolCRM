package admissionsapp

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

type programQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ExternalSISID string
	Code          string
	Active        string
}

type constituentQueryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	PrimaryEmail    string
	ExternalSISID   string
	LifecycleStage  string
	DuplicateStatus string
}

type academicTermQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ExternalSISID string
	Code          string
	Active        string
}

type duplicateReviewQueryParams struct {
	Page                   string
	Rows                   string
	OrderBy                string
	ID                     string
	SourceConstituentID    string
	CandidateConstituentID string
	MatchType              string
	Status                 string
}

func parseProgramQueryParams(r *http.Request) programQueryParams {
	values := r.URL.Query()

	return programQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("program_id"),
		ExternalSISID: values.Get("external_sis_id"),
		Code:          values.Get("code"),
		Active:        values.Get("active"),
	}
}

func parseConstituentQueryParams(r *http.Request) constituentQueryParams {
	values := r.URL.Query()

	return constituentQueryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("constituent_id"),
		PrimaryEmail:    values.Get("primary_email"),
		ExternalSISID:   values.Get("external_sis_id"),
		LifecycleStage:  values.Get("lifecycle_stage"),
		DuplicateStatus: values.Get("duplicate_status"),
	}
}

func parseAcademicTermQueryParams(r *http.Request) academicTermQueryParams {
	values := r.URL.Query()

	return academicTermQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("academic_term_id"),
		ExternalSISID: values.Get("external_sis_id"),
		Code:          values.Get("code"),
		Active:        values.Get("active"),
	}
}

func parseDuplicateReviewQueryParams(r *http.Request) duplicateReviewQueryParams {
	values := r.URL.Query()

	return duplicateReviewQueryParams{
		Page:                   values.Get("page"),
		Rows:                   values.Get("rows"),
		OrderBy:                values.Get("orderBy"),
		ID:                     values.Get("duplicate_review_id"),
		SourceConstituentID:    values.Get("source_constituent_id"),
		CandidateConstituentID: values.Get("candidate_constituent_id"),
		MatchType:              values.Get("match_type"),
		Status:                 values.Get("status"),
	}
}

func parseProgramFilter(qp programQueryParams) (admissionsbus.ProgramQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ProgramQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("program_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.Code != "" {
		filter.Code = &qp.Code
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if fieldErrors != nil {
		return admissionsbus.ProgramQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseConstituentFilter(qp constituentQueryParams) (admissionsbus.ConstituentQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ConstituentQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.PrimaryEmail != "" {
		email, err := mail.ParseAddress(qp.PrimaryEmail)
		if err != nil {
			fieldErrors.Add("primary_email", err)
		} else {
			filter.PrimaryEmail = email
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.LifecycleStage != "" {
		stage := admissionsbus.LifecycleStage(qp.LifecycleStage)
		filter.LifecycleStage = &stage
	}

	if qp.DuplicateStatus != "" {
		status := admissionsbus.DuplicateStatus(qp.DuplicateStatus)
		filter.DuplicateStatus = &status
	}

	if fieldErrors != nil {
		return admissionsbus.ConstituentQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseAcademicTermFilter(qp academicTermQueryParams) (admissionsbus.AcademicTermQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.AcademicTermQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("academic_term_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.Code != "" {
		filter.Code = &qp.Code
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if fieldErrors != nil {
		return admissionsbus.AcademicTermQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseDuplicateReviewFilter(qp duplicateReviewQueryParams) (admissionsbus.DuplicateReviewQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.DuplicateReviewQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("duplicate_review_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.SourceConstituentID != "" {
		id, err := uuid.Parse(qp.SourceConstituentID)
		if err != nil {
			fieldErrors.Add("source_constituent_id", err)
		} else {
			filter.SourceConstituentID = &id
		}
	}

	if qp.CandidateConstituentID != "" {
		id, err := uuid.Parse(qp.CandidateConstituentID)
		if err != nil {
			fieldErrors.Add("candidate_constituent_id", err)
		} else {
			filter.CandidateConstituentID = &id
		}
	}

	if qp.MatchType != "" {
		matchType := admissionsbus.DuplicateReviewMatchType(qp.MatchType)
		filter.MatchType = &matchType
	}

	if qp.Status != "" {
		status := admissionsbus.DuplicateReviewStatus(qp.Status)
		filter.Status = &status
	}

	if fieldErrors != nil {
		return admissionsbus.DuplicateReviewQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
