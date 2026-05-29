package admissionsapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"github.com/owezzy/schoolCRM/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log           *logger.Logger
	DB            *sqlx.DB
	AdmissionsBus admissionsbus.ExtBusiness
	AuditBus      auditbus.ExtBusiness
	AuthClient    authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAny := mid.Authorize(cfg.AuthClient, auth.RuleAny)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	api := newApp(cfg.AdmissionsBus, cfg.AuditBus)

	app.HandlerFunc(http.MethodGet, version, "/admissions/health", api.health, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents", api.queryConstituents, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/constituents/{constituent_id}", api.queryConstituentByID, authen, ruleAny)
	app.HandlerFunc(http.MethodPost, version, "/admissions/constituents", api.createConstituent, authen, ruleAny)
	app.HandlerFunc(http.MethodPut, version, "/admissions/constituents/{constituent_id}", api.updateConstituent, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs", api.queryPrograms, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/programs/{program_id}", api.queryProgramByID, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms", api.queryAcademicTerms, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/academic-terms/{academic_term_id}", api.queryAcademicTermByID, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications", api.queryApplications, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}", api.queryApplicationByID, authen, ruleAny)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications", api.createApplication, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/applications/{application_id}/transitions", api.queryApplicationTransitions, authen, ruleAny)
	app.HandlerFunc(http.MethodPost, version, "/admissions/applications/{application_id}/transitions", api.transitionApplicationStatus, authen, ruleAny, transaction)
	app.HandlerFunc(http.MethodGet, version, "/admissions/duplicate-reviews", api.queryDuplicateReviews, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/admissions/duplicate-reviews/{duplicate_review_id}", api.queryDuplicateReviewByID, authen, ruleAny)
	app.HandlerFunc(http.MethodPost, version, "/admissions/duplicate-reviews", api.createDuplicateReview, authen, ruleAny)
	app.HandlerFunc(http.MethodPut, version, "/admissions/duplicate-reviews/{duplicate_review_id}/resolution", api.resolveDuplicateReview, authen, ruleAny)
}
