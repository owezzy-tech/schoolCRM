package referencebus

import "github.com/owezzy/schoolCRM/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByCode, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByCode       = "a"
	OrderByCountyCode = "b"
	OrderByName       = "c"
)
