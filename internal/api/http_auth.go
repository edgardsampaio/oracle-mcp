package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/oauthex"

	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/boasvindas"
	"oracle-mcp/internal/config"
	"oracle-mcp/internal/mcp"
)

type httpAPI struct {
	svc        *auth.Service
	cfg        config.Config
	welcomeDir string
}

func (a *httpAPI) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	pair, err := a.svc.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenJSON(pair))
}

func (a *httpAPI) refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	pair, err := a.svc.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenJSON(pair))
}

func (a *httpAPI) logout(w http.ResponseWriter, r *http.Request) {
	p, token, ok := a.bearer(w, r, true)
	if !ok {
		return
	}
	if err := a.svc.Logout(r.Context(), token); err != nil {
		writeAuthErr(w, err)
		return
	}
	_ = p
	w.WriteHeader(http.StatusNoContent)
}

func (a *httpAPI) listUsers(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if size == 0 {
		size = 20
	}
	items, total, err := a.svc.ListUsers(r.Context(), p, page, size)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, userJSON(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out, "page": page, "page_size": size, "total": total})
}

func (a *httpAPI) createUser(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	var body struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Active   *bool  `json:"active"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	view, err := a.svc.CreateUser(r.Context(), p, body.Username, body.Name, body.Email, body.Password, body.Role, active)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, userJSON(view))
}

func (a *httpAPI) getUser(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	view, err := a.svc.GetUser(r.Context(), p, r.PathValue("id"))
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, userJSON(view))
}

func (a *httpAPI) patchUser(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	var body struct {
		Name   string  `json:"name"`
		Email  *string `json:"email"`
		Active *bool   `json:"active"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	current, err := a.svc.GetUser(r.Context(), p, r.PathValue("id"))
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	name := current.Name
	if body.Name != "" {
		name = body.Name
	}
	email := current.Email
	if body.Email != nil {
		email = *body.Email
	}
	active := current.Active
	if body.Active != nil {
		active = *body.Active
	}
	view, err := a.svc.UpdateUser(r.Context(), p, r.PathValue("id"), name, email, active)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, userJSON(view))
}

