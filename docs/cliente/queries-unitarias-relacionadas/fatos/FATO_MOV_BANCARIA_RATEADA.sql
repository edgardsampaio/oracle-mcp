-- FATO_MOV_BANCARIA_RATEADA
-- Fato de movimentações bancárias (entradas, saídas e transferências) rateadas por natureza e centro de custo, com conta bancária, beneficiário, histórico e dados do título baixado.
--
-- Parâmetros (nulo = filtro ignorado):
--   :data_inicio (data)
--   :data_fim (data)
--   :operacao (texto)
--   :id_filial (inteiro)
--   :id_conta_bancaria (inteiro)
--   :id_tipo_mov_bancaria (inteiro)
--   :id_natureza (inteiro)
--   :id_centro_custo (inteiro)
--   :id_fornecedor (inteiro)
--   :id_cliente (inteiro)
--   :beneficiario (texto)
--   :historico (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "Operação Mov. Bancária",
    "ID Filial",
    "Data Mov. Bancária",
    "ID Tipo Mov. Bancária",
    "Código Numerário",
    "Número Documento",
    "Situação Documento",
    "Histórico",
    "Beneficiário",
    "ID Moeda",
    "Valor Mov. Bancária",
    "ID Conta Bancária",
    "Banco Virtual",
    "ID Natureza",
    "ID Centro de Custo",
    "Rateio Financeiro",
    "Rateio Contábil",
    "ID Empresa Filial Origem",
    "ID Empresa Filial Destino",
    "ID Conta Bancária Origem",
    "ID Conta Bancária Destino",
    "Prefixo Título",
    "Número Título",
    "Parcela Título",
    "ID Tipo Título",
    "ID Fornecedor Título",
    "ID Município Fornecedor Título",
    "ID País Fornecedor Título",
    "ID Cliente Título",
    "ID Município Cliente Título",
    "Data Emissão Título",
    "Data Vencimento Título",
    "Número Lote",
    "Motivo Baixa",
    "Sequência Baixa",
    "Número Cheque",
    "ID Processo",
    "Número Registro",
    "Registro Ativo"
FROM FATO_MOV_BANCARIA_RATEADA
WHERE (:data_inicio IS NULL OR "Data Mov. Bancária" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data Mov. Bancária" <= :data_fim)
  AND (:operacao IS NULL OR "Operação Mov. Bancária" = :operacao)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR "ID Tipo Mov. Bancária" = :id_tipo_mov_bancaria)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor Título" = :id_fornecedor)
  AND (:id_cliente IS NULL OR "ID Cliente Título" = :id_cliente)
  AND (:beneficiario IS NULL OR UPPER("Beneficiário") LIKE '%' || UPPER(:beneficiario) || '%')
  AND (:historico IS NULL OR UPPER("Histórico") LIKE '%' || UPPER(:historico) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Mov. Bancária" DESC
FETCH FIRST :limite ROWS ONLY;
