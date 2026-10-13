package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadFileDefaults(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")

	cfg, err := LoadFile("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Address != defaultAddress {
		t.Fatalf("address = %s", cfg.Server.Address)
	}
	if cfg.Server.ReadTimeout != defaultReadTimeout || cfg.Server.WriteTimeout != defaultWriteTimeout {
		t.Fatalf("timeouts = %s %s", cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
	}
}

func TestLoadFileYAMLAndEnvOverride(t *testing.T) {
	clearOracleEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	body := []byte("server:\n  address: \"127.0.0.1:9090\"\n  public_base_url: \"https://mcp.exemplo.com\"\n  read_timeout: \"5s\"\n  write_timeout: \"10s\"\ndesconhecido:\n  host: \"ignorado\"\n")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ORACLE_MCP_ADDRESS", "127.0.0.1:7070")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Address != "127.0.0.1:7070" {
		t.Fatalf("address = %s", cfg.Server.Address)
	}
	if cfg.Server.PublicBaseURL != "https://mcp.exemplo.com" {
		t.Fatalf("public base url = %s", cfg.Server.PublicBaseURL)
	}
	if cfg.Server.ReadTimeout != 5*time.Second || cfg.Server.WriteTimeout != 10*time.Second {
		t.Fatalf("timeouts = %s %s", cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
	}
}

func TestLoadFileRejectsInvalid(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "sem-porta")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava erro de endereço")
	}

	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "mcp.exemplo.com")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava erro de URL")
	}

	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "0s")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava erro de timeout")
	}
}

func TestLoadFileMissing(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	if _, err := LoadFile(filepath.Join(t.TempDir(), "ausente.yaml")); err == nil {
		t.Fatal("esperava erro de arquivo ausente")
	}
}

func TestOracleRequiresCredential(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_DB_HOST", "oracle.exemplo.local")
	t.Setenv("ORACLE_MCP_DB_SERVICE", "BI")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava credencial ausente")
	}
}

func TestOracleFromYAMLAndEnv(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_DB_USER", "bi_leitura")
	t.Setenv("ORACLE_MCP_DB_PASSWORD", "segredo-local")
	t.Setenv("ORACLE_MCP_DB_SERVICE", "BI")
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	body := []byte("oracle:\n  host: \"oracle.exemplo.local\"\n  port: 1522\nbi:\n  max_rows: 25\n")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Oracle.Host != "oracle.exemplo.local" || cfg.Oracle.Port != 1522 || cfg.Oracle.User != "bi_leitura" || cfg.Oracle.Password != "segredo-local" {
		t.Fatal("oracle não aplicou yaml e ambiente")
	}
	if cfg.BI.MaxRows != 25 {
		t.Fatalf("max rows = %d", cfg.BI.MaxRows)
	}
	if cfg.Oracle.ServiceName != "BI" || cfg.Oracle.SID != "" {
		t.Fatal("service name não aplicado")
	}
}

func TestOracleSID(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_DB_HOST", "172.16.0.15")
	t.Setenv("ORACLE_MCP_DB_PORT", "1526")
	t.Setenv("ORACLE_MCP_DB_SID", "bi")
	t.Setenv("ORACLE_MCP_DB_USER", "CONSULTOR")
	t.Setenv("ORACLE_MCP_DB_PASSWORD", "segredo-local")
	cfg, err := LoadFile("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Oracle.SID != "bi" || cfg.Oracle.ServiceName != "" || cfg.Oracle.Port != 1526 {
		t.Fatal("sid não aplicado")
	}
}

func TestOracleRequiresOneTarget(t *testing.T) {
	clearOracleEnv(t)
	t.Setenv("ORACLE_MCP_ADDRESS", "")
	t.Setenv("ORACLE_MCP_PUBLIC_BASE_URL", "")
	t.Setenv("ORACLE_MCP_READ_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_WRITE_TIMEOUT", "")
	t.Setenv("ORACLE_MCP_DB_HOST", "oracle.exemplo.local")
	t.Setenv("ORACLE_MCP_DB_USER", "bi_leitura")
	t.Setenv("ORACLE_MCP_DB_PASSWORD", "segredo-local")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava service name ou sid")
	}
	t.Setenv("ORACLE_MCP_DB_SERVICE", "BI")
	t.Setenv("ORACLE_MCP_DB_SID", "bi")
	if _, err := LoadFile(""); err == nil {
		t.Fatal("esperava recusa dos dois alvos")
	}
}

func clearOracleEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"ORACLE_MCP_DB_HOST", "ORACLE_MCP_DB_PORT", "ORACLE_MCP_DB_SERVICE", "ORACLE_MCP_DB_SID",
		"ORACLE_MCP_DB_USER", "ORACLE_MCP_DB_PASSWORD", "ORACLE_MCP_DB_QUERY_TIMEOUT",
		"ORACLE_MCP_DB_MAX_OPEN", "ORACLE_MCP_DB_MAX_IDLE", "ORACLE_MCP_DB_MAX_ROWS", "ORACLE_MCP_DB_MAX_BYTES",
	} {
		t.Setenv(key, "")
	}
}
