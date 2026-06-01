// Package africastalking defines the outbound port for SMS delivery through Africa's Talking.
package africastalking

import (
	"context"
	"time"
)

// SmsGateway sends SMS messages and normalizes delivery reports.
type SmsGateway interface {
	SendSMS(ctx context.Context, req SendSMSRequest) (SendSMSResponse, error)
	ParseDeliveryReport(ctx context.Context, req DeliveryReportRequest) (DeliveryReportResponse, error)
}

// SendSMSRequest sends one SMS notification to a Kenyan phone number.
type SendSMSRequest struct {
	PhoneNumber   string
	Message       string
	SenderID      *string
	CorrelationID string
}

// SendSMSResponse contains provider message identifiers for reconciliation.
type SendSMSResponse struct {
	MessageID   string
	Status      string
	Cost        *string
	ExternalRef string
	SentAt      time.Time
}

// DeliveryReportRequest carries an Africa's Talking delivery callback in normalized form.
type DeliveryReportRequest struct {
	MessageID   string
	PhoneNumber string
	Status      string
	FailureCode *string
	ReceivedAt  time.Time
}

// DeliveryReportResponse is the durable SMS delivery event consumed by admissions workflows.
type DeliveryReportResponse struct {
	MessageID     string
	Delivered     bool
	Status        string
	FailureReason *string
	ReceivedAt    time.Time
}
