package iprs

import (
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// FixtureVerifyNationalIDRequest is a representative IPRS lookup request.
func FixtureVerifyNationalIDRequest() VerifyNationalIDRequest {
	firstName := "Amina"
	lastName := "Wanjiku"
	birthDate := time.Date(2006, time.January, 15, 0, 0, 0, 0, time.UTC)

	return VerifyNationalIDRequest{
		NationalID:  admissionsbus.MustParseKenyaNationalID("12345678"),
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &birthDate,
	}
}

// FixtureVerifyNationalIDResponse is a representative IPRS lookup response.
func FixtureVerifyNationalIDResponse() VerifyNationalIDResponse {
	birthDate := time.Date(2006, time.January, 15, 0, 0, 0, 0, time.UTC)

	return VerifyNationalIDResponse{
		Verified:    true,
		FullName:    "Amina Wanjiku",
		DateOfBirth: &birthDate,
		ExternalRef: "IPRS-VERIFY-12345678",
		VerifiedAt:  time.Date(2026, time.March, 3, 11, 0, 0, 0, time.UTC),
	}
}
