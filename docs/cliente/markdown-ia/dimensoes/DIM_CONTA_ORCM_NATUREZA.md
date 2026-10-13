# DIM_CONTA_ORCM_NATUREZA

**Tipo:** Dimensão

## Descrição

Tabela ponte (de-para) entre Contas Orçamentárias e Naturezas Financeiras.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_conta_orcamentaria` | inteiro | Não | `ID Conta Orçamentária` | igual | ID da conta orçamentária. |
| `id_natureza` | inteiro | Não | `ID Natureza` | igual | ID da natureza financeira. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Conta Orçamentária",
    "ID Natureza"
FROM DIM_CONTA_ORCM_NATUREZA
WHERE (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
ORDER BY "ID Conta Orçamentária", "ID Natureza"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Permitir cruzar orçado (por conta orçamentária) com realizado (por natureza financeira em Contas a Pagar/Receber e movimentações bancárias).
