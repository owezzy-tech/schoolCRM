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
	"github.com/owezzy/schoolCRM/business/types/name"
	"github.com/owezzy/schoolCRM/business/types/password"
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

	profile, err := a.admissionsBus.QueryApplicantProfileByConstituentID(ctx, constituent.ID)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}
	if !profile.Active {
		return errs.New(errs.Unauthenticated, fmt.Errorf("admissions applicant profile inactive"))
	}

	now := time.Now().UTC()
	expiresAt := now.Add(auth.TokenExpiry)
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   profile.UserID.String(),
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

	return a.applicantPortalTokenResponse(claims, constituent)
}

func (a *app) newApplicantPortalToken(userID uuid.UUID, constituent admissionsbus.Constituent, applicationID uuid.UUID) web.Encoder {
	now := time.Now().UTC()
	expiresAt := now.Add(auth.TokenExpiry)
	appID := ""
	if applicationID != uuid.Nil {
		appID = applicationID.String()
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: []string{admissionsbus.AdmissionsRoleApplicant.String()},
		Portal: &auth.PortalClaims{
			Scope:         auth.PortalScopeApplicant,
			ApplicationID: appID,
			ConstituentID: constituent.ID.String(),
			Email:         constituent.PrimaryEmail.Address,
		},
	}

	return a.applicantPortalTokenResponse(claims, constituent)
}

func (a *app) applicantPortalTokenResponse(claims auth.Claims, constituent admissionsbus.Constituent) web.Encoder {
	tkn, err := a.auth.GenerateToken(a.tokenKey, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	applicationID := ""
	if claims.Portal != nil {
		applicationID = claims.Portal.ApplicationID
	}

	expiresAt := time.Now().UTC().Add(auth.TokenExpiry)
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	return applicantPortalTokenResp{
		AccessToken:   tkn,
		TokenType:     "Bearer",
		ExpiresAt:     expiresAt,
		ExpiresIn:     int(time.Until(expiresAt).Seconds()),
		ApplicationID: applicationID,
		ConstituentID: constituent.ID.String(),
		ApplicantName: strings.TrimSpace(constituent.FirstName + " " + constituent.LastName),
		Email:         constituent.PrimaryEmail.Address,
	}
}

func (a *app) applicantPortalOnboard(ctx context.Context, r *http.Request) web.Encoder {
	var req applicantPortalOnboardReq
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := req.Validate(); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if a.admissionsBus == nil {
		return errs.New(errs.Internal, fmt.Errorf("admissions bus not configured"))
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	fullName, err := name.Parse(strings.TrimSpace(req.FirstName + " " + req.LastName))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	pw, err := password.ParseConfirm(req.Password, req.ConfirmPassword)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	dob, err := time.Parse(time.RFC3339, req.DateOfBirth)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("dateOfBirth is invalid: %w", err))
	}

	usr, appErr := a.findOrCreateApplicantUser(ctx, *addr, fullName, pw)
	if appErr != nil {
		return appErr
	}

	constituent, appErr := a.findOrCreateApplicantConstituent(ctx, req, *addr, dob)
	if appErr != nil {
		return appErr
	}

	if err := a.ensureApplicantProfile(ctx, usr.ID, constituent.ID); err != nil {
		return err
	}

	return a.newApplicantPortalToken(usr.ID, constituent, uuid.Nil)
}

func (a *app) findOrCreateApplicantUser(ctx context.Context, addr mail.Address, fullName name.Name, pw password.Password) (userbus.User, *errs.Error) {
	usr, err := a.userBus.QueryByEmail(ctx, addr)
	if err == nil {
		return usr, nil
	}
	if !errors.Is(err, userbus.ErrNotFound) {
		return userbus.User{}, errs.Errorf(errs.Internal, "query applicant user: %s", err)
	}

	usr, err = a.userBus.Create(ctx, uuid.UUID{}, userbus.NewUser{
		Name:     fullName,
		Email:    addr,
		Roles:    []role.Role{role.Student},
		Password: pw,
	})
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			usr, queryErr := a.userBus.QueryByEmail(ctx, addr)
			if queryErr == nil {
				return usr, nil
			}

			return userbus.User{}, errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return userbus.User{}, errs.Errorf(errs.Internal, "create applicant user: %s", err)
	}

	return usr, nil
}

