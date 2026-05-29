// Package admissionsapp maintains the app layer API for the admissions domain.
package admissionsapp

import (
	"context"
	"net/http"

	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	admissionsBus admissionsbus.ExtBusiness
}

func newApp(admissionsBus admissionsbus.ExtBusiness) *app {
	return &app{
		admissionsBus: admissionsBus,
	}
}

func (a *app) health(ctx context.Context, _ *http.Request) web.Encoder {
	health, err := a.admissionsBus.Health(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "admissions health: %s", err)
	}

	return toAppHealth(health)
}
