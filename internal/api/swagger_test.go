package api

import (
	"net/http"
	"strings"
	"testing"
)

func TestSwaggerPageAndSpec(t *testing.T) {
	ts, _, _ := testServer(t)
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	redir, err := client.Get(ts.URL + "/swagger")
	if err != nil {
		t.Fatal(err)
	}
	redir.Body.Close()
	if redir.StatusCode != http.StatusMovedPermanently || redir.Header.Get("Location") != "/swagger/" {
		t.Fatalf("redirect = %d %s", redir.StatusCode, redir.Header.Get("Location"))
	}

	page, err := http.Get(ts.URL + "/swagger/")
	if err != nil {
		t.Fatal(err)
	}
	html := readBody(t, page)
	if page.StatusCode != http.StatusOK || !strings.Contains(html, "swagger-ui-bundle.js") || !strings.Contains(html, "openapi.yaml") {
		t.Fatalf("página = %d %s", page.StatusCode, html)
	}

	spec, err := http.Get(ts.URL + "/swagger/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	body := readBody(t, spec)
	if spec.StatusCode != http.StatusOK {
		t.Fatalf("spec = %d", spec.StatusCode)
	}
	if !strings.Contains(body, "operationId: listUsers") || !strings.Contains(body, "operationId: createUser") || !strings.Contains(body, "/api/v1/users") {
		t.Fatal("spec sem listar e cadastrar")
	}

	bundle, err := http.Get(ts.URL + "/swagger/swagger-ui-bundle.js")
	if err != nil {
		t.Fatal(err)
	}
	js := readBody(t, bundle)
	if bundle.StatusCode != http.StatusOK || !strings.Contains(js, "SwaggerUIBundle") {
		t.Fatalf("bundle = %d", bundle.StatusCode)
	}
}
