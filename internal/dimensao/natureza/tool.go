// Package natureza é a dimensão DIM_NATUREZA.
// Não importa fato nem outra dimensão.
package natureza

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_natureza"

	description = `Dimensão de naturezas financeiras em hierarquia de até quatro níveis (N1 a N4) mais o nível analítico.

Classificar receitas e despesas: permite à IA traduzir termos como 'despesas com pessoal' ou 'impostos' para naturezas e agregar o financeiro em qualquer nível.

Filtros opcionais: id_natureza, codigo_natureza, nome_natureza, id_natureza_n1, nome_natureza_n1, id_natureza_n2, nome_natureza_n2, id_natureza_n3, id_natureza_n4, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Natureza N1",
    "Código Natureza N1",
    "Nome Natureza N1",
    "Natureza N1",
    "ID Natureza N2",
    "Código Natureza N2",
    "Nome Natureza N2",
    "Natureza N2",
    "ID Natureza N3",
    "Código Natureza N3",
    "Nome Natureza N3",
    "Natureza N3",
    "ID Natureza N4",
    "Código Natureza N4",
    "Nome Natureza N4",
    "Natureza N4",
    "ID Natureza",
    "Código Natureza",
    "Nome Natureza",
    "Natureza",
    "Registro Bloqueado Natureza",
    "Registro Ativo Natureza"
FROM DIM_NATUREZA
WHERE (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:codigo_natureza IS NULL OR "Código Natureza" = :codigo_natureza)
  AND (:nome_natureza IS NULL OR UPPER("Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:id_natureza_n1 IS NULL OR "ID Natureza N1" = :id_natureza_n1)
  AND (:nome_natureza_n1 IS NULL OR UPPER("Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:id_natureza_n2 IS NULL OR "ID Natureza N2" = :id_natureza_n2)
  AND (:nome_natureza_n2 IS NULL OR UPPER("Nome Natureza N2") LIKE '%' || UPPER(:nome_natureza_n2) || '%')
  AND (:id_natureza_n3 IS NULL OR "ID Natureza N3" = :id_natureza_n3)
  AND (:id_natureza_n4 IS NULL OR "ID Natureza N4" = :id_natureza_n4)
  AND (:registro_ativo IS NULL OR "Registro Ativo Natureza" = :registro_ativo)
ORDER BY "Código Natureza"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDNatureza     *int    `json:"id_natureza,omitempty"`
	CodigoNatureza *string `json:"codigo_natureza,omitempty"`
	NomeNatureza   *string `json:"nome_natureza,omitempty"`
	IDNaturezaN1   *int    `json:"id_natureza_n1,omitempty"`
	NomeNaturezaN1 *string `json:"nome_natureza_n1,omitempty"`
	IDNaturezaN2   *int    `json:"id_natureza_n2,omitempty"`
	NomeNaturezaN2 *string `json:"nome_natureza_n2,omitempty"`
	IDNaturezaN3   *int    `json:"id_natureza_n3,omitempty"`
	IDNaturezaN4   *int    `json:"id_natureza_n4,omitempty"`
	RegistroAtivo  *string `json:"registro_ativo,omitempty"`
	Limite         *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
		sql.Named("codigo_natureza", bi.NullString(in.CodigoNatureza)),
		sql.Named("nome_natureza", bi.NullString(in.NomeNatureza)),
		sql.Named("id_natureza_n1", bi.NullInt(in.IDNaturezaN1)),
		sql.Named("nome_natureza_n1", bi.NullString(in.NomeNaturezaN1)),
		sql.Named("id_natureza_n2", bi.NullInt(in.IDNaturezaN2)),
		sql.Named("nome_natureza_n2", bi.NullString(in.NomeNaturezaN2)),
		sql.Named("id_natureza_n3", bi.NullInt(in.IDNaturezaN3)),
		sql.Named("id_natureza_n4", bi.NullInt(in.IDNaturezaN4)),
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
