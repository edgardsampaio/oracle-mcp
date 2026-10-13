package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"

	"oracle-mcp/internal/config"
	"oracle-mcp/internal/identity/sqlite"
)

const (
	// ScopeMCP é o escopo do resource server.
	ScopeMCP = "mcp"
	// PermUsersRead lista usuários.
	PermUsersRead = "users:read"
	// PermUsersManage altera usuários.
	PermUsersManage = "users:manage"

	maxLoginFailures = 5
	loginWindow      = 15 * time.Minute
	authCodeTTL      = time.Minute
)

var (
	// ErrInvalidCredentials cobre usuário ausente e senha errada.
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	// ErrRateLimited indica bloqueio temporário de login.
	ErrRateLimited = errors.New("muitas tentativas")
	// ErrAdminExists impede um segundo bootstrap.
	ErrAdminExists = errors.New("já existe um administrador; a recuperação controlada não é feita por este comando")
	// ErrNotFound indica usuário ausente.
	ErrNotFound = errors.New("não encontrado")
	// ErrConflict indica username repetido.
	ErrConflict = errors.New("username já existe")
	// ErrEmailTaken indica e-mail repetido.
	ErrEmailTaken = errors.New("e-mail já existe")
	// ErrForbidden indica falta de permissão.
	ErrForbidden = errors.New("sem permissão")
	// ErrLastAdmin protege o último administrador ativo.
	ErrLastAdmin = errors.New("não é possível remover o último administrador")
	// ErrValidation indica entrada rejeitada.
	ErrValidation = errors.New("dados inválidos")
)

// OAuthError é um erro do authorization server, sem segredo.
type OAuthError struct {
	Code        string
	Description string
	Status      int
}

func (e *OAuthError) Error() string { return e.Code }

// TokenPair são os segredos emitidos uma única vez.
type TokenPair struct {
	AccessToken        string
	RefreshToken       string
	ExpiresIn          int
	MustChangePassword bool
}

// Principal é o usuário autenticado pelo access token.
type Principal struct {
	UserID             string
	Username           string
	Permissions        map[string]bool
	MustChangePassword bool
	ExpiresAt          time.Time
	AccessTokenID      string
}

// Allows informa se a permissão está presente.
func (p Principal) Allows(code string) bool { return p.Permissions[code] }

// UserView é a representação administrativa, sem segredo.
type UserView struct {
	ID                 string
	Username           string
	Name               string
	Email              string
	Roles              []string
	Active             bool
	MustChangePassword bool
}

// Service aplica identidade, sessão e OAuth.
type Service struct {
	store     *sqlite.Store
	cfg       config.Config
	logger    *slog.Logger
	dummyHash string
	now       func() time.Time
	fetchCIMD func(context.Context, string) (CIMDDocument, error)
	cache     cimdCache
}

// New prepara o serviço. O hash fictício reduz a diferença de tempo quando o usuário não existe.
func New(store *sqlite.Store, cfg config.Config, logger *slog.Logger) (*Service, error) {
	dummy, err := HashPassword(strings.Repeat("not-a-user-password", 2))
	if err != nil {
		return nil, err
	}
	s := &Service{
		store:     store,
		cfg:       cfg,
		logger:    logger,
		dummyHash: dummy,
		now:       time.Now,
	}
	s.fetchCIMD = s.fetchCIMDDefault
	return s, nil
}

// UseCIMDResolver substitui a busca do documento. Usado pelos testes.
func (s *Service) UseCIMDResolver(fn func(context.Context, string) (CIMDDocument, error)) {
	s.fetchCIMD = fn
}

func (s *Service) fetchCIMDDefault(ctx context.Context, clientID string) (CIMDDocument, error) {
	if doc, ok := s.cache.get(clientID, s.now()); ok {
		return doc, nil
	}
	doc, exp, err := (CIMDFetcher{}).Fetch(ctx, clientID)
	if err != nil {
		return CIMDDocument{}, err
	}
	if exp.After(s.now()) {
		s.cache.put(clientID, doc, exp)
	}
	return doc, nil
}

