// Package contareceber é o fato FATO_CONTA_RECEBER_RATEADO.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package contareceber

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_conta_receber_rateado"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_conta_receber_rateado_com_dimensoes"

	description = `Fato de títulos de Contas a Receber rateados por natureza e centro de custo, com datas (emissão, vencimento, última baixa), valores originais, recebidos, saldos, tributos retidos e abatimentos.

Responder perguntas de Contas a Receber: inadimplência, títulos vencidos e a vencer, valores recebidos por período e carteira por cliente, natureza ou centro de custo. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo 'ID Centro de Custo'.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: data_vencimento_inicio, data_vencimento_fim, data_emissao_inicio, data_emissao_fim, data_baixa_inicio, data_baixa_fim, id_filial, id_cliente, id_natureza, id_centro_custo, id_tipo_titulo, numero, somente_em_aberto, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Títulos de Contas a Receber rateados, enriquecidos com filial, tipo de título, cliente, natureza (N1–N4), moeda, centro de custo/contrato e calendário do vencimento real.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Centro de Custo'.

Responder perguntas de Contas a Receber com rótulos legíveis — inadimplência, vencidos, a vencer e recebidos por cliente, contrato, gestor, UF ou filial — em uma única chamada, já restrito aos centros de custo autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: data_vencimento_inicio, data_vencimento_fim, data_emissao_inicio, data_emissao_fim, data_baixa_inicio, data_baixa_fim, id_filial, id_cliente, id_natureza, id_centro_custo, id_tipo_titulo, numero, somente_em_aberto, registro_ativo, nome_empresa, nome_filial, razao_social_cliente, cnpj_cliente, nome_natureza, nome_natureza_n1, nome_gestor, nome_contrato, nome_centro_custo, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "ID Filial",
    "Prefixo",
    "Número",
    "Parcela",
    "ID Tipo Título",
    "Código Tipo Título",
    "ID Cliente",
    "Código Cliente",
    "Código Loja Cliente",
    "ID Município",
    "ID Natureza",
    "ID Moeda",
    "ID Centro de Custo",
    "Data Emissão",
    "Data Vencimento",
    "Data Vencimento Real",
    "Data Última Baixa",
    "Histórico",
    "Taxa Moeda",
    "Valor Título",
    "Valor Título Real",
    "Valor Última Baixa",
    "Valor Saldo",
    "Valor Saldo Real",
    "Valor Recebido Real",
    "Rateio",
    "Registro Ativo",
    "Valor ISS",
    "Valor INSS",
    "Valor IRRF",
    "Valor PIS",
    "Valor COFINS",
    "Valor CSLL",
    "Valor Saldo Acréscimo",
    "Valor Saldo Decréscimo",
    "Chave Título Pai",
    "Valor Total Tributos Reter",
    "Valor Líquido Título Real",
    "Valor Abatimentos",
    "Valor Tributos Retidos Baixa",
    "Valor Líquido Recebido Real",
    "Valor Líquido Saldo Real"
FROM FATO_CONTA_RECEBER_RATEADO
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = FATO_CONTA_RECEBER_RATEADO."ID Centro de Custo")
  AND (:data_vencimento_inicio IS NULL OR "Data Vencimento Real" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR "Data Vencimento Real" <= :data_vencimento_fim)
  AND (:data_emissao_inicio IS NULL OR "Data Emissão" >= :data_emissao_inicio)
  AND (:data_emissao_fim IS NULL OR "Data Emissão" <= :data_emissao_fim)
  AND (:data_baixa_inicio IS NULL OR "Data Última Baixa" >= :data_baixa_inicio)
  AND (:data_baixa_fim IS NULL OR "Data Última Baixa" <= :data_baixa_fim)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_cliente IS NULL OR "ID Cliente" = :id_cliente)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_tipo_titulo IS NULL OR "ID Tipo Título" = :id_tipo_titulo)
  AND (:numero IS NULL OR "Número" = :numero)
  AND (:somente_em_aberto IS NULL OR :somente_em_aberto <> 'S' OR "Valor Saldo Real" > 0)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Vencimento Real"
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."ID Filial",
    f."Prefixo",
    f."Número",
    f."Parcela",
    f."ID Tipo Título",
    f."Código Tipo Título",
    f."ID Cliente",
    f."Código Cliente",
    f."Código Loja Cliente",
    f."ID Município",
    f."ID Natureza",
    f."ID Moeda",
    f."ID Centro de Custo",
    f."Data Emissão",
    f."Data Vencimento",
    f."Data Vencimento Real",
    f."Data Última Baixa",
    f."Histórico",
    f."Taxa Moeda",
    f."Valor Título",
    f."Valor Título Real",
    f."Valor Última Baixa",
    f."Valor Saldo",
    f."Valor Saldo Real",
    f."Valor Recebido Real",
    f."Rateio",
    f."Registro Ativo",
    f."Valor ISS",
    f."Valor INSS",
    f."Valor IRRF",
    f."Valor PIS",
    f."Valor COFINS",
    f."Valor CSLL",
    f."Valor Saldo Acréscimo",
    f."Valor Saldo Decréscimo",
    f."Chave Título Pai",
    f."Valor Total Tributos Reter",
    f."Valor Líquido Título Real",
    f."Valor Abatimentos",
    f."Valor Tributos Retidos Baixa",
    f."Valor Líquido Recebido Real",
    f."Valor Líquido Saldo Real",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    ttit."Descrição Tipo Título",
    ttit."Tipo Título",
    cli."Cliente Razão Social",
    cli."Cliente Nome Fantasia",
    cli."Cliente CNPJ Razão Social",
    cli."CNPJ" AS "CNPJ Cliente",
    cli."Tipo Pessoa" AS "Tipo Pessoa Cliente",
    cli."Município" AS "Município Cliente",
    cli."Sigla UF" AS "UF Cliente",
    cli."Região" AS "Região Cliente",
    nat."Natureza N1",
    nat."Natureza N2",
    nat."Natureza N3",
    nat."Natureza N4",
    nat."Código Natureza",
    nat."Nome Natureza",
    nat."Natureza",
    moe."Moeda",
    ccu."Nome Gestor",
    ccu."Gestor",
    ccu."Nome Segmento",
    ccu."Segmento",
    ccu."Nome Grupo Contrato",
    ccu."Grupo Contrato",
    ccu."Código Contrato",
    ccu."Nome Contrato",
    ccu."Contrato",
    ccu."Código Centro de Custo",
    ccu."Nome Centro de Custo",
    ccu."Centro de Custo",
    ccu."Nome Cliente Final",
    ccu."Cliente Final",
    ccu."Nome Tipo Contrato",
    ccu."Tipo Contrato",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_CONTA_RECEBER_RATEADO f
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_TIPO_TITULO ttit ON ttit."ID Tipo Título" = f."ID Tipo Título"
LEFT JOIN DIM_CLIENTE cli ON cli."ID Cliente" = f."ID Cliente"
LEFT JOIN DIM_NATUREZA nat ON nat."ID Natureza" = f."ID Natureza"
LEFT JOIN DIM_MOEDA moe ON moe."ID Moeda" = f."ID Moeda"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Vencimento Real"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:data_vencimento_inicio IS NULL OR f."Data Vencimento Real" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR f."Data Vencimento Real" <= :data_vencimento_fim)
  AND (:data_emissao_inicio IS NULL OR f."Data Emissão" >= :data_emissao_inicio)
  AND (:data_emissao_fim IS NULL OR f."Data Emissão" <= :data_emissao_fim)
  AND (:data_baixa_inicio IS NULL OR f."Data Última Baixa" >= :data_baixa_inicio)
  AND (:data_baixa_fim IS NULL OR f."Data Última Baixa" <= :data_baixa_fim)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_cliente IS NULL OR f."ID Cliente" = :id_cliente)
  AND (:id_natureza IS NULL OR f."ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_tipo_titulo IS NULL OR f."ID Tipo Título" = :id_tipo_titulo)
  AND (:numero IS NULL OR f."Número" = :numero)
  AND (:somente_em_aberto IS NULL OR :somente_em_aberto <> 'S' OR f."Valor Saldo Real" > 0)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:razao_social_cliente IS NULL OR UPPER(cli."Razão Social") LIKE '%' || UPPER(:razao_social_cliente) || '%')
  AND (:cnpj_cliente IS NULL OR cli."CNPJ" = :cnpj_cliente)
  AND (:nome_natureza IS NULL OR UPPER(nat."Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:nome_natureza_n1 IS NULL OR UPPER(nat."Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Vencimento Real"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	DataVencimentoInicio *string `json:"data_vencimento_inicio,omitempty"`
	DataVencimentoFim    *string `json:"data_vencimento_fim,omitempty"`
	DataEmissaoInicio    *string `json:"data_emissao_inicio,omitempty"`
	DataEmissaoFim       *string `json:"data_emissao_fim,omitempty"`
	DataBaixaInicio      *string `json:"data_baixa_inicio,omitempty"`
	DataBaixaFim         *string `json:"data_baixa_fim,omitempty"`
	IDFilial             *int    `json:"id_filial,omitempty"`
	IDCliente            *int    `json:"id_cliente,omitempty"`
	IDNatureza           *int    `json:"id_natureza,omitempty"`
	IDCentroCusto        *int    `json:"id_centro_custo,omitempty"`
	IDTipoTitulo         *int    `json:"id_tipo_titulo,omitempty"`
	Numero               *string `json:"numero,omitempty"`
	SomenteEmAberto      *string `json:"somente_em_aberto,omitempty"`
	RegistroAtivo        *string `json:"registro_ativo,omitempty"`
	Limite               *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_vencimento_inicio, err := bi.NullDate(in.DataVencimentoInicio)
	if err != nil {
		return nil, err
	}
	data_vencimento_fim, err := bi.NullDate(in.DataVencimentoFim)
	if err != nil {
		return nil, err
	}
	data_emissao_inicio, err := bi.NullDate(in.DataEmissaoInicio)
	if err != nil {
		return nil, err
	}
	data_emissao_fim, err := bi.NullDate(in.DataEmissaoFim)
	if err != nil {
		return nil, err
	}
	data_baixa_inicio, err := bi.NullDate(in.DataBaixaInicio)
	if err != nil {
		return nil, err
	}
	data_baixa_fim, err := bi.NullDate(in.DataBaixaFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("data_vencimento_inicio", data_vencimento_inicio),
		sql.Named("data_vencimento_fim", data_vencimento_fim),
		sql.Named("data_emissao_inicio", data_emissao_inicio),
		sql.Named("data_emissao_fim", data_emissao_fim),
		sql.Named("data_baixa_inicio", data_baixa_inicio),
		sql.Named("data_baixa_fim", data_baixa_fim),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_cliente", bi.NullInt(in.IDCliente)),
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_tipo_titulo", bi.NullInt(in.IDTipoTitulo)),
		sql.Named("numero", bi.NullString(in.Numero)),
		sql.Named("somente_em_aberto", bi.NullString(in.SomenteEmAberto)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	DataVencimentoInicio *string `json:"data_vencimento_inicio,omitempty"`
	DataVencimentoFim    *string `json:"data_vencimento_fim,omitempty"`
	DataEmissaoInicio    *string `json:"data_emissao_inicio,omitempty"`
	DataEmissaoFim       *string `json:"data_emissao_fim,omitempty"`
	DataBaixaInicio      *string `json:"data_baixa_inicio,omitempty"`
	DataBaixaFim         *string `json:"data_baixa_fim,omitempty"`
	IDFilial             *int    `json:"id_filial,omitempty"`
	IDCliente            *int    `json:"id_cliente,omitempty"`
	IDNatureza           *int    `json:"id_natureza,omitempty"`
	IDCentroCusto        *int    `json:"id_centro_custo,omitempty"`
	IDTipoTitulo         *int    `json:"id_tipo_titulo,omitempty"`
	Numero               *string `json:"numero,omitempty"`
	SomenteEmAberto      *string `json:"somente_em_aberto,omitempty"`
	RegistroAtivo        *string `json:"registro_ativo,omitempty"`
	NomeEmpresa          *string `json:"nome_empresa,omitempty"`
	NomeFilial           *string `json:"nome_filial,omitempty"`
	RazaoSocialCliente   *string `json:"razao_social_cliente,omitempty"`
	CNPJCliente          *string `json:"cnpj_cliente,omitempty"`
	NomeNatureza         *string `json:"nome_natureza,omitempty"`
	NomeNaturezaN1       *string `json:"nome_natureza_n1,omitempty"`
	NomeGestor           *string `json:"nome_gestor,omitempty"`
	NomeContrato         *string `json:"nome_contrato,omitempty"`
	NomeCentroCusto      *string `json:"nome_centro_custo,omitempty"`
	Ano                  *int    `json:"ano,omitempty"`
	Mes                  *int    `json:"mes,omitempty"`
	Trimestre            *int    `json:"trimestre,omitempty"`
	Limite               *int    `json:"limite,omitempty"`
}

func (in entradaCom) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_vencimento_inicio, err := bi.NullDate(in.DataVencimentoInicio)
	if err != nil {
		return nil, err
	}
	data_vencimento_fim, err := bi.NullDate(in.DataVencimentoFim)
	if err != nil {
		return nil, err
	}
	data_emissao_inicio, err := bi.NullDate(in.DataEmissaoInicio)
	if err != nil {
		return nil, err
	}
	data_emissao_fim, err := bi.NullDate(in.DataEmissaoFim)
	if err != nil {
		return nil, err
	}
	data_baixa_inicio, err := bi.NullDate(in.DataBaixaInicio)
	if err != nil {
		return nil, err
	}
	data_baixa_fim, err := bi.NullDate(in.DataBaixaFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("data_vencimento_inicio", data_vencimento_inicio),
		sql.Named("data_vencimento_fim", data_vencimento_fim),
		sql.Named("data_emissao_inicio", data_emissao_inicio),
		sql.Named("data_emissao_fim", data_emissao_fim),
		sql.Named("data_baixa_inicio", data_baixa_inicio),
		sql.Named("data_baixa_fim", data_baixa_fim),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_cliente", bi.NullInt(in.IDCliente)),
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_tipo_titulo", bi.NullInt(in.IDTipoTitulo)),
		sql.Named("numero", bi.NullString(in.Numero)),
		sql.Named("somente_em_aberto", bi.NullString(in.SomenteEmAberto)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("nome_empresa", bi.NullString(in.NomeEmpresa)),
		sql.Named("nome_filial", bi.NullString(in.NomeFilial)),
		sql.Named("razao_social_cliente", bi.NullString(in.RazaoSocialCliente)),
		sql.Named("cnpj_cliente", bi.NullString(in.CNPJCliente)),
		sql.Named("nome_natureza", bi.NullString(in.NomeNatureza)),
		sql.Named("nome_natureza_n1", bi.NullString(in.NomeNaturezaN1)),
		sql.Named("nome_gestor", bi.NullString(in.NomeGestor)),
		sql.Named("nome_contrato", bi.NullString(in.NomeContrato)),
		sql.Named("nome_centro_custo", bi.NullString(in.NomeCentroCusto)),
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
