//go:build crud

package build

import (
	"github.com/owezzy/schoolCRM/app/domain/auditapp"
	"github.com/owezzy/schoolCRM/app/domain/checkapp"
	"github.com/owezzy/schoolCRM/app/domain/homeapp"
	"github.com/owezzy/schoolCRM/app/domain/productapp"
	"github.com/owezzy/schoolCRM/app/domain/tranapp"
	"github.com/owezzy/schoolCRM/app/domain/userapp"
	"github.com/owezzy/schoolCRM/app/sdk/mux"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Routes binds the crud routes for the sales service.
func Routes() crud {
	return crud{}
}

type crud struct{}

// Add implements the RouterAdder interface.
func (crud) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	tranapp.Routes(app, tranapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		Log:        cfg.Log,
		AuthClient: cfg.SalesConfig.AuthClient,
		DB:         cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})
}