func (a *app) findOrCreateApplicantConstituent(ctx context.Context, req applicantPortalOnboardReq, addr mail.Address, dob time.Time) (admissionsbus.Constituent, *errs.Error) {
	constituent, err := a.admissionsBus.QueryConstituentByPrimaryEmail(ctx, addr.Address)
	if err == nil {
		return constituent, nil
	}
	if !errors.Is(err, admissionsbus.ErrConstituentNotFound) {
		return admissionsbus.Constituent{}, errs.Errorf(errs.Internal, "query applicant constituent: %s", err)
	}

	constituent, err = a.admissionsBus.CreateConstituent(ctx, admissionsbus.NewConstituent{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DateOfBirth:     dob,
		PrimaryEmail:    addr,
		PrimaryPhone:    req.Phone,
		LifecycleStage:  admissionsbus.LifecycleStageApplicant,
		DuplicateStatus: admissionsbus.DuplicateStatusActive,
	})
	if err != nil {
		return admissionsbus.Constituent{}, errs.Errorf(errs.Internal, "create applicant constituent: %s", err)
	}

	return constituent, nil
}

func (a *app) ensureApplicantProfile(ctx context.Context, userID uuid.UUID, constituentID uuid.UUID) *errs.Error {
	profile, err := a.admissionsBus.QueryApplicantProfileByUserID(ctx, userID)
	if err == nil {
		if profile.ConstituentID != constituentID {
			return errs.New(errs.Aborted, fmt.Errorf("applicant user is already linked to a different constituent"))
		}
		return nil
	}
	if !errors.Is(err, admissionsbus.ErrApplicantProfileNotFound) {
		return errs.Errorf(errs.Internal, "query applicant profile: %s", err)
	}

	profile, err = a.admissionsBus.QueryApplicantProfileByConstituentID(ctx, constituentID)
	if err == nil {
		if profile.UserID != userID {
			return errs.New(errs.Aborted, fmt.Errorf("applicant constituent is already linked to a different user"))
		}
		return nil
	}
	if !errors.Is(err, admissionsbus.ErrApplicantProfileNotFound) {
		return errs.Errorf(errs.Internal, "query applicant profile by constituent: %s", err)
	}

	if _, err := a.admissionsBus.CreateApplicantProfile(ctx, admissionsbus.NewApplicantProfile{
		UserID:        userID,
		ConstituentID: constituentID,
		Active:        true,
	}); err != nil {
		return errs.Errorf(errs.Internal, "create applicant profile: %s", err)
	}

	return nil
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
	claims := mid.GetClaims(ctx)

	if claims.Portal != nil && claims.Portal.Scope == auth.PortalScopeApplicant {
		return a.authenticateApplicantPortal(ctx, claims)
	}

	usr, err := a.userBus.QueryByID(ctx, userID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	claims, err = a.claimsWithAdmissionsRoles(ctx, userID, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	resp := authenticateResp{
		UserID: userID.String(),
		Claims: claims,
		User:   toLoginUser(usr),
	}

	return resp
}

func (a *app) authenticateApplicantPortal(ctx context.Context, claims auth.Claims) web.Encoder {
	if a.admissionsBus == nil {
		return errs.New(errs.Internal, fmt.Errorf("admissions bus not configured"))
	}
	if claims.Portal == nil || claims.Portal.ConstituentID == "" {
		return errs.New(errs.Unauthenticated, fmt.Errorf("missing applicant portal constituent"))
	}

	constituentID, err := uuid.Parse(claims.Portal.ConstituentID)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("parsing constituent id: %w", err))
	}

	profile, err := a.admissionsBus.QueryApplicantProfileByConstituentID(ctx, constituentID)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}
	if !profile.Active {
		return errs.New(errs.Unauthenticated, fmt.Errorf("admissions applicant profile inactive"))
	}
	if claims.Subject != profile.UserID.String() {
		return errs.New(errs.Unauthenticated, fmt.Errorf("applicant portal subject does not match applicant profile"))
	}

	usr, err := a.userBus.QueryByID(ctx, profile.UserID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return authenticateResp{
		UserID: profile.UserID.String(),
		Claims: claims,
		User:   toLoginUser(usr),
	}
}

func (a *app) claimsWithAdmissionsRoles(ctx context.Context, userID uuid.UUID, claims auth.Claims) (auth.Claims, error) {
	if a.admissionsBus == nil {
		return claims, nil
	}

	profile, err := a.admissionsBus.QueryStaffProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, admissionsbus.ErrStaffProfileNotFound) {
			return claims, nil
		}
		return auth.Claims{}, fmt.Errorf("query admissions staff profile: %w", err)
	}

	if !profile.Active {
		return claims, nil
	}

	return appendUniqueRoles(claims, admissionsbus.AdmissionsRolesToStrings(profile.Roles)), nil
}

func appendUniqueRoles(claims auth.Claims, roles []string) auth.Claims {
	seen := make(map[string]struct{}, len(claims.Roles)+len(roles))
	for _, role := range claims.Roles {
		seen[role] = struct{}{}
	}

	for _, role := range roles {
		if _, exists := seen[role]; exists {
			continue
		}
		claims.Roles = append(claims.Roles, role)
		seen[role] = struct{}{}
	}

	return claims
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
