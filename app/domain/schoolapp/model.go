package schoolapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
)

// NamedRequest is used for school and department creation.
type NamedRequest struct {
	Name string `json:"name"`
}

func (r *NamedRequest) Decode(data []byte) error { return decode(data, r) }

type GrantRequest struct {
	UserID       uuid.UUID            `json:"userID"`
	DepartmentID *uuid.UUID           `json:"departmentID,omitempty"`
	Capability   schoolbus.Capability `json:"capability"`
}

func (r *GrantRequest) Decode(data []byte) error { return decode(data, r) }

func decode(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("expected a single JSON object")
	}
	return nil
}

type School schoolbus.School

func (s School) Encode() ([]byte, string, error) {
	data, err := json.Marshal(s)
	return data, "application/json", err
}

type Department schoolbus.Department

func (d Department) Encode() ([]byte, string, error) {
	data, err := json.Marshal(d)
	return data, "application/json", err
}

type Membership schoolbus.Membership

func (m Membership) Encode() ([]byte, string, error) {
	data, err := json.Marshal(m)
	return data, "application/json", err
}
