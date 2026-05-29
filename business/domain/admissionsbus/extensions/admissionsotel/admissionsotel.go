// Package admissionsotel provides an extension for admissionsbus that adds
// otel tracking.
package admissionsotel

import (
	"context"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/otel"
)

// Extension provides a wrapper for otel functionality around the admissionsbus.
type Extension struct {
	bus admissionsbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the admissionsbus with otel.
func NewExtension() admissionsbus.Extension {
	return func(bus admissionsbus.ExtBusiness) admissionsbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

// NewWithTx does not apply otel.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (admissionsbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Health applies otel to the admissions scaffold health check.
func (ext *Extension) Health(ctx context.Context) (admissionsbus.Health, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.health")
	defer span.End()

	return ext.bus.Health(ctx)
}
