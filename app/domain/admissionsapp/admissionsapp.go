// Package admissionsapp maintains the app layer API for the admissions domain.
package admissionsapp

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/query"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	admissionsBus admissionsbus.ExtBusiness
}

func newApp(admissionsBus admissionsbus.ExtBusiness) *app {
	return &app{
		admissionsBus: admissionsBus,
	}
}

func (a *app) health(ctx context.Context, _ *http.Request) web.Encoder {
	health, err := a.admissionsBus.Health(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "admissions health: %s", err)
	}

	return toAppHealth(health)
}

func (a *app) queryPrograms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseProgramQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseProgramFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(programOrderByFields, qp.OrderBy, admissionsbus.DefaultProgramOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	programs, err := a.admissionsBus.QueryPrograms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query programs: %s", err)
	}

	total, err := a.admissionsBus.CountPrograms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count programs: %s", err)
	}

	return query.NewResult(toAppPrograms(programs), total, page)
}

func (a *app) queryProgramByID(ctx context.Context, r *http.Request) web.Encoder {
	programID, err := uuid.Parse(web.Param(r, "program_id"))
	if err != nil {
		return errs.NewFieldErrors("program_id", err)
	}

	program, err := a.admissionsBus.QueryProgramByID(ctx, programID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query program: %s", err)
	}

	return toAppProgram(program)
}

func (a *app) queryAcademicTerms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseAcademicTermQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseAcademicTermFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(academicTermOrderByFields, qp.OrderBy, admissionsbus.DefaultAcademicTermOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	terms, err := a.admissionsBus.QueryAcademicTerms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query academic terms: %s", err)
	}

	total, err := a.admissionsBus.CountAcademicTerms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count academic terms: %s", err)
	}

	return query.NewResult(toAppAcademicTerms(terms), total, page)
}

func (a *app) queryAcademicTermByID(ctx context.Context, r *http.Request) web.Encoder {
	termID, err := uuid.Parse(web.Param(r, "academic_term_id"))
	if err != nil {
		return errs.NewFieldErrors("academic_term_id", err)
	}

	term, err := a.admissionsBus.QueryAcademicTermByID(ctx, termID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query academic term: %s", err)
	}

	return toAppAcademicTerm(term)
}
