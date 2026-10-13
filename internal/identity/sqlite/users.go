package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CountActiveAdmins conta administradores ativos.
func (s *Store) CountActiveAdmins(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE r.code = 'admin' AND u.active = 1`).Scan(&n)
	return n, err
}

// InsertUser grava o usuário e os papéis informados.
func (s *Store) InsertUser(ctx context.Context, u User, roleCodes []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, username, display_name, email, password_hash, active, must_change_password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.DisplayName, emailArg(u.Email), u.PasswordHash, boolInt(u.Active), boolInt(u.MustChangePassword), fmtTime(u.CreatedAt), fmtTime(u.UpdatedAt))
	if err != nil {
		return uniqueUserErr(err)
	}
	if err := insertRoles(ctx, tx, u.ID, roleCodes); err != nil {
		return err
	}
	return tx.Commit()
}

func insertRoles(ctx context.Context, tx *sql.Tx, userID string, roleCodes []string) error {
	for _, code := range roleCodes {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO user_roles (user_id, role_id)
			SELECT ?, id FROM roles WHERE code = ?`, userID, code)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("papel desconhecido: %s", code)
		}
	}
	return nil
}

// UserByUsername busca pelo login.
func (s *Store) UserByUsername(ctx context.Context, username string) (User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, password_hash, active, must_change_password, created_at, updated_at
		FROM users WHERE username = ? COLLATE NOCASE`, username)
	return scanUser(row)
}

// UserByID busca pelo identificador.
func (s *Store) UserByID(ctx context.Context, id string) (User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, password_hash, active, must_change_password, created_at, updated_at
		FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func scanUser(row *sql.Row) (User, error) {
	var u User
	var email sql.NullString
	var active, must int
	var created, updated string
	err := row.Scan(&u.ID, &u.Username, &u.DisplayName, &email, &u.PasswordHash, &active, &must, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}
	u.Email = email.String
	u.Active = active == 1
	u.MustChangePassword = must == 1
	u.CreatedAt, err = parseTime(created)
	if err != nil {
		return User{}, err
	}
	u.UpdatedAt, err = parseTime(updated)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

// ListUsers devolve uma página e o total.
func (s *Store) ListUsers(ctx context.Context, page, pageSize int) ([]User, int, error) {
	var total int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, username, display_name, email, password_hash, active, must_change_password, created_at, updated_at
		FROM users ORDER BY username LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		u, err := scanUserRows(rows)
		if err != nil {
			return nil, 0, err
		}
		u.PasswordHash = ""
		out = append(out, u)
	}
	return out, total, rows.Err()
}

func scanUserRows(rows *sql.Rows) (User, error) {
	var u User
	var email sql.NullString
	var active, must int
	var created, updated string
	if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &email, &u.PasswordHash, &active, &must, &created, &updated); err != nil {
		return User{}, err
	}
	var err error
	u.Email = email.String
	u.Active = active == 1
	u.MustChangePassword = must == 1
	u.CreatedAt, err = parseTime(created)
	if err != nil {
		return User{}, err
	}
	u.UpdatedAt, err = parseTime(updated)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

// UpdateUser altera nome, e-mail e situação.
func (s *Store) UpdateUser(ctx context.Context, id, displayName, email string, active bool, now time.Time) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE users SET display_name = ?, email = ?, active = ?, updated_at = ? WHERE id = ?`,
		displayName, emailArg(email), boolInt(active), fmtTime(now), id)
	if err != nil {
		return uniqueUserErr(err)
	}
	return rowsAffected(res)
}

func emailArg(email string) any {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}
	return email
}

func uniqueUserErr(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "email") {
		return fmt.Errorf("email já existe")
	}
	if strings.Contains(msg, "unique") {
		return fmt.Errorf("username já existe")
	}
	return err
}

// SetPassword troca o hash e a flag de troca obrigatória.
func (s *Store) SetPassword(ctx context.Context, id, hash string, mustChange bool, now time.Time) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE users SET password_hash = ?, must_change_password = ?, updated_at = ? WHERE id = ?`,
		hash, boolInt(mustChange), fmtTime(now), id)
	if err != nil {
		return err
	}
	return rowsAffected(res)
}

// SetRoles substitui os papéis do usuário.
func (s *Store) SetRoles(ctx context.Context, userID string, roleCodes []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = ?`, userID); err != nil {
		return err
	}
	if err := insertRoles(ctx, tx, userID, roleCodes); err != nil {
		return err
	}
	return tx.Commit()
}

// RoleCodes lista os papéis do usuário.
func (s *Store) RoleCodes(ctx context.Context, userID string) ([]string, error) {
	return queryCodes(ctx, s.db, `
		SELECT r.code FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
		ORDER BY r.code`, userID)
}

// PermissionCodes lista as permissões efetivas.
func (s *Store) PermissionCodes(ctx context.Context, userID string) ([]string, error) {
	return queryCodes(ctx, s.db, `
		SELECT DISTINCT p.code FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
		ORDER BY p.code`, userID)
}

// RoleCatalog devolve o catálogo de papéis.
func (s *Store) RoleCatalog(ctx context.Context) ([]struct{ Code, Description string }, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT code, description FROM roles ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct{ Code, Description string }
	for rows.Next() {
		var item struct{ Code, Description string }
		if err := rows.Scan(&item.Code, &item.Description); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func queryCodes(ctx context.Context, db *sql.DB, q, arg string) ([]string, error) {
	rows, err := db.QueryContext(ctx, q, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func rowsAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