func (a *httpAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	if err := a.svc.DeactivateUser(r.Context(), p, r.PathValue("id")); err != nil {
		writeAuthErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *httpAPI) setPassword(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, true)
	if !ok {
		return
	}
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	if err := a.svc.SetPassword(r.Context(), p, r.PathValue("id"), body.CurrentPassword, body.NewPassword); err != nil {
		writeAuthErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *httpAPI) setRoles(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	var body struct {
		Roles []string `json:"roles"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	view, err := a.svc.SetRoles(r.Context(), p, r.PathValue("id"), body.Roles)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, userJSON(view))
}

func (a *httpAPI) listRoles(w http.ResponseWriter, r *http.Request) {
	p, ok := a.mustAuth(w, r, false)
	if !ok {
		return
	}
	roles, err := a.svc.ListRoles(r.Context(), p)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	out := make([]any, 0, len(roles))
	for _, role := range roles {
		out = append(out, map[string]string{"code": role.Code, "description": role.Description})
	}
	writeJSON(w, http.StatusOK, map[string]any{"roles": out})
}

func (a *httpAPI) mustAuth(w http.ResponseWriter, r *http.Request, allowMustChange bool) (auth.Principal, bool) {
	p, _, ok := a.bearer(w, r, allowMustChange)
	return p, ok
}

func (a *httpAPI) bearer(w http.ResponseWriter, r *http.Request, allowMustChange bool) (auth.Principal, string, bool) {
	token := bearerToken(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "bearer ausente")
		return auth.Principal{}, "", false
	}
	p, err := a.svc.VerifyAccess(r.Context(), token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "token inválido")
		return auth.Principal{}, "", false
	}
	if p.MustChangePassword && !allowMustChange {
		writeError(w, http.StatusForbidden, "password_change_required", "troca de senha obrigatória")
		return auth.Principal{}, "", false
	}
	return p, token, true
}

func bearerToken(r *http.Request) string {
	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) == 2 && strings.EqualFold(fields[0], "bearer") {
		return fields[1]
	}
	return ""
}

func tokenJSON(pair auth.TokenPair) map[string]any {
	return map[string]any{
		"access_token":         pair.AccessToken,
		"refresh_token":        pair.RefreshToken,
		"token_type":           "Bearer",
		"expires_in":           pair.ExpiresIn,
		"must_change_password": pair.MustChangePassword,
	}
}

func userJSON(v auth.UserView) map[string]any {
	return map[string]any{
		"id": v.ID, "username": v.Username, "name": v.Name, "email": v.Email, "roles": v.Roles,
		"active": v.Active, "must_change_password": v.MustChangePassword,
	}
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "JSON inválido")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}

func writeAuthErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_grant", "credenciais inválidas")
	case errors.Is(err, auth.ErrRateLimited):
		writeError(w, http.StatusTooManyRequests, "rate_limited", "muitas tentativas")
	case errors.Is(err, auth.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden", "sem permissão")
	case errors.Is(err, auth.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "não encontrado")
	case errors.Is(err, auth.ErrEmailTaken):
		writeError(w, http.StatusConflict, "conflict", "e-mail já existe")
	case errors.Is(err, auth.ErrConflict):
		writeError(w, http.StatusConflict, "conflict", "username já existe")
	case errors.Is(err, auth.ErrLastAdmin):
		writeError(w, http.StatusConflict, "last_admin", err.Error())
	case errors.Is(err, auth.ErrValidation):
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
	default:
		var oauth *auth.OAuthError
		if errors.As(err, &oauth) {
			writeError(w, oauth.Status, oauth.Code, oauth.Description)
			return
		}
		writeError(w, http.StatusInternalServerError, "server_error", "erro interno")
	}
}

func rejectMustChange(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := sdkauth.TokenInfoFromContext(r.Context())
		if info != nil {
			if v, _ := info.Extra["must_change_password"].(bool); v {
				writeError(w, http.StatusForbidden, "password_change_required", "troca de senha obrigatória")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func protectedMCP(logger *slog.Logger, cfg config.Config, svc *auth.Service, db bi.Querier, audit bi.Auditor, welcomeDir string) http.Handler {
	inner := rejectMustChange(mcp.Handler(logger, db, audit, welcomeDir, svc, cfg.BI.MaxRows))
	return sdkauth.RequireBearerToken(func(ctx context.Context, token string, r *http.Request) (*sdkauth.TokenInfo, error) {
		p, err := svc.VerifyAccess(ctx, token)
		if err != nil {
			return nil, sdkauth.ErrInvalidToken
		}
		return &sdkauth.TokenInfo{
			Scopes: []string{auth.ScopeMCP}, Expiration: p.ExpiresAt, UserID: p.UserID,
			Extra: map[string]any{"must_change_password": p.MustChangePassword},
		}, nil
	}, &sdkauth.RequireBearerTokenOptions{
		ResourceMetadataURL: cfg.Issuer() + "/.well-known/oauth-protected-resource",
		Scopes:              []string{auth.ScopeMCP},
	})(inner)
}

func resourceMetadata(cfg config.Config) http.Handler {
	return sdkauth.ProtectedResourceMetadataHandler(&oauthex.ProtectedResourceMetadata{
		Resource:               cfg.Resource(),
		AuthorizationServers:   []string{cfg.Issuer()},
		ScopesSupported:        []string{auth.ScopeMCP},
		BearerMethodsSupported: []string{"header"},
		ResourceName:           "Oracle MCP",
	})
}

func (a *httpAPI) asMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	issuer := a.cfg.Issuer()
	writeJSON(w, http.StatusOK, map[string]any{
		"issuer":                                         issuer,
		"authorization_endpoint":                         issuer + "/oauth/authorize",
		"token_endpoint":                                 issuer + "/oauth/token",
		"registration_endpoint":                          issuer + "/oauth/register",
		"response_types_supported":                       []string{"code"},
		"grant_types_supported":                          []string{"authorization_code", "refresh_token"},
		"code_challenge_methods_supported":               []string{"S256"},
		"token_endpoint_auth_methods_supported":          []string{"none"},
		"client_id_metadata_document_supported":          true,
		"scopes_supported":                               []string{auth.ScopeMCP},
		"authorization_response_iss_parameter_supported": true,
	})
}

func (a *httpAPI) authorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "formulário inválido")
		return
	}
	in := auth.AuthorizeInput{
		ClientID: r.Form.Get("client_id"), RedirectURI: r.Form.Get("redirect_uri"),
		CodeChallenge: r.Form.Get("code_challenge"), CodeChallengeMethod: r.Form.Get("code_challenge_method"),
		State: r.Form.Get("state"), Resource: r.Form.Get("resource"), Scope: r.Form.Get("scope"),
		Username: r.Form.Get("username"), Password: r.Form.Get("password"),
	}
	if r.Method == http.MethodGet || in.Username == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, authorizePage(in, a.welcomeDir))
		return
	}
	dest, err := a.svc.Authorize(r.Context(), in)
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	http.Redirect(w, r, dest, http.StatusFound)
}

func authorizePage(in auth.AuthorizeInput, welcomeDir string) string {
	field := func(name, value string) string {
		return `<input type="hidden" name="` + html.EscapeString(name) + `" value="` + html.EscapeString(value) + `">`
	}
	return `<!DOCTYPE html><html lang="pt"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Entrar</title>
<style>
body{margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;background:#f4f4f5;color:#18181b;font:16px/1.4 system-ui,sans-serif}
main{width:min(100% - 2rem,20rem)}
img{display:block;max-width:10rem;max-height:5rem;height:auto;margin:0 auto 1.5rem;object-fit:contain}
label{display:block;margin:0 0 .75rem}
span{display:block;margin:0 0 .25rem}
input{width:100%;box-sizing:border-box;padding:.5rem .6rem;border:1px solid #d4d4d8;border-radius:6px;font:inherit;background:#fff}
button{width:100%;margin-top:.25rem;padding:.55rem .6rem;border:0;border-radius:6px;background:#18181b;color:#fff;font:inherit;cursor:pointer}
</style></head><body><main>
` + logoTag(welcomeDir) + `<form method="post" action="/oauth/authorize">
` + field("client_id", in.ClientID) + field("redirect_uri", in.RedirectURI) + field("code_challenge", in.CodeChallenge) +
		field("code_challenge_method", in.CodeChallengeMethod) + field("state", in.State) + field("resource", in.Resource) + field("scope", in.Scope) + `
<label><span>Usuário</span><input name="username" autocomplete="username"></label>
<label><span>Senha</span><input name="password" type="password" autocomplete="current-password"></label>
<button type="submit">Entrar</button>
</form></main></body></html>`
}

func logoTag(dir string) string {
	mime, legenda, data, ok := boasvindas.ConsultoriaLogo(dir)
	if !ok || len(data) == 0 {
		return ""
	}
	src := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
	return `<img src="` + html.EscapeString(src) + `" alt="` + html.EscapeString(legenda) + `">`
}

func (a *httpAPI) token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "formulário inválido")
		return
	}
	var pair auth.TokenPair
	var err error
	switch r.Form.Get("grant_type") {
	case "authorization_code":
		pair, err = a.svc.Exchange(r.Context(), r.Form.Get("client_id"), r.Form.Get("code"), r.Form.Get("redirect_uri"), r.Form.Get("code_verifier"), r.Form.Get("resource"))
	case "refresh_token":
		if resource := r.Form.Get("resource"); resource != "" && resource != a.cfg.Resource() {
			writeError(w, http.StatusBadRequest, "invalid_target", "resource inválido")
			return
		}
		pair, err = a.svc.Refresh(r.Context(), r.Form.Get("refresh_token"))
	default:
		writeError(w, http.StatusBadRequest, "unsupported_grant_type", "grant_type não suportado")
		return
	}
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": pair.AccessToken, "refresh_token": pair.RefreshToken,
		"token_type": "Bearer", "expires_in": pair.ExpiresIn, "scope": auth.ScopeMCP,
	})
}

func (a *httpAPI) register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClientName      string   `json:"client_name"`
		RedirectURIs    []string `json:"redirect_uris"`
		ApplicationType string   `json:"application_type"`
		AuthMethod      string   `json:"token_endpoint_auth_method"`
		GrantTypes      []string `json:"grant_types"`
		ResponseTypes   []string `json:"response_types"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	client, err := a.svc.RegisterClient(r.Context(), auth.ClientRegistration{
		ClientName: body.ClientName, RedirectURIs: body.RedirectURIs, ApplicationType: body.ApplicationType,
		AuthMethod: body.AuthMethod, GrantTypes: body.GrantTypes, ResponseTypes: body.ResponseTypes,
	})
	if err != nil {
		writeAuthErr(w, err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	writeJSON(w, http.StatusCreated, map[string]any{
		"client_id": client.ClientID, "client_name": client.ClientName, "redirect_uris": client.RedirectURIs,
		"application_type": client.ApplicationType, "token_endpoint_auth_method": client.AuthMethod,
		"grant_types": client.GrantTypes, "response_types": client.ResponseTypes,
	})
}
