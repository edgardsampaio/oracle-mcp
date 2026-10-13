// Package colaborador é a dimensão DIM_COLABORADOR.
// Não importa fato nem outra dimensão.
package colaborador

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta dimensão.
	ToolName = "dim_colaborador"

	description = `Dimensão de colaboradores com matrícula, nome, datas de nascimento/admissão/transferência/desligamento e a posição atual (empresa, filial, seção, função e situação). Contém dados pessoais sujeitos à LGPD.

Identificar colaboradores e sua situação atual para análises de folha de pagamento, headcount, admissões e desligamentos.

Filtros opcionais: id_colaborador, matricula, nome_colaborador, empresa_atual, filial_atual, secao_atual, funcao_atual, situacao_atual, data_admissao_inicio, data_admissao_fim, data_desligamento_inicio, data_desligamento_fim, registro_ativo, limite. O limite omitido vale 500 e não passa do teto do servidor. Esta tool só lê.`

	searchSQL = `
SELECT
    "ID Colaborador",
    "Matrícula Colaborador",
    "Nome Colaborador",
    "Data Nascimento",
    "Data Admissão",
    "Data Transferência",
    "Data Desligamento",
    "Empresa Atual",
    "Filial Atual",
    "Seção Atual",
    "Função Atual",
    "Situação Atual",
    "Registro Ativo Colaborador"
FROM DIM_COLABORADOR
WHERE (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:matricula IS NULL OR "Matrícula Colaborador" = :matricula)
  AND (:nome_colaborador IS NULL OR UPPER("Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:empresa_atual IS NULL OR "Empresa Atual" = :empresa_atual)
  AND (:filial_atual IS NULL OR "Filial Atual" = :filial_atual)
  AND (:secao_atual IS NULL OR UPPER("Seção Atual") LIKE '%' || UPPER(:secao_atual) || '%')
  AND (:funcao_atual IS NULL OR UPPER("Função Atual") LIKE '%' || UPPER(:funcao_atual) || '%')
  AND (:situacao_atual IS NULL OR "Situação Atual" = :situacao_atual)
  AND (:data_admissao_inicio IS NULL OR "Data Admissão" >= :data_admissao_inicio)
  AND (:data_admissao_fim IS NULL OR "Data Admissão" <= :data_admissao_fim)
  AND (:data_desligamento_inicio IS NULL OR "Data Desligamento" >= :data_desligamento_inicio)
  AND (:data_desligamento_fim IS NULL OR "Data Desligamento" <= :data_desligamento_fim)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Nome Colaborador"
FETCH FIRST :limite ROWS ONLY`
)

type entrada struct {
	IDColaborador          *int    `json:"id_colaborador,omitempty"`
	Matricula              *string `json:"matricula,omitempty"`
	NomeColaborador        *string `json:"nome_colaborador,omitempty"`
	EmpresaAtual           *string `json:"empresa_atual,omitempty"`
	FilialAtual            *string `json:"filial_atual,omitempty"`
	SecaoAtual             *string `json:"secao_atual,omitempty"`
	FuncaoAtual            *string `json:"funcao_atual,omitempty"`
	SituacaoAtual          *string `json:"situacao_atual,omitempty"`
	DataAdmissaoInicio     *string `json:"data_admissao_inicio,omitempty"`
	DataAdmissaoFim        *string `json:"data_admissao_fim,omitempty"`
	DataDesligamentoInicio *string `json:"data_desligamento_inicio,omitempty"`
	DataDesligamentoFim    *string `json:"data_desligamento_fim,omitempty"`
	RegistroAtivo          *string `json:"registro_ativo,omitempty"`
	Limite                 *int    `json:"limite,omitempty"`
}

func (in entrada) binds(maxRows int) ([]any, error) {
	limite, err := bi.Limite(in.Limite, maxRows)
	if err != nil {
		return nil, err
	}

	data_admissao_inicio, err := bi.NullDate(in.DataAdmissaoInicio)
	if err != nil {
		return nil, err
	}
	data_admissao_fim, err := bi.NullDate(in.DataAdmissaoFim)
	if err != nil {
		return nil, err
	}
	data_desligamento_inicio, err := bi.NullDate(in.DataDesligamentoInicio)
	if err != nil {
		return nil, err
	}
	data_desligamento_fim, err := bi.NullDate(in.DataDesligamentoFim)
	if err != nil {
		return nil, err
	}
	return []any{
		sql.Named("id_colaborador", bi.NullInt(in.IDColaborador)),
		sql.Named("matricula", bi.NullString(in.Matricula)),
		sql.Named("nome_colaborador", bi.NullString(in.NomeColaborador)),
		sql.Named("empresa_atual", bi.NullString(in.EmpresaAtual)),
		sql.Named("filial_atual", bi.NullString(in.FilialAtual)),
		sql.Named("secao_atual", bi.NullString(in.SecaoAtual)),
		sql.Named("funcao_atual", bi.NullString(in.FuncaoAtual)),
		sql.Named("situacao_atual", bi.NullString(in.SituacaoAtual)),
		sql.Named("data_admissao_inicio", data_admissao_inicio),
		sql.Named("data_admissao_fim", data_admissao_fim),
		sql.Named("data_desligamento_inicio", data_desligamento_inicio),
		sql.Named("data_desligamento_fim", data_desligamento_fim),
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
