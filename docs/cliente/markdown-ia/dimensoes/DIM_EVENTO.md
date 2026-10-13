# DIM_EVENTO

**Tipo:** Dimensão

## Descrição

Dimensão de eventos da folha de pagamento (proventos, descontos, bases) e seus tipos.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_evento` | inteiro | Não | `ID Evento` | igual | ID do evento. |
| `codigo_evento` | texto | Não | `Código Evento` | igual | Código do evento. |
| `nome_evento` | texto | Não | `Nome Evento` | contém | Busca parcial pelo nome do evento. |
| `id_tipo_evento` | inteiro | Não | `ID Tipo Evento` | igual | ID do tipo de evento. |
| `nome_tipo_evento` | texto | Não | `Nome Tipo Evento` | contém | Busca parcial pelo tipo (provento, desconto…). |
| `registro_ativo` | texto | Não | `Registro Ativo Evento` | igual | Flag de evento ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Evento",
    "Código Evento",
    "Nome Evento",
    "Evento",
    "ID Tipo Evento",
    "Código Tipo Evento",
    "Nome Tipo Evento",
    "Tipo Evento",
    "Registro Ativo Evento",
    "Registro Ativo Tipo Evento"
FROM DIM_EVENTO
WHERE (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:codigo_evento IS NULL OR "Código Evento" = :codigo_evento)
  AND (:nome_evento IS NULL OR UPPER("Nome Evento") LIKE '%' || UPPER(:nome_evento) || '%')
  AND (:id_tipo_evento IS NULL OR "ID Tipo Evento" = :id_tipo_evento)
  AND (:nome_tipo_evento IS NULL OR UPPER("Nome Tipo Evento") LIKE '%' || UPPER(:nome_tipo_evento) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Evento" = :registro_ativo)
ORDER BY "Nome Tipo Evento", "Nome Evento"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Identificar eventos de folha (ex.: salário, horas extras, INSS) para filtrar e agrupar a FATO_FOLHA_PAGAMENTO.
