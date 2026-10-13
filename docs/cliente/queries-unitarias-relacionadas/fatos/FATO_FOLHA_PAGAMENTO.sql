-- FATO_FOLHA_PAGAMENTO
-- Fato de lançamentos da folha de pagamento por colaborador, competência e evento, com valor e quantidade.
--
-- Parâmetros (nulo = filtro ignorado):
--   :ano_mes_inicio (texto)
--   :ano_mes_fim (texto)
--   :data_competencia_inicio (data)
--   :data_competencia_fim (data)
--   :periodo_pagamento (texto)
--   :id_colaborador (inteiro)
--   :id_evento (inteiro)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "Ano e Mês",
    "Data Competência",
    "Período Pagamento",
    "ID Colaborador",
    "ID Evento",
    "Valor Evento",
    "Quantidade Evento",
    "Registro Ativo"
FROM FATO_FOLHA_PAGAMENTO
WHERE (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:data_competencia_inicio IS NULL OR "Data Competência" >= :data_competencia_inicio)
  AND (:data_competencia_fim IS NULL OR "Data Competência" <= :data_competencia_fim)
  AND (:periodo_pagamento IS NULL OR "Período Pagamento" = :periodo_pagamento)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_evento IS NULL OR "ID Evento" = :id_evento)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
