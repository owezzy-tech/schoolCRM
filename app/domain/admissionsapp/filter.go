package admissionsapp

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

type programQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ExternalSISID string
	Code          string
	Active        string
}

type staffProfileQueryParams struct {
	Page    string
	Rows    string
	OrderBy string
	ID      string
	UserID  string
	Role    string
	Active  string
}

type applicantProfileQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	UserID        string
	ConstituentID string
	Active        string
}

type leadScoreRuleQueryParams struct {
	Page    string
	Rows    string
	OrderBy string
	ID      string
	Active  string
}

type leadScoreQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ConstituentID string
	Band          string
	MinScore      string
}

type constituentQueryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	PrimaryEmail    string
	ExternalSISID   string
	LifecycleStage  string
	DuplicateStatus string
}

type inquiryQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ConstituentID string
	PrimaryEmail  string
	Source        string
	Status        string
}

type academicTermQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ExternalSISID string
	Code          string
	Active        string
}

type duplicateReviewQueryParams struct {
	Page                   string
	Rows                   string
	OrderBy                string
	ID                     string
	SourceConstituentID    string
	CandidateConstituentID string
	MatchType              string
	Status                 string
}

type applicationQueryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	ConstituentID   string
	ProgramID       string
	AcademicTermID  string
	ApplicationType string
	Status          string
	ActiveOnly      string
}

type applicationFormTemplateQueryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	ProgramID       string
	AcademicTermID  string
	ApplicationType string
	Active          string
	Version         string
}

type applicationTransitionQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ApplicationID string
	ActorID       string
	FromStatus    string
	ToStatus      string
}

type checklistItemQueryParams struct {
	Page          string
	Rows          string
	OrderBy       string
	ID            string
	ApplicationID string
	Status        string
	Required      string
}

type documentQueryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	ApplicationID   string
	ChecklistItemID string
	Status          string
	UploadedByID    string
	ReviewerID      string
}

func parseProgramQueryParams(r *http.Request) programQueryParams {
	values := r.URL.Query()

	return programQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("program_id"),
		ExternalSISID: values.Get("external_sis_id"),
		Code:          values.Get("code"),
		Active:        values.Get("active"),
	}
}

func parseStaffProfileQueryParams(r *http.Request) staffProfileQueryParams {
	values := r.URL.Query()

	return staffProfileQueryParams{
		Page:    values.Get("page"),
		Rows:    values.Get("rows"),
		OrderBy: values.Get("orderBy"),
		ID:      values.Get("staff_profile_id"),
		UserID:  values.Get("user_id"),
		Role:    values.Get("role"),
		Active:  values.Get("active"),
	}
}

func parseApplicantProfileQueryParams(r *http.Request) applicantProfileQueryParams {
	values := r.URL.Query()

	return applicantProfileQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("applicant_profile_id"),
		UserID:        values.Get("user_id"),
		ConstituentID: values.Get("constituent_id"),
		Active:        values.Get("active"),
	}
}

func parseLeadScoreRuleQueryParams(r *http.Request) leadScoreRuleQueryParams {
	values := r.URL.Query()

	return leadScoreRuleQueryParams{
		Page:    values.Get("page"),
		Rows:    values.Get("rows"),
		OrderBy: values.Get("orderBy"),
		ID:      values.Get("lead_score_rule_id"),
		Active:  values.Get("active"),
	}
}

func parseLeadScoreQueryParams(r *http.Request) leadScoreQueryParams {
	values := r.URL.Query()

	return leadScoreQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("lead_score_id"),
		ConstituentID: values.Get("constituent_id"),
		Band:          values.Get("band"),
		MinScore:      values.Get("min_score"),
	}
}

func parseConstituentQueryParams(r *http.Request) constituentQueryParams {
	values := r.URL.Query()

	return constituentQueryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("constituent_id"),
		PrimaryEmail:    values.Get("primary_email"),
		ExternalSISID:   values.Get("external_sis_id"),
		LifecycleStage:  values.Get("lifecycle_stage"),
		DuplicateStatus: values.Get("duplicate_status"),
	}
}

func parseInquiryQueryParams(r *http.Request) inquiryQueryParams {
	values := r.URL.Query()

	return inquiryQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("inquiry_id"),
		ConstituentID: values.Get("constituent_id"),
		PrimaryEmail:  values.Get("primary_email"),
		Source:        values.Get("source"),
		Status:        values.Get("status"),
	}
}

func parseAcademicTermQueryParams(r *http.Request) academicTermQueryParams {
	values := r.URL.Query()

	return academicTermQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("academic_term_id"),
		ExternalSISID: values.Get("external_sis_id"),
		Code:          values.Get("code"),
		Active:        values.Get("active"),
	}
}

