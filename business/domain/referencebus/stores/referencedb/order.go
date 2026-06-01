package referencedb

import (
	"fmt"

	"github.com/owezzy/schoolCRM/business/domain/referencebus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
)

var countyOrderByFields = map[string]string{
	referencebus.OrderByCode: "code",
	referencebus.OrderByName: "name",
}

var subCountyOrderByFields = map[string]string{
	referencebus.OrderByCode:       "code",
	referencebus.OrderByCountyCode: "county_code",
	referencebus.OrderByName:       "name",
}

func countyOrderByClause(orderBy order.By) (string, error) {
	by, exists := countyOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func subCountyOrderByClause(orderBy order.By) (string, error) {
	by, exists := subCountyOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
