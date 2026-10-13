-- FATO_CONFISSAO_DIVIDA
-- Fato de saldos de confissões de dívida por fornecedor e contrato, com a data do último pagamento.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_fornecedor (inteiro)
--   :id_contrato (inteiro)
--   :valor_saldo_minimo (decimal)
--   :data_ultimo_pagamento_inicio (data)
--   :data_ultimo_pagamento_fim (data)
--   :limite (inteiro)

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
