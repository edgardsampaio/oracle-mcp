// Package centrocusto é a dimensão DIM_CENTRO_CUSTO.
// Não importa fato nem outra dimensão.
package centrocusto

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_centro_custo"

	description = `Dimensão hierárquica de centros de custo: Gestor → Segmento → Grupo de Contrato → Contrato → Centro de Custo, incluindo Cliente Final e Tipo de Contrato.

Resolver nomes de gestores, segmentos, contratos e centros de custo citados em linguagem natural para os IDs usados nos fatos, e permitir agregações por qualquer nível da hierarquia.

Filtros opcionais: id_gestor, nome_gestor, id_segmento, nome_segmento, id_grupo_contrato, id_contrato, codigo_contrato, nome_contrato, id_centro_custo, codigo_centro_custo, nome_centro_custo, id_cliente_final, nome_cliente_final, id_tipo_contrato, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Gestor",
    "Código Gestor",
    "Nome Gestor",
    "Gestor",
    "ID Segmento",
    "Código Segmento",
    "Nome Segmento",
    "Segmento",
    "ID Grupo Contrato",
    "Código Grupo Contrato",
    "Nome Grupo Contrato",
    "Grupo Contrato",
    "ID Contrato",
    "Código Contrato",
    "Nome Contrato",
    "Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo",
    "Nome Centro de Custo",
    "Centro de Custo",
    "ID Cliente Final",
    "Código Cliente Final",
    "Nome Cliente Final",
    "Cliente Final",
    "ID Tipo Contrato",
    "Código Tipo Contrato",
    "Nome Tipo Contrato",
    "Tipo Contrato",
    "Registro Ativo Gestor",
    "Registro Ativo Segmento",
    "Registro Ativo Grupo Contrato",
    "Registro Ativo Contrato",
    "Registro Ativo Centro de Custo",
    "Registro Ativo Cliente Final"
FROM DIM_CENTRO_CUSTO
WHERE (:id_gestor IS NULL OR "ID Gestor" = :id_gestor)
  AND (:nome_gestor IS NULL OR UPPER("Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:id_segmento IS NULL OR "ID Segmento" = :id_segmento)
  AND (:nome_segmento IS NULL OR UPPER("Nome Segmento") LIKE '%' || UPPER(:nome_segmento) || '%')
  AND (:id_grupo_contrato IS NULL OR "ID Grupo Contrato" = :id_grupo_contrato)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:codigo_contrato IS NULL OR "Código Contrato" = :codigo_contrato)
  AND (:nome_contrato IS NULL OR UPPER("Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:codigo_centro_custo IS NULL OR "Código Centro de Custo" = :codigo_centro_custo)
  AND (:nome_centro_custo IS NULL OR UPPER("Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:id_cliente_final IS NULL OR "ID Cliente Final" = :id_cliente_final)
  AND (:nome_cliente_final IS NULL OR UPPER("Nome Cliente Final") LIKE '%' || UPPER(:nome_cliente_final) || '%')
  AND (:id_tipo_contrato IS NULL OR "ID Tipo Contrato" = :id_tipo_contrato)
  AND (:registro_ativo IS NULL OR "Registro Ativo Centro de Custo" = :registro_ativo)
ORDER BY "Nome Gestor", "Nome Contrato", "Nome Centro de Custo"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDGestor          *int    `json:"id_gestor,omitempty"`
	NomeGestor        *string `json:"nome_gestor,omitempty"`
	IDSegmento        *int    `json:"id_segmento,omitempty"`
	NomeSegmento      *string `json:"nome_segmento,omitempty"`
	IDGrupoContrato   *int    `json:"id_grupo_contrato,omitempty"`
	IDContrato        *int    `json:"id_contrato,omitempty"`
	CodigoContrato    *string `json:"codigo_contrato,omitempty"`
	NomeContrato      *string `json:"nome_contrato,omitempty"`
	IDCentroCusto     *int    `json:"id_centro_custo,omitempty"`
	CodigoCentroCusto *string `json:"codigo_centro_custo,omitempty"`
	NomeCentroCusto   *string `json:"nome_centro_custo,omitempty"`
	IDClienteFinal    *int    `json:"id_cliente_final,omitempty"`
	NomeClienteFinal  *string `json:"nome_cliente_final,omitempty"`
	IDTipoContrato    *int    `json:"id_tipo_contrato,omitempty"`
	RegistroAtivo     *string `json:"registro_ativo,omitempty"`
	Limite            *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_gestor", bi.NullInt(in.IDGestor)),
		sql.Named("nome_gestor", bi.NullString(in.NomeGestor)),
		sql.Named("id_segmento", bi.NullInt(in.IDSegmento)),
		sql.Named("nome_segmento", bi.NullString(in.NomeSegmento)),
		sql.Named("id_grupo_contrato", bi.NullInt(in.IDGrupoContrato)),
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("codigo_contrato", bi.NullString(in.CodigoContrato)),
		sql.Named("nome_contrato", bi.NullString(in.NomeContrato)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("codigo_centro_custo", bi.NullString(in.CodigoCentroCusto)),
		sql.Named("nome_centro_custo", bi.NullString(in.NomeCentroCusto)),
		sql.Named("id_cliente_final", bi.NullInt(in.IDClienteFinal)),
		sql.Named("nome_cliente_final", bi.NullString(in.NomeClienteFinal)),
		sql.Named("id_tipo_contrato", bi.NullInt(in.IDTipoContrato)),
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
