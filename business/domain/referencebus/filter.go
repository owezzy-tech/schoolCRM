package referencebus

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	Code           *string
	CountyCode     *string
	SubCountyCode  *string
	UniversityCode *string
	ClusterCode    *string
	KNQFLevelCode  *string
	Name           *string
}
