// Package contaorcamentaria é a dimensão DIM_CONTA_ORCAMENTARIA.
// Não importa fato nem outra dimensão.
package contaorcamentaria

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_conta_orcamentaria"

	description = `Dimensão do plano de contas orçamentárias em três níveis (N1, N2, N3), com nível, tipo, condição e natureza associada.

Traduzir contas orçamentárias citadas pelo usuário para IDs e permitir análises de orçado por qualquer nível da hierarquia.

Filtros opcionais: id_conta_orcamentaria, codigo_conta, nome_conta, codigo_n1, codigo_n2, codigo_n3, nivel, tipo, codigo_natureza, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Conta Orçamentária",
    "Código Conta Orçamentária",
    "Nome Conta Orçamentária",
    "Conta Orçamentária",
    "Código Conta Orçamentária N1",
    "Nome Conta Orçamentária N1",
    "Conta Orçamentária N1",
    "Código Conta Orçamentária N2",
    "Nome Conta Orçamentária N2",
    "Conta Orçamentária N2",
    "Código Conta Orçamentária N3",
    "Nome Conta Orçamentária N3",
    "Conta Orçamentária N3",
    "Nível",
    "Tipo",
    "Condição",
    "Código Natureza",
    "Reg. Ativo Conta Orçamentária"
FROM DIM_CONTA_ORCAMENTARIA
WHERE (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:codigo_conta IS NULL OR "Código Conta Orçamentária" = :codigo_conta)
  AND (:nome_conta IS NULL OR UPPER("Nome Conta Orçamentária") LIKE '%' || UPPER(:nome_conta) || '%')
  AND (:codigo_n1 IS NULL OR "Código Conta Orçamentária N1" = :codigo_n1)
  AND (:codigo_n2 IS NULL OR "Código Conta Orçamentária N2" = :codigo_n2)
  AND (:codigo_n3 IS NULL OR "Código Conta Orçamentária N3" = :codigo_n3)
  AND (:nivel IS NULL OR "Nível" = :nivel)
  AND (:tipo IS NULL OR "Tipo" = :tipo)
  AND (:codigo_natureza IS NULL OR "Código Natureza" = :codigo_natureza)
  AND (:registro_ativo IS NULL OR "Reg. Ativo Conta Orçamentária" = :registro_ativo)
ORDER BY "Código Conta Orçamentária"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDContaOrcamentaria *int    `json:"id_conta_orcamentaria,omitempty"`
	CodigoConta         *string `json:"codigo_conta,omitempty"`
	NomeConta           *string `json:"nome_conta,omitempty"`
	CodigoN1            *string `json:"codigo_n1,omitempty"`
	CodigoN2            *string `json:"codigo_n2,omitempty"`
	CodigoN3            *string `json:"codigo_n3,omitempty"`
	Nivel               *string `json:"nivel,omitempty"`
	Tipo                *string `json:"tipo,omitempty"`
	CodigoNatureza      *string `json:"codigo_natureza,omitempty"`
	RegistroAtivo       *string `json:"registro_ativo,omitempty"`
	Limite              *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_conta_orcamentaria", bi.NullInt(in.IDContaOrcamentaria)),
		sql.Named("codigo_conta", bi.NullString(in.CodigoConta)),
		sql.Named("nome_conta", bi.NullString(in.NomeConta)),
		sql.Named("codigo_n1", bi.NullString(in.CodigoN1)),
		sql.Named("codigo_n2", bi.NullString(in.CodigoN2)),
		sql.Named("codigo_n3", bi.NullString(in.CodigoN3)),
		sql.Named("nivel", bi.NullString(in.Nivel)),
		sql.Named("tipo", bi.NullString(in.Tipo)),
		sql.Named("codigo_natureza", bi.NullString(in.CodigoNatureza)),
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
