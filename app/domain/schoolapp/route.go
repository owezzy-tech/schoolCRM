// Package schoolapp exposes school-scoped administration through JSON:API.
package schoolapp

import (
	"net/http"

	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type Config struct {
	SchoolBus  *schoolbus.Business
	AuthClient authclient.Authenticator
}

func Routes(app *web.App, cfg Config) {
	api := &appHandlers{bus: cfg.SchoolBus}
	authenticate := mid.Authenticate(cfg.AuthClient)
	app.HandlerFunc(http.MethodGet, "v1", "/me/school-memberships", api.ownMemberships, authenticate)
	app.HandlerFunc(http.MethodGet, "v1", "/schools", api.schools, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", "/schools", api.createSchool, authenticate)
	app.HandlerFunc(http.MethodGet, "v1", "/schools/{school_id}/departments", api.departments, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", "/schools/{school_id}/departments", api.createDepartment, authenticate)
	app.HandlerFunc(http.MethodGet, "v1", "/schools/{school_id}/memberships", api.memberships, authenticate)
	app.HandlerFunc(http.MethodPost, "v1", "/schools/{school_id}/memberships", api.grant, authenticate)
	app.HandlerFunc(http.MethodDelete, "v1", "/schools/{school_id}/memberships/{membership_id}", api.revoke, authenticate)
}
