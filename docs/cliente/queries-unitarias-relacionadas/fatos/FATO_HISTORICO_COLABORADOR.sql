-- FATO_HISTORICO_COLABORADOR
-- Fato mensal (snapshot por competência) da posição de cada colaborador: filial, seção, centro de custo, função, situação, forma de recebimento, banco e salário base.
--
-- Parâmetros (nulo = filtro ignorado):
--   :ano_mes_inicio (texto)
--   :ano_mes_fim (texto)
--   :id_colaborador (inteiro)
--   :id_filial (inteiro)
--   :id_secao (inteiro)
--   :id_centro_custo (inteiro)
--   :id_funcao_colaborador (inteiro)
--   :id_situacao_colaborador (inteiro)
--   :id_tipo_recebimento (inteiro)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "Ano e Mês",
    "Data Competência",
    "ID Colaborador",
    "ID Filial",
    "ID Seção",
    "ID Centro de Custo",
    "ID Função Colaborador",
    "ID Situação Colaborador",
    "ID Tipo Recebimento",
    "ID Banco",
    "Valor Salário Base",
    "Registro Ativo Colaborador"
FROM FATO_HISTORICO_COLABORADOR
WHERE (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_secao IS NULL OR "ID Seção" = :id_secao)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_funcao_colaborador IS NULL OR "ID Função Colaborador" = :id_funcao_colaborador)
  AND (:id_situacao_colaborador IS NULL OR "ID Situação Colaborador" = :id_situacao_colaborador)
  AND (:id_tipo_recebimento IS NULL OR "ID Tipo Recebimento" = :id_tipo_recebimento)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
