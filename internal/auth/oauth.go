package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"oracle-mcp/internal/identity/sqlite"
)

// AuthorizeInput é o pedido do authorization code.
type AuthorizeInput struct {
	ClientID            string
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod string
	State               string
	Resource            string
	Scope               string
	Username            string
	Password            string
}

// ClientRegistration é o corpo do DCR.
type ClientRegistration struct {
	ClientName      string
	RedirectURIs    []string
	ApplicationType string
	AuthMethod      string
	GrantTypes      []string
	ResponseTypes   []string
}

// RegisteredClient é a resposta do DCR, sem segredo.
type RegisteredClient struct {
	ClientID        string
	ClientName      string
	RedirectURIs    []string
	ApplicationType string
	AuthMethod      string
	GrantTypes      []string
	ResponseTypes   []string
}

// Authorize valida o cliente e, com senha correta, devolve a URL de redirecionamento.
func (s *Service) Authorize(ctx context.Context, in AuthorizeInput) (string, error) {
	if in.CodeChallengeMethod != "S256" || in.CodeChallenge == "" || in.State == "" {
		return "", &OAuthError{Code: "invalid_request", Description: "PKCE S256 e state são obrigatórios", Status: 400}
	}
	if in.Resource != s.cfg.Resource() {
		return "", &OAuthError{Code: "invalid_target", Description: "resource inválido", Status: 400}
	}
	scope := in.Scope
	if scope == "" {
		scope = ScopeMCP
	}
	if !scopeHas(scope, ScopeMCP) {
		return "", &OAuthError{Code: "invalid_scope", Description: "escopo não suportado", Status: 400}
	}
	redirects, err := s.clientRedirects(ctx, in.ClientID)
	if err != nil {
		return "", err
	}
	if !redirectAllowed(redirects, in.RedirectURI) {
		return "", &OAuthError{Code: "invalid_request", Description: "redirect_uri não autorizado", Status: 400}
	}
	user, err := s.authenticate(ctx, in.Username, in.Password)
	if err != nil {
		if errors.Is(err, ErrRateLimited) {
			return "", &OAuthError{Code: "access_denied", Description: "muitas tentativas", Status: 429}
		}
		return "", &OAuthError{Code: "access_denied", Description: "credenciais inválidas", Status: 401}
	}
	plain, hash, err := RandomToken()
	if err != nil {
		return "", err
	}
	now := s.now()
	if err := s.store.InsertAuthCode(ctx, hash, sqlite.AuthCode{
		ClientID: in.ClientID, UserID: user.ID, RedirectURI: in.RedirectURI,
		CodeChallenge: in.CodeChallenge, Scope: ScopeMCP, Resource: in.Resource,
		ExpiresAt: now.Add(authCodeTTL),
	}, now); err != nil {
		return "", err
	}
	_ = s.audit(ctx, user.ID, "oauth_authorize", in.ClientID, "success", nil)
	dest, err := url.Parse(in.RedirectURI)
	if err != nil {
		return "", err
	}
	q := dest.Query()
	q.Set("code", plain)
	q.Set("state", in.State)
	q.Set("iss", s.cfg.Issuer())
	dest.RawQuery = q.Encode()
	return dest.String(), nil
}

// Exchange troca o authorization code por tokens.
func (s *Service) Exchange(ctx context.Context, clientID, code, redirectURI, verifier, resource string) (TokenPair, error) {
	if resource != s.cfg.Resource() {
		return TokenPair{}, &OAuthError{Code: "invalid_target", Description: "resource inválido", Status: 400}
	}
	rec, err := s.store.ConsumeAuthCode(ctx, HashToken(code), s.now())
	if err != nil {
		return TokenPair{}, &OAuthError{Code: "invalid_grant", Description: "authorization code inválido", Status: 400}
	}
	if rec.ClientID != clientID || rec.RedirectURI != redirectURI || rec.Resource != resource {
		return TokenPair{}, &OAuthError{Code: "invalid_grant", Description: "authorization code inválido", Status: 400}
	}
	if !pkceS256(verifier, rec.CodeChallenge) {
		return TokenPair{}, &OAuthError{Code: "invalid_grant", Description: "PKCE inválido", Status: 400}
	}
	user, err := s.store.UserByID(ctx, rec.UserID)
	if err != nil || !user.Active {
		return TokenPair{}, &OAuthError{Code: "invalid_grant", Description: "usuário inativo", Status: 400}
	}
	return s.issueTokens(ctx, user, clientID)
}

