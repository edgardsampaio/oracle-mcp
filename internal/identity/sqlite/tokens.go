package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// InsertAccessToken grava só o hash.
func (s *Store) InsertAccessToken(ctx context.Context, id, userID, clientID, tokenHash, scope, audience string, mustChange bool, expires, created time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO access_tokens (id, user_id, client_id, token_hash, scope, audience, must_change_password, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, userID, clientID, tokenHash, scope, audience, boolInt(mustChange), fmtTime(expires), fmtTime(created))
	return err
}

// AccessByHash busca um access token pelo hash.
func (s *Store) AccessByHash(ctx context.Context, tokenHash string) (AccessToken, error) {
	var t AccessToken
	var must int
	var exp string
	var revoked sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, client_id, scope, audience, must_change_password, expires_at, revoked_at
		FROM access_tokens WHERE token_hash = ?`, tokenHash).
		Scan(&t.ID, &t.UserID, &t.ClientID, &t.Scope, &t.Audience, &must, &exp, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return AccessToken{}, ErrNotFound
	}
	if err != nil {
		return AccessToken{}, err
	}
	t.MustChangePassword = must == 1
	t.ExpiresAt, err = parseTime(exp)
	if err != nil {
		return AccessToken{}, err
	}
	if revoked.Valid {
		ts, err := parseTime(revoked.String)
		if err != nil {
			return AccessToken{}, err
		}
		t.RevokedAt = &ts
	}
	return t, nil
}

// InsertRefreshToken grava só o hash.
func (s *Store) InsertRefreshToken(ctx context.Context, id, userID, clientID, tokenHash string, expires, created time.Time, rotatedFrom string) error {
	var from any
	if rotatedFrom != "" {
		from = rotatedFrom
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (id, user_id, client_id, token_hash, expires_at, created_at, rotated_from_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, userID, clientID, tokenHash, fmtTime(expires), fmtTime(created), from)
	return err
}

// RefreshByHash busca um refresh token pelo hash.
func (s *Store) RefreshByHash(ctx context.Context, tokenHash string) (RefreshToken, error) {
	var t RefreshToken
	var exp string
	var revoked, rotated sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, client_id, expires_at, revoked_at, rotated_from_id
		FROM refresh_tokens WHERE token_hash = ?`, tokenHash).
		Scan(&t.ID, &t.UserID, &t.ClientID, &exp, &revoked, &rotated)
	if errors.Is(err, sql.ErrNoRows) {
		return RefreshToken{}, ErrNotFound
	}
	if err != nil {
		return RefreshToken{}, err
	}
	t.ExpiresAt, err = parseTime(exp)
	if err != nil {
		return RefreshToken{}, err
	}
	if revoked.Valid {
		ts, err := parseTime(revoked.String)
		if err != nil {
			return RefreshToken{}, err
		}
		t.RevokedAt = &ts
	}
	if rotated.Valid {
		t.RotatedFromID = rotated.String
	}
	return t, nil
}

// RefreshHasChild informa se o token já foi rotacionado.
func (s *Store) RefreshHasChild(ctx context.Context, id string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM refresh_tokens WHERE rotated_from_id = ?`, id).Scan(&n)
	return n > 0, err
}

// RevokeRefresh marca o refresh como revogado.
func (s *Store) RevokeRefresh(ctx context.Context, id string, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`, fmtTime(now), id)
	return err
}

