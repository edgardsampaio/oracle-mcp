// Package rlscentrocusto é a dimensão DIM_RLS_CENTRO_CUSTO.
// Não importa fato nem outra dimensão.
package rlscentrocusto

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_rls_centro_custo"

	description = `Tabela de segurança em nível de linha (RLS): relaciona o e-mail do usuário aos contratos e centros de custo que ele pode visualizar.

Restringir os dados retornados por qualquer consulta aos contratos/centros de custo autorizados ao usuário. O parâmetro 'email' deve ser injetado pelo servidor MCP a partir do usuário autenticado — nunca informado pela IA ou pelo usuário no chat.

Devolve só os centros de custo autorizados ao usuário autenticado. O e-mail não faz parte dos parâmetros. Filtros opcionais: id_contrato, id_centro_custo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "E-mail",
    "ID Contrato",
    "Código Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo"
FROM DIM_RLS_CENTRO_CUSTO
WHERE "E-mail" = :email
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
ORDER BY "Código Contrato", "Código Centro de Custo"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDContrato    *int `json:"id_contrato,omitempty"`
	IDCentroCusto *int `json:"id_centro_custo,omitempty"`
	Limite        *int `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
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
	email, err := bi.RequireEmail(ctx, env.Emails, userID)
	if err != nil {
		return nil, err
	}
	args, err := in.binds(env.MaxRows)
	if err != nil {
		return nil, err
	}
	args = append([]any{sql.Named("email", email)}, args...)
	return bi.QueryRows(ctx, env, ToolName, searchSQL, args)
}

func userID(req *sdk.CallToolRequest) string {
	if req == nil || req.Extra == nil || req.Extra.TokenInfo == nil {
		return ""
	}
	return req.Extra.TokenInfo.UserID
}
