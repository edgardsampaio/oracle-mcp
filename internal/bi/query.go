package bi

import (
	"context"
	"log/slog"
)

// ToolEnv é o que cada tool recebe para consultar e auditar.
type ToolEnv struct {
	DB      Querier
	Audit   Auditor
	Log     *slog.Logger
	MaxRows int
	Emails  EmailSource
}

// QueryRows executa o SQL já montado no pacote da tool.
func QueryRows(ctx context.Context, env ToolEnv, tool, statement string, args []any) ([]map[string]any, error) {
	if env.DB == nil {
		return nil, ErrUnavailable
	}
	res, err := env.DB.Query(ctx, tool, statement, args...)
	if err != nil {
		return nil, ErrUnavailable
	}
	if res.Rows == nil {
		return []map[string]any{}, nil
	}
	return res.Rows, nil
}
