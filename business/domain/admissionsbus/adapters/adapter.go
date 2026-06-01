// Package adapters defines outbound ports and shared behavior for Kenya admissions integrations.
package adapters

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AdapterName identifies a Kenya external integration in logs, receipts, and verification metadata.
type AdapterName string

const (
	AdapterKUCCPS         AdapterName = "kuccps"
	AdapterKNEC           AdapterName = "knec"
	AdapterIPRS           AdapterName = "iprs"
	AdapterMPesaDaraja    AdapterName = "mpesa_daraja"
	AdapterAfricasTalking AdapterName = "africas_talking"
	AdapterWhatsAppCloud  AdapterName = "whatsapp_cloud"
)

// ErrCircuitOpen is returned when an adapter circuit breaker is rejecting calls.
var ErrCircuitOpen = errors.New("adapter circuit open")

// Config captures vendor-neutral adapter behavior. Vendor credentials belong in concrete adapters.
type Config struct {
	Enabled        bool
	Timeout        time.Duration
	RetryPolicy    RetryPolicy
	CircuitBreaker CircuitBreakerConfig
}

// Normalize fills safe defaults. The default is CI-safe: adapters are disabled unless enabled explicitly.
func (cfg Config) Normalize() Config {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	cfg.RetryPolicy = cfg.RetryPolicy.Normalize()
	cfg.CircuitBreaker = cfg.CircuitBreaker.Normalize()

	return cfg
}

// RetryPolicy controls retry attempts for transient adapter failures.
type RetryPolicy struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
}

// Normalize fills bounded retry defaults.
func (policy RetryPolicy) Normalize() RetryPolicy {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 3
	}
	if policy.InitialWait <= 0 {
		policy.InitialWait = 100 * time.Millisecond
	}
	if policy.MaxWait <= 0 {
		policy.MaxWait = time.Second
	}
	if policy.MaxWait < policy.InitialWait {
		policy.MaxWait = policy.InitialWait
	}

	return policy
}

// CircuitBreakerConfig controls the simple adapter circuit breaker.
type CircuitBreakerConfig struct {
	FailureThreshold int
	SuccessThreshold int
	OpenDuration     time.Duration
}

// Normalize fills conservative circuit breaker defaults.
func (cfg CircuitBreakerConfig) Normalize() CircuitBreakerConfig {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.OpenDuration <= 0 {
		cfg.OpenDuration = 30 * time.Second
	}

	return cfg
}

// OperationStatus describes one adapter call for observability hooks.
type OperationStatus string

const (
	OperationStarted OperationStatus = "started"
	OperationRetried OperationStatus = "retried"
	OperationSuccess OperationStatus = "success"
	OperationFailure OperationStatus = "failure"
)

// OperationEvent is emitted by ExternalAdapter around retries and circuit-breaker transitions.
type OperationEvent struct {
	Adapter  AdapterName
	Name     string
	Status   OperationStatus
	Attempt  int
	Duration time.Duration
	Err      error
	Time     time.Time
}

// Observer receives adapter operation events. Implementations should return quickly and never mutate events.
type Observer interface {
	ObserveAdapterOperation(ctx context.Context, event OperationEvent)
}

// ObserverFunc adapts a function to Observer.
type ObserverFunc func(ctx context.Context, event OperationEvent)

// ObserveAdapterOperation implements Observer.
func (fn ObserverFunc) ObserveAdapterOperation(ctx context.Context, event OperationEvent) {
	fn(ctx, event)
}

type breakerState int

const (
	breakerClosed breakerState = iota
	breakerOpen
	breakerHalfOpen
)

// ExternalAdapter provides common retry, circuit-breaker, timeout, and observability behavior.
type ExternalAdapter struct {
	name     AdapterName
	cfg      Config
	observer Observer

	mu        sync.Mutex
	state     breakerState
	failures  int
	successes int
	openedAt  time.Time
}

// NewExternalAdapter constructs shared behavior for one concrete adapter instance.
func NewExternalAdapter(name AdapterName, cfg Config, observer Observer) *ExternalAdapter {
	return &ExternalAdapter{
		name:     name,
		cfg:      cfg.Normalize(),
		observer: observer,
	}
}

