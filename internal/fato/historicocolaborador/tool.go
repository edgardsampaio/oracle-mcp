// Package historicocolaborador é o fato FATO_HISTORICO_COLABORADOR.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package historicocolaborador

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_historico_colaborador"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_historico_colaborador_com_dimensoes"

	description = `Fato mensal (snapshot por competência) da posição de cada colaborador: filial, seção, centro de custo, função, situação, forma de recebimento, banco e salário base.

Responder perguntas de headcount e evolução de quadro: quantos colaboradores ativos por mês, por filial, seção, função ou centro de custo, e evolução do salário base. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo 'ID Centro de Custo'.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: ano_mes_inicio, ano_mes_fim, id_colaborador, id_filial, id_secao, id_centro_custo, id_funcao_colaborador, id_situacao_colaborador, id_tipo_recebimento, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Snapshot mensal dos colaboradores enriquecido com dados do colaborador, filial, seção, centro de custo/contrato, função, situação, forma de recebimento, banco e calendário da competência.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Centro de Custo'.

Responder perguntas de headcount e quadro com rótulos legíveis — ativos por mês, por função, seção, filial, contrato ou gestor, e evolução do salário base — já restrito aos centros de custo autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: ano_mes_inicio, ano_mes_fim, id_colaborador, id_filial, id_secao, id_centro_custo, id_funcao_colaborador, id_situacao_colaborador, id_tipo_recebimento, registro_ativo, nome_colaborador, nome_empresa, nome_filial, nome_secao, descricao_funcao, descricao_situacao, nome_gestor, nome_contrato, nome_centro_custo, ano, mes, trimestre, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "Ano e Mês",
    "Data Competência",
    "ID Colaborador",
    "ID Filial",
    "ID Seção",
    "ID Centro de Custo",
    "ID Função Colaborador",
    "ID Situação Colaborador",
    "ID Tipo Recebimento",
    "ID Banco",
    "Valor Salário Base",
    "Registro Ativo Colaborador"
FROM FATO_HISTORICO_COLABORADOR
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = FATO_HISTORICO_COLABORADOR."ID Centro de Custo")
  AND (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_secao IS NULL OR "ID Seção" = :id_secao)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_funcao_colaborador IS NULL OR "ID Função Colaborador" = :id_funcao_colaborador)
  AND (:id_situacao_colaborador IS NULL OR "ID Situação Colaborador" = :id_situacao_colaborador)
  AND (:id_tipo_recebimento IS NULL OR "ID Tipo Recebimento" = :id_tipo_recebimento)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."Ano e Mês",
    f."Data Competência",
    f."ID Colaborador",
    f."ID Filial",
    f."ID Seção",
    f."ID Centro de Custo",
    f."ID Função Colaborador",
    f."ID Situação Colaborador",
    f."ID Tipo Recebimento",
    f."ID Banco",
    f."Valor Salário Base",
    f."Registro Ativo Colaborador",
    colab."Matrícula Colaborador",
    colab."Nome Colaborador",
    colab."Data Admissão",
    colab."Data Desligamento",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    sec."Código Seção",
    sec."Nome Seção",
    sec."Seção",
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
    fun."Descrição Função Colaborador",
    fun."Função Colaborador",
    sit."Descrição Situação Colaborador",
    sit."Situação Colaborador",
    trec."Descrição Tipo Recebimento",
    trec."Tipo Recebimento",
    bco."Código Banco",
    bco."Nome Banco",
    bco."Banco",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_HISTORICO_COLABORADOR f
LEFT JOIN DIM_COLABORADOR colab ON colab."ID Colaborador" = f."ID Colaborador"
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_SECAO sec ON sec."ID Seção" = f."ID Seção"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_FUNCAO_COLABORADOR fun ON fun."ID Função Colaborador" = f."ID Função Colaborador"
LEFT JOIN DIM_SITUACAO_COLABORADOR sit ON sit."ID Situação Colaborador" = f."ID Situação Colaborador"
LEFT JOIN DIM_TIPO_RECEBIMENTO_COLAB trec ON trec."ID Tipo Recebimento" = f."ID Tipo Recebimento"
LEFT JOIN (SELECT DISTINCT "ID Banco", "Código Banco", "Nome Banco", "Banco" FROM DIM_BANCO_CONTA_BANCARIA) bco ON bco."ID Banco" = f."ID Banco"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Competência"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:ano_mes_inicio IS NULL OR f."Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR f."Ano e Mês" <= :ano_mes_fim)
  AND (:id_colaborador IS NULL OR f."ID Colaborador" = :id_colaborador)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_secao IS NULL OR f."ID Seção" = :id_secao)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_funcao_colaborador IS NULL OR f."ID Função Colaborador" = :id_funcao_colaborador)
  AND (:id_situacao_colaborador IS NULL OR f."ID Situação Colaborador" = :id_situacao_colaborador)
  AND (:id_tipo_recebimento IS NULL OR f."ID Tipo Recebimento" = :id_tipo_recebimento)
  AND (:registro_ativo IS NULL OR f."Registro Ativo Colaborador" = :registro_ativo)
  AND (:nome_colaborador IS NULL OR UPPER(colab."Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:nome_secao IS NULL OR UPPER(sec."Nome Seção") LIKE '%' || UPPER(:nome_secao) || '%')
  AND (:descricao_funcao IS NULL OR UPPER(fun."Descrição Função Colaborador") LIKE '%' || UPPER(:descricao_funcao) || '%')
  AND (:descricao_situacao IS NULL OR UPPER(sit."Descrição Situação Colaborador") LIKE '%' || UPPER(:descricao_situacao) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Competência" DESC, f."ID Colaborador"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	AnoMesInicio          *string `json:"ano_mes_inicio,omitempty"`
	AnoMesFim             *string `json:"ano_mes_fim,omitempty"`
	IDColaborador         *int    `json:"id_colaborador,omitempty"`
	IDFilial              *int    `json:"id_filial,omitempty"`
	IDSecao               *int    `json:"id_secao,omitempty"`
	IDCentroCusto         *int    `json:"id_centro_custo,omitempty"`
	IDFuncaoColaborador   *int    `json:"id_funcao_colaborador,omitempty"`
	IDSituacaoColaborador *int    `json:"id_situacao_colaborador,omitempty"`
	IDTipoRecebimento     *int    `json:"id_tipo_recebimento,omitempty"`
	RegistroAtivo         *string `json:"registro_ativo,omitempty"`
	Limite                *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("ano_mes_inicio", bi.NullString(in.AnoMesInicio)),
		sql.Named("ano_mes_fim", bi.NullString(in.AnoMesFim)),
		sql.Named("id_colaborador", bi.NullInt(in.IDColaborador)),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_secao", bi.NullInt(in.IDSecao)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_funcao_colaborador", bi.NullInt(in.IDFuncaoColaborador)),
		sql.Named("id_situacao_colaborador", bi.NullInt(in.IDSituacaoColaborador)),
		sql.Named("id_tipo_recebimento", bi.NullInt(in.IDTipoRecebimento)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	AnoMesInicio          *string `json:"ano_mes_inicio,omitempty"`
	AnoMesFim             *string `json:"ano_mes_fim,omitempty"`
	IDColaborador         *int    `json:"id_colaborador,omitempty"`
	IDFilial              *int    `json:"id_filial,omitempty"`
	IDSecao               *int    `json:"id_secao,omitempty"`
	IDCentroCusto         *int    `json:"id_centro_custo,omitempty"`
	IDFuncaoColaborador   *int    `json:"id_funcao_colaborador,omitempty"`
	IDSituacaoColaborador *int    `json:"id_situacao_colaborador,omitempty"`
	IDTipoRecebimento     *int    `json:"id_tipo_recebimento,omitempty"`
	RegistroAtivo         *string `json:"registro_ativo,omitempty"`
	NomeColaborador       *string `json:"nome_colaborador,omitempty"`
	NomeEmpresa           *string `json:"nome_empresa,omitempty"`
	NomeFilial            *string `json:"nome_filial,omitempty"`
	NomeSecao             *string `json:"nome_secao,omitempty"`
	DescricaoFuncao       *string `json:"descricao_funcao,omitempty"`
	DescricaoSituacao     *string `json:"descricao_situacao,omitempty"`
	NomeGestor            *string `json:"nome_gestor,omitempty"`
	NomeContrato          *string `json:"nome_contrato,omitempty"`
	NomeCentroCusto       *string `json:"nome_centro_custo,omitempty"`
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
	return []any{
		sql.Named("ano_mes_inicio", bi.NullString(in.AnoMesInicio)),
		sql.Named("ano_mes_fim", bi.NullString(in.AnoMesFim)),
		sql.Named("id_colaborador", bi.NullInt(in.IDColaborador)),
		sql.Named("id_filial", bi.NullInt(in.IDFilial)),
		sql.Named("id_secao", bi.NullInt(in.IDSecao)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_funcao_colaborador", bi.NullInt(in.IDFuncaoColaborador)),
		sql.Named("id_situacao_colaborador", bi.NullInt(in.IDSituacaoColaborador)),
		sql.Named("id_tipo_recebimento", bi.NullInt(in.IDTipoRecebimento)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("nome_colaborador", bi.NullString(in.NomeColaborador)),
		sql.Named("nome_empresa", bi.NullString(in.NomeEmpresa)),
		sql.Named("nome_filial", bi.NullString(in.NomeFilial)),
		sql.Named("nome_secao", bi.NullString(in.NomeSecao)),
		sql.Named("descricao_funcao", bi.NullString(in.DescricaoFuncao)),
		sql.Named("descricao_situacao", bi.NullString(in.DescricaoSituacao)),
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
