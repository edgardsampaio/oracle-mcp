# FATO_HISTORICO_COLABORADOR

**Tipo:** Fato

## Descrição

Fato mensal (snapshot por competência) da posição de cada colaborador: filial, seção, centro de custo, função, situação, forma de recebimento, banco e salário base.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `ano_mes_inicio` | texto | Não | `Ano e Mês` | >= | Competência inicial (mesmo formato da coluna `Ano e Mês`). |
| `ano_mes_fim` | texto | Não | `Ano e Mês` | <= | Competência final. |
| `id_colaborador` | inteiro | Não | `ID Colaborador` | igual | Colaborador. |
| `id_filial` | inteiro | Não | `ID Filial` | igual | Filial. |
| `id_secao` | inteiro | Não | `ID Seção` | igual | Seção. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | Centro de custo. |
| `id_funcao_colaborador` | inteiro | Não | `ID Função Colaborador` | igual | Função. |
| `id_situacao_colaborador` | inteiro | Não | `ID Situação Colaborador` | igual | Situação. |
| `id_tipo_recebimento` | inteiro | Não | `ID Tipo Recebimento` | igual | Forma de recebimento. |
| `registro_ativo` | texto | Não | `Registro Ativo Colaborador` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "Ano e Mês",
    "Data Competência",
    "ID Colaborador",
    "ID Filial",
    "ID Seção",
    "ID Centro de Custo",
    "ID Função Colaborador",
    "ID Situação Colaborador",
    "ID Tipo Recebimento",
    "ID Banco",
    "Valor Salário Base",
    "Registro Ativo Colaborador"
FROM FATO_HISTORICO_COLABORADOR
WHERE (:ano_mes_inicio IS NULL OR "Ano e Mês" >= :ano_mes_inicio)
  AND (:ano_mes_fim IS NULL OR "Ano e Mês" <= :ano_mes_fim)
  AND (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_secao IS NULL OR "ID Seção" = :id_secao)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_funcao_colaborador IS NULL OR "ID Função Colaborador" = :id_funcao_colaborador)
  AND (:id_situacao_colaborador IS NULL OR "ID Situação Colaborador" = :id_situacao_colaborador)
  AND (:id_tipo_recebimento IS NULL OR "ID Tipo Recebimento" = :id_tipo_recebimento)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Data Competência" DESC, "ID Colaborador"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de headcount e evolução de quadro: quantos colaboradores ativos por mês, por filial, seção, função ou centro de custo, e evolução do salário base. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo `ID Centro de Custo`.
