# FATO_MOV_BANCARIA_RATEADA

**Tipo:** Fato

## Descrição

Fato de movimentações bancárias (entradas, saídas e transferências) rateadas por natureza e centro de custo, com conta bancária, beneficiário, histórico e dados do título baixado.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `data_inicio` | data | Não | `Data Mov. Bancária` | >= | Data inicial da movimentação. |
| `data_fim` | data | Não | `Data Mov. Bancária` | <= | Data final da movimentação. |
| `operacao` | texto | Não | `Operação Mov. Bancária` | igual | Operação (entrada/saída). |
| `id_filial` | inteiro | Não | `ID Filial` | igual | Filial. |
| `id_conta_bancaria` | inteiro | Não | `ID Conta Bancária` | igual | Conta bancária. |
| `id_tipo_mov_bancaria` | inteiro | Não | `ID Tipo Mov. Bancária` | igual | Tipo de movimentação. |
| `id_natureza` | inteiro | Não | `ID Natureza` | igual | Natureza financeira. |
| `id_centro_custo` | inteiro | Não | `ID Centro de Custo` | igual | Centro de custo. |
| `id_fornecedor` | inteiro | Não | `ID Fornecedor Título` | igual | Fornecedor do título baixado. |
| `id_cliente` | inteiro | Não | `ID Cliente Título` | igual | Cliente do título baixado. |
| `beneficiario` | texto | Não | `Beneficiário` | contém | Busca parcial pelo beneficiário. |
| `historico` | texto | Não | `Histórico` | contém | Busca parcial no histórico. |
| `registro_ativo` | texto | Não | `Registro Ativo` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "Operação Mov. Bancária",
    "ID Filial",
    "Data Mov. Bancária",
    "ID Tipo Mov. Bancária",
    "Código Numerário",
    "Número Documento",
    "Situação Documento",
    "Histórico",
    "Beneficiário",
    "ID Moeda",
    "Valor Mov. Bancária",
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
    "Número Registro",
    "Registro Ativo"
FROM FATO_MOV_BANCARIA_RATEADA
WHERE (:data_inicio IS NULL OR "Data Mov. Bancária" >= :data_inicio)
  AND (:data_fim IS NULL OR "Data Mov. Bancária" <= :data_fim)
  AND (:operacao IS NULL OR "Operação Mov. Bancária" = :operacao)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR "ID Tipo Mov. Bancária" = :id_tipo_mov_bancaria)
  AND (:id_natureza IS NULL OR "ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR "ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR "ID Fornecedor Título" = :id_fornecedor)
  AND (:id_cliente IS NULL OR "ID Cliente Título" = :id_cliente)
  AND (:beneficiario IS NULL OR UPPER("Beneficiário") LIKE '%' || UPPER(:beneficiario) || '%')
  AND (:historico IS NULL OR UPPER("Histórico") LIKE '%' || UPPER(:historico) || '%')
  AND (:registro_ativo IS NULL OR "Registro Ativo" = :registro_ativo)
ORDER BY "Data Mov. Bancária" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de fluxo de caixa realizado: entradas e saídas por período, conta bancária, natureza, centro de custo, fornecedor ou cliente, além de transferências entre contas e filiais. Deve ser combinada com a DIM_RLS_CENTRO_CUSTO pelo `ID Centro de Custo`.
