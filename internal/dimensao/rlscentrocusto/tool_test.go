package rlscentrocusto

import (
	"context"
	"errors"
	"strings"
	"testing"

	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/database/oracle"
)

func TestSQLUsesCallerEmail(t *testing.T) {
	if !strings.Contains(searchSQL, `"E-mail" = :email`) || strings.Contains(searchSQL, ":email IS NULL") {
		t.Fatal(searchSQL)
	}
	db := &fakeDB{}
	_, err := Search(context.Background(), bi.ToolEnv{DB: db, Emails: stubEmail{}, MaxRows: 1000}, "user-1", entrada{})
	if !errors.Is(err, bi.ErrEmail) || db.called {
		t.Fatalf("err=%v called=%v", err, db.called)
	}
}

type stubEmail struct{}

func (stubEmail) Email(context.Context, string) (string, error) { return "", nil }

type fakeDB struct{ called bool }

func (f *fakeDB) Query(context.Context, string, string, ...any) (oracle.Result, error) {
	f.called = true
	return oracle.Result{}, nil
}
