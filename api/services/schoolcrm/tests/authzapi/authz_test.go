// Package authz_test checks every homes, products and users route against
// every caller type. The expected codes follow app/sdk/auth/rego/authorization.rego:
//
//   - rule_any: any school role (SUPER_ADMIN, SCHOOL_ADMIN, TEACHER, STUDENT, PARENT).
//   - rule_user_only: TEACHER, STUDENT or PARENT.
//   - rule_admin_only: SUPER_ADMIN or SCHOOL_ADMIN.
//   - rule_admin_or_subject: an admin, or a user role acting on its own record.
//
// Authentication failures return 401. Role denials return 403, and on record
// routes with a role-only rule the role is checked before the record is loaded,
// so the 403 reveals nothing about the record. Ownership denials on record
// routes return the same 404 as a record that does not exist, so callers
// cannot learn which records exist.
package authz_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"testing"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/domain/homeapp"
	"github.com/owezzy/schoolCRM/app/domain/productapp"
	"github.com/owezzy/schoolCRM/app/domain/userapp"
	"github.com/owezzy/schoolCRM/app/sdk/apitest"
	"github.com/owezzy/schoolCRM/business/domain/homebus"
	"github.com/owezzy/schoolCRM/business/domain/productbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
	"github.com/owezzy/schoolCRM/business/types/role"
)

// caller names the identity that sends a request.
type caller string

const (
	anonymous   caller = "anonymous"    // no token
	student     caller = "student"      // STUDENT that owns the record
	wrongUser   caller = "wronguser"    // another STUDENT
	schoolAdmin caller = "school_admin" // SCHOOL_ADMIN
)

var callers = []caller{anonymous, student, wrongUser, schoolAdmin}

// fixture is the state one matrix cell runs against. Each cell gets a new
// owner and record, so deletes and updates cannot affect other cells.
type fixture struct {
	tokens   map[caller]string
	userID   uuid.UUID
	homeID   uuid.UUID
	product  uuid.UUID
	newEmail string
}

// route describes one endpoint and the status each caller must receive.
type route struct {
	method string
	path   func(f fixture) string
	input  func(f fixture) any
	// record marks routes scoped to one record by id.
	record bool
	expect map[caller]int
}

func Test_AuthzMatrix(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_AuthzMatrix")
	s := seeder{t: t, test: test}

	// These callers are never modified by a cell, so they are shared.
	other := s.user(role.Student)
	admin := s.user(role.SchoolAdmin)

	for _, rt := range routes() {
		for _, c := range callers {
			name := fmt.Sprintf("%s %s/%s", rt.method, rt.name(), c)

			t.Run(name, func(t *testing.T) {
				f := s.fixture(other, admin)

				got, body := rt.do(t, test, f, c)
				if exp := rt.expect[c]; got != exp {
					t.Fatalf("status: got %d, want %d: %s", got, exp, body)
				}
			})
		}

		if !rt.record {
			continue
		}

		// A record the caller may not access and a record that does not
		// exist must be indistinguishable.
		t.Run(fmt.Sprintf("%s %s/wronguser-matches-missing", rt.method, rt.name()), func(t *testing.T) {
			f := s.fixture(other, admin)

			wrongCode, wrongBody := rt.do(t, test, f, wrongUser)
			missingCode, missingBody := rt.do(t, test, f.missing(), wrongUser)

			if exp := rt.expect[wrongUser]; wrongCode != exp || missingCode != exp {
				t.Fatalf("status: wronguser %d, missing %d, want both %d", wrongCode, missingCode, exp)
			}

			if !bytes.Equal(wrongBody, missingBody) {
				t.Fatalf("body differs:\nwronguser: %s\nmissing:   %s", wrongBody, missingBody)
			}
		})

		// Admins pass the policy for every record, so they also get 404 for a
		// record that does not exist.
		t.Run(fmt.Sprintf("%s %s/school_admin-missing", rt.method, rt.name()), func(t *testing.T) {
			f := s.fixture(other, admin).missing()

			if got, body := rt.do(t, test, f, schoolAdmin); got != http.StatusNotFound {
				t.Fatalf("status: got %d, want %d: %s", got, http.StatusNotFound, body)
			}
		})
	}
}

func (rt route) name() string {
	return rt.path(fixture{})
}

func (rt route) do(t *testing.T, test *apitest.Test, f fixture, c caller) (int, []byte) {
	t.Helper()

	var input any
	if rt.input != nil {
		input = rt.input(f)
	}

	return test.Do(t, rt.method, rt.path(f), f.tokens[c], input)
}

// missing returns a copy of f whose record ids do not exist.
func (f fixture) missing() fixture {
	f.userID = uuid.New()
	f.homeID = uuid.New()
	f.product = uuid.New()

	return f
}

