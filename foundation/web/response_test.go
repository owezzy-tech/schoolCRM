package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseTestModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m responseTestModel) Encode() ([]byte, string, error) {
	data, err := json.Marshal(m)
	return data, "application/json", err
}

type responseTestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e responseTestError) Error() string {
	return e.Message
}

func (e responseTestError) Encode() ([]byte, string, error) {
	data, err := json.Marshal(e)
	return data, "application/json", err
}

func (e responseTestError) HTTPStatus() int {
	return http.StatusBadRequest
}

type responseTestQuery struct {
	Items       []responseTestModel `json:"items"`
	Total       int                 `json:"total"`
	Page        int                 `json:"page"`
	RowsPerPage int                 `json:"rowsPerPage"`
}

func (q responseTestQuery) Encode() ([]byte, string, error) {
	data, err := json.Marshal(q)
	return data, "application/json", err
}

func TestRespondJSONAPISuccess(t *testing.T) {
	w := httptest.NewRecorder()
	resp := responseTestModel{ID: "constituent-1", Name: "Ada"}

	if err := Respond(context.Background(), w, resp); err != nil {
		t.Fatalf("respond: %v", err)
	}

	if got, exp := w.Code, http.StatusOK; got != exp {
		t.Fatalf("status got %d exp %d", got, exp)
	}

	if got, exp := w.Header().Get("Content-Type"), jsonAPIContentType; got != exp {
		t.Fatalf("content type got %q exp %q", got, exp)
	}

	var doc struct {
		JSONAPI jsonAPIObject `json:"jsonapi"`
		Data    struct {
			Type       string            `json:"type"`
			ID         string            `json:"id"`
			Attributes responseTestModel `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if doc.JSONAPI.Version != "1.1" {
		t.Fatalf("jsonapi version got %q", doc.JSONAPI.Version)
	}

	if doc.Data.Type != "responsetestmodel" || doc.Data.ID != "constituent-1" {
		t.Fatalf("resource got type %q id %q", doc.Data.Type, doc.Data.ID)
	}

	if doc.Data.Attributes.Name != resp.Name {
		t.Fatalf("attributes got %#v", doc.Data.Attributes)
	}

	if doc.Data.Attributes.ID != "" {
		t.Fatalf("attributes should not include id: %#v", doc.Data.Attributes)
	}
}

func TestRespondJSONAPIError(t *testing.T) {
	w := httptest.NewRecorder()
	resp := responseTestError{Code: "invalid_argument", Message: "email is required"}

	if err := Respond(context.Background(), w, resp); err != nil {
		t.Fatalf("respond: %v", err)
	}

	if got, exp := w.Code, http.StatusBadRequest; got != exp {
		t.Fatalf("status got %d exp %d", got, exp)
	}

	var doc struct {
		JSONAPI jsonAPIObject  `json:"jsonapi"`
		Errors  []jsonAPIError `json:"errors"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if doc.JSONAPI.Version != "1.1" {
		t.Fatalf("jsonapi version got %q", doc.JSONAPI.Version)
	}

	if len(doc.Errors) != 1 {
		t.Fatalf("errors got %d", len(doc.Errors))
	}

	got := doc.Errors[0]
	if got.Status != "400" || got.Code != "invalid_argument" || got.Detail != "email is required" {
		t.Fatalf("error got %#v", got)
	}
}

func TestRespondJSONAPICollection(t *testing.T) {
	w := httptest.NewRecorder()
	resp := responseTestQuery{
		Items: []responseTestModel{
			{ID: "constituent-1", Name: "Ada"},
			{ID: "constituent-2", Name: "Grace"},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}

	if err := Respond(context.Background(), w, resp); err != nil {
		t.Fatalf("respond: %v", err)
	}

	var doc struct {
		JSONAPI jsonAPIObject `json:"jsonapi"`
		Data    []struct {
			Type       string            `json:"type"`
			ID         string            `json:"id"`
			Attributes responseTestModel `json:"attributes"`
		} `json:"data"`
		Meta map[string]int `json:"meta"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if doc.JSONAPI.Version != "1.1" {
		t.Fatalf("jsonapi version got %q", doc.JSONAPI.Version)
	}

	if len(doc.Data) != 2 {
		t.Fatalf("data len got %d", len(doc.Data))
	}

	if doc.Data[0].Type != "responsetestquery" || doc.Data[0].ID != "constituent-1" {
		t.Fatalf("first resource got type %q id %q", doc.Data[0].Type, doc.Data[0].ID)
	}

	if doc.Data[0].Attributes.ID != "" {
		t.Fatalf("collection attributes should not include id: %#v", doc.Data[0].Attributes)
	}

	if doc.Meta["total"] != 2 || doc.Meta["page"] != 1 || doc.Meta["rowsPerPage"] != 10 {
		t.Fatalf("meta got %#v", doc.Meta)
	}
}

func TestRespondEncodeFailure(t *testing.T) {
	w := httptest.NewRecorder()
	resp := brokenEncoder{}

	if err := Respond(context.Background(), w, resp); err == nil {
		t.Fatal("expected error")
	}
}

type brokenEncoder struct{}

func (brokenEncoder) Encode() ([]byte, string, error) {
	return nil, "", errors.New("boom")
}
