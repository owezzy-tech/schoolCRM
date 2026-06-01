package africastalking

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestSmsGatewayContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[SendSMSRequest, SendSMSResponse]{
		{
			Name:    "sends sms fixture",
			Request: FixtureSendSMSRequest(),
			Respond: func(SendSMSRequest) (SendSMSResponse, error) {
				return FixtureSendSMSResponse(), nil
			},
			Assert: func(t *testing.T, got SendSMSResponse, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.MessageID == "" {
					t.Fatal("MessageID is empty")
				}
			},
		},
	})
}
