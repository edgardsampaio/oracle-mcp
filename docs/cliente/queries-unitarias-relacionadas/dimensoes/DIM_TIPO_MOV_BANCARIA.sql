-- DIM_TIPO_MOV_BANCARIA
-- Dimensão de tipos de movimentação bancária.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Tipo Mov. Bancária",
    "Código Tipo Mov. Bancária",
    "Descrição Tipo Mov. Bancária",
    "Tipo Mov. Bancária",
    "Reg. Ativo Tipo Mov. Bancária"
FROM DIM_TIPO_MOV_BANCARIA
WHERE (:id IS NULL OR "ID Tipo Mov. Bancária" = :id)
  AND (:codigo IS NULL OR "Código Tipo Mov. Bancária" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Mov. Bancária") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Tipo Mov. Bancária" = :registro_ativo)
ORDER BY "Descrição Tipo Mov. Bancária"
FETCH FIRST :limite ROWS ONLY;
