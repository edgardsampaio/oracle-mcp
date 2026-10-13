package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"oracle-mcp/internal/identity/sqlite"
)

// CreateUser cadastra um usuário com senha temporária.
func (s *Service) CreateUser(ctx context.Context, actor Principal, username, name, email, password, role string, active bool) (UserView, error) {
	if !actor.Allows(PermUsersManage) {
		return UserView{}, ErrForbidden
	}
	if err := validateUsername(username); err != nil {
		return UserView{}, err
	}
	if err := validateEmail(email); err != nil {
		return UserView{}, err
	}
	if err := validatePassword(password); err != nil {
		return UserView{}, err
	}
	if stringsTrim(name) == "" || !knownRole(role) {
		return UserView{}, ErrValidation
	}
	hash, err := HashPassword(password)
	if err != nil {
		return UserView{}, err
	}
	now := s.now()
	u := sqlite.User{
		ID: uuid.NewString(), Username: username, DisplayName: stringsTrim(name), Email: strings.TrimSpace(email),
		PasswordHash: hash, Active: active, MustChangePassword: true, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.store.InsertUser(ctx, u, []string{role}); err != nil {
		if strings.Contains(err.Error(), "email já existe") {
			return UserView{}, ErrEmailTaken
		}
		if strings.Contains(err.Error(), "username já existe") {
			return UserView{}, ErrConflict
		}
		return UserView{}, err
	}
	_ = s.audit(ctx, actor.UserID, "user_create", u.ID, "success", map[string]string{"role": role})
	return s.userView(ctx, u)
}

// ListUsers devolve uma página.
func (s *Service) ListUsers(ctx context.Context, actor Principal, page, pageSize int) ([]UserView, int, error) {
	if !actor.Allows(PermUsersRead) && !actor.Allows(PermUsersManage) {
		return nil, 0, ErrForbidden
	}
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, 0, ErrValidation
	}
	users, total, err := s.store.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]UserView, 0, len(users))
	for _, u := range users {
		view, err := s.userView(ctx, u)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, view)
	}
	return out, total, nil
}

// GetUser consulta um usuário.
func (s *Service) GetUser(ctx context.Context, actor Principal, id string) (UserView, error) {
	if !actor.Allows(PermUsersRead) && !actor.Allows(PermUsersManage) && actor.UserID != id {
		return UserView{}, ErrForbidden
	}
	u, err := s.store.UserByID(ctx, id)
	if err != nil {
		return UserView{}, ErrNotFound
	}
	return s.userView(ctx, u)
}

// UpdateUser altera nome, e-mail e situação.
// E-mail vazio permanece vazio. E-mail preenchido tem de ser válido.
func (s *Service) UpdateUser(ctx context.Context, actor Principal, id, name, email string, active bool) (UserView, error) {
	if !actor.Allows(PermUsersManage) {
		return UserView{}, ErrForbidden
	}
	if stringsTrim(name) == "" {
		return UserView{}, ErrValidation
	}
	email = strings.TrimSpace(email)
	if email != "" {
		if err := validateEmail(email); err != nil {
			return UserView{}, err
		}
	}
	u, err := s.store.UserByID(ctx, id)
	if err != nil {
		return UserView{}, ErrNotFound
	}
	if u.Active && !active {
		if err := s.guardLastAdmin(ctx, id); err != nil {
			return UserView{}, err
		}
	}
	now := s.now()
	if err := s.store.UpdateUser(ctx, id, stringsTrim(name), email, active, now); err != nil {
		if strings.Contains(err.Error(), "email já existe") {
			return UserView{}, ErrEmailTaken
		}
		return UserView{}, err
	}
	if !active {
		_ = s.store.RevokeUserTokens(ctx, id, now)
	}
	u, err = s.store.UserByID(ctx, id)
	if err != nil {
		return UserView{}, err
	}
	_ = s.audit(ctx, actor.UserID, "user_update", id, "success", nil)
	return s.userView(ctx, u)
}

