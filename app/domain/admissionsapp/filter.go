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
