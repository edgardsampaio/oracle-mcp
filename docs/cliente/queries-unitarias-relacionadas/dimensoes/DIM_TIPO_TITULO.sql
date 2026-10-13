-- DIM_TIPO_TITULO
-- Dimensão de tipos de título financeiro (ex.: NF, boleto, adiantamento, impostos).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Tipo Título",
    "Código Tipo Título",
    "Descrição Tipo Título",
    "Tipo Título",
    "Registro Ativo Tipo Título"
FROM DIM_TIPO_TITULO
WHERE (:id IS NULL OR "ID Tipo Título" = :id)
  AND (:codigo IS NULL OR "Código Tipo Título" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Título") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Tipo Título" = :registro_ativo)
ORDER BY "Código Tipo Título"
FETCH FIRST :limite ROWS ONLY;
