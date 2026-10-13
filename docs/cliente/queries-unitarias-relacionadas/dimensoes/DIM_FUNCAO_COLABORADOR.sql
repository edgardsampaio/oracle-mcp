-- DIM_FUNCAO_COLABORADOR
-- Dimensão de funções/cargos dos colaboradores.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id (inteiro)
--   :codigo (texto)
--   :descricao (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Função Colaborador",
    "Código Função Colaborador",
    "Descrição Função Colaborador",
    "Função Colaborador",
    "Reg. Ativo Função Colaborador"
FROM DIM_FUNCAO_COLABORADOR
WHERE (:id IS NULL OR "ID Função Colaborador" = :id)
  AND (:codigo IS NULL OR "Código Função Colaborador" = :codigo)
  AND (:descricao IS NULL OR UPPER("Descrição Função Colaborador") LIKE '%' || UPPER(:descricao) || '%')
  AND (:registro_ativo IS NULL OR "Reg. Ativo Função Colaborador" = :registro_ativo)
ORDER BY "Descrição Função Colaborador"
FETCH FIRST :limite ROWS ONLY;
