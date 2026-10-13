package api

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"oracle-mcp/internal/auth"
)

func TestAuthorizePageShowsConsultancyLogo(t *testing.T) {
	dir := t.TempDir()
	raw := []byte{0xff, 0xd8, 0xff, 0xd9}
	if err := os.WriteFile(filepath.Join(dir, "consultoria.jpeg"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	doc := `{"nome":"Assistente","logos":[{"papel":"consultoria","arquivo":"consultoria.jpeg","legenda":"Logo da consultoria"}]}`
	if err := os.WriteFile(filepath.Join(dir, "welcome.json"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	page := authorizePage(auth.AuthorizeInput{ClientID: "c"}, dir)
	want := base64.StdEncoding.EncodeToString(raw)
	if !strings.Contains(page, "data:image/jpeg;base64,"+want) || !strings.Contains(page, `alt="Logo da consultoria"`) {
		t.Fatalf("página sem logo")
	}
	if !strings.Contains(page, `name="client_id"`) || !strings.Contains(page, ">Entrar<") {
		t.Fatal("formulário incompleto")
	}

	empty := authorizePage(auth.AuthorizeInput{}, t.TempDir())
	if strings.Contains(empty, "<img") {
		t.Fatal("página sem config não deveria ter imagem")
	}
}
