-- FATO_ORCAMENTO_TOTAL — Fato com Dimensões
-- Orçamentos por versão e item enriquecidos com a conta orçamentária (N1–N3), a hierarquia do contrato e o calendário do início do período.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :id_orcamento (inteiro)
--   :codigo_orcamento (texto)
--   :numero_versao (texto)
--   :id_conta_orcamentaria (inteiro)
--   :id_contrato (inteiro)
--   :periodo_inicio (data)
--   :periodo_fim (data)
--   :nome_conta_orcamentaria (texto)
--   :nome_conta_n1 (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."ID Orçamento",
    f."Código Orçamento",
    f."Descrição Orçamento",
    f."Orçamento",
    f."Data Início Período",
    f."Data Fim Período",
    f."Número Versão",
    f."Número Item",
    f."ID Conta Orçamentária",
    f."ID Contrato",
    f."Valor Total Orçado Original",
    f."Valor Total Orçado Atual",
    f."Valor Variação Orig. x Atual",
    corc."Código Conta Orçamentária",
    corc."Nome Conta Orçamentária",
    corc."Conta Orçamentária",
    corc."Conta Orçamentária N1",
    corc."Conta Orçamentária N2",
    corc."Conta Orçamentária N3",
    ctr."Nome Gestor",
    ctr."Gestor",
    ctr."Nome Segmento",
    ctr."Segmento",
    ctr."Nome Grupo Contrato",
    ctr."Grupo Contrato",
    ctr."Código Contrato",
    ctr."Nome Contrato",
    ctr."Contrato",
    ctr."Nome Tipo Contrato",
    ctr."Tipo Contrato",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_ORCAMENTO_TOTAL f
LEFT JOIN DIM_CONTA_ORCAMENTARIA corc ON corc."ID Conta Orçamentária" = f."ID Conta Orçamentária"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Início Período"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_orcamento IS NULL OR f."ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR f."Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR f."Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR f."ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR f."Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR f."Data Início Período" <= :periodo_fim)
  AND (:nome_conta_orcamentaria IS NULL OR UPPER(corc."Nome Conta Orçamentária") LIKE '%' || UPPER(:nome_conta_orcamentaria) || '%')
  AND (:nome_conta_n1 IS NULL OR UPPER(corc."Nome Conta Orçamentária N1") LIKE '%' || UPPER(:nome_conta_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Início Período", f."Número Item"
FETCH FIRST :limite ROWS ONLY;
