// Package sqlite guarda identidades, tokens e auditoria fora do Oracle.
package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// ErrNotFound indica que o registro não existe.
var ErrNotFound = errors.New("não encontrado")

// Store é o IdentityStore em SQLite.
type Store struct {
	db *sql.DB
}

// User é uma identidade local.
type User struct {
	ID                 string
	Username           string
	DisplayName        string
	Email              string
	PasswordHash       string
	Active             bool
	MustChangePassword bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// AccessToken é um bearer opaco já persistido.
type AccessToken struct {
	ID                 string
	UserID             string
	ClientID           string
	Scope              string
	Audience           string
	MustChangePassword bool
	ExpiresAt          time.Time
	RevokedAt          *time.Time
}

// RefreshToken é um refresh opaco já persistido.
type RefreshToken struct {
	ID            string
	UserID        string
	ClientID      string
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	RotatedFromID string
}

// AuthCode é um authorization code de uso único.
type AuthCode struct {
	ClientID      string
	UserID        string
	RedirectURI   string
	CodeChallenge string
	Scope         string
	Resource      string
	ExpiresAt     time.Time
	UsedAt        *time.Time
}

// OAuthClient é um cliente registrado por DCR.
type OAuthClient struct {
	ClientID        string
	ClientName      string
	RedirectURIs    []string
	ApplicationType string
}

// AuditEvent é um registro sem segredo.
type AuditEvent struct {
	ActorID       string
	Action        string
	Target        string
	Outcome       string
	CorrelationID string
	Metadata      map[string]string
}

// Open abre o arquivo, cria o diretório e aplica as migrações.
func Open(path string) (*Store, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("caminho do sqlite: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o700); err != nil {
		return nil, fmt.Errorf("diretório de dados: %w", err)
	}
	dsn := "file:" + filepath.ToSlash(abs) + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, err
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

// Close fecha o banco.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`); err != nil {
		return err
	}
	steps := []struct {
		version int
		file    string
	}{
		{1, "migrations/001_identity.sql"},
		{2, "migrations/002_user_email.sql"},
	}
	for _, step := range steps {
		var n int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, step.version).Scan(&n); err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		body, err := migrations.ReadFile(step.file)
		if err != nil {
			return err
		}
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		if err := execScript(tx, string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migração %d: %w", step.version, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, step.version, fmtTime(time.Now())); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func execScript(tx *sql.Tx, script string) error {
	for _, stmt := range strings.Split(script, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func fmtTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
}
