package lessonbus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/name"
)

// MaxContentBytes bounds one version's structured content.
const MaxContentBytes = 256 * 1024

// Storer persists lessons. Commands run in one transaction with the
// authority check and audit record.
type Storer interface {
	WithinTx(context.Context, func(Storer) error) error
	// ShareDepartment holds the school against concurrent membership changes
	// and reports ErrNotFound unless the department belongs to the school.
	ShareDepartment(ctx context.Context, schoolID, departmentID uuid.UUID) error
	// LockPlan locks the plan for update and holds its school like ShareDepartment.
	LockPlan(context.Context, uuid.UUID) (Plan, error)
	Authority(ctx context.Context, userID, schoolID, departmentID uuid.UUID) (Authority, error)
	CreatePlan(context.Context, Plan) error
	UpdatePlan(context.Context, Plan) error
	CreateVersion(context.Context, Version) error
	UpdateVersion(context.Context, Version) error
	Plan(context.Context, uuid.UUID) (Plan, error)
	Plans(ctx context.Context, schoolID, departmentID uuid.UUID) ([]Plan, error)
	Version(ctx context.Context, planID uuid.UUID, number int) (Version, error)
	Versions(context.Context, uuid.UUID) ([]Version, error)
	Audit(context.Context, auditbus.Audit) error
}

type Business struct{ store Storer }

func NewBusiness(store Storer) *Business { return &Business{store: store} }

var lessonRoles = []schoolbus.Capability{schoolbus.Teach, schoolbus.ReviewLessons, schoolbus.ApproveLessons}

// Create starts a plan whose first version is a draft by the teaching author.
func (b *Business) Create(ctx context.Context, actorID, schoolID, departmentID uuid.UUID, draft Revision) (Plan, error) {
	draft, err := validateDraft(draft)
	if err != nil {
		return Plan{}, err
	}
	now := time.Now().UTC()
	plan := Plan{ID: uuid.New(), SchoolID: schoolID, DepartmentID: departmentID, AuthorID: actorID,
		Title: draft.Title, Status: Draft, CurrentVersion: 1, DateCreated: now, DateUpdated: now}
	err = b.store.WithinTx(ctx, func(s Storer) error {
		if err := s.ShareDepartment(ctx, schoolID, departmentID); err != nil {
			return err
		}
		if err := authorise(ctx, s, actorID, plan.SchoolID, plan.DepartmentID, schoolbus.Teach); err != nil {
			return err
		}
		if err := s.CreatePlan(ctx, plan); err != nil {
			return err
		}
		version := newVersion(plan, draft, now)
		if err := s.CreateVersion(ctx, version); err != nil {
			return err
		}
		return record(ctx, s, actorID, "lesson_version_created", version, now)
	})
	return plan, err
}

// Revise stores a new draft version based on the current one. Earlier
// versions, their decisions and any live publication stay unchanged.
func (b *Business) Revise(ctx context.Context, actorID, planID uuid.UUID, base int, draft Revision) (Version, error) {
	draft, err := validateDraft(draft)
	if err != nil {
		return Version{}, err
	}
	var version Version
	err = b.store.WithinTx(ctx, func(s Storer) error {
		plan, err := s.LockPlan(ctx, planID)
		if err != nil {
			return err
		}
		if err := authorise(ctx, s, actorID, plan.SchoolID, plan.DepartmentID, schoolbus.Teach); err != nil {
			return err
		}
		if plan.AuthorID != actorID {
			return ErrForbidden
		}
		if base != plan.CurrentVersion {
			return ErrConflict
		}
		now := time.Now().UTC()
		plan.CurrentVersion++
		plan.DateUpdated = now
		version = newVersion(plan, draft, now)
		if err := s.CreateVersion(ctx, version); err != nil {
			return err
		}
		if err := s.UpdatePlan(ctx, plan); err != nil {
			return err
		}
		return record(ctx, s, actorID, "lesson_version_created", version, now)
	})
	return version, err
}

// Submit sends the author's current draft to HOD review, fixing its content
// for every later decision.
func (b *Business) Submit(ctx context.Context, actorID, planID uuid.UUID, number int) (Version, error) {
	return b.transition(ctx, actorID, planID, number, schoolbus.Teach, func(plan *Plan, v *Version, now time.Time) (string, error) {
		if v.AuthorID != actorID {
			return "", ErrForbidden
		}
		if v.Status != Draft {
			return "", ErrConflict
		}
		v.Status = HODReview
		return "lesson_submitted", nil
	})
}

