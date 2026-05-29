package admissionsbus

import "github.com/owezzy/schoolCRM/business/sdk/order"

// DefaultProgramOrderBy represents the default way we sort programs.
var DefaultProgramOrderBy = order.NewBy(OrderByProgramID, order.ASC)

// DefaultConstituentOrderBy represents the default way we sort constituents.
var DefaultConstituentOrderBy = order.NewBy(OrderByConstituentID, order.ASC)

// DefaultAcademicTermOrderBy represents the default way we sort academic terms.
var DefaultAcademicTermOrderBy = order.NewBy(OrderByAcademicTermStartDate, order.ASC)

// DefaultDuplicateReviewOrderBy represents the default way we sort duplicate reviews.
var DefaultDuplicateReviewOrderBy = order.NewBy(OrderByDuplicateReviewDateCreated, order.ASC)

// DefaultApplicationOrderBy represents the default way we sort applications.
var DefaultApplicationOrderBy = order.NewBy(OrderByApplicationDateCreated, order.ASC)

// DefaultApplicationTransitionOrderBy represents the default way we sort application transitions.
var DefaultApplicationTransitionOrderBy = order.NewBy(OrderByApplicationTransitionDateCreated, order.ASC)

// Set of fields that programs can be ordered by.
const (
	OrderByConstituentID             = "a"
	OrderByConstituentLastName       = "b"
	OrderByConstituentPrimaryEmail   = "c"
	OrderByConstituentLifecycleStage = "d"
)

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

// Set of fields that duplicate reviews can be ordered by.
const (
	OrderByDuplicateReviewID          = "a"
	OrderByDuplicateReviewStatus      = "b"
	OrderByDuplicateReviewMatchType   = "c"
	OrderByDuplicateReviewMatchScore  = "d"
	OrderByDuplicateReviewDateCreated = "e"
)

// Set of fields that applications can be ordered by.
const (
	OrderByApplicationID          = "a"
	OrderByApplicationStatus      = "b"
	OrderByApplicationType        = "c"
	OrderByApplicationSubmittedAt = "d"
	OrderByApplicationDateCreated = "e"
)

// Set of fields that application transitions can be ordered by.
const (
	OrderByApplicationTransitionID          = "a"
	OrderByApplicationTransitionApplication = "b"
	OrderByApplicationTransitionActor       = "c"
	OrderByApplicationTransitionDateCreated = "d"
)