// CreateInitialAdmin cria o primeiro administrador e recusa um segundo.
func (s *Service) CreateInitialAdmin(ctx context.Context, username, email, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	n, err := s.store.CountActiveAdmins(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrAdminExists
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	now := s.now()
	u := sqlite.User{
		ID: uuid.NewString(), Username: username, DisplayName: username, Email: strings.TrimSpace(email),
		PasswordHash: hash, Active: true, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.store.InsertUser(ctx, u, []string{"admin"}); err != nil {
		return err
	}
	return s.audit(ctx, u.ID, "admin_bootstrap", u.ID, "success", nil)
}

// Login emite access e refresh tokens.
func (s *Service) Login(ctx context.Context, username, password string) (TokenPair, error) {
	user, err := s.authenticate(ctx, username, password)
	if err != nil {
		return TokenPair{}, err
	}
	pair, err := s.issueTokens(ctx, user, "")
	if err != nil {
		return TokenPair{}, err
	}
	_ = s.audit(ctx, user.ID, "login", user.ID, "success", nil)
	return pair, nil
}

func (s *Service) authenticate(ctx context.Context, username, password string) (sqlite.User, error) {
	since := s.now().Add(-loginWindow)
	fails, err := s.store.RecentFailedLogins(ctx, username, since)
	if err != nil {
		return sqlite.User{}, err
	}
	if fails >= maxLoginFailures {
		return sqlite.User{}, ErrRateLimited
	}
	user, err := s.store.UserByUsername(ctx, username)
	hash := s.dummyHash
	if err == nil {
		hash = user.PasswordHash
	} else if !errors.Is(err, sqlite.ErrNotFound) {
		return sqlite.User{}, err
	}
	ok, verr := VerifyPassword(password, hash)
	if verr != nil {
		return sqlite.User{}, verr
	}
	if !ok || user.ID == "" || !user.Active {
		_ = s.store.InsertLoginAttempt(ctx, username, false, s.now())
		_ = s.audit(ctx, "", "login", username, "failure", nil)
		return sqlite.User{}, ErrInvalidCredentials
	}
	if err := s.store.InsertLoginAttempt(ctx, username, true, s.now()); err != nil {
		return sqlite.User{}, err
	}
	return user, nil
}

// Refresh rotaciona o refresh token. Reuso de token revogado invalida a família.
func (s *Service) Refresh(ctx context.Context, refreshPlain string) (TokenPair, error) {
	rec, err := s.store.RefreshByHash(ctx, HashToken(refreshPlain))
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return TokenPair{}, ErrInvalidCredentials
		}
		return TokenPair{}, err
	}
	now := s.now()
	if rec.RevokedAt != nil {
		child, err := s.store.RefreshHasChild(ctx, rec.ID)
		if err != nil {
			return TokenPair{}, err
		}
		if child {
			_ = s.store.RevokeUserTokens(ctx, rec.UserID, now)
			_ = s.audit(ctx, rec.UserID, "refresh_reuse", rec.UserID, "failure", nil)
		}
		return TokenPair{}, ErrInvalidCredentials
	}
	if !rec.ExpiresAt.After(now) {
		return TokenPair{}, ErrInvalidCredentials
	}
	user, err := s.store.UserByID(ctx, rec.UserID)
	if err != nil || !user.Active {
		return TokenPair{}, ErrInvalidCredentials
	}
	if err := s.store.RevokeRefresh(ctx, rec.ID, now); err != nil {
		return TokenPair{}, err
	}
	pair, err := s.issueTokensRotated(ctx, user, rec.ClientID, rec.ID)
	if err != nil {
		return TokenPair{}, err
	}
	_ = s.audit(ctx, user.ID, "refresh", user.ID, "success", nil)
	return pair, nil
}

// Logout revoga os tokens do usuário autenticado.
func (s *Service) Logout(ctx context.Context, accessPlain string) error {
	p, err := s.VerifyAccess(ctx, accessPlain)
	if err != nil {
		return err
	}
	if err := s.store.RevokeUserTokens(ctx, p.UserID, s.now()); err != nil {
		return err
	}
	return s.audit(ctx, p.UserID, "logout", p.UserID, "success", nil)
}

