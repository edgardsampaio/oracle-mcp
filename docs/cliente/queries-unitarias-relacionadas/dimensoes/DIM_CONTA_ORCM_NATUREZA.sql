-- DIM_CONTA_ORCM_NATUREZA
-- Tabela ponte (de-para) entre Contas Orçamentárias e Naturezas Financeiras.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_conta_orcamentaria (inteiro)
--   :id_natureza (inteiro)
--   :limite (inteiro)

SELECT
    "ID Conta Orçamentária",
    "ID Natureza"
FROM DIM_CONTA_ORCM_NATUREZA
WHERE (:id_conta_orcamentaria IS NULL OR "ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
ORDER BY "ID Conta Orçamentária", "ID Natureza"
FETCH FIRST :limite ROWS ONLY;
