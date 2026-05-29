package admissionsapp

import "github.com/owezzy/schoolCRM/business/domain/admissionsbus"

var programOrderByFields = map[string]string{
	"program_id": admissionsbus.OrderByProgramID,
	"name":       admissionsbus.OrderByProgramName,
	"code":       admissionsbus.OrderByProgramCode,
	"active":     admissionsbus.OrderByProgramActive,
}

var constituentOrderByFields = map[string]string{
	"constituent_id":  admissionsbus.OrderByConstituentID,
	"last_name":       admissionsbus.OrderByConstituentLastName,
	"primary_email":   admissionsbus.OrderByConstituentPrimaryEmail,
	"lifecycle_stage": admissionsbus.OrderByConstituentLifecycleStage,
}

var academicTermOrderByFields = map[string]string{
	"academic_term_id": admissionsbus.OrderByAcademicTermID,
	"name":             admissionsbus.OrderByAcademicTermName,
	"code":             admissionsbus.OrderByAcademicTermCode,
	"start_date":       admissionsbus.OrderByAcademicTermStartDate,
	"active":           admissionsbus.OrderByAcademicTermActive,
}
