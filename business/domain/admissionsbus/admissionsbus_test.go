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

func TestCreateConstituentDefaultsKenyaNotificationPreferences(t *testing.T) {
	t.Parallel()

	store := &stubStore{constituents: map[uuid.UUID]Constituent{}}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	email := mail.Address{Address: "applicant@example.com"}

	created, err := bus.CreateConstituent(context.Background(), NewConstituent{
		FirstName:    "Ada",
		LastName:     "Applicant",
		DateOfBirth:  time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		PrimaryEmail: email,
		PrimaryPhone: "+254712345678",
	})
	if err != nil {
		t.Fatalf("CreateConstituent returned error: %v", err)
	}

	want := KenyaDefaultNotificationPreferences()
	assertNotificationPreferences(t, created.NotificationPreferences, want)
	assertNotificationPreferences(t, store.constituents[created.ID].NotificationPreferences, want)
}

func TestUpdateConstituentNotificationOptOut(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	cst := Constituent{
		ID:                      uuid.New(),
		LifecycleStage:          LifecycleStageProspect,
		DuplicateStatus:         DuplicateStatusActive,
		NotificationPreferences: KenyaDefaultNotificationPreferences(),
	}

	preferences := NotificationPreferences{
		SMSOptIn:      false,
		WhatsAppOptIn: false,
		EmailOptIn:    true,
		Priority: []NotificationChannel{
			NotificationChannelSMS,
			NotificationChannelWhatsApp,
			NotificationChannelEmail,
		},
	}
	updated, err := bus.UpdateConstituent(context.Background(), cst, UpdateConstituent{NotificationPreferences: &preferences})
	if err != nil {
		t.Fatalf("UpdateConstituent returned error: %v", err)
	}

	assertNotificationPreferences(t, updated.NotificationPreferences, preferences)
}

