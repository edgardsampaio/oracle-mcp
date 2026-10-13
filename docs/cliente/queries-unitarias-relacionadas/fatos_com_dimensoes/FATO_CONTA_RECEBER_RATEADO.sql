-- FATO_CONTA_RECEBER_RATEADO — Fato com Dimensões
-- Títulos de Contas a Receber rateados, enriquecidos com filial, tipo de título, cliente, natureza (N1–N4), moeda, centro de custo/contrato e calendário do vencimento real.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :data_vencimento_inicio (data)
--   :data_vencimento_fim (data)
--   :data_emissao_inicio (data)
--   :data_emissao_fim (data)
--   :data_baixa_inicio (data)
--   :data_baixa_fim (data)
--   :id_filial (inteiro)
--   :id_cliente (inteiro)
--   :id_natureza (inteiro)
--   :id_centro_custo (inteiro)
--   :id_tipo_titulo (inteiro)
--   :numero (texto)
--   :somente_em_aberto (texto)
--   :registro_ativo (texto)
--   :nome_empresa (texto)
--   :nome_filial (texto)
--   :razao_social_cliente (texto)
--   :cnpj_cliente (texto)
--   :nome_natureza (texto)
--   :nome_natureza_n1 (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :nome_centro_custo (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."ID Filial",
    f."Prefixo",
    f."Número",
    f."Parcela",
    f."ID Tipo Título",
    f."Código Tipo Título",
    f."ID Cliente",
    f."Código Cliente",
    f."Código Loja Cliente",
    f."ID Município",
    f."ID Natureza",
    f."ID Moeda",
    f."ID Centro de Custo",
    f."Data Emissão",
    f."Data Vencimento",
    f."Data Vencimento Real",
    f."Data Última Baixa",
    f."Histórico",
    f."Taxa Moeda",
    f."Valor Título",
    f."Valor Título Real",
    f."Valor Última Baixa",
    f."Valor Saldo",
    f."Valor Saldo Real",
    f."Valor Recebido Real",
    f."Rateio",
    f."Registro Ativo",
    f."Valor ISS",
    f."Valor INSS",
    f."Valor IRRF",
    f."Valor PIS",
    f."Valor COFINS",
    f."Valor CSLL",
    f."Valor Saldo Acréscimo",
    f."Valor Saldo Decréscimo",
    f."Chave Título Pai",
    f."Valor Total Tributos Reter",
    f."Valor Líquido Título Real",
    f."Valor Abatimentos",
    f."Valor Tributos Retidos Baixa",
    f."Valor Líquido Recebido Real",
    f."Valor Líquido Saldo Real",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    ttit."Descrição Tipo Título",
    ttit."Tipo Título",
    cli."Cliente Razão Social",
    cli."Cliente Nome Fantasia",
    cli."Cliente CNPJ Razão Social",
    cli."CNPJ" AS "CNPJ Cliente",
    cli."Tipo Pessoa" AS "Tipo Pessoa Cliente",
    cli."Município" AS "Município Cliente",
    cli."Sigla UF" AS "UF Cliente",
    cli."Região" AS "Região Cliente",
    nat."Natureza N1",
    nat."Natureza N2",
    nat."Natureza N3",
    nat."Natureza N4",
    nat."Código Natureza",
    nat."Nome Natureza",
    nat."Natureza",
    moe."Moeda",
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
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_CONTA_RECEBER_RATEADO f
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_TIPO_TITULO ttit ON ttit."ID Tipo Título" = f."ID Tipo Título"
LEFT JOIN DIM_CLIENTE cli ON cli."ID Cliente" = f."ID Cliente"
LEFT JOIN DIM_NATUREZA nat ON nat."ID Natureza" = f."ID Natureza"
LEFT JOIN DIM_MOEDA moe ON moe."ID Moeda" = f."ID Moeda"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Vencimento Real"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:data_vencimento_inicio IS NULL OR f."Data Vencimento Real" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR f."Data Vencimento Real" <= :data_vencimento_fim)
  AND (:data_emissao_inicio IS NULL OR f."Data Emissão" >= :data_emissao_inicio)
  AND (:data_emissao_fim IS NULL OR f."Data Emissão" <= :data_emissao_fim)
  AND (:data_baixa_inicio IS NULL OR f."Data Última Baixa" >= :data_baixa_inicio)
  AND (:data_baixa_fim IS NULL OR f."Data Última Baixa" <= :data_baixa_fim)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_cliente IS NULL OR f."ID Cliente" = :id_cliente)
  AND (:id_natureza IS NULL OR f."ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_tipo_titulo IS NULL OR f."ID Tipo Título" = :id_tipo_titulo)
  AND (:numero IS NULL OR f."Número" = :numero)
  AND (:somente_em_aberto IS NULL OR :somente_em_aberto <> 'S' OR f."Valor Saldo Real" > 0)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:razao_social_cliente IS NULL OR UPPER(cli."Razão Social") LIKE '%' || UPPER(:razao_social_cliente) || '%')
  AND (:cnpj_cliente IS NULL OR cli."CNPJ" = :cnpj_cliente)
  AND (:nome_natureza IS NULL OR UPPER(nat."Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:nome_natureza_n1 IS NULL OR UPPER(nat."Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Vencimento Real"
FETCH FIRST :limite ROWS ONLY;
