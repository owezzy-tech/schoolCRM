// Package lessonapp exposes lesson plan versions and their publication
// workflow through JSON:API.
package lessonapp

import (
	"net/http"

	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type Config struct {
	LessonBus  *lessonbus.Business
	AuthClient authclient.Authenticator
}

func Routes(app *web.App, cfg Config) {
	api := &appHandlers{bus: cfg.LessonBus}
	authenticate := mid.Authenticate(cfg.AuthClient)
	const department = "/schools/{school_id}/departments/{department_id}/lessons"
	const version = "/lessons/{plan_id}/versions/{version}"
	app.HandlerFunc(http.MethodGet, "v1", department, api.plans, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", department, api.create, authenticate)
	app.HandlerFunc(http.MethodGet, "v1", "/lessons/{plan_id}", api.plan, authenticate)
	app.HandlerFunc(http.MethodGet, "v1", "/lessons/{plan_id}/versions", api.versions, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", "/lessons/{plan_id}/versions", api.revise, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", version+"/submit", api.submit, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", version+"/review", api.review, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", version+"/approval", api.approve, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", version+"/publish", api.publish, authenticate)
}