// VerifyAccess confere o bearer opaco, a audiência e a situação do usuário.
func (s *Service) VerifyAccess(ctx context.Context, plain string) (Principal, error) {
	rec, err := s.store.AccessByHash(ctx, HashToken(plain))
	if err != nil {
		return Principal{}, ErrInvalidCredentials
	}
	now := s.now()
	if rec.RevokedAt != nil || !rec.ExpiresAt.After(now) || rec.Audience != s.cfg.Resource() {
		return Principal{}, ErrInvalidCredentials
	}
	if !scopeHas(rec.Scope, ScopeMCP) {
		return Principal{}, ErrInvalidCredentials
	}
	user, err := s.store.UserByID(ctx, rec.UserID)
	if err != nil || !user.Active {
		return Principal{}, ErrInvalidCredentials
	}
	perms, err := s.store.PermissionCodes(ctx, user.ID)
	if err != nil {
		return Principal{}, err
	}
	set := map[string]bool{}
	for _, p := range perms {
		set[p] = true
	}
	return Principal{
		UserID: user.ID, Username: user.Username, Permissions: set,
		MustChangePassword: rec.MustChangePassword || user.MustChangePassword,
		ExpiresAt:          rec.ExpiresAt, AccessTokenID: rec.ID,
	}, nil
}

func (s *Service) issueTokens(ctx context.Context, user sqlite.User, clientID string) (TokenPair, error) {
	return s.issueTokensRotated(ctx, user, clientID, "")
}

func (s *Service) issueTokensRotated(ctx context.Context, user sqlite.User, clientID, rotatedFrom string) (TokenPair, error) {
	now := s.now()
	accessPlain, accessHash, err := RandomToken()
	if err != nil {
		return TokenPair{}, err
	}
	refreshPlain, refreshHash, err := RandomToken()
	if err != nil {
		return TokenPair{}, err
	}
	accessExp := now.Add(s.cfg.Auth.AccessTokenTTL)
	if err := s.store.InsertAccessToken(ctx, uuid.NewString(), user.ID, clientID, accessHash, ScopeMCP, s.cfg.Resource(), user.MustChangePassword, accessExp, now); err != nil {
		return TokenPair{}, err
	}
	if err := s.store.InsertRefreshToken(ctx, uuid.NewString(), user.ID, clientID, refreshHash, now.Add(s.cfg.Auth.RefreshTokenTTL), now, rotatedFrom); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken: accessPlain, RefreshToken: refreshPlain,
		ExpiresIn: int(s.cfg.Auth.AccessTokenTTL.Seconds()), MustChangePassword: user.MustChangePassword,
	}, nil
}

func (s *Service) audit(ctx context.Context, actor, action, target, outcome string, meta map[string]string) error {
	return s.store.InsertAudit(ctx, sqlite.AuditEvent{
		ActorID: actor, Action: action, Target: target, Outcome: outcome, Metadata: meta,
	}, s.now())
}

func validatePassword(password string) error {
	if len(password) < 6 || len(password) > 128 {
		return fmt.Errorf("%w: a senha deve ter entre 6 e 128 caracteres", ErrValidation)
	}
	return nil
}

// Email devolve o e-mail gravado do usuário. Vazio significa cadastro incompleto para o RLS.
func (s *Service) Email(ctx context.Context, userID string) (string, error) {
	u, err := s.store.UserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	return u.Email, nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	at := strings.LastIndex(email, "@")
	if at <= 0 || at >= len(email)-1 || strings.ContainsAny(email, " \t\r\n") || !strings.Contains(email[at+1:], ".") {
		return fmt.Errorf("%w: e-mail inválido", ErrValidation)
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 64 {
		return fmt.Errorf("%w: username inválido", ErrValidation)
	}
	for _, r := range username {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-' {
			continue
		}
		return fmt.Errorf("%w: username inválido", ErrValidation)
	}
	return nil
}

func knownRole(role string) bool {
	return role == "admin" || role == "analyst" || role == "viewer"
}

func scopeHas(scope, want string) bool {
	for _, part := range strings.Fields(scope) {
		if part == want {
			return true
		}
	}
	return false
}
