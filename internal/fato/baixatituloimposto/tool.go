// Package baixatituloimposto é o fato FATO_BAIXA_TITULO_IMPOSTO.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package baixatituloimposto

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_baixa_titulo_imposto"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_baixa_titulo_imposto_com_dimensoes"

	description = `Fato de movimentações bancárias referentes à baixa (pagamento) de títulos de impostos, com dados do título, conta bancária, natureza, centro de custo, rateios e valores.

Responder perguntas sobre impostos pagos: quanto foi recolhido por período, por filial, por natureza/tipo de imposto, por fornecedor (órgão) ou por conta bancária. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo 'ID Centro de Custo'.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: data_inicio, data_fim, id_filial, id_conta_bancaria, id_tipo_mov_bancaria, operacao, id_natureza, id_centro_custo, id_fornecedor, id_tipo_titulo, numero_titulo, data_vencimento_inicio, data_vencimento_fim, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Baixas (pagamentos) de títulos de impostos enriquecidas com filial, tipo de movimentação, moeda, conta bancária, natureza, centro de custo/contrato, tipo de título, órgão/fornecedor do título, filiais e contas de origem/destino e calendário.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Centro de Custo'.

Responder em uma única chamada perguntas sobre impostos recolhidos com rótulos legíveis — por exemplo, total de impostos pagos por natureza e filial no trimestre, por contrato ou por banco — já restrito aos centros de custo autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: data_inicio, data_fim, id_filial, id_conta_bancaria, id_tipo_mov_bancaria, operacao, id_natureza, id_centro_custo, id_fornecedor, id_tipo_titulo, numero_titulo, data_vencimento_inicio, data_vencimento_fim, registro_ativo, nome_empresa, nome_filial, nome_natureza, nome_natureza_n1, nome_gestor, nome_contrato, nome_centro_custo, razao_social_fornecedor, cnpj_fornecedor, nome_banco, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "Operação Mov. Bancária Bx. Tit. Imposto",
    "ID Filial",
    "Data Mov. Bancária Bx. Tit. Imposto",
    "ID Tipo Mov. Bancária Bx. Tit. Imposto",
    "Código Numerário",
    "Número Documento",
    "Situação Documento",
    "Histórico",
    "Beneficiário",
    "ID Moeda",
    "Valor Mov. Bancária Bx. Tit. Imposto",
    "ID Conta Bancária",
    "Banco Virtual",
    "ID Natureza",
    "ID Centro de Custo",
    "Rateio Financeiro",
    "Rateio Contábil",
    "ID Empresa Filial Origem",
    "ID Empresa Filial Destino",
    "ID Conta Bancária Origem",
    "ID Conta Bancária Destino",
    "Prefixo Título",
    "Número Título",
    "Parcela Título",
    "ID Tipo Título",
    "ID Fornecedor Título",
    "ID Município Fornecedor Título",
    "ID País Fornecedor Título",
    "ID Cliente Título",
    "ID Município Cliente Título",
    "Data Emissão Título",
    "Data Vencimento Título",
    "Número Lote",
    "Motivo Baixa",
    "Sequência Baixa",
    "Número Cheque",
    "ID Processo",
    "Valor Título",
    "Registro Ativo"
FROM FATO_BAIXA_TITULO_IMPOSTO
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = FATO_BAIXA_TITULO_IMPOSTO."ID Centro de Custo")
  AND (:data_inicio IS NULL OR "Data Mov. Bancária Bx. Tit. Imposto" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data Mov. Bancária Bx. Tit. Imposto" <= :data_fim)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR "ID Tipo Mov. Bancária Bx. Tit. Imposto" = :id_tipo_mov_bancaria)
  AND (:operacao IS NULL OR "Operação Mov. Bancária Bx. Tit. Imposto" = :operacao)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor Título" = :id_fornecedor)
  AND (:id_tipo_titulo IS NULL OR "ID Tipo Título" = :id_tipo_titulo)
  AND (:numero_titulo IS NULL OR "Número Título" = :numero_titulo)
  AND (:data_vencimento_inicio IS NULL OR "Data Vencimento Título" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR "Data Vencimento Título" <= :data_vencimento_fim)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Mov. Bancária Bx. Tit. Imposto" DESC
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."Operação Mov. Bancária Bx. Tit. Imposto",
    f."ID Filial",
    f."Data Mov. Bancária Bx. Tit. Imposto",
    f."ID Tipo Mov. Bancária Bx. Tit. Imposto",
    f."Código Numerário",
    f."Número Documento",
    f."Situação Documento",
    f."Histórico",
    f."Beneficiário",
    f."ID Moeda",
    f."Valor Mov. Bancária Bx. Tit. Imposto",
    f."ID Conta Bancária",
    f."Banco Virtual",
    f."ID Natureza",
    f."ID Centro de Custo",
    f."Rateio Financeiro",
    f."Rateio Contábil",
    f."ID Empresa Filial Origem",
    f."ID Empresa Filial Destino",
    f."ID Conta Bancária Origem",
    f."ID Conta Bancária Destino",
    f."Prefixo Título",
    f."Número Título",
    f."Parcela Título",
    f."ID Tipo Título",
    f."ID Fornecedor Título",
    f."ID Município Fornecedor Título",
    f."ID País Fornecedor Título",
    f."ID Cliente Título",
    f."ID Município Cliente Título",
    f."Data Emissão Título",
    f."Data Vencimento Título",
    f."Número Lote",
    f."Motivo Baixa",
    f."Sequência Baixa",
    f."Número Cheque",
    f."ID Processo",
    f."Valor Título",
    f."Registro Ativo",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    tmov."Descrição Tipo Mov. Bancária",
    tmov."Tipo Mov. Bancária",
    moe."Moeda",
    cb."Banco",
    cb."Código Agência",
    cb."Número Conta Bancária",
    cb."Nome Conta Bancária",
    cb."Tipo Conta Bancária",
    cb."Conta Bancária",
    cb."Banco, Agência e Conta Bancária",
    nat."Natureza N1",
    nat."Natureza N2",
    nat."Natureza N3",
    nat."Natureza N4",
    nat."Código Natureza",
    nat."Nome Natureza",
    nat."Natureza",
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
    filo."Empresa" AS "Empresa Origem",
    filo."Filial" AS "Filial Origem",
    fild."Empresa" AS "Empresa Destino",
    fild."Filial" AS "Filial Destino",
    cbo."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Origem",
    cbd."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Destino",
    ttit."Descrição Tipo Título",
    ttit."Tipo Título",
    forn."Fornecedor Razão Social" AS "Fornecedor Razão Social Título",
    forn."Fornecedor Nome Fantasia" AS "Fornecedor Nome Fantasia Título",
    forn."Fornecedor CNPJ Razão Social" AS "Fornecedor CNPJ Razão Social Título",
    forn."CNPJ" AS "CNPJ Fornecedor Título",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor Título",
    forn."Município" AS "Município Fornecedor Título",
    forn."Sigla UF" AS "UF Fornecedor Título",
    forn."País" AS "País Fornecedor Título",
    cli."Cliente Razão Social" AS "Cliente Razão Social Título",
    cli."Cliente Nome Fantasia" AS "Cliente Nome Fantasia Título",
    cli."Cliente CNPJ Razão Social" AS "Cliente CNPJ Razão Social Título",
    cli."CNPJ" AS "CNPJ Cliente Título",
    cli."Tipo Pessoa" AS "Tipo Pessoa Cliente Título",
    cli."Município" AS "Município Cliente Título",
    cli."Sigla UF" AS "UF Cliente Título",
    cli."Região" AS "Região Cliente Título",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_BAIXA_TITULO_IMPOSTO f
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_TIPO_MOV_BANCARIA tmov ON tmov."ID Tipo Mov. Bancária" = f."ID Tipo Mov. Bancária Bx. Tit. Imposto"
LEFT JOIN DIM_MOEDA moe ON moe."ID Moeda" = f."ID Moeda"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cb ON cb."ID Conta Bancária" = f."ID Conta Bancária"
LEFT JOIN DIM_NATUREZA nat ON nat."ID Natureza" = f."ID Natureza"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_EMPRESA_FILIAL filo ON filo."ID Filial" = f."ID Empresa Filial Origem"
LEFT JOIN DIM_EMPRESA_FILIAL fild ON fild."ID Filial" = f."ID Empresa Filial Destino"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbo ON cbo."ID Conta Bancária" = f."ID Conta Bancária Origem"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbd ON cbd."ID Conta Bancária" = f."ID Conta Bancária Destino"
LEFT JOIN DIM_TIPO_TITULO ttit ON ttit."ID Tipo Título" = f."ID Tipo Título"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor Título"
LEFT JOIN DIM_CLIENTE cli ON cli."ID Cliente" = f."ID Cliente Título"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Mov. Bancária Bx. Tit. Imposto"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:data_inicio IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" >= :data_inicio)
  AND (:data_fim IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" <= :data_fim)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR f."ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR f."ID Tipo Mov. Bancária Bx. Tit. Imposto" = :id_tipo_mov_bancaria)
  AND (:operacao IS NULL OR f."Operação Mov. Bancária Bx. Tit. Imposto" = :operacao)
  AND (:id_natureza IS NULL OR f."ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor Título" = :id_fornecedor)
  AND (:id_tipo_titulo IS NULL OR f."ID Tipo Título" = :id_tipo_titulo)
  AND (:numero_titulo IS NULL OR f."Número Título" = :numero_titulo)
  AND (:data_vencimento_inicio IS NULL OR f."Data Vencimento Título" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR f."Data Vencimento Título" <= :data_vencimento_fim)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:nome_natureza IS NULL OR UPPER(nat."Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:nome_natureza_n1 IS NULL OR UPPER(nat."Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_banco IS NULL OR UPPER(cb."Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Mov. Bancária Bx. Tit. Imposto" DESC
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	DataInicio           *string `json:"data_inicio,omitempty"`
	DataFim              *string `json:"data_fim,omitempty"`
	IDFilial             *int    `json:"id_filial,omitempty"`
	IDContaBancaria      *int    `json:"id_conta_bancaria,omitempty"`
	IDTipoMovBancaria    *int    `json:"id_tipo_mov_bancaria,omitempty"`
	Operacao             *string `json:"operacao,omitempty"`
	IDNatureza           *int    `json:"id_natureza,omitempty"`
	IDCentroCusto        *int    `json:"id_centro_custo,omitempty"`
	IDFornecedor         *int    `json:"id_fornecedor,omitempty"`
	IDTipoTitulo         *int    `json:"id_tipo_titulo,omitempty"`
	NumeroTitulo         *string `json:"numero_titulo,omitempty"`
	DataVencimentoInicio *string `json:"data_vencimento_inicio,omitempty"`
	DataVencimentoFim    *string `json:"data_vencimento_fim,omitempty"`
	RegistroAtivo        *string `json:"registro_ativo,omitempty"`
	Limite               *int    `json:"limite,omitempty"`
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
	data_vencimento_inicio, err := bi.NullDate(in.DataVencimentoInicio)
	if err != nil {
		return nil, err
	}
	data_vencimento_fim, err := bi.NullDate(in.DataVencimentoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("data_inicio", data_inicio),
		sql.Named("data_fim", data_fim),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_conta_bancaria", bi.NullInt(in.IDContaBancaria)),
		sql.Named("id_tipo_mov_bancaria", bi.NullInt(in.IDTipoMovBancaria)),
		sql.Named("operacao", bi.NullString(in.Operacao)),
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("id_tipo_titulo", bi.NullInt(in.IDTipoTitulo)),
		sql.Named("numero_titulo", bi.NullString(in.NumeroTitulo)),
		sql.Named("data_vencimento_inicio", data_vencimento_inicio),
		sql.Named("data_vencimento_fim", data_vencimento_fim),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	DataInicio            *string `json:"data_inicio,omitempty"`
	DataFim               *string `json:"data_fim,omitempty"`
	IDFilial              *int    `json:"id_filial,omitempty"`
	IDContaBancaria       *int    `json:"id_conta_bancaria,omitempty"`
	IDTipoMovBancaria     *int    `json:"id_tipo_mov_bancaria,omitempty"`
	Operacao              *string `json:"operacao,omitempty"`
	IDNatureza            *int    `json:"id_natureza,omitempty"`
	IDCentroCusto         *int    `json:"id_centro_custo,omitempty"`
	IDFornecedor          *int    `json:"id_fornecedor,omitempty"`
	IDTipoTitulo          *int    `json:"id_tipo_titulo,omitempty"`
	NumeroTitulo          *string `json:"numero_titulo,omitempty"`
	DataVencimentoInicio  *string `json:"data_vencimento_inicio,omitempty"`
	DataVencimentoFim     *string `json:"data_vencimento_fim,omitempty"`
	RegistroAtivo         *string `json:"registro_ativo,omitempty"`
	NomeEmpresa           *string `json:"nome_empresa,omitempty"`
	NomeFilial            *string `json:"nome_filial,omitempty"`
	NomeNatureza          *string `json:"nome_natureza,omitempty"`
	NomeNaturezaN1        *string `json:"nome_natureza_n1,omitempty"`
	NomeGestor            *string `json:"nome_gestor,omitempty"`
	NomeContrato          *string `json:"nome_contrato,omitempty"`
	NomeCentroCusto       *string `json:"nome_centro_custo,omitempty"`
	RazaoSocialFornecedor *string `json:"razao_social_fornecedor,omitempty"`
	CNPJFornecedor        *string `json:"cnpj_fornecedor,omitempty"`
	NomeBanco             *string `json:"nome_banco,omitempty"`
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

	data_inicio, err := bi.NullDate(in.DataInicio)
	if err != nil {
		return nil, err
	}
	data_fim, err := bi.NullDate(in.DataFim)
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
	return []any{
		sql.Named("data_inicio", data_inicio),
		sql.Named("data_fim", data_fim),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_conta_bancaria", bi.NullInt(in.IDContaBancaria)),
		sql.Named("id_tipo_mov_bancaria", bi.NullInt(in.IDTipoMovBancaria)),
		sql.Named("operacao", bi.NullString(in.Operacao)),
		sql.Named("id_natureza", bi.NullInt(in.IDNatureza)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("id_tipo_titulo", bi.NullInt(in.IDTipoTitulo)),
		sql.Named("numero_titulo", bi.NullString(in.NumeroTitulo)),
		sql.Named("data_vencimento_inicio", data_vencimento_inicio),
		sql.Named("data_vencimento_fim", data_vencimento_fim),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("nome_empresa", bi.NullString(in.NomeEmpresa)),
		sql.Named("nome_filial", bi.NullString(in.NomeFilial)),
		sql.Named("nome_natureza", bi.NullString(in.NomeNatureza)),
		sql.Named("nome_natureza_n1", bi.NullString(in.NomeNaturezaN1)),
		sql.Named("nome_gestor", bi.NullString(in.NomeGestor)),
		sql.Named("nome_contrato", bi.NullString(in.NomeContrato)),
		sql.Named("nome_centro_custo", bi.NullString(in.NomeCentroCusto)),
		sql.Named("razao_social_fornecedor", bi.NullString(in.RazaoSocialFornecedor)),
		sql.Named("cnpj_fornecedor", bi.NullString(in.CNPJFornecedor)),
		sql.Named("nome_banco", bi.NullString(in.NomeBanco)),
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
