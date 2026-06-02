package whatsapp

import "time"

// FixtureSendTemplateRequest is a representative WhatsApp Cloud template request.
func FixtureSendTemplateRequest() SendTemplateRequest {
	return SendTemplateRequest{
		PhoneNumber:  "254712345678",
		TemplateName: "application_received",
		LanguageCode: "en",
		Parameters: []TemplateParameter{
			{Name: "application_ref", Value: "APP-2026-0001"},
		},
		CorrelationID: "fixture-whatsapp-0001",
	}
}

// FixtureSendMessageResponse is a representative WhatsApp Cloud message response.
func FixtureSendMessageResponse() SendMessageResponse {
	return SendMessageResponse{
		MessageID:   "wamid.HBgMMjU0NzEyMzQ1Njc4FQIAERgSFixture",
		Status:      "accepted",
		ExternalRef: "WA-CLOUD-FIXTURE-0001",
		SentAt:      time.Date(2026, time.March, 3, 14, 0, 0, 0, time.UTC),
	}
}

// FixtureDeliveryReportRequest is a representative WhatsApp delivery report.
func FixtureDeliveryReportRequest() DeliveryReportRequest {
	return DeliveryReportRequest{
		MessageID:   "wamid.HBgMMjU0NzEyMzQ1Njc4FQIAERgSFixture",
		PhoneNumber: "254712345678",
		Status:      "delivered",
		ReceivedAt:  time.Date(2026, time.March, 3, 14, 1, 0, 0, time.UTC),
	}
}

// FixtureDeliveryReportResponse is a representative normalized WhatsApp delivery report.
func FixtureDeliveryReportResponse() DeliveryReportResponse {
	return DeliveryReportResponse{
		MessageID:  "wamid.HBgMMjU0NzEyMzQ1Njc4FQIAERgSFixture",
		Delivered:  true,
		Status:     "delivered",
		ReceivedAt: time.Date(2026, time.March, 3, 14, 1, 0, 0, time.UTC),
	}
}
