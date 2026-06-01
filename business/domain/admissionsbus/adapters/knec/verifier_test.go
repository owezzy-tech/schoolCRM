package knec

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestVerifierContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[VerifyKCSEResultRequest, VerifyKCSEResultResponse]{
		{
			Name:    "verifies kcse result fixture",
			Request: FixtureVerifyKCSEResultRequest(),
			Respond: func(VerifyKCSEResultRequest) (VerifyKCSEResultResponse, error) {
				return FixtureVerifyKCSEResultResponse(), nil
			},
			Assert: func(t *testing.T, got VerifyKCSEResultResponse, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Verified {
					t.Fatal("Verified = false, want true")
				}
				if got.Result.MeanPoints == 0 {
					t.Fatal("MeanPoints = 0, want calculated points")
				}
			},
		},
	})
}
