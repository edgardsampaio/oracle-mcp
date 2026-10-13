# FATO_HISTORICO_COLABORADOR — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_HISTORICO_COLABORADOR`

## Descrição

Snapshot mensal dos colaboradores enriquecido com dados do colaborador, filial, seção, centro de custo/contrato, função, situação, forma de recebimento, banco e calendário da competência.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo `ID Centro de Custo`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `colab` | DIM_COLABORADOR | `ID Colaborador` | `ID Colaborador` | — |
| `fil` | DIM_EMPRESA_FILIAL | `ID Filial` | `ID Filial` | — |
| `sec` | DIM_SECAO | `ID Seção` | `ID Seção` | — |
| `ccu` | DIM_CENTRO_CUSTO | `ID Centro de Custo` | `ID Centro de Custo` | — |
| `fun` | DIM_FUNCAO_COLABORADOR | `ID Função Colaborador` | `ID Função Colaborador` | — |
| `sit` | DIM_SITUACAO_COLABORADOR | `ID Situação Colaborador` | `ID Situação Colaborador` | — |
| `trec` | DIM_TIPO_RECEBIMENTO_COLAB | `ID Tipo Recebimento` | `ID Tipo Recebimento` | — |
| `bco` | DIM_BANCO_CONTA_BANCARIA (distinta por banco) | `ID Banco` | `ID Banco` | banco (lista distinta de bancos da DIM_BANCO_CONTA_BANCARIA) |
| `t` | DIM_TEMPO | `Data Competência` | `Data` | calendário da `Data Competência` |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `ano_mes_inicio` | texto | Não | `f."Ano e Mês"` | >= | Competência inicial (mesmo formato da coluna `Ano e Mês`). |
| `ano_mes_fim` | texto | Não | `f."Ano e Mês"` | <= | Competência final. |
| `id_colaborador` | inteiro | Não | `f."ID Colaborador"` | igual | Colaborador. |
| `id_filial` | inteiro | Não | `f."ID Filial"` | igual | Filial. |
| `id_secao` | inteiro | Não | `f."ID Seção"` | igual | Seção. |
| `id_centro_custo` | inteiro | Não | `f."ID Centro de Custo"` | igual | Centro de custo. |
| `id_funcao_colaborador` | inteiro | Não | `f."ID Função Colaborador"` | igual | Função. |
| `id_situacao_colaborador` | inteiro | Não | `f."ID Situação Colaborador"` | igual | Situação. |
| `id_tipo_recebimento` | inteiro | Não | `f."ID Tipo Recebimento"` | igual | Forma de recebimento. |
| `registro_ativo` | texto | Não | `f."Registro Ativo Colaborador"` | igual | Flag de registro ativo. |
| `nome_colaborador` | texto | Não | `colab."Nome Colaborador"` | contém | Busca parcial pelo colaborador. |
| `nome_empresa` | texto | Não | `fil."Nome Empresa"` | contém | Busca parcial pela empresa. |
| `nome_filial` | texto | Não | `fil."Nome Filial"` | contém | Busca parcial pela filial. |
| `nome_secao` | texto | Não | `sec."Nome Seção"` | contém | Busca parcial pela seção. |
| `descricao_funcao` | texto | Não | `fun."Descrição Função Colaborador"` | contém | Busca parcial pela função. |
| `descricao_situacao` | texto | Não | `sit."Descrição Situação Colaborador"` | contém | Busca parcial pela situação. |
| `nome_gestor` | texto | Não | `ccu."Nome Gestor"` | contém | Busca parcial pelo gestor. |
| `nome_contrato` | texto | Não | `ccu."Nome Contrato"` | contém | Busca parcial pelo contrato. |
| `nome_centro_custo` | texto | Não | `ccu."Nome Centro de Custo"` | contém | Busca parcial pelo centro de custo. |
| `ano` | inteiro | Não | `t."Ano"` | igual | Ano da data de referência. |
| `mes` | inteiro | Não | `t."Mês"` | igual | Mês (1–12) da data de referência. |
| `trimestre` | inteiro | Não | `t."Trimestre"` | igual | Trimestre (1–4) da data de referência. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
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
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de headcount e quadro com rótulos legíveis — ativos por mês, por função, seção, filial, contrato ou gestor, e evolução do salário base — já restrito aos centros de custo autorizados ao usuário.
