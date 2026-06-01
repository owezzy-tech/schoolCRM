// Package mpesa defines the outbound port for M-Pesa Daraja payment operations.
package mpesa

import (
	"context"
	"time"
)

// DarajaGateway initiates STK Push requests and normalizes payment callbacks.
type DarajaGateway interface {
	InitiateSTKPush(ctx context.Context, req InitiateSTKPushRequest) (InitiateSTKPushResponse, error)
	ParseCallback(ctx context.Context, req PaymentCallbackRequest) (PaymentCallbackResponse, error)
}

// InitiateSTKPushRequest starts a customer-present payment request.
type InitiateSTKPushRequest struct {
	PhoneNumber   string
	AmountCents   int64
	AccountRef    string
	Description   string
	CallbackURL   string
	CorrelationID string
}

// InitiateSTKPushResponse contains Daraja checkout references.
type InitiateSTKPushResponse struct {
	MerchantRequestID string
	CheckoutRequestID string
	CustomerMessage   string
	ExternalRef       string
	RequestedAt       time.Time
}

// PaymentCallbackRequest carries a normalized Daraja callback payload for parsing.
type PaymentCallbackRequest struct {
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode        int
	ResultDescription string
	MpesaReceipt      *string
	PhoneNumber       string
	AmountCents       int64
	PaidAt            *time.Time
}

// PaymentCallbackResponse is the durable payment event shape consumed by admissions workflows.
type PaymentCallbackResponse struct {
	Succeeded         bool
	MerchantRequestID string
	CheckoutRequestID string
	MpesaReceipt      *string
	FailureReason     *string
	AmountCents       int64
	PaidAt            *time.Time
}
