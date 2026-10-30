package lessonbus

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func testContent() StructuredContent {
	source := uuid.New()
	return StructuredContent{SchemaVersion: 1, Framework: "kenya-cbc", Stage: "grade-1", Subject: "english",
		CurriculumRevision: "2024", DurationMinutes: 30, Objectives: []string{"Identify classroom objects"},
		Prerequisites: []string{}, Materials: []string{"Picture cards"},
		Activities:      []Activity{{Minutes: 30, Title: "Name objects", Detail: "Practise naming a desk and chair"}},
		Differentiation: "Offer picture prompts", Assessment: "Ask learners to name three objects",
		Citations: []Citation{{SourceID: source, IndexRevision: source, Page: 17, Title: "English",
			Authority: "KICD", SourceURL: "https://kicd.ac.ke/english.pdf", SourceSHA256: strings.Repeat("a", 64),
			Framework: "kenya-cbc", Stage: "grade-1", Subject: "english", Revision: "2024", EmbeddingModel: "ollama:bge-m3:pinned"}},
		Generation: &GenerationReceipt{Provider: "deepseek", ModelID: "deepseek-flash", ModelVersion: "DeepSeek-V4.1-Flash",
			CompletionID: "completion-1", GeneratedAt: time.Now().UTC()}}
}

func TestStructuredLessonBoundary(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*StructuredContent)
	}{
		{"wrong total duration", func(c *StructuredContent) { c.Activities[0].Minutes = 29 }},
		{"wrong revision", func(c *StructuredContent) { c.Citations[0].Revision = "2025" }},
		{"wrong framework", func(c *StructuredContent) { c.Citations[0].Framework = "cambridge" }},
		{"missing evidence", func(c *StructuredContent) { c.Citations = nil }},
		{"invalid source", func(c *StructuredContent) { c.Citations[0].SourceID = uuid.Nil }},
		{"invalid URL", func(c *StructuredContent) { c.Citations[0].SourceURL = "javascript:alert(1)" }},
		{"wrong model", func(c *StructuredContent) { c.Generation.ModelID = "different" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := testContent()
			tc.change(&c)
			raw, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}
			if err := validateStructuredContent(raw); !errors.Is(err, ErrInvalid) {
				t.Fatalf("got %v, want ErrInvalid", err)
			}
		})
	}
	valid := testContent()
	raw, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateStructuredContent(raw); err != nil {
		t.Fatal(err)
	}
	if err := validateStructuredContent(json.RawMessage(`{"schemaVersion":null}`)); !errors.Is(err, ErrInvalid) {
		t.Fatalf("null schema accepted: %v", err)
	}
	if err := validateStructuredContent(json.RawMessage(`{"objectives":["Legacy"]}`)); err != nil {
		t.Fatalf("legacy compatibility: %v", err)
	}
}
