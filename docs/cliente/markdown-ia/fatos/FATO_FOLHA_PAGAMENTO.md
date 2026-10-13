# FATO_FOLHA_PAGAMENTO

**Tipo:** Fato

## Descrição

Fato de lançamentos da folha de pagamento por colaborador, competência e evento, com valor e quantidade.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `ano_mes_inicio` | texto | Não | `Ano e Mês` | >= | Competência inicial (mesmo formato da coluna `Ano e Mês`). |
| `ano_mes_fim` | texto | Não | `Ano e Mês` | <= | Competência final. |
| `data_competencia_inicio` | data | Não | `Data Competência` | >= | Data de competência a partir de. |
| `data_competencia_fim` | data | Não | `Data Competência` | <= | Data de competência até. |
| `periodo_pagamento` | texto | Não | `Período Pagamento` | igual | Período de pagamento (ex.: mensal, 13º, férias). |
| `id_colaborador` | inteiro | Não | `ID Colaborador` | igual | Colaborador. |
| `id_evento` | inteiro | Não | `ID Evento` | igual | Evento de folha. |
| `registro_ativo` | texto | Não | `Registro Ativo` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
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
WHERE (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:data_competencia_inicio IS NULL OR "Data Competência" >= :data_competencia_inicio)
  AND (:data_competencia_fim IS NULL OR "Data Competência" <= :data_competencia_fim)
  AND (:periodo_pagamento IS NULL OR "Período Pagamento" = :periodo_pagamento)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas sobre custo de folha: total por competência, por evento (salário, horas extras, encargos) e por colaborador; combinada com DIM_COLABORADOR e FATO_HISTORICO_COLABORADOR para análises por filial, seção ou centro de custo.
