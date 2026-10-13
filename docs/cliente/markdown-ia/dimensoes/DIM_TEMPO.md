# DIM_TEMPO

**Tipo:** Dimensão

## Descrição

Dimensão calendário com dia útil, dia da semana, semana, mês, trimestre, semestre, ano e descrições de período.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `data_inicio` | data | Não | `Data` | >= | Data inicial. |
| `data_fim` | data | Não | `Data` | <= | Data final. |
| `ano` | inteiro | Não | `Ano` | igual | Ano. |
| `mes` | inteiro | Não | `Mês` | igual | Mês (1–12). |
| `trimestre` | inteiro | Não | `Trimestre` | igual | Trimestre (1–4). |
| `semestre` | inteiro | Não | `Semestre` | igual | Semestre (1–2). |
| `dia_util` | texto | Não | `Dia Útil` | igual | Filtra apenas dias úteis/não úteis (valor conforme convenção do DW). |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
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
```

## Objetivo

Traduzir expressões temporais ('último trimestre', 'dias úteis de março') em intervalos de datas e fornecer rótulos de período para agrupamentos.
