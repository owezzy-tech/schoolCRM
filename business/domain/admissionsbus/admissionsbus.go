// Package admissionsbus provides business access to the admissions domain.
package admissionsbus

import (
	"context"
	"fmt"

	"github.com/owezzy/schoolCRM/business/sdk/delegate"
	"github.com/owezzy/schoolCRM/business/sdk/sqldb"
	"github.com/owezzy/schoolCRM/foundation/logger"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve admissions data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Health(ctx context.Context) (Health, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Health(ctx context.Context) (Health, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for admissions access.
type Business struct {
	log        *logger.Logger
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs an admissions business API for use.
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		delegate:   delegate,
		storer:     storer,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, b.delegate, storer, b.extensions...)

	return nb, nil
}

// Health returns the current admissions context scaffold metadata.
func (b *Business) Health(ctx context.Context) (Health, error) {
	health, err := b.storer.Health(ctx)
	if err != nil {
		return Health{}, fmt.Errorf("health: %w", err)
	}

	return health, nil
}
