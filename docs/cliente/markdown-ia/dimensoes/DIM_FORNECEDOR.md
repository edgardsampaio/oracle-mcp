# DIM_FORNECEDOR

**Tipo:** Dimensão

## Descrição

Dimensão de fornecedores com dados cadastrais (razão social, nome fantasia, CNPJ/CPF, inscrições, endereço, município, UF, região e país).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_fornecedor` | inteiro | Não | `ID Fornecedor` | igual | ID do fornecedor. |
| `codigo_fornecedor` | texto | Não | `Código Fornecedor` | igual | Código do fornecedor. |
| `loja_fornecedor` | texto | Não | `Loja Fornecedor` | igual | Loja do fornecedor. |
| `razao_social` | texto | Não | `Razão Social` | contém | Busca parcial pela razão social. |
| `nome_fantasia` | texto | Não | `Nome Fantasia` | contém | Busca parcial pelo nome fantasia. |
| `cnpj` | texto | Não | `CNPJ` | igual | CNPJ/CPF completo. |
| `raiz_cnpj` | texto | Não | `Raíz CNPJ` | igual | Raiz do CNPJ. |
| `sigla_uf` | texto | Não | `Sigla UF` | igual | UF. |
| `nome_municipio` | texto | Não | `Nome Município` | contém | Busca parcial pelo município. |
| `sigla_pais` | texto | Não | `Sigla País` | igual | País (fornecedores estrangeiros). |
| `codigo_tipo_pessoa` | texto | Não | `Código Tipo Pessoa` | igual | Tipo de pessoa. |
| `registro_ativo` | texto | Não | `Registro Ativo Fornecedor` | igual | Flag de fornecedor ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Fornecedor",
    "Código Fornecedor",
    "Loja Fornecedor",
    "Razão Social",
    "Nome Fantasia",
    "Fornecedor Razão Social",
    "Fornecedor Nome Fantasia",
    "Fornecedor CNPJ Razão Social",
    "Fornecedor CNPJ Nome Fantasia",
    "Raíz CNPJ",
    "CNPJ",
    "Inscrição Estadual",
    "Inscrição Municipal",
    "Endereço",
    "Bairro",
    "CEP",
    "Código Tipo Pessoa",
    "Descrição Tipo Pessoa",
    "Tipo Pessoa",
    "ID Município",
    "Código Município",
    "Nome Município",
    "Município",
    "ID UF",
    "Sigla UF",
    "Nome UF",
    "UF",
    "Região",
    "Sigla País",
    "ID País",
    "Código País",
    "Nome País",
    "País",
    "Registro Bloqueado Fornecedor",
    "Registro Ativo Fornecedor",
    "Registro Bloqueado Pessoa",
    "Registro Ativo Pessoa",
    "Registro Ativo Município",
    "Registro Ativo UF",
    "Registro Ativo País"
FROM DIM_FORNECEDOR
WHERE (:id_fornecedor IS NULL OR "ID Fornecedor" = :id_fornecedor)
  AND (:codigo_fornecedor IS NULL OR "Código Fornecedor" = :codigo_fornecedor)
  AND (:loja_fornecedor IS NULL OR "Loja Fornecedor" = :loja_fornecedor)
  AND (:razao_social IS NULL OR UPPER("Razão Social") LIKE '%' || UPPER(:razao_social) || '%')
  AND (:nome_fantasia IS NULL OR UPPER("Nome Fantasia") LIKE '%' || UPPER(:nome_fantasia) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:sigla_uf IS NULL OR "Sigla UF" = :sigla_uf)
  AND (:nome_municipio IS NULL OR UPPER("Nome Município") LIKE '%' || UPPER(:nome_municipio) || '%')
  AND (:sigla_pais IS NULL OR "Sigla País" = :sigla_pais)
  AND (:codigo_tipo_pessoa IS NULL OR "Código Tipo Pessoa" = :codigo_tipo_pessoa)
  AND (:registro_ativo IS NULL OR "Registro Ativo Fornecedor" = :registro_ativo)
ORDER BY "Razão Social"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Localizar fornecedores por nome, documento ou localização e obter o ID usado em Contas a Pagar, processos, confissões de dívida e movimentações bancárias.
