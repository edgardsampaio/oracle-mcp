-- DIM_MOEDA
-- Dimensão de moedas utilizadas nos títulos e movimentações.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Moeda",
    "Código Moeda",
    "Nome Moeda",
    "Moeda",
    "Registro Ativo Moeda"
FROM DIM_MOEDA
WHERE (:id IS NULL OR "ID Moeda" = :id)
  AND (:codigo IS NULL OR "Código Moeda" = :codigo)
  AND (:descricao IS NULL OR UPPER("Nome Moeda") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Moeda" = :registro_ativo)
ORDER BY "Código Moeda"
FETCH FIRST :limite ROWS ONLY;
