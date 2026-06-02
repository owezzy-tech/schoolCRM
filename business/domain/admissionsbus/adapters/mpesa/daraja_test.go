package mpesa

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestDarajaGatewayContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[InitiateSTKPushRequest, InitiateSTKPushResponse]{
		{
			Name:    "initiates stk push fixture",
			Request: FixtureInitiateSTKPushRequest(),
			Respond: func(InitiateSTKPushRequest) (InitiateSTKPushResponse, error) {
				return FixtureInitiateSTKPushResponse(), nil
			},
			Assert: func(t *testing.T, got InitiateSTKPushResponse, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.CheckoutRequestID == "" {
					t.Fatal("CheckoutRequestID is empty")
				}
			},
		},
	})
}
