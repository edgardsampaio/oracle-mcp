# FATO_CONFISSAO_DIVIDA

**Tipo:** Fato

## Descrição

Fato de saldos de confissões de dívida por fornecedor e contrato, com a data do último pagamento.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_fornecedor` | inteiro | Não | `ID Fornecedor` | igual | Fornecedor. |
| `id_contrato` | inteiro | Não | `ID Contrato` | igual | Contrato. |
| `valor_saldo_minimo` | decimal | Não | `Valor Saldo` | >= | Saldo mínimo. |
| `data_ultimo_pagamento_inicio` | data | Não | `Data Último Pagamento` | >= | Último pagamento a partir de. |
| `data_ultimo_pagamento_fim` | data | Não | `Data Último Pagamento` | <= | Último pagamento até. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Fornecedor",
    "ID Contrato",
    "Valor Saldo",
    "Data Último Pagamento"
FROM FATO_CONFISSAO_DIVIDA
WHERE (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:valor_saldo_minimo IS NULL OR "Valor Saldo" >= :valor_saldo_minimo)
  AND (:data_ultimo_pagamento_inicio IS NULL OR "Data Último Pagamento" >= :data_ultimo_pagamento_inicio)
  AND (:data_ultimo_pagamento_fim IS NULL OR "Data Último Pagamento" <= :data_ultimo_pagamento_fim)
ORDER BY "Valor Saldo" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder quanto ainda se deve em confissões de dívida, a quem, em quais contratos e há quanto tempo não há pagamento.
