package authapp

import (
	"net/http"

	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	UserBus userbus.ExtBusiness
	Auth    *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	bearer := mid.Bearer(cfg.Auth)
	basic := mid.Basic(cfg.Auth, cfg.UserBus)

	api := newApp(cfg.Auth)

	app.HandlerFunc(http.MethodGet, version, "/auth/token/{kid}", api.token, basic)
	app.HandlerFunc(http.MethodGet, version, "/auth/authenticate", api.authenticate, bearer)
	app.HandlerFunc(http.MethodPost, version, "/auth/authorize", api.authorize)
}
