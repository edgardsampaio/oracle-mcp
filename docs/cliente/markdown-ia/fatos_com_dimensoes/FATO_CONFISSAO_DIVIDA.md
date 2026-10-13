# FATO_CONFISSAO_DIVIDA — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_CONFISSAO_DIVIDA`

## Descrição

Saldos de confissões de dívida enriquecidos com os dados do fornecedor, a hierarquia do contrato (gestor, segmento, grupo, tipo) e o calendário da data do último pagamento.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo `ID Contrato`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `forn` | DIM_FORNECEDOR | `ID Fornecedor` | `ID Fornecedor` | — |
| `ctr` | DIM_CENTRO_CUSTO (distinta por contrato) | `ID Contrato` | `ID Contrato` | contrato (nível contrato da DIM_CENTRO_CUSTO, sem centro de custo) |
| `t` | DIM_TEMPO | `Data Último Pagamento` | `Data` | calendário da `Data Último Pagamento` |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `id_fornecedor` | inteiro | Não | `f."ID Fornecedor"` | igual | Fornecedor. |
| `id_contrato` | inteiro | Não | `f."ID Contrato"` | igual | Contrato. |
| `valor_saldo_minimo` | decimal | Não | `f."Valor Saldo"` | >= | Saldo mínimo. |
| `data_ultimo_pagamento_inicio` | data | Não | `f."Data Último Pagamento"` | >= | Último pagamento a partir de. |
| `data_ultimo_pagamento_fim` | data | Não | `f."Data Último Pagamento"` | <= | Último pagamento até. |
| `razao_social_fornecedor` | texto | Não | `forn."Razão Social"` | contém | Busca parcial pela razão social do fornecedor. |
| `cnpj_fornecedor` | texto | Não | `forn."CNPJ"` | igual | CNPJ/CPF do fornecedor. |
| `nome_gestor` | texto | Não | `ctr."Nome Gestor"` | contém | Busca parcial pelo gestor. |
| `nome_contrato` | texto | Não | `ctr."Nome Contrato"` | contém | Busca parcial pelo contrato. |
| `ano` | inteiro | Não | `t."Ano"` | igual | Ano da data de referência. |
| `mes` | inteiro | Não | `t."Mês"` | igual | Mês (1–12) da data de referência. |
| `trimestre` | inteiro | Não | `t."Trimestre"` | igual | Trimestre (1–4) da data de referência. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    f."ID Fornecedor",
    f."ID Contrato",
    f."Valor Saldo",
    f."Data Último Pagamento",
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor",
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
FROM FATO_CONFISSAO_DIVIDA f
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
LEFT JOIN (SELECT DISTINCT "ID Contrato", "Nome Gestor", "Gestor", "Nome Segmento", "Segmento", "Nome Grupo Contrato", "Grupo Contrato", "Código Contrato", "Nome Contrato", "Contrato", "Nome Tipo Contrato", "Tipo Contrato" FROM DIM_CENTRO_CUSTO) ctr ON ctr."ID Contrato" = f."ID Contrato"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Último Pagamento"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Contrato" = f."ID Contrato")
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:id_contrato IS NULL OR f."ID Contrato" = :id_contrato)
  AND (:valor_saldo_minimo IS NULL OR f."Valor Saldo" >= :valor_saldo_minimo)
  AND (:data_ultimo_pagamento_inicio IS NULL OR f."Data Último Pagamento" >= :data_ultimo_pagamento_inicio)
  AND (:data_ultimo_pagamento_fim IS NULL OR f."Data Último Pagamento" <= :data_ultimo_pagamento_fim)
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_gestor IS NULL OR UPPER(ctr."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ctr."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Valor Saldo" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder quanto se deve em confissões de dívida por fornecedor, contrato, gestor ou segmento, e há quanto tempo não há pagamento, já restrito aos contratos autorizados ao usuário.
