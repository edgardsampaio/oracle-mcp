// Package empresafilial é a dimensão DIM_EMPRESA_FILIAL.
// Não importa fato nem outra dimensão.
package empresafilial

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_empresa_filial"

	description = `Dimensão de empresas e filiais do grupo, com códigos, nomes e CNPJ.

Resolver empresas e filiais citadas pelo usuário para os IDs usados em todos os fatos e permitir consolidação por empresa.

Filtros opcionais: id_empresa, codigo_empresa, nome_empresa, id_filial, codigo_filial, nome_filial, cnpj, raiz_cnpj, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Empresa",
    "Código Empresa",
    "Nome Empresa",
    "Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Raíz CNPJ",
    "CNPJ",
    "Registro Ativo Empresa",
    "Registro Ativo Filial"
FROM DIM_EMPRESA_FILIAL
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:codigo_empresa IS NULL OR "Código Empresa" = :codigo_empresa)
  AND (:nome_empresa IS NULL OR UPPER("Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:codigo_filial IS NULL OR "Código Filial" = :codigo_filial)
  AND (:nome_filial IS NULL OR UPPER("Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:registro_ativo IS NULL OR "Registro Ativo Filial" = :registro_ativo)
ORDER BY "Nome Empresa", "Nome Filial"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDEmpresa     *int    `json:"id_empresa,omitempty"`
	CodigoEmpresa *string `json:"codigo_empresa,omitempty"`
	NomeEmpresa   *string `json:"nome_empresa,omitempty"`
	IDFilial      *int    `json:"id_filial,omitempty"`
	CodigoFilial  *string `json:"codigo_filial,omitempty"`
	NomeFilial    *string `json:"nome_filial,omitempty"`
	CNPJ          *string `json:"cnpj,omitempty"`
	RaizCNPJ      *string `json:"raiz_cnpj,omitempty"`
	RegistroAtivo *string `json:"registro_ativo,omitempty"`
	Limite        *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_empresa", bi.NullInt(in.IDEmpresa)),
		sql.Named("codigo_empresa", bi.NullString(in.CodigoEmpresa)),
		sql.Named("nome_empresa", bi.NullString(in.NomeEmpresa)),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("codigo_filial", bi.NullString(in.CodigoFilial)),
		sql.Named("nome_filial", bi.NullString(in.NomeFilial)),
		sql.Named("cnpj", bi.NullString(in.CNPJ)),
		sql.Named("raiz_cnpj", bi.NullString(in.RaizCNPJ)),
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
