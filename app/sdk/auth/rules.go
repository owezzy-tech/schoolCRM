package auth

import (
	_ "embed"
)

// These are the current set of rules we have for auth.
const (
	RuleAuthenticate                 = "auth"
	RuleAny                          = "rule_any"
	RuleAdminOnly                    = "rule_admin_only"
	RuleUserOnly                     = "rule_user_only"
	RuleAdminOrSubject               = "rule_admin_or_subject"
	RuleSuperAdminOnly               = "rule_super_admin_only"
	RuleSchoolAdminOnly              = "rule_school_admin_only"
	RuleTeacherOnly                  = "rule_teacher_only"
	RuleStudentOnly                  = "rule_student_only"
	RuleParentOnly                   = "rule_parent_only"
	RuleRAGIngest                    = "rule_rag_ingest"
	RuleRAGDelete                    = "rule_rag_delete"
	RuleRAGQuery                     = "rule_rag_query"
	RuleAdmissionsRead               = "rule_admissions_read"
	RuleAdmissionsManageConstituents = "rule_admissions_manage_constituents"
	RuleAdmissionsManageApplications = "rule_admissions_manage_applications"
	RuleAdmissionsReviewApplications = "rule_admissions_review_applications"
	RuleAdmissionsResolveDuplicates  = "rule_admissions_resolve_duplicates"
	RuleAdmissionsManageReferences   = "rule_admissions_manage_references"
	RuleAdmissionsManageStaff        = "rule_admissions_manage_staff"
	RuleAdmissionsManageLeadScoring  = "rule_admissions_manage_lead_scoring"
)

// Package name of our rego code.
const (
	opaPackage string = "ardan.rego"
)

// Core OPA policies.
var (
	//go:embed rego/authentication.rego
	regoAuthentication string

	//go:embed rego/authorization.rego
	regoAuthorization string
)