// RevokeUserTokens revoga access e refresh ainda válidos do usuário.
func (s *Store) RevokeUserTokens(ctx context.Context, userID string, now time.Time) error {
	ts := fmtTime(now)
	if _, err := s.db.ExecContext(ctx, `UPDATE access_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`, ts, userID); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`, ts, userID)
	return err
}

// InsertAuthCode grava o hash do authorization code.
func (s *Store) InsertAuthCode(ctx context.Context, codeHash string, c AuthCode, created time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO oauth_authorization_codes
			(code_hash, client_id, user_id, redirect_uri, code_challenge, scope, resource, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		codeHash, c.ClientID, c.UserID, c.RedirectURI, c.CodeChallenge, c.Scope, c.Resource, fmtTime(c.ExpiresAt), fmtTime(created))
	return err
}

// ConsumeAuthCode marca o code como usado e o devolve.
func (s *Store) ConsumeAuthCode(ctx context.Context, codeHash string, now time.Time) (AuthCode, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return AuthCode{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var c AuthCode
	var exp string
	var used sql.NullString
	err = tx.QueryRowContext(ctx, `
		SELECT client_id, user_id, redirect_uri, code_challenge, scope, resource, expires_at, used_at
		FROM oauth_authorization_codes WHERE code_hash = ?`, codeHash).
		Scan(&c.ClientID, &c.UserID, &c.RedirectURI, &c.CodeChallenge, &c.Scope, &c.Resource, &exp, &used)
	if errors.Is(err, sql.ErrNoRows) {
		return AuthCode{}, ErrNotFound
	}
	if err != nil {
		return AuthCode{}, err
	}
	c.ExpiresAt, err = parseTime(exp)
	if err != nil {
		return AuthCode{}, err
	}
	if used.Valid || !c.ExpiresAt.After(now) {
		return AuthCode{}, ErrNotFound
	}
	if _, err := tx.ExecContext(ctx, `UPDATE oauth_authorization_codes SET used_at = ? WHERE code_hash = ?`, fmtTime(now), codeHash); err != nil {
		return AuthCode{}, err
	}
	if err := tx.Commit(); err != nil {
		return AuthCode{}, err
	}
	return c, nil
}

// InsertOAuthClient grava um cliente DCR.
func (s *Store) InsertOAuthClient(ctx context.Context, c OAuthClient, created time.Time) error {
	raw, err := json.Marshal(c.RedirectURIs)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO oauth_clients (client_id, client_name, redirect_uris, application_type, created_at)
		VALUES (?, ?, ?, ?, ?)`, c.ClientID, c.ClientName, string(raw), c.ApplicationType, fmtTime(created))
	return err
}

// OAuthClientByID busca um cliente DCR.
func (s *Store) OAuthClientByID(ctx context.Context, clientID string) (OAuthClient, error) {
	var c OAuthClient
	var raw string
	err := s.db.QueryRowContext(ctx, `
		SELECT client_id, client_name, redirect_uris, application_type FROM oauth_clients WHERE client_id = ?`, clientID).
		Scan(&c.ClientID, &c.ClientName, &raw, &c.ApplicationType)
	if errors.Is(err, sql.ErrNoRows) {
		return OAuthClient{}, ErrNotFound
	}
	if err != nil {
		return OAuthClient{}, err
	}
	if err := json.Unmarshal([]byte(raw), &c.RedirectURIs); err != nil {
		return OAuthClient{}, err
	}
	return c, nil
}

// InsertAudit grava um evento sem segredo.
func (s *Store) InsertAudit(ctx context.Context, e AuditEvent, now time.Time) error {
	if e.Metadata == nil {
		e.Metadata = map[string]string{}
	}
	raw, err := json.Marshal(e.Metadata)
	if err != nil {
		return err
	}
	var actor any
	if e.ActorID != "" {
		actor = e.ActorID
	}
	if e.CorrelationID == "" {
		e.CorrelationID = uuid.NewString()
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO audit_log (id, timestamp, actor_id, action, target, outcome, correlation_id, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.NewString(), fmtTime(now), actor, e.Action, e.Target, e.Outcome, e.CorrelationID, string(raw))
	return err
}

// RecentFailedLogins conta falhas recentes do username.
func (s *Store) RecentFailedLogins(ctx context.Context, username string, since time.Time) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM login_attempts
		WHERE username = ? COLLATE NOCASE AND success = 0 AND at >= ?`, username, fmtTime(since)).Scan(&n)
	return n, err
}

// InsertLoginAttempt registra o resultado sem a senha.
func (s *Store) InsertLoginAttempt(ctx context.Context, username string, success bool, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO login_attempts (username, at, success) VALUES (?, ?, ?)`, username, fmtTime(now), boolInt(success))
	return err
}
