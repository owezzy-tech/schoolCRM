package unittest

import (
	"context"

	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/homebus"
	"github.com/owezzy/schoolCRM/business/domain/productbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
)

// User represents an app user specified for the test.
type User struct {
	userbus.User
	Products []productbus.Product
	Homes    []homebus.Home
	Audits   []auditbus.Audit
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an unit test.
type Table struct {
	Name    string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}
