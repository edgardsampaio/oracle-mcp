// Package oracle executa consultas cadastradas em um Oracle somente leitura.
package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	go_ora "github.com/sijms/go-ora/v2"

	"oracle-mcp/internal/config"
)

// DB é o pool Oracle.
type DB struct {
	sql      *sql.DB
	timeout  time.Duration
	maxRows  int
	maxBytes int
	log      *slog.Logger
}

// Open abre o pool. A senha não entra em log.
func Open(cfg config.Oracle, limits config.BI, logger *slog.Logger) (*DB, error) {
	dsn := dataSource(cfg)
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir oracle")
	}
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &DB{sql: db, timeout: cfg.QueryTimeout, maxRows: limits.MaxRows, maxBytes: limits.MaxResponseBytes, log: logger}, nil
}

func dataSource(cfg config.Oracle) string {
	options := map[string]string{
		"TIMEOUT": strconv.Itoa(int(cfg.QueryTimeout.Seconds())),
	}
	service := cfg.ServiceName
	if cfg.SID != "" {
		service = ""
		options["SID"] = cfg.SID
	}
	return go_ora.BuildUrl(cfg.Host, cfg.Port, service, cfg.User, cfg.Password, options)
}

// Close fecha o pool.
func (db *DB) Close() error {
	if db == nil || db.sql == nil {
		return nil
	}
	return db.sql.Close()
}

// Ping verifica a conexão sem devolver detalhe ao chamador HTTP.
func (db *DB) Ping(ctx context.Context) error {
	if db == nil || db.sql == nil {
		return fmt.Errorf("oracle indisponível")
	}
	return db.sql.PingContext(ctx)
}

// Query executa uma consulta já cadastrada pelo servidor.
// O texto SQL não vem do agente e não é registrado.
func (db *DB) Query(ctx context.Context, id, statement string, args ...any) (Result, error) {
	if id == "" || statement == "" {
		return Result{}, fmt.Errorf("consulta não cadastrada")
	}
	started := time.Now()
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()
	rows, err := db.sql.QueryContext(ctx, statement, args...)
	if err != nil {
		db.logFailure(id, started)
		return Result{}, fmt.Errorf("consulta %s falhou", id)
	}
	defer rows.Close()
	result, err := Collect(rows, db.maxRows, db.maxBytes)
	if err != nil {
		db.logFailure(id, started)
		return Result{}, err
	}
	if db.log != nil {
		db.log.Info("consulta oracle", "query_id", id, "rows", len(result.Rows), "truncated", result.Truncated, "duration", time.Since(started), "outcome", "success")
	}
	return result, nil
}

func (db *DB) logFailure(id string, started time.Time) {
	if db.log == nil {
		return
	}
	db.log.Info("consulta oracle", "query_id", id, "duration", time.Since(started), "outcome", "failure")
}
