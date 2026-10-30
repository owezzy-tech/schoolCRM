package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/stores/auditdb"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus/stores/userdb"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/name"
	"github.com/owezzy/schoolCRM/business/types/password"
	"github.com/owezzy/schoolCRM/business/types/role"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// BootstrapAdmin is an operator-only, empty-database initialization command.
// Credentials are supplied through masked environment fields, never arguments.
func BootstrapAdmin(log *logger.Logger, cfg sqldb.Config, nme, email, pass string) error {
	nu, err := bootstrapUser(nme, email, pass)
	if err != nil {
		return err
	}
	db, err := sqldb.Open(cfg)
	if err != nil {
		return errors.New("opening bootstrap database failed")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	id, err := bootstrapAdmin(ctx, log, db, nu)
	if err != nil {
		return err
	}
	fmt.Println("bootstrap administrator id:", id)
	return nil
}

func bootstrapUser(nme, email, pass string) (userbus.NewUser, error) {
	n, err := name.Parse(nme)
	if err != nil {
		return userbus.NewUser{}, errors.New("a valid bootstrap name is required")
	}
	addr, err := mail.ParseAddress(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return userbus.NewUser{}, errors.New("a valid bootstrap email is required")
	}
	p, err := password.Parse(pass)
	if err != nil || len(pass) < 16 {
		// The existing password parser includes rejected input in its error.
		// Never wrap that error or expose the operator's credential.
		return userbus.NewUser{}, errors.New("bootstrap password must be 16-19 supported characters")
	}
	return userbus.NewUser{Name: n, Email: *addr, Password: p, Roles: []role.Role{role.SuperAdmin}}, nil
}

func bootstrapAdmin(ctx context.Context, log *logger.Logger, db *sqlx.DB, nu userbus.NewUser) (uuid.UUID, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return uuid.Nil, errors.New("starting bootstrap transaction failed")
	}
	defer func() { _ = tx.Rollback() }()
	// Blocks ordinary concurrent user inserts too, not only bootstrap callers.
	if _, err := tx.ExecContext(ctx, "LOCK TABLE users IN SHARE ROW EXCLUSIVE MODE"); err != nil {
		return uuid.Nil, errors.New("locking initial user creation failed")
	}
	var count int
	if err := tx.GetContext(ctx, &count, "SELECT count(*) FROM users"); err != nil {
		return uuid.Nil, errors.New("checking bootstrap database failed")
	}
	if count != 0 {
		var id uuid.UUID
		const receipt = `SELECT u.user_id FROM users u JOIN audit a ON a.obj_id=u.user_id
		WHERE lower(u.email)=lower($1) AND u.enabled AND 'SUPER_ADMIN'=ANY(u.roles)
		AND a.obj_domain='USER' AND a.action='bootstrap_admin' LIMIT 1`
		if err := tx.GetContext(ctx, &id, receipt, nu.Email.Address); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return uuid.Nil, errors.New("bootstrap requires an empty user table; existing accounts are never changed")
			}
			return uuid.Nil, errors.New("checking bootstrap receipt failed")
		}
		// A matching committed receipt is replayed without changing passwords.
		return id, nil
	}
	bus, err := userbus.NewBusiness(log, nil, userdb.NewStore(log, db)).NewWithTx(tx)
	if err != nil {
		return uuid.Nil, errors.New("binding bootstrap user transaction failed")
	}
	usr, err := bus.Create(ctx, uuid.Nil, nu)
	if err != nil {
		return uuid.Nil, errors.New("creating bootstrap administrator failed")
	}
	store, err := auditdb.NewStore(log, db).NewWithTx(tx)
	if err != nil {
		return uuid.Nil, errors.New("binding bootstrap audit transaction failed")
	}
	_, err = auditbus.NewBusiness(log, store).Create(ctx, auditbus.NewAudit{
		ObjID: usr.ID, ObjDomain: domain.User, ObjName: usr.Name, ActorID: uuid.Nil,
		Action: "bootstrap_admin", Message: "initial administrator created by deployment operator",
		Data: map[string]any{"email": usr.Email.Address, "roles": []string{role.SuperAdmin.String()}},
	})
	if err != nil {
		return uuid.Nil, errors.New("creating bootstrap audit failed; account was not committed")
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.New("bootstrap commit outcome uncertain; retry the same administrator safely")
	}
	return usr.ID, nil
}
