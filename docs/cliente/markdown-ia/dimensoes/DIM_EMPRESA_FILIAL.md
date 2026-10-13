# DIM_EMPRESA_FILIAL

**Tipo:** Dimensão

## Descrição

Dimensão de empresas e filiais do grupo, com códigos, nomes e CNPJ.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_empresa` | inteiro | Não | `ID Empresa` | igual | ID da empresa. |
| `codigo_empresa` | texto | Não | `Código Empresa` | igual | Código da empresa. |
| `nome_empresa` | texto | Não | `Nome Empresa` | contém | Busca parcial pela empresa. |
| `id_filial` | inteiro | Não | `ID Filial` | igual | ID da filial. |
| `codigo_filial` | texto | Não | `Código Filial` | igual | Código da filial. |
| `nome_filial` | texto | Não | `Nome Filial` | contém | Busca parcial pela filial. |
| `cnpj` | texto | Não | `CNPJ` | igual | CNPJ da filial. |
| `raiz_cnpj` | texto | Não | `Raíz CNPJ` | igual | Raiz do CNPJ da empresa. |
| `registro_ativo` | texto | Não | `Registro Ativo Filial` | igual | Flag de filial ativa. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Empresa",
    "Código Empresa",
    "Nome Empresa",
    "Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Raíz CNPJ",
    "CNPJ",
    "Registro Ativo Empresa",
    "Registro Ativo Filial"
FROM DIM_EMPRESA_FILIAL
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:codigo_empresa IS NULL OR "Código Empresa" = :codigo_empresa)
  AND (:nome_empresa IS NULL OR UPPER("Nome Empresa") LIKE '%' || UPPER(:nome_empresa) || '%')
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:codigo_filial IS NULL OR "Código Filial" = :codigo_filial)
  AND (:nome_filial IS NULL OR UPPER("Nome Filial") LIKE '%' || UPPER(:nome_filial) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:registro_ativo IS NULL OR "Registro Ativo Filial" = :registro_ativo)
ORDER BY "Nome Empresa", "Nome Filial"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Resolver empresas e filiais citadas pelo usuário para os IDs usados em todos os fatos e permitir consolidação por empresa.
