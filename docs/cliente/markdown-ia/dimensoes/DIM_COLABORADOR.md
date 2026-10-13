# DIM_COLABORADOR

**Tipo:** Dimensão

## Descrição

Dimensão de colaboradores com matrícula, nome, datas de nascimento/admissão/transferência/desligamento e a posição atual (empresa, filial, seção, função e situação). Contém dados pessoais sujeitos à LGPD.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_colaborador` | inteiro | Não | `ID Colaborador` | igual | ID do colaborador. |
| `matricula` | texto | Não | `Matrícula Colaborador` | igual | Matrícula. |
| `nome_colaborador` | texto | Não | `Nome Colaborador` | contém | Busca parcial pelo nome. |
| `empresa_atual` | texto | Não | `Empresa Atual` | igual | Empresa atual. |
| `filial_atual` | texto | Não | `Filial Atual` | igual | Filial atual. |
| `secao_atual` | texto | Não | `Seção Atual` | contém | Busca parcial pela seção atual. |
| `funcao_atual` | texto | Não | `Função Atual` | contém | Busca parcial pela função atual. |
| `situacao_atual` | texto | Não | `Situação Atual` | igual | Situação atual (ativo, afastado, demitido…). |
| `data_admissao_inicio` | data | Não | `Data Admissão` | >= | Admitidos a partir de. |
| `data_admissao_fim` | data | Não | `Data Admissão` | <= | Admitidos até. |
| `data_desligamento_inicio` | data | Não | `Data Desligamento` | >= | Desligados a partir de. |
| `data_desligamento_fim` | data | Não | `Data Desligamento` | <= | Desligados até. |
| `registro_ativo` | texto | Não | `Registro Ativo Colaborador` | igual | Flag de registro ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Colaborador",
    "Matrícula Colaborador",
    "Nome Colaborador",
    "Data Nascimento",
    "Data Admissão",
    "Data Transferência",
    "Data Desligamento",
    "Empresa Atual",
    "Filial Atual",
    "Seção Atual",
    "Função Atual",
    "Situação Atual",
    "Registro Ativo Colaborador"
FROM DIM_COLABORADOR
WHERE (:id_colaborador IS NULL OR "ID Colaborador" = :id_colaborador)
  AND (:matricula IS NULL OR "Matrícula Colaborador" = :matricula)
  AND (:nome_colaborador IS NULL OR UPPER("Nome Colaborador") LIKE '%' || UPPER(:nome_colaborador) || '%')
  AND (:empresa_atual IS NULL OR "Empresa Atual" = :empresa_atual)
  AND (:filial_atual IS NULL OR "Filial Atual" = :filial_atual)
  AND (:secao_atual IS NULL OR UPPER("Seção Atual") LIKE '%' || UPPER(:secao_atual) || '%')
  AND (:funcao_atual IS NULL OR UPPER("Função Atual") LIKE '%' || UPPER(:funcao_atual) || '%')
  AND (:situacao_atual IS NULL OR "Situação Atual" = :situacao_atual)
  AND (:data_admissao_inicio IS NULL OR "Data Admissão" >= :data_admissao_inicio)
  AND (:data_admissao_fim IS NULL OR "Data Admissão" <= :data_admissao_fim)
  AND (:data_desligamento_inicio IS NULL OR "Data Desligamento" >= :data_desligamento_inicio)
  AND (:data_desligamento_fim IS NULL OR "Data Desligamento" <= :data_desligamento_fim)
  AND (:registro_ativo IS NULL OR "Registro Ativo Colaborador" = :registro_ativo)
ORDER BY "Nome Colaborador"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Identificar colaboradores e sua situação atual para análises de folha de pagamento, headcount, admissões e desligamentos.
