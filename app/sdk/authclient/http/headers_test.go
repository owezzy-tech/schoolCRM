package http

import (
	"bytes"
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owezzy/schoolCRM/foundation/logger"
)

func TestCredentialsAreForwardedButNotLogged(t *testing.T) {
	const bearer = "Bearer test-private-bearer"
	const apiKey = "test-private-api-key"
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Header.Get("Authorization") != bearer || r.Header.Get("X-API-Key") != apiKey {
			t.Error("credential forwarding changed")
		}
		w.WriteHeader(stdhttp.StatusNoContent)
	}))
	defer server.Close()
	var output bytes.Buffer
	log := logger.New(&output, logger.LevelInfo, "AUTH-TEST", func(context.Context) string { return "test" })
	client, err := New(log, server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.do(context.Background(), stdhttp.MethodGet, server.URL,
		map[string]string{"Authorization": bearer, "X-API-Key": apiKey}, nil, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), bearer) || strings.Contains(output.String(), apiKey) {
		t.Fatal("request credentials entered logs")
	}
	if !strings.Contains(output.String(), "completed") {
		t.Fatal("request diagnostics lost")
	}
}
