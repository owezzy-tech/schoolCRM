// Package build manages different build options.
package build

import (
	"github.com/owezzy/schoolCRM/app/domain/authapp"
	"github.com/owezzy/schoolCRM/app/domain/checkapp"
	"github.com/owezzy/schoolCRM/app/sdk/mux"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Routes binds all the routes for the auth service.
func Routes() all {
	return all{}
}

type all struct{}

// Add implements the RouterAdder interface.
func (all) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	authapp.Routes(app, authapp.Config{
		UserBus:  cfg.BusConfig.UserBus,
		Auth:     cfg.AuthConfig.Auth,
		TokenKey: cfg.AuthConfig.TokenKey,
	})
}
