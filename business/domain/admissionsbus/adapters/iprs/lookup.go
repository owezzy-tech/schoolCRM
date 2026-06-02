// Package iprs defines the outbound port for national identity verification through an authorized aggregator.
package iprs

import (
	"context"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// Lookup verifies Kenyan national identity data through an authorized IPRS aggregator.
type Lookup interface {
	VerifyNationalID(ctx context.Context, req VerifyNationalIDRequest) (VerifyNationalIDResponse, error)
}

// VerifyNationalIDRequest identifies a national ID lookup.
type VerifyNationalIDRequest struct {
	NationalID  admissionsbus.KenyaNationalID
	FirstName   *string
	LastName    *string
	DateOfBirth *time.Time
}

// VerifyNationalIDResponse returns IPRS verification status and receipt metadata.
type VerifyNationalIDResponse struct {
	Verified    bool
	FullName    string
	DateOfBirth *time.Time
	ExternalRef string
	VerifiedAt  time.Time
}
