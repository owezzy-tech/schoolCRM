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

// CreateConstituent applies otel to Constituent creation.
func (ext *Extension) CreateConstituent(ctx context.Context, nc admissionsbus.NewConstituent) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createconstituent")
	defer span.End()

	return ext.bus.CreateConstituent(ctx, nc)
}

// UpdateConstituent applies otel to Constituent updates.
func (ext *Extension) UpdateConstituent(ctx context.Context, cst admissionsbus.Constituent, uc admissionsbus.UpdateConstituent) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.updateconstituent")
	defer span.End()

	return ext.bus.UpdateConstituent(ctx, cst, uc)
}

// QueryConstituents applies otel to Constituent queries.
func (ext *Extension) QueryConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituents")
	defer span.End()

	return ext.bus.QueryConstituents(ctx, filter, orderBy, page)
}

// CountConstituents applies otel to Constituent counts.
func (ext *Extension) CountConstituents(ctx context.Context, filter admissionsbus.ConstituentQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countconstituents")
	defer span.End()

	return ext.bus.CountConstituents(ctx, filter)
}

// QueryConstituentByID applies otel to Constituent ID lookups.
func (ext *Extension) QueryConstituentByID(ctx context.Context, constituentID uuid.UUID) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyid")
	defer span.End()

	return ext.bus.QueryConstituentByID(ctx, constituentID)
}

// QueryConstituentByPrimaryEmail applies otel to Constituent email lookups.
func (ext *Extension) QueryConstituentByPrimaryEmail(ctx context.Context, email string) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyprimaryemail")
	defer span.End()

	return ext.bus.QueryConstituentByPrimaryEmail(ctx, email)
}

// QueryConstituentByExternalSISID applies otel to Constituent SIS ID lookups.
func (ext *Extension) QueryConstituentByExternalSISID(ctx context.Context, externalSISID string) (admissionsbus.Constituent, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryconstituentbyexternalsisid")
	defer span.End()

	return ext.bus.QueryConstituentByExternalSISID(ctx, externalSISID)
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

// CreateDuplicateReview applies otel to duplicate review creation.
func (ext *Extension) CreateDuplicateReview(ctx context.Context, nr admissionsbus.NewDuplicateReview) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createduplicatereview")
	defer span.End()

	return ext.bus.CreateDuplicateReview(ctx, nr)
}

// ResolveDuplicateReview applies otel to duplicate review resolution.
func (ext *Extension) ResolveDuplicateReview(ctx context.Context, review admissionsbus.DuplicateReview, rr admissionsbus.ResolveDuplicateReview) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.resolveduplicatereview")
	defer span.End()

	return ext.bus.ResolveDuplicateReview(ctx, review, rr)
}

// QueryDuplicateReviews applies otel to duplicate review queries.
func (ext *Extension) QueryDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryduplicatereviews")
	defer span.End()

	return ext.bus.QueryDuplicateReviews(ctx, filter, orderBy, page)
}

// CountDuplicateReviews applies otel to duplicate review counts.
func (ext *Extension) CountDuplicateReviews(ctx context.Context, filter admissionsbus.DuplicateReviewQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countduplicatereviews")
	defer span.End()

	return ext.bus.CountDuplicateReviews(ctx, filter)
}

// QueryDuplicateReviewByID applies otel to duplicate review ID lookups.
func (ext *Extension) QueryDuplicateReviewByID(ctx context.Context, reviewID uuid.UUID) (admissionsbus.DuplicateReview, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryduplicatereviewbyid")
	defer span.End()

	return ext.bus.QueryDuplicateReviewByID(ctx, reviewID)
}

// CreateApplication applies otel to Application creation.
func (ext *Extension) CreateApplication(ctx context.Context, na admissionsbus.NewApplication) (admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.createapplication")
	defer span.End()

	return ext.bus.CreateApplication(ctx, na)
}

// QueryApplications applies otel to Application queries.
func (ext *Extension) QueryApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter, orderBy order.By, page page.Page) ([]admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplications")
	defer span.End()

	return ext.bus.QueryApplications(ctx, filter, orderBy, page)
}

// CountApplications applies otel to Application counts.
func (ext *Extension) CountApplications(ctx context.Context, filter admissionsbus.ApplicationQueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.countapplications")
	defer span.End()

	return ext.bus.CountApplications(ctx, filter)
}

// QueryApplicationByID applies otel to Application ID lookups.
func (ext *Extension) QueryApplicationByID(ctx context.Context, applicationID uuid.UUID) (admissionsbus.Application, error) {
	ctx, span := otel.AddSpan(ctx, "business.admissionsbus.queryapplicationbyid")
	defer span.End()

	return ext.bus.QueryApplicationByID(ctx, applicationID)
}
