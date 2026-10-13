// Package confissaodivida é o fato FATO_CONFISSAO_DIVIDA.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package confissaodivida

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_confissao_divida"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_confissao_divida_com_dimensoes"

	description = `Fato de saldos de confissões de dívida por fornecedor e contrato, com a data do último pagamento.

Responder quanto ainda se deve em confissões de dívida, a quem, em quais contratos e há quanto tempo não há pagamento.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_fornecedor, id_contrato, valor_saldo_minimo, data_ultimo_pagamento_inicio, data_ultimo_pagamento_fim, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Saldos de confissões de dívida enriquecidos com os dados do fornecedor, a hierarquia do contrato (gestor, segmento, grupo, tipo) e o calendário da data do último pagamento.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Contrato'.

Responder quanto se deve em confissões de dívida por fornecedor, contrato, gestor ou segmento, e há quanto tempo não há pagamento, já restrito aos contratos autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_fornecedor, id_contrato, valor_saldo_minimo, data_ultimo_pagamento_inicio, data_ultimo_pagamento_fim, razao_social_fornecedor, cnpj_fornecedor, nome_gestor, nome_contrato, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "ID Fornecedor",
    "ID Contrato",
    "Valor Saldo",
    "Data Último Pagamento"
FROM FATO_CONFISSAO_DIVIDA
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = FATO_CONFISSAO_DIVIDA."ID Contrato")
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:valor_saldo_minimo IS NULL OR "Valor Saldo" >= :valor_saldo_minimo)
  AND (:data_ultimo_pagamento_inicio IS NULL OR "Data Último Pagamento" >= :data_ultimo_pagamento_inicio)
  AND (:data_ultimo_pagamento_fim IS NULL OR "Data Último Pagamento" <= :data_ultimo_pagamento_fim)
ORDER BY "Valor Saldo" DESC
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."ID Fornecedor",
    f."ID Contrato",
    f."Valor Saldo",
    f."Data Último Pagamento",
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor",
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
FROM FATO_CONFISSAO_DIVIDA f
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Último Pagamento"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:valor_saldo_minimo IS NULL OR f."Valor Saldo" >= :valor_saldo_minimo)
  AND (:data_ultimo_pagamento_inicio IS NULL OR f."Data Último Pagamento" >= :data_ultimo_pagamento_inicio)
  AND (:data_ultimo_pagamento_fim IS NULL OR f."Data Último Pagamento" <= :data_ultimo_pagamento_fim)
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Valor Saldo" DESC
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDFornecedor              *int     `json:"id_fornecedor,omitempty"`
	IDContrato                *int     `json:"id_contrato,omitempty"`
	ValorSaldoMinimo          *float64 `json:"valor_saldo_minimo,omitempty"`
	DataUltimoPagamentoInicio *string  `json:"data_ultimo_pagamento_inicio,omitempty"`
	DataUltimoPagamentoFim    *string  `json:"data_ultimo_pagamento_fim,omitempty"`
	Limite                    *int     `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_ultimo_pagamento_inicio, err := bi.NullDate(in.DataUltimoPagamentoInicio)
	if err != nil {
		return nil, err
	}
	data_ultimo_pagamento_fim, err := bi.NullDate(in.DataUltimoPagamentoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("valor_saldo_minimo", bi.NullFloat(in.ValorSaldoMinimo)),
		sql.Named("data_ultimo_pagamento_inicio", data_ultimo_pagamento_inicio),
		sql.Named("data_ultimo_pagamento_fim", data_ultimo_pagamento_fim),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	IDFornecedor              *int     `json:"id_fornecedor,omitempty"`
	IDContrato                *int     `json:"id_contrato,omitempty"`
	ValorSaldoMinimo          *float64 `json:"valor_saldo_minimo,omitempty"`
	DataUltimoPagamentoInicio *string  `json:"data_ultimo_pagamento_inicio,omitempty"`
	DataUltimoPagamentoFim    *string  `json:"data_ultimo_pagamento_fim,omitempty"`
	RazaoSocialFornecedor     *string  `json:"razao_social_fornecedor,omitempty"`
	CNPJFornecedor            *string  `json:"cnpj_fornecedor,omitempty"`
	NomeGestor                *string  `json:"nome_gestor,omitempty"`
	NomeContrato              *string  `json:"nome_contrato,omitempty"`
	Ano                       *int     `json:"ano,omitempty"`
	Mes                       *int     `json:"mes,omitempty"`
	Trimestre                 *int     `json:"trimestre,omitempty"`
	Limite                    *int     `json:"limite,omitempty"`
}

func (in entradaCom) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_ultimo_pagamento_inicio, err := bi.NullDate(in.DataUltimoPagamentoInicio)
	if err != nil {
		return nil, err
	}
	data_ultimo_pagamento_fim, err := bi.NullDate(in.DataUltimoPagamentoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("id_contrato", bi.NullInt(in.IDContrato)),
		sql.Named("valor_saldo_minimo", bi.NullFloat(in.ValorSaldoMinimo)),
		sql.Named("data_ultimo_pagamento_inicio", data_ultimo_pagamento_inicio),
		sql.Named("data_ultimo_pagamento_fim", data_ultimo_pagamento_fim),
		sql.Named("razao_social_fornecedor", bi.NullString(in.RazaoSocialFornecedor)),
		sql.Named("cnpj_fornecedor", bi.NullString(in.CNPJFornecedor)),
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
