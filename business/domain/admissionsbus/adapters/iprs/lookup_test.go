package iprs

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestLookupContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[VerifyNationalIDRequest, VerifyNationalIDResponse]{
		{
			Name:    "verifies national id fixture",
			Request: FixtureVerifyNationalIDRequest(),
			Respond: func(VerifyNationalIDRequest) (VerifyNationalIDResponse, error) {
				return FixtureVerifyNationalIDResponse(), nil
			},
			Assert: func(t *testing.T, got VerifyNationalIDResponse, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Verified {
					t.Fatal("Verified = false, want true")
				}
				if got.ExternalRef == "" {
					t.Fatal("ExternalRef is empty")
				}
			},
		},
	})
}
