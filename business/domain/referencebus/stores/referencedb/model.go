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

type wardDB struct {
	Code          string    `db:"code"`
	CountyCode    string    `db:"county_code"`
	SubCountyCode string    `db:"sub_county_code"`
	Name          string    `db:"name"`
	DateCreated   time.Time `db:"date_created"`
	DateUpdated   time.Time `db:"date_updated"`
}

type universityDB struct {
	Code            string    `db:"code"`
	Name            string    `db:"name"`
	InstitutionType string    `db:"institution_type"`
	DateCreated     time.Time `db:"date_created"`
	DateUpdated     time.Time `db:"date_updated"`
}

type clusterDB struct {
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type knqfLevelDB struct {
	Code          string    `db:"code"`
	Level         int       `db:"level"`
	Name          string    `db:"name"`
	Descriptor    string    `db:"descriptor"`
	Qualification string    `db:"qualification"`
	DateCreated   time.Time `db:"date_created"`
	DateUpdated   time.Time `db:"date_updated"`
}

type programmeDB struct {
	Code           string    `db:"code"`
	UniversityCode string    `db:"university_code"`
	ClusterCode    string    `db:"cluster_code"`
	KNQFLevelCode  string    `db:"knqf_level_code"`
	Name           string    `db:"name"`
	AwardType      string    `db:"award_type"`
	DateCreated    time.Time `db:"date_created"`
	DateUpdated    time.Time `db:"date_updated"`
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

func toBusWard(db wardDB) (referencebus.Ward, error) {
	bus := referencebus.Ward{
		Code:          db.Code,
		CountyCode:    db.CountyCode,
		SubCountyCode: db.SubCountyCode,
		Name:          db.Name,
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusWards(dbWards []wardDB) ([]referencebus.Ward, error) {
	bus := make([]referencebus.Ward, len(dbWards))

	for i, dbWard := range dbWards {
		var err error
		bus[i], err = toBusWard(dbWard)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusUniversity(db universityDB) (referencebus.University, error) {
	bus := referencebus.University{
		Code:            db.Code,
		Name:            db.Name,
		InstitutionType: db.InstitutionType,
		DateCreated:     db.DateCreated.In(time.Local),
		DateUpdated:     db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusUniversities(dbUniversities []universityDB) ([]referencebus.University, error) {
	bus := make([]referencebus.University, len(dbUniversities))

	for i, dbUniversity := range dbUniversities {
		var err error
		bus[i], err = toBusUniversity(dbUniversity)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusCluster(db clusterDB) (referencebus.Cluster, error) {
	bus := referencebus.Cluster{
		Code:        db.Code,
		Name:        db.Name,
		Description: db.Description,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusClusters(dbClusters []clusterDB) ([]referencebus.Cluster, error) {
	bus := make([]referencebus.Cluster, len(dbClusters))

	for i, dbCluster := range dbClusters {
		var err error
		bus[i], err = toBusCluster(dbCluster)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusKNQFLevel(db knqfLevelDB) (referencebus.KNQFLevel, error) {
	bus := referencebus.KNQFLevel{
		Code:          db.Code,
		Level:         db.Level,
		Name:          db.Name,
		Descriptor:    db.Descriptor,
		Qualification: db.Qualification,
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusKNQFLevels(dbLevels []knqfLevelDB) ([]referencebus.KNQFLevel, error) {
	bus := make([]referencebus.KNQFLevel, len(dbLevels))

	for i, dbLevel := range dbLevels {
		var err error
		bus[i], err = toBusKNQFLevel(dbLevel)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusProgramme(db programmeDB) (referencebus.Programme, error) {
	bus := referencebus.Programme{
		Code:           db.Code,
		UniversityCode: db.UniversityCode,
		ClusterCode:    db.ClusterCode,
		KNQFLevelCode:  db.KNQFLevelCode,
		Name:           db.Name,
		AwardType:      db.AwardType,
		DateCreated:    db.DateCreated.In(time.Local),
		DateUpdated:    db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusProgrammes(dbProgrammes []programmeDB) ([]referencebus.Programme, error) {
	bus := make([]referencebus.Programme, len(dbProgrammes))

	for i, dbProgramme := range dbProgrammes {
		var err error
		bus[i], err = toBusProgramme(dbProgramme)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
