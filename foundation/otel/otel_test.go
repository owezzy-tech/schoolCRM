package otel

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace/noop"
)

func TestInitTracingEmptyHostUsesNoopProvider(t *testing.T) {
	provider, teardown, err := InitTracing(Config{
		ServiceName: "test",
		Host:        "",
	})
	if err != nil {
		t.Fatalf("InitTracing returned error: %v", err)
	}

	if _, ok := provider.(noop.TracerProvider); !ok {
		t.Fatalf("provider = %T, want noop.TracerProvider", provider)
	}

	teardown(context.Background())
}

func TestInitTracingWhitespaceHostUsesNoopProvider(t *testing.T) {
	provider, teardown, err := InitTracing(Config{
		ServiceName: "test",
		Host:        "   ",
	})
	if err != nil {
		t.Fatalf("InitTracing returned error: %v", err)
	}

	if _, ok := provider.(noop.TracerProvider); !ok {
		t.Fatalf("provider = %T, want noop.TracerProvider", provider)
	}

	teardown(context.Background())
}
