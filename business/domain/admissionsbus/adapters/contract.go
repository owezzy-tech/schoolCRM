package adapters

import "testing"

// ContractCase documents a fixture-based contract scenario for an adapter shell.
// Tests keep fixtures in code when the payload is already typed, and in testdata/*.json
// once a concrete HTTP adapter is introduced.
type ContractCase[Req any, Resp any] struct {
	Name    string
	Request Req
	Respond func(Req) (Resp, error)
	Assert  func(*testing.T, Resp, error)
}

// RunContractCases runs typed adapter fixture cases without live vendor calls.
func RunContractCases[Req any, Resp any](t *testing.T, cases []ContractCase[Req, Resp]) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()

			got, err := tc.Respond(tc.Request)
			tc.Assert(t, got, err)
		})
	}
}
