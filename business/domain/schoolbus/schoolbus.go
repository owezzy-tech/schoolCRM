package schoolbus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/name"
)

// Storer defines persistence and a transaction shared with permission checks and audit.
type Storer interface {
	WithinTx(context.Context, func(Storer) error) error
	Actor(context.Context, uuid.UUID) (Actor, error)
	LockSchool(context.Context, uuid.UUID) error
	DepartmentExists(context.Context, Scope) (bool, error)
	HasCapability(context.Context, uuid.UUID, Scope, Capability) (bool, error)
	CreateSchool(context.Context, School) error
	CreateDepartment(context.Context, Department) error
	Grant(context.Context, Membership) (Membership, error)
	Membership(context.Context, uuid.UUID, uuid.UUID) (Membership, error)
	Revoke(context.Context, uuid.UUID, time.Time) error
	Schools(context.Context, uuid.UUID, bool) ([]School, error)
	Departments(context.Context, uuid.UUID) ([]Department, error)
	Memberships(context.Context, uuid.UUID) ([]Membership, error)
	Audit(context.Context, auditbus.Audit) error
}

type Business struct{ store Storer }

func NewBusiness(store Storer) *Business { return &Business{store: store} }

func (b *Business) CreateSchool(ctx context.Context, actorID uuid.UUID, label string) (School, error) {
	label, err := validateName(label)
	if err != nil {
		return School{}, err
	}
	school := School{ID: uuid.New(), Name: label, DateCreated: time.Now().UTC()}
	err = b.store.WithinTx(ctx, func(s Storer) error {
		actor, err := enabledActor(ctx, s, actorID)
		if err != nil {
			return err
		}
		if !actor.SuperAdmin {
			return ErrForbidden
		}
		if err := s.CreateSchool(ctx, school); err != nil {
			return err
		}
		return record(ctx, s, actorID, school.ID, "school_created", school, school.DateCreated)
	})
	return school, err
}

func (b *Business) CreateDepartment(ctx context.Context, actorID, schoolID uuid.UUID, label string) (Department, error) {
	label, err := validateName(label)
	if err != nil {
		return Department{}, err
	}
	department := Department{ID: uuid.New(), SchoolID: schoolID, Name: label, DateCreated: time.Now().UTC()}
	err = b.store.WithinTx(ctx, func(s Storer) error {
		if _, err := manageSchool(ctx, s, actorID, schoolID); err != nil {
			return err
		}
		if err := s.CreateDepartment(ctx, department); err != nil {
			return err
		}
		return record(ctx, s, actorID, department.ID, "department_created", department, department.DateCreated)
	})
	return department, err
}

// Grant assigns one capability. Only SUPER_ADMIN may delegate management.
func (b *Business) Grant(ctx context.Context, actorID, userID uuid.UUID, scope Scope, capability Capability) (Membership, error) {
	if err := validateScope(scope, capability); err != nil {
		return Membership{}, err
	}
	if userID == uuid.Nil {
		return Membership{}, ErrInvalid
	}
	membership := Membership{ID: uuid.New(), SchoolID: scope.SchoolID, DepartmentID: scope.DepartmentID,
		UserID: userID, Capability: capability, Active: true, DateUpdated: time.Now().UTC()}
	err := b.store.WithinTx(ctx, func(s Storer) error {
		actor, err := manageSchool(ctx, s, actorID, scope.SchoolID)
		if err != nil {
			return err
		}
		if capability == ManageMembers && !actor.SuperAdmin {
			return ErrForbidden
		}
		target, err := enabledActor(ctx, s, userID)
		if err != nil {
			return err
		}
		if capability == ManageMembers && !target.SchoolAdmin {
			return fmt.Errorf("%w: management requires SCHOOL_ADMIN role", ErrInvalid)
		}
		if scope.DepartmentID != nil {
			ok, err := s.DepartmentExists(ctx, scope)
			if err != nil {
				return err
			}
			if !ok {
				return ErrNotFound
			}
		}
		membership, err = s.Grant(ctx, membership)
		if err != nil {
			return err
		}
		return record(ctx, s, actorID, membership.ID, "membership_granted", membership, membership.DateUpdated)
	})
	return membership, err
}

func (b *Business) Revoke(ctx context.Context, actorID, schoolID, membershipID uuid.UUID) error {
	return b.store.WithinTx(ctx, func(s Storer) error {
		actor, err := manageSchool(ctx, s, actorID, schoolID)
		if err != nil {
			return err
		}
		membership, err := s.Membership(ctx, schoolID, membershipID)
		if err != nil {
			return err
		}
		if membership.Capability == ManageMembers && !actor.SuperAdmin {
			return ErrForbidden
		}
		if !membership.Active {
			return nil
		}
		now := time.Now().UTC()
		if err := s.Revoke(ctx, membership.ID, now); err != nil {
			return err
		}
		membership.Active = false
		membership.DateUpdated = now
		return record(ctx, s, actorID, membership.ID, "membership_revoked", membership, now)
	})
}

