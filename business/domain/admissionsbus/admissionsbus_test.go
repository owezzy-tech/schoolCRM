package admissionsbus

import (
	"context"
	"errors"
	"net/mail"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/sdk/order"
	"github.com/owezzy/schoolCRM/business/sdk/page"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

func TestUpsertAcademicTermValidatesTermDateRange(t *testing.T) {
	t.Parallel()

	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
	now := time.Now()

	_, err := bus.UpsertAcademicTerm(context.Background(), UpsertAcademicTerm{
		ExternalSISID: "TERM-2026-FALL",
		Name:          "Fall 2026",
		Code:          "202609",
		StartDate:     now,
		EndDate:       now,
		Active:        true,
	})

	if !errors.Is(err, ErrInvalidTermDateRange) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidTermDateRange)
	}
}

func TestUpsertAcademicTermValidatesApplicationWindow(t *testing.T) {
	t.Parallel()

	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{})
	start := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	appStart := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	deadline := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	_, err := bus.UpsertAcademicTerm(context.Background(), UpsertAcademicTerm{
		ExternalSISID:        "TERM-2026-SPRING",
		Name:                 "Spring 2026",
		Code:                 "202601",
		StartDate:            start,
		EndDate:              end,
		ApplicationStartDate: &appStart,
		ApplicationDeadline:  &deadline,
		Active:               true,
	})

	if !errors.Is(err, ErrInvalidApplicationWindow) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidApplicationWindow)
	}
}

func TestCreateConstituentRequiresIdentityFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "applicant@example.com"}

	tests := []struct {
		name string
		nc   NewConstituent
		want error
	}{
		{
			name: "first name",
			nc: NewConstituent{
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrFirstNameRequired,
		},
		{
			name: "last name",
			nc: NewConstituent{
				FirstName:    "Ada",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrLastNameRequired,
		},
		{
			name: "date of birth",
			nc: NewConstituent{
				FirstName:    "Ada",
				LastName:     "Applicant",
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrDateOfBirthRequired,
		},
		{
			name: "primary phone",
			nc: NewConstituent{
				FirstName:    "Ada",
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
			},
			want: ErrPrimaryPhoneRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateConstituent(context.Background(), tt.nc)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateStaffProfileRequiresContextRoles(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		np   NewStaffProfile
		want error
	}{
		{
			name: "user id",
			np: NewStaffProfile{
				Roles: []AdmissionsRole{AdmissionsRoleApplicationReviewer},
			},
			want: ErrStaffProfileUserRequired,
		},
		{
			name: "role required",
			np: NewStaffProfile{
				UserID: uuid.New(),
			},
			want: ErrStaffProfileRoleRequired,
		},
		{
			name: "valid role",
			np: NewStaffProfile{
				UserID: uuid.New(),
				Roles:  []AdmissionsRole{"UNKNOWN"},
			},
			want: ErrInvalidAdmissionsRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateStaffProfile(context.Background(), tt.np)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateStaffProfileStoresAdmissionsContextRoles(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	userID := uuid.New()

	profile, err := bus.CreateStaffProfile(context.Background(), NewStaffProfile{
		UserID: userID,
		Roles: []AdmissionsRole{
			AdmissionsRoleApplicationReviewer,
			AdmissionsRoleReportViewer,
		},
		Active: true,
	})
	if err != nil {
		t.Fatalf("CreateStaffProfile returned error: %v", err)
	}

	if profile.UserID != userID {
		t.Fatalf("UserID = %s, want %s", profile.UserID, userID)
	}

	if len(store.staffProfiles) != 1 {
		t.Fatalf("stored profiles = %d, want 1", len(store.staffProfiles))
	}

	permissions := AdmissionsPermissionsToStrings(AdmissionsPermissionsForRoles(profile.Roles))
	if !containsString(permissions, AdmissionsPermissionReviewApplications.String()) {
		t.Fatalf("permissions = %v, want %s", permissions, AdmissionsPermissionReviewApplications)
	}

	if !containsString(permissions, AdmissionsPermissionResolveDuplicates.String()) {
		t.Fatalf("permissions = %v, want %s", permissions, AdmissionsPermissionResolveDuplicates)
	}
}

func TestCreateApplicantProfileLinksIdentityToConstituent(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{
			constituentID: {ID: constituentID, LifecycleStage: LifecycleStageApplicant},
		},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	userID := uuid.New()

	profile, err := bus.CreateApplicantProfile(context.Background(), NewApplicantProfile{
		UserID:        userID,
		ConstituentID: constituentID,
		Active:        true,
	})
	if err != nil {
		t.Fatalf("CreateApplicantProfile returned error: %v", err)
	}

	if profile.UserID != userID {
		t.Fatalf("UserID = %s, want %s", profile.UserID, userID)
	}

	if profile.ConstituentID != constituentID {
		t.Fatalf("ConstituentID = %s, want %s", profile.ConstituentID, constituentID)
	}

	if !profile.Active {
		t.Fatal("Active = false, want true")
	}

	stored, err := bus.QueryApplicantProfileByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("QueryApplicantProfileByUserID returned error: %v", err)
	}

	if stored.ID != profile.ID {
		t.Fatalf("stored profile ID = %s, want %s", stored.ID, profile.ID)
	}
}

func TestCreateSyncJobValidatesFrameworkFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nj   NewSyncJob
		want error
	}{
		{
			name: "name required",
			nj: NewSyncJob{
				Status:    SyncJobStatusQueued,
				Direction: SyncDirectionInbound,
			},
			want: ErrSyncJobNameRequired,
		},
		{
			name: "valid status",
			nj: NewSyncJob{
				Name:      "Nightly SIS reconciliation",
				Status:    SyncJobStatus("UNKNOWN"),
				Direction: SyncDirectionInbound,
			},
			want: ErrInvalidSyncJobStatus,
		},
		{
			name: "valid direction",
			nj: NewSyncJob{
				Name:      "Nightly SIS reconciliation",
				Status:    SyncJobStatusQueued,
				Direction: SyncDirection("UNKNOWN"),
			},
			want: ErrInvalidSyncDirection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateSyncJob(context.Background(), tt.nj)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestEnqueueSyncEventStoresApprovedRealtimeEvent(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	applicationID := uuid.New()

	event, err := bus.EnqueueSyncEvent(context.Background(), NewSyncEvent{
		EventType:    SyncEventTypeApplicationSubmission,
		Direction:    SyncDirectionOutbound,
		ResourceType: "application",
		ResourceID:   applicationID,
		PayloadHash:  "sha256:application-submission",
		AuditMessage: "Application submission queued for SIS.",
	})
	if err != nil {
		t.Fatalf("EnqueueSyncEvent returned error: %v", err)
	}

	if event.Status != SyncEventStatusQueued {
		t.Fatalf("Status = %s, want %s", event.Status, SyncEventStatusQueued)
	}

	if event.EventType != SyncEventTypeApplicationSubmission {
		t.Fatalf("EventType = %s, want %s", event.EventType, SyncEventTypeApplicationSubmission)
	}

	if len(store.syncEvents) != 1 {
		t.Fatalf("stored sync events = %d, want 1", len(store.syncEvents))
	}
}

func TestEnqueueSyncEventValidatesApprovedFieldSet(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		ne   NewSyncEvent
		want error
	}{
		{
			name: "event type",
			ne: NewSyncEvent{
				EventType:    SyncEventType("CUSTOM_FIELD_SYNC"),
				Direction:    SyncDirectionOutbound,
				ResourceType: "application",
				ResourceID:   uuid.New(),
				PayloadHash:  "sha256:custom",
			},
			want: ErrInvalidSyncEventType,
		},
		{
			name: "resource",
			ne: NewSyncEvent{
				EventType:   SyncEventTypeApplicationDecision,
				Direction:   SyncDirectionOutbound,
				ResourceID:  uuid.New(),
				PayloadHash: "sha256:decision",
			},
			want: ErrSyncEventResourceRequired,
		},
		{
			name: "payload hash",
			ne: NewSyncEvent{
				EventType:    SyncEventTypeDocumentStatus,
				Direction:    SyncDirectionOutbound,
				ResourceType: "document",
				ResourceID:   uuid.New(),
			},
			want: ErrSyncEventPayloadHashRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.EnqueueSyncEvent(context.Background(), tt.ne)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateApplicantProfileValidatesIdentityAndConstituent(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		np   NewApplicantProfile
		want error
	}{
		{
			name: "user id",
			np:   NewApplicantProfile{ConstituentID: uuid.New()},
			want: ErrApplicantProfileUserRequired,
		},
		{
			name: "constituent id",
			np:   NewApplicantProfile{UserID: uuid.New()},
			want: ErrConstituentIDRequired,
		},
		{
			name: "missing constituent",
			np:   NewApplicantProfile{UserID: uuid.New(), ConstituentID: uuid.New()},
			want: ErrConstituentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateApplicantProfile(context.Background(), tt.np)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateInquiryCreatesConstituentWithSourceAttribution(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	utmSource := "google"
	utmMedium := "cpc"
	utmCampaign := "fall-open-house"
	message := "Please send me admissions deadlines."
	email := mail.Address{Address: "new.inquiry@example.com"}

	inquiry, err := bus.CreateInquiry(context.Background(), NewInquiry{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
		Source:       "website",
		UTMSource:    &utmSource,
		UTMMedium:    &utmMedium,
		UTMCampaign:  &utmCampaign,
		Message:      &message,
	})
	if err != nil {
		t.Fatalf("CreateInquiry returned error: %v", err)
	}

	if inquiry.ConstituentID == uuid.Nil {
		t.Fatal("ConstituentID is nil, want linked constituent")
	}

	if inquiry.Status != InquiryStatusNew {
		t.Fatalf("Status = %s, want %s", inquiry.Status, InquiryStatusNew)
	}

	if inquiry.Source != "website" {
		t.Fatalf("Source = %q, want website", inquiry.Source)
	}

	if inquiry.UTMSource == nil || *inquiry.UTMSource != utmSource {
		t.Fatalf("UTMSource = %v, want %q", inquiry.UTMSource, utmSource)
	}

	constituent, err := bus.QueryConstituentByID(context.Background(), inquiry.ConstituentID)
	if err != nil {
		t.Fatalf("QueryConstituentByID returned error: %v", err)
	}

	if constituent.LifecycleStage != LifecycleStageInquiry {
		t.Fatalf("LifecycleStage = %s, want %s", constituent.LifecycleStage, LifecycleStageInquiry)
	}
}

func TestCreateInquiryLinksExistingConstituentBeforeCreate(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	email := mail.Address{Address: "existing.inquiry@example.com"}
	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{
			constituentID: {
				ID:             constituentID,
				FirstName:      "Existing",
				LastName:       "Applicant",
				DateOfBirth:    time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail:   email,
				PrimaryPhone:   "+15555550100",
				LifecycleStage: LifecycleStageProspect,
			},
		},
		constituentByEmail: map[string]Constituent{},
	}
	store.constituentByEmail[email.String()] = store.constituents[constituentID]
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	inquiry, err := bus.CreateInquiry(context.Background(), NewInquiry{
		FirstName:    "Existing",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
		Source:       "referral",
	})
	if err != nil {
		t.Fatalf("CreateInquiry returned error: %v", err)
	}

	if inquiry.ConstituentID != constituentID {
		t.Fatalf("ConstituentID = %s, want %s", inquiry.ConstituentID, constituentID)
	}

	if len(store.constituents) != 1 {
		t.Fatalf("constituents stored = %d, want 1", len(store.constituents))
	}
}

func TestCreateInquiryRequiresSourceAndIdentityFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "inquiry@example.com"}

	tests := []struct {
		name string
		ni   NewInquiry
		want error
	}{
		{
			name: "first name",
			ni: NewInquiry{
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
				Source:       "website",
			},
			want: ErrFirstNameRequired,
		},
		{
			name: "source",
			ni: NewInquiry{
				FirstName:    "Ada",
				LastName:     "Applicant",
				DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail: email,
				PrimaryPhone: "+15555550100",
			},
			want: ErrInquirySourceRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateInquiry(context.Background(), tt.ni)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateApplicationFormTemplateCreatesVersionedConfig(t *testing.T) {
	t.Parallel()

	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(uuid.New(), programID, termID)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	description := "Freshman application requirements"
	validation := `{"maxLength":120}`
	itemDescription := "Official high school transcript"

	template, err := bus.CreateApplicationFormTemplate(context.Background(), NewApplicationFormTemplate{
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
		Name:            " Freshman v1 ",
		Description:     &description,
		RequiredFields: []ApplicationFormField{
			{FieldName: "personal_statement", FieldType: "textarea", Required: true, DisplayOrder: 1, Validation: &validation},
		},
		ChecklistItems: []ApplicationChecklistTemplateItem{
			{ItemKey: "transcript", DocumentName: "High school transcript", Description: &itemDescription, Required: true, DisplayOrder: 1},
		},
		Active:   true,
		Priority: 10,
	})
	if err != nil {
		t.Fatalf("CreateApplicationFormTemplate returned error: %v", err)
	}

	if template.Version != 1 {
		t.Fatalf("Version = %d, want 1", template.Version)
	}

	if template.Name != "Freshman v1" {
		t.Fatalf("Name = %q, want trimmed Freshman v1", template.Name)
	}

	if len(template.RequiredFields) != 1 || template.RequiredFields[0].FieldName != "personal_statement" {
		t.Fatalf("RequiredFields = %#v, want personal_statement", template.RequiredFields)
	}

	if len(template.ChecklistItems) != 1 || template.ChecklistItems[0].ItemKey != "transcript" {
		t.Fatalf("ChecklistItems = %#v, want transcript", template.ChecklistItems)
	}
}

func TestUpdateApplicationFormTemplateIncrementsVersion(t *testing.T) {
	t.Parallel()

	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(uuid.New(), programID, termID)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	template, err := bus.CreateApplicationFormTemplate(context.Background(), NewApplicationFormTemplate{
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeTransfer,
		Name:            "Transfer v1",
		RequiredFields: []ApplicationFormField{
			{FieldName: "prior_college", FieldType: "text", Required: true, DisplayOrder: 1},
		},
		Active:   true,
		Priority: 20,
	})
	if err != nil {
		t.Fatalf("CreateApplicationFormTemplate returned error: %v", err)
	}

	updated, err := bus.UpdateApplicationFormTemplate(context.Background(), template, NewApplicationFormTemplate{
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeTransfer,
		Name:            "Transfer v2",
		RequiredFields: []ApplicationFormField{
			{FieldName: "prior_college", FieldType: "text", Required: true, DisplayOrder: 1},
			{FieldName: "college_gpa", FieldType: "number", Required: true, DisplayOrder: 2},
		},
		Active:   true,
		Priority: 5,
	})
	if err != nil {
		t.Fatalf("UpdateApplicationFormTemplate returned error: %v", err)
	}

	if updated.Version != 2 {
		t.Fatalf("Version = %d, want 2", updated.Version)
	}

	if len(updated.RequiredFields) != 2 {
		t.Fatalf("RequiredFields count = %d, want 2", len(updated.RequiredFields))
	}
}

func TestCreateApplicationFormTemplateValidatesConfig(t *testing.T) {
	t.Parallel()

	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(uuid.New(), programID, termID)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	tests := []struct {
		name string
		nt   NewApplicationFormTemplate
		want error
	}{
		{
			name: "name required",
			nt: NewApplicationFormTemplate{
				ProgramID:       programID,
				AcademicTermID:  termID,
				ApplicationType: ApplicationTypeFreshman,
				RequiredFields: []ApplicationFormField{
					{FieldName: "personal_statement", FieldType: "textarea", Required: true, DisplayOrder: 1},
				},
			},
			want: ErrFormTemplateNameRequired,
		},
		{
			name: "fields required",
			nt: NewApplicationFormTemplate{
				ProgramID:       programID,
				AcademicTermID:  termID,
				ApplicationType: ApplicationTypeFreshman,
				Name:            "Freshman",
			},
			want: ErrFormTemplateFieldsRequired,
		},
		{
			name: "checklist invalid",
			nt: NewApplicationFormTemplate{
				ProgramID:       programID,
				AcademicTermID:  termID,
				ApplicationType: ApplicationTypeFreshman,
				Name:            "Freshman",
				RequiredFields: []ApplicationFormField{
					{FieldName: "personal_statement", FieldType: "textarea", Required: true, DisplayOrder: 1},
				},
				ChecklistItems: []ApplicationChecklistTemplateItem{{ItemKey: "", DocumentName: "Transcript", Required: true}},
			},
			want: ErrFormTemplateChecklistInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateApplicationFormTemplate(context.Background(), tt.nt)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestAdmissionsPermissionsForRolesKeepsActionsSeparateFromMenus(t *testing.T) {
	t.Parallel()

	permissions := AdmissionsPermissionsToStrings(AdmissionsPermissionsForRoles([]AdmissionsRole{AdmissionsRoleRecruiter}))

	if !containsString(permissions, AdmissionsPermissionManageConstituents.String()) {
		t.Fatalf("permissions = %v, want %s", permissions, AdmissionsPermissionManageConstituents)
	}

	if containsString(permissions, AdmissionsPermissionReviewApplications.String()) {
		t.Fatalf("permissions = %v, did not expect %s", permissions, AdmissionsPermissionReviewApplications)
	}
}

func TestCreateLeadScoreRuleValidatesExplainableCriteria(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nr   NewLeadScoreRule
		want error
	}{
		{
			name: "name",
			nr: NewLeadScoreRule{
				Criteria: []LeadScoreCriterion{{Field: LeadScoreCriterionFieldLifecycleStage, Operator: LeadScoreCriterionOperatorEquals, Values: []string{LifecycleStageApplicant.String()}}},
			},
			want: ErrLeadScoreRuleNameRequired,
		},
		{
			name: "criteria",
			nr: NewLeadScoreRule{
				Name: "Applicant stage",
			},
			want: ErrLeadScoreCriteriaRequired,
		},
		{
			name: "points",
			nr: NewLeadScoreRule{
				Name:     "Applicant stage",
				Criteria: []LeadScoreCriterion{{Field: LeadScoreCriterionFieldLifecycleStage, Operator: LeadScoreCriterionOperatorEquals, Values: []string{LifecycleStageApplicant.String()}}},
				Points:   -1,
			},
			want: ErrInvalidLeadScorePoints,
		},
		{
			name: "criterion",
			nr: NewLeadScoreRule{
				Name:     "Applicant stage",
				Criteria: []LeadScoreCriterion{{Field: "unknown", Operator: LeadScoreCriterionOperatorEquals, Values: []string{LifecycleStageApplicant.String()}}},
			},
			want: ErrInvalidLeadScoreCriterion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateLeadScoreRule(context.Background(), tt.nr)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestRecalculateLeadScoreExplainsMatchedRulesAndBand(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{
			constituentID: {ID: constituentID, LifecycleStage: LifecycleStageApplicant, DuplicateStatus: DuplicateStatusActive},
		},
		applications: []Application{{
			ID:              uuid.New(),
			ConstituentID:   constituentID,
			ProgramID:       programID,
			AcademicTermID:  termID,
			ApplicationType: ApplicationTypeTransfer,
			Status:          ApplicationStatusSubmitted,
		}},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateLeadScoreRule(context.Background(), NewLeadScoreRule{
		Name:     "Applicant stage",
		Criteria: []LeadScoreCriterion{{Field: LeadScoreCriterionFieldLifecycleStage, Operator: LeadScoreCriterionOperatorEquals, Values: []string{LifecycleStageApplicant.String()}}},
		Points:   30,
		Active:   true,
		Priority: 1,
	})
	if err != nil {
		t.Fatalf("CreateLeadScoreRule returned error: %v", err)
	}

	_, err = bus.CreateLeadScoreRule(context.Background(), NewLeadScoreRule{
		Name:     "Submitted application",
		Criteria: []LeadScoreCriterion{{Field: LeadScoreCriterionFieldApplicationStatus, Operator: LeadScoreCriterionOperatorIn, Values: []string{ApplicationStatusSubmitted.String(), ApplicationStatusReadyForReview.String()}}},
		Points:   40,
		Active:   true,
		Priority: 2,
	})
	if err != nil {
		t.Fatalf("CreateLeadScoreRule returned error: %v", err)
	}

	score, err := bus.RecalculateLeadScoreForConstituent(context.Background(), constituentID)
	if err != nil {
		t.Fatalf("RecalculateLeadScoreForConstituent returned error: %v", err)
	}

	if score.TotalScore != 70 {
		t.Fatalf("TotalScore = %d, want 70", score.TotalScore)
	}

	if score.Band != LeadScoreBandHot {
		t.Fatalf("Band = %s, want %s", score.Band, LeadScoreBandHot)
	}

	if len(score.Breakdown) != 2 {
		t.Fatalf("Breakdown length = %d, want 2", len(score.Breakdown))
	}

	for _, result := range score.Breakdown {
		if !result.Matched {
			t.Fatalf("result %q did not match", result.Name)
		}
		if result.Reason == "" {
			t.Fatalf("result %q should explain the score", result.Name)
		}
	}

	if len(store.leadScores) != 1 {
		t.Fatalf("stored lead scores = %d, want 1", len(store.leadScores))
	}
}

func TestLeadScoreBandForTotal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		total int
		want  LeadScoreBand
	}{
		{total: 0, want: LeadScoreBandCold},
		{total: 26, want: LeadScoreBandWarm},
		{total: 51, want: LeadScoreBandHot},
		{total: 76, want: LeadScoreBandReadyToApply},
	}

	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			t.Parallel()

			if got := LeadScoreBandForTotal(tt.total); got != tt.want {
				t.Fatalf("LeadScoreBandForTotal(%d) = %s, want %s", tt.total, got, tt.want)
			}
		})
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}

func TestCreateConstituentAutoLinksExactEmailDuplicate(t *testing.T) {
	t.Parallel()

	email := mail.Address{Address: "applicant@example.com"}
	matchID := uuid.New()
	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{},
		constituentByEmail: map[string]Constituent{
			email.String(): {ID: matchID, PrimaryEmail: email},
		},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	created, err := bus.CreateConstituent(context.Background(), NewConstituent{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
	})
	if err != nil {
		t.Fatalf("CreateConstituent returned error: %v", err)
	}

	if created.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("DuplicateStatus = %s, want %s", created.DuplicateStatus, DuplicateStatusDuplicateOf)
	}

	if created.DuplicateOfID == nil || *created.DuplicateOfID != matchID {
		t.Fatalf("DuplicateOfID = %v, want %s", created.DuplicateOfID, matchID)
	}

	stored := store.constituents[created.ID]
	if stored.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("stored DuplicateStatus = %s, want %s", stored.DuplicateStatus, DuplicateStatusDuplicateOf)
	}
}

func TestCreateConstituentIgnoresMissingExactDuplicate(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "applicant@example.com"}

	_, err := bus.CreateConstituent(context.Background(), NewConstituent{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+15555550100",
	})
	if err != nil {
		t.Fatalf("CreateConstituent returned error: %v", err)
	}
}

func TestUpdateConstituentLifecycleTransitions(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	cst := Constituent{
		ID:              uuid.New(),
		LifecycleStage:  LifecycleStageProspect,
		DuplicateStatus: DuplicateStatusActive,
	}

	next := LifecycleStageInquiry
	updated, err := bus.UpdateConstituent(context.Background(), cst, UpdateConstituent{LifecycleStage: &next})
	if err != nil {
		t.Fatalf("UpdateConstituent returned error: %v", err)
	}

	if updated.LifecycleStage != LifecycleStageInquiry {
		t.Fatalf("LifecycleStage = %s, want %s", updated.LifecycleStage, LifecycleStageInquiry)
	}

	invalid := LifecycleStageEnrolled
	_, err = bus.UpdateConstituent(context.Background(), cst, UpdateConstituent{LifecycleStage: &invalid})
	if !errors.Is(err, ErrInvalidLifecycleChange) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidLifecycleChange)
	}
}

func TestCreateDuplicateReviewValidatesPair(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	id := uuid.New()

	_, err := bus.CreateDuplicateReview(context.Background(), NewDuplicateReview{
		SourceConstituentID:    id,
		CandidateConstituentID: id,
		MatchType:              DuplicateReviewMatchTypeExact,
		MatchScore:             100,
		MatchReason:            "same email",
	})

	if !errors.Is(err, ErrInvalidDuplicateReview) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidDuplicateReview)
	}
}

func TestCreateDuplicateReviewValidatesMatchData(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nr   NewDuplicateReview
		want error
	}{
		{
			name: "match type",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchType("UNKNOWN"),
				MatchScore:             100,
				MatchReason:            "same email",
			},
			want: ErrInvalidMatchType,
		},
		{
			name: "match score",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchTypeExact,
				MatchScore:             101,
				MatchReason:            "same email",
			},
			want: ErrInvalidMatchScore,
		},
		{
			name: "match reason",
			nr: NewDuplicateReview{
				SourceConstituentID:    uuid.New(),
				CandidateConstituentID: uuid.New(),
				MatchType:              DuplicateReviewMatchTypeExact,
				MatchScore:             100,
			},
			want: ErrMatchReasonRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateDuplicateReview(context.Background(), tt.nr)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestResolveDuplicateReviewLinksSourceConstituent(t *testing.T) {
	t.Parallel()

	store := &stubStore{
		constituents: map[uuid.UUID]Constituent{},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	sourceID := uuid.New()
	candidateID := uuid.New()
	store.constituents[sourceID] = Constituent{ID: sourceID, DuplicateStatus: DuplicateStatusActive}
	store.constituents[candidateID] = Constituent{ID: candidateID, DuplicateStatus: DuplicateStatusActive}

	review := DuplicateReview{
		ID:                     uuid.New(),
		SourceConstituentID:    sourceID,
		CandidateConstituentID: candidateID,
		MatchType:              DuplicateReviewMatchTypeExact,
		MatchScore:             100,
		MatchReason:            "same email",
		Status:                 DuplicateReviewStatusPending,
	}

	resolved, err := bus.ResolveDuplicateReview(context.Background(), review, ResolveDuplicateReview{
		Resolution: DuplicateReviewResolutionLink,
		ActorID:    uuid.New(),
	})
	if err != nil {
		t.Fatalf("ResolveDuplicateReview returned error: %v", err)
	}

	if resolved.Status != DuplicateReviewStatusLinked {
		t.Fatalf("Status = %s, want %s", resolved.Status, DuplicateReviewStatusLinked)
	}

	source := store.constituents[sourceID]
	if source.DuplicateStatus != DuplicateStatusDuplicateOf {
		t.Fatalf("DuplicateStatus = %s, want %s", source.DuplicateStatus, DuplicateStatusDuplicateOf)
	}

	if source.DuplicateOfID == nil || *source.DuplicateOfID != candidateID {
		t.Fatalf("DuplicateOfID = %v, want %s", source.DuplicateOfID, candidateID)
	}
}

func TestCreateApplicationPreventsDuplicateActiveApplication(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	store.applications = append(store.applications, Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
		Status:          ApplicationStatusDraft,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
	})

	if !errors.Is(err, ErrDuplicateApplication) {
		t.Fatalf("err = %v, want %v", err, ErrDuplicateApplication)
	}
}

func TestCreateApplicationAllowsClosedPriorApplication(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	store.applications = append(store.applications, Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeFreshman,
		Status:          ApplicationStatusWithdrawn,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	created, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeTransfer,
	})
	if err != nil {
		t.Fatalf("CreateApplication returned error: %v", err)
	}

	if created.Status != ApplicationStatusDraft {
		t.Fatalf("Status = %s, want %s", created.Status, ApplicationStatusDraft)
	}

	if len(store.applications) != 2 {
		t.Fatalf("application count = %d, want 2", len(store.applications))
	}
}

func TestCreateApplicationValidatesApplicationType(t *testing.T) {
	t.Parallel()

	store := newApplicationStubStore(uuid.New(), uuid.New(), uuid.New())
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   uuid.New(),
		ProgramID:       uuid.New(),
		AcademicTermID:  uuid.New(),
		ApplicationType: ApplicationType("UNKNOWN"),
	})

	if !errors.Is(err, ErrInvalidApplicationType) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidApplicationType)
	}
}

