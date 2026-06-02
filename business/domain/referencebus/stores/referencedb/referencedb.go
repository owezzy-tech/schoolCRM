// Package referencedb provides access to reference data storage.
package referencedb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/referencebus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Store manages the set of APIs for reference data database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// QueryCounties retrieves a list of existing counties from the database.
func (s *Store) QueryCounties(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.County, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		name,
		date_created,
		date_updated
	FROM
		counties`

	buf := bytes.NewBufferString(q)
	s.applyCountyFilter(filter, data, buf)

	orderByClause, err := countyOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbCounties []countyDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbCounties); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	counties, err := toBusCounties(dbCounties)
	if err != nil {
		return nil, err
	}

	return counties, nil
}

// CountCounties returns the total number of counties in the DB.
func (s *Store) CountCounties(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		counties`

	buf := bytes.NewBufferString(q)
	s.applyCountyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QuerySubCounties retrieves a list of existing sub-counties from the database.
func (s *Store) QuerySubCounties(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.SubCounty, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		county_code,
		name,
		date_created,
		date_updated
	FROM
		sub_counties`

	buf := bytes.NewBufferString(q)
	s.applySubCountyFilter(filter, data, buf)

	orderByClause, err := subCountyOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbSubCounties []subCountyDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbSubCounties); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	subCounties, err := toBusSubCounties(dbSubCounties)
	if err != nil {
		return nil, err
	}

	return subCounties, nil
}

// CountSubCounties returns the total number of sub-counties in the DB.
func (s *Store) CountSubCounties(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		sub_counties`

	buf := bytes.NewBufferString(q)
	s.applySubCountyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryWards retrieves a list of existing wards from the database.
func (s *Store) QueryWards(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Ward, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		county_code,
		sub_county_code,
		name,
		date_created,
		date_updated
	FROM
		wards`

	buf := bytes.NewBufferString(q)
	s.applyWardFilter(filter, data, buf)

	orderByClause, err := wardOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbWards []wardDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbWards); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	wards, err := toBusWards(dbWards)
	if err != nil {
		return nil, err
	}

	return wards, nil
}

// CountWards returns the total number of wards in the DB.
func (s *Store) CountWards(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		wards`

	buf := bytes.NewBufferString(q)
	s.applyWardFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryUniversities retrieves a list of existing universities from the database.
func (s *Store) QueryUniversities(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.University, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		name,
		institution_type,
		date_created,
		date_updated
	FROM
		universities`

	buf := bytes.NewBufferString(q)
	s.applyUniversityFilter(filter, data, buf)

	orderByClause, err := universityOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbUniversities []universityDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbUniversities); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	universities, err := toBusUniversities(dbUniversities)
	if err != nil {
		return nil, err
	}

	return universities, nil
}

// CountUniversities returns the total number of universities in the DB.
func (s *Store) CountUniversities(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		universities`

	buf := bytes.NewBufferString(q)
	s.applyUniversityFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryClusters retrieves a list of existing KUCCPS clusters from the database.
func (s *Store) QueryClusters(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Cluster, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		name,
		description,
		date_created,
		date_updated
	FROM
		programme_clusters`

	buf := bytes.NewBufferString(q)
	s.applyClusterFilter(filter, data, buf)

	orderByClause, err := clusterOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbClusters []clusterDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbClusters); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	clusters, err := toBusClusters(dbClusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

// CountClusters returns the total number of clusters in the DB.
func (s *Store) CountClusters(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		programme_clusters`

	buf := bytes.NewBufferString(q)
	s.applyClusterFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryKNQFLevels retrieves a list of existing KNQF levels from the database.
func (s *Store) QueryKNQFLevels(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.KNQFLevel, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		level,
		name,
		descriptor,
		qualification,
		date_created,
		date_updated
	FROM
		knqf_levels`

	buf := bytes.NewBufferString(q)
	s.applyKNQFLevelFilter(filter, data, buf)

	orderByClause, err := knqfLevelOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbLevels []knqfLevelDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbLevels); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	levels, err := toBusKNQFLevels(dbLevels)
	if err != nil {
		return nil, err
	}

	return levels, nil
}

// CountKNQFLevels returns the total number of KNQF levels in the DB.
func (s *Store) CountKNQFLevels(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		knqf_levels`

	buf := bytes.NewBufferString(q)
	s.applyKNQFLevelFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryProgrammes retrieves a list of existing programmes from the database.
func (s *Store) QueryProgrammes(ctx context.Context, filter referencebus.QueryFilter, orderBy order.By, page page.Page) ([]referencebus.Programme, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		code,
		university_code,
		cluster_code,
		knqf_level_code,
		name,
		award_type,
		date_created,
		date_updated
	FROM
		programmes`

	buf := bytes.NewBufferString(q)
	s.applyProgrammeFilter(filter, data, buf)

	orderByClause, err := programmeOrderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbProgrammes []programmeDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbProgrammes); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	programmes, err := toBusProgrammes(dbProgrammes)
	if err != nil {
		return nil, err
	}

	return programmes, nil
}

// CountProgrammes returns the total number of programmes in the DB.
func (s *Store) CountProgrammes(ctx context.Context, filter referencebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		programmes`

	buf := bytes.NewBufferString(q)
	s.applyProgrammeFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}
