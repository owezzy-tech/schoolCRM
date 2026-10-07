package lessonbus

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StructuredContent is the versioned contract for generated lesson snapshots.
// Legacy manually-authored snapshots retain their existing JSON contract.
type StructuredContent struct {
	SchemaVersion      int                `json:"schemaVersion"`
	Framework          string             `json:"framework"`
	Stage              string             `json:"stage"`
	Subject            string             `json:"subject"`
	CurriculumRevision string             `json:"curriculumRevision"`
	DurationMinutes    int                `json:"durationMinutes"`
	Objectives         []string           `json:"objectives"`
	Prerequisites      []string           `json:"prerequisites"`
	Materials          []string           `json:"materials"`
	Activities         []Activity         `json:"activities"`
	Differentiation    string             `json:"differentiation"`
	Assessment         string             `json:"assessment"`
	Citations          []Citation         `json:"citations"`
	Generation         *GenerationReceipt `json:"generation"`
	ReusedFrom         *ReuseOrigin       `json:"reusedFrom,omitempty"`
}

type Activity struct {
	Minutes int    `json:"minutes"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
}

type Citation struct {
	SourceID       uuid.UUID `json:"sourceID"`
	IndexRevision  uuid.UUID `json:"indexRevision"`
	ChunkOrdinal   int       `json:"chunkOrdinal"`
	Page           int       `json:"page"`
	Title          string    `json:"title"`
	Authority      string    `json:"authority"`
	SourceURL      string    `json:"sourceURL"`
	SourceSHA256   string    `json:"sourceSHA256"`
	Framework      string    `json:"framework"`
	Stage          string    `json:"stage"`
	Subject        string    `json:"subject"`
	Revision       string    `json:"revision"`
	EmbeddingModel string    `json:"embeddingModel"`
}

type GenerationReceipt struct {
	Provider     string    `json:"provider"`
	ModelID      string    `json:"modelID"`
	ModelVersion string    `json:"modelVersion"`
	CompletionID string    `json:"completionID"`
	GeneratedAt  time.Time `json:"generatedAt"`
}

type ReuseOrigin struct {
	PlanID  uuid.UUID `json:"planID"`
	Version int       `json:"version"`
}

func validateStructuredContent(raw json.RawMessage) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	if _, modern := fields["schemaVersion"]; !modern {
		return nil
	}
	var content StructuredContent
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&content); err != nil {
		return fmt.Errorf("%w: invalid structured lesson schema", ErrInvalid)
	}
	if content.SchemaVersion != 1 || content.DurationMinutes < 1 || content.DurationMinutes > 180 ||
		len(content.Objectives) == 0 || len(content.Activities) == 0 || len(content.Citations) == 0 ||
		content.Prerequisites == nil || content.Materials == nil || content.Generation == nil {
		return fmt.Errorf("%w: incomplete structured lesson", ErrInvalid)
	}
	switch content.Framework {
	case "kenya-cbc", "kenya-cbe", "kenya-8-4-4", "cambridge":
	default:
		return fmt.Errorf("%w: unsupported curriculum framework", ErrInvalid)
	}
	if strings.TrimSpace(content.Stage) == "" || strings.TrimSpace(content.Subject) == "" || strings.TrimSpace(content.CurriculumRevision) == "" ||
		strings.TrimSpace(content.Assessment) == "" || strings.TrimSpace(content.Differentiation) == "" {
		return fmt.Errorf("%w: structured lesson fields cannot be empty", ErrInvalid)
	}
	for _, values := range [][]string{content.Objectives, content.Prerequisites, content.Materials} {
		for _, text := range values {
			if strings.TrimSpace(text) == "" {
				return fmt.Errorf("%w: lesson list contains empty text", ErrInvalid)
			}
		}
	}
	total := 0
	for _, activity := range content.Activities {
		if activity.Minutes < 1 || activity.Minutes > 180 || strings.TrimSpace(activity.Title) == "" || strings.TrimSpace(activity.Detail) == "" {
			return fmt.Errorf("%w: invalid timed activity", ErrInvalid)
		}
		total += activity.Minutes
	}
	if total != content.DurationMinutes {
		return fmt.Errorf("%w: activity times must total lesson duration", ErrInvalid)
	}
	for _, citation := range content.Citations {
		sourceURL, err := url.Parse(citation.SourceURL)
		checksum, hashErr := hex.DecodeString(citation.SourceSHA256)
		if err != nil || sourceURL.Host == "" || (sourceURL.Scheme != "https" && sourceURL.Scheme != "http") ||
			hashErr != nil || len(checksum) != 32 || citation.SourceID == uuid.Nil || citation.IndexRevision == uuid.Nil || citation.ChunkOrdinal < 0 || citation.Page < 1 ||
			strings.TrimSpace(citation.Title) == "" || strings.TrimSpace(citation.Authority) == "" || strings.TrimSpace(citation.EmbeddingModel) == "" ||
			citation.Framework != content.Framework || citation.Stage != content.Stage || citation.Subject != content.Subject || citation.Revision != content.CurriculumRevision {
			return fmt.Errorf("%w: citation must identify the exact curriculum revision", ErrInvalid)
		}
	}
	model := content.Generation
	if model.Provider != "deepseek" || model.ModelID != "deepseek-flash" || model.ModelVersion != "DeepSeek-V4.1-Flash" ||
		model.CompletionID == "" || model.GeneratedAt.IsZero() {
		return fmt.Errorf("%w: invalid generation provenance", ErrInvalid)
	}
	if origin := content.ReusedFrom; origin != nil && (origin.PlanID == uuid.Nil || origin.Version < 1) {
		return fmt.Errorf("%w: invalid lesson reuse provenance", ErrInvalid)
	}
	return nil
}
