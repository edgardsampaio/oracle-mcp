-- DIM_EVENTO
-- Dimensão de eventos da folha de pagamento (proventos, descontos, bases) e seus tipos.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_evento (inteiro)
--   :codigo_evento (texto)
--   :nome_evento (texto)
--   :id_tipo_evento (inteiro)
--   :nome_tipo_evento (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Evento",
    "Código Evento",
    "Nome Evento",
    "Evento",
    "ID Tipo Evento",
    "Código Tipo Evento",
    "Nome Tipo Evento",
    "Tipo Evento",
    "Registro Ativo Evento",
    "Registro Ativo Tipo Evento"
FROM DIM_EVENTO
WHERE (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:codigo_evento IS NULL OR "Código Evento" = :codigo_evento)
  AND (:nome_evento IS NULL OR UPPER("Nome Evento") LIKE '%' || UPPER(:nome_evento) || '%')
  AND (:id_tipo_evento IS NULL OR "ID Tipo Evento" = :id_tipo_evento)
  AND (:nome_tipo_evento IS NULL OR UPPER("Nome Tipo Evento") LIKE '%' || UPPER(:nome_tipo_evento) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo Evento" = :registro_ativo)
ORDER BY "Nome Tipo Evento", "Nome Evento"
FETCH FIRST :limite ROWS ONLY;
