-- DIM_TIPO_RECEBIMENTO_COLAB
-- Dimensão de formas de recebimento do salário pelo colaborador (ex.: crédito em conta, cheque, PIX).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Tipo Recebimento",
    "Código Tipo Recebimento",
    "Descrição Tipo Recebimento",
    "Tipo Recebimento",
    "Reg. Ativo Tipo Recebimento"
FROM DIM_TIPO_RECEBIMENTO_COLAB
WHERE (:id IS NULL OR "ID Tipo Recebimento" = :id)
  AND (:codigo IS NULL OR "Código Tipo Recebimento" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Recebimento") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Tipo Recebimento" = :registro_ativo)
ORDER BY "Descrição Tipo Recebimento"
FETCH FIRST :limite ROWS ONLY;
