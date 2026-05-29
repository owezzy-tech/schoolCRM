// Package admissionsotel provides an extension for admissionsbus that adds
// otel tracking.
package admissionsotel

import (
	"context"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
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

// UpsertProgram applies otel to Program sync/import upserts.
func (ext *Extension) UpsertProgram(ctx context.Context, up admissionsbus.UpsertProgram) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.upsertprogram")
	defer span.End()

	return ext.bus.UpsertProgram(ctx, up)
}

// QueryPrograms applies otel to Program queries.
func (ext *Extension) QueryPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprograms")
	defer span.End()

	return ext.bus.QueryPrograms(ctx, filter, orderBy, page)
}

// CountPrograms applies otel to Program counts.
func (ext *Extension) CountPrograms(ctx context.Context, filter admissionsbus.ProgramQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countprograms")
	defer span.End()

	return ext.bus.CountPrograms(ctx, filter)
}

// QueryProgramByID applies otel to Program ID lookups.
func (ext *Extension) QueryProgramByID(ctx context.Context, programID uuid.UUID) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprogrambyid")
	defer span.End()

	return ext.bus.QueryProgramByID(ctx, programID)
}

// QueryProgramByExternalSISID applies otel to Program SIS ID lookups.
func (ext *Extension) QueryProgramByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Program, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryprogrambyexternalsisid")
	defer span.End()

	return ext.bus.QueryProgramByExternalSISID(ctx, externalSISID)
}

// UpsertAcademicTerm applies otel to AcademicTerm sync/import upserts.
func (ext *Extension) UpsertAcademicTerm(ctx context.Context, up admissionsbus.UpsertAcademicTerm) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.upsertacademicterm")
	defer span.End()

	return ext.bus.UpsertAcademicTerm(ctx, up)
}

// QueryAcademicTerms applies otel to AcademicTerm queries.
func (ext *Extension) QueryAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademicterms")
	defer span.End()

	return ext.bus.QueryAcademicTerms(ctx, filter, orderBy, page)
}

// CountAcademicTerms applies otel to AcademicTerm counts.
func (ext *Extension) CountAcademicTerms(ctx context.Context, filter admissionsbus.AcademicTermQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countacademicterms")
	defer span.End()

	return ext.bus.CountAcademicTerms(ctx, filter)
}

// QueryAcademicTermByID applies otel to AcademicTerm ID lookups.
func (ext *Extension) QueryAcademicTermByID(ctx context.Context, termID uuid.UUID) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademictermbyid")
	defer span.End()

	return ext.bus.QueryAcademicTermByID(ctx, termID)
}

// QueryAcademicTermByExternalSISID applies otel to AcademicTerm SIS ID lookups.
func (ext *Extension) QueryAcademicTermByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.AcademicTerm, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryacademictermbyexternalsisid")
	defer span.End()

	return ext.bus.QueryAcademicTermByExternalSISID(ctx, externalSISID)
}
