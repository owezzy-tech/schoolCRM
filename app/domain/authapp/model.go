package authapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/types/role"
)

type token struct {
	Token string `json:"token"`
}

// Encode implements the encoder interface.
func (t token) Encode() ([]byte, string, error) {
	data, err := json.Marshal(t)
	return data, "application/json", err
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Decode implements the decoder interface.
func (lr *loginReq) Decode(data []byte) error {
	return json.Unmarshal(data, lr)
}

func (lr loginReq) Validate() error {
	if lr.Email == "" {
		return fmt.Errorf("email is required")
	}

	if lr.Password == "" {
		return fmt.Errorf("password is required")
	}

	if _, err := mail.ParseAddress(lr.Email); err != nil {
		return fmt.Errorf("email is invalid: %w", err)
	}

	return nil
}

type loginResp struct {
	AccessToken string    `json:"accessToken"`
	TokenType   string    `json:"tokenType"`
	ExpiresAt   time.Time `json:"expiresAt"`
	ExpiresIn   int       `json:"expiresIn"`
	User        loginUser `json:"user"`
}

// Encode implements the encoder interface.
func (lr loginResp) Encode() ([]byte, string, error) {
	data, err := json.Marshal(lr)
	return data, "application/json", err
}

type loginUser struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

func toLoginUser(usr userbus.User) loginUser {
	return loginUser{
		ID:    usr.ID.String(),
		Name:  usr.Name.String(),
		Email: usr.Email.Address,
		Roles: role.ParseToString(usr.Roles),
	}
}
