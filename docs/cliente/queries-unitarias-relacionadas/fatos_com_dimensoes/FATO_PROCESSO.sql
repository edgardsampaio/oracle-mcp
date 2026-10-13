-- FATO_PROCESSO — Fato com Dimensões
-- Processos enriquecidos com o tipo de processo, centro de custo/contrato e fornecedor.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :id_processo (inteiro)
--   :codigo_processo (texto)
--   :nome_autor (texto)
--   :id_tipo_processo (inteiro)
--   :id_centro_custo (inteiro)
--   :id_fornecedor (inteiro)
--   :valor_minimo (decimal)
--   :registro_ativo (texto)
--   :descricao_tipo_processo (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :nome_centro_custo (texto)
--   :razao_social_fornecedor (texto)
--   :cnpj_fornecedor (texto)
--   :limite (inteiro)

SELECT
    f."ID Processo",
    f."Código Processo",
    f."Nome Autor",
    f."ID Tipo Processo",
    f."ID Centro de Custo",
    f."ID Fornecedor",
    f."Valor Total",
    f."Registro Ativo",
    tproc."Descrição Tipo Processo",
    tproc."Tipo Processo",
    ccu."Nome Gestor",
    ccu."Gestor",
    ccu."Nome Segmento",
    ccu."Segmento",
    ccu."Nome Grupo Contrato",
    ccu."Grupo Contrato",
    ccu."Código Contrato",
    ccu."Nome Contrato",
    ccu."Contrato",
    ccu."Código Centro de Custo",
    ccu."Nome Centro de Custo",
    ccu."Centro de Custo",
    ccu."Nome Cliente Final",
    ccu."Cliente Final",
    ccu."Nome Tipo Contrato",
    ccu."Tipo Contrato",
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor"
FROM FATO_PROCESSO f
LEFT JOIN DIM_TIPO_PROCESSO tproc ON tproc."ID Tipo Processo" = f."ID Tipo Processo"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:id_processo IS NULL OR f."ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR f."Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER(f."Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR f."ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR f."Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:descricao_tipo_processo IS NULL OR UPPER(tproc."Descrição Tipo Processo") LIKE '%' || UPPER(:descricao_tipo_processo) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
ORDER BY f."Valor Total" DESC
FETCH FIRST :limite ROWS ONLY;