func parseDuplicateReviewQueryParams(r *http.Request) duplicateReviewQueryParams {
	values := r.URL.Query()

	return duplicateReviewQueryParams{
		Page:                   values.Get("page"),
		Rows:                   values.Get("rows"),
		OrderBy:                values.Get("orderBy"),
		ID:                     values.Get("duplicate_review_id"),
		SourceConstituentID:    values.Get("source_constituent_id"),
		CandidateConstituentID: values.Get("candidate_constituent_id"),
		MatchType:              values.Get("match_type"),
		Status:                 values.Get("status"),
	}
}

func parseApplicationQueryParams(r *http.Request) applicationQueryParams {
	values := r.URL.Query()

	return applicationQueryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("application_id"),
		ConstituentID:   values.Get("constituent_id"),
		ProgramID:       values.Get("program_id"),
		AcademicTermID:  values.Get("academic_term_id"),
		ApplicationType: values.Get("application_type"),
		Status:          values.Get("status"),
		ActiveOnly:      values.Get("active_only"),
	}
}

func parseApplicationFormTemplateQueryParams(r *http.Request) applicationFormTemplateQueryParams {
	values := r.URL.Query()

	return applicationFormTemplateQueryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("form_template_id"),
		ProgramID:       values.Get("program_id"),
		AcademicTermID:  values.Get("academic_term_id"),
		ApplicationType: values.Get("application_type"),
		Active:          values.Get("active"),
		Version:         values.Get("version"),
	}
}

func parseApplicationTransitionQueryParams(r *http.Request) applicationTransitionQueryParams {
	values := r.URL.Query()

	return applicationTransitionQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("application_transition_id"),
		ApplicationID: values.Get("application_id"),
		ActorID:       values.Get("actor_id"),
		FromStatus:    values.Get("from_status"),
		ToStatus:      values.Get("to_status"),
	}
}

func parseChecklistItemQueryParams(r *http.Request) checklistItemQueryParams {
	values := r.URL.Query()

	return checklistItemQueryParams{
		Page:          values.Get("page"),
		Rows:          values.Get("rows"),
		OrderBy:       values.Get("orderBy"),
		ID:            values.Get("checklist_item_id"),
		ApplicationID: values.Get("application_id"),
		Status:        values.Get("status"),
		Required:      values.Get("required"),
	}
}

func parseDocumentQueryParams(r *http.Request) documentQueryParams {
	values := r.URL.Query()

	return documentQueryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("document_id"),
		ApplicationID:   values.Get("application_id"),
		ChecklistItemID: values.Get("checklist_item_id"),
		Status:          values.Get("status"),
		UploadedByID:    values.Get("uploaded_by_id"),
		ReviewerID:      values.Get("reviewer_id"),
	}
}

