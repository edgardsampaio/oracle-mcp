-- DIM_EMPRESA_FILIAL
-- Dimensão de empresas e filiais do grupo, com códigos, nomes e CNPJ.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_empresa (inteiro)
--   :codigo_empresa (texto)
--   :nome_empresa (texto)
--   :id_filial (inteiro)
--   :codigo_filial (texto)
--   :nome_filial (texto)
--   :cnpj (texto)
--   :raiz_cnpj (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Empresa",
    "Código Empresa",
    "Nome Empresa",
    "Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Raíz CNPJ",
    "CNPJ",
    "Registro Ativo Empresa",
    "Registro Ativo Filial"
FROM DIM_EMPRESA_FILIAL
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:codigo_empresa IS NULL OR "Código Empresa" = :codigo_empresa)
  AND (:nome_empresa IS NULL OR UPPER("Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:codigo_filial IS NULL OR "Código Filial" = :codigo_filial)
  AND (:nome_filial IS NULL OR UPPER("Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:registro_ativo IS NULL OR "Registro Ativo Filial" = :registro_ativo)
ORDER BY "Nome Empresa", "Nome Filial"
FETCH FIRST :limite ROWS ONLY;
