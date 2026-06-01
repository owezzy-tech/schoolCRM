// Package knec defines the outbound port for KNEC KCSE result verification.
package knec

import (
	"context"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// Verifier verifies KCSE result slips, portal extracts, or QR decoded payloads.
type Verifier interface {
	VerifyKCSEResult(ctx context.Context, req VerifyKCSEResultRequest) (VerifyKCSEResultResponse, error)
}

// VerifyKCSEResultRequest identifies a KCSE result verification attempt.
type VerifyKCSEResultRequest struct {
	IndexNumber  admissionsbus.KenyaKCSEIndexNumber
	ExamYear     int
	PortalOutput *string
	QRPayload    *string
}

// VerifyKCSEResultResponse returns KNEC verification metadata and normalized result data.
type VerifyKCSEResultResponse struct {
	Verified    bool
	Result      admissionsbus.KCSEResult
	Candidate   string
	ExternalRef string
	VerifiedAt  time.Time
}
