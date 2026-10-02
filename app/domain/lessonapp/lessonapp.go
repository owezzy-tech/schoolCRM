package lessonapp

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/jsonbody"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/app/sdk/query"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// bodyLimit leaves room for JSON framing around the largest permitted content.
const bodyLimit = lessonbus.MaxContentBytes + 16*1024

type appHandlers struct{ bus *lessonbus.Business }

func (a *appHandlers) create(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, departmentID, err := departmentTarget(ctx, r)
	if err != nil {
		return err
	}
	var req NewPlanRequest
	if err := jsonbody.Decode(r, &req, bodyLimit); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	plan, busErr := a.bus.Create(ctx, actorID, schoolID, departmentID, req.toBus())
	if busErr != nil {
		return appError(busErr)
	}
	return LessonPlan(plan)
}

func (a *appHandlers) revise(ctx context.Context, r *http.Request) web.Encoder {
	actorID, planID, err := planTarget(ctx, r)
	if err != nil {
		return err
	}
	var req RevisionRequest
	if err := jsonbody.Decode(r, &req, bodyLimit); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	if req.BaseVersion < 1 {
		return errs.Errorf(errs.InvalidArgument, "baseVersion is required")
	}
	version, busErr := a.bus.Revise(ctx, actorID, planID, req.BaseVersion, req.toBus())
	if busErr != nil {
		return appError(busErr)
	}
	return toLessonVersion(version)
}

func (a *appHandlers) submit(ctx context.Context, r *http.Request) web.Encoder {
	return a.step(ctx, r, a.bus.Submit)
}

func (a *appHandlers) publish(ctx context.Context, r *http.Request) web.Encoder {
	return a.step(ctx, r, a.bus.Publish)
}

func (a *appHandlers) review(ctx context.Context, r *http.Request) web.Encoder {
	return a.decide(ctx, r, a.bus.Review)
}

func (a *appHandlers) approve(ctx context.Context, r *http.Request) web.Encoder {
	return a.decide(ctx, r, a.bus.Approve)
}

type stepFunc func(ctx context.Context, actorID, planID uuid.UUID, number int) (lessonbus.Version, error)

func (a *appHandlers) step(ctx context.Context, r *http.Request, fn stepFunc) web.Encoder {
	actorID, planID, number, err := versionTarget(ctx, r)
	if err != nil {
		return err
	}
	version, busErr := fn(ctx, actorID, planID, number)
	if busErr != nil {
		return appError(busErr)
	}
	return toLessonVersion(version)
}

type decideFunc func(ctx context.Context, actorID, planID uuid.UUID, number int, decision lessonbus.Decision, feedback string) (lessonbus.Version, error)

func (a *appHandlers) decide(ctx context.Context, r *http.Request, fn decideFunc) web.Encoder {
	actorID, planID, number, err := versionTarget(ctx, r)
	if err != nil {
		return err
	}
	var req DecisionRequest
	if err := jsonbody.Decode(r, &req, 16*1024); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	version, busErr := fn(ctx, actorID, planID, number, req.Decision, req.Feedback)
	if busErr != nil {
		return appError(busErr)
	}
	return toLessonVersion(version)
}

func (a *appHandlers) plan(ctx context.Context, r *http.Request) web.Encoder {
	actorID, planID, err := planTarget(ctx, r)
	if err != nil {
		return err
	}
	plan, busErr := a.bus.Plan(ctx, actorID, planID)
	if busErr != nil {
		return appError(busErr)
	}
	return LessonPlan(plan)
}

func (a *appHandlers) plans(ctx context.Context, r *http.Request) web.Encoder {
	actorID, schoolID, departmentID, err := departmentTarget(ctx, r)
	if err != nil {
		return err
	}
	plans, busErr := a.bus.Plans(ctx, actorID, schoolID, departmentID)
	if busErr != nil {
		return appError(busErr)
	}
	items := make([]LessonPlan, len(plans))
	for i, plan := range plans {
		items[i] = LessonPlan(plan)
	}
	return collection(items)
}

func (a *appHandlers) versions(ctx context.Context, r *http.Request) web.Encoder {
	actorID, planID, err := planTarget(ctx, r)
	if err != nil {
		return err
	}
	versions, busErr := a.bus.Versions(ctx, actorID, planID)
	if busErr != nil {
		return appError(busErr)
	}
	items := make([]LessonVersion, len(versions))
	for i, version := range versions {
		items[i] = toLessonVersion(version)
	}
	return collection(items)
}

func collection[T any](items []T) query.Result[T] {
	return query.Result[T]{Items: items, Total: len(items), Page: 1, RowsPerPage: max(1, len(items))}
}

func actor(ctx context.Context) (uuid.UUID, *errs.Error) {
	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return uuid.Nil, errs.New(errs.Unauthenticated, err)
	}
	return actorID, nil
}

func uuidParam(r *http.Request, key string) (uuid.UUID, *errs.Error) {
	id, err := uuid.Parse(web.Param(r, key))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errs.Errorf(errs.InvalidArgument, "invalid %s", key)
	}
	return id, nil
}

func departmentTarget(ctx context.Context, r *http.Request) (uuid.UUID, uuid.UUID, uuid.UUID, *errs.Error) {
	actorID, err := actor(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	schoolID, err := uuidParam(r, "school_id")
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	departmentID, err := uuidParam(r, "department_id")
	return actorID, schoolID, departmentID, err
}

func planTarget(ctx context.Context, r *http.Request) (uuid.UUID, uuid.UUID, *errs.Error) {
	actorID, err := actor(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	planID, err := uuidParam(r, "plan_id")
	return actorID, planID, err
}

func versionTarget(ctx context.Context, r *http.Request) (uuid.UUID, uuid.UUID, int, *errs.Error) {
	actorID, planID, err := planTarget(ctx, r)
	if err != nil {
		return uuid.Nil, uuid.Nil, 0, err
	}
	number, parseErr := strconv.Atoi(web.Param(r, "version"))
	if parseErr != nil || number < 1 {
		return uuid.Nil, uuid.Nil, 0, errs.Errorf(errs.InvalidArgument, "invalid version")
	}
	return actorID, planID, number, nil
}

func appError(err error) *errs.Error {
	switch {
	case errors.Is(err, lessonbus.ErrForbidden):
		return errs.New(errs.PermissionDenied, lessonbus.ErrForbidden)
	case errors.Is(err, lessonbus.ErrNotFound):
		return errs.New(errs.NotFound, lessonbus.ErrNotFound)
	case errors.Is(err, lessonbus.ErrConflict):
		return errs.New(errs.Aborted, lessonbus.ErrConflict)
	case errors.Is(err, lessonbus.ErrInvalid):
		return errs.New(errs.InvalidArgument, err)
	default:
		return errs.New(errs.InternalOnlyLog, err)
	}
}
