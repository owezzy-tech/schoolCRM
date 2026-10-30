package commands

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/types/role"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"golang.org/x/crypto/bcrypt"
)

const fixturePassword = "SecretFixture12345!"

type bootstrapHash struct{}

func (bootstrapHash) Match(value driver.Value) bool {
	hash, ok := value.([]byte)
	return ok && bcrypt.CompareHashAndPassword(hash, []byte(fixturePassword)) == nil
}

type credentialFreeAudit struct{}

func (credentialFreeAudit) Match(value driver.Value) bool {
	var text string
	switch v := value.(type) {
	case string:
		text = v
	case []byte:
		text = string(v)
	default:
		return false
	}
	return strings.Contains(text, "SUPER_ADMIN") && !strings.Contains(text, fixturePassword) && !strings.Contains(text, "Password")
}

func TestBootstrapValidationMasksPasswords(t *testing.T) {
	for _, pass := range []string{"short", "A_Long_Unsupported_Secret!"} {
		_, err := bootstrapUser("Owen Adirah", "owner@example.invalid", pass)
		if err == nil || strings.Contains(err.Error(), pass) {
			t.Fatal("invalid credential was accepted or exposed")
		}
	}
	nu, err := bootstrapUser("Owen Adirah", "Owner@Example.Invalid", fixturePassword)
	if err != nil || len(nu.Roles) != 1 || nu.Roles[0] != role.SuperAdmin || nu.Email.Address != "owner@example.invalid" {
		t.Fatal("bootstrap identity or role contract failed")
	}
}

func TestBootstrapTransactions(t *testing.T) {
	for _, scenario := range []string{"create", "audit_failure", "user_failure", "commit_failure", "existing_user", "replay", "lock_failure"} {
		t.Run(scenario, func(t *testing.T) {
			raw, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			db := sqlx.NewDb(raw, "pgx")
			defer db.Close()
			log := logger.New(io.Discard, logger.LevelInfo, "TEST", func(context.Context) string { return "test" })
			nu, err := bootstrapUser("Owen Adirah", "owner@example.invalid", fixturePassword)
			if err != nil {
				t.Fatal("fixture validation failed")
			}
			mock.ExpectBegin()
			lock := mock.ExpectExec("LOCK TABLE users IN SHARE ROW EXCLUSIVE MODE")
			if scenario == "lock_failure" {
				lock.WillReturnError(errors.New("lock unavailable"))
				mock.ExpectRollback()
			} else {
				lock.WillReturnResult(sqlmock.NewResult(0, 0))
				count := 0
				if scenario == "existing_user" || scenario == "replay" {
					count = 3
				}
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
				if count != 0 {
					rows := sqlmock.NewRows([]string{"user_id"})
					if scenario == "replay" {
						rows.AddRow("11111111-1111-4111-8111-111111111111")
					}
					mock.ExpectQuery("SELECT u.user_id FROM users u JOIN audit").WithArgs(nu.Email.Address).WillReturnRows(rows)
					mock.ExpectRollback()
				} else {
					insert := mock.ExpectExec("INSERT INTO users").WithArgs(
						sqlmock.AnyArg(), "Owen Adirah", "owner@example.invalid", bootstrapHash{},
						sqlmock.AnyArg(), nil, true, sqlmock.AnyArg(), sqlmock.AnyArg(),
					)
					if scenario == "user_failure" {
						insert.WillReturnError(errors.New("insert unavailable"))
						mock.ExpectRollback()
					} else {
						insert.WillReturnResult(sqlmock.NewResult(0, 1))
						audit := mock.ExpectExec("INSERT INTO audit").WithArgs(
							sqlmock.AnyArg(), sqlmock.AnyArg(), "USER", "Owen Adirah", uuid.Nil.String(),
							"bootstrap_admin", credentialFreeAudit{}, "initial administrator created by deployment operator", sqlmock.AnyArg(),
						)
						if scenario == "audit_failure" {
							audit.WillReturnError(errors.New("audit unavailable"))
							mock.ExpectRollback()
						} else {
							audit.WillReturnResult(sqlmock.NewResult(0, 1))
							commit := mock.ExpectCommit()
							if scenario == "commit_failure" {
								commit.WillReturnError(errors.New("reply unavailable"))
							}
						}
					}
				}
			}
			id, err := bootstrapAdmin(context.Background(), log, db, nu)
			wantSuccess := scenario == "create" || scenario == "replay"
			if (err == nil) != wantSuccess || (wantSuccess && id == uuid.Nil) {
				t.Fatalf("unexpected transaction result for %s: %v", scenario, err)
			}
			if err != nil && strings.Contains(err.Error(), fixturePassword) {
				t.Fatal("credential leaked in transaction failure")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
