// Package bancocontabancaria é a dimensão DIM_BANCO_CONTA_BANCARIA.
// Não importa fato nem outra dimensão.
package bancocontabancaria

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_banco_conta_bancaria"

	description = `Dimensão de bancos e contas bancárias, com agência, número, tipo de conta e a empresa/filial à qual cada conta pertence.

Permitir que a IA identifique bancos e contas bancárias (por nome, código ou filial) e obtenha os IDs usados para filtrar e rotular as movimentações bancárias e baixas de títulos.

Filtros opcionais: id_empresa, id_filial, id_banco, codigo_banco, nome_banco, id_conta_bancaria, nome_conta, tipo_conta_bancaria, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Banco",
    "Código Banco",
    "Nome Banco",
    "Banco",
    "ID Conta Bancária",
    "Código Agência",
    "Número Conta Bancária",
    "Nome Conta Bancária",
    "Nome Reduzido Conta Bancária",
    "Tipo Conta Bancária",
    "Conta Bancária",
    "Banco, Agência e Conta Bancária",
    "ID Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Registro Ativo Banco",
    "Reg. Bloqueado Conta Bancária",
    "Registro Ativo Conta Bancária"
FROM DIM_BANCO_CONTA_BANCARIA
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_banco IS NULL OR "ID Banco" = :id_banco)
  AND (:codigo_banco IS NULL OR "Código Banco" = :codigo_banco)
  AND (:nome_banco IS NULL OR UPPER("Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:nome_conta IS NULL OR UPPER("Nome Conta Bancária") LIKE '%' || UPPER(:nome_conta) || '%')
  AND (:tipo_conta_bancaria IS NULL OR "Tipo Conta Bancária" = :tipo_conta_bancaria)
  AND (:registro_ativo IS NULL OR "Registro Ativo Conta Bancária" = :registro_ativo)
ORDER BY "Nome Banco", "Nome Conta Bancária"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDEmpresa         *int    `json:"id_empresa,omitempty"`
	IDFilial          *int    `json:"id_filial,omitempty"`
	IDBanco           *int    `json:"id_banco,omitempty"`
	CodigoBanco       *string `json:"codigo_banco,omitempty"`
	NomeBanco         *string `json:"nome_banco,omitempty"`
	IDContaBancaria   *int    `json:"id_conta_bancaria,omitempty"`
	NomeConta         *string `json:"nome_conta,omitempty"`
	TipoContaBancaria *string `json:"tipo_conta_bancaria,omitempty"`
	RegistroAtivo     *string `json:"registro_ativo,omitempty"`
	Limite            *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_empresa", bi.NullInt(in.IDEmpresa)),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_banco", bi.NullInt(in.IDBanco)),
		sql.Named("codigo_banco", bi.NullString(in.CodigoBanco)),
		sql.Named("nome_banco", bi.NullString(in.NomeBanco)),
		sql.Named("id_conta_bancaria", bi.NullInt(in.IDContaBancaria)),
		sql.Named("nome_conta", bi.NullString(in.NomeConta)),
		sql.Named("tipo_conta_bancaria", bi.NullString(in.TipoContaBancaria)),
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
