-- DIM_SECAO
-- Dimensão de seções/departamentos onde os colaboradores estão lotados.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Seção",
    "Código Seção",
    "Nome Seção",
    "Seção",
    "Registro Ativo Seção"
FROM DIM_SECAO
WHERE (:id IS NULL OR "ID Seção" = :id)
  AND (:codigo IS NULL OR "Código Seção" = :codigo)
  AND (:descricao IS NULL OR UPPER("Nome Seção") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Seção" = :registro_ativo)
ORDER BY "Nome Seção"
FETCH FIRST :limite ROWS ONLY;
