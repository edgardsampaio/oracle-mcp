# FATO_BAIXA_TITULO_IMPOSTO — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_BAIXA_TITULO_IMPOSTO`

## Descrição

Baixas (pagamentos) de títulos de impostos enriquecidas com filial, tipo de movimentação, moeda, conta bancária, natureza, centro de custo/contrato, tipo de título, órgão/fornecedor do título, filiais e contas de origem/destino e calendário.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo `ID Centro de Custo`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `fil` | DIM_EMPRESA_FILIAL | `ID Filial` | `ID Filial` | — |
| `tmov` | DIM_TIPO_MOV_BANCARIA | `ID Tipo Mov. Bancária Bx. Tit. Imposto` | `ID Tipo Mov. Bancária` | — |
| `moe` | DIM_MOEDA | `ID Moeda` | `ID Moeda` | — |
| `cb` | DIM_BANCO_CONTA_BANCARIA | `ID Conta Bancária` | `ID Conta Bancária` | — |
| `nat` | DIM_NATUREZA | `ID Natureza` | `ID Natureza` | — |
| `ccu` | DIM_CENTRO_CUSTO | `ID Centro de Custo` | `ID Centro de Custo` | — |
| `filo` | DIM_EMPRESA_FILIAL | `ID Empresa Filial Origem` | `ID Filial` | filial de origem (transferência) |
| `fild` | DIM_EMPRESA_FILIAL | `ID Empresa Filial Destino` | `ID Filial` | filial de destino (transferência) |
| `cbo` | DIM_BANCO_CONTA_BANCARIA | `ID Conta Bancária Origem` | `ID Conta Bancária` | conta de origem (transferência) |
| `cbd` | DIM_BANCO_CONTA_BANCARIA | `ID Conta Bancária Destino` | `ID Conta Bancária` | conta de destino (transferência) |
| `ttit` | DIM_TIPO_TITULO | `ID Tipo Título` | `ID Tipo Título` | — |
| `forn` | DIM_FORNECEDOR | `ID Fornecedor Título` | `ID Fornecedor` | fornecedor do título baixado |
| `cli` | DIM_CLIENTE | `ID Cliente Título` | `ID Cliente` | cliente do título baixado |
| `t` | DIM_TEMPO | `Data Mov. Bancária Bx. Tit. Imposto` | `Data` | calendário da `Data Mov. Bancária Bx. Tit. Imposto` |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `data_inicio` | data | Não | `f."Data Mov. Bancária Bx. Tit. Imposto"` | >= | Data inicial da movimentação. |
| `data_fim` | data | Não | `f."Data Mov. Bancária Bx. Tit. Imposto"` | <= | Data final da movimentação. |
| `id_filial` | inteiro | Não | `f."ID Filial"` | igual | Filial. |
| `id_conta_bancaria` | inteiro | Não | `f."ID Conta Bancária"` | igual | Conta bancária. |
| `id_tipo_mov_bancaria` | inteiro | Não | `f."ID Tipo Mov. Bancária Bx. Tit. Imposto"` | igual | Tipo de movimentação. |
| `operacao` | texto | Não | `f."Operação Mov. Bancária Bx. Tit. Imposto"` | igual | Operação (pagar/receber). |
| `id_natureza` | inteiro | Não | `f."ID Natureza"` | igual | Natureza financeira. |
| `id_centro_custo` | inteiro | Não | `f."ID Centro de Custo"` | igual | Centro de custo. |
| `id_fornecedor` | inteiro | Não | `f."ID Fornecedor Título"` | igual | Fornecedor/órgão do título. |
| `id_tipo_titulo` | inteiro | Não | `f."ID Tipo Título"` | igual | Tipo do título. |
| `numero_titulo` | texto | Não | `f."Número Título"` | igual | Número do título. |
| `data_vencimento_inicio` | data | Não | `f."Data Vencimento Título"` | >= | Vencimento do título a partir de. |
| `data_vencimento_fim` | data | Não | `f."Data Vencimento Título"` | <= | Vencimento do título até. |
| `registro_ativo` | texto | Não | `f."Registro Ativo"` | igual | Flag de registro ativo. |
| `nome_empresa` | texto | Não | `fil."Nome Empresa"` | contém | Busca parcial pela empresa. |
| `nome_filial` | texto | Não | `fil."Nome Filial"` | contém | Busca parcial pela filial. |
| `nome_natureza` | texto | Não | `nat."Nome Natureza"` | contém | Busca parcial pela natureza analítica. |
| `nome_natureza_n1` | texto | Não | `nat."Nome Natureza N1"` | contém | Busca parcial pela natureza de nível 1. |
| `nome_gestor` | texto | Não | `ccu."Nome Gestor"` | contém | Busca parcial pelo gestor. |
| `nome_contrato` | texto | Não | `ccu."Nome Contrato"` | contém | Busca parcial pelo contrato. |
| `nome_centro_custo` | texto | Não | `ccu."Nome Centro de Custo"` | contém | Busca parcial pelo centro de custo. |
| `razao_social_fornecedor` | texto | Não | `forn."Razão Social"` | contém | Busca parcial pela razão social do fornecedor. |
| `cnpj_fornecedor` | texto | Não | `forn."CNPJ"` | igual | CNPJ/CPF do fornecedor. |
| `nome_banco` | texto | Não | `cb."Nome Banco"` | contém | Busca parcial pelo banco. |
| `ano` | inteiro | Não | `t."Ano"` | igual | Ano da data de referência. |
| `mes` | inteiro | Não | `t."Mês"` | igual | Mês (1–12) da data de referência. |
| `trimestre` | inteiro | Não | `t."Trimestre"` | igual | Trimestre (1–4) da data de referência. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    f."Operação Mov. Bancária Bx. Tit. Imposto",
    f."ID Filial",
    f."Data Mov. Bancária Bx. Tit. Imposto",
    f."ID Tipo Mov. Bancária Bx. Tit. Imposto",
    f."Código Numerário",
    f."Número Documento",
    f."Situação Documento",
    f."Histórico",
    f."Beneficiário",
    f."ID Moeda",
    f."Valor Mov. Bancária Bx. Tit. Imposto",
    f."ID Conta Bancária",
    f."Banco Virtual",
    f."ID Natureza",
    f."ID Centro de Custo",
    f."Rateio Financeiro",
    f."Rateio Contábil",
    f."ID Empresa Filial Origem",
    f."ID Empresa Filial Destino",
    f."ID Conta Bancária Origem",
    f."ID Conta Bancária Destino",
    f."Prefixo Título",
    f."Número Título",
    f."Parcela Título",
    f."ID Tipo Título",
    f."ID Fornecedor Título",
    f."ID Município Fornecedor Título",
    f."ID País Fornecedor Título",
    f."ID Cliente Título",
    f."ID Município Cliente Título",
    f."Data Emissão Título",
    f."Data Vencimento Título",
    f."Número Lote",
    f."Motivo Baixa",
    f."Sequência Baixa",
    f."Número Cheque",
    f."ID Processo",
    f."Valor Título",
    f."Registro Ativo",
    fil."Código Empresa",
    fil."Nome Empresa",
    fil."Empresa",
    fil."Código Filial",
    fil."Nome Filial",
    fil."Filial",
    fil."CNPJ" AS "CNPJ Filial",
    tmov."Descrição Tipo Mov. Bancária",
    tmov."Tipo Mov. Bancária",
    moe."Moeda",
    cb."Banco",
    cb."Código Agência",
    cb."Número Conta Bancária",
    cb."Nome Conta Bancária",
    cb."Tipo Conta Bancária",
    cb."Conta Bancária",
    cb."Banco, Agência e Conta Bancária",
    nat."Natureza N1",
    nat."Natureza N2",
    nat."Natureza N3",
    nat."Natureza N4",
    nat."Código Natureza",
    nat."Nome Natureza",
    nat."Natureza",
    ccu."Nome Gestor",
    ccu."Gestor",
    ccu."Nome Segmento",
    ccu."Segmento",
    ccu."Nome Grupo Contrato",
    ccu."Grupo Contrato",
    ccu."Código Contrato",
    ccu."Nome Contrato",
    ccu."Contrato",
    ccu."Código Centro de Custo",
    ccu."Nome Centro de Custo",
    ccu."Centro de Custo",
    ccu."Nome Cliente Final",
    ccu."Cliente Final",
    ccu."Nome Tipo Contrato",
    ccu."Tipo Contrato",
    filo."Empresa" AS "Empresa Origem",
    filo."Filial" AS "Filial Origem",
    fild."Empresa" AS "Empresa Destino",
    fild."Filial" AS "Filial Destino",
    cbo."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Origem",
    cbd."Banco, Agência e Conta Bancária" AS "Banco, Agência e Conta Bancária Destino",
    ttit."Descrição Tipo Título",
    ttit."Tipo Título",
    forn."Fornecedor Razão Social" AS "Fornecedor Razão Social Título",
    forn."Fornecedor Nome Fantasia" AS "Fornecedor Nome Fantasia Título",
    forn."Fornecedor CNPJ Razão Social" AS "Fornecedor CNPJ Razão Social Título",
    forn."CNPJ" AS "CNPJ Fornecedor Título",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor Título",
    forn."Município" AS "Município Fornecedor Título",
    forn."Sigla UF" AS "UF Fornecedor Título",
    forn."País" AS "País Fornecedor Título",
    cli."Cliente Razão Social" AS "Cliente Razão Social Título",
    cli."Cliente Nome Fantasia" AS "Cliente Nome Fantasia Título",
    cli."Cliente CNPJ Razão Social" AS "Cliente CNPJ Razão Social Título",
    cli."CNPJ" AS "CNPJ Cliente Título",
    cli."Tipo Pessoa" AS "Tipo Pessoa Cliente Título",
    cli."Município" AS "Município Cliente Título",
    cli."Sigla UF" AS "UF Cliente Título",
    cli."Região" AS "Região Cliente Título",
    t."Ano",
    t."Mês",
    t."Nome do Mês",
    t."Mês e Ano",
    t."Trimestre",
    t."Semestre",
    t."Descrição Curta do Período",
    t."Dia Útil"
