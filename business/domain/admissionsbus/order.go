package admissionsbus

import "github.com/owezzy/schoolCRM/business/sdk/order"

// DefaultProgramOrderBy represents the default way we sort programs.
var DefaultProgramOrderBy = order.NewBy(OrderByProgramID, order.ASC)

// DefaultStaffProfileOrderBy represents the default way we sort staff profiles.
var DefaultStaffProfileOrderBy = order.NewBy(OrderByStaffProfileID, order.ASC)

// DefaultApplicantProfileOrderBy represents the default way we sort applicant profiles.
var DefaultApplicantProfileOrderBy = order.NewBy(OrderByApplicantProfileID, order.ASC)

// DefaultInquiryOrderBy represents the default way we sort inquiries.
var DefaultInquiryOrderBy = order.NewBy(OrderByInquiryDateCreated, order.DESC)

// DefaultLeadScoreRuleOrderBy represents the default way we sort lead score rules.
var DefaultLeadScoreRuleOrderBy = order.NewBy(OrderByLeadScoreRulePriority, order.ASC)

// DefaultLeadScoreOrderBy represents the default way we sort lead scores.
var DefaultLeadScoreOrderBy = order.NewBy(OrderByLeadScoreTotalScore, order.DESC)

// DefaultConstituentOrderBy represents the default way we sort constituents.
var DefaultConstituentOrderBy = order.NewBy(OrderByConstituentID, order.ASC)

// DefaultAcademicTermOrderBy represents the default way we sort academic terms.
var DefaultAcademicTermOrderBy = order.NewBy(OrderByAcademicTermStartDate, order.ASC)

// DefaultDuplicateReviewOrderBy represents the default way we sort duplicate reviews.
var DefaultDuplicateReviewOrderBy = order.NewBy(OrderByDuplicateReviewDateCreated, order.ASC)

// DefaultApplicationOrderBy represents the default way we sort applications.
var DefaultApplicationOrderBy = order.NewBy(OrderByApplicationDateCreated, order.ASC)

// DefaultApplicationFormTemplateOrderBy represents the default way we sort application form templates.
var DefaultApplicationFormTemplateOrderBy = order.NewBy(OrderByApplicationFormTemplatePriority, order.ASC)

// DefaultApplicationTransitionOrderBy represents the default way we sort application transitions.
var DefaultApplicationTransitionOrderBy = order.NewBy(OrderByApplicationTransitionDateCreated, order.ASC)

// DefaultChecklistItemOrderBy represents the default way we sort checklist items.
var DefaultChecklistItemOrderBy = order.NewBy(OrderByChecklistItemDisplayOrder, order.ASC)

// DefaultDocumentOrderBy represents the default way we sort documents.
var DefaultDocumentOrderBy = order.NewBy(OrderByDocumentUploadedAt, order.DESC)

// Set of fields that programs can be ordered by.
const (
	OrderByStaffProfileID          = "a"
	OrderByStaffProfileUser        = "b"
	OrderByStaffProfileDateCreated = "c"
)

// Set of fields that applicant profiles can be ordered by.
const (
	OrderByApplicantProfileID          = "a"
	OrderByApplicantProfileUser        = "b"
	OrderByApplicantProfileConstituent = "c"
	OrderByApplicantProfileDateCreated = "d"
)

// Set of fields that inquiries can be ordered by.
const (
	OrderByInquiryID          = "a"
	OrderByInquiryEmail       = "b"
	OrderByInquirySource      = "c"
	OrderByInquiryStatus      = "d"
	OrderByInquiryDateCreated = "e"
)

// Set of fields that lead score rules can be ordered by.
const (
	OrderByLeadScoreRuleID          = "a"
	OrderByLeadScoreRuleName        = "b"
	OrderByLeadScoreRulePriority    = "c"
	OrderByLeadScoreRuleDateCreated = "d"
)

// Set of fields that lead scores can be ordered by.
const (
	OrderByLeadScoreID             = "a"
	OrderByLeadScoreTotalScore     = "b"
	OrderByLeadScoreBand           = "c"
	OrderByLeadScoreRecalculatedAt = "d"
)

// Set of fields that constituents can be ordered by.
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

// Set of fields that application form templates can be ordered by.
const (
	OrderByApplicationFormTemplateID          = "a"
	OrderByApplicationFormTemplateName        = "b"
	OrderByApplicationFormTemplateType        = "c"
	OrderByApplicationFormTemplateVersion     = "d"
	OrderByApplicationFormTemplatePriority    = "e"
	OrderByApplicationFormTemplateDateCreated = "f"
)

// Set of fields that application transitions can be ordered by.
const (
	OrderByApplicationTransitionID          = "a"
	OrderByApplicationTransitionApplication = "b"
	OrderByApplicationTransitionActor       = "c"
	OrderByApplicationTransitionDateCreated = "d"
)

// Set of fields that checklist items can be ordered by.
const (
	OrderByChecklistItemID           = "a"
	OrderByChecklistItemApplication  = "b"
	OrderByChecklistItemStatus       = "c"
	OrderByChecklistItemDisplayOrder = "d"
	OrderByChecklistItemDateCreated  = "e"
)

// Set of fields that documents can be ordered by.
const (
	OrderByDocumentID            = "a"
	OrderByDocumentApplication   = "b"
	OrderByDocumentChecklistItem = "c"
	OrderByDocumentStatus        = "d"
	OrderByDocumentUploadedAt    = "e"
	OrderByDocumentReviewedAt    = "f"
)
