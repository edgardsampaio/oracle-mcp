-- FATO_BAIXA_TITULO_IMPOSTO — Fato com Dimensões
-- Baixas (pagamentos) de títulos de impostos enriquecidas com filial, tipo de movimentação, moeda, conta bancária, natureza, centro de custo/contrato, tipo de título, órgão/fornecedor do título, filiais e contas de origem/destino e calendário.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :data_inicio (data)
--   :data_fim (data)
--   :id_filial (inteiro)
--   :id_conta_bancaria (inteiro)
--   :id_tipo_mov_bancaria (inteiro)
--   :operacao (texto)
--   :id_natureza (inteiro)
--   :id_centro_custo (inteiro)
--   :id_fornecedor (inteiro)
--   :id_tipo_titulo (inteiro)
--   :numero_titulo (texto)
--   :data_vencimento_inicio (data)
--   :data_vencimento_fim (data)
--   :registro_ativo (texto)
--   :nome_empresa (texto)
--   :nome_filial (texto)
--   :nome_natureza (texto)
--   :nome_natureza_n1 (texto)
--   :nome_gestor (texto)
--   :nome_contrato (texto)
--   :nome_centro_custo (texto)
--   :razao_social_fornecedor (texto)
--   :cnpj_fornecedor (texto)
--   :nome_banco (texto)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :limite (inteiro)

SELECT
    f."Operação Mov. Bancária Bx. Tit. Imposto",
    f."ID Filial",
    f."Data Mov. Bancária Bx. Tit. Imposto",
    f."ID Tipo Mov. Bancária Bx. Tit. Imposto",
    f."Código Numerário",
    f."Número Documento",
    f."Situação Documento",
    f."Histórico",
    f."Beneficiário",
    f."ID Moeda",
    f."Valor Mov. Bancária Bx. Tit. Imposto",
    f."ID Conta Bancária",
    f."Banco Virtual",
    f."ID Natureza",
    f."ID Centro de Custo",
    f."Rateio Financeiro",
    f."Rateio Contábil",
    f."ID Empresa Filial Origem",
    f."ID Empresa Filial Destino",
    f."ID Conta Bancária Origem",
    f."ID Conta Bancária Destino",
    f."Prefixo Título",
    f."Número Título",
    f."Parcela Título",
    f."ID Tipo Título",
    f."ID Fornecedor Título",
    f."ID Município Fornecedor Título",
    f."ID País Fornecedor Título",
    f."ID Cliente Título",
    f."ID Município Cliente Título",
    f."Data Emissão Título",
    f."Data Vencimento Título",
    f."Número Lote",
    f."Motivo Baixa",
    f."Sequência Baixa",
    f."Número Cheque",
    f."ID Processo",
    f."Valor Título",
    f."Registro Ativo",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    tmov."Descrição Tipo Mov. Bancária",
    tmov."Tipo Mov. Bancária",
    moe."Moeda",
    cb."Banco",
    cb."Código Agência",
    cb."Número Conta Bancária",
    cb."Nome Conta Bancária",
    cb."Tipo Conta Bancária",
    cb."Conta Bancária",
    cb."Banco, Agência e Conta Bancária",
    nat."Natureza N1",
    nat."Natureza N2",
    nat."Natureza N3",
    nat."Natureza N4",
    nat."Código Natureza",
    nat."Nome Natureza",
    nat."Natureza",
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
    filo."Empresa" AS "Empresa Origem",
    filo."Filial" AS "Filial Origem",
    fild."Empresa" AS "Empresa Destino",
    fild."Filial" AS "Filial Destino",
    cbo."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Origem",
    cbd."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Destino",
    ttit."Descrição Tipo Título",
    ttit."Tipo Título",
    forn."Fornecedor Razão Social" AS "Fornecedor Razão Social Título",
    forn."Fornecedor Nome Fantasia" AS "Fornecedor Nome Fantasia Título",
    forn."Fornecedor CNPJ Razão Social" AS "Fornecedor CNPJ Razão Social Título",
    forn."CNPJ" AS "CNPJ Fornecedor Título",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor Título",
    forn."Município" AS "Município Fornecedor Título",
    forn."Sigla UF" AS "UF Fornecedor Título",
    forn."País" AS "País Fornecedor Título",
    cli."Cliente Razão Social" AS "Cliente Razão Social Título",
    cli."Cliente Nome Fantasia" AS "Cliente Nome Fantasia Título",
    cli."Cliente CNPJ Razão Social" AS "Cliente CNPJ Razão Social Título",
    cli."CNPJ" AS "CNPJ Cliente Título",
    cli."Tipo Pessoa" AS "Tipo Pessoa Cliente Título",
    cli."Município" AS "Município Cliente Título",
    cli."Sigla UF" AS "UF Cliente Título",
    cli."Região" AS "Região Cliente Título",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_BAIXA_TITULO_IMPOSTO f
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_TIPO_MOV_BANCARIA tmov ON tmov."ID Tipo Mov. Bancária" = f."ID Tipo Mov. Bancária Bx. Tit. Imposto"
LEFT JOIN DIM_MOEDA moe ON moe."ID Moeda" = f."ID Moeda"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cb ON cb."ID Conta Bancária" = f."ID Conta Bancária"
LEFT JOIN DIM_NATUREZA nat ON nat."ID Natureza" = f."ID Natureza"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_EMPRESA_FILIAL filo ON filo."ID Filial" = f."ID Empresa Filial Origem"
LEFT JOIN DIM_EMPRESA_FILIAL fild ON fild."ID Filial" = f."ID Empresa Filial Destino"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbo ON cbo."ID Conta Bancária" = f."ID Conta Bancária Origem"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbd ON cbd."ID Conta Bancária" = f."ID Conta Bancária Destino"
LEFT JOIN DIM_TIPO_TITULO ttit ON ttit."ID Tipo Título" = f."ID Tipo Título"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor Título"
LEFT JOIN DIM_CLIENTE cli ON cli."ID Cliente" = f."ID Cliente Título"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Mov. Bancária Bx. Tit. Imposto"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:data_inicio IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" >= :data_inicio)
  AND (:data_fim IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" <= :data_fim)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR f."ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR f."ID Tipo Mov. Bancária Bx. Tit. Imposto" = :id_tipo_mov_bancaria)
  AND (:operacao IS NULL OR f."Operação Mov. Bancária Bx. Tit. Imposto" = :operacao)
  AND (:id_natureza IS NULL OR f."ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor Título" = :id_fornecedor)
  AND (:id_tipo_titulo IS NULL OR f."ID Tipo Título" = :id_tipo_titulo)
  AND (:numero_titulo IS NULL OR f."Número Título" = :numero_titulo)
  AND (:data_vencimento_inicio IS NULL OR f."Data Vencimento Título" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR f."Data Vencimento Título" <= :data_vencimento_fim)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:nome_natureza IS NULL OR UPPER(nat."Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:nome_natureza_n1 IS NULL OR UPPER(nat."Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_banco IS NULL OR UPPER(cb."Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Mov. Bancária Bx. Tit. Imposto" DESC
FETCH FIRST :limite ROWS ONLY;
