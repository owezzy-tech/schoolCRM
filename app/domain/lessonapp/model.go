package lessonapp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
)

// NewPlanRequest is the content of a plan's first version.
type NewPlanRequest struct {
	Title         string          `json:"title"`
	Content       json.RawMessage `json:"content"`
	ChangeSummary string          `json:"changeSummary"`
	RequestID     uuid.UUID       `json:"requestID,omitempty"`
}

// RevisionRequest creates a new version from the plan's current BaseVersion.
type RevisionRequest struct {
	BaseVersion int `json:"baseVersion"`
	NewPlanRequest
}

func (r NewPlanRequest) toBus() lessonbus.Revision {
	return lessonbus.Revision{Title: r.Title, Content: r.Content, ChangeSummary: r.ChangeSummary, RequestID: r.RequestID}
}

type ReuseRequest struct {
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	RequestID uuid.UUID `json:"requestID"`
}

// DecisionRequest records an HOD review or dean approval.
type DecisionRequest struct {
	Decision lessonbus.Decision `json:"decision"`
	Feedback string             `json:"feedback"`
}

type LessonPlan lessonbus.Plan

func (p LessonPlan) Encode() ([]byte, string, error) {
	data, err := json.Marshal(p)
	return data, "application/json", err
}

// LessonVersion is identified by its plan and version number.
type LessonVersion struct {
	ID string `json:"id"`
	lessonbus.Version
}

func toLessonVersion(v lessonbus.Version) LessonVersion {
	return LessonVersion{ID: fmt.Sprintf("%s:%d", v.PlanID, v.Number), Version: v}
}

func (v LessonVersion) Encode() ([]byte, string, error) {
	data, err := json.Marshal(v)
	return data, "application/json", err
}
