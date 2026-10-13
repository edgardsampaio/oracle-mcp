# DIM_CENTRO_CUSTO

**Tipo:** Dimensão

## Descrição

Dimensão hierárquica de centros de custo: Gestor → Segmento → Grupo de Contrato → Contrato → Centro de Custo, incluindo Cliente Final e Tipo de Contrato.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_gestor` | inteiro | Não | `ID Gestor` | igual | ID do gestor. |
| `nome_gestor` | texto | Não | `Nome Gestor` | contém | Busca parcial pelo nome do gestor. |
| `id_segmento` | inteiro | Não | `ID Segmento` | igual | ID do segmento. |
| `nome_segmento` | texto | Não | `Nome Segmento` | contém | Busca parcial pelo segmento. |
| `id_grupo_contrato` | inteiro | Não | `ID Grupo Contrato` | igual | ID do grupo de contrato. |
| `id_contrato` | inteiro | Não | `ID Contrato` | igual | ID do contrato. |
| `codigo_contrato` | texto | Não | `Código Contrato` | igual | Código do contrato. |
| `nome_contrato` | texto | Não | `Nome Contrato` | contém | Busca parcial pelo nome do contrato. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | ID do centro de custo. |
| `codigo_centro_custo` | texto | Não | `Código Centro de Custo` | igual | Código do centro de custo. |
| `nome_centro_custo` | texto | Não | `Nome Centro de Custo` | contém | Busca parcial pelo centro de custo. |
| `id_cliente_final` | inteiro | Não | `ID Cliente Final` | igual | ID do cliente final. |
| `nome_cliente_final` | texto | Não | `Nome Cliente Final` | contém | Busca parcial pelo cliente final. |
| `id_tipo_contrato` | inteiro | Não | `ID Tipo Contrato` | igual | ID do tipo de contrato. |
| `registro_ativo` | texto | Não | `Registro Ativo Centro de Custo` | igual | Flag de centro de custo ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Gestor",
    "Código Gestor",
    "Nome Gestor",
    "Gestor",
    "ID Segmento",
    "Código Segmento",
    "Nome Segmento",
    "Segmento",
    "ID Grupo Contrato",
    "Código Grupo Contrato",
    "Nome Grupo Contrato",
    "Grupo Contrato",
    "ID Contrato",
    "Código Contrato",
    "Nome Contrato",
    "Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo",
    "Nome Centro de Custo",
    "Centro de Custo",
    "ID Cliente Final",
    "Código Cliente Final",
    "Nome Cliente Final",
    "Cliente Final",
    "ID Tipo Contrato",
    "Código Tipo Contrato",
    "Nome Tipo Contrato",
    "Tipo Contrato",
    "Registro Ativo Gestor",
    "Registro Ativo Segmento",
    "Registro Ativo Grupo Contrato",
    "Registro Ativo Contrato",
    "Registro Ativo Centro de Custo",
    "Registro Ativo Cliente Final"
FROM DIM_CENTRO_CUSTO
WHERE (:id_gestor IS NULL OR "ID Gestor" = :id_gestor)
  AND (:nome_gestor IS NULL OR UPPER("Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:id_segmento IS NULL OR "ID Segmento" = :id_segmento)
  AND (:nome_segmento IS NULL OR UPPER("Nome Segmento") LIKE '%' || UPPER(:nome_segmento) || '%')
  AND (:id_grupo_contrato IS NULL OR "ID Grupo Contrato" = :id_grupo_contrato)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:codigo_contrato IS NULL OR "Código Contrato" = :codigo_contrato)
  AND (:nome_contrato IS NULL OR UPPER("Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:codigo_centro_custo IS NULL OR "Código Centro de Custo" = :codigo_centro_custo)
  AND (:nome_centro_custo IS NULL OR UPPER("Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:id_cliente_final IS NULL OR "ID Cliente Final" = :id_cliente_final)
  AND (:nome_cliente_final IS NULL OR UPPER("Nome Cliente Final") LIKE '%' || UPPER(:nome_cliente_final) || '%')
  AND (:id_tipo_contrato IS NULL OR "ID Tipo Contrato" = :id_tipo_contrato)
  AND (:registro_ativo IS NULL OR "Registro Ativo Centro de Custo" = :registro_ativo)
ORDER BY "Nome Gestor", "Nome Contrato", "Nome Centro de Custo"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Resolver nomes de gestores, segmentos, contratos e centros de custo citados em linguagem natural para os IDs usados nos fatos, e permitir agregações por qualquer nível da hierarquia.
