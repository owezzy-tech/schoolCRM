package apitest

import (
	"net/http/httptest"
	"testing"

	authbuild "github.com/owezzy/schoolCRM/api/services/auth/build"
	schoolcrmbuild "github.com/owezzy/schoolCRM/api/services/schoolcrm/build"
	"github.com/owezzy/schoolCRM/app/sdk/auth"
	"github.com/owezzy/schoolCRM/app/sdk/authclient/http"
	"github.com/owezzy/schoolCRM/app/sdk/mux"
	"github.com/owezzy/schoolCRM/business/sdk/dbtest"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	db := dbtest.New(t, testName)

	// -------------------------------------------------------------------------

	auth := auth.New(auth.Config{
		Log:       db.Log,
		UserBus:   db.BusDomain.User,
		KeyLookup: &KeyStore{},
	})

	// -------------------------------------------------------------------------

	server := httptest.NewServer(mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus: db.BusDomain.User,
		},
		AuthConfig: mux.AuthConfig{
			Auth:     auth,
			TokenKey: kid,
		},
	}, authbuild.Routes()))

	authClient, err := http.New(db.Log, server.URL)
	if err != nil {
		t.Fatal("could not create authentication client")
	}

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			AuditBus:    db.BusDomain.Audit,
			UserBus:     db.BusDomain.User,
			ProductBus:  db.BusDomain.Product,
			HomeBus:     db.BusDomain.Home,
			VProductBus: db.BusDomain.VProduct,
		},
		SchoolCRMConfig: mux.SchoolCRMConfig{
			AuthClient: authClient,
		},
	}, schoolcrmbuild.Routes())

	return &Test{
		DB:   db,
		Auth: auth,
		mux:  mux,
	}
}
