package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/mcp"
)

func TestAdminAPIAndRefreshReuse(t *testing.T) {
	ts, svc, _ := testServer(t)
	ctx := context.Background()
	if err := svc.CreateInitialAdmin(ctx, "ada.admin", "ada.admin@example.com", "senha-admin-ok"); err != nil {
		t.Fatal(err)
	}
	if err := svc.CreateInitialAdmin(ctx, "outro.admin", "outro.admin@example.com", "senha-admin-ok"); err == nil {
		t.Fatal("segundo administrador foi aceito")
	}

	pair := login(t, ts.URL, "ada.admin", "senha-admin-ok")
	res := apiJSON(t, http.MethodPost, ts.URL+"/api/v1/users", pair.AccessToken, map[string]any{
		"username": "bia.analyst", "name": "Bia", "email": "bia.analyst@example.com", "password": "senha-temporaria", "role": "analyst", "active": true,
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create = %d %s", res.StatusCode, readBody(t, res))
	}
	var created map[string]any
	decode(t, res, &created)
	if created["must_change_password"] != true {
		t.Fatalf("must_change = %v", created["must_change_password"])
	}
	if created["email"] != "bia.analyst@example.com" {
		t.Fatalf("email = %v", created["email"])
	}
	id := created["id"].(string)

	listed := apiJSON(t, http.MethodGet, ts.URL+"/api/v1/users?page=1&page_size=10", pair.AccessToken, nil)
	if listed.StatusCode != http.StatusOK {
		t.Fatalf("list = %d", listed.StatusCode)
	}
	got := apiJSON(t, http.MethodGet, ts.URL+"/api/v1/users/"+id, pair.AccessToken, nil)
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get = %d", got.StatusCode)
	}
	patched := apiJSON(t, http.MethodPatch, ts.URL+"/api/v1/users/"+id, pair.AccessToken, map[string]any{"name": "Bia Silva"})
	if patched.StatusCode != http.StatusOK {
		t.Fatalf("patch = %d %s", patched.StatusCode, readBody(t, patched))
	}

	analyst := login(t, ts.URL, "bia.analyst", "senha-temporaria")
	if !analyst.MustChangePassword {
		t.Fatal("esperava troca obrigatória")
	}
	blocked := apiJSON(t, http.MethodGet, ts.URL+"/api/v1/users", analyst.AccessToken, nil)
	if blocked.StatusCode != http.StatusForbidden {
		t.Fatalf("lista com senha temporária = %d", blocked.StatusCode)
	}
	changed := apiJSON(t, http.MethodPut, ts.URL+"/api/v1/users/"+id+"/password", analyst.AccessToken, map[string]any{
		"current_password": "senha-temporaria", "new_password": "senha-nova-bia1",
	})
	if changed.StatusCode != http.StatusNoContent {
		t.Fatalf("password = %d %s", changed.StatusCode, readBody(t, changed))
	}
	analyst = login(t, ts.URL, "bia.analyst", "senha-nova-bia1")
	denied := apiJSON(t, http.MethodPost, ts.URL+"/api/v1/users", analyst.AccessToken, map[string]any{
		"username": "carla.viewer", "name": "Carla", "password": "senha-temporaria", "role": "viewer",
	})
	if denied.StatusCode != http.StatusForbidden {
		t.Fatalf("analyst create = %d", denied.StatusCode)
	}
	gone := apiJSON(t, http.MethodDelete, ts.URL+"/api/v1/users/"+id, pair.AccessToken, nil)
	if gone.StatusCode != http.StatusNoContent {
		t.Fatalf("delete = %d %s", gone.StatusCode, readBody(t, gone))
	}

	refreshed := refresh(t, ts.URL, pair.RefreshToken)
	reused := refreshRaw(t, ts.URL, pair.RefreshToken)
	if reused.StatusCode != http.StatusUnauthorized {
		t.Fatalf("reuse = %d", reused.StatusCode)
	}
	afterReuse := refreshRaw(t, ts.URL, refreshed.RefreshToken)
	if afterReuse.StatusCode != http.StatusUnauthorized {
		t.Fatalf("família após reuso = %d", afterReuse.StatusCode)
	}
}

