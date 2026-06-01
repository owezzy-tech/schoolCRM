package mpesa

import "time"

// FixtureInitiateSTKPushRequest is a representative Daraja STK Push request.
func FixtureInitiateSTKPushRequest() InitiateSTKPushRequest {
	return InitiateSTKPushRequest{
		PhoneNumber:   "254712345678",
		AmountCents:   150000,
		AccountRef:    "APP-2026-0001",
		Description:   "Application fee",
		CallbackURL:   "https://school.example.test/payments/mpesa/callback",
		CorrelationID: "fixture-payment-0001",
	}
}

// FixtureInitiateSTKPushResponse is a representative Daraja STK Push response.
func FixtureInitiateSTKPushResponse() InitiateSTKPushResponse {
	return InitiateSTKPushResponse{
		MerchantRequestID: "29115-34620561-1",
		CheckoutRequestID: "ws_CO_0303202612000012345678",
		CustomerMessage:   "Success. Request accepted for processing",
		ExternalRef:       "DARAJA-STK-FIXTURE-0001",
		RequestedAt:       time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC),
	}
}

// FixturePaymentCallbackRequest is a representative successful Daraja callback.
func FixturePaymentCallbackRequest() PaymentCallbackRequest {
	receipt := "RCD1234567"
	paidAt := time.Date(2026, time.March, 3, 12, 2, 0, 0, time.UTC)

	return PaymentCallbackRequest{
		MerchantRequestID: "29115-34620561-1",
		CheckoutRequestID: "ws_CO_0303202612000012345678",
		ResultCode:        0,
		ResultDescription: "The service request is processed successfully.",
		MpesaReceipt:      &receipt,
		PhoneNumber:       "254712345678",
		AmountCents:       150000,
		PaidAt:            &paidAt,
	}
}

// FixturePaymentCallbackResponse is a representative normalized Daraja payment event.
func FixturePaymentCallbackResponse() PaymentCallbackResponse {
	receipt := "RCD1234567"
	paidAt := time.Date(2026, time.March, 3, 12, 2, 0, 0, time.UTC)

	return PaymentCallbackResponse{
		Succeeded:         true,
		MerchantRequestID: "29115-34620561-1",
		CheckoutRequestID: "ws_CO_0303202612000012345678",
		MpesaReceipt:      &receipt,
		AmountCents:       150000,
		PaidAt:            &paidAt,
	}
}
