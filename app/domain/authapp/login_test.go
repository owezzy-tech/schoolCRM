package authapp

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"

	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
)

// failingUsers is a userbus that rejects every sign-in with the given error.
type failingUsers struct {
	userbus.ExtBusiness
	err error
}

func (f failingUsers) Authenticate(context.Context, mail.Address, string) (userbus.User, error) {
	return userbus.User{}, f.err
}

func loginWith(t *testing.T, cause error) *errs.Error {
	t.Helper()

	a := &app{userBus: failingUsers{err: cause}}
	r := httptest.NewRequest("POST", "/v1/auth/login", strings.NewReader(`{"email":"person@schoolcrm.invalid","password":"not-the-password"}`))

	resp := a.login(context.Background(), r)
	appErr, ok := resp.(*errs.Error)
	if !ok {
		t.Fatalf("login returned %T, want *errs.Error", resp)
	}
	return appErr
}

func TestLoginFailuresDoNotRevealWhetherTheAccountExists(t *testing.T) {
	unknownEmail := loginWith(t, fmt.Errorf("query: email[person@schoolcrm.invalid]: db: %w", userbus.ErrNotFound))
	wrongPassword := loginWith(t, fmt.Errorf("comparehashandpassword: %w", userbus.ErrAuthenticationFailure))

	for name, got := range map[string]*errs.Error{"unknown email": unknownEmail, "wrong password": wrongPassword} {
		if got.Code != errs.Unauthenticated {
			t.Errorf("%s: code = %s, want %s", name, got.Code, errs.Unauthenticated)
		}
		if got.Message != auth.InvalidCredentials {
			t.Errorf("%s: message = %q, want %q", name, got.Message, auth.InvalidCredentials)
		}
	}

	unknownBody, _, _ := unknownEmail.Encode()
	wrongBody, _, _ := wrongPassword.Encode()
	if string(unknownBody) != string(wrongBody) {
		t.Errorf("responses differ:\nunknown email:  %s\nwrong password: %s", unknownBody, wrongBody)
	}
	for _, leak := range []string{"query", "db:", "user not found", "comparehashandpassword", "schoolcrm.invalid"} {
		if strings.Contains(string(unknownBody), leak) {
			t.Errorf("response leaks %q: %s", leak, unknownBody)
		}
	}

	// The internal cause is kept for the error middleware's log.
	if !errors.Is(unknownEmail, userbus.ErrNotFound) {
		t.Error("unknown email: cause lost, want errors.Is(err, userbus.ErrNotFound)")
	}
	if !errors.Is(wrongPassword, userbus.ErrAuthenticationFailure) {
		t.Error("wrong password: cause lost, want errors.Is(err, userbus.ErrAuthenticationFailure)")
	}
}
