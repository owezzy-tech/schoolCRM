// Package schooldb persists school access with transactional authorisation and audit.
package schooldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/stores/auditdb"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

type Store struct {
	db    *sqlx.DB
	tx    *sqlx.Tx
	audit *auditdb.Store
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{db: db, audit: auditdb.NewStore(log, db)}
}

func (s *Store) WithinTx(ctx context.Context, fn func(schoolbus.Storer) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin school transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	audit, err := s.audit.NewWithTx(tx)
	if err != nil {
		return err
	}
	if err := fn(&Store{db: s.db, tx: tx, audit: audit}); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit school transaction: %w", err)
	}
	return nil
}

func (s *Store) Actor(ctx context.Context, id uuid.UUID) (schoolbus.Actor, error) {
	var actor schoolbus.Actor
	err := s.tx.GetContext(ctx, &actor, `SELECT enabled, 'SUPER_ADMIN' = ANY(roles) AS super_admin, 'SCHOOL_ADMIN' = ANY(roles) AS school_admin FROM users WHERE user_id = $1 FOR SHARE`, id)
	return actor, translate(err)
}

// LockSchool serialises grants/revocations with other school-scoped operations.
func (s *Store) LockSchool(ctx context.Context, id uuid.UUID) error {
	var schoolID uuid.UUID
	return translate(s.tx.GetContext(ctx, &schoolID, `SELECT school_id FROM schools WHERE school_id = $1 FOR UPDATE`, id))
}

func (s *Store) DepartmentExists(ctx context.Context, scope schoolbus.Scope) (bool, error) {
	var exists bool
	err := s.tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM school_departments WHERE school_id = $1 AND department_id = $2)`, scope.SchoolID, scope.DepartmentID)
	return exists, err
}

func (s *Store) HasCapability(ctx context.Context, userID uuid.UUID, scope schoolbus.Scope, capability schoolbus.Capability) (bool, error) {
	var exists bool
	err := s.tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM school_memberships WHERE school_id = $1 AND user_id = $2 AND department_id IS NOT DISTINCT FROM $3 AND capability = $4 AND active)`, scope.SchoolID, userID, scope.DepartmentID, capability)
	return exists, err
}

func (s *Store) CreateSchool(ctx context.Context, school schoolbus.School) error {
	_, err := s.tx.ExecContext(ctx, `INSERT INTO schools (school_id, name, date_created) VALUES ($1,$2,$3)`, school.ID, school.Name, school.DateCreated)
	return err
}

func (s *Store) CreateDepartment(ctx context.Context, department schoolbus.Department) error {
	_, err := s.tx.ExecContext(ctx, `INSERT INTO school_departments (department_id, school_id, name, date_created) VALUES ($1,$2,$3,$4)`, department.ID, department.SchoolID, department.Name, department.DateCreated)
	return err
}

func (s *Store) Grant(ctx context.Context, membership schoolbus.Membership) (schoolbus.Membership, error) {
	const q = `INSERT INTO school_memberships (membership_id, school_id, department_id, user_id, capability, active, date_updated)
	VALUES ($1,$2,$3,$4,$5,true,$6)
	ON CONFLICT (school_id, department_id, user_id, capability)
	DO UPDATE SET active = true, date_updated = EXCLUDED.date_updated
	RETURNING membership_id, school_id, department_id, user_id, capability, active, date_updated`
	var result schoolbus.Membership
	err := s.tx.GetContext(ctx, &result, q, membership.ID, membership.SchoolID, membership.DepartmentID, membership.UserID, membership.Capability, membership.DateUpdated)
	return result, err
}

func (s *Store) Membership(ctx context.Context, schoolID, membershipID uuid.UUID) (schoolbus.Membership, error) {
	var membership schoolbus.Membership
	err := s.tx.GetContext(ctx, &membership, `SELECT membership_id, school_id, department_id, user_id, capability, active, date_updated FROM school_memberships WHERE school_id = $1 AND membership_id = $2`, schoolID, membershipID)
	return membership, translate(err)
}

func (s *Store) Revoke(ctx context.Context, membershipID uuid.UUID, now time.Time) error {
	_, err := s.tx.ExecContext(ctx, `UPDATE school_memberships SET active = false, date_updated = $2 WHERE membership_id = $1`, membershipID, now)
	return err
}

func (s *Store) Schools(ctx context.Context, actorID uuid.UUID, superAdmin bool) ([]schoolbus.School, error) {
	schools := []schoolbus.School{}
	err := s.tx.SelectContext(ctx, &schools, `SELECT school_id, name, date_created FROM schools s
	WHERE $2 OR EXISTS (SELECT 1 FROM school_memberships m WHERE m.school_id = s.school_id AND m.user_id = $1 AND m.active)
	ORDER BY name, school_id`, actorID, superAdmin)
	return schools, err
}

func (s *Store) Departments(ctx context.Context, schoolID uuid.UUID) ([]schoolbus.Department, error) {
	departments := []schoolbus.Department{}
	err := s.tx.SelectContext(ctx, &departments, `SELECT department_id, school_id, name, date_created FROM school_departments WHERE school_id = $1 ORDER BY name, department_id`, schoolID)
	return departments, err
}

func (s *Store) Memberships(ctx context.Context, schoolID uuid.UUID) ([]schoolbus.Membership, error) {
	memberships := []schoolbus.Membership{}
	err := s.tx.SelectContext(ctx, &memberships, `SELECT membership_id, school_id, department_id, user_id, capability, active, date_updated FROM school_memberships WHERE school_id = $1 ORDER BY date_updated, membership_id`, schoolID)
	return memberships, err
}

func (s *Store) Audit(ctx context.Context, audit auditbus.Audit) error {
	return s.audit.Create(ctx, audit)
}

func translate(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return schoolbus.ErrNotFound
	}
	return err
}
