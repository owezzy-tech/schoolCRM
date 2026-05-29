package admissionsdb

import (
	"encoding/json"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb/dbarray"
)

type staffProfileDB struct {
	ID          uuid.UUID      `db:"staff_profile_id"`
	UserID      uuid.UUID      `db:"user_id"`
	Roles       dbarray.String `db:"admissions_roles"`
	Active      bool           `db:"is_active"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

type constituentDB struct {
	ID              uuid.UUID  `db:"constituent_id"`
	FirstName       string     `db:"first_name"`
	LastName        string     `db:"last_name"`
	PreferredName   *string    `db:"preferred_name"`
	MiddleName      *string    `db:"middle_name"`
	Suffix          *string    `db:"suffix"`
	DateOfBirth     time.Time  `db:"date_of_birth"`
	PrimaryEmail    string     `db:"primary_email"`
	PrimaryPhone    string     `db:"primary_phone"`
	ExternalSISID   *string    `db:"external_sis_id"`
	LifecycleStage  string     `db:"lifecycle_stage"`
	DuplicateStatus string     `db:"duplicate_status"`
	DuplicateOfID   *uuid.UUID `db:"duplicate_of_id"`
	SISSyncedAt     *time.Time `db:"sis_synced_at"`
	DateCreated     time.Time  `db:"date_created"`
	DateUpdated     time.Time  `db:"date_updated"`
}

type programDB struct {
	ID            uuid.UUID  `db:"program_id"`
	ExternalSISID string     `db:"external_sis_id"`
	Name          string     `db:"name"`
	Code          string     `db:"code"`
	Description   *string    `db:"description"`
	DegreeLevel   *string    `db:"degree_level"`
	Active        bool       `db:"is_active"`
	SyncedAt      *time.Time `db:"synced_at"`
	DateCreated   time.Time  `db:"date_created"`
	DateUpdated   time.Time  `db:"date_updated"`
}

type duplicateReviewDB struct {
	ID                     uuid.UUID  `db:"duplicate_review_id"`
	SourceConstituentID    uuid.UUID  `db:"source_constituent_id"`
	CandidateConstituentID uuid.UUID  `db:"candidate_constituent_id"`
	MatchType              string     `db:"match_type"`
	MatchScore             int        `db:"match_score"`
	MatchReason            string     `db:"match_reason"`
	Status                 string     `db:"status"`
	ResolvedBy             *uuid.UUID `db:"resolved_by"`
	ResolvedAt             *time.Time `db:"resolved_at"`
	ResolutionNote         *string    `db:"resolution_note"`
	DateCreated            time.Time  `db:"date_created"`
	DateUpdated            time.Time  `db:"date_updated"`
}

type applicationDB struct {
	ID                 uuid.UUID  `db:"application_id"`
	ConstituentID      uuid.UUID  `db:"constituent_id"`
	ProgramID          uuid.UUID  `db:"program_id"`
	AcademicTermID     uuid.UUID  `db:"academic_term_id"`
	ApplicationType    string     `db:"application_type"`
	Status             string     `db:"status"`
	AssignedReviewerID *uuid.UUID `db:"assigned_reviewer_id"`
	SubmittedAt        *time.Time `db:"submitted_at"`
	DateCreated        time.Time  `db:"date_created"`
	DateUpdated        time.Time  `db:"date_updated"`
}

type applicationTransitionDB struct {
	ID            uuid.UUID       `db:"application_transition_id"`
	ApplicationID uuid.UUID       `db:"application_id"`
	FromStatus    string          `db:"from_status"`
	ToStatus      string          `db:"to_status"`
	ActorID       uuid.UUID       `db:"actor_id"`
	Reason        *string         `db:"reason"`
	Note          *string         `db:"note"`
	Metadata      json.RawMessage `db:"metadata"`
	DateCreated   time.Time       `db:"date_created"`
}

func toDBStaffProfile(bus admissionsbus.StaffProfile) staffProfileDB {
	return staffProfileDB{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Roles:       admissionsbus.AdmissionsRolesToStrings(bus.Roles),
		Active:      bus.Active,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusStaffProfile(db staffProfileDB) (admissionsbus.StaffProfile, error) {
	roles, err := admissionsbus.ParseAdmissionsRoles(db.Roles)
	if err != nil {
		return admissionsbus.StaffProfile{}, err
	}

	return admissionsbus.StaffProfile{
		ID:          db.ID,
		UserID:      db.UserID,
		Roles:       roles,
		Active:      db.Active,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}, nil
}

func toBusStaffProfiles(dbs []staffProfileDB) ([]admissionsbus.StaffProfile, error) {
	bus := make([]admissionsbus.StaffProfile, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusStaffProfile(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toDBConstituent(bus admissionsbus.Constituent) constituentDB {
	return constituentDB{
		ID:              bus.ID,
		FirstName:       bus.FirstName,
		LastName:        bus.LastName,
		PreferredName:   bus.PreferredName,
		MiddleName:      bus.MiddleName,
		Suffix:          bus.Suffix,
		DateOfBirth:     bus.DateOfBirth.UTC(),
		PrimaryEmail:    bus.PrimaryEmail.String(),
		PrimaryPhone:    bus.PrimaryPhone,
		ExternalSISID:   bus.ExternalSISID,
		LifecycleStage:  bus.LifecycleStage.String(),
		DuplicateStatus: bus.DuplicateStatus.String(),
		DuplicateOfID:   bus.DuplicateOfID,
		SISSyncedAt:     utcTimePtr(bus.SISSyncedAt),
		DateCreated:     bus.DateCreated.UTC(),
		DateUpdated:     bus.DateUpdated.UTC(),
	}
}

func toBusConstituent(db constituentDB) (admissionsbus.Constituent, error) {
	email, err := mail.ParseAddress(db.PrimaryEmail)
	if err != nil {
		return admissionsbus.Constituent{}, err
	}

	return admissionsbus.Constituent{
		ID:              db.ID,
		FirstName:       db.FirstName,
		LastName:        db.LastName,
		PreferredName:   db.PreferredName,
		MiddleName:      db.MiddleName,
		Suffix:          db.Suffix,
		DateOfBirth:     db.DateOfBirth.In(time.Local),
		PrimaryEmail:    *email,
		PrimaryPhone:    db.PrimaryPhone,
		ExternalSISID:   db.ExternalSISID,
		LifecycleStage:  admissionsbus.LifecycleStage(db.LifecycleStage),
		DuplicateStatus: admissionsbus.DuplicateStatus(db.DuplicateStatus),
		DuplicateOfID:   db.DuplicateOfID,
		SISSyncedAt:     localTimePtr(db.SISSyncedAt),
		DateCreated:     db.DateCreated.In(time.Local),
		DateUpdated:     db.DateUpdated.In(time.Local),
	}, nil
}

func toBusConstituents(dbs []constituentDB) ([]admissionsbus.Constituent, error) {
	bus := make([]admissionsbus.Constituent, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusConstituent(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

type academicTermDB struct {
	ID                   uuid.UUID  `db:"academic_term_id"`
	ExternalSISID        string     `db:"external_sis_id"`
	Name                 string     `db:"name"`
	Code                 string     `db:"code"`
	TermType             *string    `db:"term_type"`
	StartDate            time.Time  `db:"start_date"`
	EndDate              time.Time  `db:"end_date"`
	ApplicationStartDate *time.Time `db:"application_start_date"`
	ApplicationDeadline  *time.Time `db:"application_deadline"`
	Active               bool       `db:"is_active"`
	SyncedAt             *time.Time `db:"synced_at"`
	DateCreated          time.Time  `db:"date_created"`
	DateUpdated          time.Time  `db:"date_updated"`
}

func toDBProgram(bus admissionsbus.Program) programDB {
	return programDB{
		ID:            bus.ID,
		ExternalSISID: bus.ExternalSISID,
		Name:          bus.Name,
		Code:          bus.Code,
		Description:   bus.Description,
		DegreeLevel:   bus.DegreeLevel,
		Active:        bus.Active,
		SyncedAt:      utcTimePtr(bus.SyncedAt),
		DateCreated:   bus.DateCreated.UTC(),
		DateUpdated:   bus.DateUpdated.UTC(),
	}
}

func toBusProgram(db programDB) admissionsbus.Program {
	return admissionsbus.Program{
		ID:            db.ID,
		ExternalSISID: db.ExternalSISID,
		Name:          db.Name,
		Code:          db.Code,
		Description:   db.Description,
		DegreeLevel:   db.DegreeLevel,
		Active:        db.Active,
		SyncedAt:      localTimePtr(db.SyncedAt),
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}
}

func toBusPrograms(dbs []programDB) []admissionsbus.Program {
	bus := make([]admissionsbus.Program, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusProgram(db)
	}

	return bus
}

func toDBAcademicTerm(bus admissionsbus.AcademicTerm) academicTermDB {
	return academicTermDB{
		ID:                   bus.ID,
		ExternalSISID:        bus.ExternalSISID,
		Name:                 bus.Name,
		Code:                 bus.Code,
		TermType:             bus.TermType,
		StartDate:            bus.StartDate.UTC(),
		EndDate:              bus.EndDate.UTC(),
		ApplicationStartDate: utcTimePtr(bus.ApplicationStartDate),
		ApplicationDeadline:  utcTimePtr(bus.ApplicationDeadline),
		Active:               bus.Active,
		SyncedAt:             utcTimePtr(bus.SyncedAt),
		DateCreated:          bus.DateCreated.UTC(),
		DateUpdated:          bus.DateUpdated.UTC(),
	}
}

func toBusAcademicTerm(db academicTermDB) admissionsbus.AcademicTerm {
	return admissionsbus.AcademicTerm{
		ID:                   db.ID,
		ExternalSISID:        db.ExternalSISID,
		Name:                 db.Name,
		Code:                 db.Code,
		TermType:             db.TermType,
		StartDate:            db.StartDate.In(time.Local),
		EndDate:              db.EndDate.In(time.Local),
		ApplicationStartDate: localTimePtr(db.ApplicationStartDate),
		ApplicationDeadline:  localTimePtr(db.ApplicationDeadline),
		Active:               db.Active,
		SyncedAt:             localTimePtr(db.SyncedAt),
		DateCreated:          db.DateCreated.In(time.Local),
		DateUpdated:          db.DateUpdated.In(time.Local),
	}
}

func toBusAcademicTerms(dbs []academicTermDB) []admissionsbus.AcademicTerm {
	bus := make([]admissionsbus.AcademicTerm, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusAcademicTerm(db)
	}

	return bus
}

func toDBDuplicateReview(bus admissionsbus.DuplicateReview) duplicateReviewDB {
	return duplicateReviewDB{
		ID:                     bus.ID,
		SourceConstituentID:    bus.SourceConstituentID,
		CandidateConstituentID: bus.CandidateConstituentID,
		MatchType:              bus.MatchType.String(),
		MatchScore:             bus.MatchScore,
		MatchReason:            bus.MatchReason,
		Status:                 bus.Status.String(),
		ResolvedBy:             bus.ResolvedBy,
		ResolvedAt:             utcTimePtr(bus.ResolvedAt),
		ResolutionNote:         bus.ResolutionNote,
		DateCreated:            bus.DateCreated.UTC(),
		DateUpdated:            bus.DateUpdated.UTC(),
	}
}

func toBusDuplicateReview(db duplicateReviewDB) admissionsbus.DuplicateReview {
	return admissionsbus.DuplicateReview{
		ID:                     db.ID,
		SourceConstituentID:    db.SourceConstituentID,
		CandidateConstituentID: db.CandidateConstituentID,
		MatchType:              admissionsbus.DuplicateReviewMatchType(db.MatchType),
		MatchScore:             db.MatchScore,
		MatchReason:            db.MatchReason,
		Status:                 admissionsbus.DuplicateReviewStatus(db.Status),
		ResolvedBy:             db.ResolvedBy,
		ResolvedAt:             localTimePtr(db.ResolvedAt),
		ResolutionNote:         db.ResolutionNote,
		DateCreated:            db.DateCreated.In(time.Local),
		DateUpdated:            db.DateUpdated.In(time.Local),
	}
}

func toBusDuplicateReviews(dbs []duplicateReviewDB) []admissionsbus.DuplicateReview {
	bus := make([]admissionsbus.DuplicateReview, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusDuplicateReview(db)
	}

	return bus
}

func toDBApplication(bus admissionsbus.Application) applicationDB {
	return applicationDB{
		ID:                 bus.ID,
		ConstituentID:      bus.ConstituentID,
		ProgramID:          bus.ProgramID,
		AcademicTermID:     bus.AcademicTermID,
		ApplicationType:    bus.ApplicationType.String(),
		Status:             bus.Status.String(),
		AssignedReviewerID: bus.AssignedReviewerID,
		SubmittedAt:        utcTimePtr(bus.SubmittedAt),
		DateCreated:        bus.DateCreated.UTC(),
		DateUpdated:        bus.DateUpdated.UTC(),
	}
}

func toBusApplication(db applicationDB) admissionsbus.Application {
	return admissionsbus.Application{
		ID:                 db.ID,
		ConstituentID:      db.ConstituentID,
		ProgramID:          db.ProgramID,
		AcademicTermID:     db.AcademicTermID,
		ApplicationType:    admissionsbus.ApplicationType(db.ApplicationType),
		Status:             admissionsbus.ApplicationStatus(db.Status),
		AssignedReviewerID: db.AssignedReviewerID,
		SubmittedAt:        localTimePtr(db.SubmittedAt),
		DateCreated:        db.DateCreated.In(time.Local),
		DateUpdated:        db.DateUpdated.In(time.Local),
	}
}

func toBusApplications(dbs []applicationDB) []admissionsbus.Application {
	bus := make([]admissionsbus.Application, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusApplication(db)
	}

	return bus
}

func toDBApplicationTransition(bus admissionsbus.ApplicationTransition) applicationTransitionDB {
	return applicationTransitionDB{
		ID:            bus.ID,
		ApplicationID: bus.ApplicationID,
		FromStatus:    bus.FromStatus.String(),
		ToStatus:      bus.ToStatus.String(),
		ActorID:       bus.ActorID,
		Reason:        bus.Reason,
		Note:          bus.Note,
		Metadata:      json.RawMessage(bus.Metadata),
		DateCreated:   bus.DateCreated.UTC(),
	}
}

func toBusApplicationTransition(db applicationTransitionDB) admissionsbus.ApplicationTransition {
	return admissionsbus.ApplicationTransition{
		ID:            db.ID,
		ApplicationID: db.ApplicationID,
		FromStatus:    admissionsbus.ApplicationStatus(db.FromStatus),
		ToStatus:      admissionsbus.ApplicationStatus(db.ToStatus),
		ActorID:       db.ActorID,
		Reason:        db.Reason,
		Note:          db.Note,
		Metadata:      []byte(db.Metadata),
		DateCreated:   db.DateCreated.In(time.Local),
	}
}

func toBusApplicationTransitions(dbs []applicationTransitionDB) []admissionsbus.ApplicationTransition {
	bus := make([]admissionsbus.ApplicationTransition, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusApplicationTransition(db)
	}

	return bus
}

func utcTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	utc := t.UTC()
	return &utc
}

func localTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	local := t.In(time.Local)
	return &local
}
