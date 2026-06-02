package whatsapp

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestCloudGatewayContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[SendTemplateRequest, SendMessageResponse]{
		{
			Name:    "sends template fixture",
			Request: FixtureSendTemplateRequest(),
			Respond: func(SendTemplateRequest) (SendMessageResponse, error) {
				return FixtureSendMessageResponse(), nil
			},
			Assert: func(t *testing.T, got SendMessageResponse, err error) {
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
