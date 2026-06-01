package kuccps

import "time"

// FixtureImportCatalogRequest is a representative KUCCPS catalogue import request.
func FixtureImportCatalogRequest() ImportCatalogRequest {
	return ImportCatalogRequest{
		AcademicYear: 2026,
		SourceRef:    "fixture-kuccps-catalog-2026.csv",
	}
}

// FixtureImportCatalogResponse is a representative KUCCPS catalogue import response.
func FixtureImportCatalogResponse() ImportCatalogResponse {
	return ImportCatalogResponse{
		AcademicYear:    2026,
		ProgrammesCount: 2,
		ClustersCount:   1,
		CutoffsCount:    2,
		ExternalRef:     "KUCCPS-CATALOG-2026-FIXTURE",
		ImportedAt:      time.Date(2026, time.March, 1, 9, 0, 0, 0, time.UTC),
	}
}

// FixtureImportPlacementsRequest is a representative KUCCPS placement import request.
func FixtureImportPlacementsRequest() ImportPlacementsRequest {
	return ImportPlacementsRequest{
		PlacementCycle: "2026/2027",
		SourceRef:      "fixture-kuccps-placements-2026.csv",
	}
}

// FixtureImportPlacementsResponse is a representative KUCCPS placement import response.
func FixtureImportPlacementsResponse() ImportPlacementsResponse {
	admissionNumber := "ADM-2026-0001"
	return ImportPlacementsResponse{
		PlacementCycle: "2026/2027",
		Placements: []PlacementRecord{
			{
				PlacementID:     "KUCCPS-2026-12345",
				KCSEIndexNumber: "12345",
				ProgrammeCode:   "UON-BCOM",
				ProgrammeName:   "Bachelor of Commerce",
				InstitutionCode: "UON",
				PlacementCycle:  "2026/2027",
				PlacementYear:   2026,
				AdmissionNumber: &admissionNumber,
			},
		},
		ExternalRef: "KUCCPS-PLACEMENTS-2026-FIXTURE",
		ImportedAt:  time.Date(2026, time.March, 2, 9, 0, 0, 0, time.UTC),
	}
}
