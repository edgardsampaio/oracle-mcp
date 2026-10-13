package bi

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLimiteClampsAndRejects(t *testing.T) {
	n, err := Limite(nil, 1000)
	if err != nil || n != DefaultLimite {
		t.Fatalf("padrão = %d %v", n, err)
	}
	asked := 5000
	n, err = Limite(&asked, 1000)
	if err != nil || n != 1000 {
		t.Fatalf("teto = %d %v", n, err)
	}
	zero := 0
	if _, err := Limite(&zero, 1000); !errors.Is(err, ErrParam) {
		t.Fatal(err)
	}
}

func TestNullBinds(t *testing.T) {
	if NullString(nil) != nil || NullInt(nil) != nil {
		t.Fatal("ausente deveria ser NULL")
	}
	blank := "  "
	if NullString(&blank) != nil {
		t.Fatal("texto vazio deveria ser NULL")
	}
	name := "adm"
	if NullString(&name) != "adm" {
		t.Fatal(NullString(&name))
	}
	id := 0
	if NullInt(&id) != 0 {
		t.Fatal("zero informado não é ausência")
	}
	when, err := NullDate(nil)
	if err != nil || when != nil {
		t.Fatalf("data ausente = %#v %v", when, err)
	}
	bad := "13/10/2026"
	if _, err := NullDate(&bad); !errors.Is(err, ErrParam) {
		t.Fatal(err)
	}
	ok := "2026-10-13"
	got, err := NullDate(&ok)
	if err != nil {
		t.Fatal(err)
	}
	day, okTime := got.(time.Time)
	if !okTime || day.Format("2006-01-02") != ok {
		t.Fatalf("data = %#v", got)
	}
}

func TestRequireEmailDoesNotInvent(t *testing.T) {
	if _, err := RequireEmail(context.Background(), stubEmail{email: ""}, "user-1"); !errors.Is(err, ErrEmail) {
		t.Fatal(err)
	}
	if _, err := RequireEmail(context.Background(), nil, "user-1"); !errors.Is(err, ErrEmail) {
		t.Fatal(err)
	}
	if _, err := RequireEmail(context.Background(), stubEmail{email: "a@b.co"}, ""); !errors.Is(err, ErrIdentity) {
		t.Fatal(err)
	}
	got, err := RequireEmail(context.Background(), stubEmail{email: " a@b.co "}, "user-1")
	if err != nil || got != "a@b.co" {
		t.Fatalf("email = %q %v", got, err)
	}
}

type stubEmail struct {
	email string
}

func (s stubEmail) Email(context.Context, string) (string, error) {
	return s.email, nil
}