func parseProgramFilter(qp programQueryParams) (admissionsbus.ProgramQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ProgramQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("program_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.Code != "" {
		filter.Code = &qp.Code
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if fieldErrors != nil {
		return admissionsbus.ProgramQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseStaffProfileFilter(qp staffProfileQueryParams) (admissionsbus.StaffProfileQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.StaffProfileQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("staff_profile_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.UserID != "" {
		userID, err := uuid.Parse(qp.UserID)
		if err != nil {
			fieldErrors.Add("user_id", err)
		} else {
			filter.UserID = &userID
		}
	}

	if qp.Role != "" {
		role := admissionsbus.AdmissionsRole(qp.Role)
		filter.Role = &role
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.StaffProfileQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseApplicantProfileFilter(qp applicantProfileQueryParams) (admissionsbus.ApplicantProfileQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ApplicantProfileQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("applicant_profile_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.UserID != "" {
		userID, err := uuid.Parse(qp.UserID)
		if err != nil {
			fieldErrors.Add("user_id", err)
		} else {
			filter.UserID = &userID
		}
	}

	if qp.ConstituentID != "" {
		constituentID, err := uuid.Parse(qp.ConstituentID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ConstituentID = &constituentID
		}
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.ApplicantProfileQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseLeadScoreRuleFilter(qp leadScoreRuleQueryParams) (admissionsbus.LeadScoreRuleQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.LeadScoreRuleQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("lead_score_rule_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.LeadScoreRuleQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseLeadScoreFilter(qp leadScoreQueryParams) (admissionsbus.LeadScoreQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.LeadScoreQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("lead_score_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ConstituentID != "" {
		constituentID, err := uuid.Parse(qp.ConstituentID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ConstituentID = &constituentID
		}
	}

	if qp.Band != "" {
		band := admissionsbus.LeadScoreBand(qp.Band)
		filter.Band = &band
	}

	if qp.MinScore != "" {
		minScore, err := strconv.Atoi(qp.MinScore)
		if err != nil {
			fieldErrors.Add("min_score", err)
		} else {
			filter.MinScore = &minScore
		}
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.LeadScoreQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseConstituentFilter(qp constituentQueryParams) (admissionsbus.ConstituentQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ConstituentQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.PrimaryEmail != "" {
		email, err := mail.ParseAddress(qp.PrimaryEmail)
		if err != nil {
			fieldErrors.Add("primary_email", err)
		} else {
			filter.PrimaryEmail = email
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.LifecycleStage != "" {
		stage := admissionsbus.LifecycleStage(qp.LifecycleStage)
		filter.LifecycleStage = &stage
	}

	if qp.DuplicateStatus != "" {
		status := admissionsbus.DuplicateStatus(qp.DuplicateStatus)
		filter.DuplicateStatus = &status
	}

	if fieldErrors != nil {
		return admissionsbus.ConstituentQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseInquiryFilter(qp inquiryQueryParams) (admissionsbus.InquiryQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.InquiryQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("inquiry_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ConstituentID != "" {
		constituentID, err := uuid.Parse(qp.ConstituentID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ConstituentID = &constituentID
		}
	}

	if qp.PrimaryEmail != "" {
		email, err := mail.ParseAddress(qp.PrimaryEmail)
		if err != nil {
			fieldErrors.Add("primary_email", err)
		} else {
			filter.PrimaryEmail = email
		}
	}

	if qp.Source != "" {
		filter.Source = &qp.Source
	}

	if qp.Status != "" {
		status := admissionsbus.InquiryStatus(qp.Status)
		filter.Status = &status
	}

	if len(fieldErrors) > 0 {
		return admissionsbus.InquiryQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseAcademicTermFilter(qp academicTermQueryParams) (admissionsbus.AcademicTermQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.AcademicTermQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("academic_term_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ExternalSISID != "" {
		filter.ExternalSISID = &qp.ExternalSISID
	}

	if qp.Code != "" {
		filter.Code = &qp.Code
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if fieldErrors != nil {
		return admissionsbus.AcademicTermQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseDuplicateReviewFilter(qp duplicateReviewQueryParams) (admissionsbus.DuplicateReviewQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.DuplicateReviewQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("duplicate_review_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.SourceConstituentID != "" {
		id, err := uuid.Parse(qp.SourceConstituentID)
		if err != nil {
			fieldErrors.Add("source_constituent_id", err)
		} else {
			filter.SourceConstituentID = &id
		}
	}

	if qp.CandidateConstituentID != "" {
		id, err := uuid.Parse(qp.CandidateConstituentID)
		if err != nil {
			fieldErrors.Add("candidate_constituent_id", err)
		} else {
			filter.CandidateConstituentID = &id
		}
	}

	if qp.MatchType != "" {
		matchType := admissionsbus.DuplicateReviewMatchType(qp.MatchType)
		filter.MatchType = &matchType
	}

	if qp.Status != "" {
		status := admissionsbus.DuplicateReviewStatus(qp.Status)
		filter.Status = &status
	}

	if fieldErrors != nil {
		return admissionsbus.DuplicateReviewQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseApplicationFilter(qp applicationQueryParams) (admissionsbus.ApplicationQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ApplicationQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("application_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ConstituentID != "" {
		id, err := uuid.Parse(qp.ConstituentID)
		if err != nil {
			fieldErrors.Add("constituent_id", err)
		} else {
			filter.ConstituentID = &id
		}
	}

	if qp.ProgramID != "" {
		id, err := uuid.Parse(qp.ProgramID)
		if err != nil {
			fieldErrors.Add("program_id", err)
		} else {
			filter.ProgramID = &id
		}
	}

	if qp.AcademicTermID != "" {
		id, err := uuid.Parse(qp.AcademicTermID)
		if err != nil {
			fieldErrors.Add("academic_term_id", err)
		} else {
			filter.AcademicTermID = &id
		}
	}

	if qp.ApplicationType != "" {
		applicationType := admissionsbus.ApplicationType(qp.ApplicationType)
		filter.ApplicationType = &applicationType
	}

	if qp.Status != "" {
		status := admissionsbus.ApplicationStatus(qp.Status)
		filter.Status = &status
	}

	if qp.ActiveOnly != "" {
		activeOnly, err := strconv.ParseBool(qp.ActiveOnly)
		if err != nil {
			fieldErrors.Add("active_only", err)
		} else {
			filter.ActiveOnly = &activeOnly
		}
	}

	if fieldErrors != nil {
		return admissionsbus.ApplicationQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseApplicationFormTemplateFilter(qp applicationFormTemplateQueryParams) (admissionsbus.ApplicationFormTemplateQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ApplicationFormTemplateQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("form_template_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ProgramID != "" {
		id, err := uuid.Parse(qp.ProgramID)
		if err != nil {
			fieldErrors.Add("program_id", err)
		} else {
			filter.ProgramID = &id
		}
	}

	if qp.AcademicTermID != "" {
		id, err := uuid.Parse(qp.AcademicTermID)
		if err != nil {
			fieldErrors.Add("academic_term_id", err)
		} else {
			filter.AcademicTermID = &id
		}
	}

	if qp.ApplicationType != "" {
		applicationType := admissionsbus.ApplicationType(qp.ApplicationType)
		filter.ApplicationType = &applicationType
	}

	if qp.Active != "" {
		active, err := strconv.ParseBool(qp.Active)
		if err != nil {
			fieldErrors.Add("active", err)
		} else {
			filter.Active = &active
		}
	}

	if qp.Version != "" {
		version, err := strconv.Atoi(qp.Version)
		if err != nil {
			fieldErrors.Add("version", err)
		} else {
			filter.Version = &version
		}
	}

	if fieldErrors != nil {
		return admissionsbus.ApplicationFormTemplateQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseApplicationTransitionFilter(qp applicationTransitionQueryParams) (admissionsbus.ApplicationTransitionQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ApplicationTransitionQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("application_transition_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ApplicationID != "" {
		id, err := uuid.Parse(qp.ApplicationID)
		if err != nil {
			fieldErrors.Add("application_id", err)
		} else {
			filter.ApplicationID = &id
		}
	}

	if qp.ActorID != "" {
		id, err := uuid.Parse(qp.ActorID)
		if err != nil {
			fieldErrors.Add("actor_id", err)
		} else {
			filter.ActorID = &id
		}
	}

	if qp.FromStatus != "" {
		status := admissionsbus.ApplicationStatus(qp.FromStatus)
		filter.FromStatus = &status
	}

	if qp.ToStatus != "" {
		status := admissionsbus.ApplicationStatus(qp.ToStatus)
		filter.ToStatus = &status
	}

	if fieldErrors != nil {
		return admissionsbus.ApplicationTransitionQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseChecklistItemFilter(qp checklistItemQueryParams) (admissionsbus.ChecklistItemQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.ChecklistItemQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("checklist_item_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ApplicationID != "" {
		id, err := uuid.Parse(qp.ApplicationID)
		if err != nil {
			fieldErrors.Add("application_id", err)
		} else {
			filter.ApplicationID = &id
		}
	}

	if qp.Status != "" {
		status := admissionsbus.DocumentStatus(qp.Status)
		filter.Status = &status
	}

	if qp.Required != "" {
		required, err := strconv.ParseBool(qp.Required)
		if err != nil {
			fieldErrors.Add("required", err)
		} else {
			filter.Required = &required
		}
	}

	if fieldErrors != nil {
		return admissionsbus.ChecklistItemQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

func parseDocumentFilter(qp documentQueryParams) (admissionsbus.DocumentQueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter admissionsbus.DocumentQueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("document_id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.ApplicationID != "" {
		id, err := uuid.Parse(qp.ApplicationID)
		if err != nil {
			fieldErrors.Add("application_id", err)
		} else {
			filter.ApplicationID = &id
		}
	}

	if qp.ChecklistItemID != "" {
		id, err := uuid.Parse(qp.ChecklistItemID)
		if err != nil {
			fieldErrors.Add("checklist_item_id", err)
		} else {
			filter.ChecklistItemID = &id
		}
	}

	if qp.Status != "" {
		status := admissionsbus.DocumentStatus(qp.Status)
		filter.Status = &status
	}

	if qp.UploadedByID != "" {
		id, err := uuid.Parse(qp.UploadedByID)
		if err != nil {
			fieldErrors.Add("uploaded_by_id", err)
		} else {
			filter.UploadedByID = &id
		}
	}

	if qp.ReviewerID != "" {
		id, err := uuid.Parse(qp.ReviewerID)
		if err != nil {
			fieldErrors.Add("reviewer_id", err)
		} else {
			filter.ReviewerID = &id
		}
	}

	if fieldErrors != nil {
		return admissionsbus.DocumentQueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
