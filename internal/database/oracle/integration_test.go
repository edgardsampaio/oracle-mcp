package oracle

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"oracle-mcp/internal/config"
)

func TestIntegrationReadOnly(t *testing.T) {
	if os.Getenv("ORACLE_MCP_ORACLE_INTEGRATION") != "1" {
		t.Skip("defina ORACLE_MCP_ORACLE_INTEGRATION=1 para consultar o Oracle autorizado")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Oracle.Enabled() {
		t.Fatal("oracle.host ausente")
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	db, err := Open(cfg.Oracle, cfg.BI, logger)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	result, err := db.Query(ctx, "dual", "SELECT 1 AS n, NULL AS z, SYSDATE AS d FROM DUAL")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("linhas = %d", len(result.Rows))
	}
}
