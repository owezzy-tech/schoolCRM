package admissionsdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

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
