// Package kuccps defines the outbound port for KUCCPS placement and catalogue imports.
package kuccps

import (
	"context"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// Importer imports KUCCPS programmes, cluster cutoffs, and placement-cycle data.
type Importer interface {
	ImportCatalog(ctx context.Context, req ImportCatalogRequest) (ImportCatalogResponse, error)
	ImportPlacements(ctx context.Context, req ImportPlacementsRequest) (ImportPlacementsResponse, error)
}

// ImportCatalogRequest scopes an annual KUCCPS programme and cutoff refresh.
type ImportCatalogRequest struct {
	AcademicYear int
	SourceRef    string
}

// ImportCatalogResponse summarizes a KUCCPS catalogue refresh.
type ImportCatalogResponse struct {
	AcademicYear    int
	ProgrammesCount int
	ClustersCount   int
	CutoffsCount    int
	ExternalRef     string
	ImportedAt      time.Time
}

// ImportPlacementsRequest scopes placement import for one KUCCPS cycle.
type ImportPlacementsRequest struct {
	PlacementCycle string
	SourceRef      string
}

// PlacementRecord is a normalized KUCCPS placement row.
type PlacementRecord struct {
	PlacementID     string
	KCSEIndexNumber string
	NationalID      *string
	ProgrammeCode   string
	ProgrammeName   string
	InstitutionCode string
	PlacementCycle  string
	PlacementYear   int
	AdmissionNumber *string
}

// ImportPlacementsResponse contains imported placements for application matching.
type ImportPlacementsResponse struct {
	PlacementCycle string
	Placements     []PlacementRecord
	ExternalRef    string
	ImportedAt     time.Time
}

// PlacementToSnapshot converts a KUCCPS placement row to the application snapshot shape.
func PlacementToSnapshot(record PlacementRecord) admissionsbus.KUCCPSPlacement {
	return admissionsbus.KUCCPSPlacement{
		PlacementID:     record.PlacementID,
		AdmissionNumber: record.AdmissionNumber,
		ProgrammeCode:   record.ProgrammeCode,
		ProgrammeName:   record.ProgrammeName,
		InstitutionCode: record.InstitutionCode,
		PlacementYear:   record.PlacementYear,
	}
}
