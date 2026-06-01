// Package authapp maintains the web based api for auth access.
package authapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/types/role"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	auth          *auth.Auth
	userBus       userbus.ExtBusiness
	admissionsBus admissionsbus.ExtBusiness
	tokenKey      string
}

func newApp(ath *auth.Auth, userBus userbus.ExtBusiness, admissionsBus admissionsbus.ExtBusiness, tokenKey string) *app {
	return &app{
		auth:          ath,
		userBus:       userBus,
		admissionsBus: admissionsBus,
		tokenKey:      tokenKey,
	}
}

func (a *app) login(ctx context.Context, r *http.Request) web.Encoder {
	var req loginReq
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Authenticate(ctx, *addr, req.Password)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	if !usr.Enabled {
		return errs.New(errs.Unauthenticated, auth.ErrUserDisabled)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(auth.TokenExpiry)
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: role.ParseToString(usr.Roles),
	}

	tkn, err := a.auth.GenerateToken(a.tokenKey, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return loginResp{
		AccessToken: tkn,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		ExpiresIn:   int(auth.TokenExpiry.Seconds()),
		User:        toLoginUser(usr),
	}
}

func (a *app) applicantPortalToken(ctx context.Context, r *http.Request) web.Encoder {
	var req applicantPortalTokenReq
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if a.admissionsBus == nil {
		return errs.New(errs.Internal, fmt.Errorf("admissions bus not configured"))
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	constituent, application, err := a.findSubmittedPortalApplication(ctx, addr.Address)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("we could not verify a submitted application for that email"))
	}

	now := time.Now().UTC()
	expiresAt := now.Add(auth.TokenExpiry)
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   application.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: []string{admissionsbus.AdmissionsRoleApplicant.String()},
		Portal: &auth.PortalClaims{
			Scope:         auth.PortalScopeApplicant,
			ApplicationID: application.ID.String(),
			ConstituentID: constituent.ID.String(),
			Email:         addr.Address,
		},
	}

	tkn, err := a.auth.GenerateToken(a.tokenKey, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return applicantPortalTokenResp{
		AccessToken:   tkn,
		TokenType:     "Bearer",
		ExpiresAt:     expiresAt,
		ExpiresIn:     int(auth.TokenExpiry.Seconds()),
		ApplicationID: application.ID.String(),
		ConstituentID: constituent.ID.String(),
		ApplicantName: strings.TrimSpace(constituent.FirstName + " " + constituent.LastName),
		Email:         constituent.PrimaryEmail.Address,
	}
}

func (a *app) findSubmittedPortalApplication(ctx context.Context, email string) (admissionsbus.Constituent, admissionsbus.Application, error) {
	constituent, err := a.admissionsBus.QueryConstituentByPrimaryEmail(ctx, email)
	if err != nil {
		return admissionsbus.Constituent{}, admissionsbus.Application{}, err
	}

	applications, err := a.admissionsBus.QueryApplications(
		ctx,
		admissionsbus.ApplicationQueryFilter{ConstituentID: &constituent.ID},
		admissionsbus.DefaultApplicationOrderBy,
		page.MustParse("1", "100"),
	)
	if err != nil {
		return admissionsbus.Constituent{}, admissionsbus.Application{}, err
	}

	var latest admissionsbus.Application
	for _, application := range applications {
		if application.SubmittedAt == nil || application.Status == admissionsbus.ApplicationStatusDraft {
			continue
		}

		if latest.ID == uuid.Nil || latest.SubmittedAt == nil || application.SubmittedAt.After(*latest.SubmittedAt) {
			latest = application
		}
	}

	if latest.ID == uuid.Nil {
		return admissionsbus.Constituent{}, admissionsbus.Application{}, admissionsbus.ErrApplicationNotFound
	}

	return constituent, latest, nil
}

func (a *app) token(ctx context.Context, r *http.Request) web.Encoder {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errs.NewFieldErrors("kid", errors.New("missing kid"))
	}

	// The BearerBasic middleware function generates the claims.
	claims := mid.GetClaims(ctx)

	tkn, err := a.auth.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return token{Token: tkn}
}

func (a *app) authenticate(ctx context.Context, r *http.Request) web.Encoder {
	// The middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	usr, err := a.userBus.QueryByID(ctx, userID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	resp := authenticateResp{
		UserID: userID.String(),
		Claims: mid.GetClaims(ctx),
		User:   toLoginUser(usr),
	}

	return resp
}

func (a *app) authorize(ctx context.Context, r *http.Request) web.Encoder {
	var auth authclient.Authorize
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Errorf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[%v] rule[%v]", auth.Claims.Roles, auth.Rule)
	}

	return nil
}
