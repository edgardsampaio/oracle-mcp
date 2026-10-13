// Package processo é o fato FATO_PROCESSO.
// Não importa dimensão. O e-mail da segurança por linha vem do usuário autenticado.
package processo

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o fato sem rótulos de dimensão.
	ToolName = "fato_processo"
	// ToolNameComDimensoes é o mesmo fato com os rótulos já unidos.
	ToolNameComDimensoes = "fato_processo_com_dimensoes"

	description = `Fato de processos (judiciais/administrativos) com autor, tipo, centro de custo, fornecedor e valor total.

Responder perguntas sobre passivo de processos: quantidade e valor por tipo, centro de custo, fornecedor ou autor. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo 'ID Centro de Custo'.

Use quando os ids já forem conhecidos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_processo, codigo_processo, nome_autor, id_tipo_processo, id_centro_custo, id_fornecedor, valor_minimo, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	descriptionCom = `Processos enriquecidos com o tipo de processo, centro de custo/contrato e fornecedor.

Todas as dimensões entram por 'LEFT JOIN', para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo 'ID Centro de Custo'.

Responder perguntas sobre passivo de processos com rótulos legíveis — quantidade e valor por tipo, contrato, gestor ou fornecedor — já restrito aos centros de custo autorizados ao usuário.

Use quando a pergunta citar nomes, períodos ou rótulos. O e-mail da segurança por linha é o do usuário autenticado e não faz parte dos parâmetros. Filtros opcionais: id_processo, codigo_processo, nome_autor, id_tipo_processo, id_centro_custo, id_fornecedor, valor_minimo, registro_ativo, descricao_tipo_processo, nome_gestor, nome_contrato, nome_centro_custo, razao_social_fornecedor, cnpj_fornecedor, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	sqlFato = `
SELECT
    "ID Processo",
    "Código Processo",
    "Nome Autor",
    "ID Tipo Processo",
    "ID Centro de Custo",
    "ID Fornecedor",
    "Valor Total",
    "Registro Ativo"
FROM FATO_PROCESSO
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = FATO_PROCESSO."ID Centro de Custo")
  AND (:id_processo IS NULL OR "ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR "Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER("Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR "ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR "Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Valor Total" DESC
FETCH FIRST :limite ROWS ONLY`

	sqlComDimensoes = `
SELECT
    f."ID Processo",
    f."Código Processo",
    f."Nome Autor",
    f."ID Tipo Processo",
    f."ID Centro de Custo",
    f."ID Fornecedor",
    f."Valor Total",
    f."Registro Ativo",
    tproc."Descrição Tipo Processo",
    tproc."Tipo Processo",
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
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor"
FROM FATO_PROCESSO f
LEFT JOIN DIM_TIPO_PROCESSO tproc ON tproc."ID Tipo Processo" = f."ID Tipo Processo"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:id_processo IS NULL OR f."ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR f."Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER(f."Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR f."ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR f."Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:descricao_tipo_processo IS NULL OR UPPER(tproc."Descrição Tipo Processo") LIKE '%' || UPPER(:descricao_tipo_processo) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
ORDER BY f."Valor Total" DESC
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDProcesso     *int     `json:"id_processo,omitempty"`
	CodigoProcesso *string  `json:"codigo_processo,omitempty"`
	NomeAutor      *string  `json:"nome_autor,omitempty"`
	IDTipoProcesso *int     `json:"id_tipo_processo,omitempty"`
	IDCentroCusto  *int     `json:"id_centro_custo,omitempty"`
	IDFornecedor   *int     `json:"id_fornecedor,omitempty"`
	ValorMinimo    *float64 `json:"valor_minimo,omitempty"`
	RegistroAtivo  *string  `json:"registro_ativo,omitempty"`
	Limite         *int     `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_processo", bi.NullInt(in.IDProcesso)),
		sql.Named("codigo_processo", bi.NullString(in.CodigoProcesso)),
		sql.Named("nome_autor", bi.NullString(in.NomeAutor)),
		sql.Named("id_tipo_processo", bi.NullInt(in.IDTipoProcesso)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("valor_minimo", bi.NullFloat(in.ValorMinimo)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("limite", limite),
	}, nil
}

type entradaCom struct {
	IDProcesso            *int     `json:"id_processo,omitempty"`
	CodigoProcesso        *string  `json:"codigo_processo,omitempty"`
	NomeAutor             *string  `json:"nome_autor,omitempty"`
	IDTipoProcesso        *int     `json:"id_tipo_processo,omitempty"`
	IDCentroCusto         *int     `json:"id_centro_custo,omitempty"`
	IDFornecedor          *int     `json:"id_fornecedor,omitempty"`
	ValorMinimo           *float64 `json:"valor_minimo,omitempty"`
	RegistroAtivo         *string  `json:"registro_ativo,omitempty"`
	DescricaoTipoProcesso *string  `json:"descricao_tipo_processo,omitempty"`
	NomeGestor            *string  `json:"nome_gestor,omitempty"`
	NomeContrato          *string  `json:"nome_contrato,omitempty"`
	NomeCentroCusto       *string  `json:"nome_centro_custo,omitempty"`
	RazaoSocialFornecedor *string  `json:"razao_social_fornecedor,omitempty"`
	CNPJFornecedor        *string  `json:"cnpj_fornecedor,omitempty"`
	Limite                *int     `json:"limite,omitempty"`
}

func (in entradaCom) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_processo", bi.NullInt(in.IDProcesso)),
		sql.Named("codigo_processo", bi.NullString(in.CodigoProcesso)),
		sql.Named("nome_autor", bi.NullString(in.NomeAutor)),
		sql.Named("id_tipo_processo", bi.NullInt(in.IDTipoProcesso)),
		sql.Named("id_centro_custo", bi.NullInt(in.IDCentroCusto)),
		sql.Named("id_fornecedor", bi.NullInt(in.IDFornecedor)),
		sql.Named("valor_minimo", bi.NullFloat(in.ValorMinimo)),
		sql.Named("registro_ativo", bi.NullString(in.RegistroAtivo)),
		sql.Named("descricao_tipo_processo", bi.NullString(in.DescricaoTipoProcesso)),
		sql.Named("nome_gestor", bi.NullString(in.NomeGestor)),
		sql.Named("nome_contrato", bi.NullString(in.NomeContrato)),
		sql.Named("nome_centro_custo", bi.NullString(in.NomeCentroCusto)),
		sql.Named("razao_social_fornecedor", bi.NullString(in.RazaoSocialFornecedor)),
		sql.Named("cnpj_fornecedor", bi.NullString(in.CNPJFornecedor)),
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