func TestCreateConstituentValidatesNotificationPriority(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	email := mail.Address{Address: "applicant@example.com"}

	tests := []struct {
		name        string
		preferences NotificationPreferences
		want        error
	}{
		{
			name: "duplicate channel",
			preferences: NotificationPreferences{
				SMSOptIn:   true,
				EmailOptIn: true,
				Priority: []NotificationChannel{
					NotificationChannelSMS,
					NotificationChannelSMS,
					NotificationChannelEmail,
				},
			},
			want: ErrNotificationPriorityDuplicate,
		},
		{
			name: "missing channel",
			preferences: NotificationPreferences{
				SMSOptIn:   true,
				EmailOptIn: true,
				Priority: []NotificationChannel{
					NotificationChannelSMS,
					NotificationChannelEmail,
				},
			},
			want: ErrNotificationPriorityIncomplete,
		},
		{
			name: "unknown channel",
			preferences: NotificationPreferences{
				SMSOptIn:   true,
				EmailOptIn: true,
				Priority: []NotificationChannel{
					NotificationChannelSMS,
					NotificationChannel("PUSH"),
					NotificationChannelEmail,
				},
			},
			want: ErrInvalidNotificationChannel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateConstituent(context.Background(), NewConstituent{
				FirstName:               "Ada",
				LastName:                "Applicant",
				DateOfBirth:             time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				PrimaryEmail:            email,
				PrimaryPhone:            "+254712345678",
				NotificationPreferences: &tt.preferences,
			})
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

func TestCreateCustomFieldDefinitionValidatesRegistryFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nd   NewCustomFieldDefinition
		want error
	}{
		{
			name: "owner",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwner("INQUIRY"),
				FieldKey:     "scholarship_level",
				Label:        "Scholarship level",
				DataType:     CustomFieldDataTypeText,
				DisplayOrder: 1,
			},
			want: ErrCustomFieldOwnerInvalid,
		},
		{
			name: "key required",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwnerConstituent,
				Label:        "Scholarship level",
				DataType:     CustomFieldDataTypeText,
				DisplayOrder: 1,
			},
			want: ErrCustomFieldKeyRequired,
		},
		{
			name: "label required",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwnerConstituent,
				FieldKey:     "scholarship_level",
				DataType:     CustomFieldDataTypeText,
				DisplayOrder: 1,
			},
			want: ErrCustomFieldLabelRequired,
		},
		{
			name: "data type",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwnerApplication,
				FieldKey:     "portfolio_score",
				Label:        "Portfolio score",
				DataType:     CustomFieldDataType("RANGE"),
				DisplayOrder: 1,
			},
			want: ErrCustomFieldDataTypeInvalid,
		},
		{
			name: "select options",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwnerApplication,
				FieldKey:     "scholarship_level",
				Label:        "Scholarship level",
				DataType:     CustomFieldDataTypeSelect,
				DisplayOrder: 1,
			},
			want: ErrCustomFieldOptionsRequired,
		},
		{
			name: "display order",
			nd: NewCustomFieldDefinition{
				Owner:        CustomFieldOwnerApplication,
				FieldKey:     "portfolio_score",
				Label:        "Portfolio score",
				DataType:     CustomFieldDataTypeNumber,
				DisplayOrder: -1,
			},
			want: ErrCustomFieldOrderInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateCustomFieldDefinition(context.Background(), tt.nd)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateCustomFieldDefinitionStoresNormalizedSeams(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	description := " Used for imports and reports "
	validation := " max:50 "

	definition, err := bus.CreateCustomFieldDefinition(context.Background(), NewCustomFieldDefinition{
		Owner:        CustomFieldOwnerConstituent,
		FieldKey:     " scholarship_level ",
		Label:        " Scholarship level ",
		Description:  &description,
		DataType:     CustomFieldDataTypeSelect,
		Options:      []string{" Full ", "Partial", "Full", ""},
		Validation:   &validation,
		Searchable:   true,
		Reportable:   true,
		Importable:   true,
		Exportable:   true,
		DisplayOrder: 4,
		Active:       true,
	})
	if err != nil {
		t.Fatalf("CreateCustomFieldDefinition returned error: %v", err)
	}

	if definition.Owner != CustomFieldOwnerConstituent {
		t.Fatalf("Owner = %s, want %s", definition.Owner, CustomFieldOwnerConstituent)
	}
	if definition.FieldKey != "scholarship_level" {
		t.Fatalf("FieldKey = %q, want scholarship_level", definition.FieldKey)
	}
	if definition.Label != "Scholarship level" {
		t.Fatalf("Label = %q, want Scholarship level", definition.Label)
	}
	if len(definition.Options) != 2 {
		t.Fatalf("Options = %v, want 2 unique options", definition.Options)
	}
	if !definition.Searchable || !definition.Reportable || !definition.Importable || !definition.Exportable {
		t.Fatalf("seams not enabled: %+v", definition)
	}
	if len(store.customFieldDefinitions) != 1 {
		t.Fatalf("stored definitions = %d, want 1", len(store.customFieldDefinitions))
	}
}

func TestSetCustomFieldValueRequiresDefinitionOwnerAndRecord(t *testing.T) {
	t.Parallel()

	definitionID := uuid.New()
	constituentID := uuid.New()
	store := &stubStore{
		customFieldDefinitions: []CustomFieldDefinition{
			{
				ID:       definitionID,
				Owner:    CustomFieldOwnerConstituent,
				FieldKey: "scholarship_level",
				Label:    "Scholarship level",
				DataType: CustomFieldDataTypeText,
				Active:   true,
			},
		},
		constituents: map[uuid.UUID]Constituent{
			constituentID: {ID: constituentID, LifecycleStage: LifecycleStageApplicant},
		},
	}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	value, err := bus.SetCustomFieldValue(context.Background(), NewCustomFieldValue{
		DefinitionID: definitionID,
		Owner:        CustomFieldOwnerConstituent,
		OwnerID:      constituentID,
		Value:        " Full scholarship ",
	})
	if err != nil {
		t.Fatalf("SetCustomFieldValue returned error: %v", err)
	}
	if value.Value != "Full scholarship" {
		t.Fatalf("Value = %q, want Full scholarship", value.Value)
	}
	if len(store.customFieldValues) != 1 {
		t.Fatalf("stored values = %d, want 1", len(store.customFieldValues))
	}

	_, err = bus.SetCustomFieldValue(context.Background(), NewCustomFieldValue{
		DefinitionID: definitionID,
		Owner:        CustomFieldOwnerApplication,
		OwnerID:      constituentID,
		Value:        "Full scholarship",
	})
	if !errors.Is(err, ErrCustomFieldOwnerInvalid) {
		t.Fatalf("err = %v, want %v", err, ErrCustomFieldOwnerInvalid)
	}
}

func TestCreateImportBatchValidatesSupportedFileTypesAndMapping(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()

	tests := []struct {
		name string
		nb   NewImportBatch
		want error
	}{
		{
			name: "file type",
			nb: NewImportBatch{
				Source:       ImportSourceManualUpload,
				FileType:     ImportFileType("PDF"),
				Target:       ImportTargetConstituents,
				Status:       ImportBatchStatusPreviewed,
				FileName:     "constituents.pdf",
				UploadedByID: uuid.New(),
				FieldMapping: map[string]string{"First Name": "firstName"},
			},
			want: ErrInvalidImportFileType,
		},
		{
			name: "uploader",
			nb: NewImportBatch{
				Source:       ImportSourceManualUpload,
				FileType:     ImportFileTypeCSV,
				Target:       ImportTargetConstituents,
				Status:       ImportBatchStatusPreviewed,
				FileName:     "constituents.csv",
				FieldMapping: map[string]string{"First Name": "firstName"},
			},
			want: ErrImportUploaderRequired,
		},
		{
			name: "rows",
			nb: NewImportBatch{
				Source:       ImportSourceManualUpload,
				FileType:     ImportFileTypeXLSX,
				Target:       ImportTargetApplications,
				Status:       ImportBatchStatusPreviewed,
				FileName:     "applications.xlsx",
				UploadedByID: uuid.New(),
				TotalRows:    1,
				ValidRows:    1,
				InvalidRows:  1,
				FieldMapping: map[string]string{"Program": "programID"},
			},
			want: ErrImportRowsInvalid,
		},
		{
			name: "mapping",
			nb: NewImportBatch{
				Source:       ImportSourceManualUpload,
				FileType:     ImportFileTypeCSV,
				Target:       ImportTargetConstituents,
				Status:       ImportBatchStatusPreviewed,
				FileName:     "constituents.csv",
				UploadedByID: uuid.New(),
				TotalRows:    1,
			},
			want: ErrImportFieldMappingRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateImportBatch(context.Background(), tt.nb)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCreateImportBatchStoresAuditReadyPreview(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	uploaderID := uuid.New()
	storageKey := "imports/constituents-2026.csv"
	reportKey := "imports/constituents-2026-invalid.csv"
	summary := "2 invalid rows require correction"

	batch, err := bus.CreateImportBatch(context.Background(), NewImportBatch{
		Source:            ImportSourceManualUpload,
		FileType:          ImportFileTypeCSV,
		Target:            ImportTargetConstituents,
		Status:            ImportBatchStatusValidationFailed,
		FileName:          " constituents-2026.csv ",
		StorageKey:        &storageKey,
		UploadedByID:      uploaderID,
		TotalRows:         10,
		ValidRows:         8,
		InvalidRows:       2,
		DuplicateRows:     1,
		FieldMapping:      map[string]string{" First Name ": " firstName ", "": "ignored"},
		InvalidReportKey:  &reportKey,
		ValidationSummary: &summary,
	})
	if err != nil {
		t.Fatalf("CreateImportBatch returned error: %v", err)
	}

	if batch.FileType != ImportFileTypeCSV {
		t.Fatalf("FileType = %s, want %s", batch.FileType, ImportFileTypeCSV)
	}
	if batch.FileName != "constituents-2026.csv" {
		t.Fatalf("FileName = %q, want constituents-2026.csv", batch.FileName)
	}
	if batch.FieldMapping["First Name"] != "firstName" {
		t.Fatalf("FieldMapping = %v, want normalized First Name mapping", batch.FieldMapping)
	}
	if len(store.importBatches) != 1 {
		t.Fatalf("stored import batches = %d, want 1", len(store.importBatches))
	}
}

func TestCreateImportInvalidRowsStoresCorrectionReportRows(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	batchID := uuid.New()
	fieldName := "primaryEmail"

	rows, err := bus.CreateImportInvalidRows(context.Background(), []NewImportInvalidRow{
		{
			BatchID:     batchID,
			RowNumber:   3,
			FieldName:   &fieldName,
			RawData:     map[string]string{" Email ": " not-an-email "},
			ErrorCode:   "INVALID_EMAIL",
			ErrorDetail: "Primary email must be valid",
		},
	})
	if err != nil {
		t.Fatalf("CreateImportInvalidRows returned error: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rows[0].RawData["Email"] != "not-an-email" {
		t.Fatalf("RawData = %v, want normalized email data", rows[0].RawData)
	}
	if len(store.importInvalidRows) != 1 {
		t.Fatalf("stored invalid rows = %d, want 1", len(store.importInvalidRows))
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
				Adapter:   IntegrationAdapterKUCCPS,
				Operation: "BATCH_PLACEMENT_PULL",
				Status:    SyncJobStatus("UNKNOWN"),
				Direction: SyncDirectionInbound,
			},
			want: ErrInvalidSyncJobStatus,
		},
		{
			name: "valid direction",
			nj: NewSyncJob{
				Name:      "Nightly SIS reconciliation",
				Adapter:   IntegrationAdapterKUCCPS,
				Operation: "BATCH_PLACEMENT_PULL",
				Status:    SyncJobStatusQueued,
				Direction: SyncDirection("UNKNOWN"),
			},
			want: ErrInvalidSyncDirection,
		},
		{
			name: "adapter required",
			nj: NewSyncJob{
				Name:      "Nightly SIS reconciliation",
				Operation: "BATCH_PLACEMENT_PULL",
				Status:    SyncJobStatusQueued,
				Direction: SyncDirectionInbound,
			},
			want: ErrInvalidIntegrationAdapter,
		},
		{
			name: "operation required",
			nj: NewSyncJob{
				Name:      "Nightly SIS reconciliation",
				Adapter:   IntegrationAdapterKUCCPS,
				Status:    SyncJobStatusQueued,
				Direction: SyncDirectionInbound,
			},
			want: ErrSyncJobOperationRequired,
		},
		{
			name: "max attempts invalid",
			nj: NewSyncJob{
				Name:        "Nightly SIS reconciliation",
				Adapter:     IntegrationAdapterKUCCPS,
				Operation:   "BATCH_PLACEMENT_PULL",
				Status:      SyncJobStatusQueued,
				Direction:   SyncDirectionInbound,
				MaxAttempts: -1,
			},
			want: ErrInvalidMaxAttempts,
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

func TestCreateSyncJobDefaultsAdapterRetryMetadata(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	job, err := bus.CreateSyncJob(context.Background(), NewSyncJob{
		Name:      " KUCCPS nightly placement pull ",
		Adapter:   IntegrationAdapterKUCCPS,
		Operation: " BATCH_PLACEMENT_PULL ",
		Status:    SyncJobStatusQueued,
		Direction: SyncDirectionInbound,
	})
	if err != nil {
		t.Fatalf("CreateSyncJob returned error: %v", err)
	}

	if job.Name != "KUCCPS nightly placement pull" {
		t.Fatalf("Name = %q, want trimmed name", job.Name)
	}
	if job.Operation != "BATCH_PLACEMENT_PULL" {
		t.Fatalf("Operation = %q, want trimmed operation", job.Operation)
	}
	if job.Adapter != IntegrationAdapterKUCCPS {
		t.Fatalf("Adapter = %s, want %s", job.Adapter, IntegrationAdapterKUCCPS)
	}
	if job.MaxAttempts != 3 {
		t.Fatalf("MaxAttempts = %d, want 3", job.MaxAttempts)
	}
	if len(store.syncJobs) != 1 {
		t.Fatalf("stored sync jobs = %d, want 1", len(store.syncJobs))
	}
}

func TestUpdateSyncJobRetryStateTransitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		initial    SyncJobStatus
		update     UpdateSyncJob
		wantStatus SyncJobStatus
		wantErr    error
	}{
		{
			name:       "queued to running",
			initial:    SyncJobStatusQueued,
			update:     UpdateSyncJob{Status: SyncJobStatusRunning, AttemptCount: 1},
			wantStatus: SyncJobStatusRunning,
		},
		{
			name:       "running to succeeded",
			initial:    SyncJobStatusRunning,
			update:     UpdateSyncJob{Status: SyncJobStatusSucceeded, AttemptCount: 1},
			wantStatus: SyncJobStatusSucceeded,
		},
		{
			name:       "running to retry ready",
			initial:    SyncJobStatusRunning,
			update:     UpdateSyncJob{Status: SyncJobStatusRetryReady, AttemptCount: 1, Retryable: true},
			wantStatus: SyncJobStatusRetryReady,
		},
		{
			name:       "retry ready to running",
			initial:    SyncJobStatusRetryReady,
			update:     UpdateSyncJob{Status: SyncJobStatusRunning, AttemptCount: 2},
			wantStatus: SyncJobStatusRunning,
		},
		{
			name:    "succeeded cannot rerun",
			initial: SyncJobStatusSucceeded,
			update:  UpdateSyncJob{Status: SyncJobStatusRunning, AttemptCount: 1},
			wantErr: ErrInvalidSyncJobTransition,
		},
		{
			name:    "attempts cannot exceed max",
			initial: SyncJobStatusRunning,
			update:  UpdateSyncJob{Status: SyncJobStatusFailed, AttemptCount: 4},
			wantErr: ErrMaxAttemptsExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := &stubStore{}
			bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
			job := SyncJob{
				ID:          uuid.New(),
				Name:        "KNEC result verification",
				Adapter:     IntegrationAdapterKNEC,
				Operation:   "RESULT_VERIFICATION",
				Status:      tt.initial,
				Direction:   SyncDirectionInbound,
				MaxAttempts: 3,
			}

			updated, err := bus.UpdateSyncJob(context.Background(), job, tt.update)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if updated.Status != tt.wantStatus {
				t.Fatalf("Status = %s, want %s", updated.Status, tt.wantStatus)
			}
			if len(store.syncJobs) != 1 {
				t.Fatalf("stored sync jobs = %d, want 1", len(store.syncJobs))
			}
		})
	}
}

func TestQuerySyncJobsFiltersByAdapterIsolation(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	kuccps := SyncJob{ID: uuid.New(), Name: "KUCCPS", Adapter: IntegrationAdapterKUCCPS, Operation: "PLACEMENT_PULL", Status: SyncJobStatusRunning, Direction: SyncDirectionInbound, MaxAttempts: 3}
	knec := SyncJob{ID: uuid.New(), Name: "KNEC", Adapter: IntegrationAdapterKNEC, Operation: "RESULT_VERIFICATION", Status: SyncJobStatusFailed, Direction: SyncDirectionInbound, MaxAttempts: 3}
	store.syncJobs = []SyncJob{kuccps, knec}
	adapter := IntegrationAdapterKNEC
	status := SyncJobStatusFailed

	jobs, err := bus.QuerySyncJobs(context.Background(), SyncJobQueryFilter{Adapter: &adapter, Status: &status}, DefaultSyncJobOrderBy, page.MustParse("1", "10"))
	if err != nil {
		t.Fatalf("QuerySyncJobs returned error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("jobs = %d, want 1", len(jobs))
	}
	if jobs[0].ID != knec.ID {
		t.Fatalf("job ID = %s, want KNEC failed job %s", jobs[0].ID, knec.ID)
	}
}

func TestEnqueueSyncEventStoresApprovedRealtimeEvent(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	applicationID := uuid.New()

	event, err := bus.EnqueueSyncEvent(context.Background(), NewSyncEvent{
		Adapter:      IntegrationAdapterKUCCPS,
		Operation:    "APPLICATION_SUBMISSION",
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
	if event.Adapter != IntegrationAdapterKUCCPS {
		t.Fatalf("Adapter = %s, want %s", event.Adapter, IntegrationAdapterKUCCPS)
	}
	if event.MaxAttempts != 3 {
		t.Fatalf("MaxAttempts = %d, want 3", event.MaxAttempts)
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
				Adapter:      IntegrationAdapterKUCCPS,
				Operation:    "APPLICATION_SUBMISSION",
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
				Adapter:     IntegrationAdapterKUCCPS,
				Operation:   "APPLICATION_DECISION",
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
				Adapter:      IntegrationAdapterKUCCPS,
				Operation:    "DOCUMENT_STATUS",
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
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
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
		ApplicationType: ApplicationTypeDiploma,
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
		ApplicationType: ApplicationTypeDiploma,
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
				ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
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
				ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
				Name:            "Freshman",
			},
			want: ErrFormTemplateFieldsRequired,
		},
		{
			name: "checklist invalid",
			nt: NewApplicationFormTemplate{
				ProgramID:       programID,
				AcademicTermID:  termID,
				ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
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
			ApplicationType: ApplicationTypeDiploma,
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

func assertNotificationPreferences(t *testing.T, got NotificationPreferences, want NotificationPreferences) {
	t.Helper()

	if got.SMSOptIn != want.SMSOptIn {
		t.Fatalf("SMSOptIn = %t, want %t", got.SMSOptIn, want.SMSOptIn)
	}
	if got.WhatsAppOptIn != want.WhatsAppOptIn {
		t.Fatalf("WhatsAppOptIn = %t, want %t", got.WhatsAppOptIn, want.WhatsAppOptIn)
	}
	if got.EmailOptIn != want.EmailOptIn {
		t.Fatalf("EmailOptIn = %t, want %t", got.EmailOptIn, want.EmailOptIn)
	}
	if len(got.Priority) != len(want.Priority) {
		t.Fatalf("Priority length = %d, want %d", len(got.Priority), len(want.Priority))
	}
	for i, channel := range want.Priority {
		if got.Priority[i] != channel {
			t.Fatalf("Priority[%d] = %s, want %s", i, got.Priority[i], channel)
		}
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
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
		Status:          ApplicationStatusDraft,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
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
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
		Status:          ApplicationStatusWithdrawn,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	created, err := bus.CreateApplication(context.Background(), NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
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

func TestCreateEventStoresNormalizedData(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)
	registrationDeadline := time.Date(2026, time.June, 5, 9, 0, 0, 0, time.UTC)

	event, err := bus.CreateEvent(context.Background(), NewEvent{
		Title:                   "  Spring Open Day  ",
		Type:                    EventTypeOpenDay,
		Status:                  EventStatusUpcoming,
		Description:             "  Campus-wide visit day  ",
		StartTime:               time.Date(2026, time.June, 10, 9, 0, 0, 0, time.UTC),
		EndTime:                 time.Date(2026, time.June, 10, 15, 0, 0, 0, time.UTC),
		Location:                "  Main Campus  ",
		Capacity:                250,
		RegistrationDeadline:    &registrationDeadline,
		AutoConfirmationEnabled: true,
		AutoReminderEnabled:     true,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}

	if event.Title != "Spring Open Day" {
		t.Fatalf("Title = %q, want trimmed title", event.Title)
	}

	if event.Description != "Campus-wide visit day" {
		t.Fatalf("Description = %q, want trimmed description", event.Description)
	}

	if event.Location != "Main Campus" {
		t.Fatalf("Location = %q, want trimmed location", event.Location)
	}

	if len(store.events) != 1 {
		t.Fatalf("event count = %d, want 1", len(store.events))
	}
}

func TestCreateEventValidatesRequiredFields(t *testing.T) {
	t.Parallel()

	bus := newTestBusiness()
	start := time.Date(2026, time.June, 10, 9, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	tests := []struct {
		name string
		ne   NewEvent
		want error
	}{
		{
			name: "title",
			ne: NewEvent{
				Type:        EventTypeOpenDay,
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Location:    "Main Campus",
				Capacity:    10,
			},
			want: ErrEventTitleRequired,
		},
		{
			name: "type",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventType("unknown"),
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Location:    "Main Campus",
				Capacity:    10,
			},
			want: ErrInvalidEventType,
		},
		{
			name: "status",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventTypeOpenDay,
				Status:      EventStatus("unknown"),
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Location:    "Main Campus",
				Capacity:    10,
			},
			want: ErrInvalidEventStatus,
		},
		{
			name: "description",
			ne: NewEvent{
				Title:     "Open Day",
				Type:      EventTypeOpenDay,
				Status:    EventStatusUpcoming,
				StartTime: start,
				EndTime:   end,
				Location:  "Main Campus",
				Capacity:  10,
			},
			want: ErrEventDescriptionRequired,
		},
		{
			name: "date range",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventTypeOpenDay,
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   end,
				EndTime:     start,
				Location:    "Main Campus",
				Capacity:    10,
			},
			want: ErrEventDateRangeInvalid,
		},
		{
			name: "location",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventTypeOpenDay,
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Capacity:    10,
			},
			want: ErrEventLocationRequired,
		},
		{
			name: "capacity",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventTypeOpenDay,
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Location:    "Main Campus",
				Capacity:    -1,
			},
			want: ErrEventCapacityInvalid,
		},
		{
			name: "registration deadline after start",
			ne: NewEvent{
				Title:       "Open Day",
				Type:        EventTypeOpenDay,
				Status:      EventStatusUpcoming,
				Description: "Description",
				StartTime:   start,
				EndTime:     end,
				Location:    "Main Campus",
				Capacity:    10,
				RegistrationDeadline: func() *time.Time {
					deadline := start.Add(time.Hour)
					return &deadline
				}(),
			},
			want: ErrEventDateRangeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := bus.CreateEvent(context.Background(), tt.ne)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestUpdateEventReplacesMutableFields(t *testing.T) {
	t.Parallel()

	store := &stubStore{}
	start := time.Date(2026, time.June, 10, 9, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	event := Event{
		ID:          uuid.New(),
		Title:       "Open Day",
		Type:        EventTypeOpenDay,
		Status:      EventStatusDraft,
		Description: "Description",
		StartTime:   start,
		EndTime:     end,
		Location:    "Main Campus",
		Capacity:    100,
		DateCreated: start,
		DateUpdated: start,
	}
	store.events = append(store.events, event)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	updated, err := bus.UpdateEvent(context.Background(), event, NewEvent{
		Title:                   " Updated Webinar ",
		Type:                    EventTypeWebinar,
		Status:                  EventStatusUpcoming,
		Description:             " Updated description ",
		StartTime:               start.Add(24 * time.Hour),
		EndTime:                 end.Add(24 * time.Hour),
		Location:                " Zoom ",
		IsVirtual:               true,
		Capacity:                250,
		AutoConfirmationEnabled: true,
	})
	if err != nil {
		t.Fatalf("UpdateEvent returned error: %v", err)
	}

	if updated.ID != event.ID {
		t.Fatalf("ID changed from %s to %s", event.ID, updated.ID)
	}

	if updated.Title != "Updated Webinar" || updated.Location != "Zoom" {
		t.Fatalf("updated event = %+v, want trimmed mutable fields", updated)
	}

	if updated.Type != EventTypeWebinar || updated.Status != EventStatusUpcoming {
		t.Fatalf("type/status = %s/%s, want %s/%s", updated.Type, updated.Status, EventTypeWebinar, EventStatusUpcoming)
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

func TestUpdateApplicationDraftReplacesApplicantEditableFields(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	application := Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
		Status:          ApplicationStatusDraft,
	}
	store.applications = append(store.applications, application)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	updated, err := bus.UpdateApplicationDraft(context.Background(), application, NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
		KCSEResult: &ApplicationKCSEResult{
			IndexNumber: "12345678901",
			ExamYear:    2025,
			MeanGrade:   "b+",
			MeanPoints:  68,
		},
	})
	if err != nil {
		t.Fatalf("UpdateApplicationDraft returned error: %v", err)
	}

	if updated.ID != application.ID || updated.Status != ApplicationStatusDraft {
		t.Fatalf("identity/status changed: id=%s status=%s", updated.ID, updated.Status)
	}
	if updated.ApplicationType != ApplicationTypeSelfSponsoredUndergrad {
		t.Fatalf("ApplicationType = %s, want %s", updated.ApplicationType, ApplicationTypeSelfSponsoredUndergrad)
	}
	if updated.KCSEResult == nil || updated.KCSEResult.MeanGrade != "B+" {
		t.Fatalf("KCSEResult = %+v, want normalized B+", updated.KCSEResult)
	}
}

func TestUpdateApplicationDraftRequiresDraftStatus(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.UpdateApplicationDraft(context.Background(), Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
		Status:          ApplicationStatusSubmitted,
	}, NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
	})

	if !errors.Is(err, ErrApplicationNotDraft) {
		t.Fatalf("err = %v, want %v", err, ErrApplicationNotDraft)
	}
}

func TestUpdateApplicationDraftRejectsConstituentChange(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.UpdateApplicationDraft(context.Background(), Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
		Status:          ApplicationStatusDraft,
	}, NewApplication{
		ConstituentID:   uuid.New(),
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
	})

	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("err = %v, want %v", err, ErrApplicationNotFound)
	}
}

func TestUpdateApplicationDraftPreventsDuplicateActiveApplication(t *testing.T) {
	t.Parallel()

	constituentID := uuid.New()
	programID := uuid.New()
	termID := uuid.New()
	store := newApplicationStubStore(constituentID, programID, termID)
	application := Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       uuid.New(),
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
		Status:          ApplicationStatusDraft,
	}
	store.programs[application.ProgramID] = Program{ID: application.ProgramID, Active: true}
	store.applications = append(store.applications, application, Application{
		ID:              uuid.New(),
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeSelfSponsoredUndergrad,
		Status:          ApplicationStatusSubmitted,
	})
	bus := NewBusiness(logger.New(ioDiscard{}, logger.LevelInfo, "TEST", func(context.Context) string { return "" }), nil, store)

	_, err := bus.UpdateApplicationDraft(context.Background(), application, NewApplication{
		ConstituentID:   constituentID,
		ProgramID:       programID,
		AcademicTermID:  termID,
		ApplicationType: ApplicationTypeDiploma,
	})

	if !errors.Is(err, ErrDuplicateApplication) {
		t.Fatalf("err = %v, want %v", err, ErrDuplicateApplication)
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
	events                 []Event
	eventRegistrations     []EventRegistration
	leadScoreRules         []LeadScoreRule
	leadScores             []LeadScore
	applicationTemplates   []ApplicationFormTemplate
	customFieldDefinitions []CustomFieldDefinition
	customFieldValues      []CustomFieldValue
	constituents           map[uuid.UUID]Constituent
	constituentByEmail     map[string]Constituent
	duplicateReviews       []DuplicateReview
	programs               map[uuid.UUID]Program
	terms                  map[uuid.UUID]AcademicTerm
	applications           []Application
	applicationTransitions []ApplicationTransition
	checklistItems         []ChecklistItem
	documents              []Document
	importBatches          []ImportBatch
	importInvalidRows      []ImportInvalidRow
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

func (s *stubStore) QueryEvents(_ context.Context, filter EventQueryFilter, _ order.By, _ page.Page) ([]Event, error) {
	events := make([]Event, 0, len(s.events))
	for _, event := range s.events {
		if filter.ID != nil && event.ID != *filter.ID {
			continue
		}
		if filter.Type != nil && event.Type != *filter.Type {
			continue
		}
		if filter.Status != nil && event.Status != *filter.Status {
			continue
		}
		if filter.Virtual != nil && event.IsVirtual != *filter.Virtual {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *stubStore) CountEvents(_ context.Context, filter EventQueryFilter) (int, error) {
	events, err := s.QueryEvents(context.Background(), filter, order.By{}, page.Page{})
	if err != nil {
		return 0, err
	}

	return len(events), nil
}

func (s *stubStore) QueryEventByID(_ context.Context, eventID uuid.UUID) (Event, error) {
	for _, event := range s.events {
		if event.ID == eventID {
			return event, nil
		}
	}

	return Event{}, ErrEventNotFound
}

func (s *stubStore) CreateEvent(_ context.Context, event Event) error {
	s.events = append(s.events, event)
	return nil
}

func (s *stubStore) UpdateEvent(_ context.Context, event Event) error {
	for i, existing := range s.events {
		if existing.ID == event.ID {
			s.events[i] = event
			return nil
		}
	}

	s.events = append(s.events, event)
	return nil
}

func (s *stubStore) QueryEventRegistrations(_ context.Context, filter EventRegistrationQueryFilter, _ order.By, _ page.Page) ([]EventRegistration, error) {
	registrations := make([]EventRegistration, 0, len(s.eventRegistrations))
	for _, registration := range s.eventRegistrations {
		if filter.ID != nil && registration.ID != *filter.ID {
			continue
		}
		if filter.EventID != nil && registration.EventID != *filter.EventID {
			continue
		}
		if filter.ConstituentID != nil {
			if registration.ConstituentID == nil || *registration.ConstituentID != *filter.ConstituentID {
				continue
			}
		}
		if filter.Status != nil && registration.Status != *filter.Status {
			continue
		}
		if filter.MatchStatus != nil && registration.MatchStatus != *filter.MatchStatus {
			continue
		}
		if filter.Source != nil && registration.Source != *filter.Source {
			continue
		}
		registrations = append(registrations, registration)
	}

	return registrations, nil
}

func (s *stubStore) CountEventRegistrations(_ context.Context, filter EventRegistrationQueryFilter) (int, error) {
	registrations, err := s.QueryEventRegistrations(context.Background(), filter, order.By{}, page.Page{})
	if err != nil {
		return 0, err
	}

	return len(registrations), nil
}

func (s *stubStore) QueryEventRegistrationByID(_ context.Context, registrationID uuid.UUID) (EventRegistration, error) {
	for _, registration := range s.eventRegistrations {
		if registration.ID == registrationID {
			return registration, nil
		}
	}

	return EventRegistration{}, ErrEventRegistrationNotFound
}

func (s *stubStore) CreateEventRegistration(_ context.Context, registration EventRegistration) error {
	s.eventRegistrations = append(s.eventRegistrations, registration)
	return nil
}

func (s *stubStore) UpdateEventRegistration(_ context.Context, registration EventRegistration) error {
	for i, existing := range s.eventRegistrations {
		if existing.ID == registration.ID {
			s.eventRegistrations[i] = registration
			return nil
		}
	}

	s.eventRegistrations = append(s.eventRegistrations, registration)
	return nil
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

func (s *stubStore) QueryConstituentByNationalID(context.Context, string) (Constituent, error) {
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByUPI(context.Context, string) (Constituent, error) {
	return Constituent{}, ErrConstituentNotFound
}

func (s *stubStore) QueryConstituentByKCSEIndexNumber(context.Context, string) (Constituent, error) {
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

func (s *stubStore) CreateCustomFieldDefinition(_ context.Context, definition CustomFieldDefinition) error {
	s.customFieldDefinitions = append(s.customFieldDefinitions, definition)
	return nil
}

func (s *stubStore) UpdateCustomFieldDefinition(_ context.Context, definition CustomFieldDefinition) error {
	for i, existing := range s.customFieldDefinitions {
		if existing.ID == definition.ID {
			s.customFieldDefinitions[i] = definition
			return nil
		}
	}
	s.customFieldDefinitions = append(s.customFieldDefinitions, definition)
	return nil
}

func (s *stubStore) QueryCustomFieldDefinitions(_ context.Context, filter CustomFieldDefinitionQueryFilter, _ order.By, _ page.Page) ([]CustomFieldDefinition, error) {
	var definitions []CustomFieldDefinition
	for _, definition := range s.customFieldDefinitions {
		if filter.ID != nil && definition.ID != *filter.ID {
			continue
		}
		if filter.Owner != nil && definition.Owner != *filter.Owner {
			continue
		}
		if filter.Active != nil && definition.Active != *filter.Active {
			continue
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (s *stubStore) CountCustomFieldDefinitions(_ context.Context, filter CustomFieldDefinitionQueryFilter) (int, error) {
	definitions, err := s.QueryCustomFieldDefinitions(context.Background(), filter, order.By{}, page.Page{})
	if err != nil {
		return 0, err
	}
	return len(definitions), nil
}

func (s *stubStore) QueryCustomFieldDefinitionByID(_ context.Context, definitionID uuid.UUID) (CustomFieldDefinition, error) {
	for _, definition := range s.customFieldDefinitions {
		if definition.ID == definitionID {
			return definition, nil
		}
	}
	return CustomFieldDefinition{}, ErrCustomFieldDefinitionNotFound
}

func (s *stubStore) SetCustomFieldValue(_ context.Context, value CustomFieldValue) error {
	for i, existing := range s.customFieldValues {
		if existing.DefinitionID == value.DefinitionID && existing.Owner == value.Owner && existing.OwnerID == value.OwnerID {
			value.ID = existing.ID
			value.DateCreated = existing.DateCreated
			s.customFieldValues[i] = value
			return nil
		}
	}
	s.customFieldValues = append(s.customFieldValues, value)
	return nil
}

func (s *stubStore) QueryCustomFieldValues(_ context.Context, filter CustomFieldValueQueryFilter, _ order.By, _ page.Page) ([]CustomFieldValue, error) {
	var values []CustomFieldValue
	for _, value := range s.customFieldValues {
		if filter.ID != nil && value.ID != *filter.ID {
			continue
		}
		if filter.DefinitionID != nil && value.DefinitionID != *filter.DefinitionID {
			continue
		}
		if filter.Owner != nil && value.Owner != *filter.Owner {
			continue
		}
		if filter.OwnerID != nil && value.OwnerID != *filter.OwnerID {
			continue
		}
		values = append(values, value)
	}
	return values, nil
}

func (s *stubStore) CountCustomFieldValues(_ context.Context, filter CustomFieldValueQueryFilter) (int, error) {
	values, err := s.QueryCustomFieldValues(context.Background(), filter, order.By{}, page.Page{})
	if err != nil {
		return 0, err
	}
	return len(values), nil
}

func (s *stubStore) QueryCustomFieldValueByID(_ context.Context, valueID uuid.UUID) (CustomFieldValue, error) {
	for _, value := range s.customFieldValues {
		if value.ID == valueID {
			return value, nil
		}
	}
	return CustomFieldValue{}, ErrCustomFieldValueNotFound
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

func (s *stubStore) CreateImportBatch(_ context.Context, batch ImportBatch) error {
	s.importBatches = append(s.importBatches, batch)
	return nil
}

func (s *stubStore) UpdateImportBatch(_ context.Context, batch ImportBatch) error {
	for i, existing := range s.importBatches {
		if existing.ID == batch.ID {
			s.importBatches[i] = batch
			return nil
		}
	}
	s.importBatches = append(s.importBatches, batch)
	return nil
}

func (s *stubStore) QueryImportBatches(context.Context, ImportBatchQueryFilter, order.By, page.Page) ([]ImportBatch, error) {
	return s.importBatches, nil
}

func (s *stubStore) CountImportBatches(context.Context, ImportBatchQueryFilter) (int, error) {
	return len(s.importBatches), nil
}

func (s *stubStore) QueryImportBatchByID(_ context.Context, batchID uuid.UUID) (ImportBatch, error) {
	for _, batch := range s.importBatches {
		if batch.ID == batchID {
			return batch, nil
		}
	}
	return ImportBatch{}, ErrImportBatchNotFound
}

func (s *stubStore) CreateImportInvalidRows(_ context.Context, rows []ImportInvalidRow) error {
	s.importInvalidRows = append(s.importInvalidRows, rows...)
	return nil
}

func (s *stubStore) QueryImportInvalidRows(context.Context, ImportInvalidRowQueryFilter, order.By, page.Page) ([]ImportInvalidRow, error) {
	return s.importInvalidRows, nil
}

func (s *stubStore) CountImportInvalidRows(context.Context, ImportInvalidRowQueryFilter) (int, error) {
	return len(s.importInvalidRows), nil
}

func (s *stubStore) QueryImportInvalidRowByID(_ context.Context, rowID uuid.UUID) (ImportInvalidRow, error) {
	for _, row := range s.importInvalidRows {
		if row.ID == rowID {
			return row, nil
		}
	}
	return ImportInvalidRow{}, ErrImportInvalidRowNotFound
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

func (s *stubStore) QuerySyncJobs(_ context.Context, filter SyncJobQueryFilter, _ order.By, _ page.Page) ([]SyncJob, error) {
	jobs := make([]SyncJob, 0, len(s.syncJobs))
	for _, job := range s.syncJobs {
		if filter.Adapter != nil && job.Adapter != *filter.Adapter {
			continue
		}
		if filter.Status != nil && job.Status != *filter.Status {
			continue
		}
		if filter.Direction != nil && job.Direction != *filter.Direction {
			continue
		}
		if filter.Retryable != nil && job.Retryable != *filter.Retryable {
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
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

func (s *stubStore) QuerySyncEvents(_ context.Context, filter SyncEventQueryFilter, _ order.By, _ page.Page) ([]SyncEvent, error) {
	events := make([]SyncEvent, 0, len(s.syncEvents))
	for _, event := range s.syncEvents {
		if filter.Adapter != nil && event.Adapter != *filter.Adapter {
			continue
		}
		if filter.Status != nil && event.Status != *filter.Status {
			continue
		}
		events = append(events, event)
	}

	return events, nil
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
