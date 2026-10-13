-- FATO_CONTA_PAGAR_RATEADO
-- Fato de títulos de Contas a Pagar rateados por natureza e centro de custo, com datas (emissão, vencimento, última baixa), valores originais, pagos, saldos, tributos retidos, acréscimos/decréscimos e indicadores de protesto, confissão de dívida, processo e dação.
--
-- Parâmetros (nulo = filtro ignorado):
--   :data_vencimento_inicio (data)
--   :data_vencimento_fim (data)
--   :data_emissao_inicio (data)
--   :data_emissao_fim (data)
--   :data_baixa_inicio (data)
--   :data_baixa_fim (data)
--   :id_filial (inteiro)
--   :id_fornecedor (inteiro)
--   :id_natureza (inteiro)
--   :id_centro_custo (inteiro)
--   :id_tipo_titulo (inteiro)
--   :chave_titulo (texto)
--   :numero (texto)
--   :somente_em_aberto (texto)
--   :confissao_divida (texto)
--   :protesto (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "Chave Título",
    "ID Filial",
    "Prefixo",
    "Número",
    "Parcela",
    "ID Tipo Título",
    "Código Tipo Título",
    "ID Fornecedor",
    "Código Fornecedor",
    "Código Loja Fornecedor",
    "ID Município",
    "ID País",
    "ID Natureza",
    "ID Moeda",
    "ID Centro de Custo",
    "Data Emissão",
    "Data Vencimento",
    "Data Vencimento Real",
    "Data Última Baixa",
    "Histórico",
    "Taxa Moeda",
    "Valor Título",
    "Valor Título Real",
    "Valor Última Baixa",
    "Valor Saldo",
    "Valor Saldo Real",
    "Valor Pago Real",
    "Rateio",
    "Protesto",
    "Valor Protesto",
    "Confissão Dívida",
    "Processo",
    "ID Processo",
    "Dação",
    "Adiantamento Cheque",
    "Adiantamento Compensado",
    "Registro Ativo",
    "Valor ISS",
    "Valor INSS",
    "Valor IRRF",
    "Valor PIS",
    "Valor COFINS",
    "Valor CSLL",
    "Valor Saldo Acréscimo",
    "Valor Saldo Decréscimo",
    "Chave Título Pai",
    "Valor Multa",
    "Valor Juros",
    "Valor Desconto",
    "Deduz INSS",
    "Valor SEST",
    "Valor Base",
    "Valor Acréscimo",
    "Valor Decréscimo",
    "Valor Saldo Real Final",
    "Valor Total Tributos Reter",
    "Valor Total Tributos Reter Emissão",
    "Valor Total Tributos Reter Baixa",
    "Valor Líquido Título Real",
    "Valor Abatimentos",
    "Valor Tributos Retidos Baixa",
    "Valor Tributos Pendentes Retenção Baixa",
    "Valor Líquido Pago Real",
    "Valor Líquido Saldo Real"
FROM FATO_CONTA_PAGAR_RATEADO
WHERE (:data_vencimento_inicio IS NULL OR "Data Vencimento Real" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR "Data Vencimento Real" <= :data_vencimento_fim)
  AND (:data_emissao_inicio IS NULL OR "Data Emissão" >= :data_emissao_inicio)
  AND (:data_emissao_fim IS NULL OR "Data Emissão" <= :data_emissao_fim)
  AND (:data_baixa_inicio IS NULL OR "Data Última Baixa" >= :data_baixa_inicio)
  AND (:data_baixa_fim IS NULL OR "Data Última Baixa" <= :data_baixa_fim)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_tipo_titulo IS NULL OR "ID Tipo Título" = :id_tipo_titulo)
  AND (:chave_titulo IS NULL OR "Chave Título" = :chave_titulo)
  AND (:numero IS NULL OR "Número" = :numero)
  AND (:somente_em_aberto IS NULL OR :somente_em_aberto <> 'S' OR "Valor Saldo Real" > 0)
  AND (:confissao_divida IS NULL OR "Confissão Dívida" = :confissao_divida)
  AND (:protesto IS NULL OR "Protesto" = :protesto)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Vencimento Real"
FETCH FIRST :limite ROWS ONLY;