func TestTransitionApplicationStatusAllowsHappyPath(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	app := Application{
		ID:     uuid.New(),
		Status: ApplicationStatusDraft,
	}
	reason := "submitted by applicant"
	actorID := uuid.New()

	updated, transition, err := bus.TransitionApplicationStatus(context.Background(), app, NewApplicationTransition{
		ToStatus: ApplicationStatusSubmitted,
		ActorID:  actorID,
		Reason:   &reason,
	})
	if err != nil {
		t.Fatalf("TransitionApplicationStatus returned error: %v", err)
	}

	if updated.Status != ApplicationStatusSubmitted {
		t.Fatalf("Status = %s, want %s", updated.Status, ApplicationStatusSubmitted)
	}

	if updated.SubmittedAt == nil {
		t.Fatal("SubmittedAt is nil, want timestamp")
	}

	if transition.FromStatus != ApplicationStatusDraft || transition.ToStatus != ApplicationStatusSubmitted || transition.ActorID != actorID {
		t.Fatalf("transition = %+v, want from DRAFT to SUBMITTED by actor", transition)
	}

	if len(store.applicationTransitions) != 1 {
		t.Fatalf("transition count = %d, want 1", len(store.applicationTransitions))
	}
}

func TestTransitionApplicationStatusRejectsInvalidTransition(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	app := Application{
		ID:     uuid.New(),
		Status: ApplicationStatusDraft,
	}

	_, _, err := bus.TransitionApplicationStatus(context.Background(), app, NewApplicationTransition{
		ToStatus: ApplicationStatusAdmitted,
		ActorID:  uuid.New(),
	})

	if !errors.Is(err, ErrInvalidApplicationTransition) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidApplicationTransition)
	}
}

