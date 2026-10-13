-- DIM_CONTA_ORCAMENTARIA
-- Dimensão do plano de contas orçamentárias em três níveis (N1, N2, N3), com nível, tipo, condição e natureza associada.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_conta_orcamentaria (inteiro)
--   :codigo_conta (texto)
--   :nome_conta (texto)
--   :codigo_n1 (texto)
--   :codigo_n2 (texto)
--   :codigo_n3 (texto)
--   :nivel (texto)
--   :tipo (texto)
--   :codigo_natureza (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

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
