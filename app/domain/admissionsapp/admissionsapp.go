// Package admissionsapp maintains the app layer API for the admissions domain.
package admissionsapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/app/sdk/query"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/name"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	admissionsBus admissionsbus.ExtBusiness
	auditBus      auditbus.ExtBusiness
}

func newApp(admissionsBus admissionsbus.ExtBusiness, auditBus auditbus.ExtBusiness) *app {
	return &app{
		admissionsBus: admissionsBus,
		auditBus:      auditBus,
	}
}

func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	admissionsBus, err := a.admissionsBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	return &app{
		admissionsBus: admissionsBus,
		auditBus:      a.auditBus,
	}, nil
}

func (a *app) health(ctx context.Context, _ *http.Request) web.Encoder {
	health, err := a.admissionsBus.Health(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "admissions health: %s", err)
	}

	return toAppHealth(health)
}

func (a *app) createConstituent(ctx context.Context, r *http.Request) web.Encoder {
	var app NewConstituent
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nc, err := toBusNewConstituent(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cst, err := a.admissionsBus.CreateConstituent(ctx, nc)
	if err != nil {
		return errs.Errorf(errs.Internal, "create constituent: %s", err)
	}

	return toAppConstituent(cst)
}

func (a *app) createInquiry(ctx context.Context, r *http.Request) web.Encoder {
	var app NewInquiry
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ni, err := toBusNewInquiry(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	inquiry, err := a.admissionsBus.CreateInquiry(ctx, ni)
	if err != nil {
		return errs.Errorf(errs.Internal, "create inquiry: %s", err)
	}

	return toAppInquiry(inquiry)
}

func (a *app) queryInquiries(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseInquiryQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseInquiryFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(inquiryOrderByFields, qp.OrderBy, admissionsbus.DefaultInquiryOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	inquiries, err := a.admissionsBus.QueryInquiries(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query inquiries: %s", err)
	}

	total, err := a.admissionsBus.CountInquiries(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count inquiries: %s", err)
	}

	return query.NewResult(toAppInquiries(inquiries), total, page)
}

func (a *app) queryInquiryByID(ctx context.Context, r *http.Request) web.Encoder {
	inquiryID, err := uuid.Parse(web.Param(r, "inquiry_id"))
	if err != nil {
		return errs.NewFieldErrors("inquiry_id", err)
	}

	inquiry, err := a.admissionsBus.QueryInquiryByID(ctx, inquiryID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query inquiry: %s", err)
	}

	return toAppInquiry(inquiry)
}

func (a *app) queryEvents(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseEventQueryParams(r)

	pg, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseEventFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(eventOrderByFields, qp.OrderBy, admissionsbus.DefaultEventOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	events, err := a.admissionsBus.QueryEvents(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Errorf(errs.Internal, "query events: %s", err)
	}

	total, err := a.admissionsBus.CountEvents(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count events: %s", err)
	}

	appEvents := make([]Event, len(events))
	for i, event := range events {
		registrations, regErr := a.admissionsBus.QueryEventRegistrations(ctx, admissionsbus.EventRegistrationQueryFilter{EventID: &event.ID}, admissionsbus.DefaultEventRegistrationOrderBy, page.MustParse("1", "5000"))
		if regErr != nil {
			return errs.Errorf(errs.Internal, "query event registrations: %s", regErr)
		}
		appEvents[i] = toAppEvent(event, registrations)
	}

	return query.NewResult(appEvents, total, pg)
}

func (a *app) queryEventByID(ctx context.Context, r *http.Request) web.Encoder {
	eventID, err := uuid.Parse(web.Param(r, "event_id"))
	if err != nil {
		return errs.NewFieldErrors("event_id", err)
	}

	event, err := a.admissionsBus.QueryEventByID(ctx, eventID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query event: %s", err)
	}

	registrations, err := a.admissionsBus.QueryEventRegistrations(ctx, admissionsbus.EventRegistrationQueryFilter{EventID: &eventID}, admissionsbus.DefaultEventRegistrationOrderBy, page.MustParse("1", "5000"))
	if err != nil {
		return errs.Errorf(errs.Internal, "query event registrations: %s", err)
	}

	return toAppEvent(event, registrations)
}

func (a *app) queryEventRegistrations(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseEventRegistrationQueryParams(r)
	if qp.EventID == "" {
		qp.EventID = web.Param(r, "event_id")
	}

	pag, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseEventRegistrationFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(eventRegistrationOrderByFields, qp.OrderBy, admissionsbus.DefaultEventRegistrationOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	registrations, err := a.admissionsBus.QueryEventRegistrations(ctx, filter, orderBy, pag)
	if err != nil {
		return errs.Errorf(errs.Internal, "query event registrations: %s", err)
	}

	total, err := a.admissionsBus.CountEventRegistrations(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count event registrations: %s", err)
	}

	return query.NewResult(toAppEventRegistrations(registrations), total, pag)
}

func (a *app) registerForApplicantEvent(ctx context.Context, r *http.Request) web.Encoder {
	var app NewEventRegistration
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	eventID, err := uuid.Parse(web.Param(r, "event_id"))
	if err != nil {
		return errs.NewFieldErrors("event_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	app.EventID = eventID.String()
	constituentID := profile.ConstituentID.String()
	app.ConstituentID = &constituentID

	registrationInput, err := toBusNewEventRegistration(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	registration, err := a.admissionsBus.RegisterForEvent(ctx, registrationInput)
	if err != nil {
		return errs.Errorf(errs.Internal, "register for event: %s", err)
	}

	return toAppEventRegistration(registration)
}

func (a *app) registerForEvent(ctx context.Context, r *http.Request) web.Encoder {
	var app NewEventRegistration
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	eventID, err := uuid.Parse(web.Param(r, "event_id"))
	if err != nil {
		return errs.NewFieldErrors("event_id", err)
	}

	app.EventID = eventID.String()

	registrationInput, err := toBusNewEventRegistration(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	registration, err := a.admissionsBus.RegisterForEvent(ctx, registrationInput)
	if err != nil {
		return errs.Errorf(errs.Internal, "register for event: %s", err)
	}

	return toAppEventRegistration(registration)
}

func (a *app) checkInEventRegistration(ctx context.Context, r *http.Request) web.Encoder {
	var app NewEventCheckIn
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	registrationID, err := uuid.Parse(web.Param(r, "event_registration_id"))
	if err != nil {
		return errs.NewFieldErrors("event_registration_id", err)
	}

	registration, err := a.admissionsBus.QueryEventRegistrationByID(ctx, registrationID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query event registration: %s", err)
	}

	checkIn, err := toBusNewEventCheckIn(NewEventCheckIn{RegistrationID: registrationID.String()}, mid.GetSubjectID(ctx))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.CheckInEventRegistration(ctx, registration, checkIn)
	if err != nil {
		return errs.Errorf(errs.Internal, "check in event registration: %s", err)
	}

	return toAppEventRegistration(updated)
}

func (a *app) queryApplicantEvents(ctx context.Context, r *http.Request) web.Encoder {
	return a.queryEvents(ctx, r)
}

func (a *app) queryApplicantEventByID(ctx context.Context, r *http.Request) web.Encoder {
	return a.queryEventByID(ctx, r)
}

func (a *app) createStaffProfile(ctx context.Context, r *http.Request) web.Encoder {
	var app NewStaffProfile
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	np, err := toBusNewStaffProfile(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profile, err := a.admissionsBus.CreateStaffProfile(ctx, np)
	if err != nil {
		return errs.Errorf(errs.Internal, "create staff profile: %s", err)
	}

	return toAppStaffProfile(profile)
}

func (a *app) updateStaffProfile(ctx context.Context, r *http.Request) web.Encoder {
	var app NewStaffProfile
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profileID, err := uuid.Parse(web.Param(r, "staff_profile_id"))
	if err != nil {
		return errs.NewFieldErrors("staff_profile_id", err)
	}

	profile, err := a.admissionsBus.QueryStaffProfileByID(ctx, profileID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query staff profile: %s", err)
	}

	np, err := toBusNewStaffProfile(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.UpdateStaffProfile(ctx, profile, np)
	if err != nil {
		return errs.Errorf(errs.Internal, "update staff profile: %s", err)
	}

	return toAppStaffProfile(updated)
}

func (a *app) queryStaffProfiles(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseStaffProfileQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseStaffProfileFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(staffProfileOrderByFields, qp.OrderBy, admissionsbus.DefaultStaffProfileOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	profiles, err := a.admissionsBus.QueryStaffProfiles(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query staff profiles: %s", err)
	}

	total, err := a.admissionsBus.CountStaffProfiles(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count staff profiles: %s", err)
	}

	return query.NewResult(toAppStaffProfiles(profiles), total, page)
}

func (a *app) queryStaffProfileByID(ctx context.Context, r *http.Request) web.Encoder {
	profileID, err := uuid.Parse(web.Param(r, "staff_profile_id"))
	if err != nil {
		return errs.NewFieldErrors("staff_profile_id", err)
	}

	profile, err := a.admissionsBus.QueryStaffProfileByID(ctx, profileID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query staff profile: %s", err)
	}

	return toAppStaffProfile(profile)
}

func (a *app) createApplicantProfile(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicantProfile
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	np, err := toBusNewApplicantProfile(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profile, err := a.admissionsBus.CreateApplicantProfile(ctx, np)
	if err != nil {
		return errs.Errorf(errs.Internal, "create applicant profile: %s", err)
	}

	return toAppApplicantProfile(profile)
}

func (a *app) updateApplicantProfile(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicantProfile
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profileID, err := uuid.Parse(web.Param(r, "applicant_profile_id"))
	if err != nil {
		return errs.NewFieldErrors("applicant_profile_id", err)
	}

	profile, err := a.admissionsBus.QueryApplicantProfileByID(ctx, profileID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant profile: %s", err)
	}

	np, err := toBusNewApplicantProfile(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.UpdateApplicantProfile(ctx, profile, np)
	if err != nil {
		return errs.Errorf(errs.Internal, "update applicant profile: %s", err)
	}

	return toAppApplicantProfile(updated)
}

func (a *app) queryApplicantProfiles(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicantProfileQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicantProfileFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(applicantProfileOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicantProfileOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	profiles, err := a.admissionsBus.QueryApplicantProfiles(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant profiles: %s", err)
	}

	total, err := a.admissionsBus.CountApplicantProfiles(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant profiles: %s", err)
	}

	return query.NewResult(toAppApplicantProfiles(profiles), total, page)
}

func (a *app) queryApplicantProfileByID(ctx context.Context, r *http.Request) web.Encoder {
	profileID, err := uuid.Parse(web.Param(r, "applicant_profile_id"))
	if err != nil {
		return errs.NewFieldErrors("applicant_profile_id", err)
	}

	profile, err := a.admissionsBus.QueryApplicantProfileByID(ctx, profileID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant profile: %s", err)
	}

	return toAppApplicantProfile(profile)
}

func (a *app) queryCurrentApplicantProfile(ctx context.Context, _ *http.Request) web.Encoder {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	profile, err := a.admissionsBus.QueryApplicantProfileByUserID(ctx, userID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant profile: %s", err)
	}

	return toAppApplicantProfile(profile)
}

func (a *app) createLeadScoreRule(ctx context.Context, r *http.Request) web.Encoder {
	var app NewLeadScoreRule
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rule, err := a.admissionsBus.CreateLeadScoreRule(ctx, toBusNewLeadScoreRule(app))
	if err != nil {
		return errs.Errorf(errs.Internal, "create lead score rule: %s", err)
	}

	return toAppLeadScoreRule(rule)
}

func (a *app) updateLeadScoreRule(ctx context.Context, r *http.Request) web.Encoder {
	var app NewLeadScoreRule
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ruleID, err := uuid.Parse(web.Param(r, "lead_score_rule_id"))
	if err != nil {
		return errs.NewFieldErrors("lead_score_rule_id", err)
	}

	rule, err := a.admissionsBus.QueryLeadScoreRuleByID(ctx, ruleID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead score rule: %s", err)
	}

	updated, err := a.admissionsBus.UpdateLeadScoreRule(ctx, rule, toBusNewLeadScoreRule(app))
	if err != nil {
		return errs.Errorf(errs.Internal, "update lead score rule: %s", err)
	}

	return toAppLeadScoreRule(updated)
}

func (a *app) queryLeadScoreRules(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseLeadScoreRuleQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseLeadScoreRuleFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(leadScoreRuleOrderByFields, qp.OrderBy, admissionsbus.DefaultLeadScoreRuleOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	rules, err := a.admissionsBus.QueryLeadScoreRules(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead score rules: %s", err)
	}

	total, err := a.admissionsBus.CountLeadScoreRules(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count lead score rules: %s", err)
	}

	return query.NewResult(toAppLeadScoreRules(rules), total, page)
}

func (a *app) queryLeadScoreRuleByID(ctx context.Context, r *http.Request) web.Encoder {
	ruleID, err := uuid.Parse(web.Param(r, "lead_score_rule_id"))
	if err != nil {
		return errs.NewFieldErrors("lead_score_rule_id", err)
	}

	rule, err := a.admissionsBus.QueryLeadScoreRuleByID(ctx, ruleID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead score rule: %s", err)
	}

	return toAppLeadScoreRule(rule)
}

func (a *app) queryLeadScores(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseLeadScoreQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseLeadScoreFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(leadScoreOrderByFields, qp.OrderBy, admissionsbus.DefaultLeadScoreOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	scores, err := a.admissionsBus.QueryLeadScores(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead scores: %s", err)
	}

	total, err := a.admissionsBus.CountLeadScores(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count lead scores: %s", err)
	}

	return query.NewResult(toAppLeadScores(scores), total, page)
}

func (a *app) queryLeadScoreByID(ctx context.Context, r *http.Request) web.Encoder {
	scoreID, err := uuid.Parse(web.Param(r, "lead_score_id"))
	if err != nil {
		return errs.NewFieldErrors("lead_score_id", err)
	}

	score, err := a.admissionsBus.QueryLeadScoreByID(ctx, scoreID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead score: %s", err)
	}

	return toAppLeadScore(score)
}

func (a *app) queryLeadScoreByConstituentID(ctx context.Context, r *http.Request) web.Encoder {
	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	score, err := a.admissionsBus.QueryLeadScoreByConstituentID(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query lead score: %s", err)
	}

	return toAppLeadScore(score)
}

func (a *app) recalculateLeadScoreForConstituent(ctx context.Context, r *http.Request) web.Encoder {
	txApp, err := a.newWithTx(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "new transaction app: %s", err)
	}

	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	score, err := txApp.admissionsBus.RecalculateLeadScoreForConstituent(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "recalculate lead score: %s", err)
	}

	return toAppLeadScore(score)
}

func (a *app) updateConstituent(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateConstituent
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	cst, err := a.admissionsBus.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituent: %s", err)
	}

	uc, err := toBusUpdateConstituent(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.UpdateConstituent(ctx, cst, uc)
	if err != nil {
		return errs.Errorf(errs.Internal, "update constituent: %s", err)
	}

	return toAppConstituent(updated)
}

func (a *app) queryConstituents(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseConstituentQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseConstituentFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(constituentOrderByFields, qp.OrderBy, admissionsbus.DefaultConstituentOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	constituents, err := a.admissionsBus.QueryConstituents(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituents: %s", err)
	}

	total, err := a.admissionsBus.CountConstituents(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count constituents: %s", err)
	}

	return query.NewResult(toAppConstituents(constituents), total, page)
}

func (a *app) queryConstituentByID(ctx context.Context, r *http.Request) web.Encoder {
	constituentID, err := uuid.Parse(web.Param(r, "constituent_id"))
	if err != nil {
		return errs.NewFieldErrors("constituent_id", err)
	}

	cst, err := a.admissionsBus.QueryConstituentByID(ctx, constituentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query constituent: %s", err)
	}

	return toAppConstituent(cst)
}

func (a *app) queryPrograms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseProgramQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseProgramFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(programOrderByFields, qp.OrderBy, admissionsbus.DefaultProgramOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	programs, err := a.admissionsBus.QueryPrograms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query programs: %s", err)
	}

	total, err := a.admissionsBus.CountPrograms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count programs: %s", err)
	}

	return query.NewResult(toAppPrograms(programs), total, page)
}

func (a *app) queryProgramByID(ctx context.Context, r *http.Request) web.Encoder {
	programID, err := uuid.Parse(web.Param(r, "program_id"))
	if err != nil {
		return errs.NewFieldErrors("program_id", err)
	}

	program, err := a.admissionsBus.QueryProgramByID(ctx, programID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query program: %s", err)
	}

	return toAppProgram(program)
}

func (a *app) queryApplicantPrograms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseProgramQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseProgramFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}
	active := true
	filter.Active = &active

	orderBy, err := order.Parse(programOrderByFields, qp.OrderBy, admissionsbus.DefaultProgramOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	programs, err := a.admissionsBus.QueryPrograms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant programs: %s", err)
	}

	total, err := a.admissionsBus.CountPrograms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant programs: %s", err)
	}

	return query.NewResult(toAppPrograms(programs), total, page)
}

func (a *app) queryApplicantProgramByID(ctx context.Context, r *http.Request) web.Encoder {
	programID, err := uuid.Parse(web.Param(r, "program_id"))
	if err != nil {
		return errs.NewFieldErrors("program_id", err)
	}

	program, err := a.admissionsBus.QueryProgramByID(ctx, programID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant program: %s", err)
	}
	if !program.Active {
		return errs.New(errs.NotFound, admissionsbus.ErrProgramNotFound)
	}

	return toAppProgram(program)
}

func (a *app) queryAcademicTerms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseAcademicTermQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseAcademicTermFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(academicTermOrderByFields, qp.OrderBy, admissionsbus.DefaultAcademicTermOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	terms, err := a.admissionsBus.QueryAcademicTerms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query academic terms: %s", err)
	}

	total, err := a.admissionsBus.CountAcademicTerms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count academic terms: %s", err)
	}

	return query.NewResult(toAppAcademicTerms(terms), total, page)
}

func (a *app) queryAcademicTermByID(ctx context.Context, r *http.Request) web.Encoder {
	termID, err := uuid.Parse(web.Param(r, "academic_term_id"))
	if err != nil {
		return errs.NewFieldErrors("academic_term_id", err)
	}

	term, err := a.admissionsBus.QueryAcademicTermByID(ctx, termID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query academic term: %s", err)
	}

	return toAppAcademicTerm(term)
}

func (a *app) queryApplicantAcademicTerms(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseAcademicTermQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseAcademicTermFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}
	active := true
	filter.Active = &active

	orderBy, err := order.Parse(academicTermOrderByFields, qp.OrderBy, admissionsbus.DefaultAcademicTermOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	terms, err := a.admissionsBus.QueryAcademicTerms(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant academic terms: %s", err)
	}

	total, err := a.admissionsBus.CountAcademicTerms(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant academic terms: %s", err)
	}

	return query.NewResult(toAppAcademicTerms(terms), total, page)
}

func (a *app) queryApplicantAcademicTermByID(ctx context.Context, r *http.Request) web.Encoder {
	termID, err := uuid.Parse(web.Param(r, "academic_term_id"))
	if err != nil {
		return errs.NewFieldErrors("academic_term_id", err)
	}

	term, err := a.admissionsBus.QueryAcademicTermByID(ctx, termID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant academic term: %s", err)
	}
	if !term.Active {
		return errs.New(errs.NotFound, admissionsbus.ErrAcademicTermNotFound)
	}

	return toAppAcademicTerm(term)
}

func (a *app) createDuplicateReview(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDuplicateReview
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nr, err := toBusNewDuplicateReview(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	review, err := a.admissionsBus.CreateDuplicateReview(ctx, nr)
	if err != nil {
		return errs.Errorf(errs.Internal, "create duplicate review: %s", err)
	}

	return toAppDuplicateReview(review)
}

func (a *app) resolveDuplicateReview(ctx context.Context, r *http.Request) web.Encoder {
	var app ResolveDuplicateReview
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	reviewID, err := uuid.Parse(web.Param(r, "duplicate_review_id"))
	if err != nil {
		return errs.NewFieldErrors("duplicate_review_id", err)
	}

	review, err := a.admissionsBus.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate review: %s", err)
	}

	rr := toBusResolveDuplicateReview(app, mid.GetSubjectID(ctx))
	resolved, err := a.admissionsBus.ResolveDuplicateReview(ctx, review, rr)
	if err != nil {
		return errs.Errorf(errs.Internal, "resolve duplicate review: %s", err)
	}

	if a.auditBus != nil {
		na := auditbus.NewAudit{
			ObjID:     resolved.ID,
			ObjDomain: domain.Admissions,
			ObjName:   name.MustParse("Duplicate Review"),
			ActorID:   rr.ActorID,
			Action:    "duplicate_" + rr.Resolution.String(),
			Data:      toAppDuplicateReview(resolved),
			Message:   "duplicate review resolved",
		}

		if _, err := a.auditBus.Create(ctx, na); err != nil {
			return errs.Errorf(errs.Internal, "audit duplicate review: %s", err)
		}
	}

	return toAppDuplicateReview(resolved)
}

func (a *app) queryDuplicateReviews(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseDuplicateReviewQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseDuplicateReviewFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(duplicateReviewOrderByFields, qp.OrderBy, admissionsbus.DefaultDuplicateReviewOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	reviews, err := a.admissionsBus.QueryDuplicateReviews(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate reviews: %s", err)
	}

	total, err := a.admissionsBus.CountDuplicateReviews(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count duplicate reviews: %s", err)
	}

	return query.NewResult(toAppDuplicateReviews(reviews), total, page)
}

func (a *app) queryDuplicateReviewByID(ctx context.Context, r *http.Request) web.Encoder {
	reviewID, err := uuid.Parse(web.Param(r, "duplicate_review_id"))
	if err != nil {
		return errs.NewFieldErrors("duplicate_review_id", err)
	}

	review, err := a.admissionsBus.QueryDuplicateReviewByID(ctx, reviewID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query duplicate review: %s", err)
	}

	return toAppDuplicateReview(review)
}

func (a *app) createApplication(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplication
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	na, err := toBusNewApplication(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	application, err := a.admissionsBus.CreateApplication(ctx, na)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrDuplicateApplication) {
			return errs.New(errs.Aborted, admissionsbus.ErrDuplicateApplication)
		}
		return errs.Errorf(errs.Internal, "create application: %s", err)
	}

	return toAppApplication(application)
}

func (a *app) createApplicantApplication(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplication
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	na, err := toBusNewApplication(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	if na.ConstituentID != profile.ConstituentID {
		return errs.New(errs.PermissionDenied, admissionsbus.ErrApplicationNotFound)
	}
	na.AssignedReviewerID = nil

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	application, err := a.admissionsBus.CreateApplication(ctx, na)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrDuplicateApplication) {
			return errs.New(errs.Aborted, admissionsbus.ErrDuplicateApplication)
		}
		return errs.Errorf(errs.Internal, "create applicant application: %s", err)
	}

	return toAppApplication(application)
}

func (a *app) updateApplicantApplication(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplication
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	application, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID)
	if appErr != nil {
		return appErr
	}

	na, err := toBusNewApplication(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	if na.ConstituentID != profile.ConstituentID {
		return errs.New(errs.PermissionDenied, admissionsbus.ErrApplicationNotFound)
	}
	na.AssignedReviewerID = nil

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	updated, err := a.admissionsBus.UpdateApplicationDraft(ctx, application, na)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrDuplicateApplication) {
			return errs.New(errs.Aborted, admissionsbus.ErrDuplicateApplication)
		}
		if errors.Is(err, admissionsbus.ErrApplicationNotDraft) {
			return errs.New(errs.FailedPrecondition, admissionsbus.ErrApplicationNotDraft)
		}
		return errs.Errorf(errs.Internal, "update applicant application: %s", err)
	}

	return toAppApplication(updated)
}

func (a *app) queryApplications(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(applicationOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	applications, err := a.admissionsBus.QueryApplications(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applications: %s", err)
	}

	total, err := a.admissionsBus.CountApplications(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applications: %s", err)
	}

	return query.NewResult(toAppApplications(applications), total, page)
}

func (a *app) queryApplicationByID(ctx context.Context, r *http.Request) web.Encoder {
	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	application, err := a.admissionsBus.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application: %s", err)
	}

	return toAppApplication(application)
}

func (a *app) queryApplicantApplications(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}
	filter.ConstituentID = &profile.ConstituentID

	orderBy, err := order.Parse(applicationOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	applications, err := a.admissionsBus.QueryApplications(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant applications: %s", err)
	}

	total, err := a.admissionsBus.CountApplications(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant applications: %s", err)
	}

	return query.NewResult(toAppApplications(applications), total, page)
}

func (a *app) queryApplicantApplicationByID(ctx context.Context, r *http.Request) web.Encoder {
	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	application, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID)
	if appErr != nil {
		return appErr
	}

	return toAppApplication(application)
}

func (a *app) createApplicationFormTemplate(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicationFormTemplate
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nt, err := toBusNewApplicationFormTemplate(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	template, err := a.admissionsBus.CreateApplicationFormTemplate(ctx, nt)
	if err != nil {
		return errs.Errorf(errs.Internal, "create application form template: %s", err)
	}

	return toAppApplicationFormTemplate(template)
}

func (a *app) updateApplicationFormTemplate(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicationFormTemplate
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	templateID, err := uuid.Parse(web.Param(r, "form_template_id"))
	if err != nil {
		return errs.NewFieldErrors("form_template_id", err)
	}

	template, err := a.admissionsBus.QueryApplicationFormTemplateByID(ctx, templateID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application form template: %s", err)
	}

	nt, err := toBusNewApplicationFormTemplate(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	updated, err := a.admissionsBus.UpdateApplicationFormTemplate(ctx, template, nt)
	if err != nil {
		return errs.Errorf(errs.Internal, "update application form template: %s", err)
	}

	return toAppApplicationFormTemplate(updated)
}

func (a *app) queryApplicationFormTemplates(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationFormTemplateQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationFormTemplateFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(applicationFormTemplateOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationFormTemplateOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	templates, err := a.admissionsBus.QueryApplicationFormTemplates(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application form templates: %s", err)
	}

	total, err := a.admissionsBus.CountApplicationFormTemplates(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count application form templates: %s", err)
	}

	return query.NewResult(toAppApplicationFormTemplates(templates), total, page)
}

func (a *app) queryApplicationFormTemplateByID(ctx context.Context, r *http.Request) web.Encoder {
	templateID, err := uuid.Parse(web.Param(r, "form_template_id"))
	if err != nil {
		return errs.NewFieldErrors("form_template_id", err)
	}

	template, err := a.admissionsBus.QueryApplicationFormTemplateByID(ctx, templateID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application form template: %s", err)
	}

	return toAppApplicationFormTemplate(template)
}

func (a *app) queryApplicantApplicationFormTemplates(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationFormTemplateQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationFormTemplateFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}
	active := true
	filter.Active = &active

	orderBy, err := order.Parse(applicationFormTemplateOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationFormTemplateOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	templates, err := a.admissionsBus.QueryApplicationFormTemplates(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant application form templates: %s", err)
	}

	total, err := a.admissionsBus.CountApplicationFormTemplates(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant application form templates: %s", err)
	}

	return query.NewResult(toAppApplicationFormTemplates(templates), total, page)
}

func (a *app) queryApplicantApplicationFormTemplateByID(ctx context.Context, r *http.Request) web.Encoder {
	templateID, err := uuid.Parse(web.Param(r, "form_template_id"))
	if err != nil {
		return errs.NewFieldErrors("form_template_id", err)
	}

	template, err := a.admissionsBus.QueryApplicationFormTemplateByID(ctx, templateID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant application form template: %s", err)
	}
	if !template.Active {
		return errs.New(errs.NotFound, admissionsbus.ErrFormTemplateNotFound)
	}

	return toAppApplicationFormTemplate(template)
}

func (a *app) createCustomFieldDefinition(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCustomFieldDefinition
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	definition, err := a.admissionsBus.CreateCustomFieldDefinition(ctx, toBusNewCustomFieldDefinition(app))
	if err != nil {
		return errs.Errorf(errs.Internal, "create custom field definition: %s", err)
	}

	return toAppCustomFieldDefinition(definition)
}

func (a *app) updateCustomFieldDefinition(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCustomFieldDefinition
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	definitionID, err := uuid.Parse(web.Param(r, "custom_field_definition_id"))
	if err != nil {
		return errs.NewFieldErrors("custom_field_definition_id", err)
	}

	definition, err := a.admissionsBus.QueryCustomFieldDefinitionByID(ctx, definitionID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query custom field definition: %s", err)
	}

	updated, err := a.admissionsBus.UpdateCustomFieldDefinition(ctx, definition, toBusNewCustomFieldDefinition(app))
	if err != nil {
		return errs.Errorf(errs.Internal, "update custom field definition: %s", err)
	}

	return toAppCustomFieldDefinition(updated)
}

func (a *app) queryCustomFieldDefinitions(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseCustomFieldDefinitionQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseCustomFieldDefinitionFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(customFieldDefinitionOrderByFields, qp.OrderBy, admissionsbus.DefaultCustomFieldDefinitionOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	definitions, err := a.admissionsBus.QueryCustomFieldDefinitions(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query custom field definitions: %s", err)
	}

	total, err := a.admissionsBus.CountCustomFieldDefinitions(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count custom field definitions: %s", err)
	}

	return query.NewResult(toAppCustomFieldDefinitions(definitions), total, page)
}

func (a *app) queryCustomFieldDefinitionByID(ctx context.Context, r *http.Request) web.Encoder {
	definitionID, err := uuid.Parse(web.Param(r, "custom_field_definition_id"))
	if err != nil {
		return errs.NewFieldErrors("custom_field_definition_id", err)
	}

	definition, err := a.admissionsBus.QueryCustomFieldDefinitionByID(ctx, definitionID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query custom field definition: %s", err)
	}

	return toAppCustomFieldDefinition(definition)
}

func (a *app) setCustomFieldValue(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCustomFieldValue
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nv, err := toBusNewCustomFieldValue(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	value, err := a.admissionsBus.SetCustomFieldValue(ctx, nv)
	if err != nil {
		return errs.Errorf(errs.Internal, "set custom field value: %s", err)
	}

	return toAppCustomFieldValue(value)
}

func (a *app) setApplicantCustomFieldValue(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCustomFieldValue
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nv, err := toBusNewCustomFieldValue(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	if appErr := a.ensureApplicantOwnsCustomFieldOwner(ctx, nv.Owner, nv.OwnerID, profile.ConstituentID); appErr != nil {
		return appErr
	}

	value, err := a.admissionsBus.SetCustomFieldValue(ctx, nv)
	if err != nil {
		return errs.Errorf(errs.Internal, "set applicant custom field value: %s", err)
	}

	return toAppCustomFieldValue(value)
}

func (a *app) queryCustomFieldValues(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseCustomFieldValueQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseCustomFieldValueFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(customFieldValueOrderByFields, qp.OrderBy, admissionsbus.DefaultCustomFieldValueOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	values, err := a.admissionsBus.QueryCustomFieldValues(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query custom field values: %s", err)
	}

	total, err := a.admissionsBus.CountCustomFieldValues(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count custom field values: %s", err)
	}

	return query.NewResult(toAppCustomFieldValues(values), total, page)
}

func (a *app) queryCustomFieldValueByID(ctx context.Context, r *http.Request) web.Encoder {
	valueID, err := uuid.Parse(web.Param(r, "custom_field_value_id"))
	if err != nil {
		return errs.NewFieldErrors("custom_field_value_id", err)
	}

	value, err := a.admissionsBus.QueryCustomFieldValueByID(ctx, valueID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query custom field value: %s", err)
	}

	return toAppCustomFieldValue(value)
}

func (a *app) transitionApplicationStatus(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicationTransition
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	application, err := a.admissionsBus.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application: %s", err)
	}

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nt := toBusNewApplicationTransition(app, mid.GetSubjectID(ctx))
	updated, transition, err := a.admissionsBus.TransitionApplicationStatus(ctx, application, nt)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrInvalidApplicationTransition) {
			return errs.New(errs.FailedPrecondition, admissionsbus.ErrInvalidApplicationTransition)
		}
		return errs.Errorf(errs.Internal, "transition application: %s", err)
	}

	if a.auditBus != nil {
		na := auditbus.NewAudit{
			ObjID:     updated.ID,
			ObjDomain: domain.Admissions,
			ObjName:   name.MustParse("Application"),
			ActorID:   nt.ActorID,
			Action:    "application_transition",
			Data:      toAppApplicationTransition(transition),
			Message:   fmt.Sprintf("application status changed from %s to %s", transition.FromStatus, transition.ToStatus),
		}

		if _, err := a.auditBus.Create(ctx, na); err != nil {
			return errs.Errorf(errs.Internal, "audit application transition: %s", err)
		}
	}

	return toAppApplication(updated)
}

func (a *app) transitionApplicantApplicationStatus(ctx context.Context, r *http.Request) web.Encoder {
	var app NewApplicationTransition
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}

	application, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID)
	if appErr != nil {
		return appErr
	}

	nt := toBusNewApplicationTransition(app, mid.GetSubjectID(ctx))
	if !isApplicantApplicationTransition(nt.ToStatus) {
		return errs.New(errs.PermissionDenied, admissionsbus.ErrInvalidApplicationTransition)
	}

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	updated, transition, err := a.admissionsBus.TransitionApplicationStatus(ctx, application, nt)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrInvalidApplicationTransition) {
			return errs.New(errs.FailedPrecondition, admissionsbus.ErrInvalidApplicationTransition)
		}
		return errs.Errorf(errs.Internal, "transition applicant application: %s", err)
	}

	if a.auditBus != nil {
		na := auditbus.NewAudit{
			ObjID:     updated.ID,
			ObjDomain: domain.Admissions,
			ObjName:   name.MustParse("Application"),
			ActorID:   nt.ActorID,
			Action:    "applicant_application_transition",
			Data:      toAppApplicationTransition(transition),
			Message:   fmt.Sprintf("applicant changed application status from %s to %s", transition.FromStatus, transition.ToStatus),
		}

		if _, err := a.auditBus.Create(ctx, na); err != nil {
			return errs.Errorf(errs.Internal, "audit applicant application transition: %s", err)
		}
	}

	return toAppApplication(updated)
}

func (a *app) queryApplicationTransitions(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseApplicationTransitionQueryParams(r)
	if qp.ApplicationID == "" {
		qp.ApplicationID = web.Param(r, "application_id")
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseApplicationTransitionFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(applicationTransitionOrderByFields, qp.OrderBy, admissionsbus.DefaultApplicationTransitionOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	transitions, err := a.admissionsBus.QueryApplicationTransitions(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query application transitions: %s", err)
	}

	total, err := a.admissionsBus.CountApplicationTransitions(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count application transitions: %s", err)
	}

	return query.NewResult(toAppApplicationTransitions(transitions), total, page)
}

func (a *app) createChecklistItem(ctx context.Context, r *http.Request) web.Encoder {
	var app NewChecklistItem
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	ni := toBusNewChecklistItem(app, applicationID)
	item, err := a.admissionsBus.CreateChecklistItem(ctx, ni)
	if err != nil {
		return errs.Errorf(errs.Internal, "create checklist item: %s", err)
	}

	return toAppChecklistItem(item)
}

func (a *app) updateChecklistItem(ctx context.Context, r *http.Request) web.Encoder {
	var app NewChecklistItem
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	itemID, err := uuid.Parse(web.Param(r, "checklist_item_id"))
	if err != nil {
		return errs.NewFieldErrors("checklist_item_id", err)
	}

	item, err := a.admissionsBus.QueryChecklistItemByID(ctx, itemID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query checklist item: %s", err)
	}

	ni := toBusNewChecklistItem(app, applicationID)
	updated, err := a.admissionsBus.UpdateChecklistItem(ctx, item, ni)
	if err != nil {
		return errs.Errorf(errs.Internal, "update checklist item: %s", err)
	}

	return toAppChecklistItem(updated)
}

func (a *app) queryChecklistItems(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseChecklistItemQueryParams(r)
	if qp.ApplicationID == "" {
		qp.ApplicationID = web.Param(r, "application_id")
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseChecklistItemFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(checklistItemOrderByFields, qp.OrderBy, admissionsbus.DefaultChecklistItemOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	items, err := a.admissionsBus.QueryChecklistItems(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query checklist items: %s", err)
	}

	total, err := a.admissionsBus.CountChecklistItems(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count checklist items: %s", err)
	}

	return query.NewResult(toAppChecklistItems(items), total, page)
}

func (a *app) queryApplicantChecklistItems(ctx context.Context, r *http.Request) web.Encoder {
	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}
	if _, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID); appErr != nil {
		return appErr
	}

	qp := parseChecklistItemQueryParams(r)
	qp.ApplicationID = applicationID.String()

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseChecklistItemFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(checklistItemOrderByFields, qp.OrderBy, admissionsbus.DefaultChecklistItemOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	items, err := a.admissionsBus.QueryChecklistItems(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant checklist items: %s", err)
	}

	total, err := a.admissionsBus.CountChecklistItems(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant checklist items: %s", err)
	}

	return query.NewResult(toAppChecklistItems(items), total, page)
}

func (a *app) queryChecklistItemByID(ctx context.Context, r *http.Request) web.Encoder {
	itemID, err := uuid.Parse(web.Param(r, "checklist_item_id"))
	if err != nil {
		return errs.NewFieldErrors("checklist_item_id", err)
	}

	item, err := a.admissionsBus.QueryChecklistItemByID(ctx, itemID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query checklist item: %s", err)
	}

	return toAppChecklistItem(item)
}

func (a *app) createDocument(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDocument
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nd, err := toBusNewDocument(app, applicationID, mid.GetSubjectID(ctx))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	document, err := a.admissionsBus.CreateDocument(ctx, nd)
	if err != nil {
		return errs.Errorf(errs.Internal, "create document: %s", err)
	}

	if err := a.auditDocument(ctx, document, nd.UploadedByID, "document_upload", "admissions document uploaded"); err != nil {
		return err
	}

	return toAppDocument(document)
}

func (a *app) createApplicantDocument(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDocument
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}
	if _, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID); appErr != nil {
		return appErr
	}

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nd, err := toBusNewDocument(app, applicationID, mid.GetSubjectID(ctx))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	document, err := a.admissionsBus.CreateDocument(ctx, nd)
	if err != nil {
		return errs.Errorf(errs.Internal, "create applicant document: %s", err)
	}

	if err := a.auditDocument(ctx, document, nd.UploadedByID, "applicant_document_upload", "applicant admissions document uploaded"); err != nil {
		return err
	}

	return toAppDocument(document)
}

func (a *app) verifyDocument(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDocumentVerification
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	documentID, err := uuid.Parse(web.Param(r, "document_id"))
	if err != nil {
		return errs.NewFieldErrors("document_id", err)
	}

	document, err := a.admissionsBus.QueryDocumentByID(ctx, documentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query document: %s", err)
	}

	a, err = a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nv := toBusNewDocumentVerification(app, mid.GetSubjectID(ctx))
	updated, err := a.admissionsBus.VerifyDocument(ctx, document, nv)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrDocumentStatusNotReviewable) {
			return errs.New(errs.FailedPrecondition, admissionsbus.ErrDocumentStatusNotReviewable)
		}
		return errs.Errorf(errs.Internal, "verify document: %s", err)
	}

	if err := a.auditDocument(ctx, updated, nv.ReviewerID, "document_"+nv.Status.String(), "admissions document verification recorded"); err != nil {
		return err
	}

	return toAppDocument(updated)
}

func (a *app) queryDocuments(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseDocumentQueryParams(r)
	if qp.ApplicationID == "" {
		qp.ApplicationID = web.Param(r, "application_id")
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseDocumentFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(documentOrderByFields, qp.OrderBy, admissionsbus.DefaultDocumentOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	documents, err := a.admissionsBus.QueryDocuments(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query documents: %s", err)
	}

	total, err := a.admissionsBus.CountDocuments(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count documents: %s", err)
	}

	return query.NewResult(toAppDocuments(documents), total, page)
}

func (a *app) queryApplicantDocuments(ctx context.Context, r *http.Request) web.Encoder {
	applicationID, err := uuid.Parse(web.Param(r, "application_id"))
	if err != nil {
		return errs.NewFieldErrors("application_id", err)
	}

	profile, appErr := a.currentApplicantProfile(ctx)
	if appErr != nil {
		return appErr
	}
	if _, appErr := a.ownedApplication(ctx, applicationID, profile.ConstituentID); appErr != nil {
		return appErr
	}

	qp := parseDocumentQueryParams(r)
	qp.ApplicationID = applicationID.String()

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseDocumentFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(documentOrderByFields, qp.OrderBy, admissionsbus.DefaultDocumentOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	documents, err := a.admissionsBus.QueryDocuments(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query applicant documents: %s", err)
	}

	total, err := a.admissionsBus.CountDocuments(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count applicant documents: %s", err)
	}

	return query.NewResult(toAppDocuments(documents), total, page)
}

func (a *app) queryDocumentByID(ctx context.Context, r *http.Request) web.Encoder {
	documentID, err := uuid.Parse(web.Param(r, "document_id"))
	if err != nil {
		return errs.NewFieldErrors("document_id", err)
	}

	document, err := a.admissionsBus.QueryDocumentByID(ctx, documentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query document: %s", err)
	}

	if err := a.auditDocument(ctx, document, mid.GetSubjectID(ctx), "document_view", "admissions document metadata viewed"); err != nil {
		return err
	}

	return toAppDocument(document)
}

func (a *app) downloadDocument(ctx context.Context, r *http.Request) web.Encoder {
	documentID, err := uuid.Parse(web.Param(r, "document_id"))
	if err != nil {
		return errs.NewFieldErrors("document_id", err)
	}

	document, err := a.admissionsBus.QueryDocumentByID(ctx, documentID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query document: %s", err)
	}

	if err := a.auditDocument(ctx, document, mid.GetSubjectID(ctx), "document_download", "admissions document download requested"); err != nil {
		return err
	}

	return toAppDocument(document)
}

func (a *app) createImportBatch(ctx context.Context, r *http.Request) web.Encoder {
	var app NewImportBatch
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nb := toBusNewImportBatch(app, mid.GetSubjectID(ctx))
	batch, err := a.admissionsBus.CreateImportBatch(ctx, nb)
	if err != nil {
		return errs.Errorf(errs.Internal, "create import batch: %s", err)
	}

	if err := a.auditImportBatch(ctx, batch, nb.UploadedByID, "import_batch_create", "admissions import batch recorded"); err != nil {
		return err
	}

	return toAppImportBatch(batch)
}

func (a *app) queryImportBatches(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseImportBatchQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseImportBatchFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(importBatchOrderByFields, qp.OrderBy, admissionsbus.DefaultImportBatchOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	batches, err := a.admissionsBus.QueryImportBatches(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query import batches: %s", err)
	}

	total, err := a.admissionsBus.CountImportBatches(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count import batches: %s", err)
	}

	return query.NewResult(toAppImportBatches(batches), total, page)
}

func (a *app) queryImportBatchByID(ctx context.Context, r *http.Request) web.Encoder {
	batchID, err := uuid.Parse(web.Param(r, "import_batch_id"))
	if err != nil {
		return errs.NewFieldErrors("import_batch_id", err)
	}

	batch, err := a.admissionsBus.QueryImportBatchByID(ctx, batchID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query import batch: %s", err)
	}

	return toAppImportBatch(batch)
}

func (a *app) createImportInvalidRows(ctx context.Context, r *http.Request) web.Encoder {
	var app NewImportInvalidRows
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	batchID, err := uuid.Parse(web.Param(r, "import_batch_id"))
	if err != nil {
		return errs.NewFieldErrors("import_batch_id", err)
	}

	rows, err := a.admissionsBus.CreateImportInvalidRows(ctx, toBusNewImportInvalidRows(app, batchID))
	if err != nil {
		return errs.Errorf(errs.Internal, "create import invalid rows: %s", err)
	}

	return query.NewResult(toAppImportInvalidRows(rows), len(rows), page.MustParse("1", "100"))
}

func (a *app) queryImportInvalidRows(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseImportInvalidRowQueryParams(r)
	if qp.BatchID == "" {
		qp.BatchID = web.Param(r, "import_batch_id")
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseImportInvalidRowFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(importInvalidRowOrderByFields, qp.OrderBy, admissionsbus.DefaultImportInvalidRowOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	rows, err := a.admissionsBus.QueryImportInvalidRows(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query import invalid rows: %s", err)
	}

	total, err := a.admissionsBus.CountImportInvalidRows(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count import invalid rows: %s", err)
	}

	return query.NewResult(toAppImportInvalidRows(rows), total, page)
}

func (a *app) queryImportInvalidRowByID(ctx context.Context, r *http.Request) web.Encoder {
	rowID, err := uuid.Parse(web.Param(r, "import_invalid_row_id"))
	if err != nil {
		return errs.NewFieldErrors("import_invalid_row_id", err)
	}

	row, err := a.admissionsBus.QueryImportInvalidRowByID(ctx, rowID)
	if err != nil {
		return errs.Errorf(errs.Internal, "query import invalid row: %s", err)
	}

	return toAppImportInvalidRow(row)
}

func (a *app) downloadImportInvalidRows(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseImportInvalidRowQueryParams(r)
	qp.BatchID = web.Param(r, "import_batch_id")

	filter, err := parseImportInvalidRowFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	rows, err := a.admissionsBus.QueryImportInvalidRows(ctx, filter, admissionsbus.DefaultImportInvalidRowOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return errs.Errorf(errs.Internal, "query import invalid rows download: %s", err)
	}

	return query.NewResult(toAppImportInvalidRows(rows), len(rows), page.MustParse("1", "100"))
}

func (a *app) auditDocument(ctx context.Context, document admissionsbus.Document, actorID uuid.UUID, action string, message string) *errs.Error {
	if a.auditBus == nil {
		return nil
	}

	na := auditbus.NewAudit{
		ObjID:     document.ID,
		ObjDomain: domain.Admissions,
		ObjName:   name.MustParse("Document"),
		ActorID:   actorID,
		Action:    action,
		Data:      toAppDocument(document),
		Message:   message,
	}

	if _, err := a.auditBus.Create(ctx, na); err != nil {
		return errs.Errorf(errs.Internal, "audit document: %s", err)
	}

	return nil
}

func (a *app) auditImportBatch(ctx context.Context, batch admissionsbus.ImportBatch, actorID uuid.UUID, action string, message string) *errs.Error {
	if a.auditBus == nil {
		return nil
	}

	na := auditbus.NewAudit{
		ObjID:     batch.ID,
		ObjDomain: domain.Admissions,
		ObjName:   name.MustParse("Import Batch"),
		ActorID:   actorID,
		Action:    action,
		Data:      toAppImportBatch(batch),
		Message:   message,
	}

	if _, err := a.auditBus.Create(ctx, na); err != nil {
		return errs.Errorf(errs.Internal, "audit import batch: %s", err)
	}

	return nil
}

func (a *app) currentApplicantProfile(ctx context.Context) (admissionsbus.ApplicantProfile, *errs.Error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return admissionsbus.ApplicantProfile{}, errs.New(errs.Unauthenticated, err)
	}

	profile, err := a.admissionsBus.QueryApplicantProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrApplicantProfileNotFound) {
			return admissionsbus.ApplicantProfile{}, errs.New(errs.Unauthenticated, admissionsbus.ErrApplicantProfileNotFound)
		}
		return admissionsbus.ApplicantProfile{}, errs.Errorf(errs.Internal, "query applicant profile: %s", err)
	}
	if !profile.Active {
		return admissionsbus.ApplicantProfile{}, errs.New(errs.PermissionDenied, fmt.Errorf("applicant profile is inactive"))
	}

	return profile, nil
}

func (a *app) ownedApplication(ctx context.Context, applicationID uuid.UUID, constituentID uuid.UUID) (admissionsbus.Application, *errs.Error) {
	application, err := a.admissionsBus.QueryApplicationByID(ctx, applicationID)
	if err != nil {
		return admissionsbus.Application{}, errs.Errorf(errs.Internal, "query applicant application: %s", err)
	}
	if application.ConstituentID != constituentID {
		return admissionsbus.Application{}, errs.New(errs.PermissionDenied, admissionsbus.ErrApplicationNotFound)
	}

	return application, nil
}

func (a *app) ensureApplicantOwnsCustomFieldOwner(ctx context.Context, owner admissionsbus.CustomFieldOwner, ownerID uuid.UUID, constituentID uuid.UUID) *errs.Error {
	switch owner {
	case admissionsbus.CustomFieldOwnerConstituent:
		if ownerID != constituentID {
			return errs.New(errs.PermissionDenied, admissionsbus.ErrApplicationNotFound)
		}
		return nil
	case admissionsbus.CustomFieldOwnerApplication:
		_, err := a.ownedApplication(ctx, ownerID, constituentID)
		return err
	default:
		return errs.New(errs.InvalidArgument, fmt.Errorf("invalid custom field owner"))
	}
}

func isApplicantApplicationTransition(status admissionsbus.ApplicationStatus) bool {
	switch status {
	case admissionsbus.ApplicationStatusSubmitted,
		admissionsbus.ApplicationStatusWithdrawn:
		return true
	default:
		return false
	}
}
