# DIM_MOEDA

**Tipo:** Dimensão

## Descrição

Dimensão de moedas utilizadas nos títulos e movimentações.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Moeda` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Moeda` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Nome Moeda` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Registro Ativo Moeda` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Moeda",
    "Código Moeda",
    "Nome Moeda",
    "Moeda",
    "Registro Ativo Moeda"
FROM DIM_MOEDA
WHERE (:id IS NULL OR "ID Moeda" = :id)
  AND (:codigo IS NULL OR "Código Moeda" = :codigo)
  AND (:descricao IS NULL OR UPPER("Nome Moeda") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Moeda" = :registro_ativo)
ORDER BY "Código Moeda"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Identificar a moeda de títulos e movimentos, evitando somar valores em moedas diferentes.