// Name returns the configured adapter name.
func (adapter *ExternalAdapter) Name() AdapterName {
	return adapter.name
}

// Config returns normalized adapter configuration.
func (adapter *ExternalAdapter) Config() Config {
	return adapter.cfg
}

// Execute runs an adapter operation through timeout, retry, circuit-breaker, and observer hooks.
func (adapter *ExternalAdapter) Execute(ctx context.Context, operation string, fn func(context.Context) error, isTransient func(error) bool) error {
	if adapter == nil {
		return errors.New("external adapter is nil")
	}

	if !adapter.cfg.Enabled {
		return errors.New("external adapter disabled")
	}

	if err := adapter.beforeExecute(); err != nil {
		adapter.observe(ctx, OperationEvent{
			Adapter: adapter.name,
			Name:    operation,
			Status:  OperationFailure,
			Err:     err,
			Time:    time.Now(),
		})

		return err
	}

	policy := adapter.cfg.RetryPolicy.Normalize()
	start := time.Now()
	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		opCtx, cancel := context.WithTimeout(ctx, adapter.cfg.Timeout)
		adapter.observe(opCtx, OperationEvent{
			Adapter: adapter.name,
			Name:    operation,
			Status:  OperationStarted,
			Attempt: attempt,
			Time:    time.Now(),
		})

		err := fn(opCtx)
		cancel()
		if err == nil {
			adapter.recordSuccess()
			adapter.observe(ctx, OperationEvent{
				Adapter:  adapter.name,
				Name:     operation,
				Status:   OperationSuccess,
				Attempt:  attempt,
				Duration: time.Since(start),
				Time:     time.Now(),
			})

			return nil
		}

		lastErr = err
		if isTransient == nil || !isTransient(err) || attempt == policy.MaxAttempts {
			break
		}

		wait := retryWait(policy, attempt)
		adapter.observe(ctx, OperationEvent{
			Adapter:  adapter.name,
			Name:     operation,
			Status:   OperationRetried,
			Attempt:  attempt,
			Duration: time.Since(start),
			Err:      err,
			Time:     time.Now(),
		})

		select {
		case <-ctx.Done():
			lastErr = ctx.Err()
			break
		case <-time.After(wait):
		}
	}

	adapter.recordFailure()
	wrapped := fmt.Errorf("%s %s: %w", adapter.name, operation, lastErr)
	adapter.observe(ctx, OperationEvent{
		Adapter:  adapter.name,
		Name:     operation,
		Status:   OperationFailure,
		Attempt:  policy.MaxAttempts,
		Duration: time.Since(start),
		Err:      wrapped,
		Time:     time.Now(),
	})

	return wrapped
}

func (adapter *ExternalAdapter) beforeExecute() error {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	if adapter.state != breakerOpen {
		return nil
	}

	if time.Since(adapter.openedAt) < adapter.cfg.CircuitBreaker.OpenDuration {
		return ErrCircuitOpen
	}

	adapter.state = breakerHalfOpen
	adapter.successes = 0

	return nil
}

func (adapter *ExternalAdapter) recordSuccess() {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	adapter.failures = 0
	if adapter.state != breakerHalfOpen {
		return
	}

	adapter.successes++
	if adapter.successes >= adapter.cfg.CircuitBreaker.SuccessThreshold {
		adapter.state = breakerClosed
		adapter.successes = 0
	}
}

func (adapter *ExternalAdapter) recordFailure() {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	adapter.successes = 0
	adapter.failures++
	if adapter.failures >= adapter.cfg.CircuitBreaker.FailureThreshold || adapter.state == breakerHalfOpen {
		adapter.state = breakerOpen
		adapter.openedAt = time.Now()
	}
}

func (adapter *ExternalAdapter) observe(ctx context.Context, event OperationEvent) {
	if adapter.observer == nil {
		return
	}

	adapter.observer.ObserveAdapterOperation(ctx, event)
}

func retryWait(policy RetryPolicy, attempt int) time.Duration {
	wait := policy.InitialWait
	for range attempt - 1 {
		wait *= 2
		if wait >= policy.MaxWait {
			return policy.MaxWait
		}
	}

	return wait
}
