// Package authapp maintains the web based api for auth access.
package authapp

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/types/role"
	"github.com/owezzy/schoolCRM/foundation/web"
)

type app struct {
	auth     *auth.Auth
	userBus  userbus.ExtBusiness
	tokenKey string
}

func newApp(ath *auth.Auth, userBus userbus.ExtBusiness, tokenKey string) *app {
	return &app{
		auth:     ath,
		userBus:  userBus,
		tokenKey: tokenKey,
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
