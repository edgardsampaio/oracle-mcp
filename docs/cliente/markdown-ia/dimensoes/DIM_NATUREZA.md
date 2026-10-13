# DIM_NATUREZA

**Tipo:** Dimensão

## Descrição

Dimensão de naturezas financeiras em hierarquia de até quatro níveis (N1 a N4) mais o nível analítico.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_natureza` | inteiro | Não | `ID Natureza` | igual | ID da natureza analítica. |
| `codigo_natureza` | texto | Não | `Código Natureza` | igual | Código da natureza analítica. |
| `nome_natureza` | texto | Não | `Nome Natureza` | contém | Busca parcial pela natureza analítica. |
| `id_natureza_n1` | inteiro | Não | `ID Natureza N1` | igual | ID do nível 1. |
| `nome_natureza_n1` | texto | Não | `Nome Natureza N1` | contém | Busca parcial no nível 1. |
| `id_natureza_n2` | inteiro | Não | `ID Natureza N2` | igual | ID do nível 2. |
| `nome_natureza_n2` | texto | Não | `Nome Natureza N2` | contém | Busca parcial no nível 2. |
| `id_natureza_n3` | inteiro | Não | `ID Natureza N3` | igual | ID do nível 3. |
| `id_natureza_n4` | inteiro | Não | `ID Natureza N4` | igual | ID do nível 4. |
| `registro_ativo` | texto | Não | `Registro Ativo Natureza` | igual | Flag de natureza ativa. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Natureza N1",
    "Código Natureza N1",
    "Nome Natureza N1",
    "Natureza N1",
    "ID Natureza N2",
    "Código Natureza N2",
    "Nome Natureza N2",
    "Natureza N2",
    "ID Natureza N3",
    "Código Natureza N3",
    "Nome Natureza N3",
    "Natureza N3",
    "ID Natureza N4",
    "Código Natureza N4",
    "Nome Natureza N4",
    "Natureza N4",
    "ID Natureza",
    "Código Natureza",
    "Nome Natureza",
    "Natureza",
    "Registro Bloqueado Natureza",
    "Registro Ativo Natureza"
FROM DIM_NATUREZA
WHERE (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:codigo_natureza IS NULL OR "Código Natureza" = :codigo_natureza)
  AND (:nome_natureza IS NULL OR UPPER("Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:id_natureza_n1 IS NULL OR "ID Natureza N1" = :id_natureza_n1)
  AND (:nome_natureza_n1 IS NULL OR UPPER("Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:id_natureza_n2 IS NULL OR "ID Natureza N2" = :id_natureza_n2)
  AND (:nome_natureza_n2 IS NULL OR UPPER("Nome Natureza N2") LIKE '%' || UPPER(:nome_natureza_n2) || '%')
  AND (:id_natureza_n3 IS NULL OR "ID Natureza N3" = :id_natureza_n3)
  AND (:id_natureza_n4 IS NULL OR "ID Natureza N4" = :id_natureza_n4)
  AND (:registro_ativo IS NULL OR "Registro Ativo Natureza" = :registro_ativo)
ORDER BY "Código Natureza"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Classificar receitas e despesas: permite à IA traduzir termos como 'despesas com pessoal' ou 'impostos' para naturezas e agregar o financeiro em qualquer nível.
