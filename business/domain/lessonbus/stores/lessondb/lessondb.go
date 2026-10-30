// Package lessondb persists lesson plans, their immutable versions and audit
// records in one transaction.
package lessondb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/stores/auditdb"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
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

func (s *Store) WithinTx(ctx context.Context, fn func(lessonbus.Storer) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin lesson transaction: %w", err)
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
		return fmt.Errorf("commit lesson transaction: %w", err)
	}
	return nil
}

// planColumns reads a plan with the title and status of its current version.
const planColumns = `p.plan_id, p.school_id, p.department_id, p.author_id, v.title, v.status,
	p.current_version, p.published_version, p.date_created, p.date_updated
	FROM lesson_plans p JOIN lesson_plan_versions v ON v.plan_id = p.plan_id AND v.version_number = p.current_version`

// visiblePlanColumns reads a plan as viewer $1 sees it: its newest version
// that is not another author's unsubmitted draft. Plans with no such version
// produce no row.
const visiblePlanColumns = `p.plan_id, p.school_id, p.department_id, p.author_id, v.title, v.status,
	v.version_number AS current_version, p.published_version, p.date_created, p.date_updated
	FROM lesson_plans p CROSS JOIN LATERAL (
		SELECT title, status, version_number FROM lesson_plan_versions v
		WHERE v.plan_id = p.plan_id AND (v.status <> 'draft' OR v.author_id = $1)
		ORDER BY v.version_number DESC LIMIT 1) v`

const versionColumns = `plan_id, version_number, title, content, change_summary, author_id, status, reviewer_id, reviewed_at,
	approver_id, approved_at, published_at, review_feedback, approval_feedback, date_created`

// ShareDepartment takes a shared school lock, which membership grants and
// revocations (FOR UPDATE) must wait for, and vice versa.
func (s *Store) ShareDepartment(ctx context.Context, schoolID, departmentID uuid.UUID) error {
	var id uuid.UUID
	return translate(s.tx.GetContext(ctx, &id, `SELECT s.school_id FROM schools s
	JOIN school_departments d ON d.school_id = s.school_id AND d.department_id = $2
	WHERE s.school_id = $1 FOR SHARE OF s`, schoolID, departmentID))
}

func (s *Store) LockPlan(ctx context.Context, planID uuid.UUID) (lessonbus.Plan, error) {
	var plan lessonbus.Plan
	err := s.tx.GetContext(ctx, &plan, `SELECT `+planColumns+`
	JOIN schools s ON s.school_id = p.school_id
	WHERE p.plan_id = $1 FOR UPDATE OF p FOR SHARE OF s`, planID)
	return plan, translate(err)
}

// Authority reads the enabled flag under a shared user lock, so disabling the
// account waits for this transaction, then the active department capabilities.
func (s *Store) Authority(ctx context.Context, userID, schoolID, departmentID uuid.UUID) (lessonbus.Authority, error) {
	var authority lessonbus.Authority
	if err := s.tx.GetContext(ctx, &authority.Enabled, `SELECT enabled FROM users WHERE user_id = $1 FOR SHARE`, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authority, nil
		}
		return authority, err
	}
	err := s.tx.SelectContext(ctx, &authority.Capabilities, `SELECT capability FROM school_memberships
	WHERE user_id = $1 AND school_id = $2 AND department_id = $3 AND active`, userID, schoolID, departmentID)
	return authority, err
}

func (s *Store) CreatePlan(ctx context.Context, plan lessonbus.Plan) error {
	_, err := s.tx.ExecContext(ctx, `INSERT INTO lesson_plans
	(plan_id, school_id, department_id, author_id, current_version, published_version, date_created, date_updated)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, plan.ID, plan.SchoolID, plan.DepartmentID, plan.AuthorID,
		plan.CurrentVersion, plan.PublishedVersion, plan.DateCreated, plan.DateUpdated)
	return err
}

func (s *Store) UpdatePlan(ctx context.Context, plan lessonbus.Plan) error {
	_, err := s.tx.ExecContext(ctx, `UPDATE lesson_plans SET current_version = $2, published_version = $3, date_updated = $4
	WHERE plan_id = $1`, plan.ID, plan.CurrentVersion, plan.PublishedVersion, plan.DateUpdated)
	return err
}

func (s *Store) CreateVersion(ctx context.Context, v lessonbus.Version) error {
	_, err := s.tx.ExecContext(ctx, `INSERT INTO lesson_plan_versions (`+versionColumns+`)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`, v.PlanID, v.Number, v.Title, []byte(v.Content),
		v.ChangeSummary, v.AuthorID, v.Status, v.ReviewerID, v.ReviewedAt, v.ApproverID, v.ApprovedAt, v.PublishedAt,
		v.ReviewFeedback, v.ApprovalFeedback, v.DateCreated)
	return err
}

// UpdateVersion writes workflow decisions only; content is immutable.
func (s *Store) UpdateVersion(ctx context.Context, v lessonbus.Version) error {
	_, err := s.tx.ExecContext(ctx, `UPDATE lesson_plan_versions SET status = $3, reviewer_id = $4, reviewed_at = $5,
	approver_id = $6, approved_at = $7, published_at = $8, review_feedback = $9, approval_feedback = $10
	WHERE plan_id = $1 AND version_number = $2`, v.PlanID, v.Number, v.Status, v.ReviewerID, v.ReviewedAt,
		v.ApproverID, v.ApprovedAt, v.PublishedAt, v.ReviewFeedback, v.ApprovalFeedback)
	return err
}

func (s *Store) Plan(ctx context.Context, planID, viewerID uuid.UUID) (lessonbus.Plan, error) {
	var plan lessonbus.Plan
	err := s.tx.GetContext(ctx, &plan, `SELECT `+visiblePlanColumns+` WHERE p.plan_id = $2`, viewerID, planID)
	return plan, translate(err)
}

func (s *Store) Plans(ctx context.Context, schoolID, departmentID, viewerID uuid.UUID) ([]lessonbus.Plan, error) {
	plans := []lessonbus.Plan{}
	err := s.tx.SelectContext(ctx, &plans, `SELECT `+visiblePlanColumns+`
	WHERE p.school_id = $2 AND p.department_id = $3 ORDER BY p.date_updated DESC, p.plan_id`, viewerID, schoolID, departmentID)
	return plans, err
}

func (s *Store) Version(ctx context.Context, planID uuid.UUID, number int) (lessonbus.Version, error) {
	var version lessonbus.Version
	err := s.tx.GetContext(ctx, &version, `SELECT `+versionColumns+` FROM lesson_plan_versions
	WHERE plan_id = $1 AND version_number = $2`, planID, number)
	return version, translate(err)
}

func (s *Store) Versions(ctx context.Context, planID, viewerID uuid.UUID) ([]lessonbus.Version, error) {
	versions := []lessonbus.Version{}
	err := s.tx.SelectContext(ctx, &versions, `SELECT `+versionColumns+` FROM lesson_plan_versions
	WHERE plan_id = $1 AND (status <> 'draft' OR author_id = $2) ORDER BY version_number DESC`, planID, viewerID)
	return versions, err
}

func (s *Store) Audit(ctx context.Context, audit auditbus.Audit) error {
	return s.audit.Create(ctx, audit)
}

func translate(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return lessonbus.ErrNotFound
	}
	return err
}
