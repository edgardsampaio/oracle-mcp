package bi

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"oracle-mcp/internal/identity/sqlite"
)

func TestStoreAuditorKeepsMetadataSmall(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.db")
	store, err := sqlite.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := StoreAuditor(store).Record(context.Background(), "user-1", "dim_centro_custo", "failure", time.Millisecond); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", "file:"+filepath.ToSlash(path)+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	var metadata string
	if err := db.QueryRow(`SELECT metadata FROM audit_log`).Scan(&metadata); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(metadata, "duration_ms") || strings.Contains(metadata, "SELECT") || strings.Contains(metadata, "Administrativo") {
		t.Fatalf("metadata = %s", metadata)
	}
}
