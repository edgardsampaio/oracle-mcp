// Package cliente é a dimensão DIM_CLIENTE.
// Não importa fato nem outra dimensão.
package cliente

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_cliente"

	description = `Dimensão de clientes com dados cadastrais (razão social, nome fantasia, CNPJ/CPF, inscrições, endereço, município, UF e região).

Localizar clientes por nome, documento ou localização e obter o ID usado em Contas a Receber e movimentações bancárias.

Filtros opcionais: id_cliente, codigo_cliente, loja_cliente, razao_social, nome_fantasia, cnpj, raiz_cnpj, sigla_uf, nome_municipio, regiao, codigo_tipo_pessoa, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Cliente",
    "Código Cliente",
    "Loja Cliente",
    "Razão Social",
    "Nome Fantasia",
    "Cliente Razão Social",
    "Cliente Nome Fantasia",
    "Cliente CNPJ Razão Social",
    "Cliente CNPJ Nome Fantasia",
    "Raíz CNPJ",
    "CNPJ",
    "Inscrição Estadual",
    "Inscrição Municipal",
    "Endereço",
    "Bairro",
    "CEP",
    "Código Tipo Pessoa",
    "Descrição Tipo Pessoa",
    "Tipo Pessoa",
    "ID Município",
    "Código Município",
    "Nome Município",
    "Município",
    "ID UF",
    "Sigla UF",
    "Nome UF",
    "UF",
    "Região",
    "Sigla País",
    "Registro Bloqueado Cliente",
    "Registro Ativo Cliente",
    "Registro Bloqueado Pessoa",
    "Registro Ativo Pessoa",
    "Registro Ativo Município",
    "Registro Ativo UF"
FROM DIM_CLIENTE
WHERE (:id_cliente IS NULL OR "ID Cliente" = :id_cliente)
  AND (:codigo_cliente IS NULL OR "Código Cliente" = :codigo_cliente)
  AND (:loja_cliente IS NULL OR "Loja Cliente" = :loja_cliente)
  AND (:razao_social IS NULL OR UPPER("Razão Social") LIKE '%' || UPPER(:razao_social) || '%')
  AND (:nome_fantasia IS NULL OR UPPER("Nome Fantasia") LIKE '%' || UPPER(:nome_fantasia) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:sigla_uf IS NULL OR "Sigla UF" = :sigla_uf)
  AND (:nome_municipio IS NULL OR UPPER("Nome Município") LIKE '%' || UPPER(:nome_municipio) || '%')
  AND (:regiao IS NULL OR "Região" = :regiao)
  AND (:codigo_tipo_pessoa IS NULL OR "Código Tipo Pessoa" = :codigo_tipo_pessoa)
  AND (:registro_ativo IS NULL OR "Registro Ativo Cliente" = :registro_ativo)
ORDER BY "Razão Social"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDCliente        *int    `json:"id_cliente,omitempty"`
	CodigoCliente    *string `json:"codigo_cliente,omitempty"`
	LojaCliente      *string `json:"loja_cliente,omitempty"`
	RazaoSocial      *string `json:"razao_social,omitempty"`
	NomeFantasia     *string `json:"nome_fantasia,omitempty"`
	CNPJ             *string `json:"cnpj,omitempty"`
	RaizCNPJ         *string `json:"raiz_cnpj,omitempty"`
	SiglaUF          *string `json:"sigla_uf,omitempty"`
	NomeMunicipio    *string `json:"nome_municipio,omitempty"`
	Regiao           *string `json:"regiao,omitempty"`
	CodigoTipoPessoa *string `json:"codigo_tipo_pessoa,omitempty"`
	RegistroAtivo    *string `json:"registro_ativo,omitempty"`
	Limite           *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_cliente", bi.NullInt(in.IDCliente)),
		sql.Named("codigo_cliente", bi.NullString(in.CodigoCliente)),
		sql.Named("loja_cliente", bi.NullString(in.LojaCliente)),
		sql.Named("razao_social", bi.NullString(in.RazaoSocial)),
		sql.Named("nome_fantasia", bi.NullString(in.NomeFantasia)),
		sql.Named("cnpj", bi.NullString(in.CNPJ)),
		sql.Named("raiz_cnpj", bi.NullString(in.RaizCNPJ)),
		sql.Named("sigla_uf", bi.NullString(in.SiglaUF)),
		sql.Named("nome_municipio", bi.NullString(in.NomeMunicipio)),
		sql.Named("regiao", bi.NullString(in.Regiao)),
		sql.Named("codigo_tipo_pessoa", bi.NullString(in.CodigoTipoPessoa)),
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
