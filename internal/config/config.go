package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultAddress      = "127.0.0.1:8080"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 30 * time.Second
	defaultSQLitePath   = "data/identity.db"
	defaultAccessTTL    = 15 * time.Minute
	defaultRefreshTTL   = 168 * time.Hour
)

// Config é a configuração validada do processo.
// write_timeout é aceito e guardado para limite por requisição nas fases
// seguintes. O http.Server não o aplica: o transporte MCP stateful mantém
// um GET de eventos aberto.
type Config struct {
	Server Server
	Auth   Auth
	Oracle Oracle
	BI     BI
}

// Server reúne o bind HTTP e os tempos limite.
type Server struct {
	Address       string
	PublicBaseURL string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// Auth reúne o SQLite de identidade e a vida dos tokens.
type Auth struct {
	SQLitePath      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Oracle é a fonte somente leitura. Host vazio desliga o adaptador.
type Oracle struct {
	Host         string
	Port         int
	ServiceName  string
	SID          string
	User         string
	Password     string
	QueryTimeout time.Duration
	MaxOpen      int
	MaxIdle      int
	UserEnv      string
	PasswordEnv  string
}

// Enabled informa se a conexão Oracle deve ser aberta.
func (o Oracle) Enabled() bool {
	return strings.TrimSpace(o.Host) != ""
}

// BI guarda os limites das consultas cadastradas.
type BI struct {
	MaxRows          int
	MaxResponseBytes int
}

type fileConfig struct {
	Server struct {
		Address       string `yaml:"address"`
		PublicBaseURL string `yaml:"public_base_url"`
		ReadTimeout   string `yaml:"read_timeout"`
		WriteTimeout  string `yaml:"write_timeout"`
	} `yaml:"server"`
	Auth struct {
		SQLitePath      string `yaml:"sqlite_path"`
		AccessTokenTTL  string `yaml:"access_token_ttl"`
		RefreshTokenTTL string `yaml:"refresh_token_ttl"`
	} `yaml:"auth"`
	Oracle struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		ServiceName  string `yaml:"service_name"`
		SID          string `yaml:"sid"`
		UserEnv      string `yaml:"username_env"`
		PasswordEnv  string `yaml:"password_env"`
		QueryTimeout string `yaml:"query_timeout"`
		MaxOpen      int    `yaml:"max_open_connections"`
		MaxIdle      int    `yaml:"max_idle_connections"`
	} `yaml:"oracle"`
	BI struct {
		MaxRows          int `yaml:"max_rows"`
		MaxResponseBytes int `yaml:"max_response_bytes"`
	} `yaml:"bi"`
}

// Load lê ORACLE_MCP_CONFIG, aplica variáveis de ambiente e valida.
func Load() (Config, error) {
	return LoadFile(os.Getenv("ORACLE_MCP_CONFIG"))
}

// LoadFile lê um YAML opcional. path vazio usa só padrão e ambiente.
func LoadFile(path string) (Config, error) {
	cfg := Config{
		Server: Server{
			Address:      defaultAddress,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		Auth: Auth{
			SQLitePath:      defaultSQLitePath,
			AccessTokenTTL:  defaultAccessTTL,
			RefreshTokenTTL: defaultRefreshTTL,
		},
		Oracle: Oracle{
			Port:         1521,
			QueryTimeout: 30 * time.Second,
			MaxOpen:      10,
			MaxIdle:      5,
			UserEnv:      "ORACLE_MCP_DB_USER",
			PasswordEnv:  "ORACLE_MCP_DB_PASSWORD",
		},
		BI: BI{MaxRows: 1000, MaxResponseBytes: 2 << 20},
	}
	if path != "" {
		if err := applyFile(&cfg, path); err != nil {
			return Config{}, err
		}
	}
	if err := applyEnv(&cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func applyFile(cfg *Config, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ler configuração %s: %w", path, err)
	}
	var file fileConfig
	if err := yaml.Unmarshal(raw, &file); err != nil {
		return fmt.Errorf("interpretar configuração %s: %w", path, err)
	}
	if file.Server.Address != "" {
		cfg.Server.Address = file.Server.Address
	}
	if file.Server.PublicBaseURL != "" {
		cfg.Server.PublicBaseURL = file.Server.PublicBaseURL
	}
	if file.Server.ReadTimeout != "" {
		d, err := time.ParseDuration(file.Server.ReadTimeout)
		if err != nil {
			return fmt.Errorf("server.read_timeout: %w", err)
		}
		cfg.Server.ReadTimeout = d
	}
	if file.Server.WriteTimeout != "" {
		d, err := time.ParseDuration(file.Server.WriteTimeout)
		if err != nil {
			return fmt.Errorf("server.write_timeout: %w", err)
		}
		cfg.Server.WriteTimeout = d
	}
	if file.Auth.SQLitePath != "" {
		cfg.Auth.SQLitePath = file.Auth.SQLitePath
	}
	if file.Auth.AccessTokenTTL != "" {
		d, err := time.ParseDuration(file.Auth.AccessTokenTTL)
		if err != nil {
			return fmt.Errorf("auth.access_token_ttl: %w", err)
		}
		cfg.Auth.AccessTokenTTL = d
	}
	if file.Auth.RefreshTokenTTL != "" {
		d, err := time.ParseDuration(file.Auth.RefreshTokenTTL)
		if err != nil {
			return fmt.Errorf("auth.refresh_token_ttl: %w", err)
		}
		cfg.Auth.RefreshTokenTTL = d
	}
	return applyOracleFile(cfg, file)
}

func applyEnv(cfg *Config) error {
	if v := os.Getenv("ORACLE_MCP_ADDRESS"); v != "" {
		cfg.Server.Address = v
	}
	if v := os.Getenv("ORACLE_MCP_PUBLIC_BASE_URL"); v != "" {
		cfg.Server.PublicBaseURL = v
	}
	if v := os.Getenv("ORACLE_MCP_READ_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_READ_TIMEOUT: %w", err)
		}
		cfg.Server.ReadTimeout = d
	}
	if v := os.Getenv("ORACLE_MCP_WRITE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_WRITE_TIMEOUT: %w", err)
		}
		cfg.Server.WriteTimeout = d
	}
	if v := os.Getenv("ORACLE_MCP_SQLITE_PATH"); v != "" {
		cfg.Auth.SQLitePath = v
	}
	if v := os.Getenv("ORACLE_MCP_ACCESS_TOKEN_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_ACCESS_TOKEN_TTL: %w", err)
		}
		cfg.Auth.AccessTokenTTL = d
	}
	if v := os.Getenv("ORACLE_MCP_REFRESH_TOKEN_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("ORACLE_MCP_REFRESH_TOKEN_TTL: %w", err)
		}
		cfg.Auth.RefreshTokenTTL = d
	}
	return applyOracleEnv(cfg)
}

