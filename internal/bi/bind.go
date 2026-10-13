package bi

import (
	"context"
	"errors"
	"strings"
	"time"

	"oracle-mcp/internal/identity/sqlite"
)

const (
	// DefaultLimite é o teto quando a tool não recebe limite.
	DefaultLimite = 500
	fallbackMax   = 1000
)

// EmailSource lê o e-mail do usuário já autenticado.
type EmailSource interface {
	Email(ctx context.Context, userID string) (string, error)
}

// Limite devolve o FETCH FIRST. Omitido vale 500 e não passa de maxRows.
func Limite(v *int, maxRows int) (int, error) {
	if maxRows <= 0 {
		maxRows = fallbackMax
	}
	n := DefaultLimite
	if v != nil {
		n = *v
	}
	if n < 1 {
		return 0, ErrParam
	}
	if n > maxRows {
		n = maxRows
	}
	return n, nil
}

// NullString omite texto vazio. O Oracle recebe NULL e ignora o filtro.
func NullString(v *string) any {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return s
}

// NullInt omite inteiro ausente. Zero informado segue como zero.
func NullInt(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

// NullFloat omite decimal ausente. Zero informado segue como zero.
func NullFloat(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

// NullDate aceita AAAA-MM-DD. Ausente vira NULL.
func NullDate(v *string) (any, error) {
	if v == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, ErrParam
	}
	return t, nil
}

// RequireEmail lê o e-mail do cadastro. Não aceita valor vindo da tool.
func RequireEmail(ctx context.Context, src EmailSource, userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", ErrIdentity
	}
	if src == nil {
		return "", ErrEmail
	}
	email, err := src.Email(ctx, userID)
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return "", ErrIdentity
		}
		return "", ErrUnavailable
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrEmail
	}
	return email, nil
}
