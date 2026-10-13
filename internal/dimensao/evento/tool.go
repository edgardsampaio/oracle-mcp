// Package evento é a dimensão DIM_EVENTO.
// Não importa fato nem outra dimensão.
package evento

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_evento"

	description = `Dimensão de eventos da folha de pagamento (proventos, descontos, bases) e seus tipos.

Identificar eventos de folha (ex.: salário, horas extras, INSS) para filtrar e agrupar a FATO_FOLHA_PAGAMENTO.

Filtros opcionais: id_evento, codigo_evento, nome_evento, id_tipo_evento, nome_tipo_evento, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Evento",
    "Código Evento",
    "Nome Evento",
    "Evento",
    "ID Tipo Evento",
    "Código Tipo Evento",
    "Nome Tipo Evento",
    "Tipo Evento",
    "Registro Ativo Evento",
    "Registro Ativo Tipo Evento"
FROM DIM_EVENTO
WHERE (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:codigo_evento IS NULL OR "Código Evento" = :codigo_evento)
  AND (:nome_evento IS NULL OR UPPER("Nome Evento") LIKE '%' || UPPER(:nome_evento) || '%')
  AND (:id_tipo_evento IS NULL OR "ID Tipo Evento" = :id_tipo_evento)
  AND (:nome_tipo_evento IS NULL OR UPPER("Nome Tipo Evento") LIKE '%' || UPPER(:nome_tipo_evento) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Evento" = :registro_ativo)
ORDER BY "Nome Tipo Evento", "Nome Evento"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDEvento       *int    `json:"id_evento,omitempty"`
	CodigoEvento   *string `json:"codigo_evento,omitempty"`
	NomeEvento     *string `json:"nome_evento,omitempty"`
	IDTipoEvento   *int    `json:"id_tipo_evento,omitempty"`
	NomeTipoEvento *string `json:"nome_tipo_evento,omitempty"`
	RegistroAtivo  *string `json:"registro_ativo,omitempty"`
	Limite         *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_evento", bi.NullInt(in.IDEvento)),
		sql.Named("codigo_evento", bi.NullString(in.CodigoEvento)),
		sql.Named("nome_evento", bi.NullString(in.NomeEvento)),
		sql.Named("id_tipo_evento", bi.NullInt(in.IDTipoEvento)),
		sql.Named("nome_tipo_evento", bi.NullString(in.NomeTipoEvento)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type saida struct {
	Itens []map[string]any `json:"itens"`
}

// Register publica a tool no servidor MCP.
func Register(server *sdk.Server, env bi.ToolEnv) {
	sdk.AddTool(server, &sdk.Tool{Name: ToolName, Description: description}, func(ctx context.Context, req *sdk.CallToolRequest, in entrada) (*sdk.CallToolResult, saida, error) {
		rows, err := Search(ctx, env, userID(req), in)
		if rows == nil {
			rows = []map[string]any{}
		}
		return nil, saida{Itens: rows}, err
	})
}

// Search executa a consulta cadastrada desta dimensão.
func Search(ctx context.Context, env bi.ToolEnv, userID string, in entrada) (rows []map[string]any, err error) {
	start := time.Now()
	defer bi.Finish(ctx, env.Audit, env.Log, userID, ToolName, start, &err)
	if userID == "" {
		return nil, bi.ErrIdentity
	}
	args, err := in.binds(env.MaxRows)
	if err != nil {
		return nil, err
	}
	return bi.QueryRows(ctx, env, ToolName, searchSQL, args)
}

func userID(req *sdk.CallToolRequest) string {
	if req == nil || req.Extra == nil || req.Extra.TokenInfo == nil {
		return ""
	}
	return req.Extra.TokenInfo.UserID
}
