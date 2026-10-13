# FATO_ORCAMENTO_TOTAL — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_ORCAMENTO_TOTAL`

## Descrição

Orçamentos por versão e item enriquecidos com a conta orçamentária (N1–N3), a hierarquia do contrato e o calendário do início do período.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo `ID Contrato`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `corc` | DIM_CONTA_ORCAMENTARIA | `ID Conta Orçamentária` | `ID Conta Orçamentária` | — |
| `ctr` | DIM_CENTRO_CUSTO (distinta por contrato) | `ID Contrato` | `ID Contrato` | contrato (nível contrato da DIM_CENTRO_CUSTO, sem centro de custo) |
| `t` | DIM_TEMPO | `Data Início Período` | `Data` | calendário da `Data Início Período` |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `id_orcamento` | inteiro | Não | `f."ID Orçamento"` | igual | Orçamento. |
| `codigo_orcamento` | texto | Não | `f."Código Orçamento"` | igual | Código do orçamento. |
| `numero_versao` | texto | Não | `f."Número Versão"` | igual | Versão do orçamento. |
| `id_conta_orcamentaria` | inteiro | Não | `f."ID Conta Orçamentária"` | igual | Conta orçamentária. |
| `id_contrato` | inteiro | Não | `f."ID Contrato"` | igual | Contrato. |
| `periodo_inicio` | data | Não | — | regra própria | Início do intervalo de interesse (retorna orçamentos cujo período se sobrepõe). |
| `periodo_fim` | data | Não | — | regra própria | Fim do intervalo de interesse. |
| `nome_conta_orcamentaria` | texto | Não | `corc."Nome Conta Orçamentária"` | contém | Busca parcial pela conta orçamentária. |
| `nome_conta_n1` | texto | Não | `corc."Nome Conta Orçamentária N1"` | contém | Busca parcial pela conta de nível 1. |
| `nome_gestor` | texto | Não | `ctr."Nome Gestor"` | contém | Busca parcial pelo gestor. |
| `nome_contrato` | texto | Não | `ctr."Nome Contrato"` | contém | Busca parcial pelo contrato. |
| `ano` | inteiro | Não | `t."Ano"` | igual | Ano da data de referência. |
| `mes` | inteiro | Não | `t."Mês"` | igual | Mês (1–12) da data de referência. |
| `trimestre` | inteiro | Não | `t."Trimestre"` | igual | Trimestre (1–4) da data de referência. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    f."ID Orçamento",
    f."Código Orçamento",
    f."Descrição Orçamento",
    f."Orçamento",
    f."Data Início Período",
    f."Data Fim Período",
    f."Número Versão",
    f."Número Item",
    f."ID Conta Orçamentária",
    f."ID Contrato",
    f."Valor Total Orçado Original",
    f."Valor Total Orçado Atual",
    f."Valor Variação Orig. x Atual",
    corc."Código Conta Orçamentária",
    corc."Nome Conta Orçamentária",
    corc."Conta Orçamentária",
    corc."Conta Orçamentária N1",
    corc."Conta Orçamentária N2",
    corc."Conta Orçamentária N3",
    ctr."Nome Gestor",
    ctr."Gestor",
    ctr."Nome Segmento",
    ctr."Segmento",
    ctr."Nome Grupo Contrato",
    ctr."Grupo Contrato",
    ctr."Código Contrato",
    ctr."Nome Contrato",
    ctr."Contrato",
    ctr."Nome Tipo Contrato",
    ctr."Tipo Contrato",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_ORCAMENTO_TOTAL f
LEFT JOIN DIM_CONTA_ORCAMENTARIA corc ON corc."ID Conta Orçamentária" = f."ID Conta Orçamentária"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Início Período"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_orcamento IS NULL OR f."ID Orçamento" = :id_orcamento)
  AND (:codigo_orcamento IS NULL OR f."Código Orçamento" = :codigo_orcamento)
  AND (:numero_versao IS NULL OR f."Número Versão" = :numero_versao)
  AND (:id_conta_orcamentaria IS NULL OR f."ID Conta Orçamentária" = :id_conta_orcamentaria)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:periodo_inicio IS NULL OR f."Data Fim Período" >= :periodo_inicio)
  AND (:periodo_fim IS NULL OR f."Data Início Período" <= :periodo_fim)
  AND (:nome_conta_orcamentaria IS NULL OR UPPER(corc."Nome Conta Orçamentária") LIKE '%' || UPPER(:nome_conta_orcamentaria) || '%')
  AND (:nome_conta_n1 IS NULL OR UPPER(corc."Nome Conta Orçamentária N1") LIKE '%' || UPPER(:nome_conta_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Início Período", f."Número Item"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas de orçamento com rótulos legíveis — orçado original x atual por conta, contrato, gestor ou segmento — já restrito aos contratos autorizados ao usuário.