func TestMCPRequiresBearer(t *testing.T) {
	ts, svc, cfg := testServer(t)
	res, err := http.Post(ts.URL+"/mcp", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("sem bearer = %d", res.StatusCode)
	}

	ctx := context.Background()
	if err := svc.CreateInitialAdmin(ctx, "ada.admin", "ada.admin@example.com", "senha-admin-ok"); err != nil {
		t.Fatal(err)
	}
	pair := login(t, ts.URL, "ada.admin", "senha-admin-ok")
	client := sdk.NewClient(&sdk.Implementation{Name: "oracle-mcp-test", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, &sdk.StreamableClientTransport{
		Endpoint:   ts.URL + "/mcp",
		HTTPClient: &http.Client{Transport: bearerTripper{token: pair.AccessToken, base: http.DefaultTransport}},
	}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != mcp.ToolCount() {
		t.Fatalf("tools = %d", len(tools.Tools))
	}
	called, err := session.CallTool(ctx, &sdk.CallToolParams{Name: "dim_centro_custo"})
	if err != nil {
		t.Fatal(err)
	}
	if !called.IsError {
		t.Fatal("esperava consulta indisponível")
	}
	body, err := json.Marshal(called)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "consulta indisponível") || strings.Contains(string(body), "SELECT") {
		t.Fatalf("resposta do lookup = %s", body)
	}
	meta, err := http.Get(ts.URL + "/.well-known/oauth-authorization-server")
	if err != nil {
		t.Fatal(err)
	}
	defer meta.Body.Close()
	var doc map[string]any
	if err := json.NewDecoder(meta.Body).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if doc["client_id_metadata_document_supported"] != true || doc["registration_endpoint"] == "" {
		t.Fatalf("metadata = %#v", doc)
	}
	if doc["issuer"] != cfg.Issuer() {
		t.Fatalf("issuer = %v", doc["issuer"])
	}
}

func TestPKCEWithDCRAndCIMD(t *testing.T) {
	ts, svc, cfg := testServer(t)
	ctx := context.Background()
	if err := svc.CreateInitialAdmin(ctx, "ada.admin", "ada.admin@example.com", "senha-admin-ok"); err != nil {
		t.Fatal(err)
	}
	redirect := "http://127.0.0.1:9/callback"
	reg := apiJSON(t, http.MethodPost, ts.URL+"/oauth/register", "", map[string]any{
		"client_name": "Cliente nativo", "redirect_uris": []string{redirect}, "application_type": "native",
	})
	if reg.StatusCode != http.StatusCreated {
		t.Fatalf("dcr = %d %s", reg.StatusCode, readBody(t, reg))
	}
	var client map[string]any
	decode(t, reg, &client)
	code := authorize(t, ts.URL, client["client_id"].(string), redirect, cfg.Resource(), "ada.admin", "senha-admin-ok")
	token := exchange(t, ts.URL, client["client_id"].(string), code, redirect, cfg.Resource())
	if token["access_token"] == "" {
		t.Fatal("access token vazio")
	}

	cimdID := "https://client.example/app.json"
	svc.UseCIMDResolver(func(context.Context, string) (auth.CIMDDocument, error) {
		return auth.CIMDDocument{ClientID: cimdID, ClientName: "CIMD", RedirectURIs: []string{redirect}}, nil
	})
	cimdCode := authorize(t, ts.URL, cimdID, redirect, cfg.Resource(), "ada.admin", "senha-admin-ok")
	cimdToken := exchange(t, ts.URL, cimdID, cimdCode, redirect, cfg.Resource())
	if cimdToken["access_token"] == "" {
		t.Fatal("token CIMD vazio")
	}

	bad := postForm(t, ts.URL+"/oauth/authorize", url.Values{
		"client_id": {cimdID}, "redirect_uri": {"http://127.0.0.1:9/outro"},
		"code_challenge": {pkceChallenge(pkceVerifier)}, "code_challenge_method": {"S256"},
		"state": {"xyz"}, "resource": {cfg.Resource()}, "username": {"ada.admin"}, "password": {"senha-admin-ok"},
	})
	if bad.StatusCode != http.StatusBadRequest {
		t.Fatalf("redirect fora da lista = %d", bad.StatusCode)
	}
	bad.Body.Close()
}

func TestClientRedirects(t *testing.T) {
	ts, svc, cfg := testServer(t)
	ctx := context.Background()
	if err := svc.CreateInitialAdmin(ctx, "ada.admin", "ada.admin@example.com", "senha-admin-ok"); err != nil {
		t.Fatal(err)
	}
	register := func(name string, uris []string) string {
		t.Helper()
		res := apiJSON(t, http.MethodPost, ts.URL+"/oauth/register", "", map[string]any{
			"client_name": name, "redirect_uris": uris,
		})
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("dcr %s = %d %s", name, res.StatusCode, readBody(t, res))
		}
		var client map[string]any
		decode(t, res, &client)
		return client["client_id"].(string)
	}
	cursorID := register("Cursor", []string{
		"cursor://anysphere.cursor-mcp/oauth/callback",
		"https://www.cursor.com/agents/mcp/oauth/callback",
		"http://localhost:8787/callback",
	})
	for _, redirect := range []string{
		"http://localhost:8787/callback",
		"http://127.0.0.1:9999/callback",
	} {
		authorize(t, ts.URL, cursorID, redirect, cfg.Resource(), "ada.admin", "senha-admin-ok")
	}
	onlyWeb := register("Cursor web", []string{"https://www.cursor.com/agents/mcp/oauth/callback"})
	authorize(t, ts.URL, onlyWeb, "http://localhost:8787/callback", cfg.Resource(), "ada.admin", "senha-admin-ok")

	claudeID := register("Claude", []string{"https://claude.ai/api/mcp/auth_callback"})
	authorize(t, ts.URL, claudeID, "https://claude.ai/api/mcp/auth_callback", cfg.Resource(), "ada.admin", "senha-admin-ok")

	chatID := register("ChatGPT", []string{"https://chatgpt.com/connector_platform_oauth_redirect"})
	authorize(t, ts.URL, chatID, "https://chatgpt.com/connector_platform_oauth_redirect", cfg.Resource(), "ada.admin", "senha-admin-ok")

	lanID := register("LAN", []string{"http://192.168.0.10:9/callback"})
	authorize(t, ts.URL, lanID, "http://192.168.0.10:40/callback", cfg.Resource(), "ada.admin", "senha-admin-ok")

	for _, uri := range []string{"http://evil.example/callback", "javascript:alert(1)"} {
		res := apiJSON(t, http.MethodPost, ts.URL+"/oauth/register", "", map[string]any{
			"client_name": "recusado", "redirect_uris": []string{uri},
		})
		body := readBody(t, res)
		if res.StatusCode == http.StatusCreated {
			t.Fatalf("%s foi aceito: %s", uri, body)
		}
	}
}

