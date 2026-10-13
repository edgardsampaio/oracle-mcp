// Package tipotitulo é a dimensão DIM_TIPO_TITULO.
// Não importa fato nem outra dimensão.
package tipotitulo

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_tipo_titulo"

	description = `Dimensão de tipos de título financeiro (ex.: NF, boleto, adiantamento, impostos).

Classificar títulos de Contas a Pagar, Contas a Receber e baixas por tipo — por exemplo, separar adiantamentos e impostos de títulos comerciais.

Filtros opcionais: id, codigo, descricao, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Tipo Título",
    "Código Tipo Título",
    "Descrição Tipo Título",
    "Tipo Título",
    "Registro Ativo Tipo Título"
FROM DIM_TIPO_TITULO
WHERE (:id IS NULL OR "ID Tipo Título" = :id)
  AND (:codigo IS NULL OR "Código Tipo Título" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Título") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Tipo Título" = :registro_ativo)
ORDER BY "Código Tipo Título"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	ID            *int    `json:"id,omitempty"`
	Codigo        *string `json:"codigo,omitempty"`
	Descricao     *string `json:"descricao,omitempty"`
	RegistroAtivo *string `json:"registro_ativo,omitempty"`
	Limite        *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id", bi.NullInt(in.ID)),
		sql.Named("codigo", bi.NullString(in.Codigo)),
		sql.Named("descricao", bi.NullString(in.Descricao)),
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
