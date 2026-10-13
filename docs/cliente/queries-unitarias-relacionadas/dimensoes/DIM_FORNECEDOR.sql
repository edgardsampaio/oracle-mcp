-- DIM_FORNECEDOR
-- Dimensão de fornecedores com dados cadastrais (razão social, nome fantasia, CNPJ/CPF, inscrições, endereço, município, UF, região e país).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_fornecedor (inteiro)
--   :codigo_fornecedor (texto)
--   :loja_fornecedor (texto)
--   :razao_social (texto)
--   :nome_fantasia (texto)
--   :cnpj (texto)
--   :raiz_cnpj (texto)
--   :sigla_uf (texto)
--   :nome_municipio (texto)
--   :sigla_pais (texto)
--   :codigo_tipo_pessoa (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Fornecedor",
    "Código Fornecedor",
    "Loja Fornecedor",
    "Razão Social",
    "Nome Fantasia",
    "Fornecedor Razão Social",
    "Fornecedor Nome Fantasia",
    "Fornecedor CNPJ Razão Social",
    "Fornecedor CNPJ Nome Fantasia",
    "Raíz CNPJ",
    "CNPJ",
    "Inscrição Estadual",
    "Inscrição Municipal",
    "Endereço",
    "Bairro",
    "CEP",
    "Código Tipo Pessoa",
    "Descrição Tipo Pessoa",
    "Tipo Pessoa",
    "ID Município",
    "Código Município",
    "Nome Município",
    "Município",
    "ID UF",
    "Sigla UF",
    "Nome UF",
    "UF",
    "Região",
    "Sigla País",
    "ID País",
    "Código País",
    "Nome País",
    "País",
    "Registro Bloqueado Fornecedor",
    "Registro Ativo Fornecedor",
    "Registro Bloqueado Pessoa",
    "Registro Ativo Pessoa",
    "Registro Ativo Município",
    "Registro Ativo UF",
    "Registro Ativo País"
FROM DIM_FORNECEDOR
WHERE (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:codigo_fornecedor IS NULL OR "Código Fornecedor" = :codigo_fornecedor)
  AND (:loja_fornecedor IS NULL OR "Loja Fornecedor" = :loja_fornecedor)
  AND (:razao_social IS NULL OR UPPER("Razão Social") LIKE '%' || UPPER(:razao_social) || '%')
  AND (:nome_fantasia IS NULL OR UPPER("Nome Fantasia") LIKE '%' || UPPER(:nome_fantasia) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:sigla_uf IS NULL OR "Sigla UF" = :sigla_uf)
  AND (:nome_municipio IS NULL OR UPPER("Nome Município") LIKE '%' || UPPER(:nome_municipio) || '%')
  AND (:sigla_pais IS NULL OR "Sigla País" = :sigla_pais)
  AND (:codigo_tipo_pessoa IS NULL OR "Código Tipo Pessoa" = :codigo_tipo_pessoa)
  AND (:registro_ativo IS NULL OR "Registro Ativo Fornecedor" = :registro_ativo)
ORDER BY "Razão Social"
FETCH FIRST :limite ROWS ONLY;
