package lessonbus

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
)

// Reuse creates an author-owned draft from an exact permitted source version.
// Source decisions and publication pointers are never copied.
func (b *Business) Reuse(ctx context.Context, actorID, sourceID uuid.UUID, number int, title string, requestID uuid.UUID) (Plan, error) {
	if number < 1 || requestID == uuid.Nil {
		return Plan{}, ErrInvalid
	}
	var result Plan
	err := b.store.WithinTx(ctx, func(s Storer) error {
		source, err := s.LockPlan(ctx, sourceID)
		if err != nil {
			return err
		}
		if err := authorise(ctx, s, actorID, source.SchoolID, source.DepartmentID, schoolbus.Teach); err != nil {
			return err
		}
		version, err := s.Version(ctx, sourceID, number)
		if err != nil {
			return err
		}
		if version.Status == Draft && version.AuthorID != actorID {
			return ErrNotFound
		}
		if title == "" {
			title = version.Title
		}
		var content map[string]json.RawMessage
		if err := json.Unmarshal(version.Content, &content); err != nil {
			return fmt.Errorf("decode reused content: %w", err)
		}
		provenance, err := json.Marshal(struct {
			PlanID  uuid.UUID `json:"planID"`
			Version int       `json:"version"`
		}{sourceID, number})
		if err != nil {
			return err
		}
		content["reusedFrom"] = provenance
		snapshot, err := json.Marshal(content)
		if err != nil {
			return err
		}
		draft, err := validateDraft(Revision{Title: title, Content: snapshot,
			ChangeSummary: fmt.Sprintf("Reused from plan %s, version %d", sourceID, number), RequestID: requestID})
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		result = Plan{ID: uuid.New(), SchoolID: source.SchoolID, DepartmentID: source.DepartmentID,
			AuthorID: actorID, Title: draft.Title, Status: Draft, CurrentVersion: 1,
			DateCreated: now, DateUpdated: now}
		result, err = createPlan(ctx, s, result, draft, sourceID, number, now)
		return err
	})
	return result, err
}

// createPlan runs inside the caller's already-authorised transaction.
func createPlan(ctx context.Context, s Storer, plan Plan, draft Revision, sourceID uuid.UUID, sourceVersion int, now time.Time) (Plan, error) {
	inputHash, err := creationHash(plan, draft, sourceID, sourceVersion)
	if err != nil {
		return Plan{}, err
	}
	if draft.RequestID != uuid.Nil {
		previous, err := s.Creation(ctx, plan.AuthorID, draft.RequestID)
		if err == nil {
			if previous.InputHash != inputHash {
				return Plan{}, fmt.Errorf("%w: requestID already used with different input", ErrConflict)
			}
			return previous.Result, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return Plan{}, err
		}
	}
	if err := s.CreatePlan(ctx, plan); err != nil {
		return Plan{}, err
	}
	version := newVersion(plan, draft, now)
	if err := s.CreateVersion(ctx, version); err != nil {
		return Plan{}, err
	}
	if err := record(ctx, s, plan.AuthorID, "lesson_version_created", version, now); err != nil {
		return Plan{}, err
	}
	if draft.RequestID != uuid.Nil {
		if err := s.RecordCreation(ctx, Creation{ActorID: plan.AuthorID, RequestID: draft.RequestID, InputHash: inputHash, Result: plan}); err != nil {
			return Plan{}, err
		}
	}
	return plan, nil
}

func creationHash(plan Plan, draft Revision, sourceID uuid.UUID, sourceVersion int) (string, error) {
	// Canonical object ordering survives JSONB checkpoint round trips.
	var content any
	decoder := json.NewDecoder(bytes.NewReader(draft.Content))
	decoder.UseNumber()
	if err := decoder.Decode(&content); err != nil {
		return "", err
	}
	input := struct {
		SchoolID      uuid.UUID
		DepartmentID  uuid.UUID
		Title         string
		ChangeSummary string
		Content       any
		SourceID      uuid.UUID
		SourceVersion int
	}{plan.SchoolID, plan.DepartmentID, draft.Title, draft.ChangeSummary, content, sourceID, sourceVersion}
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
