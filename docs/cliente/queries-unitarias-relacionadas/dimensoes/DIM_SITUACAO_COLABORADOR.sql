-- DIM_SITUACAO_COLABORADOR
-- Dimensão de situações do colaborador (ativo, afastado, férias, demitido etc.).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Situação Colaborador",
    "Código Situação Colaborador",
    "Descrição Situação Colaborador",
    "Situação Colaborador",
    "Reg. Ativo Situação Colab."
FROM DIM_SITUACAO_COLABORADOR
WHERE (:id IS NULL OR "ID Situação Colaborador" = :id)
  AND (:codigo IS NULL OR "Código Situação Colaborador" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Situação Colaborador") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Situação Colab." = :registro_ativo)
ORDER BY "Descrição Situação Colaborador"
FETCH FIRST :limite ROWS ONLY;
