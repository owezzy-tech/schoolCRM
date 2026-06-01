package kuccps

import (
	"testing"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/adapters"
)

func TestImporterContractFixtures(t *testing.T) {
	t.Parallel()

	adapters.RunContractCases(t, []adapters.ContractCase[ImportPlacementsRequest, ImportPlacementsResponse]{
		{
			Name:    "imports placement cycle fixture",
			Request: FixtureImportPlacementsRequest(),
			Respond: func(ImportPlacementsRequest) (ImportPlacementsResponse, error) {
				return FixtureImportPlacementsResponse(), nil
			},
			Assert: func(t *testing.T, got ImportPlacementsResponse, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.PlacementCycle != "2026/2027" {
					t.Fatalf("PlacementCycle = %q, want %q", got.PlacementCycle, "2026/2027")
				}
				if len(got.Placements) != 1 {
					t.Fatalf("Placements length = %d, want 1", len(got.Placements))
				}
			},
		},
	})
}
