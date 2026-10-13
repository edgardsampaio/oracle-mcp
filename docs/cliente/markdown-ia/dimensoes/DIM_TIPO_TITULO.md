# DIM_TIPO_TITULO

**Tipo:** Dimensão

## Descrição

Dimensão de tipos de título financeiro (ex.: NF, boleto, adiantamento, impostos).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Tipo Título` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Tipo Título` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Tipo Título` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Registro Ativo Tipo Título` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Tipo Título",
    "Código Tipo Título",
    "Descrição Tipo Título",
    "Tipo Título",
    "Registro Ativo Tipo Título"
FROM DIM_TIPO_TITULO
WHERE (:id IS NULL OR "ID Tipo Título" = :id)
  AND (:codigo IS NULL OR "Código Tipo Título" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Título") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Tipo Título" = :registro_ativo)
ORDER BY "Código Tipo Título"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Classificar títulos de Contas a Pagar, Contas a Receber e baixas por tipo — por exemplo, separar adiantamentos e impostos de títulos comerciais.
