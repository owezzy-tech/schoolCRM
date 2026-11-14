package errs_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/owezzy/schoolCRM/app/sdk/errs"
)

func TestMaskHidesTheCauseFromClients(t *testing.T) {
	cause := errors.New("db: user not found")
	err := errs.Mask(errs.Unauthenticated, "Invalid email or password.", cause)

	if err.Error() != "Invalid email or password." {
		t.Errorf("Error() = %q, want the client-safe message", err.Error())
	}

	body, _, encErr := err.Encode()
	if encErr != nil {
		t.Fatalf("Encode: %v", encErr)
	}
	if strings.Contains(string(body), "user not found") {
		t.Errorf("encoded body leaks the cause: %s", body)
	}

	if !errors.Is(err, cause) || err.Unwrap() != cause {
		t.Error("the cause must stay reachable for logging")
	}
	if err.HTTPStatus() != 401 {
		t.Errorf("HTTPStatus() = %d, want 401", err.HTTPStatus())
	}
}

func TestNewAndErrorfHaveNoCause(t *testing.T) {
	if errs.New(errs.Internal, errors.New("x")).Unwrap() != nil {
		t.Error("New: Unwrap() should be nil")
	}
	if errs.Errorf(errs.Internal, "x").Unwrap() != nil {
		t.Error("Errorf: Unwrap() should be nil")
	}
}