func (b *Business) Schools(ctx context.Context, actorID uuid.UUID) ([]School, error) {
	var schools []School
	err := b.store.WithinTx(ctx, func(s Storer) error {
		actor, err := enabledActor(ctx, s, actorID)
		if err != nil {
			return err
		}
		schools, err = s.Schools(ctx, actorID, actor.SuperAdmin)
		return err
	})
	return schools, err
}

func (b *Business) Departments(ctx context.Context, actorID, schoolID uuid.UUID) ([]Department, error) {
	var departments []Department
	err := b.store.WithinTx(ctx, func(s Storer) error {
		if err := readSchool(ctx, s, actorID, schoolID); err != nil {
			return err
		}
		var err error
		departments, err = s.Departments(ctx, schoolID)
		return err
	})
	return departments, err
}

func (b *Business) Memberships(ctx context.Context, actorID, schoolID uuid.UUID) ([]Membership, error) {
	var memberships []Membership
	err := b.store.WithinTx(ctx, func(s Storer) error {
		if _, err := manageSchool(ctx, s, actorID, schoolID); err != nil {
			return err
		}
		var err error
		memberships, err = s.Memberships(ctx, schoolID)
		return err
	})
	return memberships, err
}

func validateName(label string) (string, error) {
	label = strings.TrimSpace(label)
	if !utf8.ValidString(label) || label == "" || utf8.RuneCountInString(label) > 200 {
		return "", fmt.Errorf("%w: name must contain 1 to 200 characters", ErrInvalid)
	}
	return label, nil
}

func validateScope(scope Scope, capability Capability) error {
	if scope.SchoolID == uuid.Nil {
		return ErrInvalid
	}
	switch capability {
	case ManageMembers:
		if scope.DepartmentID != nil {
			return fmt.Errorf("%w: management is school-scoped", ErrInvalid)
		}
	case Teach, ReviewLessons, ApproveLessons:
		if scope.DepartmentID == nil || *scope.DepartmentID == uuid.Nil {
			return fmt.Errorf("%w: department is required", ErrInvalid)
		}
	default:
		return fmt.Errorf("%w: unknown capability", ErrInvalid)
	}
	return nil
}

func enabledActor(ctx context.Context, s Storer, actorID uuid.UUID) (Actor, error) {
	actor, err := s.Actor(ctx, actorID)
	if err != nil {
		return Actor{}, err
	}
	if !actor.Enabled {
		return Actor{}, ErrForbidden
	}
	return actor, nil
}

// manageSchool locks the school and returns the actor once they are SUPER_ADMIN
// or a SCHOOL_ADMIN holding that school's active management delegation.
func manageSchool(ctx context.Context, s Storer, actorID, schoolID uuid.UUID) (Actor, error) {
	if schoolID == uuid.Nil {
		return Actor{}, ErrInvalid
	}
	if err := s.LockSchool(ctx, schoolID); err != nil {
		return Actor{}, err
	}
	actor, err := enabledActor(ctx, s, actorID)
	if err != nil {
		return Actor{}, err
	}
	if actor.SuperAdmin {
		return actor, nil
	}
	if !actor.SchoolAdmin {
		return Actor{}, ErrForbidden
	}
	ok, err := s.HasCapability(ctx, actorID, Scope{SchoolID: schoolID}, ManageMembers)
	if err != nil {
		return Actor{}, err
	}
	if !ok {
		return Actor{}, ErrForbidden
	}
	return actor, nil
}

func readSchool(ctx context.Context, s Storer, actorID, schoolID uuid.UUID) error {
	if schoolID == uuid.Nil {
		return ErrInvalid
	}
	if err := s.LockSchool(ctx, schoolID); err != nil {
		return err
	}
	actor, err := enabledActor(ctx, s, actorID)
	if err != nil {
		return err
	}
	if actor.SuperAdmin {
		return nil
	}
	schools, err := s.Schools(ctx, actorID, false)
	if err != nil {
		return err
	}
	for _, school := range schools {
		if school.ID == schoolID {
			return nil
		}
	}
	return ErrForbidden
}

func record(ctx context.Context, s Storer, actorID, resourceID uuid.UUID, action string, resource any, now time.Time) error {
	data, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("marshal school audit: %w", err)
	}
	return s.Audit(ctx, auditbus.Audit{ID: uuid.New(), ObjID: resourceID, ObjDomain: domain.School,
		ObjName: name.MustParse("school access"), ActorID: actorID, Action: action, Data: data, Timestamp: now})
}
