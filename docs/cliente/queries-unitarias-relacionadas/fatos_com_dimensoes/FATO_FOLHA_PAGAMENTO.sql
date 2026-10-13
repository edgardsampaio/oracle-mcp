-- FATO_FOLHA_PAGAMENTO — Fato com Dimensões
-- Lançamentos da folha de pagamento enriquecidos com os dados do colaborador, o evento e seu tipo, e o calendário da competência.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :ano_mes_inicio (texto)
--   :ano_mes_fim (texto)
--   :data_competencia_inicio (data)
--   :data_competencia_fim (data)
--   :periodo_pagamento (texto)
--   :id_colaborador (inteiro)
--   :id_evento (inteiro)
--   :registro_ativo (texto)
--   :nome_colaborador (texto)
--   :matricula (texto)
--   :nome_evento (texto)
--   :nome_tipo_evento (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."Ano e Mês",
    f."Data Competência",
    f."Período Pagamento",
    f."ID Colaborador",
    f."ID Evento",
    f."Valor Evento",
    f."Quantidade Evento",
    f."Registro Ativo",
    colab."Matrícula Colaborador",
    colab."Nome Colaborador",
    colab."Data Admissão",
    colab."Data Desligamento",
    colab."Situação Atual",
    evt."Código Evento",
    evt."Nome Evento",
    evt."Evento",
    evt."Nome Tipo Evento",
    evt."Tipo Evento",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_FOLHA_PAGAMENTO f
LEFT JOIN DIM_COLABORADOR colab ON colab."ID Colaborador" = f."ID Colaborador"
LEFT JOIN DIM_EVENTO evt ON evt."ID Evento" = f."ID Evento"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Competência"
WHERE EXISTS (SELECT 1 FROM FATO_HISTORICO_COLABORADOR h JOIN DIM_RLS_CENTRO_CUSTO rls ON rls."ID Centro de Custo" = h."ID Centro de Custo" WHERE rls."E-mail" = :email AND h."ID Colaborador" = f."ID Colaborador" AND h."Ano e Mês" = f."Ano e Mês")
  AND (:ano_mes_inicio IS NULL OR f."Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR f."Ano e Mês" <= :ano_mes_fim)
  AND (:data_competencia_inicio IS NULL OR f."Data Competência" >= :data_competencia_inicio)
  AND (:data_competencia_fim IS NULL OR f."Data Competência" <= :data_competencia_fim)
  AND (:periodo_pagamento IS NULL OR f."Período Pagamento" = :periodo_pagamento)
  AND (:id_colaborador IS NULL OR f."ID Colaborador" = :id_colaborador)
  AND (:id_evento IS NULL OR f."ID Evento" = :id_evento)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_colaborador IS NULL OR UPPER(colab."Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:matricula IS NULL OR colab."Matrícula Colaborador" = :matricula)
  AND (:nome_evento IS NULL OR UPPER(evt."Nome Evento") LIKE '%' || UPPER(:nome_evento) || '%')
  AND (:nome_tipo_evento IS NULL OR UPPER(evt."Nome Tipo Evento") LIKE '%' || UPPER(:nome_tipo_evento) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Competência" DESC, f."ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
