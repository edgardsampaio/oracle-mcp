-- FATO_PROCESSO
-- Fato de processos (judiciais/administrativos) com autor, tipo, centro de custo, fornecedor e valor total.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_processo (inteiro)
--   :codigo_processo (texto)
--   :nome_autor (texto)
--   :id_tipo_processo (inteiro)
--   :id_centro_custo (inteiro)
--   :id_fornecedor (inteiro)
--   :valor_minimo (decimal)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Processo",
    "Código Processo",
    "Nome Autor",
    "ID Tipo Processo",
    "ID Centro de Custo",
    "ID Fornecedor",
    "Valor Total",
    "Registro Ativo"
FROM FATO_PROCESSO
WHERE (:id_processo IS NULL OR "ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR "Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER("Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR "ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR "Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Valor Total" DESC
FETCH FIRST :limite ROWS ONLY;
