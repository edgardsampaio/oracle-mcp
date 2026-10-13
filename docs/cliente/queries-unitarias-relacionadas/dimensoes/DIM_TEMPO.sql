-- DIM_TEMPO
-- Dimensão calendário com dia útil, dia da semana, semana, mês, trimestre, semestre, ano e descrições de período.
--
-- Parâmetros (nulo = filtro ignorado):
--   :data_inicio (data)
--   :data_fim (data)
--   :ano (inteiro)
--   :mes (inteiro)
--   :trimestre (inteiro)
--   :semestre (inteiro)
--   :dia_util (texto)
--   :limite (inteiro)

SELECT
    "Data",
    "Dia Útil",
    "Dia da Semana",
    "Nome do Dia da Semana",
    "Semana do Mês",
    "Descrição da Semana do Mês",
    "Semana do Ano",
    "Descrição da Semana do Ano",
    "Descrição Curta do Período",
    "Descrição Longa do Período",
    "Ano e Mês",
    "Mês e Ano",
    "Ano",
    "Mês",
    "Nome do Mês",
    "Trimestre",
    "Descrição do Trimestre",
    "Descrição Curta do Trimestre",
    "Descrição Longa do Trimestre",
    "Semestre",
    "Descrição do Semestre",
    "Descrição Curta do Semestre",
    "Descrição Longa do Semestre"
FROM DIM_TEMPO
WHERE (:data_inicio IS NULL OR "Data" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data" <= :data_fim)
  AND (:ano IS NULL OR "Ano" = :ano)
  AND (:mes IS NULL OR "Mês" = :mes)
  AND (:trimestre IS NULL OR "Trimestre" = :trimestre)
  AND (:semestre IS NULL OR "Semestre" = :semestre)
  AND (:dia_util IS NULL OR "Dia Útil" = :dia_util)
ORDER BY "Data"
FETCH FIRST :limite ROWS ONLY;
