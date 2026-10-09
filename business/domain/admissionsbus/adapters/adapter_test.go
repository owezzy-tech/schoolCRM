package adapters

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExternalAdapterRetriesTransientFailures(t *testing.T) {
	t.Parallel()

	transientErr := errors.New("transient")
	calls := 0
	adapter := NewExternalAdapter(AdapterKNEC, Config{
		Enabled: true,
		RetryPolicy: RetryPolicy{
			MaxAttempts: 3,
			InitialWait: time.Nanosecond,
			MaxWait:     time.Nanosecond,
		},
	}, nil)

	err := adapter.Execute(context.Background(), "verify", func(context.Context) error {
		calls++
		if calls < 3 {
			return transientErr
		}

		return nil
	}, func(err error) bool {
		return errors.Is(err, transientErr)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}

func TestExternalAdapterOpensCircuitAfterFailures(t *testing.T) {
	t.Parallel()

	adapter := NewExternalAdapter(AdapterMPesaDaraja, Config{
		Enabled: true,
		RetryPolicy: RetryPolicy{
			MaxAttempts: 1,
		},
		CircuitBreaker: CircuitBreakerConfig{
			FailureThreshold: 1,
			SuccessThreshold: 1,
			OpenDuration:     time.Minute,
		},
	}, nil)

	err := adapter.Execute(context.Background(), "stk_push", func(context.Context) error {
		return errors.New("provider down")
	}, func(error) bool { return false })
	if err == nil {
		t.Fatal("expected first failure")
	}

	err = adapter.Execute(context.Background(), "stk_push", func(context.Context) error {
		return nil
	}, func(error) bool { return false })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("err = %v, want %v", err, ErrCircuitOpen)
	}
}

func TestExternalAdapterEmitsObserverEvents(t *testing.T) {
	t.Parallel()

	var events []OperationEvent
	adapter := NewExternalAdapter(AdapterCelcomAfrica, Config{Enabled: true}, ObserverFunc(func(_ context.Context, event OperationEvent) {
		events = append(events, event)
	}))

	err := adapter.Execute(context.Background(), "send_sms", func(context.Context) error {
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("events length = %d, want 2", len(events))
	}
	if events[0].Status != OperationStarted || events[1].Status != OperationSuccess {
		t.Fatalf("event statuses = %s, %s; want %s, %s", events[0].Status, events[1].Status, OperationStarted, OperationSuccess)
	}
}

func TestExternalAdapterStopsRetriesWhenCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	transientErr := errors.New("transient")
	calls := 0
	adapter := NewExternalAdapter(AdapterKNEC, Config{
		Enabled: true,
		RetryPolicy: RetryPolicy{
			MaxAttempts: 3,
			InitialWait: time.Minute,
			MaxWait:     time.Minute,
		},
	}, nil)

	err := adapter.Execute(ctx, "verify", func(context.Context) error {
		calls++
		if calls == 1 {
			cancel()
			return transientErr
		}
		return nil
	}, func(err error) bool {
		return errors.Is(err, transientErr)
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}
