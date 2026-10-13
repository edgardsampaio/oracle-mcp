# DIM_FUNCAO_COLABORADOR

**Tipo:** Dimensão

## Descrição

Dimensão de funções/cargos dos colaboradores.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id` | inteiro | Não | `ID Função Colaborador` | igual | ID interno do registro. |
| `codigo` | texto | Não | `Código Função Colaborador` | igual | Código do registro no sistema de origem. |
| `descricao` | texto | Não | `Descrição Função Colaborador` | contém | Busca parcial (case-insensitive) pela descrição. |
| `registro_ativo` | texto | Não | `Reg. Ativo Função Colaborador` | igual | Filtra pela flag de registro ativo (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Função Colaborador",
    "Código Função Colaborador",
    "Descrição Função Colaborador",
    "Função Colaborador",
    "Reg. Ativo Função Colaborador"
FROM DIM_FUNCAO_COLABORADOR
WHERE (:id IS NULL OR "ID Função Colaborador" = :id)
  AND (:codigo IS NULL OR "Código Função Colaborador" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Função Colaborador") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Função Colaborador" = :registro_ativo)
ORDER BY "Descrição Função Colaborador"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Resolver funções citadas pelo usuário para filtrar o histórico de colaboradores (headcount e salários por função).
