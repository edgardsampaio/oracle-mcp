// Package orcamento é o fato FATO_ORCAMENTO_TOTAL.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package orcamento

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_orcamento_total"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_orcamento_total_com_dimensoes"

	description = `Fato de orçamentos por versão, item, conta orçamentária e contrato, com valor orçado original, atual e a variação entre eles.

Responder perguntas de orçamento: quanto foi orçado por conta, contrato ou período, revisões entre versões e, via DIM_CONTA_ORCM_NATUREZA, comparações orçado x realizado.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_orcamento, codigo_orcamento, numero_versao, id_conta_orcamentaria, id_contrato, periodo_inicio, periodo_fim, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Orçamentos por versão e item enriquecidos com a conta orçamentária (N1–N3), a hierarquia do contrato e o calendário do início do período.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Contrato'.

Responder perguntas de orçamento com rótulos legíveis — orçado original x atual por conta, contrato, gestor ou segmento — já restrito aos contratos autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_orcamento, codigo_orcamento, numero_versao, id_conta_orcamentaria, id_contrato, periodo_inicio, periodo_fim, nome_conta_orcamentaria, nome_conta_n1, nome_gestor, nome_contrato, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "ID Orçamento",
    "Código Orçamento",
    "Descrição Orçamento",
    "Orçamento",
    "Data Início Período",
    "Data Fim Período",
    "Número Versão",
    "Número Item",
    "ID Conta Orçamentária",
    "ID Contrato",
    "Valor Total Orçado Original",
    "Valor Total Orçado Atual",
    "Valor Variação Orig. x Atual"
FROM FATO_ORCAMENTO_TOTAL
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = FATO_ORCAMENTO_TOTAL."ID Contrato")
  AND (:id_orcamento IS NULL OR "ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR "Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR "Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR "Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR "Data Início Período" <= :periodo_fim)
ORDER BY "Data Início Período", "Número Item"
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."ID Orçamento",
    f."Código Orçamento",
    f."Descrição Orçamento",
    f."Orçamento",
    f."Data Início Período",
    f."Data Fim Período",
    f."Número Versão",
    f."Número Item",
    f."ID Conta Orçamentária",
    f."ID Contrato",
    f."Valor Total Orçado Original",
    f."Valor Total Orçado Atual",
    f."Valor Variação Orig. x Atual",
    corc."Código Conta Orçamentária",
    corc."Nome Conta Orçamentária",
    corc."Conta Orçamentária",
    corc."Conta Orçamentária N1",
    corc."Conta Orçamentária N2",
    corc."Conta Orçamentária N3",
    ctr."Nome Gestor",
    ctr."Gestor",
    ctr."Nome Segmento",
    ctr."Segmento",
    ctr."Nome Grupo Contrato",
    ctr."Grupo Contrato",
    ctr."Código Contrato",
    ctr."Nome Contrato",
    ctr."Contrato",
    ctr."Nome Tipo Contrato",
    ctr."Tipo Contrato",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_ORCAMENTO_TOTAL f
LEFT JOIN DIM_CONTA_ORCAMENTARIA corc ON corc."ID Conta Orçamentária" = f."ID Conta Orçamentária"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Início Período"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_orcamento IS NULL OR f."ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR f."Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR f."Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR f."ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR f."Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR f."Data Início Período" <= :periodo_fim)
  AND (:nome_conta_orcamentaria IS NULL OR UPPER(corc."Nome Conta Orçamentária") LIKE '%' || UPPER(:nome_conta_orcamentaria) || '%')
  AND (:nome_conta_n1 IS NULL OR UPPER(corc."Nome Conta Orçamentária N1") LIKE '%' || UPPER(:nome_conta_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Início Período", f."Número Item"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDOrcamento         *int    `json:"id_orcamento,omitempty"`
	CodigoOrcamento     *string `json:"codigo_orcamento,omitempty"`
	NumeroVersao        *string `json:"numero_versao,omitempty"`
	IDContaOrcamentaria *int    `json:"id_conta_orcamentaria,omitempty"`
	IDContrato          *int    `json:"id_contrato,omitempty"`
	PeriodoInicio       *string `json:"periodo_inicio,omitempty"`
	PeriodoFim          *string `json:"periodo_fim,omitempty"`
	Limite              *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	periodo_inicio, err := bi.NullDate(in.PeriodoInicio)
	if err != nil {
		return nil, err
	}
	periodo_fim, err := bi.NullDate(in.PeriodoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_orcamento", bi.NullInt(in.IDOrcamento)),
		sql.Named("codigo_orcamento", bi.NullString(in.CodigoOrcamento)),
		sql.Named("numero_versao", bi.NullString(in.NumeroVersao)),
		sql.Named("id_conta_orcamentaria", bi.NullInt(in.IDContaOrcamentaria)),
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("periodo_inicio", periodo_inicio),
		sql.Named("periodo_fim", periodo_fim),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	IDOrcamento           *int    `json:"id_orcamento,omitempty"`
	CodigoOrcamento       *string `json:"codigo_orcamento,omitempty"`
	NumeroVersao          *string `json:"numero_versao,omitempty"`
	IDContaOrcamentaria   *int    `json:"id_conta_orcamentaria,omitempty"`
	IDContrato            *int    `json:"id_contrato,omitempty"`
	PeriodoInicio         *string `json:"periodo_inicio,omitempty"`
	PeriodoFim            *string `json:"periodo_fim,omitempty"`
	NomeContaOrcamentaria *string `json:"nome_conta_orcamentaria,omitempty"`
	NomeContaN1           *string `json:"nome_conta_n1,omitempty"`
	NomeGestor            *string `json:"nome_gestor,omitempty"`
	NomeContrato          *string `json:"nome_contrato,omitempty"`
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

	periodo_inicio, err := bi.NullDate(in.PeriodoInicio)
	if err != nil {
		return nil, err
	}
	periodo_fim, err := bi.NullDate(in.PeriodoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_orcamento", bi.NullInt(in.IDOrcamento)),
		sql.Named("codigo_orcamento", bi.NullString(in.CodigoOrcamento)),
		sql.Named("numero_versao", bi.NullString(in.NumeroVersao)),
		sql.Named("id_conta_orcamentaria", bi.NullInt(in.IDContaOrcamentaria)),
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("periodo_inicio", periodo_inicio),
		sql.Named("periodo_fim", periodo_fim),
		sql.Named("nome_conta_orcamentaria", bi.NullString(in.NomeContaOrcamentaria)),
		sql.Named("nome_conta_n1", bi.NullString(in.NomeContaN1)),
		sql.Named("nome_gestor", bi.NullString(in.NomeGestor)),
		sql.Named("nome_contrato", bi.NullString(in.NomeContrato)),
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
