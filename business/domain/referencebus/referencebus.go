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
	QueryWards(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Ward, error)
	CountWards(ctx context.Context, filter QueryFilter) (int, error)
	QueryUniversities(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]University, error)
	CountUniversities(ctx context.Context, filter QueryFilter) (int, error)
	QueryClusters(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Cluster, error)
	CountClusters(ctx context.Context, filter QueryFilter) (int, error)
	QueryKNQFLevels(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]KNQFLevel, error)
	CountKNQFLevels(ctx context.Context, filter QueryFilter) (int, error)
	QueryProgrammes(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Programme, error)
	CountProgrammes(ctx context.Context, filter QueryFilter) (int, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	QueryCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]County, error)
	CountCounties(ctx context.Context, filter QueryFilter) (int, error)
	QuerySubCounties(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]SubCounty, error)
	CountSubCounties(ctx context.Context, filter QueryFilter) (int, error)
	QueryWards(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Ward, error)
	CountWards(ctx context.Context, filter QueryFilter) (int, error)
	QueryUniversities(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]University, error)
	CountUniversities(ctx context.Context, filter QueryFilter) (int, error)
	QueryClusters(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Cluster, error)
	CountClusters(ctx context.Context, filter QueryFilter) (int, error)
	QueryKNQFLevels(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]KNQFLevel, error)
	CountKNQFLevels(ctx context.Context, filter QueryFilter) (int, error)
	QueryProgrammes(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Programme, error)
	CountProgrammes(ctx context.Context, filter QueryFilter) (int, error)
	QueryCountyByCode(ctx context.Context, code string) (County, error)
	QuerySubCountiesByCounty(ctx context.Context, countyCode string) ([]SubCounty, error)
	QueryWardsBySubCounty(ctx context.Context, subCountyCode string) ([]Ward, error)
	QueryProgrammesByCluster(ctx context.Context, clusterCode string) ([]Programme, error)
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

// QueryWards retrieves a list of existing wards.
func (b *Business) QueryWards(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Ward, error) {
	wards, err := b.storer.QueryWards(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query wards: %w", err)
	}

	return wards, nil
}

// CountWards returns the total number of wards.
func (b *Business) CountWards(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountWards(ctx, filter)
}

// QueryUniversities retrieves a list of existing universities.
func (b *Business) QueryUniversities(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]University, error) {
	universities, err := b.storer.QueryUniversities(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query universities: %w", err)
	}

	return universities, nil
}

// CountUniversities returns the total number of universities.
func (b *Business) CountUniversities(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountUniversities(ctx, filter)
}

// QueryClusters retrieves a list of existing KUCCPS clusters.
func (b *Business) QueryClusters(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Cluster, error) {
	clusters, err := b.storer.QueryClusters(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query clusters: %w", err)
	}

	return clusters, nil
}

// CountClusters returns the total number of clusters.
func (b *Business) CountClusters(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountClusters(ctx, filter)
}

// QueryKNQFLevels retrieves a list of existing KNQF levels.
func (b *Business) QueryKNQFLevels(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]KNQFLevel, error) {
	levels, err := b.storer.QueryKNQFLevels(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query KNQF levels: %w", err)
	}

	return levels, nil
}

// CountKNQFLevels returns the total number of KNQF levels.
func (b *Business) CountKNQFLevels(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountKNQFLevels(ctx, filter)
}

// QueryProgrammes retrieves a list of existing programmes.
func (b *Business) QueryProgrammes(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Programme, error) {
	programmes, err := b.storer.QueryProgrammes(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query programmes: %w", err)
	}

	return programmes, nil
}

// CountProgrammes returns the total number of programmes.
func (b *Business) CountProgrammes(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.CountProgrammes(ctx, filter)
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

// QueryWardsBySubCounty retrieves wards for a sub-county code.
func (b *Business) QueryWardsBySubCounty(ctx context.Context, subCountyCode string) ([]Ward, error) {
	filter := QueryFilter{SubCountyCode: &subCountyCode}

	wards, err := b.QueryWards(ctx, filter, DefaultOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return nil, fmt.Errorf("query wards by sub-county: %w", err)
	}

	return wards, nil
}

// QueryProgrammesByCluster retrieves programmes for a KUCCPS cluster code.
func (b *Business) QueryProgrammesByCluster(ctx context.Context, clusterCode string) ([]Programme, error) {
	filter := QueryFilter{ClusterCode: &clusterCode}

	programmes, err := b.QueryProgrammes(ctx, filter, DefaultOrderBy, page.MustParse("1", "100"))
	if err != nil {
		return nil, fmt.Errorf("query programmes by cluster: %w", err)
	}

	return programmes, nil
}
