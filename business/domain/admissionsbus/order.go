package admissionsbus

import "github.com/owezzy/schoolCRM/business/sdk/order"

// DefaultProgramOrderBy represents the default way we sort programs.
var DefaultProgramOrderBy = order.NewBy(OrderByProgramID, order.ASC)

// DefaultAcademicTermOrderBy represents the default way we sort academic terms.
var DefaultAcademicTermOrderBy = order.NewBy(OrderByAcademicTermStartDate, order.ASC)

// Set of fields that programs can be ordered by.
const (
	OrderByProgramID     = "a"
	OrderByProgramName   = "b"
	OrderByProgramCode   = "c"
	OrderByProgramActive = "d"
)

// Set of fields that academic terms can be ordered by.
const (
	OrderByAcademicTermID        = "a"
	OrderByAcademicTermName      = "b"
	OrderByAcademicTermCode      = "c"
	OrderByAcademicTermStartDate = "d"
	OrderByAcademicTermActive    = "e"
)
