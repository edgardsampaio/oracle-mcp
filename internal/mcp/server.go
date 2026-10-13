// Package mcp expõe o transporte Streamable HTTP do servidor.
package mcp

import (
	"log/slog"
	"net/http"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"oracle-mcp/internal/bi"
)

const (
	serverName    = "oracle-mcp"
	serverVersion = "0.1.0"
	// sessionTimeout encerra sessão ociosa. Zero deixaria a sessão para sempre.
	sessionTimeout = 30 * time.Minute
)

// Handler devolve o endpoint /mcp com sessão stateful.
// Stateless fica desligado: o protocolo 2026-07-28 sobre HTTP exige esse modo,
// e clientes anteriores negociam para baixo a partir do handshake.
func Handler(logger *slog.Logger, db bi.Querier, audit bi.Auditor, welcomeDir string, emails bi.EmailSource, maxRows int) http.Handler {
	server := sdk.NewServer(&sdk.Implementation{
		Name:    serverName,
		Title:   "Oracle MCP",
		Version: serverVersion,
	}, &sdk.ServerOptions{Logger: logger})
	registerTools(server, bi.ToolEnv{DB: db, Audit: audit, Log: logger, MaxRows: maxRows, Emails: emails}, welcomeDir)
	return sdk.NewStreamableHTTPHandler(func(*http.Request) *sdk.Server {
		return server
	}, &sdk.StreamableHTTPOptions{
		Logger:         logger,
		SessionTimeout: sessionTimeout,
	})
}