// DeactivateUser faz o soft delete.
func (s *Service) DeactivateUser(ctx context.Context, actor Principal, id string) error {
	if !actor.Allows(PermUsersManage) {
		return ErrForbidden
	}
	u, err := s.store.UserByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := s.guardLastAdmin(ctx, id); err != nil {
		return err
	}
	now := s.now()
	if err := s.store.UpdateUser(ctx, id, u.DisplayName, u.Email, false, now); err != nil {
		return err
	}
	if err := s.store.RevokeUserTokens(ctx, id, now); err != nil {
		return err
	}
	return s.audit(ctx, actor.UserID, "user_deactivate", id, "success", nil)
}

// SetPassword altera a senha própria ou redefine como administrador.
func (s *Service) SetPassword(ctx context.Context, actor Principal, id, current, next string) error {
	self := actor.UserID == id
	if !self && !actor.Allows(PermUsersManage) {
		return ErrForbidden
	}
	if err := validatePassword(next); err != nil {
		return err
	}
	u, err := s.store.UserByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	mustChange := !self
	if self {
		ok, err := VerifyPassword(current, u.PasswordHash)
		if err != nil {
			return err
		}
		if !ok {
			return ErrInvalidCredentials
		}
	}
	hash, err := HashPassword(next)
	if err != nil {
		return err
	}
	now := s.now()
	if err := s.store.SetPassword(ctx, id, hash, mustChange, now); err != nil {
		return err
	}
	if err := s.store.RevokeUserTokens(ctx, id, now); err != nil {
		return err
	}
	return s.audit(ctx, actor.UserID, "password_change", id, "success", nil)
}

// SetRoles substitui os papéis.
func (s *Service) SetRoles(ctx context.Context, actor Principal, id string, roles []string) (UserView, error) {
	if !actor.Allows(PermUsersManage) {
		return UserView{}, ErrForbidden
	}
	if len(roles) == 0 {
		return UserView{}, ErrValidation
	}
	for _, role := range roles {
		if !knownRole(role) {
			return UserView{}, ErrValidation
		}
	}
	if _, err := s.store.UserByID(ctx, id); err != nil {
		return UserView{}, ErrNotFound
	}
	keepsAdmin := false
	for _, role := range roles {
		if role == "admin" {
			keepsAdmin = true
		}
	}
	if !keepsAdmin {
		if err := s.guardLastAdmin(ctx, id); err != nil {
			return UserView{}, err
		}
	}
	if err := s.store.SetRoles(ctx, id, roles); err != nil {
		return UserView{}, err
	}
	u, err := s.store.UserByID(ctx, id)
	if err != nil {
		return UserView{}, err
	}
	_ = s.audit(ctx, actor.UserID, "user_roles", id, "success", nil)
	return s.userView(ctx, u)
}

// ListRoles devolve o catálogo.
func (s *Service) ListRoles(ctx context.Context, actor Principal) ([]struct{ Code, Description string }, error) {
	if !actor.Allows(PermUsersRead) && !actor.Allows(PermUsersManage) {
		return nil, ErrForbidden
	}
	return s.store.RoleCatalog(ctx)
}

func (s *Service) userView(ctx context.Context, u sqlite.User) (UserView, error) {
	roles, err := s.store.RoleCodes(ctx, u.ID)
	if err != nil {
		return UserView{}, err
	}
	if roles == nil {
		roles = []string{}
	}
	return UserView{
		ID: u.ID, Username: u.Username, Name: u.DisplayName, Email: u.Email, Roles: roles,
		Active: u.Active, MustChangePassword: u.MustChangePassword,
	}, nil
}

func (s *Service) guardLastAdmin(ctx context.Context, userID string) error {
	roles, err := s.store.RoleCodes(ctx, userID)
	if err != nil {
		return err
	}
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
		}
	}
	if !isAdmin {
		return nil
	}
	n, err := s.store.CountActiveAdmins(ctx)
	if err != nil {
		return err
	}
	if n <= 1 {
		return ErrLastAdmin
	}
	return nil
}

func stringsTrim(s string) string {
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
