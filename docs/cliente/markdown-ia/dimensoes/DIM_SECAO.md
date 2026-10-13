# DIM_SECAO

**Tipo:** Dimensão

## Descrição

Dimensão de seções/departamentos onde os colaboradores estão lotados.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Seção` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Seção` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Nome Seção` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Registro Ativo Seção` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Seção",
    "Código Seção",
    "Nome Seção",
    "Seção",
    "Registro Ativo Seção"
FROM DIM_SECAO
WHERE (:id IS NULL OR "ID Seção" = :id)
  AND (:codigo IS NULL OR "Código Seção" = :codigo)
  AND (:descricao IS NULL OR UPPER("Nome Seção") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Seção" = :registro_ativo)
ORDER BY "Nome Seção"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Resolver seções citadas pelo usuário para filtrar o histórico de colaboradores e a folha por departamento.
