-- DIM_BANCO_CONTA_BANCARIA
-- Dimensão de bancos e contas bancárias, com agência, número, tipo de conta e a empresa/filial à qual cada conta pertence.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_empresa (inteiro)
--   :id_filial (inteiro)
--   :id_banco (inteiro)
--   :codigo_banco (texto)
--   :nome_banco (texto)
--   :id_conta_bancaria (inteiro)
--   :nome_conta (texto)
--   :tipo_conta_bancaria (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Banco",
    "Código Banco",
    "Nome Banco",
    "Banco",
    "ID Conta Bancária",
    "Código Agência",
    "Número Conta Bancária",
    "Nome Conta Bancária",
    "Nome Reduzido Conta Bancária",
    "Tipo Conta Bancária",
    "Conta Bancária",
    "Banco, Agência e Conta Bancária",
    "ID Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Registro Ativo Banco",
    "Reg. Bloqueado Conta Bancária",
    "Registro Ativo Conta Bancária"
FROM DIM_BANCO_CONTA_BANCARIA
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_banco IS NULL OR "ID Banco" = :id_banco)
  AND (:codigo_banco IS NULL OR "Código Banco" = :codigo_banco)
  AND (:nome_banco IS NULL OR UPPER("Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:nome_conta IS NULL OR UPPER("Nome Conta Bancária") LIKE '%' || UPPER(:nome_conta) || '%')
  AND (:tipo_conta_bancaria IS NULL OR "Tipo Conta Bancária" = :tipo_conta_bancaria)
  AND (:registro_ativo IS NULL OR "Registro Ativo Conta Bancária" = :registro_ativo)
ORDER BY "Nome Banco", "Nome Conta Bancária"
FETCH FIRST :limite ROWS ONLY;
