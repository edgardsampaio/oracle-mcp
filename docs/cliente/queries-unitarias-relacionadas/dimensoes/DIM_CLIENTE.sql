-- DIM_CLIENTE
-- Dimensão de clientes com dados cadastrais (razão social, nome fantasia, CNPJ/CPF, inscrições, endereço, município, UF e região).
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_cliente (inteiro)
--   :codigo_cliente (texto)
--   :loja_cliente (texto)
--   :razao_social (texto)
--   :nome_fantasia (texto)
--   :cnpj (texto)
--   :raiz_cnpj (texto)
--   :sigla_uf (texto)
--   :nome_municipio (texto)
--   :regiao (texto)
--   :codigo_tipo_pessoa (texto)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Cliente",
    "Código Cliente",
    "Loja Cliente",
    "Razão Social",
    "Nome Fantasia",
    "Cliente Razão Social",
    "Cliente Nome Fantasia",
    "Cliente CNPJ Razão Social",
    "Cliente CNPJ Nome Fantasia",
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
    "Registro Bloqueado Cliente",
    "Registro Ativo Cliente",
    "Registro Bloqueado Pessoa",
    "Registro Ativo Pessoa",
    "Registro Ativo Município",
    "Registro Ativo UF"
FROM DIM_CLIENTE
WHERE (:id_cliente IS NULL OR "ID Cliente" = :id_cliente)
  AND (:codigo_cliente IS NULL OR "Código Cliente" = :codigo_cliente)
  AND (:loja_cliente IS NULL OR "Loja Cliente" = :loja_cliente)
  AND (:razao_social IS NULL OR UPPER("Razão Social") LIKE '%' || UPPER(:razao_social) || '%')
  AND (:nome_fantasia IS NULL OR UPPER("Nome Fantasia") LIKE '%' || UPPER(:nome_fantasia) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:sigla_uf IS NULL OR "Sigla UF" = :sigla_uf)
  AND (:nome_municipio IS NULL OR UPPER("Nome Município") LIKE '%' || UPPER(:nome_municipio) || '%')
  AND (:regiao IS NULL OR "Região" = :regiao)
  AND (:codigo_tipo_pessoa IS NULL OR "Código Tipo Pessoa" = :codigo_tipo_pessoa)
  AND (:registro_ativo IS NULL OR "Registro Ativo Cliente" = :registro_ativo)
ORDER BY "Razão Social"
FETCH FIRST :limite ROWS ONLY;
