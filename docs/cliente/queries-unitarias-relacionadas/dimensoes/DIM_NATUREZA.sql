-- DIM_NATUREZA
-- Dimensão de naturezas financeiras em hierarquia de até quatro níveis (N1 a N4) mais o nível analítico.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_natureza (inteiro)
--   :codigo_natureza (texto)
--   :nome_natureza (texto)
--   :id_natureza_n1 (inteiro)
--   :nome_natureza_n1 (texto)
--   :id_natureza_n2 (inteiro)
--   :nome_natureza_n2 (texto)
--   :id_natureza_n3 (inteiro)
--   :id_natureza_n4 (inteiro)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Natureza N1",
    "Código Natureza N1",
    "Nome Natureza N1",
    "Natureza N1",
    "ID Natureza N2",
    "Código Natureza N2",
    "Nome Natureza N2",
    "Natureza N2",
    "ID Natureza N3",
    "Código Natureza N3",
    "Nome Natureza N3",
    "Natureza N3",
    "ID Natureza N4",
    "Código Natureza N4",
    "Nome Natureza N4",
    "Natureza N4",
    "ID Natureza",
    "Código Natureza",
    "Nome Natureza",
    "Natureza",
    "Registro Bloqueado Natureza",
    "Registro Ativo Natureza"
FROM DIM_NATUREZA
WHERE (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:codigo_natureza IS NULL OR "Código Natureza" = :codigo_natureza)
  AND (:nome_natureza IS NULL OR UPPER("Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:id_natureza_n1 IS NULL OR "ID Natureza N1" = :id_natureza_n1)
  AND (:nome_natureza_n1 IS NULL OR UPPER("Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:id_natureza_n2 IS NULL OR "ID Natureza N2" = :id_natureza_n2)
  AND (:nome_natureza_n2 IS NULL OR UPPER("Nome Natureza N2") LIKE '%' || UPPER(:nome_natureza_n2) || '%')
  AND (:id_natureza_n3 IS NULL OR "ID Natureza N3" = :id_natureza_n3)
  AND (:id_natureza_n4 IS NULL OR "ID Natureza N4" = :id_natureza_n4)
  AND (:registro_ativo IS NULL OR "Registro Ativo Natureza" = :registro_ativo)
ORDER BY "Código Natureza"
FETCH FIRST :limite ROWS ONLY;
