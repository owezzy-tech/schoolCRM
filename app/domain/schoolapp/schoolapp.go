package schoolapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/jsonbody"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/app/sdk/query"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type appHandlers struct{ bus *schoolbus.Business }

func (a *appHandlers) createSchool(ctx context.Context, r *http.Request) web.Encoder {
	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}
	var req NamedRequest
	if err := decodeRequest(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	school, err := a.bus.CreateSchool(ctx, actorID, req.Name)
	if err != nil {
		return appError(err)
	}
	return School(school)
}

func (a *appHandlers) createDepartment(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, err := requestScope(ctx, r)
	if err != nil {
		return err
	}
	var req NamedRequest
	if err := decodeRequest(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	department, busErr := a.bus.CreateDepartment(ctx, actorID, schoolID, req.Name)
	if busErr != nil {
		return appError(busErr)
	}
	return Department(department)
}

func (a *appHandlers) grant(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, err := requestScope(ctx, r)
	if err != nil {
		return err
	}
	var req GrantRequest
	if err := decodeRequest(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	membership, busErr := a.bus.Grant(ctx, actorID, req.UserID, schoolbus.Scope{SchoolID: schoolID, DepartmentID: req.DepartmentID}, req.Capability)
	if busErr != nil {
		return appError(busErr)
	}
	return Membership(membership)
}

func (a *appHandlers) revoke(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, err := requestScope(ctx, r)
	if err != nil {
		return err
	}
	membershipID, parseErr := uuid.Parse(web.Param(r, "membership_id"))
	if parseErr != nil || membershipID == uuid.Nil {
		return errs.Errorf(errs.InvalidArgument, "invalid membership ID")
	}
	if err := a.bus.Revoke(ctx, actorID, schoolID, membershipID); err != nil {
		return appError(err)
	}
	return nil
}

func (a *appHandlers) schools(ctx context.Context, r *http.Request) web.Encoder {
	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}
	schools, err := a.bus.Schools(ctx, actorID)
	if err != nil {
		return appError(err)
	}
	return query.Result[schoolbus.School]{Items: schools, Total: len(schools), Page: 1, RowsPerPage: max(1, len(schools))}
}

func (a *appHandlers) departments(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, err := requestScope(ctx, r)
	if err != nil {
		return err
	}
	departments, busErr := a.bus.Departments(ctx, actorID, schoolID)
	if busErr != nil {
		return appError(busErr)
	}
	return query.Result[schoolbus.Department]{Items: departments, Total: len(departments), Page: 1, RowsPerPage: max(1, len(departments))}
}

func (a *appHandlers) memberships(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, err := requestScope(ctx, r)
	if err != nil {
		return err
	}
	memberships, busErr := a.bus.Memberships(ctx, actorID, schoolID)
	if busErr != nil {
		return appError(busErr)
	}
	return query.Result[schoolbus.Membership]{Items: memberships, Total: len(memberships), Page: 1, RowsPerPage: max(1, len(memberships))}
}

func (a *appHandlers) ownMemberships(ctx context.Context, r *http.Request) web.Encoder {
	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}
	memberships, err := a.bus.OwnMemberships(ctx, actorID)
	if err != nil {
		return appError(err)
	}
	return query.Result[schoolbus.Membership]{Items: memberships, Total: len(memberships), Page: 1, RowsPerPage: max(1, len(memberships))}
}

func requestScope(ctx context.Context, r *http.Request) (uuid.UUID, uuid.UUID, *errs.Error) {
	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, errs.New(errs.Unauthenticated, err)
	}
	schoolID, err := uuid.Parse(web.Param(r, "school_id"))
	if err != nil || schoolID == uuid.Nil {
		return uuid.Nil, uuid.Nil, errs.Errorf(errs.InvalidArgument, "invalid school ID")
	}
	return actorID, schoolID, nil
}

func decodeRequest(r *http.Request, v any) error {
	return jsonbody.Decode(r, v, 16*1024)
}

func appError(err error) *errs.Error {
	switch {
	case errors.Is(err, schoolbus.ErrForbidden):
		return errs.New(errs.PermissionDenied, schoolbus.ErrForbidden)
	case errors.Is(err, schoolbus.ErrNotFound):
		return errs.New(errs.NotFound, schoolbus.ErrNotFound)
	case errors.Is(err, schoolbus.ErrInvalid):
		return errs.New(errs.InvalidArgument, err)
	default:
		return errs.New(errs.InternalOnlyLog, err)
	}
}
