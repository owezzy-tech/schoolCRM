// Package whatsapp defines the outbound port for WhatsApp Cloud messaging.
package whatsapp

import (
	"context"
	"time"
)

// CloudGateway sends approved templates, session messages, and normalizes delivery reports.
type CloudGateway interface {
	SendTemplate(ctx context.Context, req SendTemplateRequest) (SendMessageResponse, error)
	SendSessionMessage(ctx context.Context, req SendSessionMessageRequest) (SendMessageResponse, error)
	ParseDeliveryReport(ctx context.Context, req DeliveryReportRequest) (DeliveryReportResponse, error)
}

// TemplateParameter is one WhatsApp template parameter.
type TemplateParameter struct {
	Name  string
	Value string
}

// SendTemplateRequest sends an approved WhatsApp template.
type SendTemplateRequest struct {
	PhoneNumber   string
	TemplateName  string
	LanguageCode  string
	Parameters    []TemplateParameter
	CorrelationID string
}

// SendSessionMessageRequest sends a free-form WhatsApp session message.
type SendSessionMessageRequest struct {
	PhoneNumber   string
	Message       string
	CorrelationID string
}

// SendMessageResponse contains WhatsApp message identifiers for reconciliation.
type SendMessageResponse struct {
	MessageID   string
	Status      string
	ExternalRef string
	SentAt      time.Time
}

// DeliveryReportRequest carries a normalized WhatsApp delivery webhook event.
type DeliveryReportRequest struct {
	MessageID   string
	PhoneNumber string
	Status      string
	FailureCode *string
	ReceivedAt  time.Time
}

// DeliveryReportResponse is the durable WhatsApp delivery event consumed by admissions workflows.
type DeliveryReportResponse struct {
	MessageID     string
	Delivered     bool
	Status        string
	FailureReason *string
	ReceivedAt    time.Time
}
