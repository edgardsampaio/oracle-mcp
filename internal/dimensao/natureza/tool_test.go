package natureza

import (
	"context"
	"strings"
	"testing"

	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/database/oracle"
)

func TestSearchKeepsHierarchy(t *testing.T) {
	if !strings.Contains(description, "natureza") {
		t.Fatal("descrição vazia")
	}
	db := &fakeDB{rows: oracle.Result{Rows: []map[string]any{{
		"ID Natureza":   "5",
		"Nome Natureza": "Energia",
	}}}}
	nome := "energia"
	rows, err := Search(context.Background(), bi.ToolEnv{DB: db, MaxRows: 1000}, "user-1", entrada{NomeNatureza: &nome})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["Nome Natureza"] != "Energia" {
		t.Fatalf("rows = %#v", rows)
	}
	if ToolName != "dim_natureza" {
		t.Fatal(ToolName)
	}
	if !strings.Contains(db.statement, "DIM_NATUREZA") || !strings.Contains(db.statement, ":nome_natureza") || !strings.Contains(db.statement, "Nome Natureza N1") {
		t.Fatalf("sql = %s", db.statement)
	}
	if strings.Contains(db.statement, ":1") || strings.Contains(strings.ToUpper(db.statement), "SELECT *") {
		t.Fatalf("sql = %s", db.statement)
	}
}

type fakeDB struct {
	rows      oracle.Result
	statement string
}

func (f *fakeDB) Query(_ context.Context, _, statement string, _ ...any) (oracle.Result, error) {
	f.statement = statement
	return f.rows, nil
}
