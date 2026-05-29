//go:build !crud && !reporting

// Package build manages different build options.
package build

import (
	"github.com/owezzy/schoolCRM/app/domain/admissionsapp"
	"github.com/owezzy/schoolCRM/app/domain/auditapp"
	"github.com/owezzy/schoolCRM/app/domain/checkapp"
	"github.com/owezzy/schoolCRM/app/domain/homeapp"
	"github.com/owezzy/schoolCRM/app/domain/productapp"
	"github.com/owezzy/schoolCRM/app/domain/rawapp"
	"github.com/owezzy/schoolCRM/app/domain/tranapp"
	"github.com/owezzy/schoolCRM/app/domain/userapp"
	"github.com/owezzy/schoolCRM/app/domain/vproductapp"
	"github.com/owezzy/schoolCRM/app/sdk/mux"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Routes binds all the routes for the SchoolCRM service.
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

	homeapp.Routes(app, homeapp.Config{
		Log:        cfg.Log,
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SchoolCRMConfig.AuthClient,
	})

	admissionsapp.Routes(app, admissionsapp.Config{
		Log:           cfg.Log,
		AdmissionsBus: cfg.BusConfig.AdmissionsBus,
		AuditBus:      cfg.BusConfig.AuditBus,
		AuthClient:    cfg.SchoolCRMConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		Log:        cfg.Log,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SchoolCRMConfig.AuthClient,
	})

	rawapp.Routes(app)

	tranapp.Routes(app, tranapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SchoolCRMConfig.AuthClient,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SchoolCRMConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.SchoolCRMConfig.AuthClient,
	})

	vproductapp.Routes(app, vproductapp.Config{
		Log:         cfg.Log,
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.SchoolCRMConfig.AuthClient,
	})
}
