-- DIM_CENTRO_CUSTO
-- Dimensão hierárquica de centros de custo: Gestor → Segmento → Grupo de Contrato → Contrato → Centro de Custo, incluindo Cliente Final e Tipo de Contrato.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_gestor (inteiro)
--   :nome_gestor (texto)
--   :id_segmento (inteiro)
--   :nome_segmento (texto)
--   :id_grupo_contrato (inteiro)
--   :id_contrato (inteiro)
--   :codigo_contrato (texto)
--   :nome_contrato (texto)
--   :id_centro_custo (inteiro)
--   :codigo_centro_custo (texto)
--   :nome_centro_custo (texto)
--   :id_cliente_final (inteiro)
--   :nome_cliente_final (texto)
--   :id_tipo_contrato (inteiro)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Gestor",
    "Código Gestor",
    "Nome Gestor",
    "Gestor",
    "ID Segmento",
    "Código Segmento",
    "Nome Segmento",
    "Segmento",
    "ID Grupo Contrato",
    "Código Grupo Contrato",
    "Nome Grupo Contrato",
    "Grupo Contrato",
    "ID Contrato",
    "Código Contrato",
    "Nome Contrato",
    "Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo",
    "Nome Centro de Custo",
    "Centro de Custo",
    "ID Cliente Final",
    "Código Cliente Final",
    "Nome Cliente Final",
    "Cliente Final",
    "ID Tipo Contrato",
    "Código Tipo Contrato",
    "Nome Tipo Contrato",
    "Tipo Contrato",
    "Registro Ativo Gestor",
    "Registro Ativo Segmento",
    "Registro Ativo Grupo Contrato",
    "Registro Ativo Contrato",
    "Registro Ativo Centro de Custo",
    "Registro Ativo Cliente Final"
FROM DIM_CENTRO_CUSTO
WHERE (:id_gestor IS NULL OR "ID Gestor" = :id_gestor)
  AND (:nome_gestor IS NULL OR UPPER("Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:id_segmento IS NULL OR "ID Segmento" = :id_segmento)
  AND (:nome_segmento IS NULL OR UPPER("Nome Segmento") LIKE '%' || UPPER(:nome_segmento) || '%')
  AND (:id_grupo_contrato IS NULL OR "ID Grupo Contrato" = :id_grupo_contrato)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:codigo_contrato IS NULL OR "Código Contrato" = :codigo_contrato)
  AND (:nome_contrato IS NULL OR UPPER("Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:codigo_centro_custo IS NULL OR "Código Centro de Custo" = :codigo_centro_custo)
  AND (:nome_centro_custo IS NULL OR UPPER("Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:id_cliente_final IS NULL OR "ID Cliente Final" = :id_cliente_final)
  AND (:nome_cliente_final IS NULL OR UPPER("Nome Cliente Final") LIKE '%' || UPPER(:nome_cliente_final) || '%')
  AND (:id_tipo_contrato IS NULL OR "ID Tipo Contrato" = :id_tipo_contrato)
  AND (:registro_ativo IS NULL OR "Registro Ativo Centro de Custo" = :registro_ativo)
ORDER BY "Nome Gestor", "Nome Contrato", "Nome Centro de Custo"
FETCH FIRST :limite ROWS ONLY;
