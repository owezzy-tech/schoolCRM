package referencebus_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/referencebus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
)

type stubStore struct {
	countyFilter    referencebus.QueryFilter
	subCountyFilter referencebus.QueryFilter
	wardFilter      referencebus.QueryFilter
	programmeFilter referencebus.QueryFilter
	countyPage      page.Page
	subCountyPage   page.Page
	wardPage        page.Page
	programmePage   page.Page
}

func (s *stubStore) QueryCounties(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.County, error) {
	s.countyFilter = filter
	s.countyPage = page

	return []referencebus.County{{
		Code:        "1",
		Name:        "Mombasa",
		DateCreated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountCounties(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QuerySubCounties(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.SubCounty, error) {
	s.subCountyFilter = filter
	s.subCountyPage = page

	return []referencebus.SubCounty{{
		Code:        "1",
		CountyCode:  "1",
		Name:        "Changamwe",
		DateCreated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountSubCounties(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QueryWards(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Ward, error) {
	s.wardFilter = filter
	s.wardPage = page

	return []referencebus.Ward{{
		Code:          "W001",
		CountyCode:    "1",
		SubCountyCode: "1",
		Name:          "Port Reitz",
		DateCreated:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountWards(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QueryUniversities(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.University, error) {
	return []referencebus.University{stubUniversity()}, nil
}

func (s *stubStore) CountUniversities(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QueryClusters(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Cluster, error) {
	return []referencebus.Cluster{{
		Code:        "CL02",
		Name:        "Engineering and Technology",
		Description: "Programmes whose placement cluster emphasizes mathematics, physics and technical sciences.",
		DateCreated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountClusters(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QueryKNQFLevels(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.KNQFLevel, error) {
	return []referencebus.KNQFLevel{{
		Code:          "KNQF-7",
		Level:         7,
		Name:          "KNQF Level 7",
		Descriptor:    "Broad professional knowledge and analytical skills.",
		Qualification: "Bachelor Degree",
		DateCreated:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountKNQFLevels(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func (s *stubStore) QueryProgrammes(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Programme, error) {
	s.programmeFilter = filter
	s.programmePage = page

	return []referencebus.Programme{{
		Code:           "JKUAT-BENG-CIVIL",
		UniversityCode: "JKUAT",
		ClusterCode:    "CL02",
		KNQFLevelCode:  "KNQF-7",
		Name:           "Bachelor of Science in Civil Engineering",
		AwardType:      "BACHELOR",
		DateCreated:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}

func (s *stubStore) CountProgrammes(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	return 1, nil
}

func stubUniversity() referencebus.University {
	return referencebus.University{
		Code:            "UON",
		Name:            "University of Nairobi",
		InstitutionType: "PUBLIC_UNIVERSITY",
		DateCreated:     time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateUpdated:     time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestQueryCountyByCode(t *testing.T) {
	ctx := context.Background()
	store := &stubStore{}
	bus := referencebus.NewBusiness(store)

	county, err := bus.QueryCountyByCode(ctx, "1")
	if err != nil {
		t.Fatalf("QueryCountyByCode error: %v", err)
	}

	if county.Code != "1" {
		t.Fatalf("expected county code 1, got %q", county.Code)
	}

	if store.countyFilter.Code == nil || *store.countyFilter.Code != "1" {
		t.Fatalf("expected county filter code 1, got %#v", store.countyFilter.Code)
	}

	if store.countyPage.RowsPerPage() != 1 {
		t.Fatalf("expected single-row lookup, got %d", store.countyPage.RowsPerPage())
	}
}

func TestQuerySubCountiesByCounty(t *testing.T) {
	ctx := context.Background()
	store := &stubStore{}
	bus := referencebus.NewBusiness(store)

	subCounties, err := bus.QuerySubCountiesByCounty(ctx, "1")
	if err != nil {
		t.Fatalf("QuerySubCountiesByCounty error: %v", err)
	}

	if len(subCounties) != 1 {
		t.Fatalf("expected 1 sub-county, got %d", len(subCounties))
	}

	if store.subCountyFilter.CountyCode == nil || *store.subCountyFilter.CountyCode != "1" {
		t.Fatalf("expected county filter code 1, got %#v", store.subCountyFilter.CountyCode)
	}

	if store.subCountyPage.RowsPerPage() != 100 {
		t.Fatalf("expected maximum page size lookup, got %d", store.subCountyPage.RowsPerPage())
	}
}

func TestQueryWardsBySubCounty(t *testing.T) {
	ctx := context.Background()
	store := &stubStore{}
	bus := referencebus.NewBusiness(store)

	wards, err := bus.QueryWardsBySubCounty(ctx, "1")
	if err != nil {
		t.Fatalf("QueryWardsBySubCounty error: %v", err)
	}

	if len(wards) != 1 {
		t.Fatalf("expected 1 ward, got %d", len(wards))
	}

	if store.wardFilter.SubCountyCode == nil || *store.wardFilter.SubCountyCode != "1" {
		t.Fatalf("expected sub-county filter code 1, got %#v", store.wardFilter.SubCountyCode)
	}

	if store.wardPage.RowsPerPage() != 100 {
		t.Fatalf("expected maximum page size lookup, got %d", store.wardPage.RowsPerPage())
	}
}

func TestQueryProgrammesByCluster(t *testing.T) {
	ctx := context.Background()
	store := &stubStore{}
	bus := referencebus.NewBusiness(store)

	programmes, err := bus.QueryProgrammesByCluster(ctx, "CL02")
	if err != nil {
		t.Fatalf("QueryProgrammesByCluster error: %v", err)
	}

	if len(programmes) != 1 {
		t.Fatalf("expected 1 programme, got %d", len(programmes))
	}

	if store.programmeFilter.ClusterCode == nil || *store.programmeFilter.ClusterCode != "CL02" {
		t.Fatalf("expected cluster filter code CL02, got %#v", store.programmeFilter.ClusterCode)
	}

	if store.programmePage.RowsPerPage() != 100 {
		t.Fatalf("expected maximum page size lookup, got %d", store.programmePage.RowsPerPage())
	}
}

func TestReferenceCatalogInterfacesAreReadOnly(t *testing.T) {
	for _, iface := range []struct {
		name string
		typ  reflect.Type
	}{
		{name: "Storer", typ: reflect.TypeOf((*referencebus.Storer)(nil)).Elem()},
		{name: "ExtBusiness", typ: reflect.TypeOf((*referencebus.ExtBusiness)(nil)).Elem()},
	} {
		for i := 0; i < iface.typ.NumMethod(); i++ {
			method := iface.typ.Method(i).Name
			if strings.HasPrefix(method, "Create") || strings.HasPrefix(method, "Update") || strings.HasPrefix(method, "Delete") {
				t.Fatalf("%s exposes mutable catalog method %q", iface.name, method)
			}
		}
	}
}
