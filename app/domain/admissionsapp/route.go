package admissionsapp

import (
	"net/http"

	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log           *logger.Logger
	AdmissionsBus admissionsbus.ExtBusiness
	AuthClient    authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAny := mid.Authorize(cfg.AuthClient, auth.RuleAny)

	api := newApp(cfg.AdmissionsBus)

	app.HandlerFunc(http.MethodGet, version, "/admissions/health", api.health, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents", api.queryConstituents, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents/{constituent_id}", api.queryConstituentByID, authen, ruleAny)
	app.HandlerFunc(http.MethodPost, version, "/admissions/constituents", api.createConstituent, authen, ruleAny)
	app.HandlerFunc(http.MethodPut, version, "/admissions/constituents/{constituent_id}", api.updateConstituent, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs", api.queryPrograms, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs/{program_id}", api.queryProgramByID, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms", api.queryAcademicTerms, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms/{academic_term_id}", api.queryAcademicTermByID, authen, ruleAny)
}