func TestTransitionApplicationStatusAllowsWithdrawalFromReviewStates(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	allowed := []ApplicationStatus{
		ApplicationStatusDraft,
		ApplicationStatusSubmitted,
		ApplicationStatusAwaitingDocuments,
		ApplicationStatusReadyForReview,
		ApplicationStatusInReview,
	}

	for _, from := range allowed {
		t.Run(from.String(), func(t *testing.T) {
			t.Parallel()

			_, _, err := bus.TransitionApplicationStatus(context.Background(), Application{ID: uuid.New(), Status: from}, NewApplicationTransition{
				ToStatus: ApplicationStatusWithdrawn,
				ActorID:  uuid.New(),
			})
			if err != nil {
				t.Fatalf("TransitionApplicationStatus returned error: %v", err)
			}
		})
	}
}

func TestTransitionApplicationStatusRequiresActor(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	_, _, err := bus.TransitionApplicationStatus(context.Background(), Application{ID: uuid.New(), Status: ApplicationStatusDraft}, NewApplicationTransition{
		ToStatus: ApplicationStatusSubmitted,
	})

	if !errors.Is(err, ErrApplicationActorRequired) {
		t.Fatalf("err = %v, want %v", err, ErrApplicationActorRequired)
	}
}

func TestCreateChecklistItemTiesItemToApplication(t *testing.T) {
	t.Parallel()

	applicationID := uuid.New()
	store := &stubStore{
		applications: []Application{{ID: applicationID, Status: ApplicationStatusAwaitingDocuments}},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	item, err := bus.CreateChecklistItem(context.Background(), NewChecklistItem{
		ApplicationID: applicationID,
		ItemKey:       "transcript",
		DocumentName:  "Official transcript",
		Required:      true,
		DisplayOrder:  1,
	})
	if err != nil {
		t.Fatalf("CreateChecklistItem returned error: %v", err)
	}

	if item.ApplicationID != applicationID {
		t.Fatalf("ApplicationID = %s, want %s", item.ApplicationID, applicationID)
	}

	if item.Status != DocumentStatusPendingReview {
		t.Fatalf("Status = %s, want %s", item.Status, DocumentStatusPendingReview)
	}

	if len(store.checklistItems) != 1 {
		t.Fatalf("stored checklist items = %d, want 1", len(store.checklistItems))
	}
}

func TestCreateDocumentTiesDocumentToApplicationChecklistItem(t *testing.T) {
	t.Parallel()

	applicationID := uuid.New()
	checklistItemID := uuid.New()
	uploadedByID := uuid.New()
	store := &stubStore{
		checklistItems: []ChecklistItem{{ID: checklistItemID, ApplicationID: applicationID, Status: DocumentStatusPendingReview}},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	document, err := bus.CreateDocument(context.Background(), NewDocument{
		ApplicationID:   applicationID,
		ChecklistItemID: checklistItemID,
		FileName:        "transcript.pdf",
		ContentType:     "application/pdf",
		SizeBytes:       2048,
		StorageKey:      "admissions/applications/app/transcript.pdf",
		UploadedByID:    uploadedByID,
	})
	if err != nil {
		t.Fatalf("CreateDocument returned error: %v", err)
	}

	if document.ApplicationID != applicationID || document.ChecklistItemID != checklistItemID {
		t.Fatalf("document linkage = (%s, %s), want (%s, %s)", document.ApplicationID, document.ChecklistItemID, applicationID, checklistItemID)
	}

	if document.Status != DocumentStatusPendingReview {
		t.Fatalf("Status = %s, want %s", document.Status, DocumentStatusPendingReview)
	}

	if document.StorageKey == "" {
		t.Fatal("StorageKey is empty")
	}

	if len(store.documents) != 1 {
		t.Fatalf("stored documents = %d, want 1", len(store.documents))
	}
}

func TestVerifyDocumentSupportsAcceptedRejectedWaived(t *testing.T) {
	t.Parallel()

	applicationID := uuid.New()
	checklistItemID := uuid.New()
	reviewerID := uuid.New()
	statuses := []DocumentStatus{DocumentStatusAccepted, DocumentStatusRejected, DocumentStatusWaived}

	for _, status := range statuses {
		status := status
		t.Run(status.String(), func(t *testing.T) {
			t.Parallel()

			documentID := uuid.New()
			store := &stubStore{
				checklistItems: []ChecklistItem{{ID: checklistItemID, ApplicationID: applicationID, Status: DocumentStatusPendingReview}},
				documents: []Document{{
					ID:              documentID,
					ApplicationID:   applicationID,
					ChecklistItemID: checklistItemID,
					Status:          DocumentStatusPendingReview,
				}},
			}
			bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

			updated, err := bus.VerifyDocument(context.Background(), store.documents[0], NewDocumentVerification{
				Status:     status,
				ReviewerID: reviewerID,
			})
			if err != nil {
				t.Fatalf("VerifyDocument returned error: %v", err)
			}

			if updated.Status != status {
				t.Fatalf("Status = %s, want %s", updated.Status, status)
			}

			if updated.ReviewerID == nil || *updated.ReviewerID != reviewerID {
				t.Fatalf("ReviewerID = %v, want %s", updated.ReviewerID, reviewerID)
			}

			if updated.ReviewedAt == nil {
				t.Fatal("ReviewedAt is nil")
			}

			if store.checklistItems[0].Status != status {
				t.Fatalf("checklist item status = %s, want %s", store.checklistItems[0].Status, status)
			}
		})
	}
}

func TestVerifyDocumentRejectsNonReviewStatus(t *testing.T) {
	t.Parallel()

	_, err := newTestBusiness().VerifyDocument(context.Background(), Document{ID: uuid.New()}, NewDocumentVerification{
		Status:     DocumentStatusPendingReview,
		ReviewerID: uuid.New(),
	})

	if !errors.Is(err, ErrDocumentStatusNotReviewable) {
		t.Fatalf("err = %v, want %v", err, ErrDocumentStatusNotReviewable)
	}
}

func newTestBusiness() ExtBusiness {
	return NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, &stubStore{
		constituents:       map[uuid.UUID]Constituent{},
		constituentByEmail: map[string]Constituent{},
	})
}

func newApplicationStubStore(constituentID uuid.UUID, programID uuid.UUID, termID uuid.UUID) *stubStore {
	return &stubStore{
		constituents: map[uuid.UUID]Constituent{
			constituentID: {ID: constituentID, LifecycleStage: LifecycleStageApplicant, DuplicateStatus: DuplicateStatusActive},
		},
		programs: map[uuid.UUID]Program{
			programID: {ID: programID, Active: true},
		},
		terms: map[uuid.UUID]AcademicTerm{
			termID: {ID: termID, Active: true},
		},
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

type stubStore struct {
	staffProfiles          []StaffProfile
	applicantProfiles      []ApplicantProfile
	inquiries              []Inquiry
	leadScoreRules         []LeadScoreRule
	leadScores             []LeadScore
	applicationTemplates   []ApplicationFormTemplate
	constituents           map[uuid.UUID]Constituent
	constituentByEmail     map[string]Constituent
	duplicateReviews       []DuplicateReview
	programs               map[uuid.UUID]Program
	terms                  map[uuid.UUID]AcademicTerm
	applications           []Application
	applicationTransitions []ApplicationTransition
	checklistItems         []ChecklistItem
	documents              []Document
	syncJobs               []SyncJob
	syncEvents             []SyncEvent
}

func (s *stubStore) NewWithTx(sqldb.CommitRollbacker) (Storer, error) {
	return s, nil
}

func (s *stubStore) Health(context.Context) (Health, error) {
	return Health{}, nil
}

func (s *stubStore) CreateStaffProfile(_ context.Context, profile StaffProfile) error {
	s.staffProfiles = append(s.staffProfiles, profile)
	return nil
}

func (s *stubStore) UpdateStaffProfile(_ context.Context, profile StaffProfile) error {
	for i, existing := range s.staffProfiles {
		if existing.ID == profile.ID {
			s.staffProfiles[i] = profile
			return nil
		}
	}
	s.staffProfiles = append(s.staffProfiles, profile)
	return nil
}

func (s *stubStore) QueryStaffProfiles(context.Context, StaffProfileQueryFilter, order.By, page.Page) ([]StaffProfile, error) {
	return s.staffProfiles, nil
}

func (s *stubStore) CountStaffProfiles(context.Context, StaffProfileQueryFilter) (int, error) {
	return len(s.staffProfiles), nil
}

func (s *stubStore) QueryStaffProfileByID(_ context.Context, profileID uuid.UUID) (StaffProfile, error) {
	for _, profile := range s.staffProfiles {
		if profile.ID == profileID {
			return profile, nil
		}
	}
	return StaffProfile{}, ErrStaffProfileNotFound
}

func (s *stubStore) QueryStaffProfileByUserID(_ context.Context, userID uuid.UUID) (StaffProfile, error) {
	for _, profile := range s.staffProfiles {
		if profile.UserID == userID {
			return profile, nil
		}
	}
	return StaffProfile{}, ErrStaffProfileNotFound
}

func (s *stubStore) CreateApplicantProfile(_ context.Context, profile ApplicantProfile) error {
	s.applicantProfiles = append(s.applicantProfiles, profile)
	return nil
}

func (s *stubStore) UpdateApplicantProfile(_ context.Context, profile ApplicantProfile) error {
	for i, existing := range s.applicantProfiles {
		if existing.ID == profile.ID {
			s.applicantProfiles[i] = profile
			return nil
		}
	}
	s.applicantProfiles = append(s.applicantProfiles, profile)
	return nil
}

func (s *stubStore) QueryApplicantProfiles(context.Context, ApplicantProfileQueryFilter, order.By, page.Page) ([]ApplicantProfile, error) {
	return s.applicantProfiles, nil
}

func (s *stubStore) CountApplicantProfiles(context.Context, ApplicantProfileQueryFilter) (int, error) {
	return len(s.applicantProfiles), nil
}

func (s *stubStore) QueryApplicantProfileByID(_ context.Context, profileID uuid.UUID) (ApplicantProfile, error) {
	for _, profile := range s.applicantProfiles {
		if profile.ID == profileID {
			return profile, nil
		}
	}
	return ApplicantProfile{}, ErrApplicantProfileNotFound
}

func (s *stubStore) QueryApplicantProfileByUserID(_ context.Context, userID uuid.UUID) (ApplicantProfile, error) {
	for _, profile := range s.applicantProfiles {
		if profile.UserID == userID {
			return profile, nil
		}
	}
	return ApplicantProfile{}, ErrApplicantProfileNotFound
}

func (s *stubStore) QueryApplicantProfileByConstituentID(_ context.Context, constituentID uuid.UUID) (ApplicantProfile, error) {
	for _, profile := range s.applicantProfiles {
		if profile.ConstituentID == constituentID {
			return profile, nil
		}
	}
	return ApplicantProfile{}, ErrApplicantProfileNotFound
}

func (s *stubStore) CreateInquiry(_ context.Context, inquiry Inquiry) error {
	s.inquiries = append(s.inquiries, inquiry)
	return nil
}

func (s *stubStore) UpdateInquiry(_ context.Context, inquiry Inquiry) error {
	for i, existing := range s.inquiries {
		if existing.ID == inquiry.ID {
			s.inquiries[i] = inquiry
			return nil
		}
	}
	s.inquiries = append(s.inquiries, inquiry)
	return nil
}

func (s *stubStore) QueryInquiries(context.Context, InquiryQueryFilter, order.By, page.Page) ([]Inquiry, error) {
	return s.inquiries, nil
}

func (s *stubStore) CountInquiries(context.Context, InquiryQueryFilter) (int, error) {
	return len(s.inquiries), nil
}

func (s *stubStore) QueryInquiryByID(_ context.Context, inquiryID uuid.UUID) (Inquiry, error) {
	for _, inquiry := range s.inquiries {
		if inquiry.ID == inquiryID {
			return inquiry, nil
		}
	}
	return Inquiry{}, ErrInquiryNotFound
}

func (s *stubStore) CreateLeadScoreRule(_ context.Context, rule LeadScoreRule) error {
	s.leadScoreRules = append(s.leadScoreRules, rule)
	return nil
}

func (s *stubStore) UpdateLeadScoreRule(_ context.Context, rule LeadScoreRule) error {
	for i, existing := range s.leadScoreRules {
		if existing.ID == rule.ID {
			s.leadScoreRules[i] = rule
			return nil
		}
	}
	s.leadScoreRules = append(s.leadScoreRules, rule)
	return nil
}

func (s *stubStore) QueryLeadScoreRules(_ context.Context, filter LeadScoreRuleQueryFilter, _ order.By, _ page.Page) ([]LeadScoreRule, error) {
	var rules []LeadScoreRule
	for _, rule := range s.leadScoreRules {
		if filter.Active != nil && rule.Active != *filter.Active {
			continue
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (s *stubStore) CountLeadScoreRules(context.Context, LeadScoreRuleQueryFilter) (int, error) {
	return len(s.leadScoreRules), nil
}

func (s *stubStore) QueryLeadScoreRuleByID(_ context.Context, ruleID uuid.UUID) (LeadScoreRule, error) {
	for _, rule := range s.leadScoreRules {
		if rule.ID == ruleID {
			return rule, nil
		}
	}
	return LeadScoreRule{}, ErrLeadScoreRuleNotFound
}

func (s *stubStore) UpsertLeadScore(_ context.Context, score LeadScore) error {
	for i, existing := range s.leadScores {
		if existing.ConstituentID == score.ConstituentID {
			s.leadScores[i] = score
			return nil
		}
	}
	s.leadScores = append(s.leadScores, score)
	return nil
}

func (s *stubStore) QueryLeadScores(context.Context, LeadScoreQueryFilter, order.By, page.Page) ([]LeadScore, error) {
	return s.leadScores, nil
}

func (s *stubStore) CountLeadScores(context.Context, LeadScoreQueryFilter) (int, error) {
	return len(s.leadScores), nil
}

func (s *stubStore) QueryLeadScoreByID(_ context.Context, scoreID uuid.UUID) (LeadScore, error) {
	for _, score := range s.leadScores {
		if score.ID == scoreID {
			return score, nil
		}
	}
	return LeadScore{}, ErrLeadScoreNotFound
}

func (s *stubStore) QueryLeadScoreByConstituentID(_ context.Context, constituentID uuid.UUID) (LeadScore, error) {
	for _, score := range s.leadScores {
		if score.ConstituentID == constituentID {
			return score, nil
		}
	}
	return LeadScore{}, ErrLeadScoreNotFound
}

func (s *stubStore) CreateConstituent(_ context.Context, cst Constituent) error {
	if s.constituents != nil {
		s.constituents[cst.ID] = cst
	}

	if s.constituentByEmail != nil {
		if _, exists := s.constituentByEmail[cst.PrimaryEmail.String()]; !exists {
			s.constituentByEmail[cst.PrimaryEmail.String()] = cst
		}
	}

	return nil
}

func (s *stubStore) UpdateConstituent(_ context.Context, cst Constituent) error {
	if s.constituents != nil {
		s.constituents[cst.ID] = cst
	}
	return nil
}

func (s *stubStore) QueryConstituents(context.Context, ConstituentQueryFilter, order.By, page.Page) ([]Constituent, error) {
	return nil, nil
}

func (s *stubStore) CountConstituents(context.Context, ConstituentQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryConstituentByID(_ context.Context, id uuid.UUID) (Constituent, error) {
	if s.constituents != nil {
		cst, exists := s.constituents[id]
		if exists {
			return cst, nil
		}
	}
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByPrimaryEmail(_ context.Context, email string) (Constituent, error) {
	if s.constituentByEmail != nil {
		cst, exists := s.constituentByEmail[email]
		if exists {
			return cst, nil
		}
	}

	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByExternalSISID(context.Context, string) (Constituent, error) {
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) UpsertProgram(context.Context, Program) error {
	return nil
}

func (s *stubStore) QueryPrograms(context.Context, ProgramQueryFilter, order.By, page.Page) ([]Program, error) {
	return nil, nil
}

func (s *stubStore) CountPrograms(context.Context, ProgramQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryProgramByID(_ context.Context, id uuid.UUID) (Program, error) {
	if s.programs != nil {
		program, exists := s.programs[id]
		if exists {
			return program, nil
		}
	}
	return Program{}, ErrProgramNotFound
}

func (s *stubStore) QueryProgramByExternalSISID(context.Context, string) (Program, error) {
	return Program{}, nil
}

func (s *stubStore) UpsertAcademicTerm(context.Context, AcademicTerm) error {
	return nil
}

func (s *stubStore) QueryAcademicTerms(context.Context, AcademicTermQueryFilter, order.By, page.Page) ([]AcademicTerm, error) {
	return nil, nil
}

func (s *stubStore) CountAcademicTerms(context.Context, AcademicTermQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryAcademicTermByID(_ context.Context, id uuid.UUID) (AcademicTerm, error) {
	if s.terms != nil {
		term, exists := s.terms[id]
		if exists {
			return term, nil
		}
	}
	return AcademicTerm{}, ErrAcademicTermNotFound
}

func (s *stubStore) QueryAcademicTermByExternalSISID(context.Context, string) (AcademicTerm, error) {
	return AcademicTerm{}, nil
}

func (s *stubStore) CreateDuplicateReview(_ context.Context, review DuplicateReview) error {
	s.duplicateReviews = append(s.duplicateReviews, review)
	return nil
}

func (s *stubStore) UpdateDuplicateReview(context.Context, DuplicateReview) error {
	return nil
}

func (s *stubStore) QueryDuplicateReviews(context.Context, DuplicateReviewQueryFilter, order.By, page.Page) ([]DuplicateReview, error) {
	return nil, nil
}

func (s *stubStore) CountDuplicateReviews(context.Context, DuplicateReviewQueryFilter) (int, error) {
	return 0, nil
}

func (s *stubStore) QueryDuplicateReviewByID(context.Context, uuid.UUID) (DuplicateReview, error) {
	return DuplicateReview{}, nil
}

func (s *stubStore) CreateApplication(_ context.Context, app Application) error {
	s.applications = append(s.applications, app)
	return nil
}

func (s *stubStore) UpdateApplication(_ context.Context, app Application) error {
	for i, existing := range s.applications {
		if existing.ID == app.ID {
			s.applications[i] = app
			return nil
		}
	}
	s.applications = append(s.applications, app)
	return nil
}

func (s *stubStore) QueryApplications(context.Context, ApplicationQueryFilter, order.By, page.Page) ([]Application, error) {
	return s.applications, nil
}

func (s *stubStore) CountApplications(context.Context, ApplicationQueryFilter) (int, error) {
	return len(s.applications), nil
}

func (s *stubStore) QueryApplicationByID(_ context.Context, id uuid.UUID) (Application, error) {
	for _, app := range s.applications {
		if app.ID == id {
			return app, nil
		}
	}
	return Application{}, ErrApplicationNotFound
}

func (s *stubStore) QueryActiveApplicationByTuple(_ context.Context, constituentID uuid.UUID, academicTermID uuid.UUID, programID uuid.UUID) (Application, error) {
	for _, app := range s.applications {
		if app.ConstituentID == constituentID && app.AcademicTermID == academicTermID && app.ProgramID == programID && isApplicationActive(app.Status) {
			return app, nil
		}
	}
	return Application{}, ErrApplicationNotFound
}

func (s *stubStore) CreateApplicationFormTemplate(_ context.Context, template ApplicationFormTemplate) error {
	s.applicationTemplates = append(s.applicationTemplates, template)
	return nil
}

func (s *stubStore) UpdateApplicationFormTemplate(_ context.Context, template ApplicationFormTemplate) error {
	for i, existing := range s.applicationTemplates {
		if existing.ID == template.ID {
			s.applicationTemplates[i] = template
			return nil
		}
	}
	s.applicationTemplates = append(s.applicationTemplates, template)
	return nil
}

func (s *stubStore) QueryApplicationFormTemplates(context.Context, ApplicationFormTemplateQueryFilter, order.By, page.Page) ([]ApplicationFormTemplate, error) {
	return s.applicationTemplates, nil
}

func (s *stubStore) CountApplicationFormTemplates(context.Context, ApplicationFormTemplateQueryFilter) (int, error) {
	return len(s.applicationTemplates), nil
}

func (s *stubStore) QueryApplicationFormTemplateByID(_ context.Context, templateID uuid.UUID) (ApplicationFormTemplate, error) {
	for _, template := range s.applicationTemplates {
		if template.ID == templateID {
			return template, nil
		}
	}
	return ApplicationFormTemplate{}, ErrFormTemplateNotFound
}

func (s *stubStore) CreateApplicationTransition(_ context.Context, transition ApplicationTransition) error {
	s.applicationTransitions = append(s.applicationTransitions, transition)
	return nil
}

func (s *stubStore) QueryApplicationTransitions(context.Context, ApplicationTransitionQueryFilter, order.By, page.Page) ([]ApplicationTransition, error) {
	return s.applicationTransitions, nil
}

func (s *stubStore) CountApplicationTransitions(context.Context, ApplicationTransitionQueryFilter) (int, error) {
	return len(s.applicationTransitions), nil
}

func (s *stubStore) CreateChecklistItem(_ context.Context, item ChecklistItem) error {
	s.checklistItems = append(s.checklistItems, item)
	return nil
}

func (s *stubStore) UpdateChecklistItem(_ context.Context, item ChecklistItem) error {
	for i, existing := range s.checklistItems {
		if existing.ID == item.ID {
			s.checklistItems[i] = item
			return nil
		}
	}
	s.checklistItems = append(s.checklistItems, item)
	return nil
}

func (s *stubStore) QueryChecklistItems(context.Context, ChecklistItemQueryFilter, order.By, page.Page) ([]ChecklistItem, error) {
	return s.checklistItems, nil
}

func (s *stubStore) CountChecklistItems(context.Context, ChecklistItemQueryFilter) (int, error) {
	return len(s.checklistItems), nil
}

func (s *stubStore) QueryChecklistItemByID(_ context.Context, itemID uuid.UUID) (ChecklistItem, error) {
	for _, item := range s.checklistItems {
		if item.ID == itemID {
			return item, nil
		}
	}
	return ChecklistItem{}, ErrChecklistItemNotFound
}

func (s *stubStore) CreateDocument(_ context.Context, document Document) error {
	s.documents = append(s.documents, document)
	return nil
}

func (s *stubStore) UpdateDocument(_ context.Context, document Document) error {
	for i, existing := range s.documents {
		if existing.ID == document.ID {
			s.documents[i] = document
			return nil
		}
	}
	s.documents = append(s.documents, document)
	return nil
}

func (s *stubStore) QueryDocuments(context.Context, DocumentQueryFilter, order.By, page.Page) ([]Document, error) {
	return s.documents, nil
}

func (s *stubStore) CountDocuments(context.Context, DocumentQueryFilter) (int, error) {
	return len(s.documents), nil
}

func (s *stubStore) QueryDocumentByID(_ context.Context, documentID uuid.UUID) (Document, error) {
	for _, document := range s.documents {
		if document.ID == documentID {
			return document, nil
		}
	}
	return Document{}, ErrDocumentNotFound
}

func (s *stubStore) CreateSyncJob(_ context.Context, job SyncJob) error {
	s.syncJobs = append(s.syncJobs, job)
	return nil
}

func (s *stubStore) UpdateSyncJob(_ context.Context, job SyncJob) error {
	for i, existing := range s.syncJobs {
		if existing.ID == job.ID {
			s.syncJobs[i] = job
			return nil
		}
	}
	s.syncJobs = append(s.syncJobs, job)
	return nil
}

func (s *stubStore) QuerySyncJobs(context.Context, SyncJobQueryFilter, order.By, page.Page) ([]SyncJob, error) {
	return s.syncJobs, nil
}

func (s *stubStore) CountSyncJobs(context.Context, SyncJobQueryFilter) (int, error) {
	return len(s.syncJobs), nil
}

func (s *stubStore) QuerySyncJobByID(_ context.Context, jobID uuid.UUID) (SyncJob, error) {
	for _, job := range s.syncJobs {
		if job.ID == jobID {
			return job, nil
		}
	}
	return SyncJob{}, ErrSyncJobNotFound
}

func (s *stubStore) CreateSyncEvent(_ context.Context, event SyncEvent) error {
	s.syncEvents = append(s.syncEvents, event)
	return nil
}

func (s *stubStore) UpdateSyncEvent(_ context.Context, event SyncEvent) error {
	for i, existing := range s.syncEvents {
		if existing.ID == event.ID {
			s.syncEvents[i] = event
			return nil
		}
	}
	s.syncEvents = append(s.syncEvents, event)
	return nil
}

func (s *stubStore) QuerySyncEvents(context.Context, SyncEventQueryFilter, order.By, page.Page) ([]SyncEvent, error) {
	return s.syncEvents, nil
}

func (s *stubStore) CountSyncEvents(context.Context, SyncEventQueryFilter) (int, error) {
	return len(s.syncEvents), nil
}

func (s *stubStore) QuerySyncEventByID(_ context.Context, eventID uuid.UUID) (SyncEvent, error) {
	for _, event := range s.syncEvents {
		if event.ID == eventID {
			return event, nil
		}
	}
	return SyncEvent{}, ErrSyncEventNotFound
}
