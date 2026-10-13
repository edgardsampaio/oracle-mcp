-- FATO_ORCAMENTO_TOTAL
-- Fato de orçamentos por versão, item, conta orçamentária e contrato, com valor orçado original, atual e a variação entre eles.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_orcamento (inteiro)
--   :codigo_orcamento (texto)
--   :numero_versao (texto)
--   :id_conta_orcamentaria (inteiro)
--   :id_contrato (inteiro)
--   :periodo_inicio (data)
--   :periodo_fim (data)
--   :limite (inteiro)

SELECT
    "ID Orçamento",
    "Código Orçamento",
    "Descrição Orçamento",
    "Orçamento",
    "Data Início Período",
    "Data Fim Período",
    "Número Versão",
    "Número Item",
    "ID Conta Orçamentária",
    "ID Contrato",
    "Valor Total Orçado Original",
    "Valor Total Orçado Atual",
    "Valor Variação Orig. x Atual"
FROM FATO_ORCAMENTO_TOTAL
WHERE (:id_orcamento IS NULL OR "ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR "Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR "Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR "Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR "Data Início Período" <= :periodo_fim)
ORDER BY "Data Início Período", "Número Item"
FETCH FIRST :limite ROWS ONLY;
