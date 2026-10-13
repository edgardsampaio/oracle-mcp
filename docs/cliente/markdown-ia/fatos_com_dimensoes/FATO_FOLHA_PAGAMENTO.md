# FATO_FOLHA_PAGAMENTO — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_FOLHA_PAGAMENTO`

## Descrição

Lançamentos da folha de pagamento enriquecidos com os dados do colaborador, o evento e seu tipo, e o calendário da competência.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada via `FATO_HISTORICO_COLABORADOR` (centro de custo do colaborador na mesma competência), já que a folha não tem `ID Centro de Custo`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `colab` | DIM_COLABORADOR | `ID Colaborador` | `ID Colaborador` | — |
| `evt` | DIM_EVENTO | `ID Evento` | `ID Evento` | — |
| `t` | DIM_TEMPO | `Data Competência` | `Data` | calendário da `Data Competência` |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `ano_mes_inicio` | texto | Não | `f."Ano e Mês"` | >= | Competência inicial (mesmo formato da coluna `Ano e Mês`). |
| `ano_mes_fim` | texto | Não | `f."Ano e Mês"` | <= | Competência final. |
| `data_competencia_inicio` | data | Não | `f."Data Competência"` | >= | Data de competência a partir de. |
| `data_competencia_fim` | data | Não | `f."Data Competência"` | <= | Data de competência até. |
| `periodo_pagamento` | texto | Não | `f."Período Pagamento"` | igual | Período de pagamento (ex.: mensal, 13º, férias). |
| `id_colaborador` | inteiro | Não | `f."ID Colaborador"` | igual | Colaborador. |
| `id_evento` | inteiro | Não | `f."ID Evento"` | igual | Evento de folha. |
| `registro_ativo` | texto | Não | `f."Registro Ativo"` | igual | Flag de registro ativo. |
| `nome_colaborador` | texto | Não | `colab."Nome Colaborador"` | contém | Busca parcial pelo colaborador. |
| `matricula` | texto | Não | `colab."Matrícula Colaborador"` | igual | Matrícula do colaborador. |
| `nome_evento` | texto | Não | `evt."Nome Evento"` | contém | Busca parcial pelo evento. |
| `nome_tipo_evento` | texto | Não | `evt."Nome Tipo Evento"` | contém | Busca parcial pelo tipo de evento (provento, desconto…). |
| `ano` | inteiro | Não | `t."Ano"` | igual | Ano da data de referência. |
| `mes` | inteiro | Não | `t."Mês"` | igual | Mês (1–12) da data de referência. |
| `trimestre` | inteiro | Não | `t."Trimestre"` | igual | Trimestre (1–4) da data de referência. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
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
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de custo de folha com rótulos legíveis — total por evento, tipo de evento, colaborador e competência — já restrito aos colaboradores lotados em centros de custo autorizados ao usuário.
