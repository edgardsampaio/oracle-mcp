// Package tempo é a dimensão DIM_TEMPO.
// Não importa fato nem outra dimensão.
package tempo

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_tempo"

	description = `Dimensão calendário com dia útil, dia da semana, semana, mês, trimestre, semestre, ano e descrições de período.

Traduzir expressões temporais ('último trimestre', 'dias úteis de março') em intervalos de datas e fornecer rótulos de período para agrupamentos.

Filtros opcionais: data_inicio, data_fim, ano, mes, trimestre, semestre, dia_util, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "Data",
    "Dia Útil",
    "Dia da Semana",
    "Nome do Dia da Semana",
    "Semana do Mês",
    "Descrição da Semana do Mês",
    "Semana do Ano",
    "Descrição da Semana do Ano",
    "Descrição Curta do Período",
    "Descrição Longa do Período",
    "Ano e Mês",
    "Mês e Ano",
    "Ano",
    "Mês",
    "Nome do Mês",
    "Trimestre",
    "Descrição do Trimestre",
    "Descrição Curta do Trimestre",
    "Descrição Longa do Trimestre",
    "Semestre",
    "Descrição do Semestre",
    "Descrição Curta do Semestre",
    "Descrição Longa do Semestre"
FROM DIM_TEMPO
WHERE (:data_inicio IS NULL OR "Data" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data" <= :data_fim)
  AND (:ano IS NULL OR "Ano" = :ano)
  AND (:mes IS NULL OR "Mês" = :mes)
  AND (:trimestre IS NULL OR "Trimestre" = :trimestre)
  AND (:semestre IS NULL OR "Semestre" = :semestre)
  AND (:dia_util IS NULL OR "Dia Útil" = :dia_util)
ORDER BY "Data"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	DataInicio *string `json:"data_inicio,omitempty"`
	DataFim    *string `json:"data_fim,omitempty"`
	Ano        *int    `json:"ano,omitempty"`
	Mes        *int    `json:"mes,omitempty"`
	Trimestre  *int    `json:"trimestre,omitempty"`
	Semestre   *int    `json:"semestre,omitempty"`
	DiaUtil    *string `json:"dia_util,omitempty"`
	Limite     *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_inicio, err := bi.NullDate(in.DataInicio)
	if err != nil {
		return nil, err
	}
	data_fim, err := bi.NullDate(in.DataFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("data_inicio", data_inicio),
		sql.Named("data_fim", data_fim),
		sql.Named("ano", bi.NullInt(in.Ano)),
		sql.Named("mes", bi.NullInt(in.Mes)),
		sql.Named("trimestre", bi.NullInt(in.Trimestre)),
		sql.Named("semestre", bi.NullInt(in.Semestre)),
		sql.Named("dia_util", bi.NullString(in.DiaUtil)),
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
