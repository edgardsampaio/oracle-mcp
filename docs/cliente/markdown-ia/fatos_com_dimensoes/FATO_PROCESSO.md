# FATO_PROCESSO — Fato com Dimensões

**Tipo:** Fato com JOIN de dimensões · **Fato base:** `FATO_PROCESSO`

## Descrição

Processos enriquecidos com o tipo de processo, centro de custo/contrato e fornecedor.

Todas as dimensões entram por `LEFT JOIN`, para que nenhuma linha do fato se perca quando faltar cadastro. A segurança por linha é aplicada pelo `ID Centro de Custo`.

**Relacionamentos**

| Alias | Dimensão | Coluna no fato | Coluna na dimensão | Papel |
|---|---|---|---|---|
| `tproc` | DIM_TIPO_PROCESSO | `ID Tipo Processo` | `ID Tipo Processo` | — |
| `ccu` | DIM_CENTRO_CUSTO | `ID Centro de Custo` | `ID Centro de Custo` | — |
| `forn` | DIM_FORNECEDOR | `ID Fornecedor` | `ID Fornecedor` | — |

## Parâmetros

Todos os filtros são opcionais, exceto `email`: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas. Os parâmetros com alias `f.` filtram o fato; os demais filtram atributos das dimensões.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `email` | texto | Sim | `DIM_RLS_CENTRO_CUSTO."E-mail"` | igual | E-mail do usuário autenticado, injetado pelo servidor MCP (nunca pela IA). Aplica a segurança por linha. |
| `id_processo` | inteiro | Não | `f."ID Processo"` | igual | Processo. |
| `codigo_processo` | texto | Não | `f."Código Processo"` | igual | Código/número do processo. |
| `nome_autor` | texto | Não | `f."Nome Autor"` | contém | Busca parcial pelo autor. |
| `id_tipo_processo` | inteiro | Não | `f."ID Tipo Processo"` | igual | Tipo de processo. |
| `id_centro_custo` | inteiro | Não | `f."ID Centro de Custo"` | igual | Centro de custo. |
| `id_fornecedor` | inteiro | Não | `f."ID Fornecedor"` | igual | Fornecedor. |
| `valor_minimo` | decimal | Não | `f."Valor Total"` | >= | Valor total mínimo. |
| `registro_ativo` | texto | Não | `f."Registro Ativo"` | igual | Flag de registro ativo. |
| `descricao_tipo_processo` | texto | Não | `tproc."Descrição Tipo Processo"` | contém | Busca parcial pelo tipo de processo. |
| `nome_gestor` | texto | Não | `ccu."Nome Gestor"` | contém | Busca parcial pelo gestor. |
| `nome_contrato` | texto | Não | `ccu."Nome Contrato"` | contém | Busca parcial pelo contrato. |
| `nome_centro_custo` | texto | Não | `ccu."Nome Centro de Custo"` | contém | Busca parcial pelo centro de custo. |
| `razao_social_fornecedor` | texto | Não | `forn."Razão Social"` | contém | Busca parcial pela razão social do fornecedor. |
| `cnpj_fornecedor` | texto | Não | `forn."CNPJ"` | igual | CNPJ/CPF do fornecedor. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    f."ID Processo",
    f."Código Processo",
    f."Nome Autor",
    f."ID Tipo Processo",
    f."ID Centro de Custo",
    f."ID Fornecedor",
    f."Valor Total",
    f."Registro Ativo",
    tproc."Descrição Tipo Processo",
    tproc."Tipo Processo",
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
    forn."Fornecedor Razão Social",
    forn."Fornecedor Nome Fantasia",
    forn."Fornecedor CNPJ Razão Social",
    forn."CNPJ" AS "CNPJ Fornecedor",
    forn."Tipo Pessoa" AS "Tipo Pessoa Fornecedor",
    forn."Município" AS "Município Fornecedor",
    forn."Sigla UF" AS "UF Fornecedor",
    forn."País" AS "País Fornecedor"
FROM FATO_PROCESSO f
LEFT JOIN DIM_TIPO_PROCESSO tproc ON tproc."ID Tipo Processo" = f."ID Tipo Processo"
LEFT JOIN DIM_CENTRO_CUSTO ccu ON ccu."ID Centro de Custo" = f."ID Centro de Custo"
LEFT JOIN DIM_FORNECEDOR forn ON forn."ID Fornecedor" = f."ID Fornecedor"
WHERE EXISTS (SELECT 1 FROM DIM_RLS_CENTRO_CUSTO rls WHERE rls."E-mail" = :email AND rls."ID Centro de Custo" = f."ID Centro de Custo")
  AND (:id_processo IS NULL OR f."ID Processo" = :id_processo)
  AND (:codigo_processo IS NULL OR f."Código Processo" = :codigo_processo)
  AND (:nome_autor IS NULL OR UPPER(f."Nome Autor") LIKE '%' || UPPER(:nome_autor) || '%')
  AND (:id_tipo_processo IS NULL OR f."ID Tipo Processo" = :id_tipo_processo)
  AND (:id_centro_custo IS NULL OR f."ID Centro de Custo" = :id_centro_custo)
  AND (:id_fornecedor IS NULL OR f."ID Fornecedor" = :id_fornecedor)
  AND (:valor_minimo IS NULL OR f."Valor Total" >= :valor_minimo)
  AND (:registro_ativo IS NULL OR f."Registro Ativo" = :registro_ativo)
  AND (:descricao_tipo_processo IS NULL OR UPPER(tproc."Descrição Tipo Processo") LIKE '%' || UPPER(:descricao_tipo_processo) || '%')
  AND (:nome_gestor IS NULL OR UPPER(ccu."Nome Gestor") LIKE '%' || UPPER(:nome_gestor) || '%')
  AND (:nome_contrato IS NULL OR UPPER(ccu."Nome Contrato") LIKE '%' || UPPER(:nome_contrato) || '%')
  AND (:nome_centro_custo IS NULL OR UPPER(ccu."Nome Centro de Custo") LIKE '%' || UPPER(:nome_centro_custo) || '%')
  AND (:razao_social_fornecedor IS NULL OR UPPER(forn."Razão Social") LIKE '%' || UPPER(:razao_social_fornecedor) || '%')
  AND (:cnpj_fornecedor IS NULL OR forn."CNPJ" = :cnpj_fornecedor)
ORDER BY f."Valor Total" DESC
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Responder perguntas sobre passivo de processos com rótulos legíveis — quantidade e valor por tipo, contrato, gestor ou fornecedor — já restrito aos centros de custo autorizados ao usuário.
