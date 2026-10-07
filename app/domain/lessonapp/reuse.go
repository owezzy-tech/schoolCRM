package lessonapp

import (
	"context"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/jsonbody"
	"github.com/owezzy/schoolCRM/foundation/web"
	"net/http"
)

func (a *appHandlers) reuse(ctx context.Context, r *http.Request) web.Encoder {
	actorID, planID, err := planTarget(ctx, r)
	if err != nil {
		return err
	}
	var req ReuseRequest
	if err := jsonbody.Decode(r, &req, 16*1024); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	plan, busErr := a.bus.Reuse(ctx, actorID, planID, req.Version, req.Title, req.RequestID)
	if busErr != nil {
		return appError(busErr)
	}
	return LessonPlan(plan)
}
