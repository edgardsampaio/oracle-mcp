-- DIM_COLABORADOR
-- Dimensão de colaboradores com matrícula, nome, datas de nascimento/admissão/transferência/desligamento e a posição atual (empresa, filial, seção, função e situação). Contém dados pessoais sujeitos à LGPD.
--
-- Parâmetros (nulo = filtro ignorado):
--   :id_colaborador (inteiro)
--   :matricula (texto)
--   :nome_colaborador (texto)
--   :empresa_atual (texto)
--   :filial_atual (texto)
--   :secao_atual (texto)
--   :funcao_atual (texto)
--   :situacao_atual (texto)
--   :data_admissao_inicio (data)
--   :data_admissao_fim (data)
--   :data_desligamento_inicio (data)
--   :data_desligamento_fim (data)
--   :registro_ativo (texto)
--   :limite (inteiro)

SELECT
    "ID Colaborador",
    "Matrícula Colaborador",
    "Nome Colaborador",
    "Data Nascimento",
    "Data Admissão",
    "Data Transferência",
    "Data Desligamento",
    "Empresa Atual",
    "Filial Atual",
    "Seção Atual",
    "Função Atual",
    "Situação Atual",
    "Registro Ativo Colaborador"
FROM DIM_COLABORADOR
WHERE (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:matricula IS NULL OR "Matrícula Colaborador" = :matricula)
  AND (:nome_colaborador IS NULL OR UPPER("Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:empresa_atual IS NULL OR "Empresa Atual" = :empresa_atual)
  AND (:filial_atual IS NULL OR "Filial Atual" = :filial_atual)
  AND (:secao_atual IS NULL OR UPPER("Seção Atual") LIKE '%' || UPPER(:secao_atual) || '%')
  AND (:funcao_atual IS NULL OR UPPER("Função Atual") LIKE '%' || UPPER(:funcao_atual) || '%')
  AND (:situacao_atual IS NULL OR "Situação Atual" = :situacao_atual)
  AND (:data_admissao_inicio IS NULL OR "Data Admissão" >= :data_admissao_inicio)
  AND (:data_admissao_fim IS NULL OR "Data Admissão" <= :data_admissao_fim)
  AND (:data_desligamento_inicio IS NULL OR "Data Desligamento" >= :data_desligamento_inicio)
  AND (:data_desligamento_fim IS NULL OR "Data Desligamento" <= :data_desligamento_fim)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Nome Colaborador"
FETCH FIRST :limite ROWS ONLY;
