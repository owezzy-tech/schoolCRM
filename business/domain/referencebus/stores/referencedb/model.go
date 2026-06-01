package referencedb

import (
	"time"

	"github.com/owezzy/schoolCRM/business/domain/referencebus"
)

type countyDB struct {
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type subCountyDB struct {
	Code        string    `db:"code"`
	CountyCode  string    `db:"county_code"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toBusCounty(db countyDB) (referencebus.County, error) {
	bus := referencebus.County{
		Code:        db.Code,
		Name:        db.Name,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusCounties(dbCounties []countyDB) ([]referencebus.County, error) {
	bus := make([]referencebus.County, len(dbCounties))

	for i, dbCounty := range dbCounties {
		var err error
		bus[i], err = toBusCounty(dbCounty)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusSubCounty(db subCountyDB) (referencebus.SubCounty, error) {
	bus := referencebus.SubCounty{
		Code:        db.Code,
		CountyCode:  db.CountyCode,
		Name:        db.Name,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusSubCounties(dbSubCounties []subCountyDB) ([]referencebus.SubCounty, error) {
	bus := make([]referencebus.SubCounty, len(dbSubCounties))

	for i, dbSubCounty := range dbSubCounties {
		var err error
		bus[i], err = toBusSubCounty(dbSubCounty)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
