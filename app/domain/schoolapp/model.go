package schoolapp

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/schoolbus"
)

// NamedRequest is used for school and department creation.
type NamedRequest struct {
	Name string `json:"name"`
}

type GrantRequest struct {
	UserID       uuid.UUID            `json:"userID"`
	DepartmentID *uuid.UUID           `json:"departmentID,omitempty"`
	Capability   schoolbus.Capability `json:"capability"`
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
