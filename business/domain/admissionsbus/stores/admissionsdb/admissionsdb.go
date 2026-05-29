// Package admissionsdb contains admissions related persistence functionality.
package admissionsdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Store manages the set of APIs for admissions database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (admissionsbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Health returns the current admissions store scaffold metadata.
func (s *Store) Health(ctx context.Context) (admissionsbus.Health, error) {
	s.log.Info(ctx, "admissions scaffold health", "status", "ready")

	return admissionsbus.Health{
		Context:    admissionsbus.DomainName,
		Status:     "ready",
		Aggregates: admissionsbus.AggregateNames(),
	}, nil
}
