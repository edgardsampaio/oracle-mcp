-- DIM_TIPO_PROCESSO
-- Dimensão de tipos de processo (ex.: trabalhista, cível, tributário).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Tipo Processo",
    "Código Tipo Processo",
    "Descrição Tipo Processo",
    "Tipo Processo",
    "Registro Ativo Tipo Processo"
FROM DIM_TIPO_PROCESSO
WHERE (:id IS NULL OR "ID Tipo Processo" = :id)
  AND (:codigo IS NULL OR "Código Tipo Processo" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Tipo Processo") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Tipo Processo" = :registro_ativo)
ORDER BY "Descrição Tipo Processo"
FETCH FIRST :limite ROWS ONLY;
