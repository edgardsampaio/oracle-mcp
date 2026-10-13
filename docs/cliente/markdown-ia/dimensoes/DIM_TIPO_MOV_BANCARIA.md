# DIM_TIPO_MOV_BANCARIA

**Tipo:** Dimensão

## Descrição

Dimensão de tipos de movimentação bancária.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Tipo Mov. Bancária` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Tipo Mov. Bancária` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Tipo Mov. Bancária` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Reg. Ativo Tipo Mov. Bancária` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Tipo Mov. Bancária",
    "Código Tipo Mov. Bancária",
    "Descrição Tipo Mov. Bancária",
    "Tipo Mov. Bancária",
    "Reg. Ativo Tipo Mov. Bancária"
FROM DIM_TIPO_MOV_BANCARIA
WHERE (:id IS NULL OR "ID Tipo Mov. Bancária" = :id)
  AND (:codigo IS NULL OR "Código Tipo Mov. Bancária" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Mov. Bancária") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Tipo Mov. Bancária" = :registro_ativo)
ORDER BY "Descrição Tipo Mov. Bancária"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Classificar movimentações bancárias por tipo (ex.: transferência, tarifa, aplicação) na FATO_MOV_BANCARIA_RATEADA e na FATO_BAIXA_TITULO_IMPOSTO.
