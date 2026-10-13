# FATO_BAIXA_TITULO_IMPOSTO

**Tipo:** Fato

## Descrição

Fato de movimentações bancárias referentes à baixa (pagamento) de títulos de impostos, com dados do título, conta bancária, natureza, centro de custo, rateios e valores.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `data_inicio` | data | Não | `Data Mov. Bancária Bx. Tit. Imposto` | >= | Data inicial da movimentação. |
| `data_fim` | data | Não | `Data Mov. Bancária Bx. Tit. Imposto` | <= | Data final da movimentação. |
| `id_filial` | inteiro | Não | `ID Filial` | igual | Filial. |
| `id_conta_bancaria` | inteiro | Não | `ID Conta Bancária` | igual | Conta bancária. |
| `id_tipo_mov_bancaria` | inteiro | Não | `ID Tipo Mov. Bancária Bx. Tit. Imposto` | igual | Tipo de movimentação. |
| `operacao` | texto | Não | `Operação Mov. Bancária Bx. Tit. Imposto` | igual | Operação (pagar/receber). |
| `id_natureza` | inteiro | Não | `ID Natureza` | igual | Natureza financeira. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | Centro de custo. |
| `id_fornecedor` | inteiro | Não | `ID Fornecedor Título` | igual | Fornecedor/órgão do título. |
| `id_tipo_titulo` | inteiro | Não | `ID Tipo Título` | igual | Tipo do título. |
| `numero_titulo` | texto | Não | `Número Título` | igual | Número do título. |
| `data_vencimento_inicio` | data | Não | `Data Vencimento Título` | >= | Vencimento do título a partir de. |
| `data_vencimento_fim` | data | Não | `Data Vencimento Título` | <= | Vencimento do título até. |
| `registro_ativo` | texto | Não | `Registro Ativo` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "Operação Mov. Bancária Bx. Tit. Imposto",
    "ID Filial",
    "Data Mov. Bancária Bx. Tit. Imposto",
    "ID Tipo Mov. Bancária Bx. Tit. Imposto",
    "Código Numerário",
    "Número Documento",
    "Situação Documento",
    "Histórico",
    "Beneficiário",
    "ID Moeda",
    "Valor Mov. Bancária Bx. Tit. Imposto",
    "ID Conta Bancária",
    "Banco Virtual",
    "ID Natureza",
    "ID Centro de Custo",
    "Rateio Financeiro",
    "Rateio Contábil",
    "ID Empresa Filial Origem",
    "ID Empresa Filial Destino",
    "ID Conta Bancária Origem",
    "ID Conta Bancária Destino",
    "Prefixo Título",
    "Número Título",
    "Parcela Título",
    "ID Tipo Título",
    "ID Fornecedor Título",
    "ID Município Fornecedor Título",
    "ID País Fornecedor Título",
    "ID Cliente Título",
    "ID Município Cliente Título",
    "Data Emissão Título",
    "Data Vencimento Título",
    "Número Lote",
    "Motivo Baixa",
    "Sequência Baixa",
    "Número Cheque",
    "ID Processo",
    "Valor Título",
    "Registro Ativo"
FROM FATO_BAIXA_TITULO_IMPOSTO
WHERE (:data_inicio IS NULL OR "Data Mov. Bancária Bx. Tit. Imposto" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data Mov. Bancária Bx. Tit. Imposto" <= :data_fim)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR "ID Tipo Mov. Bancária Bx. Tit. Imposto" = :id_tipo_mov_bancaria)
  AND (:operacao IS NULL OR "Operação Mov. Bancária Bx. Tit. Imposto" = :operacao)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor Título" = :id_fornecedor)
  AND (:id_tipo_titulo IS NULL OR "ID Tipo Título" = :id_tipo_titulo)
  AND (:numero_titulo IS NULL OR "Número Título" = :numero_titulo)
  AND (:data_vencimento_inicio IS NULL OR "Data Vencimento Título" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR "Data Vencimento Título" <= :data_vencimento_fim)
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Mov. Bancária Bx. Tit. Imposto" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas sobre impostos pagos: quanto foi recolhido por período, por filial, por natureza/tipo de imposto, por fornecedor (órgão) ou por conta bancária. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo `ID Centro de Custo`.
