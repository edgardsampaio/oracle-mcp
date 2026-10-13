# FATO_CONTA_PAGAR_RATEADO

**Tipo:** Fato

## Descrição

Fato de títulos de Contas a Pagar rateados por natureza e centro de custo, com datas (emissão, vencimento, última baixa), valores originais, pagos, saldos, tributos retidos, acréscimos/decréscimos e indicadores de protesto, confissão de dívida, processo e dação.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `data_vencimento_inicio` | data | Não | `Data Vencimento Real` | >= | Vencimento real a partir de. |
| `data_vencimento_fim` | data | Não | `Data Vencimento Real` | <= | Vencimento real até. |
| `data_emissao_inicio` | data | Não | `Data Emissão` | >= | Emissão a partir de. |
| `data_emissao_fim` | data | Não | `Data Emissão` | <= | Emissão até. |
| `data_baixa_inicio` | data | Não | `Data Última Baixa` | >= | Última baixa a partir de. |
| `data_baixa_fim` | data | Não | `Data Última Baixa` | <= | Última baixa até. |
| `id_filial` | inteiro | Não | `ID Filial` | igual | Filial. |
| `id_fornecedor` | inteiro | Não | `ID Fornecedor` | igual | Fornecedor. |
| `id_natureza` | inteiro | Não | `ID Natureza` | igual | Natureza financeira. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | Centro de custo. |
| `id_tipo_titulo` | inteiro | Não | `ID Tipo Título` | igual | Tipo do título. |
| `chave_titulo` | texto | Não | `Chave Título` | igual | Chave única do título. |
| `numero` | texto | Não | `Número` | igual | Número do título. |
| `somente_em_aberto` | texto | Não | — | regra própria | 'S' retorna apenas títulos com saldo em aberto. |
| `confissao_divida` | texto | Não | `Confissão Dívida` | igual | Flag de título em confissão de dívida. |
| `protesto` | texto | Não | `Protesto` | igual | Flag de título protestado. |
| `registro_ativo` | texto | Não | `Registro Ativo` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "Chave Título",
    "ID Filial",
    "Prefixo",
    "Número",
    "Parcela",
    "ID Tipo Título",
    "Código Tipo Título",
    "ID Fornecedor",
    "Código Fornecedor",
    "Código Loja Fornecedor",
    "ID Município",
    "ID País",
    "ID Natureza",
    "ID Moeda",
    "ID Centro de Custo",
    "Data Emissão",
    "Data Vencimento",
    "Data Vencimento Real",
    "Data Última Baixa",
    "Histórico",
    "Taxa Moeda",
    "Valor Título",
    "Valor Título Real",
    "Valor Última Baixa",
    "Valor Saldo",
    "Valor Saldo Real",
    "Valor Pago Real",
    "Rateio",
    "Protesto",
    "Valor Protesto",
    "Confissão Dívida",
    "Processo",
    "ID Processo",
    "Dação",
    "Adiantamento Cheque",
    "Adiantamento Compensado",
    "Registro Ativo",
    "Valor ISS",
    "Valor INSS",
    "Valor IRRF",
    "Valor PIS",
    "Valor COFINS",
    "Valor CSLL",
    "Valor Saldo Acréscimo",
    "Valor Saldo Decréscimo",
    "Chave Título Pai",
    "Valor Multa",
    "Valor Juros",
    "Valor Desconto",
    "Deduz INSS",
    "Valor SEST",
    "Valor Base",
    "Valor Acréscimo",
    "Valor Decréscimo",
    "Valor Saldo Real Final",
    "Valor Total Tributos Reter",
    "Valor Total Tributos Reter Emissão",
    "Valor Total Tributos Reter Baixa",
    "Valor Líquido Título Real",
    "Valor Abatimentos",
    "Valor Tributos Retidos Baixa",
    "Valor Tributos Pendentes Retenção Baixa",
    "Valor Líquido Pago Real",
    "Valor Líquido Saldo Real"
FROM FATO_CONTA_PAGAR_RATEADO
WHERE (:data_vencimento_inicio IS NULL OR "Data Vencimento Real" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR "Data Vencimento Real" <= :data_vencimento_fim)
  AND (:data_emissao_inicio IS NULL OR "Data Emissão" >= :data_emissao_inicio)
  AND (:data_emissao_fim IS NULL OR "Data Emissão" <= :data_emissao_fim)
  AND (:data_baixa_inicio IS NULL OR "Data Última Baixa" >= :data_baixa_inicio)
  AND (:data_baixa_fim IS NULL OR "Data Última Baixa" <= :data_baixa_fim)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_tipo_titulo IS NULL OR "ID Tipo Título" = :id_tipo_titulo)
  AND (:chave_titulo IS NULL OR "Chave Título" = :chave_titulo)
  AND (:numero IS NULL OR "Número" = :numero)
  AND (:somente_em_aberto IS NULL OR :somente_em_aberto <> 'S' OR "Valor Saldo Real" > 0)
  AND (:confissao_divida IS NULL OR "Confissão Dívida" = :confissao_divida)
  AND (:protesto IS NULL OR "Protesto" = :protesto)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Vencimento Real"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de Contas a Pagar: títulos em aberto, vencidos e a vencer, valores pagos por período, endividamento por fornecedor, natureza ou centro de custo e tributos retidos. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo `ID Centro de Custo`.
