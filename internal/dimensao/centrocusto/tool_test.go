package centrocusto

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/database/oracle"
)

func TestSearchUsesHierarchySQL(t *testing.T) {
	db := &fakeDB{rows: oracle.Result{Rows: []map[string]any{{
		"ID Centro de Custo":   "10",
		"Nome Centro de Custo": "Administrativo",
	}}}}
	nome := "adm"
	env := bi.ToolEnv{DB: db, MaxRows: 1000}
	rows, err := Search(context.Background(), env, "user-1", entrada{NomeCentroCusto: &nome})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["Nome Centro de Custo"] != "Administrativo" {
		t.Fatalf("rows = %#v", rows)
	}
	if !strings.Contains(db.statement, "DIM_CENTRO_CUSTO") || !strings.Contains(db.statement, ":nome_centro_custo") || !strings.Contains(db.statement, "Nome Gestor") {
		t.Fatalf("sql = %s", db.statement)
	}
	if strings.Contains(db.statement, ":1") || strings.Contains(strings.ToUpper(db.statement), "SELECT *") {
		t.Fatalf("sql = %s", db.statement)
	}
	named, ok := db.argsNamed("nome_centro_custo")
	if !ok || named != "adm" {
		t.Fatalf("busca = %#v", db.args)
	}
	limite, ok := db.argsNamed("limite")
	if !ok || limite != 500 {
		t.Fatalf("limite = %#v", limite)
	}
}

func TestSearchAuditsFailure(t *testing.T) {
	audit := &memAudit{}
	_, err := Search(context.Background(), bi.ToolEnv{Audit: audit, MaxRows: 1000}, "user-1", entrada{})
	if !errors.Is(err, bi.ErrUnavailable) {
		t.Fatal(err)
	}
	if len(audit.events) != 1 || audit.events[0].tool != ToolName || audit.events[0].actor != "user-1" || audit.events[0].outcome != "failure" {
		t.Fatalf("audit = %#v", audit.events)
	}
}

type memAudit struct {
	events []auditEvent
}

type auditEvent struct {
	actor, tool, outcome string
}

func (m *memAudit) Record(_ context.Context, actor, tool, outcome string, duration time.Duration) error {
	m.events = append(m.events, auditEvent{actor: actor, tool: tool, outcome: outcome})
	_ = duration
	return nil
}

type fakeDB struct {
	rows      oracle.Result
	statement string
	args      []any
}

func (f *fakeDB) Query(_ context.Context, _, statement string, args ...any) (oracle.Result, error) {
	f.statement = statement
	f.args = args
	return f.rows, nil
}

func (f *fakeDB) argsNamed(name string) (any, bool) {
	for _, arg := range f.args {
		named, ok := arg.(sql.NamedArg)
		if ok && named.Name == name {
			return named.Value, true
		}
	}
	return nil, false
}