FROM FATO_BAIXA_TITULO_IMPOSTO f
LEFT JOIN DIM_EMPRESA_FILIAL fil ON fil."ID Filial" = f."ID Filial"
LEFT JOIN DIM_TIPO_MOV_BANCARIA tmov ON tmov."ID Tipo Mov. Bancária" = f."ID Tipo Mov. Bancária Bx. Tit. Imposto"
LEFT JOIN DIM_MOEDA moe ON moe."ID Moeda" = f."ID Moeda"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cb ON cb."ID Conta Bancária" = f."ID Conta Bancária"
LEFT JOIN DIM_NATUREZA nat ON nat."ID Natureza" = f."ID Natureza"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_EMPRESA_FILIAL filo ON filo."ID Filial" = f."ID Empresa Filial Origem"
LEFT JOIN DIM_EMPRESA_FILIAL fild ON fild."ID Filial" = f."ID Empresa Filial Destino"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbo ON cbo."ID Conta Bancária" = f."ID Conta Bancária Origem"
LEFT JOIN DIM_BANCO_CONTA_BANCARIA cbd ON cbd."ID Conta Bancária" = f."ID Conta Bancária Destino"
LEFT JOIN DIM_TIPO_TITULO ttit ON ttit."ID Tipo Título" = f."ID Tipo Título"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor Título"
LEFT JOIN DIM_CLIENTE cli ON cli."ID Cliente" = f."ID Cliente Título"
LEFT JOIN DIM_TEMPO t ON t."Data" = f."Data Mov. Bancária Bx. Tit. Imposto"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:data_inicio IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" >= :data_inicio)
  AND (:data_fim IS NULL OR f."Data Mov. Bancária Bx. Tit. Imposto" <= :data_fim)
  AND (:id_filial IS NULL OR f."ID Filial" = :id_filial)
  AND (:id_conta_bancaria IS NULL OR f."ID Conta Bancária" = :id_conta_bancaria)
  AND (:id_tipo_mov_bancaria IS NULL OR f."ID Tipo Mov. Bancária Bx. Tit. Imposto" = :id_tipo_mov_bancaria)
  AND (:operacao IS NULL OR f."Operação Mov. Bancária Bx. Tit. Imposto" = :operacao)
  AND (:id_natureza IS NULL OR f."ID Natureza" = :id_natureza)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor Título" = :id_fornecedor)
  AND (:id_tipo_titulo IS NULL OR f."ID Tipo Título" = :id_tipo_titulo)
  AND (:numero_titulo IS NULL OR f."Número Título" = :numero_titulo)
  AND (:data_vencimento_inicio IS NULL OR f."Data Vencimento Título" >= :data_vencimento_inicio)
  AND (:data_vencimento_fim IS NULL OR f."Data Vencimento Título" <= :data_vencimento_fim)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:nome_empresa IS NULL OR UPPER(fil."Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:nome_filial IS NULL OR UPPER(fil."Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:nome_natureza IS NULL OR UPPER(nat."Nome Natureza") LIKE '%' || UPPER(:nome_natureza) || '%')
  AND (:nome_natureza_n1 IS NULL OR UPPER(nat."Nome Natureza N1") LIKE '%' || UPPER(:nome_natureza_n1) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
  AND (:nome_banco IS NULL OR UPPER(cb."Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:ano IS NULL OR t."Ano" = :ano)
  AND (:mes IS NULL OR t."Mês" = :mes)
  AND (:trimestre IS NULL OR t."Trimestre" = :trimestre)
ORDER BY f."Data Mov. Bancária Bx. Tit. Imposto" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder em uma única chamada perguntas sobre impostos recolhidos com rótulos legíveis — por exemplo, total de impostos pagos por natureza e filial no trimestre, por contrato ou por banco — já restrito aos centros de custo autorizados ao usuário.
