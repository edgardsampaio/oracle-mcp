package confissaodivida

import (
	"context"
	"errors"
	"strings"
	"testing"

	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/database/oracle"
)

func TestSQLRestringePorEmail(t *testing.T) {
	for _, statement := range []string{sqlFato, sqlComDimensoes} {
		if !strings.Contains(statement, "DIM_RLS_CENTRO_CUSTO") || !strings.Contains(statement, ":email") {
			t.Fatalf("sem rls: %s", statement[:80])
		}
		if strings.Contains(strings.ToUpper(statement), "SELECT *") {
			t.Fatal("select *")
		}
	}
}

func TestSemEmailNaoConsulta(t *testing.T) {
	db := &fakeDB{}
	env := bi.ToolEnv{DB: db, Emails: stubEmail{}, MaxRows: 1000}
	if _, err := Search(context.Background(), env, "user-1", entrada{}); !errors.Is(err, bi.ErrEmail) || db.called {
		t.Fatalf("fato err=%v called=%v", err, db.called)
	}
	db.called = false
	if _, err := SearchComDimensoes(context.Background(), env, "user-1", entradaCom{}); !errors.Is(err, bi.ErrEmail) || db.called {
		t.Fatalf("com dimensões err=%v called=%v", err, db.called)
	}
}

type stubEmail struct{}

func (stubEmail) Email(context.Context, string) (string, error) { return "", nil }

type fakeDB struct {
	called bool
}

func (f *fakeDB) Query(context.Context, string, string, ...any) (oracle.Result, error) {
	f.called = true
	return oracle.Result{}, nil
}
