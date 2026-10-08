package commands

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus/stores/admissionsdb"
	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/stores/auditdb"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus/stores/userdb"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/business/types/domain"
	"github.com/owezzy/schoolCRM/business/types/role"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// AdmissionsAdmin provisions the initial admissions context for an enabled
// platform administrator. It never changes global roles or existing profiles.
func AdmissionsAdmin(log *logger.Logger, cfg sqldb.Config, email string) error {
	addr, err := mail.ParseAddress(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return errors.New("a valid administrator email is required")
	}
	db, err := sqldb.Open(cfg)
	if err != nil {
		return errors.New("opening admissions provisioning database failed")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	id, err := provisionAdmissionsAdmin(ctx, log, db, addr.Address)
	if err != nil {
		return err
	}
	fmt.Println("admissions administrator profile id:", id)
	return nil
}

func provisionAdmissionsAdmin(ctx context.Context, log *logger.Logger, db *sqlx.DB, email string) (uuid.UUID, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return uuid.Nil, errors.New("starting admissions provisioning transaction failed")
	}
	defer func() { _ = tx.Rollback() }()
	// Serialize provisioning and hold the identity stable through the grant.
	var userID uuid.UUID
	if err := tx.GetContext(ctx, &userID, "SELECT user_id FROM users WHERE lower(email)=lower($1) FOR UPDATE", email); err != nil {
		return uuid.Nil, errors.New("existing administrator identity required")
	}
	users, err := userbus.NewBusiness(log, nil, userdb.NewStore(log, db)).NewWithTx(tx)
	if err != nil {
		return uuid.Nil, errors.New("binding identity transaction failed")
	}
	usr, err := users.QueryByID(ctx, userID)
	if err != nil || !usr.Enabled {
		return uuid.Nil, errors.New("enabled platform administrator required")
	}
	admin := false
	for _, r := range usr.Roles {
		admin = admin || r.Equal(role.SuperAdmin)
	}
	if !admin {
		return uuid.Nil, errors.New("SUPER_ADMIN identity required; global roles are never changed")
	}
	bus, err := admissionsbus.NewBusiness(log, nil, admissionsdb.NewStore(log, db)).NewWithTx(tx)
	if err != nil {
		return uuid.Nil, errors.New("binding admissions transaction failed")
	}
	profile, err := bus.QueryStaffProfileByUserID(ctx, userID)
	if err == nil {
		if profile.Active && len(profile.Roles) == 1 && profile.Roles[0] == admissionsbus.AdmissionsRoleAdmin {
			return profile.ID, nil
		}
		return uuid.Nil, errors.New("existing admissions profile differs; use supported staff management")
	}
	if !errors.Is(err, admissionsbus.ErrStaffProfileNotFound) {
		return uuid.Nil, errors.New("checking existing admissions profile failed")
	}
	profile, err = bus.CreateStaffProfile(ctx, admissionsbus.NewStaffProfile{
		UserID: userID, Roles: []admissionsbus.AdmissionsRole{admissionsbus.AdmissionsRoleAdmin}, Active: true,
	})
	if err != nil {
		return uuid.Nil, errors.New("creating admissions administrator profile failed")
	}
	store, err := auditdb.NewStore(log, db).NewWithTx(tx)
	if err != nil {
		return uuid.Nil, errors.New("binding admissions audit transaction failed")
	}
	_, err = auditbus.NewBusiness(log, store).Create(ctx, auditbus.NewAudit{
		ObjID: profile.ID, ObjDomain: domain.Admissions, ObjName: usr.Name, ActorID: uuid.Nil,
		Action: "bootstrap_admissions_admin", Message: "admissions administrator provisioned by deployment operator",
		Data: map[string]any{"userID": userID, "roles": []string{admissionsbus.AdmissionsRoleAdmin.String()}},
	})
	if err != nil {
		return uuid.Nil, errors.New("creating admissions audit failed; profile was not committed")
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.New("admissions provisioning commit outcome uncertain; retry the same identity safely")
	}
	return profile.ID, nil
}