// Review records the HOD decision. The reviewer cannot be the author.
func (b *Business) Review(ctx context.Context, actorID, planID uuid.UUID, number int, decision Decision, feedback string) (Version, error) {
	next, feedback, err := validateDecision(decision, feedback, DeanApproval)
	if err != nil {
		return Version{}, err
	}
	return b.transition(ctx, actorID, planID, number, schoolbus.ReviewLessons, func(plan *Plan, v *Version, now time.Time) (string, error) {
		if v.Status != HODReview {
			return "", ErrConflict
		}
		if v.AuthorID == actorID {
			return "", ErrForbidden
		}
		v.Status, v.ReviewerID, v.ReviewedAt, v.ReviewFeedback = next, &actorID, &now, feedback
		return "lesson_reviewed", nil
	})
}

// Approve records the dean decision. The dean cannot be the author or the HOD
// reviewer of this version.
func (b *Business) Approve(ctx context.Context, actorID, planID uuid.UUID, number int, decision Decision, feedback string) (Version, error) {
	next, feedback, err := validateDecision(decision, feedback, Approved)
	if err != nil {
		return Version{}, err
	}
	return b.transition(ctx, actorID, planID, number, schoolbus.ApproveLessons, func(plan *Plan, v *Version, now time.Time) (string, error) {
		if v.Status != DeanApproval {
			return "", ErrConflict
		}
		if v.AuthorID == actorID || *v.ReviewerID == actorID {
			return "", ErrForbidden
		}
		v.Status, v.ApproverID, v.ApprovedAt, v.ApprovalFeedback = next, &actorID, &now, feedback
		return "lesson_approved", nil
	})
}

// Publish makes the author's approved current version live. Repeating a
// completed publication changes nothing.
func (b *Business) Publish(ctx context.Context, actorID, planID uuid.UUID, number int) (Version, error) {
	return b.transition(ctx, actorID, planID, number, schoolbus.Teach, func(plan *Plan, v *Version, now time.Time) (string, error) {
		if v.AuthorID != actorID {
			return "", ErrForbidden
		}
		if v.Status == Published {
			return "", nil
		}
		if v.Status != Approved {
			return "", ErrConflict
		}
		v.Status, v.PublishedAt = Published, &now
		plan.PublishedVersion = &v.Number
		return "lesson_published", nil
	})
}

func (b *Business) Plan(ctx context.Context, actorID, planID uuid.UUID) (Plan, error) {
	var plan Plan
	err := b.store.WithinTx(ctx, func(s Storer) error {
		var err error
		if plan, err = s.Plan(ctx, planID); err != nil {
			return err
		}
		return authorise(ctx, s, actorID, plan.SchoolID, plan.DepartmentID, lessonRoles...)
	})
	return plan, err
}

func (b *Business) Plans(ctx context.Context, actorID, schoolID, departmentID uuid.UUID) ([]Plan, error) {
	var plans []Plan
	err := b.store.WithinTx(ctx, func(s Storer) error {
		if err := authorise(ctx, s, actorID, schoolID, departmentID, lessonRoles...); err != nil {
			return err
		}
		var err error
		plans, err = s.Plans(ctx, schoolID, departmentID)
		return err
	})
	return plans, err
}

func (b *Business) Versions(ctx context.Context, actorID, planID uuid.UUID) ([]Version, error) {
	var versions []Version
	err := b.store.WithinTx(ctx, func(s Storer) error {
		plan, err := s.Plan(ctx, planID)
		if err != nil {
			return err
		}
		if err := authorise(ctx, s, actorID, plan.SchoolID, plan.DepartmentID, lessonRoles...); err != nil {
			return err
		}
		versions, err = s.Versions(ctx, planID)
		return err
	})
	return versions, err
}

