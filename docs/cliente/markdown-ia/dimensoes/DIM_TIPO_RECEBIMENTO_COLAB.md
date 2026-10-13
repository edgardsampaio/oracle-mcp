# DIM_TIPO_RECEBIMENTO_COLAB

**Tipo:** Dimensão

## Descrição

Dimensão de formas de recebimento do salário pelo colaborador (ex.: crédito em conta, cheque, PIX).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Tipo Recebimento` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Tipo Recebimento` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Tipo Recebimento` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Reg. Ativo Tipo Recebimento` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Tipo Recebimento",
    "Código Tipo Recebimento",
    "Descrição Tipo Recebimento",
    "Tipo Recebimento",
    "Reg. Ativo Tipo Recebimento"
FROM DIM_TIPO_RECEBIMENTO_COLAB
WHERE (:id IS NULL OR "ID Tipo Recebimento" = :id)
  AND (:codigo IS NULL OR "Código Tipo Recebimento" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Recebimento") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Tipo Recebimento" = :registro_ativo)
ORDER BY "Descrição Tipo Recebimento"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Filtrar e agrupar o histórico de colaboradores pela forma de recebimento.
