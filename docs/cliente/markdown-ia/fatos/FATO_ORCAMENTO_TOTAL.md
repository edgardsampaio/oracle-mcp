# FATO_ORCAMENTO_TOTAL

**Tipo:** Fato

## Descrição

Fato de orçamentos por versão, item, conta orçamentária e contrato, com valor orçado original, atual e a variação entre eles.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_orcamento` | inteiro | Não | `ID Orçamento` | igual | Orçamento. |
| `codigo_orcamento` | texto | Não | `Código Orçamento` | igual | Código do orçamento. |
| `numero_versao` | texto | Não | `Número Versão` | igual | Versão do orçamento. |
| `id_conta_orcamentaria` | inteiro | Não | `ID Conta Orçamentária` | igual | Conta orçamentária. |
| `id_contrato` | inteiro | Não | `ID Contrato` | igual | Contrato. |
| `periodo_inicio` | data | Não | — | regra própria | Início do intervalo de interesse (retorna orçamentos cujo período se sobrepõe). |
| `periodo_fim` | data | Não | — | regra própria | Fim do intervalo de interesse. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
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
WHERE (:id_orcamento IS NULL OR "ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR "Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR "Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR "Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR "Data Início Período" <= :periodo_fim)
ORDER BY "Data Início Período", "Número Item"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de orçamento: quanto foi orçado por conta, contrato ou período, revisões entre versões e, via DIM_CONTA_ORCM_NATUREZA, comparações orçado x realizado.
