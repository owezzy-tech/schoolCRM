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
