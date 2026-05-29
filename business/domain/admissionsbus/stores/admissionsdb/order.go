package admissionsdb

import (
	"fmt"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/order"
)

var programOrderByFields = map[string]string{
	admissionsbus.OrderByProgramID:     "program_id",
	admissionsbus.OrderByProgramName:   "name",
	admissionsbus.OrderByProgramCode:   "code",
	admissionsbus.OrderByProgramActive: "is_active",
}

var constituentOrderByFields = map[string]string{
	admissionsbus.OrderByConstituentID:             "constituent_id",
	admissionsbus.OrderByConstituentLastName:       "last_name",
	admissionsbus.OrderByConstituentPrimaryEmail:   "primary_email",
	admissionsbus.OrderByConstituentLifecycleStage: "lifecycle_stage",
}

var academicTermOrderByFields = map[string]string{
	admissionsbus.OrderByAcademicTermID:        "academic_term_id",
	admissionsbus.OrderByAcademicTermName:      "name",
	admissionsbus.OrderByAcademicTermCode:      "code",
	admissionsbus.OrderByAcademicTermStartDate: "start_date",
	admissionsbus.OrderByAcademicTermActive:    "is_active",
}

func programOrderByClause(orderBy order.By) (string, error) {
	by, exists := programOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func constituentOrderByClause(orderBy order.By) (string, error) {
	by, exists := constituentOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}

func academicTermOrderByClause(orderBy order.By) (string, error) {
	by, exists := academicTermOrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
