package boasvindas

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

func TestGreetReadsTextAndImage(t *testing.T) {
	dir := t.TempDir()
	png := []byte{0x89, 'P', 'N', 'G'}
	if err := os.WriteFile(filepath.Join(dir, "consultoria.png"), png, 0o644); err != nil {
		t.Fatal(err)
	}
	writeDoc(t, dir, []logo{
		{Papel: "consultoria", Arquivo: "consultoria.png", Legenda: "Consultoria"},
		{Papel: "cliente", Arquivo: "cliente.png", Legenda: "Cliente"},
	})
	got, err := Greet(context.Background(), dir, nil, nil, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Texto, "não invente número") || !strings.Contains(got.Texto, "cliente.png") {
		t.Fatalf("texto = %s", got.Texto)
	}
	if len(got.Imagens) != 1 || got.Imagens[0].MIME != "image/png" || string(got.Imagens[0].Dados) != string(png) {
		t.Fatalf("imagens = %#v", got.Imagens)
	}
}

func TestToolResultKeepsBytesOutOfStructuredOutput(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G'}
	call, out := toolResult(Result{
		Texto: "Olá.",
		Imagens: []Image{{
			Papel: "consultoria", Legenda: "Consultoria", MIME: "image/png", Dados: png,
		}},
	})
	raw, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "Dados") || strings.Contains(string(raw), "Imagens") {
		t.Fatalf("resultado estruturado = %s", raw)
	}
	if len(call.Content) != 2 {
		t.Fatalf("conteúdo = %d", len(call.Content))
	}
	img, ok := call.Content[1].(*sdk.ImageContent)
	if !ok || img.MIMEType != "image/png" || string(img.Data) != string(png) {
		t.Fatalf("imagem = %#v", call.Content[1])
	}
}

func TestGreetRejectsParentPath(t *testing.T) {
	root := t.TempDir()
	secret := []byte("segredo-png")
	if err := os.WriteFile(filepath.Join(root, "segredo.png"), secret, 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "config")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeDoc(t, dir, []logo{{Papel: "consultoria", Arquivo: "../segredo.png", Legenda: "Fora"}})
	got, err := Greet(context.Background(), dir, nil, nil, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Imagens) != 0 || strings.Contains(got.Texto, "segredo-png") {
		t.Fatalf("vazou = %#v", got)
	}
	if !strings.Contains(got.Texto, "não encontrada") {
		t.Fatalf("texto = %s", got.Texto)
	}
}

func TestGreetRequiresIdentity(t *testing.T) {
	_, err := Greet(context.Background(), t.TempDir(), nil, nil, "")
	if !errors.Is(err, bi.ErrIdentity) {
		t.Fatal(err)
	}
}

func writeDoc(t *testing.T, dir string, logos []logo) {
	t.Helper()
	doc := fileDoc{
		Nome:          "Assistente de BI",
		Personalidade: "não invente número",
		Mensagem:      "Olá.",
		Logos:         logos,
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, fileName), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}
