package admissionsapp

import (
	"net/http"
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
