// Package referencebus provides business access to reference data.
package referencebus

import (
	"context"
	"fmt"

	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
)

// Storer interface declares the behavior this package needs to retrieve
// reference data.
type Storer interface {
	QueryCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]County, error)
	CountCounties(ctx context.Context, filter QueryFilter) (int, error)
	QuerySubCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]SubCounty, error)
	CountSubCounties(ctx context.Context, filter QueryFilter) (int, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	QueryCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]County, error)
	CountCounties(ctx context.Context, filter QueryFilter) (int, error)
	QuerySubCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]SubCounty, error)
	CountSubCounties(ctx context.Context, filter QueryFilter) (int, error)
	QueryCountyByCode(ctx context.Context, code string) (County, error)
	QuerySubCountiesByCounty(ctx context.Context, countyCode string) ([]SubCounty, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for reference data access.
type Business struct {
	storer Storer
}

// NewBusiness constructs a reference data business API for use.
func NewBusiness(storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		storer: storer,
	})

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// QueryCounties retrieves a list of existing counties.
func (b *Business) QueryCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]County, error) {
	counties, err := b.storer.QueryCounties(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query counties: %w", err)
	}

	return counties, nil
}

// CountCounties returns the total number of counties.
func (b *Business) CountCounties(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountCounties(ctx, filter)
}

// QuerySubCounties retrieves a list of existing sub-counties.
func (b *Business) QuerySubCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]SubCounty, error) {
	subCounties, err := b.storer.QuerySubCounties(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query sub-counties: %w", err)
	}

	return subCounties, nil
}

// CountSubCounties returns the total number of sub-counties.
func (b *Business) CountSubCounties(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountSubCounties(ctx, filter)
}

// QueryCountyByCode retrieves a single county by code.
func (b *Business) QueryCountyByCode(ctx context.Context, code string) (County, error) {
	filter := QueryFilter{Code: &code}

	counties, err := b.QueryCounties(ctx, filter, DefaultOrderBy, page.MustParse("1", "1"))
	if err != nil {
		return County{}, fmt.Errorf("query county by code: %w", err)
	}

	if len(counties) == 0 {
		return County{}, fmt.Errorf("county %q not found", code)
	}

	return counties[0], nil
}

// QuerySubCountiesByCounty retrieves sub-counties for a county code.
func (b *Business) QuerySubCountiesByCounty(ctx context.Context, countyCode string) ([]SubCounty, error) {
	filter := QueryFilter{CountyCode: &countyCode}

	subCounties, err := b.QuerySubCounties(ctx, filter, DefaultOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return nil, fmt.Errorf("query sub-counties by county: %w", err)
	}

	return subCounties, nil
}
