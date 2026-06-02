package referencebus

import (
	"time"
)

// County represents an individual Kenya county reference row.
type County struct {
	Code        string
	Name        string
	DateCreated time.Time
	DateUpdated time.Time
}

// SubCounty represents an individual Kenya sub-county reference row.
type SubCounty struct {
	Code        string
	CountyCode  string
	Name        string
	DateCreated time.Time
	DateUpdated time.Time
}

// Ward represents a Kenya ward reference row.
type Ward struct {
	Code          string
	CountyCode    string
	SubCountyCode string
	Name          string
	DateCreated   time.Time
	DateUpdated   time.Time
}

// University represents an accredited Kenya university reference row.
type University struct {
	Code            string
	Name            string
	InstitutionType string
	DateCreated     time.Time
	DateUpdated     time.Time
}

// Cluster represents a KUCCPS subject cluster reference row.
type Cluster struct {
	Code        string
	Name        string
	Description string
	DateCreated time.Time
	DateUpdated time.Time
}

// KNQFLevel represents a Kenya National Qualifications Framework level.
type KNQFLevel struct {
	Code          string
	Level         int
	Name          string
	Descriptor    string
	Qualification string
	DateCreated   time.Time
	DateUpdated   time.Time
}

// Programme represents a Kenya admissions programme reference row.
type Programme struct {
	Code           string
	UniversityCode string
	ClusterCode    string
	KNQFLevelCode  string
	Name           string
	AwardType      string
	DateCreated    time.Time
	DateUpdated    time.Time
}
