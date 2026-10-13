# DIM_RLS_CENTRO_CUSTO

**Tipo:** Dimensão

## Descrição

Tabela de segurança em nível de linha (RLS): relaciona o e-mail do usuário aos contratos e centros de custo que ele pode visualizar.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `E-mail` | igual | E-mail do usuário autenticado (injetado pelo servidor). |
| `id_contrato` | inteiro | Não | `ID Contrato` | igual | ID do contrato. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | ID do centro de custo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "E-mail",
    "ID Contrato",
    "Código Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo"
FROM DIM_RLS_CENTRO_CUSTO
WHERE (:email IS NULL OR "E-mail" = :email)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
ORDER BY "Código Contrato", "Código Centro de Custo"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Restringir os dados retornados por qualquer consulta aos contratos/centros de custo autorizados ao usuário. O parâmetro `email` deve ser injetado pelo servidor MCP a partir do usuário autenticado — nunca informado pela IA ou pelo usuário no chat.
