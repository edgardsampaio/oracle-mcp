# DIM_CONTA_ORCAMENTARIA

**Tipo:** Dimensão

## Descrição

Dimensão do plano de contas orçamentárias em três níveis (N1, N2, N3), com nível, tipo, condição e natureza associada.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_conta_orcamentaria` | inteiro | Não | `ID Conta Orçamentária` | igual | ID da conta. |
| `codigo_conta` | texto | Não | `Código Conta Orçamentária` | igual | Código da conta. |
| `nome_conta` | texto | Não | `Nome Conta Orçamentária` | contém | Busca parcial pelo nome da conta. |
| `codigo_n1` | texto | Não | `Código Conta Orçamentária N1` | igual | Código do nível 1. |
| `codigo_n2` | texto | Não | `Código Conta Orçamentária N2` | igual | Código do nível 2. |
| `codigo_n3` | texto | Não | `Código Conta Orçamentária N3` | igual | Código do nível 3. |
| `nivel` | texto | Não | `Nível` | igual | Nível da conta. |
| `tipo` | texto | Não | `Tipo` | igual | Tipo da conta (sintética/analítica). |
| `codigo_natureza` | texto | Não | `Código Natureza` | igual | Natureza financeira vinculada. |
| `registro_ativo` | texto | Não | `Reg. Ativo Conta Orçamentária` | igual | Flag de conta ativa. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Conta Orçamentária",
    "Código Conta Orçamentária",
    "Nome Conta Orçamentária",
    "Conta Orçamentária",
    "Código Conta Orçamentária N1",
    "Nome Conta Orçamentária N1",
    "Conta Orçamentária N1",
    "Código Conta Orçamentária N2",
    "Nome Conta Orçamentária N2",
    "Conta Orçamentária N2",
    "Código Conta Orçamentária N3",
    "Nome Conta Orçamentária N3",
    "Conta Orçamentária N3",
    "Nível",
    "Tipo",
    "Condição",
    "Código Natureza",
    "Reg. Ativo Conta Orçamentária"
FROM DIM_CONTA_ORCAMENTARIA
WHERE (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:codigo_conta IS NULL OR "Código Conta Orçamentária" = :codigo_conta)
  AND (:nome_conta IS NULL OR UPPER("Nome Conta Orçamentária") LIKE '%' || UPPER(:nome_conta) || '%')
  AND (:codigo_n1 IS NULL OR "Código Conta Orçamentária N1" = :codigo_n1)
  AND (:codigo_n2 IS NULL OR "Código Conta Orçamentária N2" = :codigo_n2)
  AND (:codigo_n3 IS NULL OR "Código Conta Orçamentária N3" = :codigo_n3)
  AND (:nivel IS NULL OR "Nível" = :nivel)
  AND (:tipo IS NULL OR "Tipo" = :tipo)
  AND (:codigo_natureza IS NULL OR "Código Natureza" = :codigo_natureza)
  AND (:registro_ativo IS NULL OR "Reg. Ativo Conta Orçamentária" = :registro_ativo)
ORDER BY "Código Conta Orçamentária"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Traduzir contas orçamentárias citadas pelo usuário para IDs e permitir análises de orçado por qualquer nível da hierarquia.
