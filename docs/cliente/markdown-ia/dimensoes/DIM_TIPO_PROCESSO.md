# DIM_TIPO_PROCESSO

**Tipo:** Dimensão

## Descrição

Dimensão de tipos de processo (ex.: trabalhista, cível, tributário).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Tipo Processo` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Tipo Processo` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Tipo Processo` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Registro Ativo Tipo Processo` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Tipo Processo",
    "Código Tipo Processo",
    "Descrição Tipo Processo",
    "Tipo Processo",
    "Registro Ativo Tipo Processo"
FROM DIM_TIPO_PROCESSO
WHERE (:id IS NULL OR "ID Tipo Processo" = :id)
  AND (:codigo IS NULL OR "Código Tipo Processo" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Processo") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Tipo Processo" = :registro_ativo)
ORDER BY "Descrição Tipo Processo"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Classificar os processos da FATO_PROCESSO e os títulos vinculados a processos.
