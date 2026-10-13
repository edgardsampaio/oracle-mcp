# FATO_PROCESSO

**Tipo:** Fato

## Descrição

Fato de processos (judiciais/administrativos) com autor, tipo, centro de custo, fornecedor e valor total.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_processo` | inteiro | Não | `ID Processo` | igual | Processo. |
| `codigo_processo` | texto | Não | `Código Processo` | igual | Código/número do processo. |
| `nome_autor` | texto | Não | `Nome Autor` | contém | Busca parcial pelo autor. |
| `id_tipo_processo` | inteiro | Não | `ID Tipo Processo` | igual | Tipo de processo. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | Centro de custo. |
| `id_fornecedor` | inteiro | Não | `ID Fornecedor` | igual | Fornecedor. |
| `valor_minimo` | decimal | Não | `Valor Total` | >= | Valor total mínimo. |
| `registro_ativo` | texto | Não | `Registro Ativo` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Processo",
    "Código Processo",
    "Nome Autor",
    "ID Tipo Processo",
    "ID Centro de Custo",
    "ID Fornecedor",
    "Valor Total",
    "Registro Ativo"
FROM FATO_PROCESSO
WHERE (:id_processo IS NULL OR "ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR "Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER("Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR "ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR "Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Valor Total" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas sobre passivo de processos: quantidade e valor por tipo, centro de custo, fornecedor ou autor. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo `ID Centro de Custo`.
