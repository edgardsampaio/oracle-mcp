package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"oracle-mcp/internal/api"
	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/config"
	"oracle-mcp/internal/database/oracle"
	"oracle-mcp/internal/identity/sqlite"
)

var errUsage = errors.New("use: oracle-mcp [admin create]")

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger, os.Args[1:]); err != nil {
		logger.Error("falha", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, args []string) error {
	if len(args) > 0 && args[0] == "admin" {
		return runAdmin(logger, args[1:])
	}
	if len(args) > 0 {
		return errUsage
	}
	if err := prepareRuntime(); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := applyRuntimePaths(&cfg); err != nil {
		return err
	}
	store, err := sqlite.Open(cfg.Auth.SQLitePath)
	if err != nil {
		return err
	}
	defer store.Close()
	svc, err := auth.New(store, cfg, logger)
	if err != nil {
		return err
	}
	admins, err := store.CountActiveAdmins(context.Background())
	if err != nil {
		return err
	}
	username, err := ensureInitialAdmin(context.Background(), admins, os.Stderr, func() (string, string, string, error) {
		return promptAdmin(os.Stdin, os.Stderr)
	}, svc.CreateInitialAdmin)
	if err != nil {
		return err
	}
	if username != "" {
		logger.Info("administrador criado", "username", username)
	}
	var ping api.Pinger
	var querier bi.Querier
	if cfg.Oracle.Enabled() {
		odb, err := oracle.Open(cfg.Oracle, cfg.BI, logger)
		if err != nil {
			return err
		}
		defer odb.Close()
		ping = odb
		querier = odb
	}
	exeDir, err := executableDir()
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:              cfg.Server.Address,
		Handler:           api.NewHandler(logger, cfg, svc, ping, querier, bi.StoreAuditor(store), filepath.Join(exeDir, "config")),
		ReadHeaderTimeout: cfg.Server.ReadTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
	}
	logger.Info("escutando", "addr", cfg.Server.Address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
