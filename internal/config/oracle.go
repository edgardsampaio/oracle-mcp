package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func applyOracleFile(cfg *Config, file fileConfig) error {
	if file.Oracle.Host != "" {
		cfg.Oracle.Host = file.Oracle.Host
	}
	if file.Oracle.Port != 0 {
		cfg.Oracle.Port = file.Oracle.Port
	}
	if file.Oracle.ServiceName != "" {
		cfg.Oracle.ServiceName = file.Oracle.ServiceName
	}
	if file.Oracle.SID != "" {
		cfg.Oracle.SID = file.Oracle.SID
	}
	if file.Oracle.UserEnv != "" {
		cfg.Oracle.UserEnv = file.Oracle.UserEnv
	}
	if file.Oracle.PasswordEnv != "" {
		cfg.Oracle.PasswordEnv = file.Oracle.PasswordEnv
	}
	if file.Oracle.QueryTimeout != "" {
		d, err := time.ParseDuration(file.Oracle.QueryTimeout)
		if err != nil {
			return fmt.Errorf("oracle.query_timeout: %w", err)
		}
		cfg.Oracle.QueryTimeout = d
	}
	if file.Oracle.MaxOpen != 0 {
		cfg.Oracle.MaxOpen = file.Oracle.MaxOpen
	}
	if file.Oracle.MaxIdle != 0 {
		cfg.Oracle.MaxIdle = file.Oracle.MaxIdle
	}
	if file.BI.MaxRows != 0 {
		cfg.BI.MaxRows = file.BI.MaxRows
	}
	if file.BI.MaxResponseBytes != 0 {
		cfg.BI.MaxResponseBytes = file.BI.MaxResponseBytes
	}
	return nil
}

func applyOracleEnv(cfg *Config) error {
	if v := os.Getenv("ORACLE_MCP_DB_HOST"); v != "" {
		cfg.Oracle.Host = v
	}
	if v := os.Getenv("ORACLE_MCP_DB_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_PORT: %w", err)
		}
		cfg.Oracle.Port = n
	}
	if v := os.Getenv("ORACLE_MCP_DB_SERVICE"); v != "" {
		cfg.Oracle.ServiceName = v
	}
	if v := os.Getenv("ORACLE_MCP_DB_SID"); v != "" {
		cfg.Oracle.SID = v
	}
	if v := os.Getenv("ORACLE_MCP_DB_QUERY_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_QUERY_TIMEOUT: %w", err)
		}
		cfg.Oracle.QueryTimeout = d
	}
	if v := os.Getenv("ORACLE_MCP_DB_MAX_OPEN"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_MAX_OPEN: %w", err)
		}
		cfg.Oracle.MaxOpen = n
	}
	if v := os.Getenv("ORACLE_MCP_DB_MAX_IDLE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_MAX_IDLE: %w", err)
		}
		cfg.Oracle.MaxIdle = n
	}
	if v := os.Getenv("ORACLE_MCP_DB_MAX_ROWS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_MAX_ROWS: %w", err)
		}
		cfg.BI.MaxRows = n
	}
	if v := os.Getenv("ORACLE_MCP_DB_MAX_BYTES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_DB_MAX_BYTES: %w", err)
		}
		cfg.BI.MaxResponseBytes = n
	}
	return nil
}

func (cfg *Config) validateOracle() error {
	if cfg.BI.MaxRows <= 0 || cfg.BI.MaxResponseBytes <= 0 {
		return fmt.Errorf("limites de consulta devem ser positivos")
	}
	if !cfg.Oracle.Enabled() {
		return nil
	}
	if cfg.Oracle.Port < 1 || cfg.Oracle.Port > 65535 {
		return fmt.Errorf("oracle.port inválida")
	}
	if cfg.Oracle.ServiceName != "" && cfg.Oracle.SID != "" {
		return fmt.Errorf("informe oracle.service_name ou oracle.sid")
	}
	if cfg.Oracle.ServiceName == "" && cfg.Oracle.SID == "" {
		return fmt.Errorf("oracle.service_name ou oracle.sid é obrigatório")
	}
	if cfg.Oracle.QueryTimeout <= 0 || cfg.Oracle.MaxOpen <= 0 || cfg.Oracle.MaxIdle < 0 || cfg.Oracle.MaxIdle > cfg.Oracle.MaxOpen {
		return fmt.Errorf("pool Oracle inválido")
	}
	cfg.Oracle.User = os.Getenv(cfg.Oracle.UserEnv)
	cfg.Oracle.Password = os.Getenv(cfg.Oracle.PasswordEnv)
	if cfg.Oracle.User == "" || cfg.Oracle.Password == "" {
		return fmt.Errorf("credencial Oracle ausente nas variáveis %s e %s", cfg.Oracle.UserEnv, cfg.Oracle.PasswordEnv)
	}
	return nil
}
