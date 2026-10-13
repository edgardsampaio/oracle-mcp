-- FATO_CONFISSAO_DIVIDA — Fato com Dimensões
-- Saldos de confissões de dívida enriquecidos com os dados do fornecedor, a hierarquia do contrato (gestor, segmento, grupo, tipo) e o calendário da data do último pagamento.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :id_fornecedor (inteiro)
--   :id_contrato (inteiro)
--   :valor_saldo_minimo (decimal)
--   :data_ultimo_pagamento_inicio (data)
--   :data_ultimo_pagamento_fim (data)
--   :razao_social_fornecedor (texto)
--   :cnpj_fornecedor (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."ID Fornecedor",
    f."ID Contrato",
    f."Valor Saldo",
    f."Data Último Pagamento",
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor",
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
FROM FATO_CONFISSAO_DIVIDA f
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Último Pagamento"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:valor_saldo_minimo IS NULL OR f."Valor Saldo" >= :valor_saldo_minimo)
  AND (:data_ultimo_pagamento_inicio IS NULL OR f."Data Último Pagamento" >= :data_ultimo_pagamento_inicio)
  AND (:data_ultimo_pagamento_fim IS NULL OR f."Data Último Pagamento" <= :data_ultimo_pagamento_fim)
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Valor Saldo" DESC
FETCH FIRST :limite ROWS ONLY;