const pkceVerifier = "verifier-verifier-verifier-verifier"

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

type tokenResponse struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	MustChangePassword bool   `json:"must_change_password"`
}

func login(t *testing.T, base, username, password string) tokenResponse {
	t.Helper()
	res := apiJSON(t, http.MethodPost, base+"/api/v1/auth/login", "", map[string]string{"username": username, "password": password})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login %s = %d %s", username, res.StatusCode, readBody(t, res))
	}
	var pair tokenResponse
	decode(t, res, &pair)
	return pair
}

func refresh(t *testing.T, base, token string) tokenResponse {
	t.Helper()
	res := refreshRaw(t, base, token)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("refresh = %d %s", res.StatusCode, readBody(t, res))
	}
	var pair tokenResponse
	decode(t, res, &pair)
	return pair
}

func refreshRaw(t *testing.T, base, token string) *http.Response {
	t.Helper()
	return apiJSON(t, http.MethodPost, base+"/api/v1/auth/refresh", "", map[string]string{"refresh_token": token})
}

func authorize(t *testing.T, base, clientID, redirect, resource, username, password string) string {
	t.Helper()
	res := postForm(t, base+"/oauth/authorize", url.Values{
		"client_id": {clientID}, "redirect_uri": {redirect},
		"code_challenge": {pkceChallenge(pkceVerifier)}, "code_challenge_method": {"S256"},
		"state": {"xyz"}, "resource": {resource}, "username": {username}, "password": {password},
	})
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		t.Fatalf("authorize = %d %s", res.StatusCode, readBody(t, res))
	}
	loc, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if loc.Query().Get("state") != "xyz" || loc.Query().Get("iss") == "" || loc.Query().Get("code") == "" {
		t.Fatalf("location = %s", res.Header.Get("Location"))
	}
	return loc.Query().Get("code")
}

func exchange(t *testing.T, base, clientID, code, redirect, resource string) map[string]any {
	t.Helper()
	res := postForm(t, base+"/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "client_id": {clientID}, "code": {code},
		"redirect_uri": {redirect}, "code_verifier": {pkceVerifier}, "resource": {resource},
	})
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("token = %d %s", res.StatusCode, readBody(t, res))
	}
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body
}

func apiJSON(t *testing.T, method, rawURL, token string, payload any) *http.Response {
	t.Helper()
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		t.Fatal(err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func postForm(t *testing.T, rawURL string, form url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func decode(t *testing.T, res *http.Response, dst any) {
	t.Helper()
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		t.Fatal(err)
	}
}

func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	return string(b)
}

type bearerTripper struct {
	token string
	base  http.RoundTripper
}

func (b bearerTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	clone := r.Clone(r.Context())
	clone.Header.Set("Authorization", "Bearer "+b.token)
	base := b.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clone)
}
