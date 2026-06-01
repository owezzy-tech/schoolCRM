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
