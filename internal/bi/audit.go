package bi

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"oracle-mcp/internal/database/oracle"
	"oracle-mcp/internal/identity/sqlite"
)

// Querier executa a consulta que a tool já declarou.
type Querier interface {
	Query(ctx context.Context, id, statement string, args ...any) (oracle.Result, error)
}

// Auditor grava ferramenta, ator, duração e outcome. Sem SQL e sem linhas.
type Auditor interface {
	Record(ctx context.Context, actor, tool, outcome string, duration time.Duration) error
}

// Finish grava o outcome da tool. Sem SQL e sem linhas.
func Finish(ctx context.Context, audit Auditor, log *slog.Logger, actor, tool string, start time.Time, errp *error) {
	if audit == nil {
		return
	}
	outcome := "success"
	if errp != nil && *errp != nil {
		outcome = "failure"
	}
	if err := audit.Record(ctx, actor, tool, outcome, time.Since(start)); err != nil && log != nil {
		log.Info("auditoria", "tool", tool, "outcome", "failure")
	}
}

type storeAuditor struct {
	store *sqlite.Store
}

// StoreAuditor grava em audit_log. O metadata traz só a duração.
func StoreAuditor(store *sqlite.Store) Auditor {
	return storeAuditor{store: store}
}

func (a storeAuditor) Record(ctx context.Context, actor, tool, outcome string, duration time.Duration) error {
	if a.store == nil {
		return nil
	}
	return a.store.InsertAudit(ctx, sqlite.AuditEvent{
		ActorID:  actor,
		Action:   "tool",
		Target:   tool,
		Outcome:  outcome,
		Metadata: map[string]string{"duration_ms": strconv.FormatInt(duration.Milliseconds(), 10)},
	}, time.Now())
}
