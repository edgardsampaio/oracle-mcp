package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"oracle-mcp/internal/auth"
)

func TestAdminCreateRejectsExtraArgs(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	err := runAdmin(logger, []string{"create", "--password", "segredo"})
	if err == nil {
		t.Fatal("esperava recusa de argumento")
	}
}

func TestBootstrapSkipsWhenAdminExists(t *testing.T) {
	var out bytes.Buffer
	_, err := ensureInitialAdmin(context.Background(), 1, &out, func() (string, string, string, error) {
		t.Fatal("não deveria pedir senha")
		return "", "", "", nil
	}, func(context.Context, string, string, string) error {
		t.Fatal("não deveria criar")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBootstrapRetriesMismatchThenCreates(t *testing.T) {
	var out bytes.Buffer
	calls := 0
	var gotUser, gotPass string
	user, err := bootstrapAdmin(context.Background(), &out, func() (string, string, string, error) {
		calls++
		if calls == 1 {
			return "", "", "", errPasswordMismatch
		}
		return "ada.admin", "ada.admin@example.com", "senha-admin-ok", nil
	}, func(_ context.Context, username, email, password string) error {
		gotUser, gotPass = username, password
		if email != "ada.admin@example.com" {
			t.Fatalf("e-mail = %s", email)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if user != "ada.admin" || gotUser != "ada.admin" || gotPass != "senha-admin-ok" {
		t.Fatalf("usuário %q senha recebida %q", gotUser, gotPass)
	}
	if !bytes.Contains(out.Bytes(), []byte("as senhas não conferem")) {
		t.Fatalf("saída: %s", out.String())
	}
}

func TestBootstrapRetriesInvalidPassword(t *testing.T) {
	var out bytes.Buffer
	calls := 0
	_, err := bootstrapAdmin(context.Background(), &out, func() (string, string, string, error) {
		calls++
		return "ada.admin", "ada.admin@example.com", "senha-admin-ok", nil
	}, func(context.Context, string, string, string) error {
		if calls == 1 {
			return errors.Join(auth.ErrValidation, errors.New("a senha deve ter entre 6 e 128 caracteres"))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("pedidos: %d", calls)
	}
}

func TestBootstrapStopsWhenTerminalMissing(t *testing.T) {
	var out bytes.Buffer
	_, err := bootstrapAdmin(context.Background(), &out, func() (string, string, string, error) {
		return "", "", "", errors.New("execute em um terminal; a senha não é aceita por argumento")
	}, func(context.Context, string, string, string) error {
		t.Fatal("não deveria criar")
		return nil
	})
	if err == nil {
		t.Fatal("esperava falha de terminal")
	}
}