func (cfg *Config) validate() error {
	if _, _, err := net.SplitHostPort(cfg.Server.Address); err != nil {
		return fmt.Errorf("server.address: %w", err)
	}
	if cfg.Server.ReadTimeout <= 0 {
		return fmt.Errorf("server.read_timeout deve ser positivo")
	}
	if cfg.Server.WriteTimeout <= 0 {
		return fmt.Errorf("server.write_timeout deve ser positivo")
	}
	if cfg.Auth.AccessTokenTTL <= 0 || cfg.Auth.RefreshTokenTTL <= 0 {
		return fmt.Errorf("a vida dos tokens deve ser positiva")
	}
	if strings.TrimSpace(cfg.Auth.SQLitePath) == "" {
		return fmt.Errorf("auth.sqlite_path é obrigatório")
	}
	if cfg.Server.PublicBaseURL == "" {
		host, _, err := net.SplitHostPort(cfg.Server.Address)
		if err != nil {
			return err
		}
		if host == "0.0.0.0" || host == "::" || host == "" {
			return fmt.Errorf("server.public_base_url é obrigatório quando o bind não é um host concreto")
		}
		cfg.Server.PublicBaseURL = "http://" + cfg.Server.Address
	}
	u, err := url.Parse(cfg.Server.PublicBaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("server.public_base_url deve ser uma URL http ou https absoluta")
	}
	cfg.Server.PublicBaseURL = strings.TrimRight(cfg.Server.PublicBaseURL, "/")
	return cfg.validateOracle()
}

// Issuer é a URL pública sem barra final.
func (cfg Config) Issuer() string {
	return strings.TrimRight(cfg.Server.PublicBaseURL, "/")
}

// Resource é a audiência do endpoint MCP.
func (cfg Config) Resource() string {
	return cfg.Issuer() + "/mcp"
}
