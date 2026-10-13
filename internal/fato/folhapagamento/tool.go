// Package folhapagamento é o fato FATO_FOLHA_PAGAMENTO.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package folhapagamento

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_folha_pagamento"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_folha_pagamento_com_dimensoes"

	description = `Fato de lançamentos da folha de pagamento por colaborador, competência e evento, com valor e quantidade.

Responder perguntas sobre custo de folha: total por competência, por evento (salário, horas extras, encargos) e por colaborador; combinada com DIM_COLABORADOR e FATO_HISTORICO_COLABORADOR para análises por filial, seção ou centro de custo.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: ano_mes_inicio, ano_mes_fim, data_competencia_inicio, data_competencia_fim, periodo_pagamento, id_colaborador, id_evento, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Lançamentos da folha de pagamento enriquecidos com os dados do colaborador, o evento e seu tipo, e o calendário da competência.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada via 'FATO_HISTORICO_COLABORADOR' (centro de custo do colaborador na mesma competência), já que a folha não tem 'ID Centro de Custo'.

Responder perguntas de custo de folha com rótulos legíveis — total por evento, tipo de evento, colaborador e competência — já restrito aos colaboradores lotados em centros de custo autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: ano_mes_inicio, ano_mes_fim, data_competencia_inicio, data_competencia_fim, periodo_pagamento, id_colaborador, id_evento, registro_ativo, nome_colaborador, matricula, nome_evento, nome_tipo_evento, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "Ano e Mês",
    "Data Competência",
    "Período Pagamento",
    "ID Colaborador",
    "ID Evento",
    "Valor Evento",
    "Quantidade Evento",
    "Registro Ativo"
FROM FATO_FOLHA_PAGAMENTO
WHERE EXISTS (SELECT 1 FROM FATO_HISTORICO_COLABORADOR h JOIN DIM_RLS_CENTRO_CUSTO rls ON rls."ID Centro de Custo" = h."ID Centro de Custo" WHERE rls."E-mail" = :email AND h."ID Colaborador" = FATO_FOLHA_PAGAMENTO."ID Colaborador" AND h."Ano e Mês" = FATO_FOLHA_PAGAMENTO."Ano e Mês")
  AND (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:data_competencia_inicio IS NULL OR "Data Competência" >= :data_competencia_inicio)
  AND (:data_competencia_fim IS NULL OR "Data Competência" <= :data_competencia_fim)
  AND (:periodo_pagamento IS NULL OR "Período Pagamento" = :periodo_pagamento)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."Ano e Mês",
    f."Data Competência",
    f."Período Pagamento",
    f."ID Colaborador",
    f."ID Evento",
    f."Valor Evento",
    f."Quantidade Evento",
    f."Registro Ativo",
    colab."Matrícula Colaborador",
    colab."Nome Colaborador",
    colab."Data Admissão",
    colab."Data Desligamento",
    colab."Situação Atual",
    evt."Código Evento",
    evt."Nome Evento",
    evt."Evento",
    evt."Nome Tipo Evento",
    evt."Tipo Evento",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_FOLHA_PAGAMENTO f
LEFT JOIN DIM_COLABORADOR colab ON colab."ID Colaborador" = f."ID Colaborador"
LEFT JOIN DIM_EVENTO evt ON evt."ID Evento" = f."ID Evento"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Competência"
WHERE EXISTS (SELECT 1 FROM FATO_HISTORICO_COLABORADOR h JOIN DIM_RLS_CENTRO_CUSTO rls ON rls."ID Centro de Custo" = h."ID Centro de Custo" WHERE rls."E-mail" = :email AND h."ID Colaborador" = f."ID Colaborador" AND h."Ano e Mês" = f."Ano e Mês")
  AND (:ano_mes_inicio IS NULL OR f."Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR f."Ano e Mês" <= :ano_mes_fim)
  AND (:data_competencia_inicio IS NULL OR f."Data Competência" >= :data_competencia_inicio)
  AND (:data_competencia_fim IS NULL OR f."Data Competência" <= :data_competencia_fim)
  AND (:periodo_pagamento IS NULL OR f."Período Pagamento" = :periodo_pagamento)
  AND (:id_colaborador IS NULL OR f."ID Colaborador" = :id_colaborador)
  AND (:id_evento IS NULL OR f."ID Evento" = :id_evento)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_colaborador IS NULL OR UPPER(colab."Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:matricula IS NULL OR colab."Matrícula Colaborador" = :matricula)
  AND (:nome_evento IS NULL OR UPPER(evt."Nome Evento") LIKE '%' || UPPER(:nome_evento) || '%')
  AND (:nome_tipo_evento IS NULL OR UPPER(evt."Nome Tipo Evento") LIKE '%' || UPPER(:nome_tipo_evento) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Competência" DESC, f."ID Colaborador"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	AnoMesInicio          *string `json:"ano_mes_inicio,omitempty"`
	AnoMesFim             *string `json:"ano_mes_fim,omitempty"`
	DataCompetenciaInicio *string `json:"data_competencia_inicio,omitempty"`
	DataCompetenciaFim    *string `json:"data_competencia_fim,omitempty"`
	PeriodoPagamento      *string `json:"periodo_pagamento,omitempty"`
	IDColaborador         *int    `json:"id_colaborador,omitempty"`
	IDEvento              *int    `json:"id_evento,omitempty"`
	RegistroAtivo         *string `json:"registro_ativo,omitempty"`
	Limite                *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_competencia_inicio, err := bi.NullDate(in.DataCompetenciaInicio)
	if err != nil {
		return nil, err
	}
	data_competencia_fim, err := bi.NullDate(in.DataCompetenciaFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("ano_mes_inicio", bi.NullString(in.AnoMesInicio)),
		sql.Named("ano_mes_fim", bi.NullString(in.AnoMesFim)),
		sql.Named("data_competencia_inicio", data_competencia_inicio),
		sql.Named("data_competencia_fim", data_competencia_fim),
		sql.Named("periodo_pagamento", bi.NullString(in.PeriodoPagamento)),
		sql.Named("id_colaborador", bi.NullInt(in.IDColaborador)),
		sql.Named("id_evento", bi.NullInt(in.IDEvento)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	AnoMesInicio          *string `json:"ano_mes_inicio,omitempty"`
	AnoMesFim             *string `json:"ano_mes_fim,omitempty"`
	DataCompetenciaInicio *string `json:"data_competencia_inicio,omitempty"`
	DataCompetenciaFim    *string `json:"data_competencia_fim,omitempty"`
	PeriodoPagamento      *string `json:"periodo_pagamento,omitempty"`
	IDColaborador         *int    `json:"id_colaborador,omitempty"`
	IDEvento              *int    `json:"id_evento,omitempty"`
	RegistroAtivo         *string `json:"registro_ativo,omitempty"`
	NomeColaborador       *string `json:"nome_colaborador,omitempty"`
	Matricula             *string `json:"matricula,omitempty"`
	NomeEvento            *string `json:"nome_evento,omitempty"`
	NomeTipoEvento        *string `json:"nome_tipo_evento,omitempty"`
	Ano                   *int    `json:"ano,omitempty"`
	Mes                   *int    `json:"mes,omitempty"`
	Trimestre             *int    `json:"trimestre,omitempty"`
	Limite                *int    `json:"limite,omitempty"`
}

func (in entradaCom) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_competencia_inicio, err := bi.NullDate(in.DataCompetenciaInicio)
	if err != nil {
		return nil, err
	}
	data_competencia_fim, err := bi.NullDate(in.DataCompetenciaFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("ano_mes_inicio", bi.NullString(in.AnoMesInicio)),
		sql.Named("ano_mes_fim", bi.NullString(in.AnoMesFim)),
		sql.Named("data_competencia_inicio", data_competencia_inicio),
		sql.Named("data_competencia_fim", data_competencia_fim),
		sql.Named("periodo_pagamento", bi.NullString(in.PeriodoPagamento)),
		sql.Named("id_colaborador", bi.NullInt(in.IDColaborador)),
		sql.Named("id_evento", bi.NullInt(in.IDEvento)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("nome_colaborador", bi.NullString(in.NomeColaborador)),
		sql.Named("matricula", bi.NullString(in.Matricula)),
		sql.Named("nome_evento", bi.NullString(in.NomeEvento)),
		sql.Named("nome_tipo_evento", bi.NullString(in.NomeTipoEvento)),
		sql.Named("ano", bi.NullInt(in.Ano)),
		sql.Named("mes", bi.NullInt(in.Mes)),
		sql.Named("trimestre", bi.NullInt(in.Trimestre)),
		sql.Named("limite", limite),
	}, nil
}

type saida struct {
	Itens []map[string]any `json:"itens"`
}

// Register publica as duas tools deste fato.
func Register(server *sdk.Server, env bi.ToolEnv) {
	sdk.AddTool(server, &sdk.Tool{Name: ToolName, Description: description}, func(ctx context.Context, req *sdk.CallToolRequest, in entrada) (*sdk.CallToolResult, saida, error) {
		rows, err := Search(ctx, env, userID(req), in)
		if rows == nil {
			rows = []map[string]any{}
		}
		return nil, saida{Itens: rows}, err
	})
	sdk.AddTool(server, &sdk.Tool{Name: ToolNameComDimensoes, Description: descriptionCom}, func(ctx context.Context, req *sdk.CallToolRequest, in entradaCom) (*sdk.CallToolResult, saida, error) {
		rows, err := SearchComDimensoes(ctx, env, userID(req), in)
		if rows == nil {
			rows = []map[string]any{}
		}
		return nil, saida{Itens: rows}, err
	})
}

// Search executa o fato restrito aos centros de custo do usuário.
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
	return bi.QueryRows(ctx, env, ToolName, sqlFato, args)
}

// SearchComDimensoes executa o fato com os rótulos das dimensões.
func SearchComDimensoes(ctx context.Context, env bi.ToolEnv, userID string, in entradaCom) (rows []map[string]any, err error) {
	start := time.Now()
	defer bi.Finish(ctx, env.Audit, env.Log, userID, ToolNameComDimensoes, start, &err)
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
	return bi.QueryRows(ctx, env, ToolNameComDimensoes, sqlComDimensoes, args)
}

func userID(req *sdk.CallToolRequest) string {
	if req == nil || req.Extra == nil || req.Extra.TokenInfo == nil {
		return ""
	}
	return req.Extra.TokenInfo.UserID
}