// RegisterClient executa o Dynamic Client Registration para cliente público.
func (s *Service) RegisterClient(ctx context.Context, in ClientRegistration) (RegisteredClient, error) {
	if in.AuthMethod != "" && in.AuthMethod != "none" {
		return RegisteredClient{}, &OAuthError{Code: "invalid_client_metadata", Description: "apenas cliente público", Status: 400}
	}
	appType := in.ApplicationType
	if appType == "" {
		appType = "native"
	}
	if appType != "native" && appType != "web" {
		return RegisteredClient{}, &OAuthError{Code: "invalid_client_metadata", Description: "application_type inválido", Status: 400}
	}
	if strings.TrimSpace(in.ClientName) == "" || len(in.RedirectURIs) == 0 || len(in.RedirectURIs) > 10 {
		return RegisteredClient{}, &OAuthError{Code: "invalid_client_metadata", Description: "client_name e redirect_uris são obrigatórios", Status: 400}
	}
	for _, uri := range in.RedirectURIs {
		if err := validateRedirect(uri, appType); err != nil {
			return RegisteredClient{}, &OAuthError{Code: "invalid_redirect_uri", Description: err.Error(), Status: 400}
		}
	}
	if len(in.GrantTypes) > 0 && !subset(in.GrantTypes, []string{"authorization_code", "refresh_token"}) {
		return RegisteredClient{}, &OAuthError{Code: "invalid_client_metadata", Description: "grant_type não suportado", Status: 400}
	}
	if len(in.ResponseTypes) > 0 && !subset(in.ResponseTypes, []string{"code"}) {
		return RegisteredClient{}, &OAuthError{Code: "invalid_client_metadata", Description: "response_type não suportado", Status: 400}
	}
	clientID := uuid.NewString()
	now := s.now()
	if err := s.store.InsertOAuthClient(ctx, sqlite.OAuthClient{
		ClientID: clientID, ClientName: strings.TrimSpace(in.ClientName),
		RedirectURIs: in.RedirectURIs, ApplicationType: appType,
	}, now); err != nil {
		return RegisteredClient{}, err
	}
	_ = s.audit(ctx, "", "oauth_register", clientID, "success", map[string]string{"application_type": appType})
	grants := in.GrantTypes
	if len(grants) == 0 {
		grants = []string{"authorization_code", "refresh_token"}
	}
	responses := in.ResponseTypes
	if len(responses) == 0 {
		responses = []string{"code"}
	}
	return RegisteredClient{
		ClientID: clientID, ClientName: strings.TrimSpace(in.ClientName), RedirectURIs: in.RedirectURIs,
		ApplicationType: appType, AuthMethod: "none", GrantTypes: grants, ResponseTypes: responses,
	}, nil
}

func (s *Service) clientRedirects(ctx context.Context, clientID string) ([]string, error) {
	if isCIMDID(clientID) {
		doc, err := s.fetchCIMD(ctx, clientID)
		if err != nil {
			return nil, &OAuthError{Code: "invalid_client", Description: "documento CIMD rejeitado", Status: 400}
		}
		return doc.RedirectURIs, nil
	}
	client, err := s.store.OAuthClientByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return nil, &OAuthError{Code: "invalid_client", Description: "cliente desconhecido", Status: 400}
		}
		return nil, err
	}
	return client.RedirectURIs, nil
}

func pkceS256(verifier, challenge string) bool {
	if verifier == "" {
		return false
	}
	sum := sha256.Sum256([]byte(verifier))
	got := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(got), []byte(challenge)) == 1
}

func containsExact(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}

// cursorRedirects é o conjunto que o Cursor envia no registro e depois usa um deles na autorização.
var cursorRedirects = []string{
	"cursor://anysphere.cursor-mcp/oauth/callback",
	"https://www.cursor.com/agents/mcp/oauth/callback",
	"http://localhost:8787/callback",
}

func redirectAllowed(registered []string, requested string) bool {
	if containsExact(registered, requested) || localHTTPMatch(registered, requested) {
		return true
	}
	if !containsExact(cursorRedirects, requested) {
		return false
	}
	for _, item := range registered {
		if containsExact(cursorRedirects, item) {
			return true
		}
	}
	return false
}

func localHTTPMatch(registered []string, requested string) bool {
	req, err := url.Parse(requested)
	if err != nil || req.Scheme != "http" || req.Fragment != "" || !localHost(req.Hostname()) {
		return false
	}
	for _, item := range registered {
		reg, err := url.Parse(item)
		if err != nil || reg.Scheme != "http" || !localHost(reg.Hostname()) {
			continue
		}
		if reg.EscapedPath() != req.EscapedPath() || reg.RawQuery != req.RawQuery {
			continue
		}
		if sameLocalHost(reg.Hostname(), req.Hostname()) {
			return true
		}
	}
	return false
}

func sameLocalHost(a, b string) bool {
	if a == b {
		return true
	}
	if loopbackHost(a) && loopbackHost(b) {
		return true
	}
	ia, ib := net.ParseIP(a), net.ParseIP(b)
	return ia != nil && ib != nil && ia.Equal(ib)
}

func subset(got, allow []string) bool {
	ok := map[string]bool{}
	for _, item := range allow {
		ok[item] = true
	}
	for _, item := range got {
		if !ok[item] {
			return false
		}
	}
	return true
}

func validateRedirect(raw, _ string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Fragment != "" || u.User != nil {
		return errors.New("redirect_uri inválido")
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		if u.Host == "" {
			return errors.New("redirect_uri inválido")
		}
		return nil
	case "http":
		if u.Host != "" && localHost(u.Hostname()) {
			return nil
		}
		return errors.New("redirect_uri não permitido para este application_type")
	default:
		if !privateScheme(u.Scheme) || u.Host == "" && u.Opaque == "" && u.Path == "" {
			return errors.New("redirect_uri inválido")
		}
		return nil
	}
}

func privateScheme(scheme string) bool {
	switch strings.ToLower(scheme) {
	case "javascript", "data", "file", "blob", "vbscript", "about":
		return false
	default:
		return scheme != ""
	}
}

func localHost(host string) bool {
	if loopbackHost(host) {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsPrivate() || ip.IsLinkLocalUnicast())
}

func loopbackHost(host string) bool {
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
