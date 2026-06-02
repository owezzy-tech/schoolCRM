package celcomafrica

import "time"

// FixtureSendSMSRequest is a representative Celcom Africa SMS request.
func FixtureSendSMSRequest() SendSMSRequest {
	senderID := "SCHOOLCRM"

	return SendSMSRequest{
		PhoneNumber:   "254712345678",
		Message:       "Your application APP-2026-0001 has been received.",
		SenderID:      &senderID,
		CorrelationID: "fixture-sms-0001",
	}
}

// FixtureSendSMSResponse is a representative Celcom Africa SMS response.
func FixtureSendSMSResponse() SendSMSResponse {
	cost := "KES 0.80"

	return SendSMSResponse{
		MessageID:   "CELCOM-123456789",
		Status:      "Success",
		Cost:        &cost,
		ExternalRef: "CELCOM-SMS-FIXTURE-0001",
		SentAt:      time.Date(2026, time.March, 3, 13, 0, 0, 0, time.UTC),
	}
}

// FixtureDeliveryReportRequest is a representative Celcom Africa delivery report.
func FixtureDeliveryReportRequest() DeliveryReportRequest {
	return DeliveryReportRequest{
		MessageID:   "CELCOM-123456789",
		PhoneNumber: "254712345678",
		Status:      "Delivered",
		ReceivedAt:  time.Date(2026, time.March, 3, 13, 1, 0, 0, time.UTC),
	}
}

// FixtureDeliveryReportResponse is a representative normalized SMS delivery report.
func FixtureDeliveryReportResponse() DeliveryReportResponse {
	return DeliveryReportResponse{
		MessageID:  "CELCOM-123456789",
		Delivered:  true,
		Status:     "Delivered",
		ReceivedAt: time.Date(2026, time.March, 3, 13, 1, 0, 0, time.UTC),
	}
}
