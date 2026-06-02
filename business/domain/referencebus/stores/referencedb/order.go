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

var wardOrderByFields = map[string]string{
	referencebus.OrderByCode:          "code",
	referencebus.OrderByCountyCode:    "county_code",
	referencebus.OrderBySubCountyCode: "sub_county_code",
	referencebus.OrderByName:          "name",
}

var universityOrderByFields = map[string]string{
	referencebus.OrderByCode: "code",
	referencebus.OrderByName: "name",
}

var clusterOrderByFields = map[string]string{
	referencebus.OrderByCode: "code",
	referencebus.OrderByName: "name",
}

var knqfLevelOrderByFields = map[string]string{
	referencebus.OrderByCode:      "code",
	referencebus.OrderByKNQFLevel: "level",
	referencebus.OrderByName:      "name",
}

var programmeOrderByFields = map[string]string{
	referencebus.OrderByCode:           "code",
	referencebus.OrderByUniversityCode: "university_code",
	referencebus.OrderByClusterCode:    "cluster_code",
	referencebus.OrderByKNQFLevelCode:  "knqf_level_code",
	referencebus.OrderByName:           "name",
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

func wardOrderByClause(orderBy order.By) (string, error) {
	by, exists := wardOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func universityOrderByClause(orderBy order.By) (string, error) {
	by, exists := universityOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func clusterOrderByClause(orderBy order.By) (string, error) {
	by, exists := clusterOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func knqfLevelOrderByClause(orderBy order.By) (string, error) {
	by, exists := knqfLevelOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func programmeOrderByClause(orderBy order.By) (string, error) {
	by, exists := programmeOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
