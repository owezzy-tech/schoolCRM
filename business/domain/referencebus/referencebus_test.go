package referencebus_test

import (
	"context"
	"testing"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/referencebus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
)

type stubStore struct {
	countyFilter    referencebus.QueryFilter
	subCountyFilter referencebus.QueryFilter
	countyPage      page.Page
	subCountyPage   page.Page
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
