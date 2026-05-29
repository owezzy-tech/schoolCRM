// Package admissionsapp maintains the app layer API for the admissions domain.
package admissionsapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/app/sdk/query"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/name"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	admissionsBus admissionsbus.ExtBusiness
	auditBus      auditbus.ExtBusiness
}

func newApp(admissionsBus admissionsbus.ExtBusiness, auditBus auditbus.ExtBusiness) *app {
	return &app{
		admissionsBus: admissionsBus,
		auditBus:      auditBus,
	}
}

func (a *app) health(ctx context.Context, _ *http.Request) web.Encoder {
	health, err := a.admissionsBus.Health(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "admissions health: %s", err)
	}

	return toAppHealth(health)
}

func (a *app) createConstituent(ctx context.Context, r *http.Request) web.Encoder {
	var app NewConstituent
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nc, err := toBusNewConstituent(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cst, err := a.admissionsBus.CreateConstituent(ctx, nc)
	if err != nil {
		return errs.Errorf(errs.Internal, "create constituent: %s", err)
	}

	return toAppConstituent(cst)
}

func (a *app) updateConstituent(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateConstituent
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	cst, err := a.admissionsBus.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituent: %s", err)
	}

	uc, err := toBusUpdateConstituent(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.UpdateConstituent(ctx, cst, uc)
	if err != nil {
		return errs.Errorf(errs.Internal, "update constituent: %s", err)
	}

	return toAppConstituent(updated)
}

func (a *app) queryConstituents(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseConstituentQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseConstituentFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(constituentOrderByFields, qp.OrderBy, admissionsbus.DefaultConstituentOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	constituents, err := a.admissionsBus.QueryConstituents(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituents: %s", err)
	}

	total, err := a.admissionsBus.CountConstituents(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count constituents: %s", err)
	}

	return query.NewResult(toAppConstituents(constituents), total, page)
}

func (a *app) queryConstituentByID(ctx context.Context, r *http.Request) web.Encoder {
	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	cst, err := a.admissionsBus.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituent: %s", err)
	}

	return toAppConstituent(cst)
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

func (a *app) createDuplicateReview(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDuplicateReview
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nr, err := toBusNewDuplicateReview(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	review, err := a.admissionsBus.CreateDuplicateReview(ctx, nr)
	if err != nil {
		return errs.Errorf(errs.Internal, "create duplicate review: %s", err)
	}

	return toAppDuplicateReview(review)
}

func (a *app) resolveDuplicateReview(ctx context.Context, r *http.Request) web.Encoder {
	var app ResolveDuplicateReview
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	reviewID, err := uuid.Parse(web.Param(r, "duplicate_review_id"))
	if err != nil {
		return errs.NewFieldErrors("duplicate_review_id", err)
	}

	review, err := a.admissionsBus.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate review: %s", err)
	}

	rr := toBusResolveDuplicateReview(app, mid.GetSubjectID(ctx))
	resolved, err := a.admissionsBus.ResolveDuplicateReview(ctx, review, rr)
	if err != nil {
		return errs.Errorf(errs.Internal, "resolve duplicate review: %s", err)
	}

	if a.auditBus != nil {
		na := auditbus.NewAudit{
			ObjID:     resolved.ID,
			ObjDomain: domain.Admissions,
			ObjName:   name.MustParse("Duplicate Review"),
			ActorID:   rr.ActorID,
			Action:    "duplicate_" + rr.Resolution.String(),
			Data:      toAppDuplicateReview(resolved),
			Message:   "duplicate review resolved",
		}

		if _, err := a.auditBus.Create(ctx, na); err != nil {
			return errs.Errorf(errs.Internal, "audit duplicate review: %s", err)
		}
	}

	return toAppDuplicateReview(resolved)
}

func (a *app) queryDuplicateReviews(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseDuplicateReviewQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseDuplicateReviewFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(duplicateReviewOrderByFields, qp.OrderBy, admissionsbus.DefaultDuplicateReviewOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	reviews, err := a.admissionsBus.QueryDuplicateReviews(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate reviews: %s", err)
	}

	total, err := a.admissionsBus.CountDuplicateReviews(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count duplicate reviews: %s", err)
	}

	return query.NewResult(toAppDuplicateReviews(reviews), total, page)
}

func (a *app) queryDuplicateReviewByID(ctx context.Context, r *http.Request) web.Encoder {
	reviewID, err := uuid.Parse(web.Param(r, "duplicate_review_id"))
	if err != nil {
		return errs.NewFieldErrors("duplicate_review_id", err)
	}

	review, err := a.admissionsBus.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate review: %s", err)
	}

	return toAppDuplicateReview(review)
}

func (a *app) createApplication(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplication
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	na, err := toBusNewApplication(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	application, err := a.admissionsBus.CreateApplication(ctx, na)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrDuplicateApplication) {
			return errs.New(errs.Aborted, admissionsbus.ErrDuplicateApplication)
		}
		return errs.Errorf(errs.Internal, "create application: %s", err)
	}

	return toAppApplication(application)
}

func (a *app) queryApplications(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(applicationOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	applications, err := a.admissionsBus.QueryApplications(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applications: %s", err)
	}

	total, err := a.admissionsBus.CountApplications(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applications: %s", err)
	}

	return query.NewResult(toAppApplications(applications), total, page)
}

func (a *app) queryApplicationByID(ctx context.Context, r *http.Request) web.Encoder {
	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	application, err := a.admissionsBus.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application: %s", err)
	}

	return toAppApplication(application)
}
