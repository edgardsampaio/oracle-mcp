// Package contaorcmnatureza é a dimensão DIM_CONTA_ORCM_NATUREZA.
// Não importa fato nem outra dimensão.
package contaorcmnatureza

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_conta_orcm_natureza"

	description = `Tabela ponte (de-para) entre Contas Orçamentárias e Naturezas Financeiras.

Permitir cruzar orçado (por conta orçamentária) com realizado (por natureza financeira em Contas a Pagar/Receber e movimentações bancárias).

Filtros opcionais: id_conta_orcamentaria, id_natureza, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Conta Orçamentária",
    "ID Natureza"
FROM DIM_CONTA_ORCM_NATUREZA
WHERE (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
ORDER BY "ID Conta Orçamentária", "ID Natureza"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDContaOrcamentaria *int `json:"id_conta_orcamentaria,omitempty"`
	IDNatureza          *int `json:"id_natureza,omitempty"`
	Limite              *int `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_conta_orcamentaria", bi.NullInt(in.IDContaOrcamentaria)),
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
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
