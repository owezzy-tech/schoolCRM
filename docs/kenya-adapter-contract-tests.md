# Kenya Adapter Contract-Test Convention

Kenya external adapters are outbound ports under `business/domain/admissionsbus/adapters`. The current slice defines typed ports, shared adapter behavior, and fixture tests only. CI must not call KUCCPS, KNEC, IPRS, M-Pesa Daraja, Celcom Africa, or WhatsApp Cloud.

## Package layout

Each adapter package owns its narrow port, typed request/response models, and representative fixtures:

```text
business/domain/admissionsbus/adapters/
├── adapter.go              # shared config, retry, circuit breaker, observability hooks
├── contract.go             # generic fixture contract-test runner
├── kuccps/                 # Importer port
├── knec/                   # Verifier port
├── iprs/                   # Lookup port
├── mpesa/                  # DarajaGateway port
├── celcomafrica/           # SmsGateway port
└── whatsapp/               # CloudGateway port
```

## Rules

- Keep port methods `context.Context` first.
- Keep interfaces narrow and provider-specific.
- Keep domain-facing request/response models typed; add JSON wire DTOs only when a concrete HTTP adapter is introduced.
- Use `adapters.ExternalAdapter` for shared timeout, retry, circuit-breaker, and observability behavior.
- Use `adapters.Config.Normalize()` for CI-safe defaults. Adapters are disabled unless explicitly enabled by concrete wiring.
- Fixture tests must use static fixture functions or local `testdata/*.json` files.
- CI tests must never require vendor credentials and must never call live vendor endpoints.

## Fixture-test pattern

Use `adapters.RunContractCases` for typed port fixtures:

```go
func TestGatewayContractFixtures(t *testing.T) {
    t.Parallel()

    adapters.RunContractCases(t, []adapters.ContractCase[Request, Response]{
        {
            Name:    "representative success fixture",
            Request: FixtureRequest(),
            Respond: func(Request) (Response, error) {
                return FixtureResponse(), nil
            },
            Assert: func(t *testing.T, got Response, err error) {
                t.Helper()
                if err != nil {
                    t.Fatalf("unexpected error: %v", err)
                }
                if got.ExternalRef == "" {
                    t.Fatal("ExternalRef is empty")
                }
            },
        },
    })
}
```

When concrete HTTP clients are added, keep the same contract cases but replace the fixture responder with `httptest.Server` handlers backed by `testdata/*.json` request and response payloads.

## Required adapter fixture coverage

Every adapter shell must keep at least one representative success fixture:

- `kuccps`: placement-cycle import.
- `knec`: KCSE result verification.
- `iprs`: national ID verification.
- `mpesa`: STK Push initiation.
- `celcomafrica`: SMS send request.
- `whatsapp`: template send request.

Future concrete adapters should add non-retryable client-error and retryable transient-error fixtures without changing these port contracts.