// transition applies one workflow step to the plan's current version. It
// locks the plan, rechecks current authority, rejects stale versions, and
// commits the version, plan pointer and audit record together. A step that
// returns no action is an idempotent repeat and writes nothing.
func (b *Business) transition(ctx context.Context, actorID, planID uuid.UUID, number int, capability schoolbus.Capability,
	step func(plan *Plan, v *Version, now time.Time) (string, error)) (Version, error) {
	var version Version
	err := b.store.WithinTx(ctx, func(s Storer) error {
		plan, err := s.LockPlan(ctx, planID)
		if err != nil {
			return err
		}
		if err := authorise(ctx, s, actorID, plan.SchoolID, plan.DepartmentID, capability); err != nil {
			return err
		}
		if number < 1 || number > plan.CurrentVersion {
			return ErrNotFound
		}
		if number != plan.CurrentVersion {
			return ErrConflict
		}
		if version, err = s.Version(ctx, planID, number); err != nil {
			return err
		}
		now := time.Now().UTC()
		action, err := step(&plan, &version, now)
		if err != nil || action == "" {
			return err
		}
		plan.DateUpdated = now
		if err := s.UpdateVersion(ctx, version); err != nil {
			return err
		}
		if err := s.UpdatePlan(ctx, plan); err != nil {
			return err
		}
		return record(ctx, s, actorID, action, version, now)
	})
	return version, err
}

func authorise(ctx context.Context, s Storer, actorID, schoolID, departmentID uuid.UUID, capabilities ...schoolbus.Capability) error {
	authority, err := s.Authority(ctx, actorID, schoolID, departmentID)
	if err != nil {
		return err
	}
	if !authority.has(capabilities...) {
		return ErrForbidden
	}
	return nil
}

func newVersion(plan Plan, draft Revision, now time.Time) Version {
	return Version{PlanID: plan.ID, Number: plan.CurrentVersion, Title: draft.Title, Content: draft.Content,
		ChangeSummary: draft.ChangeSummary, AuthorID: plan.AuthorID, Status: Draft, DateCreated: now}
}

func validateDraft(draft Revision) (Revision, error) {
	draft.Title = strings.TrimSpace(draft.Title)
	if !utf8.ValidString(draft.Title) || draft.Title == "" || utf8.RuneCountInString(draft.Title) > 200 {
		return Revision{}, fmt.Errorf("%w: title must contain 1 to 200 characters", ErrInvalid)
	}
	content := bytes.TrimSpace(draft.Content)
	if len(content) > MaxContentBytes || !json.Valid(content) || content[0] != '{' {
		return Revision{}, fmt.Errorf("%w: content must be a JSON object of at most 256 KiB", ErrInvalid)
	}
	draft.Content = content
	draft.ChangeSummary = strings.TrimSpace(draft.ChangeSummary)
	if !utf8.ValidString(draft.ChangeSummary) || utf8.RuneCountInString(draft.ChangeSummary) > 1000 {
		return Revision{}, fmt.Errorf("%w: change summary must be at most 1000 characters", ErrInvalid)
	}
	return draft, nil
}

// validateDecision returns the status an accepted decision moves to. Requests
// for changes must explain what to change.
func validateDecision(decision Decision, feedback string, accepted Status) (Status, string, error) {
	feedback = strings.TrimSpace(feedback)
	if !utf8.ValidString(feedback) || utf8.RuneCountInString(feedback) > 2000 {
		return "", "", fmt.Errorf("%w: feedback must be at most 2000 characters", ErrInvalid)
	}
	switch decision {
	case Accept:
		return accepted, feedback, nil
	case RequestChanges:
		if feedback == "" {
			return "", "", fmt.Errorf("%w: feedback is required when requesting changes", ErrInvalid)
		}
		return ChangesRequested, feedback, nil
	default:
		return "", "", fmt.Errorf("%w: unknown decision", ErrInvalid)
	}
}

// record audits the version's workflow metadata; content stays in the
// immutable version row it identifies.
func record(ctx context.Context, s Storer, actorID uuid.UUID, action string, version Version, now time.Time) error {
	version.Content = nil
	data, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("marshal lesson audit: %w", err)
	}
	return s.Audit(ctx, auditbus.Audit{ID: uuid.New(), ObjID: version.PlanID, ObjDomain: domain.Lesson,
		ObjName: name.MustParse("lesson plan"), ActorID: actorID, Action: action, Data: data, Timestamp: now})
}
