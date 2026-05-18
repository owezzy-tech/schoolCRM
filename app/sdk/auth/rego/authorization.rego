package ardan.rego

import rego.v1

role_super_admin := "SUPER_ADMIN"

role_school_admin := "SCHOOL_ADMIN"

role_teacher := "TEACHER"

role_student := "STUDENT"

role_parent := "PARENT"

role_all := {role_super_admin, role_school_admin, role_teacher, role_student, role_parent}

role_admin_compat := {role_super_admin, role_school_admin}

role_user_compat := {role_teacher, role_student, role_parent}

default rule_any := false

rule_any if {
	claim_roles := {role | some role in input.Roles}
	input_roles := role_all & claim_roles
	count(input_roles) > 0
}

default rule_admin_only := false

rule_admin_only if {
	claim_roles := {role | some role in input.Roles}
	input_admin := role_admin_compat & claim_roles
	count(input_admin) > 0
}

default rule_user_only := false

rule_user_only if {
	claim_roles := {role | some role in input.Roles}
	input_user := role_user_compat & claim_roles
	count(input_user) > 0
}

default rule_admin_or_subject := false

rule_admin_or_subject if {
	claim_roles := {role | some role in input.Roles}
	input_admin := role_admin_compat & claim_roles
	count(input_admin) > 0
} else if {
	claim_roles := {role | some role in input.Roles}
	input_user := role_user_compat & claim_roles
	count(input_user) > 0
	input.UserID == input.Subject
}

default rule_super_admin_only := false

rule_super_admin_only if {
	claim_roles := {role | some role in input.Roles}
	input_super_admin := {role_super_admin} & claim_roles
	count(input_super_admin) > 0
}

default rule_school_admin_only := false

rule_school_admin_only if {
	claim_roles := {role | some role in input.Roles}
	input_school_admin := {role_school_admin} & claim_roles
	count(input_school_admin) > 0
}

default rule_teacher_only := false

rule_teacher_only if {
	claim_roles := {role | some role in input.Roles}
	input_teacher := {role_teacher} & claim_roles
	count(input_teacher) > 0
}

default rule_student_only := false

rule_student_only if {
	claim_roles := {role | some role in input.Roles}
	input_student := {role_student} & claim_roles
	count(input_student) > 0
}

default rule_parent_only := false

rule_parent_only if {
	claim_roles := {role | some role in input.Roles}
	input_parent := {role_parent} & claim_roles
	count(input_parent) > 0
}

default rule_rag_ingest := false

rule_rag_ingest if {
	claim_roles := {role | some role in input.Roles}
	allowed := {role_super_admin, role_school_admin, role_teacher}
	input_allowed := allowed & claim_roles
	count(input_allowed) > 0
}

default rule_rag_delete := false

rule_rag_delete if {
	claim_roles := {role | some role in input.Roles}
	allowed := {role_super_admin, role_school_admin}
	input_allowed := allowed & claim_roles
	count(input_allowed) > 0
}

default rule_rag_query := false

rule_rag_query if {
	claim_roles := {role | some role in input.Roles}
	input_roles := role_all & claim_roles
	count(input_roles) > 0
}
