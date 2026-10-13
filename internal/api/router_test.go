package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/config"
	"oracle-mcp/internal/identity/sqlite"
)

func TestHealth(t *testing.T) {
	ts, _, _ := testServer(t)
	cases := map[string]string{
		"/api/v1/health/live":  "live",
		"/api/v1/health/ready": "ready",
	}
	for path, want := range cases {
		res, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("%s status = %d", path, res.StatusCode)
		}
		var payload map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("%s body: %v", path, err)
		}
		if payload["status"] != want {
			t.Fatalf("%s status = %q", path, payload["status"])
		}
	}
}

func TestReadyDoesNotLeakOracleError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := config.Config{
		Server: config.Server{Address: "127.0.0.1:8080", PublicBaseURL: "http://mcp.test", ReadTimeout: time.Second, WriteTimeout: time.Second},
		Auth:   config.Auth{SQLitePath: filepath.Join(t.TempDir(), "identity.db"), AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour},
		BI:     config.BI{MaxRows: 10, MaxResponseBytes: 1000},
	}
	store, err := sqlite.Open(cfg.Auth.SQLitePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	svc, err := auth.New(store, cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(NewHandler(logger, cfg, svc, errPinger{}, nil, nil, ""))
	t.Cleanup(ts.Close)
	res, err := http.Get(ts.URL + "/api/v1/health/ready")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if strings.Contains(string(body), "secret") || strings.Contains(string(body), "oracle.internal") {
		t.Fatalf("body vazou detalhe: %s", body)
	}
}

type errPinger struct{}

func (errPinger) Ping(context.Context) error {
	return errors.New("falha em oracle.internal senha=secret")
}

func testServer(t *testing.T) (*httptest.Server, *auth.Service, config.Config) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := config.Config{
		Server: config.Server{
			Address: "127.0.0.1:8080", PublicBaseURL: "http://mcp.test",
			ReadTimeout: time.Second, WriteTimeout: time.Second,
		},
		Auth: config.Auth{
			SQLitePath:     filepath.Join(t.TempDir(), "identity.db"),
			AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour,
		},
	}
	store, err := sqlite.Open(cfg.Auth.SQLitePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	svc, err := auth.New(store, cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(NewHandler(logger, cfg, svc, nil, nil, bi.StoreAuditor(store), ""))
	t.Cleanup(ts.Close)
	return ts, svc, cfg
}
