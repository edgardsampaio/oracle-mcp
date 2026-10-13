# DIM_SITUACAO_COLABORADOR

**Tipo:** Dimensão

## Descrição

Dimensão de situações do colaborador (ativo, afastado, férias, demitido etc.).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Situação Colaborador` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Situação Colaborador` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Situação Colaborador` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Reg. Ativo Situação Colab.` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Situação Colaborador",
    "Código Situação Colaborador",
    "Descrição Situação Colaborador",
    "Situação Colaborador",
    "Reg. Ativo Situação Colab."
FROM DIM_SITUACAO_COLABORADOR
WHERE (:id IS NULL OR "ID Situação Colaborador" = :id)
  AND (:codigo IS NULL OR "Código Situação Colaborador" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Situação Colaborador") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Situação Colab." = :registro_ativo)
ORDER BY "Descrição Situação Colaborador"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Filtrar o histórico de colaboradores por situação, viabilizando contagem de ativos, afastados e desligados.