func routes() []route {
	const (
		ok        = http.StatusOK
		noContent = http.StatusNoContent
		unauth    = http.StatusUnauthorized
		forbidden = http.StatusForbidden
		notFound  = http.StatusNotFound
	)

	homePath := func(f fixture) string { return "/v1/homes/" + id(f.homeID) }
	productPath := func(f fixture) string { return "/v1/products/" + id(f.product) }
	userPath := func(f fixture) string { return "/v1/users/" + id(f.userID) }

	newHome := func(fixture) any {
		return &homeapp.NewHome{
			Type: "SINGLE FAMILY",
			Address: homeapp.NewAddress{
				Address1: "123 Mocking Bird Lane",
				ZipCode:  "35810",
				City:     "Huntsville",
				State:    "AL",
				Country:  "US",
			},
		}
	}
	updateHome := func(fixture) any {
		return &homeapp.UpdateHome{Type: dbtest.StringPointer("CONDO")}
	}
	newProduct := func(fixture) any {
		return &productapp.NewProduct{Name: "Guitar", Cost: 10.34, Quantity: 10}
	}
	updateProduct := func(fixture) any {
		return &productapp.UpdateProduct{Name: dbtest.StringPointer("Drums")}
	}
	newUser := func(f fixture) any {
		return &userapp.NewUser{
			Name:            "Bill Kennedy",
			Email:           f.newEmail,
			Roles:           []string{"STUDENT"},
			Department:      "ITO",
			Password:        "123",
			PasswordConfirm: "123",
		}
	}
	updateUser := func(fixture) any {
		return &userapp.UpdateUser{
			Name:            dbtest.StringPointer("Jack Kennedy"),
			Password:        dbtest.StringPointer("123"),
			PasswordConfirm: dbtest.StringPointer("123"),
		}
	}
	updateRole := func(fixture) any {
		return &userapp.UpdateUserRole{Roles: []string{"STUDENT"}}
	}

	return []route{
		// homes ---------------------------------------------------------------
		{
			method: http.MethodGet, path: func(fixture) string { return "/v1/homes" },
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: ok, schoolAdmin: ok},
		},
		{
			method: http.MethodGet, path: homePath, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodPost, path: func(fixture) string { return "/v1/homes" }, input: newHome,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: ok, schoolAdmin: forbidden},
		},
		{
			method: http.MethodPut, path: homePath, input: updateHome, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodDelete, path: homePath, record: true,
			expect: map[caller]int{anonymous: unauth, student: noContent, wrongUser: notFound, schoolAdmin: noContent},
		},

		// products ------------------------------------------------------------
		{
			method: http.MethodGet, path: func(fixture) string { return "/v1/products" },
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: ok, schoolAdmin: ok},
		},
		{
			method: http.MethodGet, path: productPath, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodPost, path: func(fixture) string { return "/v1/products" }, input: newProduct,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: ok, schoolAdmin: forbidden},
		},
		{
			method: http.MethodPut, path: productPath, input: updateProduct, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodDelete, path: productPath, record: true,
			expect: map[caller]int{anonymous: unauth, student: noContent, wrongUser: notFound, schoolAdmin: noContent},
		},

		// users ---------------------------------------------------------------
		{
			method: http.MethodGet, path: func(fixture) string { return "/v1/users" },
			expect: map[caller]int{anonymous: unauth, student: forbidden, wrongUser: forbidden, schoolAdmin: ok},
		},
		{
			method: http.MethodGet, path: userPath, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodPost, path: func(fixture) string { return "/v1/users" }, input: newUser,
			expect: map[caller]int{anonymous: unauth, student: forbidden, wrongUser: forbidden, schoolAdmin: ok},
		},
		{
			// rule_admin_only on a record route: the role is checked before
			// the record is loaded, so a STUDENT gets 403 for any id.
			method: http.MethodPut, path: func(f fixture) string { return "/v1/users/role/" + id(f.userID) }, input: updateRole, record: true,
			expect: map[caller]int{anonymous: unauth, student: forbidden, wrongUser: forbidden, schoolAdmin: ok},
		},
		{
			method: http.MethodPut, path: userPath, input: updateUser, record: true,
			expect: map[caller]int{anonymous: unauth, student: ok, wrongUser: notFound, schoolAdmin: ok},
		},
		{
			method: http.MethodDelete, path: userPath, record: true,
			expect: map[caller]int{anonymous: unauth, student: noContent, wrongUser: notFound, schoolAdmin: noContent},
		},
	}
}

// id prints a placeholder for uuid.Nil so subtest names stay readable.
func id(v uuid.UUID) string {
	if v == uuid.Nil {
		return "{id}"
	}

	return v.String()
}

// =============================================================================

type seeded struct {
	user  userbus.User
	token string
}

type seeder struct {
	t    *testing.T
	test *apitest.Test
}

// user creates a user with a unique email. The shared userbus seed helper
// picks emails from a small random range, which can collide across the many
// users this test creates.
func (s seeder) user(rle role.Role) seeded {
	s.t.Helper()

	ctx := context.Background()
	bus := s.test.DB.BusDomain.User

	nu := userbus.TestNewUsers(1, rle)[0]
	nu.Email = mail.Address{Address: uniqueEmail()}

	usr, err := bus.Create(ctx, uuid.UUID{}, nu)
	if err != nil {
		s.t.Fatalf("seeding user: %s", err)
	}

	return seeded{
		user:  usr,
		token: apitest.Token(bus, s.test.Auth, usr.Email.Address),
	}
}

func (s seeder) fixture(other seeded, admin seeded) fixture {
	s.t.Helper()

	ctx := context.Background()
	bus := s.test.DB.BusDomain
	owner := s.user(role.Student)

	hmes, err := homebus.TestGenerateSeedHomes(ctx, 1, bus.Home, owner.user.ID)
	if err != nil {
		s.t.Fatalf("seeding home: %s", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 1, bus.Product, owner.user.ID)
	if err != nil {
		s.t.Fatalf("seeding product: %s", err)
	}

	return fixture{
		tokens: map[caller]string{
			anonymous:   "",
			student:     owner.token,
			wrongUser:   other.token,
			schoolAdmin: admin.token,
		},
		userID:   owner.user.ID,
		homeID:   hmes[0].ID,
		product:  prds[0].ID,
		newEmail: uniqueEmail(),
	}
}

func uniqueEmail() string {
	return fmt.Sprintf("authz-%s@schoolcrm.local", uuid.NewString())
}
