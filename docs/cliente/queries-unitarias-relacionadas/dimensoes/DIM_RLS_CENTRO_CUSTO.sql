-- DIM_RLS_CENTRO_CUSTO
-- Tabela de segurança em nível de linha (RLS): relaciona o e-mail do usuário aos contratos e centros de custo que ele pode visualizar.
--
-- Parâmetros (nulo = filtro ignorado):
--   :email (texto, obrigatório)
--   :id_contrato (inteiro)
--   :id_centro_custo (inteiro)
--   :limite (inteiro)

SELECT
    "E-mail",
    "ID Contrato",
    "Código Contrato",
    "ID Centro de Custo",
    "Código Centro de Custo"
FROM DIM_RLS_CENTRO_CUSTO
WHERE (:email IS NULL OR "E-mail" = :email)
  AND (:id_contrato IS NULL OR "ID Contrato" = :id_contrato)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
ORDER BY "Código Contrato", "Código Centro de Custo"
FETCH FIRST :limite ROWS ONLY;
