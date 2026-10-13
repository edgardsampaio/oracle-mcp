// Package api monta as rotas HTTP do processo.
package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/config"
)

// Pinger verifica uma dependência sem expor o erro ao cliente.
type Pinger interface {
	Ping(context.Context) error
}

// NewHandler devolve o roteador com health e o endpoint MCP.
// ping nulo deixa o readiness só com o processo. Com Oracle configurado, o ping entra no ready.
func NewHandler(logger *slog.Logger, cfg config.Config, svc *auth.Service, ping Pinger, db bi.Querier, audit bi.Auditor, welcomeDir string) http.Handler {
	mux := http.NewServeMux()
	api := &httpAPI{svc: svc, cfg: cfg, welcomeDir: welcomeDir}
	mux.HandleFunc("GET /api/v1/health/live", health("live"))
	mux.HandleFunc("GET /api/v1/health/ready", ready(ping))
	mux.HandleFunc("POST /api/v1/auth/login", api.login)
	mux.HandleFunc("POST /api/v1/auth/refresh", api.refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", api.logout)
	mux.HandleFunc("GET /api/v1/users", api.listUsers)
	mux.HandleFunc("POST /api/v1/users", api.createUser)
	mux.HandleFunc("GET /api/v1/users/{id}", api.getUser)
	mux.HandleFunc("PATCH /api/v1/users/{id}", api.patchUser)
	mux.HandleFunc("DELETE /api/v1/users/{id}", api.deleteUser)
	mux.HandleFunc("PUT /api/v1/users/{id}/password", api.setPassword)
	mux.HandleFunc("PUT /api/v1/users/{id}/roles", api.setRoles)
	mux.HandleFunc("GET /api/v1/roles", api.listRoles)
	mux.Handle("/.well-known/oauth-protected-resource", resourceMetadata(cfg))
	mux.Handle("/.well-known/oauth-protected-resource/mcp", resourceMetadata(cfg))
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", api.asMetadata)
	mux.HandleFunc("GET /oauth/authorize", api.authorize)
	mux.HandleFunc("POST /oauth/authorize", api.authorize)
	mux.HandleFunc("POST /oauth/token", api.token)
	mux.HandleFunc("POST /oauth/register", api.register)
	mux.Handle("/mcp", protectedMCP(logger, cfg, svc, db, audit, welcomeDir))
	mux.Handle("GET /swagger", swaggerHandler())
	mux.Handle("GET /swagger/", swaggerHandler())
	return logRequests(logger, mux)
}

func health(status string) http.HandlerFunc {
	body := []byte(`{"status":"` + status + `"}` + "\n")
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

func ready(ping Pinger) http.HandlerFunc {
	ok := health("ready")
	return func(w http.ResponseWriter, r *http.Request) {
		if ping != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			if err := ping.Ping(ctx); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("{\"status\":\"unavailable\"}\n"))
				return
			}
		}
		ok(w, r)
	}
}

func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"duration", time.Since(start),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
