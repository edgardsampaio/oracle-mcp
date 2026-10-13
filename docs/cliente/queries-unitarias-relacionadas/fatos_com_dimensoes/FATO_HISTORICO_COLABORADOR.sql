-- FATO_HISTORICO_COLABORADOR — Fato com Dimensões
-- Snapshot mensal dos colaboradores enriquecido com dados do colaborador, filial, seção, centro de custo/contrato, função, situação, forma de recebimento, banco e calendário da competência.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
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
--   :nome_colaborador (texto)
--   :nome_empresa (texto)
--   :nome_filial (texto)
--   :nome_secao (texto)
--   :descricao_funcao (texto)
--   :descricao_situacao (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :nome_centro_custo (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."Ano e Mês",
    f."Data Competência",
    f."ID Colaborador",
    f."ID Filial",
    f."ID Seção",
    f."ID Centro de Custo",
    f."ID Função Colaborador",
    f."ID Situação Colaborador",
    f."ID Tipo Recebimento",
    f."ID Banco",
    f."Valor Salário Base",
    f."Registro Ativo Colaborador",
    colab."Matrícula Colaborador",
    colab."Nome Colaborador",
    colab."Data Admissão",
    colab."Data Desligamento",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    sec."Código Seção",
    sec."Nome Seção",
    sec."Seção",
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
    fun."Descrição Função Colaborador",
    fun."Função Colaborador",
    sit."Descrição Situação Colaborador",
    sit."Situação Colaborador",
    trec."Descrição Tipo Recebimento",
    trec."Tipo Recebimento",
    bco."Código Banco",
    bco."Nome Banco",
    bco."Banco",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_HISTORICO_COLABORADOR f
LEFT JOIN DIM_COLABORADOR colab ON colab."ID Colaborador" = f."ID Colaborador"
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_SECAO sec ON sec."ID Seção" = f."ID Seção"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_FUNCAO_COLABORADOR fun ON fun."ID Função Colaborador" = f."ID Função Colaborador"
LEFT JOIN DIM_SITUACAO_COLABORADOR sit ON sit."ID Situação Colaborador" = f."ID Situação Colaborador"
LEFT JOIN DIM_TIPO_RECEBIMENTO_COLAB trec ON trec."ID Tipo Recebimento" = f."ID Tipo Recebimento"
LEFT JOIN (SELECT DISTINCT "ID Banco", "Código Banco", "Nome Banco", "Banco" FROM DIM_BANCO_CONTA_BANCARIA) bco ON bco."ID Banco" = f."ID Banco"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Competência"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:ano_mes_inicio IS NULL OR f."Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR f."Ano e Mês" <= :ano_mes_fim)
  AND (:id_colaborador IS NULL OR f."ID Colaborador" = :id_colaborador)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_secao IS NULL OR f."ID Seção" = :id_secao)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_funcao_colaborador IS NULL OR f."ID Função Colaborador" = :id_funcao_colaborador)
  AND (:id_situacao_colaborador IS NULL OR f."ID Situação Colaborador" = :id_situacao_colaborador)
  AND (:id_tipo_recebimento IS NULL OR f."ID Tipo Recebimento" = :id_tipo_recebimento)
  AND (:registro_ativo IS NULL OR f."Registro Ativo Colaborador" = :registro_ativo)
  AND (:nome_colaborador IS NULL OR UPPER(colab."Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:nome_secao IS NULL OR UPPER(sec."Nome Seção") LIKE '%' || UPPER(:nome_secao) || '%')
  AND (:descricao_funcao IS NULL OR UPPER(fun."Descrição Função Colaborador") LIKE '%' || UPPER(:descricao_funcao) || '%')
  AND (:descricao_situacao IS NULL OR UPPER(sit."Descrição Situação Colaborador") LIKE '%' || UPPER(:descricao_situacao) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Competência" DESC, f."ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
