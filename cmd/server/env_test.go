package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	body := "# comentario\nORACLE_MCP_DB_SID=bi\nORACLE_MCP_DB_PASSWORD=segredo$\nJA_DEFINIDA=arquivo\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JA_DEFINIDA", "ambiente")
	t.Setenv("ORACLE_MCP_DB_SID", "")
	t.Setenv("ORACLE_MCP_DB_PASSWORD", "")
	loadDotEnv(path)
	if os.Getenv("ORACLE_MCP_DB_SID") != "bi" {
		t.Fatal("sid não carregou")
	}
	if os.Getenv("ORACLE_MCP_DB_PASSWORD") != "segredo$" {
		t.Fatal("senha não preservou o caractere final")
	}
	if os.Getenv("JA_DEFINIDA") != "ambiente" {
		t.Fatal("variável do ambiente foi substituída")
	}
}
